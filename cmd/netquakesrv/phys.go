package main

import "github.com/matttproud/go-quake/cvar"

func init() {
	commands.Add("v_cshift", noImpl)
	commands.Add("bf", noImpl)
	commands.Add("centerview", noImpl)

	cvars.NewFloat("sv_friction", 4, cvar.ServerSide)
	cvars.NewFloat("sv_stopspeed", 100)
	cvars.NewFloat("sv_gravity", 800, cvar.ServerSide)
	cvars.NewFloat("sv_maxvelocity", 2000)
	cvars.NewFloat("sv_nostep", 0)
	cvars.NewFloat("edgefriction", 2)
	cvars.NewFloat("sv_maxspeed", 320, cvar.ServerSide)
	cvars.NewFloat("sv_accelerate", 10)
	cvars.NewFloat("sv_idealpitchscale", 0.8)
	cvars.NewFloat("sv_aim", 0.93)
	cvars.NewFloat("sv_nostep", 0)
}
