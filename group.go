package rs

import (
	"encoding/gob"
	"io"
)

type Group struct {
	name string
	gid  int
}

type group struct {
	Name string
	GID  int
}

// WriteTo
func (me *Group) WriteTo(w io.Writer) (int64, error) {
	g := group{
		Name: me.name,
		GID:  me.gid,
	}
	return 0, gob.NewEncoder(w).Encode(&g)
}

// ReadFrom
func (me *Group) ReadFrom(r io.Reader) (int64, error) {
	var g group
	err := gob.NewDecoder(r).Decode(&g)
	me.name = g.Name
	me.gid = g.GID
	return 0, err
}
