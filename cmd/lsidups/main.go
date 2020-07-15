package main

import (
	"flag"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/vitali-fedulov/images"
)

type Image struct {
	fp      string
	imgHash []float32
	imgSize image.Point
}

// checks if path has hidden elements, unix only
func isHidden(fp string) bool {
	for _, element := range strings.Split(fp, string(filepath.Separator)) {
		if strings.HasPrefix(element, ".") && element != "." {
			return true
		}
	}
	return false
}

func getFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(fp string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			abs, err := filepath.Abs(fp)
			if err == nil {
				files = append(files, abs)
			}
		}
		return nil
	})
	return files, err
}

func filterFiles(slice []string, condition func(string) bool) []string {
	var newSlice []string
	for _, element := range slice {
		if condition(element) {
			newSlice = append(newSlice, element)
		}
	}
	return newSlice
}

// takes file pathes list and list of extensions, retuns pathes with thouse extensions
func filterExt(files []string, searchExt []string) []string {
	return filterFiles(files, func(fp string) bool {
		for _, ext := range searchExt {
			fpext := strings.ToLower(filepath.Ext(fp))
			if strings.Contains(fpext, ext) {
				return true
			}
		}
		return false
	})
}

func hashImage(fp string) Image {
	pic, err := images.Open(fp)
	if err != nil {
		panic(err)
	}
	imgHash, imgSize := images.Hash(pic)
	return Image{fp, imgHash, imgSize}
}

func worker(jobs <-chan string, results chan<- Image, wg *sync.WaitGroup) {
	defer wg.Done()
	for fp := range jobs {
		results <- hashImage(fp)
	}
}

// Contains tells whether a contains x.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func main() {
	// start := time.Now()

	searchExt := []string{".jpg", ".jpeg", ".png", ".gif"}
	p := flag.String("p", ".", "find duplicate images in given directory")
	flag.Parse()

	files, _ := getFiles(*p)
	files = filterFiles(files, func(fp string) bool { return !isHidden(fp) })

	files = filterExt(files, searchExt)

	// fmt.Println("got files", time.Since(start))

	// start = time.Now()

	numJobs := len(files)
	jobs := make(chan string, numJobs)
	results := make(chan Image, numJobs)

	var wg sync.WaitGroup
	for w := 1; w <= runtime.NumCPU(); w++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	for _, fp := range files {
		jobs <- fp
	}
	close(jobs)

	wg.Wait()
	// yay, antipatterns!
	close(results)

	var pics []Image
	for pic := range results {
		pics = append(pics, pic)
	}

	// start := time.Now()
	var dups []string
	for _, pic := range pics {
		for _, ipic := range pics {
			if ipic.fp != pic.fp {
				if images.Similar(ipic.imgHash, pic.imgHash, ipic.imgSize, pic.imgSize) {
					if !Contains(dups, pic.fp) {
						dups = append(dups, pic.fp)
					}
					if !Contains(dups, ipic.fp) {
						dups = append(dups, ipic.fp)
					}
				}
			}
		}
	}

	// fmt.Println("search took", time.Since(start))
	for _, fp := range dups {
		fmt.Println(fp)
	}
}
