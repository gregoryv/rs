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

func TestSystem_Import(t *testing.T) {
	var (
		sysA    = NewSystem()
		asRootA = Root.UseAudited(sysA, t)
		exportA bytes.Buffer

		sysB    = NewSystem()
		asRootB = Root.Use(sysB)
	)
	// Make it a bit different from default system
	asRootA.Exec("/bin/mkdir /tmp/test")
	// Export
	sysA.Export(&exportA)
	// Import
	t.Log(exportA.String())
	if err := sysB.Import("/", &exportA); err != nil {
		t.Fatal(err)
	}

	var (
		got     = bytes.NewBufferString("\n")
		exp     = bytes.NewBufferString("\n")
		exportB bytes.Buffer
	)
	asRootA.Fexec(got, "/bin/ls -R -l /")
	asRootB.Fexec(exp, "/bin/ls -R -l /")

	sysB.Export(&exportB)
	t.Log(exportB.String())
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
