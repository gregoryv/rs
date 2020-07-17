package rs

import (
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/nugo"
)

func Test_operation_String(t *testing.T) {
	var (
		assert = asserter.New(t)
	)
	assert(OpRead.String() == "read")
	assert(OpWrite.String() == "write")
	assert(OpExec.String() == "exec")
	assert(operation(10).String() == "")
}

func Test_operation_Modes(t *testing.T) {
	ok := func(n, u, g, o nugo.NodeMode) {
		t.Helper()
		assert := asserter.New(t)
		assert(n != u).Error("anonymous eq user")
		assert(n != g).Error("anonymous eq group")
		assert(n != o).Error("anonymous eq other")
		assert(u != g).Error("user eq group")
		assert(u != o).Error("user eq other")
		assert(g != o).Error("group eq other")
	}
	ok(OpRead.Modes())
	ok(OpWrite.Modes())
	ok(OpExec.Modes())
}

func Test_operation_Modes_bad(t *testing.T) {
	defer func() {
		e := recover()
		if e == nil {
			t.Error("should panic")
		}
	}()
	operation(10).Modes()
}
