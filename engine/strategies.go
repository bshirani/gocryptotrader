package engine

import (
	"fmt"
	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/portfolio/strategies"
	"os"
)

// loads active strategies from the database
func SetupStrategies(cfg *config.Config) (slit []strategies.Handler) {
	// load the strategies from the condif
	count := 0
	fmt.Println("strategies", len(cfg.TradeManager.Strategies))
	for _, cs := range cfg.TradeManager.Strategies {
		count += 1
		strat, _ := strategies.LoadStrategyByName(cs.Capture)

		// fmt.Println("load pair", cs.Pair)

		pair, err := currency.NewPairFromString(cs.Pair.Symbol)
		if err != nil {
			fmt.Println("error hydrating pair:", pair, err)
			os.Exit(123)
			return
		}

		strat.SetID(count)
		strat.SetWeight(cs.Weight)
		strat.SetDirection(cs.Side)
		strat.SetPair(pair)
		strat.SetDefaults()
		// fmt.Println("created strategy", strat.GetPair(), strat.GetDirection(), strat.Name(), strat.GetWeight())
		slit = append(slit, strat)
	}

	if len(slit) == 0 {
		panic("no strategies loaded")
	}
	return slit
}

func isActivePair(pairs currency.Pairs, mypair currency.Pair) bool {
	for _, p := range pairs {
		if arePairsEqual(p, mypair) {
			return true
		}
	}
	return false
}

func printStrategies(strategies []strategies.Handler) {
	for _, x := range strategies {
		fmt.Println("Loaded Strategy:", x.Name(), x.GetPair(), x.GetDirection(), x.GetWeight())
	}
}
