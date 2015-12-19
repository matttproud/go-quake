package main

import "github.com/matttproud/go-quake/cvar"

var sysTicRate *cvar.Float

func init() {
	cvars.NewFloat("host_framerate", 0)
	cvars.NewFloat("host_speeds", 0)

	sysTicRate, _ = cvars.NewFloat("sys_ticrate", 0.05)
	cvars.NewFloat("serverprofile", 0)

	cvars.NewFloat("fraglimit", 0, cvar.ServerSide)
	cvars.NewFloat("timelimit", 0, cvar.ServerSide)
	cvars.NewFloat("teamplay", 0, cvar.ServerSide)

	cvars.NewFloat("samelevel", 0)
	cvars.NewFloat("noexit", 0, cvar.ServerSide)

	cvars.NewFloat("developer", 0)

	cvars.NewFloat("skill", 1)
	cvars.NewFloat("deathmatch", 0)
	cvars.NewFloat("coop", 0)

	cvars.NewFloat("pausable", 1)

	cvars.NewFloat("temp1", 0)
}
