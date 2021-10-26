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

	outpath := "portfolio_analysis.json"
	fmt.Println("saving analysis to", outpath)
	pf.Save(outpath)
	weights := "pf_weighted.json"
	fmt.Println("saving pf weights to", weights)
	pf.Weights.Save(weights)
}
