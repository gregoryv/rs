/*
Package rs provides a resource system which enforces unix style access control.

Resources are stored as nugo.Nodes and can either have a []byte slice
as source or implement the Executable interface. Using the Save and
Load syscalls, structs are gob encoded and decoded to an access
controlled resource.

Anonymous account has uid,gid 0,0 whereas the Root account 1,1.

*/
package rs

import (
	"bufio"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gregoryv/fox"
	"github.com/gregoryv/nexus"
	"github.com/gregoryv/nugo"
)

// NewSystem returns a system with installed resources resembling a
// unix filesystem.
func NewSystem() *System {
	sys := &System{
		mounts:   make(map[string]*nugo.Node),
		Accounts: []*Account{},
		Groups:   []*Group{},
	}
	asRoot := Root.Use(sys)
	asRoot.mount("/", nugo.ModeDir|nugo.ModeSort|nugo.ModeDistinct)
	installSys(sys)
	return sys
}

// installSys creates default resources on the system. Should only be
// called once on one system.
func installSys(sys *System) {
	asRoot := Root.Use(sys)
	asRoot.Mkdir("/bin", 01755)
	asRoot.Mkdir("/etc", 00755)
	asRoot.Mkdir("/etc/accounts", 00755)
	asRoot.Mkdir("/etc/groups", 00755)
	asRoot.Mkdir("/tmp", 07777)
	asRoot.Install("/bin/chmod", ExecFunc(Chmod), 00755)
	asRoot.Install("/bin/chown", &Chown{}, 00755)
	asRoot.Install("/bin/ls", ExecFunc(Ls), 01755)
	asRoot.Install("/bin/mkacc", ExecFunc(Mkacc), 00755)
	asRoot.Install("/bin/mkdir", ExecFunc(Mkdir), 00755)
	asRoot.Install("/bin/secure", ExecFunc(Secure), 00755)

	asRoot.AddAccount(Anonymous)
	asRoot.AddAccount(Root)
}

type System struct {
	mounts   map[string]*nugo.Node
	Accounts []*Account
	Groups   []*Group

	auditer fox.Logger // Used audit Syscall.Exec calls

	lm           sync.RWMutex
	lastModified time.Time // last time a resource was modified
}

// touch synced update of lastModified field
func (me *System) touch() {
	me.lm.Lock()
	me.lastModified = time.Now()
	me.lm.Unlock()
}

// LastModified returns last time resources state was modified.
func (me *System) LastModified() time.Time {
	me.lm.RLock()
	defer me.lm.RUnlock()
	return me.lastModified
}

// NextUID returns next available uid
func (me *System) NextUID() int {
	var uid int
	for _, acc := range me.Accounts {
		if acc.UID > uid {
			uid = acc.UID
		}
	}
	return uid + 1
}

// accountByUID
func (me *System) accountByUID(uid int) (*Account, error) {
	for _, acc := range me.Accounts {
		if acc.UID == uid {
			return acc, nil
		}
	}
	return nil, fmt.Errorf("uid %v not found", uid)
}

// NextGID returns next available gid
func (me *System) NextGID() int {
	var gid int
	for _, acc := range me.Accounts {
		for _, id := range acc.Groups {
			if id > gid {
				gid = id
			}
		}
	}
	return gid + 1
}

// groupByGID
func (me *System) groupByGID(gid int) (*Group, error) {
	for _, group := range me.Groups {
		if group.gid == gid {
			return group, nil
		}
	}
	return nil, fmt.Errorf("gid %v not found", gid)
}

// SetAuditer sets the auditer for Syscall.Exec calls
func (me *System) SetAuditer(auditer fox.Logger) *System {
	me.auditer = auditer
	return me
}

func (me *System) mount(rn *nugo.Node) error {
	abspath := path.Clean(rn.Name)
	if _, found := me.mounts[abspath]; found {
		return fmt.Errorf("mount: %s already exists", abspath)
	}
	me.mounts[abspath] = rn
	me.touch()
	return nil
}

// rootNode returns the mounting point of the abspath. Currently only
// "/" is available.
func (me *System) rootNode(abspath string) *nugo.Node {
	rn := me.mounts["/"]
	for p, n := range me.mounts {
		if strings.Index(abspath, p) == 0 {
			if len(n.Name) > len(rn.Name) {
				rn = n
			}
		}
	}
	return rn
}

// Export
func (me *System) Export(w io.Writer) error {
	root := me.rootNode("/")
	exp := NodeExporter(w)
	exp(root, "/", nil)
	root.Walk(NodeExporter(w))
	return nil
}

// NodeExporter writes each node with it's src as gob encoded base64 string
func NodeExporter(writer io.Writer) nugo.Visitor {
	p, _ := nexus.NewPrinter(writer)
	return func(node *nugo.Node, abspath string, w *nugo.Walker) {
		p.Print(uint32(node.Mode), node.UID, node.GID, " ", abspath)
		if !node.IsDir() {
			p.Print(" ")
			b64 := base64.NewEncoder(base64.StdEncoding, p)
			switch content := node.Content.(type) {
			case []byte:
				b64.Write(content)
			default:
				// all executables must be structs for this to work
				if content != nil {
					gob.NewEncoder(b64).Encode(content)
				}
			}
		}
		p.Println()
	}
}

func Import(r io.Reader) (*System, error) {
	scanner := bufio.NewScanner(r)
	var rn *nugo.Node
	for scanner.Scan() {
		n := nugo.NewNode("undef")
		var (
			src     string
			abspath string
			modeStr string
		)

		_, err := fmt.Sscanf(scanner.Text(), "%s %d %d %s %s",
			&modeStr,
			&n.UID,
			&n.GID,
			&abspath,
			&src,
		)
		if err != nil && err != io.EOF {
			return nil, err
		}
		mode, err := strconv.ParseUint(modeStr, 10, 32)
		if err != nil {
			return nil, err
		}
		n.Mode = nugo.NodeMode(mode)
		if rn == nil {
			n.Name = abspath
			rn = n
			continue
		}
		n.Name = path.Base(abspath)

		if src != "EOF" {
			b64 := base64.NewDecoder(base64.StdEncoding, strings.NewReader(src))
			content := make([]byte, 1000)
			gob.NewDecoder(b64).Decode(&content)
			n.Content = content
		}
		parent, err := rn.Find(path.Dir(abspath))
		if err != nil {
			return nil, err
		}
		parent.Add(n)
	}
	sys := NewSystem()
	sys.mount(rn)
	sys.touch()
	return sys, nil
}

func isSibling(lastpath, abspath string) bool {
	return lastpath == path.Dir(abspath)
}

func isChild(lastpath, abspath string) bool {
	return lastpath == path.Dir(path.Dir(abspath))
}
