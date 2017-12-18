package main

import (
	"flag"
	"log"
	"os"

	kc "github.com/laher/kc/internal"
)

func main() {
	var (
		fs      = flag.NewFlagSet("kcr", flag.ExitOnError)
		verbose = fs.Bool("v", false, "verbose")
		file    = fs.String("f", "", "file to apply")
	)
	contexts, args := kc.Contexts(os.Args[1:])
	fs.Parse(args)
	for _, context := range contexts {
		if len(contexts) > 1 || *verbose {
			log.Printf("context: %s", context)
		}

		e, err := replace(context, *verbose, *file, fs.Args())
		if err != nil {
			log.Printf("Error: %s", err)
			os.Exit(e)
		}
	}
	if *verbose {
		log.Print("done")
	}
}

func replace(context string, verbose bool, file string, args []string) (int, error) {
	args = append([]string{"replace", "--cascade", "--force", "-f", file}, args...)
	return kc.Run(kc.PrepKC(context, args...), verbose)
}
