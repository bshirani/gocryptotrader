package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	gateioDownloadURL = "https://download.gatedata.org"
	// gateioPathFormat = "/${biz}/${type}/${year}${month}/${market}-${year}${month}.csv.gz"
	gateioPathFormat = "/%s/%s/%s/%s-%s.csv.gz"
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

	// continue
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
	// os.Exit(1)

	// set the starting date (in any way you wish)
	start, _ := time.Parse("2006-01-02", "2021-08-01")
	t1 := start.AddDate(-3, 0, 0)
	// handle error

	// set d to starting date and keep adding 1 day to it as long as month doesn't change
	for _, p := range pairs {
		fmt.Printf("%s%10s%s...", string(colorCyan), p, string(colorReset))

		for d := start; d.After(t1); d = d.AddDate(0, -1, 0) {
			var monthYear string
			if int(d.Month()) < 10 {
				monthYear = fmt.Sprintf("%d0%d", d.Year(), d.Month())
			} else {
				monthYear = fmt.Sprintf("%d%-d", d.Year(), d.Month())
			}

			path := fmt.Sprintf(gateioDownloadURL+gateioPathFormat, "spot", "candlesticks_1m", monthYear, p, monthYear)

			// only if file does not exist
			// does file exist?

			filename := fmt.Sprintf("/home/bijan/work/crypto/gateiodata/%s/%s-%s.csv.gz", p, p, monthYear)
			filename404 := fmt.Sprintf("/home/bijan/work/crypto/gateiodata/%s/%s-%s.404", p, p, monthYear)
			if _, err := os.Stat(filename404); !errors.Is(err, os.ErrNotExist) {
				fmt.Printf("%s%s", string(colorRed), "❌")
				continue
			} else if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
				// create dir if necessary
				// fmt.Println("wget", path, filename)
				newpath := filepath.Join(".", p)
				os.MkdirAll(newpath, os.ModePerm)
				// continue
			} else {
				// file exists
				fmt.Printf("%s%s", string(colorGreen), "✅")
				continue
			}

			res, e := http.Get(path)
			// if e != nil {
			// 	fmt.Println("error!!", e)
			// 	continue
			// }
			if res.StatusCode > 299 {
				fmt.Printf("%s%s", string(colorRed), "❌")
				f, e := os.Create(filename404)
				if e != nil {
					panic(e)
				}
				f.Close()
				continue
			} else {
				fmt.Printf("%s%s", string(colorGreen), "✅")
			}
			defer res.Body.Close()
			io.ReadAll(res.Body)
			res.Body.Close()
			f, e := os.Create(filename)
			if e != nil {
				panic(e)
			}
			defer f.Close()
			f.ReadFrom(res.Body)
			// time.Sleep(time.Millisecond * 500)

			if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
				fmt.Println("FILE DOES NOT EXIST", filename)
				os.Exit(123)
			}
		}
		fmt.Println(string(colorReset))
	}

}
