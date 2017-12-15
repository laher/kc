package main

import (
	"flag"
	"log"
	"os"

	"github.com/laher/kc/internal"
)

var (
	context = flag.String("c", "", "kubectl context")
	verbose = flag.Bool("v", false, "verbose")
)

func main() {
	flag.Parse()
	e, err := kc.Logg(*context, *verbose, false, 0, flag.Args())
	if err != nil {
		log.Printf("Error: %s", err)
	}
	os.Exit(e)
}
