package main

import (
	"flag"
	"log"
	"os"

	kc "github.com/laher/kc/internal"
)

func main() {
	var (
		fs      = flag.NewFlagSet("kca", flag.ExitOnError)
		verbose = fs.Bool("v", false, "verbose")
		file    = fs.String("f", "", "file to apply")
	)
	contexts, args := kc.Contexts(os.Args[1:])
	fs.Parse(args)
	for _, context := range contexts {
		if len(contexts) > 1 || *verbose {
			log.Printf("context: %s", context)
		}
		e, err := apply(context, *verbose, *file)
		if err != nil {
			log.Printf("Error: %s", err)
			os.Exit(e)
		}
	}
	if *verbose {
		log.Print("done")
	}
}

func apply(context string, verbose bool, file string) (int, error) {
	args := []string{"apply", "-f", file}
	return kc.Run(kc.PrepKC(context, args...), verbose)
}
