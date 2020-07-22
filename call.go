package rs

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/gregoryv/fox"
	"github.com/gregoryv/nugo"
)

type Syscall struct {
	*System // todo should probably hide this
	acc     *Account
	auditer fox.Logger // used to audit who executes what
}

// SetGroup
func (me *Syscall) SetGroup(abspath string, gid int) error {
	n, err := me.stat(abspath)
	if err != nil {
		return wrap("SetGroup", err)
	}
	if !me.acc.Owns(n) && me.acc != Root {
		return fmt.Errorf("SetGroup: %v not owner of %s", me.acc.uid, abspath)
	}
	n.SetGID(gid)
	return nil
}

// SetOwner
func (me *Syscall) SetOwner(abspath string, uid int) error {
	n, err := me.stat(abspath)
	if err != nil {
		return wrap("SetOwner", err)
	}
	if !me.acc.Owns(n) && me.acc != Root {
		return fmt.Errorf("SetOwner: %v not owner of %s", me.acc.uid, abspath)
	}
	n.SetUID(uid)
	return nil
}

// SetMode sets the mode of abspath if the caller is the owner or Root.
// Only permissions bits can be set for now.
func (me *Syscall) SetMode(abspath string, mode Mode) error {
	n, err := me.stat(abspath)
	if err != nil {
		return wrap("SetMode", err)
	}
	if !me.acc.Owns(n) && me.acc != Root {
		return fmt.Errorf("SetMode: %v not owner of %s", me.acc.uid, abspath)
	}
	if nugo.NodeMode(mode) > nugo.ModePerm {
		return fmt.Errorf("SetMode: invalid mode")
	}
	n.SetPerm(nugo.NodeMode(mode)) // todo add SetMode
	return nil
}

// RemoveAll
func (me *Syscall) RemoveAll(abspath string) error {
	n, err := me.stat(abspath)
	if err != nil {
		return wrap("RemoveAll", err)
	}
	n.Parent().DelChild(n.Name())
	return nil
}

// Open resource for reading. Underlying source must be string or []byte.
// If resource is open for writing this call blocks.
func (me *Syscall) Open(abspath string) (*Resource, error) {
	n, err := me.stat(abspath)
	if err != nil {
		return nil, fmt.Errorf("Open: %s", err)
	}
	if err := me.acc.permitted(OpRead, n); err != nil {
		return nil, wrap("Open", err)
	}
	r := newResource(n, OpRead)
	r.unlock = n.RUnlock
	src := n.Source()
	switch src := src.(type) {
	case []byte:
		r.buf = bytes.NewBuffer(src)
	default:
		// todo figure out how to read Any source
		return nil, fmt.Errorf("Open: %s(%T) non readable source", abspath, src)
	}
	// Resource must be closed to unlock
	n.RLock()
	return r, nil
}

// Create returns a new resource for writing. Fails if existing
// resource is directory. Caller must close resource.
func (me *Syscall) Create(abspath string) (*Resource, error) {
	rif, _ := me.Stat(abspath)
	if rif != nil && rif.IsDir() == nil {
		return nil, fmt.Errorf("Create: %s is a directory", abspath)
	}
	dir, name := path.Split(abspath)
	parent, err := me.Stat(dir)
	if err != nil {
		return nil, wrap("Create", err)
	}
	if err := me.acc.permitted(OpWrite, parent.node); err != nil {
		return nil, wrap("Create", err)
	}
	n := parent.node.Make(name)
	n.SetPerm(00644)
	n.UnsetMode(nugo.ModeDir)
	n.Lock()
	r := newResource(n, OpWrite)
	r.buf = &bytes.Buffer{}
	r.unlock = n.Unlock
	return r, nil
}

// SaveAs save src to the given abspath. Fails if abspath already exists.
func (me *Syscall) SaveAs(abspath string, src interface{}) error {
	if _, err := me.Stat(abspath); err == nil {
		return fmt.Errorf("SaveAs: %s exists", abspath)
	}
	w, err := me.Create(abspath)
	if err != nil {
		return wrap("SaveAs", err)
	}
	defer w.Close()
	return wrap("SaveAs", gob.NewEncoder(w).Encode(src))
}

// Save save src to the given abspath. Overwrites existing resource.
// If src implements io.WriterTo interface that is used otherwise it's gob encoded.
func (me *Syscall) Save(abspath string, src interface{}) error {
	rif, _ := me.Stat(abspath)
	if rif != nil && rif.IsDir() == nil {
		return fmt.Errorf("Save: %s is directory", abspath)
	}
	w, err := me.Create(abspath)
	if err != nil {
		return wrap("Save", err)
	}
	defer w.Close()
	switch src := src.(type) {
	case io.WriterTo:
		_, err := src.WriteTo(w)
		return err
	default:
		return wrap("Save", gob.NewEncoder(w).Encode(src))
	}
}

// Load loads the resource from abspath. If res implements
// io.ReaderFrom that is used otherwise gob.Decoded.
func (me *Syscall) Load(res interface{}, abspath string) error {
	r, err := me.Open(abspath)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	switch res := res.(type) {
	case io.ReaderFrom:
		_, err := res.ReadFrom(r)
		return err
	default:
		return wrap("Load", gob.NewDecoder(r).Decode(res))
	}
}

// LoadAccount
func (me *Syscall) LoadAccount(acc *Account, name string) error {
	return me.Load(acc, "/etc/accounts/"+name+".account")
}

// Install resource at the absolute path
func (me *Syscall) Install(abspath string, cmd Executable, mode nugo.NodeMode,
) (*ResInfo, error) {
	dir, name := path.Split(abspath)
	parent, err := me.Stat(dir)
	if err != nil {
		return nil, wrap("Install", err)
	}
	if err := me.acc.permitted(OpWrite, parent.node); err != nil {
		return nil, wrap("Install", err)
	}
	n := parent.node.Make(name)
	n.SetPerm(mode)
	n.SetSource(cmd)
	n.UnsetMode(nugo.ModeDir)
	return &ResInfo{node: n}, nil
}

// Fexec creates and executes a new command and directs the output to
// the given writer.
func (me *Syscall) Fexec(w io.Writer, abspath string, args ...string) error {
	cmd := NewCmd(abspath, args...)
	cmd.Out = w
	return me.ExecCmd(cmd)
}

// Exec splits the cli on whitespace and executes the first as
// absolute path and the rest as arguments
func (me *Syscall) Exec(cli string) error {
	parts := strings.Split(cli, " ")
	return me.ExecCmd(NewCmd(parts[0], parts[1:]...))
}

// ExecCmd executes the given command. Fails if e.g. resource is not
// Executable. All exec calls are audited if system has an auditer
// configured.
func (me *Syscall) ExecCmd(cmd *Cmd) error {
	n, err := me.stat(cmd.Abspath)
	if err != nil {
		return err
	}
	switch src := n.Source().(type) {
	case Executable:
		// If needed setuid can be checked and enforced here
		cmd.Sys = me
		err = src.Exec(cmd)
		if me.auditer != nil {
			msg := fmt.Sprintf("%v %s", me.acc.uid, cmd.String())
			if err != nil {
				// don't audit the actual error message, leave that to
				// other form of logging
				msg = fmt.Sprintf("%s ERR", msg)
			}
			me.auditer.Log(msg)
		}
		return err
	default:
		return fmt.Errorf("Cannot run %T", src)
	}
}

type Mode nugo.NodeMode

// AddAccount adds a new account to the system. Name and uid must be
// unique.
func (me *Syscall) AddAccount(acc *Account) error {
	for _, existing := range me.System.accounts {
		if existing.uid == acc.uid {
			return fmt.Errorf("uid exists")
		}
		if existing.name == acc.name {
			return fmt.Errorf("name exists")
		}
	}
	me.System.accounts = append(me.System.accounts, acc)
	abspath := fmt.Sprintf("/etc/accounts/%s.account", acc.name)
	return me.Save(abspath, acc)
}

// Mkdir creates the absolute path whith a given mode where the parent
// must exist.
func (me *Syscall) Mkdir(abspath string, mode Mode) (*ResInfo, error) {
	dir, name := path.Split(abspath)
	parent, err := me.stat(dir)
	if err != nil {
		return nil, fmt.Errorf("Mkdir: %w", err)
	}
	if err := me.acc.permitted(OpWrite, parent); err != nil {
		return nil, fmt.Errorf("Mkdir: %w", err)
	}
	n := parent.Make(name)
	n.SetPerm(nugo.NodeMode(mode))
	return &ResInfo{node: n}, nil
}

// Stat returns the node of the abspath if account is allowed to reach
// it, ie. all nodes up to it must have execute flags set.
func (me *Syscall) Stat(abspath string) (*ResInfo, error) {
	n, err := me.stat(abspath)
	if err != nil {
		return nil, fmt.Errorf("Stat %v", err)
	}
	return &ResInfo{node: n}, nil
}

// stat returns the node of the abspath if account is allowed to reach
// it, ie. all nodes up to it must have execute mode set.
func (me *Syscall) stat(abspath string) (*nugo.Node, error) {
	rn := me.rootNode(abspath)
	n, err := rn.Find(abspath)
	if err != nil {
		return nil, err
	}
	parent := n.Parent()
	// check each parent for access
	for parent != nil {
		if err := me.acc.permitted(OpExec, parent); err != nil {
			return nil, fmt.Errorf("%s uid:%d: %v", abspath, me.acc.uid, err)
		}
		parent = parent.Parent()
	}
	return n, nil
}

func wrap(prefix string, err error) error {
	if err != nil {
		return fmt.Errorf("%s: %w", prefix, err)
	}
	return nil
}

// mount creates a root node for the given path.
func (me *Syscall) mount(abspath string, mode nugo.NodeMode) error {
	rn := nugo.NewRootNode(abspath, mode)
	rn.SetSeal(me.acc.uid, me.acc.gid(), 01755)
	return me.System.mount(rn)
}

// Visitor is called during a walk with a specific node and the
// absolute path to that node. Use the given Walker to stop if needed.
type Visitor func(child *ResInfo, abspath string, w *nugo.Walker)
