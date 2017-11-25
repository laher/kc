package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/alecthomas/kingpin"
)

type subcommand struct {
	command   *kingpin.CmdClause
	context   *string
	remainder *[]string
}

func sub(c *kingpin.CmdClause) subcommand {
	s := subcommand{
		command:   c,
		context:   c.Flag("context", "Context").Short('c').String(),
		remainder: remainingArgs(c.Arg("remainder", "Remaining kubectl args")),
	}
	return s
}

var (
	getCommand      = kingpin.Command("g", "Get").Alias("get")
	getPodCommand   = sub(getCommand.Command("p", "pods").Alias("pod").Alias("pods"))
	versionsCommand = sub(kingpin.Command("v", "Get pod versions").Alias("version"))
	execCommand     = sub(kingpin.Command("x", "exec on a pod").Alias("exec"))
	shCommand       = sub(kingpin.Command("sh", "Shell (bash) onto a box").Alias("bash"))
	logCommand      = sub(kingpin.Command("l", "log").Alias("log"))
	logFollow       = logCommand.command.Flag("follow", "follow").Short('f').Bool()
	logTail         = logCommand.command.Flag("tail", "tail").Short('t').Default("-1").Int()
	tailCommand     = sub(kingpin.Command("t", "tail log").Alias("tail"))
)

type rka []string

func (i *rka) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *rka) String() string {
	return ""
}

func (i *rka) IsCumulative() bool {
	return true
}

func remainingArgs(s kingpin.Settings) *[]string {
	target := new([]string)
	s.SetValue((*rka)(target))
	return target
}

func main() {
	ex, err := kc()
	if err != nil {
		log.Printf("Error: %v", err)
	}
	os.Exit(ex)
}
func kc() (int, error) {
	switch kingpin.Parse() {
	case "g p":
		return run(prepKC(getPodCommand, "get", "pod"))

	case "v":
		return versions()

	case "x":
		return execPod()

	case "sh":
		return shell()

	case "l":
		return logg(logCommand, *logFollow, *logTail)

	case "t":
		return logg(tailCommand, true, 1)
	default:
		return 1, errors.New("Unsupported function")
	}
}

func execPod() (int, error) {
	if len(*execCommand.remainder) < 1 {
		log.Printf("Need an identifier for exec")
		os.Exit(1)
	}
	return run(prepKC(execCommand,
		append([]string{"exec", "-it"}, *execCommand.remainder...)...))
}

func pod(cmd subcommand, p string) string {
	switch {
	case strings.Contains(p, "="):
		gpC := prepKC(cmd, "get", "pod", "-o=name", "--selector", p)
		r, w := io.Pipe()
		gpC.Stdout = w
		br := bufio.NewReader(r)
		name := ""
		go func() {
			for {
				b, _, err := br.ReadLine()
				if err != nil {
					//done
					return
				}
				name = string(b)
			}
		}()
		ex, err := run(gpC)
		if err != nil {
			log.Printf("Error fetching pod %s\n", err)
			os.Exit(ex)
		}
		return name
	default:
		return p
	}
}

func logg(c subcommand, follow bool, tail int) (int, error) {
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
	return run(prepKC(c, allArgs...))
}

func shell() (int, error) {
	if len(*shCommand.remainder) < 1 {
		log.Printf("Need an identifier for exec")
		os.Exit(1)
	}
	p := (*shCommand.remainder)[0]
	pod := ""
	args := (*shCommand.remainder)[1:]
	switch {
	case strings.HasPrefix(p, "d/"):
		pod = p[2:]
	default:
		pod = p
	}
	return run(prepKC(shCommand, append([]string{"exec", "-it", pod, "bash"}, args...)...))
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
	return run(cmd)
}

func prepKC(sc subcommand, args ...string) *exec.Cmd {
	allArgs := []string{"kubectl"}
	if *sc.context != "" {
		allArgs = append(allArgs, "--context", *sc.context)
	}
	return prep(append(allArgs, args...)...)
}

func prep(args ...string) *exec.Cmd {
	p, err := exec.LookPath(args[0])
	if err != nil {
		log.Printf("Couldn't find exe %s - %s", p, err)
	}
	cmd := exec.Command(args[0])

	cmd.Args = args

	// redirect IO
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd
}

var verbose = false

func run(cmd *exec.Cmd) (int, error) {
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
