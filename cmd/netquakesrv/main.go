// netquakesrv is a classic "Net Quake" dedicated server.
package main

import (
	"github.com/matttproud/go-quake/command"
	"github.com/matttproud/go-quake/cvar"
)

var (
	cvars    = cvar.New()
	commands = command.New()
)

// sys_linux.c:354
func main() {
}
