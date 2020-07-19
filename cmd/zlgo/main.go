package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/kortschak/zalgo"
)

var DESC string = os.Args[0] + `

  Is a tool for generating zalgo text

`

var EXAMPLES string = `
Examples:
  zalgofy stdio
    echo 'some text....' | zlgo
`

func usage() {
	fmt.Fprint(flag.CommandLine.Output(), DESC)

	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %[1]s:\n", os.Args[0])

	flag.PrintDefaults()
	fmt.Fprint(flag.CommandLine.Output(), EXAMPLES)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		str := strings.TrimSpace(scanner.Text())
		fmt.Fprintln(zalgo.Corruption, str)
	}
}
