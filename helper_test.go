package rs

import (
	"bufio"
	"fmt"
	"io"
	"testing"
)

var John = NewAccount("john", 2)
var Eva = NewAccount("eva", 3)

// test struct
type Alien struct {
	Name string
}

// optional prefix
func newWriter(t *testing.T, prefix ...string) io.Writer {
	r, w := io.Pipe()
	p := "--|"
	if len(prefix) > 0 {
		p = fmt.Sprintf("%s |", prefix[0])
	}
	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			t.Logf("%s %s", p, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			t.Logf("%s %s", p, scanner.Text())
		}
	}()
	return w
}
