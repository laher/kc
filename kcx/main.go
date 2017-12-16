package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"sync"

	kc "github.com/laher/kc/internal"
)

func main() {
	var (
		fs        = flag.NewFlagSet("kcx", flag.ExitOnError)
		context   string
		verbose   = fs.Bool("v", false, "verbose")
		label     = fs.String("l", "", "select pod by label")
		container = fs.String("c", "", "container")
	)
	args := os.Args
	if len(args) > 1 && !strings.HasPrefix(args[1], "-") {
		context = args[1]
		args = args[2:]
	}
	fs.Parse(args)
	e, err := exec(context, *verbose, *label, *container, fs.Args())
	if err != nil {
		log.Printf("Error: %s", err)
	}
	os.Exit(e)

}

func exec(context string, verbose bool, label, container string, args []string) (int, error) {
	kcArgs := []string{"exec"}
	if container != "" {
		kcArgs = append(kcArgs, "-c", container)
	}
	var pods []string
	if label != "" {
		pods = kc.PodsByLabel(context, verbose, label)
	} else {
		if len(args) < 1 {
			log.Printf("Need an identifier for log")
			os.Exit(1)
		}
		pods = []string{args[0]}
		args = args[1:]
	}
	wg := sync.WaitGroup{}
	for _, pod := range pods {
		allArgs := []string{"exec"}
		if len(pods) == 1 {
			allArgs = append(allArgs, "-it", pod)
			allArgs = append(allArgs, args...)
			return kc.Run(kc.PrepKC(context, allArgs...), verbose)
		}
		//no support for interactive mode with multiple-targets
		go func(pod string) {
			wg.Add(1)
			allArgs = append(allArgs, pod)
			allArgs = append(allArgs, args...)
			_, err := kc.Run(kc.PrepKC(context, kcArgs...), verbose)
			if err != nil {
				log.Printf("error: %v", err)
			}
		}(pod)
	}
	wg.Wait()
	return 0, nil
}
