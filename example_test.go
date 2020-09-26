package rs

import (
	"fmt"
	"os"
)

func Example_defaultResourceSystem() {
	sys := NewSystem()
	asRoot := Root.Use(sys)
	asRoot.Fexec(os.Stdout, "/bin/ls -R -l /")
	// output:
	// d--xrwxr-xr-x 1 1 /bin
	// ----rwxr-xr-x 1 1 /bin/chmod
	// ----rwxr-xr-x 1 1 /bin/chown
	// ---xrwxr-xr-x 1 1 /bin/ls
	// ----rwxr-xr-x 1 1 /bin/mkacc
	// ----rwxr-xr-x 1 1 /bin/mkdir
	// ----rwxr-xr-x 1 1 /bin/secure
	// d---rwxr-xr-x 1 1 /etc
	// d---rwxr-xr-x 1 1 /etc/accounts
	// ----rw-r--r-- 1 1 /etc/accounts/anonymous
	// ----rw-r--r-- 1 1 /etc/accounts/root
	// d---rwxr-xr-x 1 1 /etc/groups
	// ----rw-r--r-- 1 1 /etc/groups/anonymous
	// ----rw-r--r-- 1 1 /etc/groups/root
	// drwxrwxrwxrwx 1 1 /tmp
}

func Example_saveAndLoadResource() {
	sys := NewSystem()
	asRoot := Root.Use(sys)
	asRoot.Exec("/bin/mkdir /tmp/aliens")
	asRoot.Save("/tmp/aliens/green.gob", &Alien{Name: "Mr Green"})
	var alien Alien
	asRoot.Load(&alien, "/tmp/aliens/green.gob")
	fmt.Printf("%#v", alien)
	// output:
	// rs.Alien{Name:"Mr Green"}
}

// When exporting a system each node is serialized using base64 encoding of the content, if any.
func Example_exportSystem() {
	sys := NewSystem()
	sys.Export(os.Stdout)
	// output:
	// 4026532845 1 1 /
	// 3758097389 1 1 /bin
	// 3758096877 1 1 /etc
	// 3758096877 1 1 /etc/accounts
	// 1610613156 1 1 /etc/accounts/anonymous Mv+DAwEBB0FjY291bnQB/4QAAQMBBE5hbWUBDAABA1VJRAEEAAEGR3JvdXBzAf+GAAAAE/+FAgEBBVtdaW50Af+GAAEEAAAR/4QBCWFub255bW91cwIB
	// 1610613156 1 1 /etc/accounts/root Mv+DAwEBB0FjY291bnQB/4QAAQMBBE5hbWUBDAABA1VJRAEEAAEGR3JvdXBzAf+GAAAAE/+FAgEBBVtdaW50Af+GAAEEAAAO/4QBBHJvb3QBAgEB
	// 3758096877 1 1 /etc/groups
	// 1610613156 1 1 /etc/groups/anonymous JP+BAwEBBWdyb3VwAf+CAAECAQROYW1lAQwAAQNHSUQBBAAAAA7/ggEJYW5vbnltb3Vz
	// 1610613156 1 1 /etc/groups/root JP+BAwEBBWdyb3VwAf+CAAECAQROYW1lAQwAAQNHSUQBBAAAAAv/ggEEcm9vdAEC
	// 3758100479 1 1 /tmp
}
