package rs

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"github.com/gregoryv/nugo"
)

// Ls lists resources
func Ls(cmd *Cmd) ExecErr {
	flags := flag.NewFlagSet("ls", flag.ContinueOnError)
	longList := flags.Bool("l", false, "use a long listing format")
	jsonFmt := flags.Bool("json", false, "write json")
	recursive := flags.Bool("R", false, "recursive")
	flags.SetOutput(cmd.Out)
	if err := flags.Parse(cmd.Args); err != nil {
		return err
	}
	abspath := flags.Arg(0)
	res, err := cmd.Sys.Stat(abspath)
	if err != nil {
		return err
	}
	var format formatter
	switch {
	case *jsonFmt:
		jf := &jsonFormat{
			recursive: *recursive,
			long:      *longList,
			out:       cmd.Out,
		}
		jf.Open()
		defer jf.Close()
		format = jf
	default:
		format = &textFormat{
			recursive: *recursive,
			long:      *longList,
			out:       cmd.Out,
		}
	}
	w := NewWalker(cmd.Sys)
	w.SetRecursive(*recursive)
	if err := res.IsDir(); err != nil {
		format.Visit(res, abspath, w.w)
		return nil
	}
	return w.Walk(&ResInfo{res.node}, format.Visit)
}

type formatter interface {
	Visit(c *ResInfo, abspath string, w *nugo.Walker)
}

type jsonFormat struct {
	recursive bool
	long      bool
	out       io.Writer
	separator string
	*json.Encoder
}

// Open
func (me *jsonFormat) Open()  { fmt.Fprint(me.out, "[") }
func (me *jsonFormat) Close() { fmt.Fprint(me.out, "]") }

func (me *jsonFormat) Visit(c *ResInfo, abspath string, w *nugo.Walker) {
	fmt.Fprint(me.out, me.separator)
	fmt.Fprint(me.out, "{")
	fmt.Fprintf(me.out, "%q: %q", "name", c.Name())
	if me.long {
		seal := c.node.Seal()
		fmt.Fprintf(me.out, ", %q: %q", "mode", seal.Mode)
		fmt.Fprintf(me.out, `, %q: "%v"`, "uid", seal.UID)
		fmt.Fprintf(me.out, `, %q: "%v"`, "gid", seal.GID)
	}
	fmt.Fprint(me.out, "}")
	me.separator = ","
}

type textFormat struct {
	recursive bool
	long      bool
	out       io.Writer
}

func (me *textFormat) Visit(c *ResInfo, abspath string, w *nugo.Walker) {
	switch {
	case me.recursive && !me.long:
		fmt.Fprintf(me.out, "%s\n", abspath)
	case me.recursive && me.long:
		fmt.Fprintf(me.out, "%s %s\n", c.node.Seal(), abspath)
	case !me.recursive && !me.long:
		fmt.Fprintf(me.out, "%s\n", c.Name())
	case !me.recursive && me.long:
		fmt.Fprintf(me.out, "%s\n", c.node)
	}
}
