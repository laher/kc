package kc

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func HandleQuit() chan struct{} {
	c := make(chan os.Signal, 1)
	d := make(chan struct{})
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-c
		// sig is a ^C, handle it
		log.Printf("Received signal %v", sig.String())
		close(d)
	}()

	return d
}
