package rs

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gregoryv/nugo"
)

var (
	Anonymous = NewAccount("anonymous", 0)
	Root      = NewAccount("root", 1)
)

// NewAccount returns a new account with the given uid as both uid and
// group id.
func NewAccount(Name string, uid int) *Account {
	return &Account{
		Name:   Name,
		UID:    uid,
		Groups: []int{uid},
	}
}

type Account struct {
	Name string
	UID  int

	mu     sync.Mutex
	Groups []int
}

// gid returns the first group id of the account
func (me *Account) gid() int { return me.Groups[0] }

// todo hide as command
func (me *Account) joinGroup(gid int) {
	for _, id := range me.Groups {
		if id == gid {
			return
		}
	}
	me.mu.Lock()
	me.Groups = append(me.Groups, gid)
	me.mu.Unlock()
}

// todo hide as command
func (me *Account) leaveGroup(gid int) {
	for i, id := range me.Groups {
		if id == gid {
			me.mu.Lock()
			me.Groups = append(me.Groups[:i], me.Groups[i+1:]...)
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

// owns returns true if the account uid mathes the given id
func (me *Account) owns(s nugo.Sealed) bool {
	return me.UID == s.Seal().UID
}

// permitted returns error if account does not have operation
// permission to the given seal.
func (me *Account) permitted(op operation, s nugo.Sealed) error {
	if me.UID == Root.UID {
		return nil
	}
	n, u, g, o := op.Modes()
	seal := s.Seal()
	switch {
	case me.UID == 0 && (seal.Mode&n == n): // anonymous
	case me.UID == seal.UID && (seal.Mode&u == u): // owner
	case me.member(seal.GID) && (seal.Mode&g == g): // group
	case me.UID > 0 && seal.Mode&o == o: // other
	default:
		return fmt.Errorf("%v %v denied", seal, op)
	}
	return nil
}

var ErrPermissionDenied = errors.New("permission denied")

func (me *Account) member(gid int) bool {
	for _, id := range me.Groups {
		if id == gid {
			return true
		}
	}
	return false
}
