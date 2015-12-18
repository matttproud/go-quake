package main

import "github.com/matttproud/go-quake/bsp"

var noVis [bsp.MaxMapLeafs / 8]byte

func init() {
	for i := range noVis[:] {
		noVis[i] = 0xff
	}
}
