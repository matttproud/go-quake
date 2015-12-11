package main

func init() {
	commands.Add("togglemenu", noImpl)
	commands.Add("menu_main", noImpl)
	commands.Add("menu_singleplayer", noImpl)
	commands.Add("menu_load", noImpl)
	commands.Add("menu_save", noImpl)
	commands.Add("menu_multiplayer", noImpl)
	commands.Add("menu_setup", noImpl)
	commands.Add("menu_options", noImpl)
	commands.Add("menu_keys", noImpl)
	commands.Add("menu_video", noImpl)
	commands.Add("help", noImpl)
	commands.Add("menu_quit", noImpl)
}
