package rs

import (
	"flag"
	"fmt"

	"github.com/gregoryv/nugo"
)

// Ls lists resources
func Ls(cmd *Cmd) ExecErr {
	flags := flag.NewFlagSet("ls", flag.ContinueOnError)
	longList := flags.Bool("l", false, "use a long listing format")
	recursive := flags.Bool("R", false, "recursive")
	flags.SetOutput(cmd.Out)
	if err := flags.Parse(cmd.Args); err != nil {
		return err
	}
	abspath := flags.Arg(0)
	visitor := func(c *ResInfo, abspath string, w *nugo.Walker) {
		switch {
		case *recursive:
			fmt.Fprintf(cmd.Out, "%s\n", abspath)
		default:
			fmt.Fprintf(cmd.Out, "%s\n", c.Name())
		}
	}
	if *longList {
		visitor = func(c *ResInfo, abspath string, w *nugo.Walker) {
			switch {
			case *recursive:
				fmt.Fprintf(cmd.Out, "%s %s\n", c.node.Seal(), abspath)
			default:
				fmt.Fprintf(cmd.Out, "%s\n", c.node)
			}
		}
	}
	res, err := cmd.Sys.Stat(abspath)
	if err != nil {
		return err
	}
	w := NewWalker(cmd.Sys)
	w.SetRecursive(*recursive)
	if err := res.IsDir(); err != nil {
		visitor(res, abspath, w.w)
		return nil
	}
	return w.Walk(&ResInfo{res.node}, visitor)
}
