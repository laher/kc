package main

import (
	"flag"
	"log"
	"os"

	kc "github.com/laher/kc/internal"
)

var (
	context    = flag.String("c", "", "kubectl context")
	verbose    = flag.Bool("v", false, "verbose")
	deployment = flag.String("d", "", "deployment")
)

func main() {
	flag.Parse()
	e, err := bounce(*context, *verbose, *deployment, flag.Args())
	if err != nil {
		log.Printf("Error: %s", err)
	}
	os.Exit(e)
}

func bounce(context string, verbose bool, deployment string, args []string) (int, error) {
	kcargs := []string{"scale", "--current-replicas=1", "--replicas=0", "deploy", deployment}
	ex, err := kc.Run(kc.PrepKC(context, kcargs...), verbose)
	if err != nil {
		return ex, err
	}
	kcargs = []string{"scale", "--current-replicas=0", "--replicas=1", "deploy", deployment}
	return kc.Run(kc.PrepKC(context, kcargs...), verbose)
}
