package rs

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestLs(t *testing.T) {
	sys := NewSystem()
	asRoot := Root.Use(sys)
	asJohn := NewAccount("john", 2).Use(sys)
	var buf bytes.Buffer
	ok, bad := asserter.NewErrors(t)
	bad(asRoot.Exec("/bin/ls -xx"))
	bad(asRoot.Exec("/bin/ls /nosuch"))
	// ls directory is covered by Examples
	// ls file
	asRoot.Fexec(&buf, "/bin/ls", "-l", "/etc/accounts/root.acc")
	exp := "----rw-r--r-- 1 1 root.acc\n"
	assert := asserter.New(t)
	assert().Equals(buf.String(), exp)

	// only list accessible
	buf.Reset()
	n, _ := asRoot.stat("/etc")
	n.SetPerm(0)
	ok(asJohn.Fexec(&buf, "/bin/ls", "-R", "/"))
	if strings.Contains(buf.String(), "/etc/accounts") {
		t.Error("listed /etc")
	}

	// only list accessible in long format
	buf.Reset()
	ok(asJohn.Fexec(&buf, "/bin/ls", "-R", "-l", "/"))
	if strings.Contains(buf.String(), "/etc/accounts") {
		t.Error("listed /etc")
	}
}

func Example_Ls_fix() {
	asRoot := Root.Use(NewSystem())
	asRoot.Exec("/bin/mkdir /a")
	asRoot.Exec("/bin/mkdir /a/b")
	asRoot.Fexec(os.Stdout, "/bin/ls", "-R", "/a")
	// output:
	// /a/b
}

func ExampleLs() {
	Anonymous.Use(NewSystem()).Fexec(os.Stdout, "/bin/ls", "/")
	// output:
	// bin
	// etc
	// tmp
}

func ExampleLs_longListFormat() {
	Anonymous.Use(NewSystem()).Fexec(os.Stdout, "/bin/ls", "-l", "/")
	// output:
	// d--xrwxr-xr-x 1 1 bin
	// d---rwxr-xr-x 1 1 etc
	// drwxrwxrwxrwx 1 1 tmp
}

func ExampleLs_help() {
	asRoot := Root.Use(NewSystem())
	asRoot.Fexec(os.Stdout, "/bin/ls", "-h")
	// output:
	// Usage of ls:
	//   -R	recursive
	//   -l	use a long listing format
}
