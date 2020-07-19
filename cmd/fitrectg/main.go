package main

import (
	"flag"
	"fmt"
	"math"
	"os"
)

var DESC string = os.Args[0] + `

  Fits some rectangle (meant for images) into given rectangle while
  preserving aspect ratio. Outputs new width and height as WxH.

`

var EXAMPLES string = `
Examples:
  calculate new dimensions of image with width = 3600 and height = 2404
  to fit into 456x490 rectangle
    fitrectg -w 3600 -h 2404 -fw 456 -fh 490
`

func usage() {
	fmt.Fprint(flag.CommandLine.Output(), DESC)

	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %[1]s:\n", os.Args[0])

	flag.PrintDefaults()
	fmt.Fprint(flag.CommandLine.Output(), EXAMPLES)
}

func fitIntoRect(W float64, H float64, rectW float64, rectH float64) (int, int) {
	s1 := rectW / W
	s2 := rectH / H
	if s1 < s2 {
		W = W * s1
		H = H * s1
	} else {
		W = W * s2
		H = H * s2
	}

	return int(math.Round(W)), int(math.Round(H))
}

func main() {
	flag.Usage = usage
	W := flag.Float64("w", 0, "current image width")
	H := flag.Float64("h", 0, "current image height")
	rectW := flag.Float64("fw", 0, "rectangle width")
	rectH := flag.Float64("fh", 0, "rectangle height")

	flag.Parse()

	if *W == 0 || *H == 0 || *rectH == 0 || *rectW == 0 {
		flag.Usage()
		os.Exit(1)
	}

	outW, outH := fitIntoRect(*W, *H, *rectW, *rectH)
	fmt.Printf("%dx%d\n", outW, outH)
}
