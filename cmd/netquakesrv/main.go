// netquakesrv is a classic "Net Quake" dedicated server.
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"golang.org/x/net/context"

	"github.com/matttproud/go-quake/command"
	"github.com/matttproud/go-quake/cvar"
	"github.com/matttproud/go-quake/prog"
)

var (
	cvars    = cvar.New()
	commands = command.New()

	server *Server
)

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handleInterrupt(ctx, cancel)
	log.Println("Starting inspection subsystems ...")
	if err := Inspect(); err != nil {
		log.Println(err)
		return
	}
	log.Println("[DONE] Starting inspection subsystems")
	log.Println("Finding game assets ...")
	assets, err := GamePath()
	if err != nil {
		log.Println(err)
		return
	}
	defer assets.Close()
	log.Println("[DONE] Finding game assets")
	log.Println("Preparing the game virtual machine ...")
	progs, err := assets.Load("progs.dat")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = prog.Open(progs)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("[DONE] Preparing the game virtual machine")
	log.Println("Beginning listening for new clients ...")
	conn, err := Listen()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("[DONE] Beginning listening for new clients ...")
	log.Println("Running main loop ...")
	sessions := make(SessionRegistry)
	serverCancel := func() {
		sessions.Close()
		cancel()
	}
	server = &Server{
		State:      Waiting,
		Conn:       conn,
		MaxPlayers: 1,
		Sessions:   sessions,
		closeSig:   make(chan struct{}),
		Cancel:     serverCancel,
	}
	defer server.Close()
	if err := server.Loop(ctx); err != nil {
		log.Println(err)
		return
	}
	log.Println("[DONE] Running main loop")
}

//func loop(ctx context.Context) error {
//	last := time.Now()
//	for {
//		select {
//		case <-ctx.Done():
//			return ctx.Err()
//		default:
//			//	if err := hostFrame(last); err != nil {
//			//		return err
//			//	}
//			took := time.Since(last)
//			sleep := time.Duration(float32(time.Second) * sysTicRate.Get())
//			time.Sleep(sleep - took)
//			last = time.Now()
//		}
//	}
//	return nil
//}

func handleInterrupt(ctx context.Context, cancel func()) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill)
	select {
	case <-ctx.Done():
	case sig := <-sigCh:
		log.Printf("Received %v; terminating ...", sig)
		cancel()
	}
}
