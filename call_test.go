package rs

import (
	"fmt"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestSyscall_SetGroup_asRoot(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	ok, bad := asserter.NewErrors(t)
	ok(asRoot.SetGroup("/tmp", 2))
	bad(asRoot.SetGroup("/nosuch", 2))
}

func TestSyscall_SetGroup_asAnonymous(t *testing.T) {
	asAnonymous := Anonymous.Use(NewSystem())
	_, bad := asserter.NewErrors(t)
	bad(asAnonymous.SetGroup("/tmp", 2))
}

func TestSyscall_SetOwner_asRoot(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	ok, bad := asserter.NewErrors(t)
	ok(asRoot.SetOwner("/tmp", 2))
	bad(asRoot.SetOwner("/nosuch", 2))
}

func TestSyscall_SetOwner_asAnonymous(t *testing.T) {
	asAnonymous := Anonymous.Use(NewSystem())
	_, bad := asserter.NewErrors(t)
	bad(asAnonymous.SetOwner("/tmp", 2))
}

func TestSyscall_SetMode_asRoot(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	ok, bad := asserter.NewErrors(t)
	ok(asRoot.SetMode("/tmp", 0))
	bad(asRoot.SetMode("/nosuch", 0))
}

func TestSyscall_SetMode_asJohn(t *testing.T) {
	asJohn := John.Use(NewSystem())
	_, bad := asserter.NewErrors(t)
	bad(asJohn.SetMode("/etc", 0))
}

func TestSyscall_SetMode_asAnonymous(t *testing.T) {
	asAnonymous := Anonymous.Use(NewSystem())
	_, bad := asserter.NewErrors(t)
	bad(asAnonymous.SetMode("/etc", 0))
}

func TestSyscall_RemoveAll_asRoot(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	asRoot.SaveAs("/tmp/alien", &Alien{Name: "RemoveAll"})
	ok, bad := asserter.NewErrors(t)
	ok(asRoot.RemoveAll("/tmp/alien"))
	bad(asRoot.RemoveAll("/tmp/nosuch"))
}

func TestSyscall_RemoveAll_asAnonymous(t *testing.T) {
	asAnonymous := Anonymous.Use(NewSystem())
	_, bad := asserter.NewErrors(t)
	bad(asAnonymous.RemoveAll("/etc/accounts.gob"))
}

func TestSyscall_Load_asRoot(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	ok, _ := asserter.NewFatalErrors(t)
	assert := asserter.New(t)
	alien := Alien{Name: "Mr green"}
	ok(asRoot.Save("/thing.gob", &alien))
	var got Alien
	ok(asRoot.Load(&got, "/thing.gob"))
	assert().Equals(got, alien)
}

func TestSyscall_Load_errors_asRoot(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	_, bad := asserter.NewErrors(t)
	var x interface{}
	bad(asRoot.Load(&x, "/nosuch"))
	bad(asRoot.Load(&x, "/bin/mkdir"))
}

func TestSyscall_Save_asRoot(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	ok, _ := asserter.NewFatalErrors(t)
	ok(asRoot.Save("/thing.gob", &Alien{Name: "Mr green"}))
	ok(asRoot.Save("/thing.gob", &Alien{Name: "Mr red"}))
}

func TestSyscall_Save_errors_asRoot(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	_, bad := asserter.NewErrors(t)
	bad(asRoot.Save("/nosuch/thing.gob", &Alien{}))
	bad(asRoot.Save("/", &Alien{}))
}

func TestSyscall_SaveAs_asRoot(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	ok, bad := asserter.NewErrors(t)
	ok(asRoot.SaveAs("/thing.gob", &Alien{}))
	bad(asRoot.SaveAs("/thing.gob", &Alien{}))
}

func TestSyscall_SaveAs_errors_asRoot(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	_, bad := asserter.NewErrors(t)
	bad(asRoot.SaveAs("/nosuch/thing.gob", &Alien{}))
	bad(asRoot.SaveAs("/", &Alien{}))
}

func TestSyscall_Open_asRoot(t *testing.T) {
	sys := NewSystem()
	asRoot := Root.Use(sys)
	must, _ := asserter.NewFatalErrors(t)
	must(asRoot.Save("/tmp/alien.gob", &Alien{Name: "x"}))
	ok, bad := asserter.NewMixed(t)
	ok(asRoot.Open("/tmp/alien.gob")).Log("owner created")
	bad(asRoot.Open("/nosuch")).Log("missing resource")
	res, _ := asRoot.Open("/tmp/alien.gob")
	bad(res.Write([]byte(""))).Log("write to readonly")
}

func TestSyscall_Open_asAnonymous(t *testing.T) {
	sys := NewSystem()
	asRoot := Root.Use(sys)
	asAnonymous := Anonymous.Use(sys)
	must, _ := asserter.NewFatalErrors(t)
	must(asRoot.Save("/tmp/alien.gob", &Alien{Name: "x"}))
	ok, bad := asserter.NewMixed(t)
	ok(asRoot.Open("/tmp/alien.gob")).Log("owner created")
	bad(asAnonymous.Open("/tmp/alien.gob")).Log("inadequate permission")
}

func TestSyscall_Create_asRoot(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	must, _ := asserter.NewFatalErrors(t)
	res, err := asRoot.Create("/file")
	must(err)
	ok, bad := asserter.NewMixed(t)
	ok(asRoot.Create("/file")).Log("overwrite")
	bad(res.Read([]byte{})).Log("write only")
	bad(asRoot.Create("/")).Log("directory")
}

func TestSyscall_Create_asAnonymous(t *testing.T) {
	sys := NewSystem()
	asRoot := Root.Use(sys)
	asAnonymous := Anonymous.Use(sys)
	must, _ := asserter.NewFatalMixed(t)
	must(asRoot.Create("/file"))
	_, bad := asserter.NewMixed(t)
	bad(asAnonymous.Create("/file"))
}

func TestSyscall_Mkdir_asRoot(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	ok, bad := asserter.NewMixed(t)
	ok(asRoot.Mkdir("/adir", 0))
	bad(asRoot.Mkdir("/nosuch/whatever", 0))
}

func TestSyscall_Mkdir_asAnonymous(t *testing.T) {
	asAnonymous := Anonymous.Use(NewSystem())
	_, bad := asserter.NewMixed(t)
	bad(asAnonymous.Mkdir("/whatever", 0))
}

func TestSyscall_Exec_asRoot(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	ok, bad := asserter.NewErrors(t)
	ok(asRoot.Exec("/bin/mkdir /tmp"))
	bad(asRoot.Exec("/bin/nosuch/mkdir /tmp"))
	bad(asRoot.Exec("/bin")).Log("not executable")
	bad(asRoot.Exec("/bin/mkdir -nosuch")).Log("bad flag")
}

func ExampleSyscall_Stat() {
	sys := Anonymous.Use(NewSystem())
	_, err := sys.Stat("/etc/accounts/root")
	fmt.Println(err)
	// output:
	// Stat /etc/accounts/root uid:0: d---rwxr-xr-x 1 1 exec denied
}

func TestSystem_Stat_asRoot(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	ok, bad := asserter.NewMixed(t)
	ok(asRoot.Stat("/bin"))
	bad(asRoot.Stat("/nothing"))
}

func TestSystem_Install_asRoot(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	ok, bad := asserter.NewMixed(t)
	ok(asRoot.Install("/bin/x", nil, 0))
	bad(asRoot.Install("/bin/nosuchdir/x", nil, 0))
}

func TestSystem_Install_asAnonymous(t *testing.T) {
	asAnonymous := Anonymous.Use(NewSystem())
	ok, bad := asserter.NewMixed(t)
	ok(asAnonymous.Install("/tmp/x", nil, 0))
	bad(asAnonymous.Install("/bin/x", nil, 0))
}

func TestSyscall_AddAccount_asRoot(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	var eva Account
	ok, _ := asserter.NewFatalErrors(t)
	ok(asRoot.AddAccount(NewAccount("eva", 3)))
	ok(asRoot.LoadAccount(&eva, "eva"))
	assert := asserter.New(t)
	assert().Equals(eva.UID, 3)
}

func TestSyscall_AddAccount_withEmail(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	ok, _ := asserter.NewFatalErrors(t)
	ok(asRoot.AddAccount(NewAccount("john@example.com", 3)))
	var acc Account
	ok(asRoot.LoadAccount(&acc, "john@example.com"))
}

func TestSyscall_AddAccount_bad(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	_, bad := asserter.NewFatalErrors(t)
	bad(asRoot.AddAccount(Root))
}

func TestSyscall_AddAccount_asJohn(t *testing.T) {
	asJohn := John.Use(NewSystem())
	_, bad := asserter.NewFatalErrors(t)
	bad(asJohn.AddAccount(John))
}

func TestSyscall_joinGroup_asRoot(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	ok, bad := asserter.NewFatalErrors(t)
	ok(asRoot.joinGroup(&Group{Name: "new", gid: 100}))
	bad(asRoot.joinGroup(&Group{Name: "new", gid: 100}))
	bad(asRoot.joinGroup(&Group{Name: "new", gid: 101}))
}

func TestSyscall_joinGroup_asJohn(t *testing.T) {
	asJohn := John.Use(NewSystem())
	_, bad := asserter.NewFatalErrors(t)
	bad(asJohn.joinGroup(&Group{Name: "new", gid: 100}))
}
