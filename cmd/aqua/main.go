package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	_ "golang.org/x/image/webp"

	"github.com/corona10/goimagehash"
	"github.com/vitali-fedulov/images"
)

func tell_webp_iccp(err error) {
	if fmt.Sprint(err) == "webp: invalid format" {
		fmt.Printf("A webp format error?\n" +
			"Check if your image has an ICC profile, and if it does - this could be related to\n" +
			"   https://github.com/golang/go/issues/60437#issuecomment-1563939784\n\n" +
			"(It'll probably work fine if you convert it to png.)\n\n")
	}
}

func main() {
	var fp1, fp2 string
	flag.StringVar(&fp1, "a", "", "file a")
	flag.StringVar(&fp2, "b", "", "file b")

	flag.Parse()

	if fp1 == "" || fp2 == "" {
		os.Exit(1)
	}

	file1, err := os.Open(fp1)
	if err != nil {
		panic(err)
	} else {
		defer file1.Close()
	}

	file2, err := os.Open(fp2)
	if err != nil {
		panic(err)
	} else {
		defer file1.Close()
	}

	img1, _, err := image.Decode(file1)
	if err != nil {
		tell_webp_iccp(err)
		panic(err)
	}

	img2, _, err := image.Decode(file2)
	if err != nil {
		tell_webp_iccp(err)
		panic(err)
	}

	hash1, _ := goimagehash.AverageHash(img1)
	hash2, _ := goimagehash.AverageHash(img2)
	distance, _ := hash1.Distance(hash2)

	fmt.Printf("average: %v\n", distance)

	hash1, _ = goimagehash.DifferenceHash(img1)
	hash2, _ = goimagehash.DifferenceHash(img2)
	distance, _ = hash1.Distance(hash2)

	fmt.Printf("diff: %v\n", distance)

	hash1, _ = goimagehash.PerceptionHash(img1)
	hash2, _ = goimagehash.PerceptionHash(img2)
	distance, _ = hash1.Distance(hash2)

	fmt.Printf("phash: %v\n", distance)

	width, height := 16, 16
	hash3, _ := goimagehash.ExtPerceptionHash(img1, width, height)
	hash4, _ := goimagehash.ExtPerceptionHash(img2, width, height)
	distance, _ = hash3.Distance(hash4)

	fmt.Printf("phash big: %v\n", distance)
	fmt.Printf(" a size: %v\n", hash3.Bits())
	fmt.Printf(" b size: %v\n", hash4.Bits())

	hashA, imgSizeA := images.Hash(img1)
	hashB, imgSizeB := images.Hash(img2)

	if images.Similar(hashA, hashB, imgSizeA, imgSizeB) {
		fmt.Println("Images are similar.")
	} else {
		fmt.Println("Images are distinct.")
	}
}
