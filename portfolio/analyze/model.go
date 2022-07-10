package analyze

import (
	"encoding/json"
	"fmt"
	"gocryptotrader/common"
	"gocryptotrader/database/repository/livetrade"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/shopspring/decimal"
)

const (
	host = "http://localhost:8000"
)

func (p *PortfolioAnalysis) BacktestModel() {
	lastCSV, _ := common.LastFileInDir("results/fcsv")
	// SelectFeatures(lastCSV)
	GetPredictions(lastCSV)
}

func (p *PortfolioAnalysis) SelectFeatures(filename string) {
	url := fmt.Sprintf("%s/select_features?file=../results/fcsv/%s", host, filename)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)
	for key, value := range result {
		fmt.Println(key, value.(float64))
	}
}

func GetPredictions(filename string) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Could not get working directory. Error: %v.\n", err)
		os.Exit(1)
	}
	lastBT, _ := common.LastFileInDir(BacktestResults)
	btPath := filepath.Join(wd, BacktestResults, lastBT)
	fmt.Println("lastbt", lastBT)
	trades, err := livetrade.LoadJSON(btPath)
	if err != nil {
		panic(err)
	}
	for _, t := range trades {
		t.Prediction = 0
	}
	url := fmt.Sprintf(
		"%s/predict?file=../results/fcsv/%s",
		host,
		filename)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)

	// zero the predictions from earlier runs
	for key, value := range result {
		for _, t := range trades {
			id, _ := strconv.ParseInt(key, 10, 32)
			if t.ID == int(id) {
				t.Prediction = value.(float64)
				t.PredictionAmount = decimal.NewFromFloat(t.Prediction).Mul(t.Amount)
				break
			}
		}
	}
	livetrade.WriteJSON(trades, btPath)
}
