package rs

import (
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/nugo"
)

func TestAccount_joinGroup(t *testing.T) {
	acc := NewAccount("root", 1)
	acc.joinGroup(2)
	acc.joinGroup(2) // nop, already there
	if len(acc.groups) != 2 {
		t.Fail()
	}
}

func TestAccount_DelGroup(t *testing.T) {
	acc := NewAccount("root", 1)
	acc.joinGroup(2)
	acc.DelGroup(2)
	if len(acc.groups) != 1 {
		t.Fail()
	}
}

func TestAccount_Owns(t *testing.T) {
	acc := NewAccount("root", 1)
	n := nugo.NewNode("x")
	n.SetUID(2)
	if acc.Owns(n) {
		t.Error("uid 1 Owns uid 2")
	}
}

func TestAccount_permittedAnonymous(t *testing.T) {
	var (
		perm    = Anonymous.permitted
		ok, bad = asserter.NewErrors(t)
	)
	ok(perm(OpRead, sealed(1, 1, 07000)))
	ok(perm(OpRead, sealed(1, 1, 04000)))
	ok(perm(OpWrite, sealed(1, 1, 02000)))
	ok(perm(OpExec, sealed(1, 1, 01000)))
	bad(perm(OpExec, sealed(1, 1, 02000)))
	bad(perm(OpExec, sealed(1, 1, 00000)))
}

func TestAccount_permittedRoot(t *testing.T) {
	var (
		ok, _ = asserter.NewErrors(t)
		perm  = Root.permitted
	)
	// root is special in that it always has full access
	ok(perm(OpRead, sealed(1, 1, 00000)))
	ok(perm(OpWrite, sealed(1, 2, 00000)))
	ok(perm(OpExec, sealed(0, 0, 00000)))
}

func TestAccount_permittedOther(t *testing.T) {
	perm := NewAccount("john", 2).permitted
	ok, _ := asserter.NewErrors(t)
	ok(perm(OpRead, sealed(2, 2, 00400)))
	ok(perm(OpRead, sealed(3, 2, 00040)))
	ok(perm(OpRead, sealed(1, 1, 00004)))
}

// sealed returns a sealed node
func sealed(uid, gid int, perm nugo.NodeMode) *nugo.Node {
	n := nugo.NewNode("x")
	n.SetUID(uid)
	n.SetGID(gid)
	n.SetPerm(perm)
	return n
}
