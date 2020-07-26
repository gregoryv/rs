package rs

import (
	"encoding/gob"
	"io"
)

type Group struct {
	Name string
	gid  int
}

type group struct {
	Name string
	GID  int
}

// WriteTo
func (me *Group) WriteTo(w io.Writer) (int64, error) {
	g := group{
		Name: me.Name,
		GID:  me.gid,
	}
	return 0, gob.NewEncoder(w).Encode(&g)
}

// ReadFrom
func (me *Group) ReadFrom(r io.Reader) (int64, error) {
	var g group
	err := gob.NewDecoder(r).Decode(&g)
	me.Name = g.Name
	me.gid = g.GID
	return 0, err
}
