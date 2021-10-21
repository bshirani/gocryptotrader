package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
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
	baseDir          = "/home/bijan/work/crypto/gateiodata"
	colorReset       = "\033[0m"

	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

func main() {
	start, _ := time.Parse("2006-01-02", "2021-08-01")
	t1 := start.AddDate(-3, 0, 0)
	finished := make(chan bool)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			<-finished
			os.Exit(0)
		}
	}()

	for _, p := range symbols() {
		fmt.Printf("%s%10s%s...", string(colorCyan), p, string(colorReset))
		for d := start; d.After(t1); d = d.AddDate(0, -1, 0) {
			go worker(d, p, finished)
			<-finished
		}
		fmt.Println(string(colorReset))
	}
}

func worker(d time.Time, p string, finished chan bool) {
	defer func() {
		finished <- true
	}()
	var monthYear string
	if int(d.Month()) < 10 {
		monthYear = fmt.Sprintf("%d0%d", d.Year(), d.Month())
	} else {
		monthYear = fmt.Sprintf("%d%-d", d.Year(), d.Month())
	}

	path := fmt.Sprintf(gateioDownloadURL+gateioPathFormat, "spot", "candlesticks_1m", monthYear, p, monthYear)

	// only if file does not exist
	// does file exist?

	csvFilename := fmt.Sprintf("%s/%s/%s-%s.csv", baseDir, p, p, monthYear)
	filename := fmt.Sprintf("%s/%s/%s-%s.csv.gz", baseDir, p, p, monthYear)
	filename404 := fmt.Sprintf("%s/%s/%s-%s.404", baseDir, p, p, monthYear)

	if _, err := os.Stat(csvFilename); !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("%s%s", string(colorGreen), "E")
		return

	} else if _, err := os.Stat(filename404); !errors.Is(err, os.ErrNotExist) {
		// 404ed already
		fmt.Printf("%s%s", string(colorRed), "4")
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
		return
	}

	res, e := http.Get(path)
	if res.StatusCode > 299 {
		fmt.Printf("%s%s", string(colorRed), "F")
		f, e := os.Create(filename404)
		if e != nil {
			panic(e)
		}
		f.Close()
		return
	} else {
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

func symbols() []string {
	file, err := os.Open("./symbols.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	pairs := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pairs = append(pairs, scanner.Text())
	}
	fmt.Println("loaded", len(pairs))
	return pairs
}
