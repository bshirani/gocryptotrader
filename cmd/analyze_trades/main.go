package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gocryptotrader/config"
	"gocryptotrader/currency/coinmarketcap"
	"gocryptotrader/engine"
	"gocryptotrader/portfolio/analyze"
)

type Syncer struct {
	bot *engine.Engine
	cmc *coinmarketcap.Coinmarketcap
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Could not get working directory. Error: %v.\n", err)
		os.Exit(1)
	}
	configPath := filepath.Join(wd, "../confs/dev/backtest.json")
	cfg, err := config.ReadConfigFromFile(configPath)
	if err != nil {
		fmt.Printf("Could not read config. Error: %v. Path: %s\n", err, configPath)
		os.Exit(1)
	}

	pf := &analyze.PortfolioAnalysis{
		Config: cfg,
	}
	err = pf.Analyze("")
	if err != nil {
		fmt.Println("error analyzeTrades", err)
	}

	outpath := "portfolio_analysis.json"
	fmt.Println("saving analysis to", outpath)
	pf.Save(outpath)
	weights := "pf_weighted.json"
	fmt.Println("saving pf weights to", weights)
	pf.Weights.Save(weights)

	allPath := filepath.Join(wd, "../confs/dev/strategy/all.strat")
	pf.SaveAllStrategiesConfigFile(allPath)
}
