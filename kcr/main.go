package main

import (
	"flag"
	"log"
	"os"

	kc "github.com/laher/kc/internal"
)

var (
	context = flag.String("c", "", "kubectl context")
	verbose = flag.Bool("v", false, "verbose")
	file    = flag.String("f", "", "file to apply")
)

func main() {
	flag.Parse()
	e, err := replace(*context, *verbose, *file)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	os.Exit(e)
}

func replace(context string, verbose bool, file string) (int, error) {
	args := []string{"replace", "--cascade", "--force", "-f", file}
	return kc.Run(kc.PrepKC(context, args...), verbose)
}
