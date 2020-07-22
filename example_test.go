package rs

import "os"

func Example_defaultResourceSystem() {
	Root.Use(NewSystem()).Fexec(os.Stdout, "/bin/ls", "-R", "-l", "/")
	// output:
	// d--xrwxr-xr-x 1 1 /bin
	// ----rwxr-xr-x 1 1 /bin/chmod
	// ----rwxr-xr-x 1 1 /bin/chown
	// ---xrwxr-xr-x 1 1 /bin/ls
	// ----rwxr-xr-x 1 1 /bin/mkacc
	// ----rwxr-xr-x 1 1 /bin/mkdir
	// d---rwxr-xr-x 1 1 /etc
	// d---rwxr-xr-x 1 1 /etc/accounts
	// ----rw-r--r-- 1 1 /etc/accounts/anonymous.acc
	// ----rw-r--r-- 1 1 /etc/accounts/root.acc
	// drwxrwxrwxrwx 1 1 /tmp
}
