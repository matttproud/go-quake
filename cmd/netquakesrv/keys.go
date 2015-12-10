package main

func init() {
	commands.Add("bind", noImpl)
	commands.Add("unbind", noImpl)
	commands.Add("unbindall", noImpl)
}
