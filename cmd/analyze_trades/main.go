package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

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

	filename := fmt.Sprintf(
		"portfolio_analysis_%v.json",
		time.Now().Format("2006-01-02-15-04-05"))
	filename = filepath.Join(wd, "results", filename)
	pf.Save(filename)

	prodWeighted := filepath.Join(wd, "../confs/dev/strategy/prod.strat")
	fmt.Println("saving", len(pf.Weights.Strategies), "pf weights to", prodWeighted)
	pf.Weights.Save(prodWeighted)

	allPath := filepath.Join(wd, "../confs/dev/strategy/all.strat")
	pf.SaveAllStrategiesConfigFile(allPath)
}
