package rs

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/fox"
	"github.com/gregoryv/golden"
	"github.com/gregoryv/nugo"
)

func TestSystem_Export(t *testing.T) {
	sys := NewSystem()
	var buf bytes.Buffer
	ok, _ := asserter.NewErrors(t)
	ok(sys.Export(&buf))
	golden.Assert(t, buf.String())
}

func TestImport(t *testing.T) {
	var firstExport bytes.Buffer
	NewSystem().Export(&firstExport)
	// using copy of buffer content for later comparison
	importedSys, err := Import(strings.NewReader(firstExport.String()))
	if err != nil {
		t.Fatal(err)
	}
	var secondExport bytes.Buffer
	importedSys.Export(&secondExport)
	assert := asserter.New(t)
	assert().Equals(firstExport.String(), secondExport.String())
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

func Example_saveAndLoadResource() {
	asRoot := Root.Use(NewSystem())
	asRoot.Exec("/bin/mkdir /tmp/aliens")
	asRoot.Save("/tmp/aliens/green.gob", &Alien{Name: "Mr Green"})
	var alien Alien
	asRoot.Load(&alien, "/tmp/aliens/green.gob")
	fmt.Printf("%#v", alien)
	// output:
	// rs.Alien{Name:"Mr Green"}
}
