package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/laher/kc/internal"
)

func main() {
	var (
		fs      = flag.NewFlagSet("count", flag.ExitOnError)
		context string
		verbose = fs.Bool("v", false, "verbose")
		label   = fs.String("l", "", "select pod by label")
	)
	args := os.Args
	log.Printf("args: %v\n", args)
	if len(args) > 1 && !strings.HasPrefix(args[1], "-") {
		context = args[1]
		args = args[2:]
	}
	log.Printf("args: %v\n", args)
	if err := fs.Parse(args); err != nil {
		log.Printf("Error: %s", err)
		os.Exit(1)
	}
	e, err := kc.Logg(context, *verbose, false, -1, *label, fs.Args())
	if err != nil {
		log.Printf("Error: %s", err)
	}
	os.Exit(e)
}
