package kc

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func Contexts(args []string) ([]string, []string) {
	var (
		context string
	)
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		context = args[0]
		args = args[1:]
	}
	contexts := strings.Split(context, ",")

	return contexts, args
}

//PodsByLabel resolves pod names using a selector
func PodsByLabel(context string, verbose bool, label string) []string {
	gpC := PrepKC(context, "get", "pod", `-o=jsonpath={range .items[*]}{@.metadata.name}{"\n"}{end}`, "--selector", label)
	r, w := io.Pipe()
	gpC.Stdout = w
	scanner := bufio.NewScanner(r)
	names := []string{}
	go func() {
		ex, err := Run(gpC, verbose)
		if err != nil {
			log.Printf("Error fetching pod %s\n", err)
			os.Exit(ex)
		}
		w.Close()
	}()
	for scanner.Scan() {
		name := scanner.Text()
		if strings.TrimSpace(name) != "" {
			names = append(names, name)
		}
	}

	if err := scanner.Err(); err != nil {
		//done
		log.Printf("Error scanning output %s\n", err)
		os.Exit(1)
	}

	return names
}

func PrepKC(context string, args ...string) *exec.Cmd {
	allArgs := []string{"kubectl"}
	if context != "" {
		allArgs = append(allArgs, "--context", context)
	}
	return Prep(append(allArgs, args...)...)
}

func Prep(args ...string) *exec.Cmd {
	p, err := exec.LookPath(args[0])
	if err != nil {
		log.Printf("Couldn't find exe %s - %s", p, err)
	}
	cmd := exec.Command(args[0])
	cmd.Args = args
	// redirect output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// don't redirect stdin except for interactive sessions (exec -it). It can prevent Ctrl-C (quit) from working
	return cmd
}

func Run(cmd *exec.Cmd, verbose bool) (int, error) {
	if verbose {
		log.Printf("Running cmd: %s", cmd.Args)
	}

	err := cmd.Start()
	if err != nil {
		log.Printf("Launch error: %s", err)
		return 1, err
	}
	if verbose {
		log.Printf("Waiting for command to finish...")
	}
	err = cmd.Wait()
	if err != nil {
		if verbose {
			log.Printf("Command exited with error: %v", err)
		}
	} else {
		if verbose {
			log.Printf("Command completed without error")
		}
	}
	if err != nil {
		if e2, ok := err.(*exec.ExitError); ok { // there is error code
			processState, ok2 := e2.Sys().(syscall.WaitStatus)
			if ok2 {
				errcode := processState.ExitStatus()
				log.Printf("%s returned exit status: %d", cmd.Args[0], errcode)
				return errcode, err
			}
		}
		return 1, err
	}
	return 0, nil
}
