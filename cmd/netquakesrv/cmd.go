package main

func init() {
	commands.Add("stuffcmds", noImpl)
	commands.Add("exec", noImpl)
	commands.Add("echo", noImpl)
	commands.Add("alias", noImpl)
	commands.Add("cmd", noImpl)
	commands.Add("wait", noImpl)
}
