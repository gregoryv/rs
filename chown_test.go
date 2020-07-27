package rs

import (
	"os"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestChown_asRoot(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	ok, bad := asserter.NewFatalErrors(t)
	ok(asRoot.Exec("/bin/mkacc john"))
	ok(asRoot.Exec("/bin/chown john /tmp"))
	ok(asRoot.Exec("/bin/chown john:john /tmp"))
	bad(asRoot.Exec("/bin/chown john:nosuch /tmp"))
	bad(asRoot.Exec("/bin/chown")).Log("missing args")
	bad(asRoot.Exec("/bin/chown john /nosuch")).Log("missing resource")
	bad(asRoot.Exec("/bin/chown clark /tmp")).Log("account missing")
}

func TestChown_asAnonymous(t *testing.T) {
	sys := NewSystem()
	asRoot := Root.Use(sys)
	asAnonymous := Anonymous.Use(sys)
	ok, bad := asserter.NewFatalErrors(t)
	ok(asRoot.Exec("/bin/mkacc john"))
	bad(asAnonymous.Exec("/bin/chown john:john /tmp"))
}

func TestChown_asJohn(t *testing.T) {
	asJohn := John.Use(NewSystem())
	_, bad := asserter.NewFatalErrors(t)
	bad(asJohn.Exec("/bin/chown root:root /tmp"))
}

func ExampleChown_help() {
	asRoot := Root.Use(NewSystem())
	asRoot.Fexec(os.Stdout, "/bin/chown -h")
	// output:
	// Usage: chown OWNER ...paths
}
