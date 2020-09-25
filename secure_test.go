package rs

import (
	"os"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestSecure_createsCredentials_asRoot(t *testing.T) {
	sys := NewSystem()
	asRoot := Root.Use(sys)
	asRoot.AddAccount(John)
	must, _ := asserter.NewFatalErrors(t)
	must(asRoot.Exec("/bin/secure -a john -s mysecret"))
	_, bad := asserter.NewErrors(t)
	oK, _ := asserter.NewMixed(t)
	bad(asRoot.Exec("/bin/secure -a eva -s x"))
	bad(asRoot.Exec("/bin/secure -a john"))
	oK(asRoot.Stat("/etc/credentials"))
}

func TestSecure_login_asRoot(t *testing.T) {
	sys := NewSystem()
	asRoot := Root.Use(sys)
	asRoot.AddAccount(John)
	asRoot.AddAccount(Eva)
	must, _ := asserter.NewFatalErrors(t)
	must(asRoot.Exec("/bin/secure -a john -s secret"))
	must(asRoot.Exec("/bin/secure -a john -s secret1"))
	must(asRoot.Exec("/bin/secure -a eva -s X"))
	ok, bad := asserter.NewErrors(t)
	ok(asRoot.Exec("/bin/secure -c -a john -s secret"))
	ok(asRoot.Exec("/bin/secure -c -a john -s secret1"))
	bad(asRoot.Exec("/bin/secure -c -a john -s hack"))
}

func ExampleSecure() {
	asRoot := Root.Use(NewSystem())
	asRoot.Fexec(os.Stdout, "/bin/secure -h")
	// output:
	// Usage of secure:
	//   -a string
	//     	account name
	//   -c	check if secret is valid
	//   -s string
	//     	secret
}
