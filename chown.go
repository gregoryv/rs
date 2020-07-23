package rs

import (
	"flag"
	"fmt"
	"io"
	"strings"

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
	if len(flags.Args()) < 2 {
		return fmt.Errorf("chown: missing path")
	}
	uid, gid, err := me.parseOwner(cmd.Sys, flags.Arg(0))
	if err != nil {
		return err
	}
	for _, path := range cmd.Args[1:] {
		if err := cmd.Sys.SetOwner(path, uid); err != nil {
			return fmt.Errorf("chown: %w", err)
		}
		if gid > -1 {
			// if SetOwner worked so should this
			cmd.Sys.SetGroup(path, gid)
		}
	}
	return nil
}

func (me *Chown) WriteUsage(w io.Writer) {
	p, _ := nexus.NewPrinter(w)
	p.Println("Usage: chown OWNER ...paths")
}

// parseOwner parses OWNER[:GROUP]
func (me *Chown) parseOwner(sys *Syscall, v string) (uid int, gid int, err error) {
	uid = -1
	gid = -1
	parts := strings.Split(v, ":")
	owner := parts[0]
	var acc Account
	if err = sys.LoadAccount(&acc, owner); err != nil {
		return
	}
	uid = acc.uid
	if len(parts) == 2 {
		groupName := parts[1]
		var group Group
		if err = sys.LoadGroup(&group, groupName); err != nil {
			return
		}
		gid = group.gid
	}
	return
}
