package main

import (
	"flag"
	"log"
	"os"
	"strings"

	kc "github.com/laher/kc/internal"
)

func main() {
	var (
		fs      = flag.NewFlagSet("kca", flag.ExitOnError)
		context string
		verbose = fs.Bool("v", false, "verbose")
		file    = fs.String("f", "", "file to apply")
	)
	args := os.Args
	if len(args) > 1 && !strings.HasPrefix(args[1], "-") {
		context = args[1]
		args = args[2:]
	}
	fs.Parse(args)
	e, err := apply(context, *verbose, *file)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	os.Exit(e)
}

func apply(context string, verbose bool, file string) (int, error) {
	args := []string{"apply", "-f", file}
	return kc.Run(kc.PrepKC(context, args...), verbose)
}
