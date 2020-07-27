package rs

import (
	"os"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestMkacc(t *testing.T) {
	sys := NewSystem()
	asRoot := Root.Use(sys)
	asAnonymous := Anonymous.Use(sys)
	ok, bad := asserter.NewErrors(t)
	ok(asRoot.Exec("/bin/mkacc -h"))
	ok(asRoot.Exec("/bin/mkacc john")) // automatically use uid,gid 2,2
	bad(asRoot.Exec("/bin/mkacc -uid 2 -gid 2 john")).Log("same uid")
	bad(asRoot.Exec("/bin/mkacc -uid 3 -gid 3 john")).Log("same name")
	bad(asRoot.Exec("/bin/mkacc -uid k -gid 3 john")).Log("uid not int")
	bad(asRoot.Exec("/bin/mkacc")).Log("bad name")
	bad(asRoot.Exec("/bin/mkacc -uid 1 john")).Log("bad uid")
	bad(asRoot.Exec("/bin/mkacc -uid 3 -gid 1 john")).Log("bad gid")
	bad(asAnonymous.Exec("/bin/mkacc -uid 4 -git 4 eva")).Log("unauthorized")
}

func ExampleMkacc_help() {
	asRoot := Root.Use(NewSystem())
	asRoot.Fexec(os.Stdout, "/bin/mkacc -h")
	// output:
	// Usage of mkacc:
	//   -gid int
	//     	optional gid of the new account (default -1)
	//   -uid int
	//     	optional uid of the new account (default -1)
}
