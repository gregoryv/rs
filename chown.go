package rs

import (
	"flag"
	"fmt"
	"io"

	"github.com/gregoryv/nexus"
)

type Chown struct{}

func (me *Chown) Exec(cmd *Cmd) ExecErr {
	flags := flag.NewFlagSet("chown", flag.ContinueOnError)
	flags.Usage = func() { me.WriteUsage(cmd.Out) }
	flags.SetOutput(cmd.Out)
	if err := flags.Parse(cmd.Args); err != nil {
		return err
	}

	if len(cmd.Args) < 2 {
		return fmt.Errorf("chown: missing path")
	}
	owner := cmd.Args[0]
	var acc Account
	err := cmd.Sys.Load(&acc, "/etc/accounts/"+owner+".acc")
	if err != nil {
		return err
	}
	for _, path := range cmd.Args[1:] {
		if err := cmd.Sys.SetOwner(path, acc.uid); err != nil {
			return fmt.Errorf("chown: %w", err)
		}
	}
	return nil
}

// WriteUsage
func (me *Chown) WriteUsage(w io.Writer) {
	p, _ := nexus.NewPrinter(w)
	p.Println("Usage: chown OWNER ...paths")
}
