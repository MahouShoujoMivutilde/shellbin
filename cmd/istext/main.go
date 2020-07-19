package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

var DESC string = os.Args[0] + `

  Checks if file is a text file, and
    if it is
      prints filepath and exits with 0,
    else
      exits with 1

  Designed to be used as filter for fd,
  it is also much faster than
    file --mime-type -b file.txt + case text/*...
  ...shenanigans.

`

var EXAMPLES string = `
Examples:
  check file is a text file
    istext file.txt && echo 'this is text file' || echo 'this is not text'

  find only only text files with fd
    fd -t f -x istext {}
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
	input := flag.Arg(0)

	if input == "" {
		flag.Usage()
		os.Exit(1)
	}

	mime, _ := mimetype.DetectFile(input)

	for _, pattern := range []string{"text/", "/xml", "application/json"} {
		if strings.Contains(strings.ToLower(mime.String()), pattern) {
			fmt.Println(input)
			os.Exit(0)
		}
	}

	os.Exit(1)
}
