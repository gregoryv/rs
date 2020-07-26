package rs

import (
	"flag"
	"fmt"
)

// Mkacc creates an account.
func Mkacc(cmd *Cmd) ExecErr {
	flags := flag.NewFlagSet("mkacc", flag.ContinueOnError)
	uid := flags.Int("uid", -1, "optional uid of the new account")
	gid := flags.Int("gid", -1, "optional gid of the new account")
	flags.SetOutput(cmd.Out)
	if err := flags.Parse(cmd.Args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}
	Name := flags.Arg(0)
	if Name == "" {
		return fmt.Errorf("missing account name")
	}
	if *uid == -1 {
		*uid = cmd.Sys.NextUID()
	}
	if *gid == -1 {
		*gid = cmd.Sys.NextGID()
	}
	if *uid < 2 {
		return fmt.Errorf("invalid uid")
	}
	if *gid < 2 {
		return fmt.Errorf("invalid gid")
	}
	acc := NewAccount(Name, *uid)
	acc.Groups[0] = *gid
	return cmd.Sys.AddAccount(acc)
}
