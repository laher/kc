package main

import (
	"flag"
	"log"
	"os"

	"github.com/laher/kc/internal"
)

func main() {
	var (
		fs          = flag.NewFlagSet("kc", flag.ExitOnError)
		verbose     = fs.Bool("v", false, "verbose")
		interactive = fs.Bool("i", false, "interactive")
	)
	contexts, args := kc.Contexts(os.Args[1:])
	fs.Parse(args)
	for _, context := range contexts {
		if len(contexts) > 1 || *verbose {
			log.Printf("context: %s", context)
		}
		e, err := kctl(context, *verbose, fs.Args(), *interactive)
		if err != nil {
			log.Printf("Error: %s", err)
			os.Exit(e)
		}
	}
	if *verbose {
		log.Print("done")
	}
}

func kctl(context string, verbose bool, args []string, interactive bool) (int, error) {
	kcArgs := []string{}
	kcArgs = append(kcArgs, args...)
	cmd := kc.PrepKC(context, kcArgs...)
	if interactive {
		cmd.Stdin = os.Stdin
	}
	return kc.Run(cmd, verbose)
}
