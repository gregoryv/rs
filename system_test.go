package rs

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/fox"
	"github.com/gregoryv/golden"
	"github.com/gregoryv/nugo"
)

func TestSystem_LastModified(t *testing.T) {
	sys := NewSystem()
	before := sys.LastModified()
	var x time.Time
	assert := asserter.New(t)
	assert(x != before).Error(before)
	sys.touch()
	after := sys.LastModified()
	assert(before != after).Error(before, after)
}

func TestSystem_Export(t *testing.T) {
	sys := NewSystem()
	var buf bytes.Buffer
	ok := asserter.Wrap(t).Ok
	ok(sys.Export(&buf))
	golden.Assert(t, buf.String())
}

func Test_data_persistence(t *testing.T) {
	t.Log(`System should be the same after export and import`)
	var (
		sysA    = NewSystem()
		asRootA = Root.Use(sysA)
		exportA bytes.Buffer
	)
	// Make it a bit different from default system
	asRootA.Exec("/bin/mkdir /tmp/test")
	asRootA.Exec("/bin/mkacc -uid 2 -gid 2 john")
	err := asRootA.Exec("/bin/secure -a john -s mysecret")
	if err != nil {
		t.Fatal(err)
	}

	// Export
	sysA.Export(&exportA)
	exp := exportA.String() // save for later as it's read when imported

	// Import
	sysB := NewSystem()
	if err := sysB.Import("/", &exportA); err != nil {
		t.Fatal(err)
	}

	var exportB bytes.Buffer
	sysB.Export(&exportB)
	got := exportB.String()
	equals := asserter.Wrap(t).Equals
	equals(got, exp)
}

func TestSystem_groupByGID(t *testing.T) {
	sys := NewSystem()
	ok, bad := asserter.NewMixed(t)
	ok(sys.groupByGID(0))
	bad(sys.groupByGID(99))
}

func TestSystem_accountByUID(t *testing.T) {
	sys := NewSystem()
	ok, bad := asserter.NewMixed(t)
	ok(sys.accountByUID(0))
	bad(sys.accountByUID(99))
}

func TestSystem_SetAuditer(t *testing.T) {
	var (
		buf     bytes.Buffer
		sys     = NewSystem().SetAuditer(fox.NewSyncLog(&buf))
		asRoot  = Root.Use(sys) // use after the auditer is set
		asJohn  = John.Use(sys)
		ok, bad = asserter.NewErrors(t)
	)
	bad(asJohn.Exec("/bin/mkdir /etc/s"))
	ok(asJohn.Exec("/bin/mkdir /tmp/s"))
	ok(asRoot.Exec("/bin/mkdir /etc/x"))
	if buf.String() == "" {
		t.Error("expected audit")
	}
	if !strings.Contains(buf.String(), "ERR") {
		t.Error(buf.String())
	}
}

func TestSystem_rootNode(t *testing.T) {
	sys := NewSystem()
	sys.mount(nugo.NewRoot("/mnt"))
	sys.mount(nugo.NewRoot("/mnt/usb"))
	if rn := sys.rootNode("/mnt/usb/some/path"); rn.Name != "/mnt/usb" {
		t.Fail()
	}
	if rn := sys.rootNode("/nosuch/dir"); rn.Name != "/" {
		t.Fail()
	}
}

func TestSystem_mount(t *testing.T) {
	sys := NewSystem()
	ok, bad := asserter.NewErrors(t)
	bad(sys.mount(nugo.NewRoot("/")))
	ok(sys.mount(nugo.NewRoot("/mnt/usb")))
}
