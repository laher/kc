package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func execPod() (int, error) {
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
	pod := pod(execCommand, p)
	args = append(args, pod)
	return run(prepKC(execCommand,
		append(args, r...)...), *execCommand.verbose)
}

func logg(c *subcommand, follow bool, tail int) (int, error) {
	if len(*c.remainder) < 1 {
		log.Printf("Need an identifier for log")
		os.Exit(1)
	}
	p := (*c.remainder)[0]
	allArgs := []string{"log"}
	if follow {
		allArgs = append(allArgs, "-f")
	}
	if tail != -1 {
		allArgs = append(allArgs, "--tail", strconv.Itoa(tail))
	}
	pod := pod(c, p)
	allArgs = append(allArgs, pod)
	args := (*c.remainder)[1:]
	allArgs = append(allArgs, args...)
	return run(prepKC(c, allArgs...), *c.verbose)
}

func bounce() (int, error) {
	args := []string{"scale", "--current-replicas=1", "--replicas=0", "deploy", *bounceDeployment}
	ex, err := run(prepKC(bounceCommand, args...), *bounceCommand.verbose)
	if err != nil {
		return ex, err
	}
	args = []string{"scale", "--current-replicas=0", "--replicas=1", "deploy", *bounceDeployment}
	return run(prepKC(bounceCommand, args...), *bounceCommand.verbose)
}

func apply() (int, error) {
	args := []string{"apply", "-f", *applyFile}
	return run(prepKC(applyCommand, args...), *applyCommand.verbose)
}

func replace() (int, error) {
	args := []string{"replace", "--cascade", "--force", "-f", *replaceFile}
	return run(prepKC(replaceCommand, args...), *replaceCommand.verbose)
}

func shell() (int, error) {
	if len(*shCommand.remainder) < 1 {
		log.Printf("Need an identifier for exec")
		os.Exit(1)
	}
	args := []string{"exec", "-it"}
	if *execContainer != "" {
		args = append(args, "-c", *execContainer)
	}
	p := (*shCommand.remainder)[0]
	pod := pod(shCommand, p)
	args = append(args, pod, "sh")
	args = append(args, (*shCommand.remainder)[1:]...)
	return run(prepKC(shCommand, args...), *shCommand.verbose)
}

func versions() (int, error) {
	cmd := prepKC(versionsCommand,
		"get", "pod", "-o", "jsonpath='{range .items[*].spec.containers[*]}{.name}{\"\\t\"}{.image}{\"\\n\"}{end}'")
	r, w := io.Pipe()
	cmd.Stdout = w
	br := bufio.NewReader(r)
	go func() {
		for {
			b, _, err := br.ReadLine()
			if err != nil {
				//done
				return
			}
			s := string(b)
			parts := strings.Split(s, "\t")
			fmt.Print(parts[0], "\t")
			if len(parts) > 1 {
				parts2 := strings.Split(parts[1], "/")
				if len(parts2) > 1 {
					fmt.Println(parts2[1])
				} else {

					fmt.Println(parts2[0])
				}
			}
		}
	}()
	return run(cmd, *versionsCommand.verbose)
}
