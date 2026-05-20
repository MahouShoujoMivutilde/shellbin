package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"time"
)

var DESC string = os.Args[0] + `
  Measures reaction time.

  Launch the program. Wait for the GO! sign. Press enter as fast as you can.

`

func usage() {
	fmt.Fprint(flag.CommandLine.Output(), DESC)

	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %[1]s:\n", os.Args[0])

	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	// grace period
	time.Sleep(5 * time.Second)

	// wait 0.0 to 5.0 seconds
	time.Sleep(time.Duration(rand.Float64()) * 5 * time.Second)

	fmt.Println("GO!")

	start := time.Now()

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(time.Since(start))
}
