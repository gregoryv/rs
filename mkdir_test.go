package rs

import "os"

func ExampleMkdir_help() {
	asRoot := Root.Use(NewSystem())
	asRoot.Fexec(os.Stdout, "/bin/mkdir", "-h")
	// output:
	// Usage of mkdir:
	//   -m uint
	//     	mode for new directory (default 493)

}
