package main

import (
	"bufio"
	"errors"
	"fmt"
	"gocryptotrader/currency"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	baseDir          = "/home/bijan/work/crypto/kraken_data"
	baseCmd          = "dbseed candle file --exchange %s --base %s --quote %s --interval 60 --asset spot --filename %s"
	finishedFilename = "finished.log"
)

func main() {
	finished := finishedFiles()

	files, err := ioutil.ReadDir(baseDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), "_1.csv") {
			for _, fin := range finished {
				if strings.EqualFold(fin, f.Name()) {
					continue
				}
			}
			task(f.Name())
		}
	}
}

func task(fileName string) {
	// fmt.Println(fileName)
	dirName := strings.Split(fileName, "-")[0]
	c, err := currency.NewPairFromString(dirName)
	if err != nil {
		fmt.Println("cant find currency", err)
	}
	cmd := fmt.Sprintf(baseCmd, "kraken", c.Base.String(), c.Quote.String(), path.Join(baseDir, fileName))
	fmt.Println(cmd)
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Println("err running cmd", cmd, out, err)
		os.Exit(2)
	} else {
		fmt.Printf("%s", out)
	}
	markFileFinished(fileName)
}

func markFileFinished(symbol string) {
	f, err := os.OpenFile(finishedFilename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	output := fmt.Sprintf("%s\n", symbol)
	if _, err := f.WriteString(output); err != nil {
		log.Println(err)
	}
}

func finishedFiles() []string {
	if _, err := os.Stat(finishedFilename); errors.Is(err, os.ErrNotExist) {
		f, e := os.Create(finishedFilename)
		if e != nil {
			panic(e)
		}
		f.Close()
	}

	file, err := os.Open(finishedFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	pairs := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pairs = append(pairs, scanner.Text())
	}
	return pairs
}
