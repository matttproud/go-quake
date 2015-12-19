package main

import (
	"github.com/matttproud/go-quake/prog"

	. "github.com/matttproud/go-quake/qtype"
)

const maxEntLeafs = 16

type EntityState struct {
	Origin     Vec3
	Angles     Vec3
	ModelIndex int32
	Frame      int32
	ColorMap   int32
	Skin       int32
	Effects    int32
}

type Edict struct {
	// free
	// area

	LeafCount int32
	LeafNums  [maxEntLeafs]int16
	State     EntityState
	V         prog.EntVars
}

type Player struct {
}
