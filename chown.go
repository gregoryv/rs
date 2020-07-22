package rs

import "fmt"

func Chown(cmd *Cmd) ExecErr {
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
