package engine

import (
	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/portfolio/strategies"
)

// loads active strategies from the database
func SetupStrategies(cfg []*config.StrategySetting, exch string) (slit []strategies.Handler) {
	count := 0
	for _, cs := range cfg {
		count += 1
		strat, _ := strategies.LoadStrategyByName(cs.Capture)

		strat.SetID(count)
		strat.SetWeight(cs.Weight)
		strat.SetDirection(cs.Side)
		pair := currency.GetPairTranslation(exch, cs.Pair)
		// fmt.Println("setting pair", pair, "for exchange", exch)
		strat.SetPair(pair)
		strat.SetName(cs.Capture)
		strat.SetDefaults()
		slit = append(slit, strat)
	}

	if len(slit) == 0 {
		panic("no strategies loaded")
	}
	return slit
}

func isActivePair(pairs currency.Pairs, mypair currency.Pair) bool {
	for _, p := range pairs {
		if currency.ArePairsEqual(p, mypair) {
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
