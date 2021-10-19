package main

import (
	"bufio"
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
	baseDir = "/home/bijan/work/crypto/gateiodata"
	baseCmd = "dbseed candle file --exchange %s --base %s --quote %s --interval 60 --asset spot --filename %s"
)

func main() {
	finished := finishedSymbols()

	files, err := ioutil.ReadDir(baseDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			// if not finished
			files, _ = ioutil.ReadDir(path.Join(baseDir, file.Name()))
			for _, f := range files {
				if strings.HasSuffix(f.Name(), ".csv") {
					for _, fin := range finished {
						if strings.EqualFold(fin, f.Name()) {
							// fmt.Println("already processed")
							continue
						}
					}
					task(f.Name())
				}
			}
		}
	}
}

func task(fileName string) {
	fmt.Println(fileName)
	dirName := strings.Split(fileName, "-")[0]
	c, err := currency.NewPairFromString(dirName)
	if err != nil {
		fmt.Println("cant find currency", err)
	}
	cmd := fmt.Sprintf(baseCmd, "gateio", c.Base.String(), c.Quote.String(), path.Join(baseDir, dirName, fileName))
	// fmt.Println(cmd)
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Println("err running cmd", cmd, out, err)
		os.Exit(2)
	}
	markFileFinished(fileName)
}

func markFileFinished(symbol string) {
	f, err := os.OpenFile("./finished.log",
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

func finishedSymbols() []string {
	file, err := os.Open("./finished.log")
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
