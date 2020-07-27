package rs

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

// NewCmd returns a new command.
func NewCmd(abspath string, args ...string) *Cmd {
	return &Cmd{
		Abspath: abspath,
		Args:    args,
		Out:     ioutil.Discard,
	}
}

type Cmd struct {
	Abspath string // of the command
	Args    []string

	// Access to system with a specific account
	Sys *Syscall

	In  io.Reader
	Out io.Writer
}

// String returns the command with its arguments
func (me *Cmd) String() string {
	return fmt.Sprintf("%s %s", me.Abspath, strings.Join(me.Args, " "))
}

type Executable interface {
	Exec(*Cmd) error
}

type ExecFunc func(*Cmd) error

func (me ExecFunc) Exec(cmd *Cmd) error {
	return me(cmd)
}
