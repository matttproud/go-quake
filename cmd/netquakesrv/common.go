package main

import (
	"errors"

	"github.com/matttproud/go-quake/cvar"
)

func noImpl(...string) error {
	return errors.New("not implemented")
}

func init() {
	cvars.NewFloat("registered", 0)
	cvars.NewString("cmdline", "", cvar.ServerSide)
}
