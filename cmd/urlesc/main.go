package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

var DESC string = os.Args[0] + `

  Is a tool for escaping file path to make it safe for urls

`

var EXAMPLES string = `
Examples:
  escape path
    echo 'some/path/внезапно!@@"/dir' | urlesc
`

func usage() {
	fmt.Fprint(flag.CommandLine.Output(), DESC)

	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %[1]s:\n", os.Args[0])

	flag.PrintDefaults()
	fmt.Fprint(flag.CommandLine.Output(), EXAMPLES)
}

func main() {
	unescPtr := flag.Bool("u", false, "unescape uri instead")
	flag.Usage = usage
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		str := strings.TrimSpace(scanner.Text())
		u, err := url.Parse(str)
		if err == nil {
			if *unescPtr {
				str, err := url.PathUnescape(str)
				if err == nil {
					fmt.Println(str)
				}
			} else {
				fmt.Println(u.EscapedPath())
			}
		}
	}
}
