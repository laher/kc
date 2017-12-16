package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/laher/kc/internal"
)

func main() {
	var (
		fs      = flag.NewFlagSet("kcv", flag.ExitOnError)
		context string
		verbose = fs.Bool("v", false, "verbose")
	)
	args := os.Args
	if len(args) > 1 && !strings.HasPrefix(args[1], "-") {
		context = args[1]
		args = args[2:]
	}
	fs.Parse(args)
	contexts := strings.Split(context, ",")
	for _, context := range contexts {
		if len(contexts) > 1 || *verbose {
			log.Printf("context: %s", context)
		}
		e, err := getpods(context, *verbose, fs.Args())
		if err != nil {
			log.Printf("Error: %s", err)
			os.Exit(e)
		}
	}
	if *verbose {
		log.Print("done")
	}
}

func getpods(context string, verbose bool, args []string) (int, error) {
	kcArgs := []string{"get", "pod"}
	kcArgs = append(kcArgs, args...)
	cmd := kc.PrepKC(context, kcArgs...)
	return kc.Run(cmd, verbose)
}
