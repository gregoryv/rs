package rs

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/gregoryv/nugo"
)

var (
	Anonymous = NewAccount("anonymous", 0)
	Root      = NewAccount("root", 1)
)

// NewAccount returns a new account with the given uid as both uid and
// group id.
func NewAccount(name string, uid int) *Account {
	return &Account{
		name:   name,
		uid:    uid,
		groups: []int{uid},
	}
}

type Account struct {
	name string
	uid  int

	mu     sync.Mutex
	groups []int
}

type account struct {
	Name   string
	UID    int
	Groups []int
}

// WriteTo
func (me *Account) WriteTo(w io.Writer) (int64, error) {
	a := account{
		Name:   me.name,
		UID:    me.uid,
		Groups: me.groups,
	}
	return 0, gob.NewEncoder(w).Encode(&a)
}

// ReadFrom
func (me *Account) ReadFrom(r io.Reader) (int64, error) {
	var a account
	err := gob.NewDecoder(r).Decode(&a)
	me.name = a.Name
	me.uid = a.UID
	me.groups = a.Groups
	return 0, err
}

// gid returns the first group id of the account
func (my *Account) gid() int { return my.groups[0] }

// todo hide as command
func (me *Account) joinGroup(gid int) {
	for _, id := range me.groups {
		if id == gid {
			return
		}
	}
	me.mu.Lock()
	me.groups = append(me.groups, gid)
	me.mu.Unlock()
}

// todo hide as command
func (me *Account) leaveGroup(gid int) {
	for i, id := range me.groups {
		if id == gid {
			me.mu.Lock()
			me.groups = append(me.groups[:i], me.groups[i+1:]...)
			me.mu.Unlock()
			return
		}
	}
}

// Use returns a Syscall struct for accessing the system.
func (me *Account) Use(sys *System) *Syscall {
	return &Syscall{
		System:  sys,
		acc:     me,
		auditer: sys.auditer,
	}
}

// Owns returns tru if the account uid mathes the given id
func (me *Account) Owns(s nugo.Sealed) bool {
	return me.uid == s.Seal().UID
}

// permitted returns error if account does not have operation
// permission to the given seal.
func (my *Account) permitted(op operation, s nugo.Sealed) error {
	if my.uid == Root.uid {
		return nil
	}
	n, u, g, o := op.Modes()
	seal := s.Seal()
	switch {
	case my.uid == 0 && (seal.Mode&n == n): // anonymous
	case my.uid == seal.UID && (seal.Mode&u == u): // owner
	case my.member(seal.GID) && (seal.Mode&g == g): // group
	case my.uid > 0 && seal.Mode&o == o: // other
	default:
		return fmt.Errorf("%v %v denied", seal, op)
	}
	return nil
}

var ErrPermissionDenied = errors.New("permission denied")

func (my *Account) member(gid int) bool {
	for _, id := range my.groups {
		if id == gid {
			return true
		}
	}
	return false
}
