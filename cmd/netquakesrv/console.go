package main

func init() {
	commands.Add("toggleconsole", noImpl)
	commands.Add("messagemode", noImpl)
	commands.Add("messagemode2", noImpl)
	commands.Add("clear", noImpl)

	cvars.NewFloat("con_notifytime", 3)
}
