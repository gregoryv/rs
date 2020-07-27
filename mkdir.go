package rs

import (
	"flag"
)

// Mkdir creates directories
func Mkdir(cmd *Cmd) error {
	flags := flag.NewFlagSet("mkdir", flag.ContinueOnError)
	flags.SetOutput(cmd.Out)
	mode := flags.Uint("m", 00755, "mode for new directory")
	if err := flags.Parse(cmd.Args); err != nil {
		return err
	}
	abspath := flags.Arg(0)
	_, err := cmd.Sys.Mkdir(abspath, Mode(*mode))
	return err
}
