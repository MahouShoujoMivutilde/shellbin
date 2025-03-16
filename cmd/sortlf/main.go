package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/djherbis/times"
)

// TAKEN FROM LF: START

const (
	ignorecase bool = true
	ignoredia  bool = true
)

var normMap map[rune]rune

func init() {
	normMap = make(map[rune]rune)

	// (not only) european
	appendTransliterate(
		"ěřůøĉĝĥĵŝŭèùÿėįųāēīūļķņģőűëïąćęłńśźżõșțčďĺľňŕšťýžéíñóúüåäöçîşûğăâđêôơưáàãảạ",
		"eruocghjsueuyeiuaeiulkngoueiacelnszzostcdllnrstyzeinouuaaocisugaadeoouaaaaa",
	)

	// Vietnamese
	appendTransliterate(
		"áạàảãăắặằẳẵâấậầẩẫéẹèẻẽêếệềểễiíịìỉĩoóọòỏõôốộồổỗơớợờởỡúụùủũưứựừửữyýỵỳỷỹđ",
		"aaaaaaaaaaaaaaaaaeeeeeeeeeeeiiiiiioooooooooooooooooouuuuuuuuuuuyyyyyyd",
	)
}

func appendTransliterate(base, norm string) {
	normRunes := []rune(norm)
	baseRunes := []rune(base)

	lenNorm := len(normRunes)
	lenBase := len(baseRunes)
	if lenNorm != lenBase {
		panic(
			"Base and normalized strings have differend length: base=" + strconv.Itoa(
				lenBase,
			) + ", norm=" + strconv.Itoa(
				lenNorm,
			),
		) // programmer error in constant length
	}

	for i := 0; i < lenBase; i++ {
		normMap[baseRunes[i]] = normRunes[i]

		baseUpper := unicode.ToUpper(baseRunes[i])
		normUpper := unicode.ToUpper(normRunes[i])

		normMap[baseUpper] = normUpper
	}
}

// Remove diacritics and make lowercase.
func removeDiacritics(baseString string) string {
	var normalizedRunes []rune
	for _, baseRune := range baseString {
		if normRune, ok := normMap[baseRune]; ok {
			normalizedRunes = append(normalizedRunes, normRune)
		} else {
			normalizedRunes = append(normalizedRunes, baseRune)
		}
	}
	return string(normalizedRunes)
}

func normalize(s1, s2 string) (string, string) {
	if ignorecase {
		s1 = strings.ToLower(s1)
		s2 = strings.ToLower(s2)
	}
	if ignoredia {
		s1 = removeDiacritics(s1)
		s2 = removeDiacritics(s2)
	}
	return s1, s2
}

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

// This function compares two strings for natural sorting which takes into
// account values of numbers in strings. For example, '2' is less than '10',
// and similarly 'foo2bar' is less than 'foo10bar', but 'bar2bar' is greater
// than 'foo10bar'.
func naturalLess(s1, s2 string) bool {
	lo1, lo2, hi1, hi2 := 0, 0, 0, 0
	for {
		if hi1 >= len(s1) {
			return hi2 != len(s2)
		}

		if hi2 >= len(s2) {
			return false
		}

		isDigit1 := isDigit(s1[hi1])
		isDigit2 := isDigit(s2[hi2])

		for lo1 = hi1; hi1 < len(s1) && isDigit(s1[hi1]) == isDigit1; hi1++ {
		}

		for lo2 = hi2; hi2 < len(s2) && isDigit(s2[hi2]) == isDigit2; hi2++ {
		}

		if s1[lo1:hi1] == s2[lo2:hi2] {
			continue
		}

		if isDigit1 && isDigit2 {
			num1, err1 := strconv.Atoi(s1[lo1:hi1])
			num2, err2 := strconv.Atoi(s2[lo2:hi2])

			if err1 == nil && err2 == nil {
				return num1 < num2
			}
		}

		return s1[lo1:hi1] < s2[lo2:hi2]
	}
}

type linkState byte

const (
	notLink linkState = iota
	working
	broken
)

type file struct {
	os.FileInfo
	linkState  linkState
	linkTarget string
	path       string
	dirCount   int
	accessTime time.Time
	changeTime time.Time
	ext        string
}

func readdir(path string) ([]*file, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()

	files := make([]*file, 0, len(names))
	for _, fname := range names {
		fpath := filepath.Join(path, fname)

		lstat, err := os.Lstat(fpath)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			continue
		}

		var linkState linkState
		var linkTarget string

		if lstat.Mode()&os.ModeSymlink != 0 {
			stat, err := os.Stat(fpath)
			if err == nil {
				linkState = working
				lstat = stat
			} else {
				linkState = broken
			}
			linkTarget, _ = os.Readlink(fpath)
		}

		ts := times.Get(lstat)
		at := ts.AccessTime()
		var ct time.Time
		// from times docs: ChangeTime() panics unless HasChangeTime() is true
		if ts.HasChangeTime() {
			ct = ts.ChangeTime()
		} else {
			// fall back to ModTime if ChangeTime cannot be determined
			ct = lstat.ModTime()
		}

		// returns an empty string if extension could not be determined
		// i.e. directories, filenames without extensions
		ext := filepath.Ext(fpath)

		files = append(files, &file{
			FileInfo:   lstat,
			linkState:  linkState,
			linkTarget: linkTarget,
			path:       fpath,
			dirCount:   -1,
			accessTime: at,
			changeTime: ct,
			ext:        ext,
		})
	}

	return files, err
}

func sorted(path string, sortType string) []*file {
	files, err := readdir(path)
	if err != nil {
		panic(err)
	}

	switch sortType {
	case "natural":
		sort.SliceStable(files, func(i, j int) bool {
			s1, s2 := normalize(files[i].Name(), files[j].Name())
			return naturalLess(s1, s2)
		})
	case "name":
		sort.SliceStable(files, func(i, j int) bool {
			s1, s2 := normalize(files[i].Name(), files[j].Name())
			return s1 < s2
		})
	case "size":
		sort.SliceStable(files, func(i, j int) bool {
			return files[i].Size() < files[j].Size()
		})
	case "time":
		sort.SliceStable(files, func(i, j int) bool {
			return files[i].ModTime().Before(files[j].ModTime())
		})
	case "atime":
		sort.SliceStable(files, func(i, j int) bool {
			return files[i].accessTime.Before(files[j].accessTime)
		})
	case "ctime":
		sort.SliceStable(files, func(i, j int) bool {
			return files[i].changeTime.Before(files[j].changeTime)
		})
	case "ext":
		sort.SliceStable(files, func(i, j int) bool {
			ext1, ext2 := normalize(files[i].ext, files[j].ext)

			// if the extension could not be determined (directories, files without)
			// use a zero byte so that these files can be ranked higher
			if ext1 == "" {
				ext1 = "\x00"
			}
			if ext2 == "" {
				ext2 = "\x00"
			}

			name1, name2 := normalize(files[i].Name(), files[j].Name())

			// in order to also have natural sorting with the filenames
			// combine the name with the ext but have the ext at the front
			return ext1 < ext2 || ext1 == ext2 && name1 < name2
		})
	}

	return files
}

// TAKEN FROM LF: END

func main() {
	// XXX: account for setlocal ... sortby
	// Right now it SILENTLY supersedes lf_sort
	flag.Usage = func() {
		fmt.Println("sortlf <diretory>")
		fmt.Println("	like `ls`, but with the sorting algo from `lf`")
		fmt.Println("	respects `lf_sortby`, `lf_reverse` and `lf_hidden` env. variables")
	}
	flag.Parse()
	path := flag.Arg(0)

	if path == "" {
		panic("no directory path provided, see sortlf -h")
	}

	info, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	if !info.IsDir() {
		panic(fmt.Sprintf("%s is not a directory", path))
	}

	sortType := os.Getenv("lf_sortby")
	reverse := os.Getenv("lf_reverse")
	hidden := os.Getenv("lf_hidden")
	if sortType == "" {
		panic("lf_sortby is empty")
	}
	if sortType == "" {
		panic("lf_reverse is empty")
	}
	if hidden == "" {
		panic("lf_hidden is empty")
	}

	files := sorted(path, sortType)

	if reverse == "true" {
		slices.Reverse(files)
	}

	for _, f := range files {
		basename := f.Name()
		if hidden == "false" && strings.HasPrefix(basename, ".") {
			continue
		}
		fmt.Println(basename)
	}
}
