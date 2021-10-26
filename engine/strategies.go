package engine

import (
	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/portfolio/strategies"
)

// loads active strategies from the database
func SetupStrategies(cfg []*config.StrategySetting, liveMode bool) (slit []strategies.Handler) {
	count := 0
	for _, cs := range cfg {
		count += 1
		strat, _ := strategies.LoadStrategyByName(cs.Capture)

		strat.SetID(count)
		strat.SetWeight(cs.Weight)
		strat.SetDirection(cs.Side)
		strat.SetPair(cs.Pair)
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
	// for _, x := range strategies {
	// 	fmt.Println("Loaded Strategy:", x.Name(), x.GetPair(), x.GetDirection(), x.GetWeight())
	// }
}
