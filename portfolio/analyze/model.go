package analyze

import (
	"encoding/json"
	"fmt"
	"gocryptotrader/common"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	host = "http://localhost:8000"
)

func BacktestModel() {
	lastCSV, _ := common.LastFileInDir("results/fcsv")
	// SelectFeatures(lastCSV)
	GetPredictions(lastCSV)
}

func SelectFeatures(filename string) {
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
	url := fmt.Sprintf(
		"%s/predict?file=../results/fcsv/%s",
		host,
		filename)
	fmt.Println(url)
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
