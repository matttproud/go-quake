package main

import "github.com/matttproud/go-quake/cvar"

func init() {
	commands.Add("edict", noImpl)
	commands.Add("edicts", noImpl)
	commands.Add("edictcount", noImpl)
	commands.Add("profile", noImpl)

	cvars.NewFloat("nomonsters", 0)
	cvars.NewFloat("gamecfg", 0)
	cvars.NewFloat("scratch1", 0)
	cvars.NewFloat("scratch2", 0)
	cvars.NewFloat("scratch3", 0)
	cvars.NewFloat("scratch4", 0)
	cvars.NewFloat("savedgamecfg", 0, cvar.Saved)
	cvars.NewFloat("saved1", 0, cvar.Saved)
	cvars.NewFloat("saved2", 0, cvar.Saved)
	cvars.NewFloat("saved3", 0, cvar.Saved)
	cvars.NewFloat("saved4", 0, cvar.Saved)
}
