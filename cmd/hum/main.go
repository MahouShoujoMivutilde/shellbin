package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

var DESC string = os.Args[0] + `

  Is a tool for humanizing various things (now supports time)

  echo thing | hum (time) [-t]

`

func usage() {
	fmt.Fprint(flag.CommandLine.Output(), DESC)

	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %[1]s:\n", os.Args[0])

	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage

	timeFormat := flag.String("t",
		"2006-01-02 15:04:05.999999999 -0700",
		// defaults to stat human readable format
		"time format, fill with Mon Jan 2 15:04:05 MST 2006,"+
			" see https://golang.org/src/time/format.go")

	flag.Parse()

	mode := flag.Arg(0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		str := strings.TrimSpace(scanner.Text())

		switch mode {
		case "time":
			t, err := time.Parse(*timeFormat, str)
			if err != nil {
				panic(err)
			}

			fmt.Println(humanize.Time(t))
		}
	}
}
