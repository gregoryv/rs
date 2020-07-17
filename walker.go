package rs

import (
	"github.com/gregoryv/nugo"
)

func NewWalker(sys *Syscall) *Walker {
	return &Walker{
		w:   nugo.NewWalker(),
		sys: sys,
	}
}

type Walker struct {
	w   *nugo.Walker
	sys *Syscall
}

// SetRecursive
func (me *Walker) SetRecursive(r bool) { me.w.SetRecursive(r) }

func (me *Walker) Walk(res *ResInfo, fn Visitor) error {
	// wrap the visitor with access control
	visitor := func(n *nugo.Node, abspath string, w *nugo.Walker) {
		if me.sys.acc.permitted(OpExec, n) != nil {
			w.SkipChild()
		}
		c := &ResInfo{n}
		fn(c, abspath, w)
	}
	me.w.Walk(res.node, visitor)
	return nil
}
