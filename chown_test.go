package rs

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func TestChown(t *testing.T) {
	sys := NewSystem()
	asRoot := Root.Use(sys)
	ok, bad := asserter.NewErrors(t)

	bad(asRoot.Exec("/bin/chown john /tmp"))
	ok(asRoot.Exec("/bin/mkacc john"))
	ok(asRoot.Exec("/bin/chown john /tmp"))
}
