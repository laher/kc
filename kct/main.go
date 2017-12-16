package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/laher/kc/internal"
)

func main() {
	var (
		fs      = flag.NewFlagSet("kcl", flag.ExitOnError)
		verbose = fs.Bool("v", false, "verbose")
		label   = fs.String("l", "", "select pod by label")
	)
	contexts, args := kc.Contexts(os.Args[1:])
	fs.Parse(args)
	wg := sync.WaitGroup{}
	for _, context := range contexts {
		if len(contexts) > 1 || *verbose {
			log.Printf("context: %s", context)
		}
		e, err := kc.Logg(context, *verbose, true, 1, *label, fs.Args(), &wg)
		if err != nil {
			log.Printf("Error: %s", err)
			os.Exit(e)
		}
	}
	wg.Wait()
	if *verbose {
		log.Print("done")
	}
}
