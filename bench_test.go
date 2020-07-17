package rs

import (
	"fmt"
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
