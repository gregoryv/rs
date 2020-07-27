package rs

import (
	"flag"
	"fmt"
)

// Chmod command sets mode of a resource.
func Chmod(cmd *Cmd) error {
	flags := flag.NewFlagSet("chmod", flag.ContinueOnError)
	mode := flags.Uint("m", 0, "mode")
	flags.SetOutput(cmd.Out)
	if err := flags.Parse(cmd.Args); err != nil {
		return err
	}
	abspath := flags.Arg(0)
	if abspath == "" {
		return fmt.Errorf("missing abspath")
	}
	return cmd.Sys.SetMode(abspath, Mode(*mode))
}
