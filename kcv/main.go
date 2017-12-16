package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
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

		e, err := versions(context, *verbose, fs.Args())
		if err != nil {
			log.Printf("Error: %s", err)
			os.Exit(e)
		}
	}
	if *verbose {
		log.Print("done")
	}
}

func versions(context string, verbose bool, args []string) (int, error) {
	kcArgs := []string{"get", "pod", "-o", "jsonpath='{range .items[*].spec.containers[*]}{.name}{\"\\t\"}{.image}{\"\\n\"}{end}'"}
	kcArgs = append(kcArgs, args...)
	cmd := kc.PrepKC(context, kcArgs...)
	r, w := io.Pipe()
	cmd.Stdout = w
	br := bufio.NewReader(r)
	go func() {
		existing := map[string]struct{}{}
		for {
			b, _, err := br.ReadLine()
			if err != nil {
				//done
				return
			}
			s := string(b)
			parts := strings.Split(s, "\t")

			name := parts[0]
			image := ""
			if len(parts) > 1 {
				parts2 := strings.Split(parts[1], "/")
				if len(parts2) > 1 {
					image = parts2[1]
				} else {
					image = parts2[0]

				}
			}
			record := fmt.Sprintf("%s\t%s", name, image)
			_, exists := existing[record]
			if !exists {
				fmt.Println(record)
				existing[record] = struct{}{}
			}
		}
	}()
	return kc.Run(cmd, verbose)
}
