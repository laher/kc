package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	kc "github.com/laher/kc/internal"
)

func main() {
	var (
		fs         = flag.NewFlagSet("kcb", flag.ExitOnError)
		verbose    = fs.Bool("v", false, "verbose")
		deployment = flag.String("d", "", "deployment")
		current    = flag.Int("c", 1, "current replicas")
	)
	contexts, args := kc.Contexts(os.Args[1:])
	fs.Parse(args)
	for _, context := range contexts {
		if len(contexts) > 1 || *verbose {
			log.Printf("context: %s", context)
		}
		e, err := bounce(context, *verbose, *deployment, *current, flag.Args())
		if err != nil {
			log.Printf("Error: %s", err)
			os.Exit(e)
		}
	}
	if *verbose {
		log.Print("done")
	}
}

func bounce(context string, verbose bool, deployment string, current int, args []string) (int, error) {
	if deployment == "" {
		return 1, errors.New("Resource name cannot be empty")
	}
	kcargs := []string{"scale", fmt.Sprintf("--current-replicas=%d", current), "--replicas=0", "deploy", deployment}
	ex, err := kc.Run(kc.PrepKC(context, kcargs...), verbose)
	if err != nil {
		return ex, err
	}
	kcargs = []string{"scale", "--current-replicas=0", fmt.Sprintf("--replicas=%d", current), "deploy", deployment}
	return kc.Run(kc.PrepKC(context, kcargs...), verbose)
}
