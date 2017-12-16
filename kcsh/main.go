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
		fs        = flag.NewFlagSet("kcsh", flag.ExitOnError)
		context   string
		verbose   = fs.Bool("v", false, "verbose")
		label     = fs.String("l", "", "select pod by label")
		container = fs.String("c", "", "container")
		sh        = fs.String("sh", "sh", "shell")
	)
	args := os.Args
	if len(args) > 1 && !strings.HasPrefix(args[1], "-") {
		context = args[1]
		args = args[2:]
	}
	fs.Parse(args)
	e, err := shell(context, *verbose, *label, *container, *sh, fs.Args())
	if err != nil {
		log.Printf("Error: %s", err)
	}
	os.Exit(e)

}

func shell(context string, verbose bool, label, container string, shell string, args []string) (int, error) {
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
			allArgs = append(allArgs, "-it", pod, shell)
			allArgs = append(allArgs, args...)
			return kc.Run(kc.PrepKC(context, allArgs...), verbose)
		}
		//no support for interactive mode with multiple-targets
		go func(pod string) {
			wg.Add(1)
			allArgs = append(allArgs, pod, shell)
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
