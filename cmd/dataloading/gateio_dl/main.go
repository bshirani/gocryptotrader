package main

import (
	"errors"
	"fmt"
	"gocryptotrader/common"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

const (
	gateioDownloadURL = "https://download.gatedata.org"
	// gateioPathFormat = "/${biz}/${type}/${year}${month}/${market}-${year}${month}.csv.gz"
	gateioPathFormat = "/%s/%s/%s/%s-%s.csv.gz"
	baseDir          = "/Users/bijan/gateiodata"
	colorReset       = "\033[0m"
	startDate        = "2022-03-01"

	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

func main() {
	start, _ := time.Parse("2006-01-02", startDate)
	t1 := start.AddDate(-300, 0, 0)
	finished := make(chan string)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			<-finished
			os.Exit(0)
		}
	}()

	// fmt.Println("loaded symbols", len(common.Symbols())
	for _, p := range common.Symbols() {
		var results []string
		fmt.Printf("%s%10s%s...", string(colorCyan), p, string(colorReset))
		for d := start; d.After(t1); d = d.AddDate(0, -1, 0) {
			go worker(d, p, finished)
			lr := <-finished
			results = append(results, lr)
			// check if last 3 results have been failures
			if len(results) > 3 {
				if checkAllFailed(results[len(results)-3:]) {
					break
				}
			}

			// fmt.Println(results)
		}
		fmt.Println(string(colorReset))
	}
}

func checkAllFailed(results []string) bool {
	for _, x := range results[len(results)-3:] {
		if x != "4" && x != "F" {
			return false
		}
	}
	return true
}

func worker(d time.Time, p string, finished chan string) {
	var monthYear string
	if int(d.Month()) < 10 {
		monthYear = fmt.Sprintf("%d0%d", d.Year(), d.Month())
	} else {
		monthYear = fmt.Sprintf("%d%-d", d.Year(), d.Month())
	}

	path := fmt.Sprintf(gateioDownloadURL+gateioPathFormat, "spot", "candlesticks_1m", monthYear, p, monthYear)
	// fmt.Println("path", path)

	// only if file does not exist
	// does file exist?

	csvFilename := fmt.Sprintf("%s/%s/%s-%s.csv", baseDir, p, p, monthYear)
	filename := fmt.Sprintf("%s/%s/%s-%s.csv.gz", baseDir, p, p, monthYear)
	filename404 := fmt.Sprintf("%s/%s/%s-%s.404", baseDir, p, p, monthYear)

	if _, err := os.Stat(csvFilename); !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("%s%s", string(colorGreen), "E")
		defer func() {
			finished <- "E"
		}()
		return

	} else if _, err := os.Stat(filename404); !errors.Is(err, os.ErrNotExist) {
		// 404ed already
		fmt.Printf("%s%s", string(colorRed), "4")
		defer func() {
			finished <- "4"
		}()
		return
	} else if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		// file does not exist
		// create dir if necessary
		// fmt.Println("wget", path, filename)
		newpath := filepath.Join(baseDir, p)
		os.MkdirAll(newpath, os.ModePerm)
	} else {
		// file exists

		fmt.Printf("%s%s", string(colorGreen), "E")
		defer func() {
			finished <- "E"
		}()
		return
	}

	res, e := http.Get(path)
	if res.StatusCode > 299 {
		fmt.Printf("%s%s", string(colorRed), "F")
		// fmt.Println("failed", path)
		defer func() {
			finished <- "F"
		}()
		f, e := os.Create(filename404)
		if e != nil {
			panic(e)
		}
		f.Close()
		return
	} else {
		defer func() {
			finished <- "S"
		}()
		fmt.Printf("%s%s", string(colorGreen), "S")
	}
	defer res.Body.Close()

	f, e := os.Create(filename)
	if e != nil {
		panic(e)
	}
	defer f.Close()
	io.Copy(f, res.Body)
	// time.Sleep(time.Millisecond * 500)

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		fmt.Println("FILE DOES NOT EXIST", filename)
		os.Exit(123)
	}

}
