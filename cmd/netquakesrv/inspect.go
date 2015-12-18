package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var httpAddr string

func init() {
	flag.StringVar(&httpAddr, "http-addr", "localhost:8081", "HTTP address for server inspection")
}

func Inspect() error {
	if httpAddr == "" {
		return nil
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- http.ListenAndServe(httpAddr, nil)
	}()
	select {
	case err := <-errCh:
		return err
	case <-time.After(250 * time.Millisecond):
		// Assume no issues after reasonable delay.
		return nil
	}
}
