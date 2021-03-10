package rs

import (
	"fmt"
	"strings"
	"testing"
)

func BenchmarkSyscall_Mkdir(b *testing.B) {
	asRoot := Root.Use(NewSystem())
	for i := 0; i < b.N; i++ {
		asRoot.Mkdir(fmt.Sprintf("/dir%d", i), 0)
	}
}

func Benchmark_Stat(b *testing.B) {
	asRoot := Root.Use(NewSystem())
	for i := 0; i < b.N; i++ {
		asRoot.Stat("/etc/accounts.gob")
	}
}

func Test_mem(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	asRoot := Root.Use(NewSystem())

	kb := 1024
	fileSize := 4 * kb
	content := strings.Repeat("x", fileSize)
	n := 4000
	for i := 0; i < n; i++ {
		dir := fmt.Sprintf("/tmp/dir%d", i)

		asRoot.Mkdir(dir, 00755)
		filename := fmt.Sprintf("%s/file.gob", dir)

		err := asRoot.Save(filename, content)
		if err != nil {
			t.Fatal(err)
		}

	}
}
