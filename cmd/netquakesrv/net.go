package main

import (
	"flag"
	"os"

	"github.com/matttproud/go-quake/cvar"
)

var (
	flagRecord bool
	flagPort   int
	flagListen bool

	cvHostname          *cvar.String
	cvNetMessageTimeout *cvar.Float
)

func isErrTransient(err error) bool {
	if err == nil {
		return true
	}
	type timeouter interface {
		Timeout() bool
	}
	if tout, ok := err.(timeouter); ok {
		if tout.Timeout() {
			return true
		}
	}
	type temporaryer interface {
		Temporary() bool
	}
	if temp, ok := err.(temporaryer); ok {
		if temp.Temporary() {
			return true
		}
	}
	return false
}

func init() {
	flag.BoolVar(&flagRecord, "record", false, "whether to record a demo")
	flag.IntVar(&flagPort, "port", 8080, "serving port")
	flag.BoolVar(&flagListen, "listen", false, "whether to listen for connections")

	commands.Add("slist", noImpl)
	commands.Add("listen", noImpl)
	commands.Add("maxplayers", func(args ...string) error {
		return server.cmdMaxPlayers(args...)
	})
	commands.Add("port", noImpl)
	commands.Add("net_stats", noImpl)
	commands.Add("test", noImpl)
	commands.Add("test2", noImpl)

	cvNetMessageTimeout, _ = cvars.NewFloat("net_messagetimeout", 300)
	cvHostname, _ = cvars.NewString("hostname", "UNNAMED")
	// omitted many modem and IPX settings

	hn, err := os.Hostname()
	if err == nil {
		cvHostname.Set(hn)
	}
}
