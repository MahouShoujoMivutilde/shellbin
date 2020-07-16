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

type dups struct {
	list []string
	m    sync.Mutex
}

func dupsSearch(pics <-chan Image, ipics *[]Image, duplicates *dups, wg *sync.WaitGroup) {
	defer wg.Done()
	for pic := range pics {
		for _, ipic := range *ipics {
			if ipic.fp != pic.fp {
				if images.Similar(ipic.imgHash, pic.imgHash, ipic.imgSize, pic.imgSize) {
					duplicates.m.Lock()
					if !in.ContainsStr(duplicates.list, pic.fp) {
						duplicates.list = append(duplicates.list, pic.fp)
					}
					if !in.ContainsStr(duplicates.list, ipic.fp) {
						duplicates.list = append(duplicates.list, ipic.fp)
					}
					duplicates.m.Unlock()
				}
			}
		}
	}
}

func main() {
	flag.Parse()

	start := time.Now()

	files, _ := filehelpers.GetFiles(input)
	files = filehelpers.FilterFiles(files, func(fp string) bool {
		return !filehelpers.IsHidden(fp)
	})

	// making sure it's image formats go supports
	files = filehelpers.FilterExt(files, searchExt)

	if verbose {
		fmt.Printf("> found %d images, took %s\n", len(files), time.Since(start))
	}

	start = time.Now()

	// calculating image similarity hashes
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
	// yay, antipatterns! (actually it's ok when you sure)
	close(results)

	if verbose {
		fmt.Printf("> processed images, took %s\n", time.Since(start))
	}

	var pics []Image
	for pic := range results {
		pics = append(pics, pic)
	}

	start = time.Now()

	// searching for similar images
	var duplicates dups
	picschan := make(chan Image, len(pics))

	for w := 1; w <= runtime.NumCPU(); w++ {
		wg.Add(1)
		go dupsSearch(picschan, &pics, &duplicates, &wg)
	}

	for _, pic := range pics {
		picschan <- pic
	}
	close(picschan)

	wg.Wait()

	for _, fp := range duplicates.list {
		fmt.Println(fp)
	}

	if verbose {
		fmt.Printf("> search of similar images took %s\n", time.Since(start))
	}
}
