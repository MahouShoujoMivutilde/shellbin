package main

import (
	"flag"
	"fmt"
	"image"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/MahouShoujoMivutilde/shellbin/internal/filehelpers"
	"github.com/MahouShoujoMivutilde/shellbin/internal/in"
	"github.com/vitali-fedulov/images"
)

type extensions []string

func (e *extensions) String() string {
	return strings.Join(*e, ",")
}

func (e *extensions) Set(val string) error {
	*e = extensions{}
	for _, ext := range strings.Split(val, ",") {
		*e = append(*e, ext)
	}
	return nil
}

var searchExt extensions
var input string
var verbose bool

func init() {
	searchExt = extensions{".jpg", ".jpeg", ".png", ".gif"}
	flag.Var(&searchExt, "e", "image extensions (with dots) to look for")
	flag.StringVar(&input, "p", ".", "find duplicate images in given directory")
	flag.BoolVar(&verbose, "v", false, "show time it took to complete key parts of the search")
}

type Image struct {
	fp      string
	imgHash []float32
	imgSize image.Point
}

func makeImage(fp string) Image {
	pic, err := images.Open(fp)
	if err != nil {
		panic(err)
	}
	imgHash, imgSize := images.Hash(pic)
	return Image{fp, imgHash, imgSize}
}

func imageMaker(jobs <-chan string, results chan<- Image, wg *sync.WaitGroup) {
	defer wg.Done()
	for fp := range jobs {
		results <- makeImage(fp)
	}
}

func main() {
	flag.Parse()

	start := time.Now()

	files, _ := filehelpers.GetFiles(input)
	files = filehelpers.FilterFiles(files, func(fp string) bool {
		return !filehelpers.IsHidden(fp)
	})

	files = filehelpers.FilterExt(files, searchExt)

	if verbose {
		fmt.Printf("> found %d images, took %s\n", len(files), time.Since(start))
	}

	start = time.Now()

	numJobs := len(files)
	jobs := make(chan string, numJobs)
	results := make(chan Image, numJobs)

	var wg sync.WaitGroup
	for w := 1; w <= runtime.NumCPU(); w++ {
		wg.Add(1)
		go imageMaker(jobs, results, &wg)
	}

	for _, fp := range files {
		jobs <- fp
	}
	close(jobs)

	wg.Wait()
	// yay, antipatterns!
	close(results)

	if verbose {
		fmt.Printf("> processed image parameters, took %s\n", time.Since(start))
	}

	var pics []Image
	for pic := range results {
		pics = append(pics, pic)
	}

	start = time.Now()
	var dups []string
	for _, pic := range pics {
		for _, ipic := range pics {
			if ipic.fp != pic.fp {
				if images.Similar(ipic.imgHash, pic.imgHash, ipic.imgSize, pic.imgSize) {
					if !in.ContainsStr(dups, pic.fp) {
						dups = append(dups, pic.fp)
					}
					if !in.ContainsStr(dups, ipic.fp) {
						dups = append(dups, ipic.fp)
					}
				}
			}
		}
	}

	for _, fp := range dups {
		fmt.Println(fp)
	}

	if verbose {
		fmt.Printf("> search of similar images took %s\n", time.Since(start))
	}
}
