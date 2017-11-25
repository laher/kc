package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/alecthomas/kingpin"
)

var (
	getCommand      = kingpin.Command("g", "Get").Alias("get")
	getPodCommand   = sub(getCommand.Command("p", "Get pods").Alias("pod").Alias("pods"))
	versionsCommand = sub(kingpin.Command("v", "Get pod versions").Alias("version"))
	execCommand     = sub(kingpin.Command("x", "Exec onto a pod").Alias("exec"))
	execContainer   = execCommand.command.Flag("co", "container").String()
	shCommand       = sub(kingpin.Command("sh", "Shell (bash) onto a box").Alias("bash"))
	shContainer     = shCommand.command.Flag("co", "container").String()
	logCommand      = sub(kingpin.Command("l", "log").Alias("log"))
	logFollow       = logCommand.command.Flag("follow", "follow").Short('f').Bool()
	logTail         = logCommand.command.Flag("tail", "tail").Short('t').Default("-1").Int()
	tailCommand     = sub(kingpin.Command("t", "tail log").Alias("tail"))
)

type subcommand struct {
	command   *kingpin.CmdClause
	context   *string
	remainder *[]string
	verbose   *bool
}

func sub(c *kingpin.CmdClause) subcommand {
	s := subcommand{
		command:   c,
		context:   c.Flag("context", "Context").Short('c').String(),
		remainder: remainingArgs(c.Arg("remainder", "Remaining kubectl args")),
		verbose:   c.Flag("verbose", "Show verbose loggng").Short('v').Bool(),
	}
	return s
}

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
		return run(prepKC(getPodCommand, "get", "pod"), *getPodCommand.verbose)

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

func pod(c subcommand, p string) string {
	switch {
	case strings.Contains(p, "="):
		gpC := prepKC(c, "get", "pod", "-o=name", "--selector", p)
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
		ex, err := run(gpC, *c.verbose)
		if err != nil {
			log.Printf("Error fetching pod %s\n", err)
			os.Exit(ex)
		}
		if strings.HasPrefix(name, "pod/") {
			return name[4:]
		}
		return name
	default:
		return p
	}
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

func run(cmd *exec.Cmd, verbose bool) (int, error) {
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
