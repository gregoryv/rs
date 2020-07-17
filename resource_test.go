package rs

import (
	"bytes"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestResInfo_Name(t *testing.T) {
	var (
		sys    = NewSystem()
		asRoot = Root.Use(sys)
		rif, _ = asRoot.Stat("/")
	)
	if rif.Name() != "/" {
		t.Error("name failed")
	}
}

func TestResInfo_IsDir(t *testing.T) {
	var (
		sys     = NewSystem()
		asRoot  = Root.Use(sys)
		dir, _  = asRoot.Stat("/")
		file, _ = asRoot.Stat("/bin/mkdir")
		ok, bad = asserter.NewErrors(t)
	)
	ok(dir.IsDir())
	bad(file.IsDir())
}

func TestResource_Read(t *testing.T) {
	var (
		ok, bad = asserter.NewMixed(t)
		b       = make([]byte, 10)
	)
	r := &Resource{op: OpRead, buf: bytes.NewBufferString("hello")}
	ok(r.Read(b))
	r = &Resource{op: OpRead}
	bad(r.Read(b))
}
