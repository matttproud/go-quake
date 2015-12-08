package main

import "github.com/matttproud/go-quake/cvar"

func init() {
	commands.Add("v_cshift", noImpl)
	commands.Add("bf", noImpl)
	commands.Add("centerview", noImpl)

	cvars.NewFloat("lcd_x", 0)
	cvars.NewFloat("lcd_yaw", 0)

	cvars.NewFloat("scr_ofsx", 0)
	cvars.NewFloat("scr_ofsy", 0)
	cvars.NewFloat("scr_ofsz", 0)

	cvars.NewFloat("cl_rollspeed", 200)
	cvars.NewFloat("cl_rollangle", 2.0)

	cvars.NewFloat("cl_bob", 0.02)
	cvars.NewFloat("cl_bobup", 0.5)

	cvars.NewFloat("v_kicktime", 0.5)
	cvars.NewFloat("v_kickroll", 0.5)
	cvars.NewFloat("v_kickpitch", 0.5)

	cvars.NewFloat("v_iyaw_cycle", 2)
	cvars.NewFloat("v_iroll_cycle", 0.5)
	cvars.NewFloat("v_ipitch_cycle", 1)
	cvars.NewFloat("v_iyaw_level", 0.3)
	cvars.NewFloat("v_iroll_level", 0.1)
	cvars.NewFloat("v_ipitch_level", 0.3)

	cvars.NewFloat("v_idlescale", 0)

	cvars.NewFloat("crosshair", 0, cvar.Saved)
	cvars.NewFloat("cl_crossx", 0)
	cvars.NewFloat("cl_crossy", 0)

	cvars.NewFloat("gl_cshiftpercent", 100)

	cvars.NewFloat("gamma", 1, cvar.Saved)
}
