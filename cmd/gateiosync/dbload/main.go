package main

import (
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
	fmt.Println(1)

	files, err := ioutil.ReadDir(baseDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			pairPath := path.Join(baseDir, file.Name())
			files, _ = ioutil.ReadDir(pairPath)
			os.Chdir(pairPath)
			newDir, _ := os.Getwd()
			fmt.Println("changing dir to", pairPath, newDir)

			c, err := currency.NewPairFromString(file.Name())
			if err != nil {
				fmt.Println("cant find currency", err)
			}
			for _, file = range files {
				if strings.Contains(file.Name(), ".csv") {
					cmd := fmt.Sprintf(baseCmd, "gateio", c.Base.String(), c.Quote.String(), path.Join(pairPath, file.Name()))
					fmt.Println(cmd)
					out, err := exec.Command("bash", "-c", cmd).Output()
					if err != nil {
						fmt.Println("err running cmd", cmd, out, err)
						os.Exit(2)
					}
				}
			}
		}
	}
}
