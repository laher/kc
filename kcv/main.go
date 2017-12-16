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
	e, err := versions(context, *verbose)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	os.Exit(e)
}

func versions(context string, verbose bool) (int, error) {
	cmd := kc.PrepKC(context,
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
	return kc.Run(cmd, verbose)
}
