package kc

import (
	"log"
	"os"
	"strconv"
)

func Logg(context string, verbose, follow bool, tail int, args []string) (int, error) {
	if len(args) < 1 {
		log.Printf("Need an identifier for log")
		os.Exit(1)
	}
	p := (args)[0]
	allArgs := []string{"log"}
	if follow {
		allArgs = append(allArgs, "-f")
	}
	if tail != -1 {
		allArgs = append(allArgs, "--tail", strconv.Itoa(tail))
	}
	pod := Pod(context, verbose, p)
	allArgs = append(allArgs, pod)
	allArgs = append(allArgs, args[1:]...)
	return Run(PrepKC(context, allArgs...), verbose)
}
