package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func main() {
	unescPtr := flag.Bool("u", false, "unescape uri instead")
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
