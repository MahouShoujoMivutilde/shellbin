package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kortschak/zalgo"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		str := strings.TrimSpace(scanner.Text())
		fmt.Fprintln(zalgo.Corruption, str)
	}
}
