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
			pairPath := path.Join(baseDir, file.Name())
			files, _ = ioutil.ReadDir(pairPath)

			var alreadyProcessed bool
			c, err := currency.NewPairFromString(file.Name())
			if err != nil {
				fmt.Println("cant find currency", err)
			}
			for _, file = range files {
				for _, x := range finished {
					if strings.EqualFold(x, file.Name()) {
						alreadyProcessed = true
						continue
					}
				}
				if !alreadyProcessed {
					if strings.HasSuffix(file.Name(), ".csv") {
						cmd := fmt.Sprintf(baseCmd, "gateio", c.Base.String(), c.Quote.String(), path.Join(pairPath, file.Name()))
						fmt.Println(cmd)
						out, err := exec.Command("bash", "-c", cmd).Output()
						if err != nil {
							fmt.Println("err running cmd", cmd, out, err)
							os.Exit(2)
						}
						markFileFinished(file.Name())
					}
				}
				alreadyProcessed = false
			}
		}
	}
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
	fmt.Println("loaded", len(pairs))
	return pairs
}
