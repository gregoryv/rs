package rs

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type Executable interface {
	Exec(*Cmd) ExecErr
}

type ExecFunc func(*Cmd) ExecErr

func (me ExecFunc) Exec(cmd *Cmd) ExecErr { return me(cmd) }

type ExecErr error

// NewCmd returns a new command.
func NewCmd(abspath string, args ...string) *Cmd {
	return &Cmd{
		Abspath: abspath, Args: args, Out: ioutil.Discard}
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
