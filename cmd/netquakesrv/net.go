package main

import "flag"

var (
	flagRecord bool
	flagPort   int
	flagListen bool
)

func init() {
	flag.BoolVar(&flagRecord, "record", false, "whether to record a demo")
	flag.IntVar(&flagPort, "port", 8080, "serving port")
	flag.BoolVar(&flagListen, "listen", false, "whether to listen for connections")

	commands.Add("slist", noImpl)
	commands.Add("listen", noImpl)
	commands.Add("maxplayers", noImpl)
	commands.Add("port", noImpl)

	cvars.NewFloat("net_messagetimeout", 300)
	cvars.NewString("hostname", "UNNAMED")
	// omitted many modem and IPX settings
}
