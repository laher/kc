package kc

import (
	"log"
	"os"
	"strconv"
)

func Logg(context string, verbose, follow bool, tail int, label string, args []string) (int, error) {
	var pod string
	if label != "" {
		pod = Pod(context, verbose, label)
	} else {
		if len(args) < 1 {
			log.Printf("Need an identifier for log")
			os.Exit(1)
		}
		pod = args[0]
		args = args[1:]
	}
	allArgs := []string{"log"}
	if follow {
		allArgs = append(allArgs, "-f")
	}
	if tail != -1 {
		allArgs = append(allArgs, "--tail", strconv.Itoa(tail))
	}
	allArgs = append(allArgs, pod)
	allArgs = append(allArgs, args...)
	return Run(PrepKC(context, allArgs...), verbose)
}
