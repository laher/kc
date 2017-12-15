package kc

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/alecthomas/kingpin"
)

var (
	getCommand       = kingpin.Command("g", "Get").Alias("get")
	getPodCommand    = Sub(getCommand.Command("p", "Get pods").Alias("pod").Alias("pods"))
	execCommand      = Sub(kingpin.Command("x", "Exec onto a pod").Alias("exec"))
	execContainer    = execCommand.command.Flag("co", "container").String()
	shCommand        = Sub(kingpin.Command("sh", "Shell (bash) onto a box").Alias("bash"))
	shContainer      = shCommand.command.Flag("co", "container").String()
	logCommand       = Sub(kingpin.Command("l", "log").Alias("log"))
	logFollow        = logCommand.command.Flag("follow", "follow").Short('f').Bool()
	logTail          = logCommand.command.Flag("tail", "tail").Short('t').Default("-1").Int()
	tailCommand      = Sub(kingpin.Command("t", "tail log").Alias("tail"))
	applyCommand     = Sub(kingpin.Command("a", "apply using a k8s file").Alias("apply"))
	applyFile        = applyCommand.command.Arg("file", "k8s file name").String()
	replaceCommand   = Sub(kingpin.Command("r", "replace using a k8s file").Alias("replace"))
	replaceFile      = replaceCommand.command.Arg("file", "k8s file name").String()
	bounceCommand    = Sub(kingpin.Command("b", "Bounce a deployment").Alias("bounce"))
	bounceDeployment = bounceCommand.command.Arg("deployment", "deployment name").String()
)

func init() {

	//	Remainder(getPodCommand, versionsCommand, execCommand, shCommand, logCommand, tailCommand, applyCommand, replaceCommand, bounceCommand)
}

/*
func Main() {
	ex, err := kc()
	if err != nil {
		log.Printf("Error: %v", err)
	}
	os.Exit(ex)
}

func Kc() (int, error) {
	switch kingpin.Parse() {
	case "g p": //not really important IMO
		return Run(PrepKC(getPodCommand, "get", "pod"), *getPodCommand.verbose)
	case "v":
		return Versions()
	case "x":
		return ExecPod()
	case "sh":
		return Shell()
	case "l":
		return Logg(logCommand, *logFollow, *logTail)
	case "t":
		return Logg(tailCommand, true, 1)
	case "a":
		return Apply()
	case "r":
		return Replace()
	case "b":
		return Bounce()
	default:
		return 1, errors.New("Unsupported subcommand")
	}
}
*/

type subcommand struct {
	command   *kingpin.CmdClause
	context   *string
	remainder *[]string
	verbose   *bool
}

func Remainder(cs ...*subcommand) {
	for _, c := range cs {
		c.remainder = RemainingArgs(c.command.Arg("remainder", "Remaining kubectl args"))
	}
}

func Sub(c *kingpin.CmdClause) *subcommand {
	s := &subcommand{
		command: c,
		context: c.Flag("context", "Context").Short('c').String(),
		verbose: c.Flag("verbose", "Show verbose loggng").Short('v').Bool(),
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

func RemainingArgs(s kingpin.Settings) *[]string {
	target := new([]string)
	s.SetValue((*rka)(target))
	return target
}

//resolve pod name using a selector if necessary
func Pod(c *subcommand, p string) string {
	switch {
	case strings.Contains(p, "="):
		gpC := PrepKC(*c.context, "get", "pod", "-o=name", "--selector", p)
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
		ex, err := Run(gpC, *c.verbose)
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

	// redirect IO
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
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
