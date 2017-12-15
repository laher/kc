package kc

import (
	"log"
	"os"
	"strconv"
)

func ExecPod() (int, error) {
	if len(*execCommand.remainder) < 1 {
		log.Printf("Need an identifier for exec")
		os.Exit(1)
	}
	args := []string{"exec", "-it"}
	if *execContainer != "" {
		args = append(args, "-c", *execContainer)
	}
	p := (*execCommand.remainder)[0]
	r := (*execCommand.remainder)[1:]
	pod := Pod(*execCommand.context, *execCommand.verbose, p)
	args = append(args, pod)
	return Run(PrepKC(*execCommand.context,
		append(args, r...)...), *execCommand.verbose)
}

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

func Bounce() (int, error) {
	args := []string{"scale", "--current-replicas=1", "--replicas=0", "deploy", *bounceDeployment}
	ex, err := Run(PrepKC(*bounceCommand.context, args...), *bounceCommand.verbose)
	if err != nil {
		return ex, err
	}
	args = []string{"scale", "--current-replicas=0", "--replicas=1", "deploy", *bounceDeployment}
	return Run(PrepKC(*bounceCommand.context, args...), *bounceCommand.verbose)
}

func Apply() (int, error) {
	args := []string{"apply", "-f", *applyFile}
	return Run(PrepKC(*applyCommand.context, args...), *applyCommand.verbose)
}

func Replace() (int, error) {
	args := []string{"replace", "--cascade", "--force", "-f", *replaceFile}
	return Run(PrepKC(*replaceCommand.context, args...), *replaceCommand.verbose)
}

func Shell() (int, error) {
	if len(*shCommand.remainder) < 1 {
		log.Printf("Need an identifier for exec")
		os.Exit(1)
	}
	args := []string{"exec", "-it"}
	if *execContainer != "" {
		args = append(args, "-c", *execContainer)
	}
	p := (*shCommand.remainder)[0]
	pod := Pod(*shCommand.context, *shCommand.verbose, p)
	args = append(args, pod, "sh")
	args = append(args, (*shCommand.remainder)[1:]...)
	return Run(PrepKC(*shCommand.context, args...), *shCommand.verbose)
}
