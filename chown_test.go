package rs

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func TestChown_asRoot(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	ok, bad := asserter.NewFatalErrors(t)
	ok(asRoot.Exec("/bin/mkacc john"))
	ok(asRoot.Exec("/bin/chown john /tmp"))
	bad(asRoot.Exec("/bin/chown")).Log("missing args")
	bad(asRoot.Exec("/bin/chown john /nosuch")).Log("missing resource")
	bad(asRoot.Exec("/bin/chown clark /tmp")).Log("account missing")
}
