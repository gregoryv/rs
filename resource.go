package rs

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gregoryv/nugo"
)

// ResInfo describes a resource and is returned by Stat
type ResInfo struct {
	node *nugo.Node
}

// Name returns the name of the file
func (me *ResInfo) Name() string { return me.node.Name() }

// IsDir returns nil if the resource is a directory
func (me *ResInfo) IsDir() error {
	if !me.node.IsDir() {
		return fmt.Errorf("IsDir: %s not a directory", me.node.Name())
	}
	return nil
}

func newResource(n *nugo.Node, op operation) *Resource {
	return &Resource{
		node: n,
		op:   op,
	}
}

// Resource wraps access to the underlying node
type Resource struct {
	node   *nugo.Node
	op     operation
	unlock func()
	io.Reader
	buf *bytes.Buffer // used for writing
}

// Read reads from the underlying source. Fails if not readable or
// resource is in write mode.
func (me *Resource) Read(b []byte) (int, error) {
	if me.writeOnly() {
		return 0, fmt.Errorf("Read: %s write only", me.node.Name())
	}
	if me.buf == nil {
		return 0, fmt.Errorf("Read: unreadable source")
	}
	return me.buf.Read(b)
}

// Write writes to the resource. Is not flushed until closed.
func (me *Resource) Write(p []byte) (int, error) {
	if me.readOnly() {
		return 0, fmt.Errorf("Write: %s read only", me.node.Name())
	}
	return me.buf.Write(p)
}

// Close closes the resource. If resource is in write mode the written
// buffer is flushed.
func (me *Resource) Close() error {
	if me.writeOnly() {
		me.node.Content = me.buf.Bytes()
	}
	me.buf = nil
	me.unlock()
	return nil
}

// read + write should not be possible
func (me *Resource) readOnly() bool  { return me.op&OpRead != 0 }
func (me *Resource) writeOnly() bool { return me.op&OpWrite != 0 }
