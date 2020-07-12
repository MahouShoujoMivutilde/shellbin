package main

import (
	"flag"
	"fmt"
	"math"
	"os"
)

func usage() {
	fmt.Printf("%[1]s - fit image into given rectangle while preserving aspect ratio\nOutputs new width and height as WxH.", os.Args[0])
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
