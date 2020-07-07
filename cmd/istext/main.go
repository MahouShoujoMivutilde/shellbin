package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

func usage() {
	fmt.Printf("%[1]s - check if file is a text file,\n" +
				"and if it is - print filepath\n" +
				"Usage of %[1]s:\n" +
				"    %[1]s filename.txt\n", os.Args[0])
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
