package main

import (
	"fmt"

	"gocryptotrader/currency/coinmarketcap"
	"gocryptotrader/engine"
	"gocryptotrader/portfolio/analyze"
)

type Syncer struct {
	bot *engine.Engine
	cmc *coinmarketcap.Coinmarketcap
}

func main() {
	pf := &analyze.PortfolioAnalysis{}
	err := pf.Analyze("")
	if err != nil {
		fmt.Println("error analyzeTrades", err)
	}
	pf.PrintResults()
}
