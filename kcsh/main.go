package main

import (
	"flag"
	"log"
	"os"

	kc "github.com/laher/kc/internal"
)

var (
	context   = flag.String("c", "", "kubectl context")
	container = flag.String("co", "", "container")
	verbose   = flag.Bool("v", false, "verbose")
)

func main() {
	flag.Parse()
	e, err := shell(*context, *verbose, *container, flag.Args())
	if err != nil {
		log.Printf("Error: %s", err)
	}
	os.Exit(e)

}

func shell(context string, verbose bool, container string, args []string) (int, error) {
	if len(args) < 1 {
		log.Printf("Need an identifier for exec")
		os.Exit(1)
	}
	kcArgs := []string{"exec", "-it"}
	if container != "" {
		kcArgs = append(kcArgs, "-c", container)
	}
	p := args[0]
	pod := kc.Pod(context, verbose, p)
	kcArgs = append(kcArgs, pod, "sh")
	kcArgs = append(kcArgs, args[1:]...)
	return kc.Run(kc.PrepKC(context, kcArgs...), verbose)
}
