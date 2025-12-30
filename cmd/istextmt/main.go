package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/gabriel-vasile/mimetype"
)

var DESC string = os.Args[0] + `

  Checks if file is a text file, and
    if it is
      prints filepath
    else
      doesn't print anything

  Designed to be used as filter for fd,
  it is also much faster than
    file --mime-type -b file.txt + case text/*...
  ...shenanigans.

`

var EXAMPLES string = `
Examples:
  check file is a text file
    echo "file.txt" | istextmt
    if it is, it will print it's name, otherwise - skip

  find only only text files with fd
    fd -t f | istextmt
`

var textPatt = []string{
	"text/", "/xml", "application/json",
	"application/postscript", "application/rss+xml",
	"application/atom+xml", "application/javascript",
	"application/x-python",
}

func usage() {
	fmt.Fprint(flag.CommandLine.Output(), DESC)

	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %[1]s:\n", os.Args[0])

	flag.PrintDefaults()
	fmt.Fprint(flag.CommandLine.Output(), EXAMPLES)
}

func getTextFiles(filesIn <-chan string, textFilesOut chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for file := range filesIn {
		if fpAbs, err := filepath.Abs(file); err == nil {
			mime, _ := mimetype.DetectFile(fpAbs)
			for _, pattern := range textPatt {
				if strings.Contains(strings.ToLower(mime.String()), pattern) {
					textFilesOut <- fpAbs
					break
				}
			}
		}
	}
}

func getMimes(filesIn <-chan string) <-chan string {
	var wg sync.WaitGroup

	textFilesOut := make(chan string, 8192)

	for w := 1; w <= runtime.NumCPU(); w++ {
		wg.Add(1)
		go getTextFiles(filesIn, textFilesOut, &wg)
	}

	go func() {
		wg.Wait()
		close(textFilesOut)
	}()

	return textFilesOut
}

func main() {
	flag.Usage = usage
	flag.Parse()

	filesIn := make(chan string, 8192)

	// setup processing
	textFilesOut := getMimes(filesIn)

	// async get input
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			file := strings.TrimSpace(scanner.Text())
			filesIn <- file
		}
		close(filesIn)
	}()

	// print only the files that are text by mime
	for fp := range textFilesOut {
		fmt.Println(fp)
	}
}
