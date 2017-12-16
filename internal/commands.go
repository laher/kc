package kc

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

func Logg(context string, verbose, follow bool, tail int, label string, args []string, wg *sync.WaitGroup) (int, error) {
	var pods []string
	if label != "" {
		pods = PodsByLabel(context, verbose, label)
	} else {
		if len(args) < 1 {
			log.Printf("Need an identifier for log")
			os.Exit(1)
		}
		pods = []string{args[0]}
		args = args[1:]
	}
	for _, pod := range pods {
		allArgs := []string{"log"}
		if follow {
			allArgs = append(allArgs, "-f")
		}
		if tail != -1 {
			allArgs = append(allArgs, "--tail", strconv.Itoa(tail))
		}
		allArgs = append(allArgs, pod)
		allArgs = append(allArgs, args...)
		if len(pods) == 1 {
			return Run(PrepKC(context, allArgs...), verbose)
		}
		wg.Add(1)
		go func(pod string) {
			cmd := PrepKC(context, allArgs...)
			cmd.Stdout = nil
			p, err := cmd.StdoutPipe()
			if err != nil {
				log.Printf("Stdout error: %s", err)
			}
			go func() {
				b := bufio.NewScanner(p)
				for b.Scan() {
					os.Stdout.WriteString(fmt.Sprintf("[%s]: %s\n", pod, b.Text()))
				}
			}()
			Run(cmd, verbose)
			wg.Done()
		}(pod)
	}
	return 0, nil
}
