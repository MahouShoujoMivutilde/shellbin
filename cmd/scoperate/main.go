package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var DESC string = os.Args[0] + `
  Track the amount of new-calls per ID within some timeframe

`

var EXAMPLES string = `
Examples:
  get current rate of new-call's for given ID (e.g. PID/name/etc) within last 500ms
    scoperate -id ID -action get-rate

  register new-call for given ID to increase its rate
    scoperate -id ID -action new-call
`

const gBasePath string = "/tmp/scoperate.id."
const gTimeFrame time.Duration = 500 * time.Millisecond

func usage() {
	fmt.Fprint(flag.CommandLine.Output(), DESC)

	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %[1]s:\n", os.Args[0])

	flag.PrintDefaults()
	fmt.Fprint(flag.CommandLine.Output(), EXAMPLES)
}

func loadTimes(id string) ([]time.Time, error) {
	times := make([]time.Time, 10)

	file, err := os.Open(gBasePath + id)
	if err != nil {
		return times, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&times)
	return times, err
}

func saveTimes(id string, times []time.Time) error {
	err := os.MkdirAll(filepath.Dir(gBasePath+id), 0700)
	if err != nil {
		return err
	}

	file, err := os.Create(gBasePath + id)
	if err != nil {
		return err
	}

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(&times)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Usage = usage
	id := flag.String("id", "", "limit scope of actions to given ID; meant for distinguishing different instances of lf")
	action := flag.String("action", "", "pick action: new-call or get-rate")
	flag.Parse()

	if *id == "" || (*action != "new-call" && *action != "get-rate") {
		flag.Usage()
		os.Exit(1)
	}

	if *action == "new-call" {
		times, _ := loadTimes(*id)
		now := time.Now()
		times = append(times, now)

		newTimes := make([]time.Time, 0)
		for _, t := range times {
			if now.Sub(t) <= gTimeFrame {
				newTimes = append(newTimes, t)
			}
		}

		err := saveTimes(*id, newTimes)
		if err != nil {
			panic(fmt.Sprintf("Call not saved! 👿 - %s", err))
		}

	} else if *action == "get-rate" {
		times, _ := loadTimes(*id)
		now := time.Now()

		newTimes := make([]time.Time, 0)
		for _, t := range times {
			if now.Sub(t) <= gTimeFrame && !t.IsZero() {
				newTimes = append(newTimes, t)
			}
		}

		fmt.Println(float64(len(newTimes)) / gTimeFrame.Seconds())
	}
}
