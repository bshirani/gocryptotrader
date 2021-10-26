package engine

import (
	"fmt"
	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/portfolio/strategies"
	"os"
)

// loads active strategies from the database
func SetupStrategies(exMgr iExchangeManager, cfg []*config.StrategySetting, liveMode bool) (slit []strategies.Handler) {
	count := 0
	for _, cs := range cfg {
		count += 1
		strat, _ := strategies.LoadStrategyByName(cs.Capture)

		var pair currency.Pair
		var err error

		if liveMode {
			pair, err = currency.NewPairFromString(cs.Pair.Symbol)
		} else {
			if cs.Pair.BacktestSymbol == "" {
				fmt.Println("backtest symbol cannot be blank")
				os.Exit(123)
			}
			pair, err = currency.NewPairFromString(cs.Pair.BacktestSymbol)

		}

		if err != nil {
			fmt.Println("error hydrating pair:", pair, "error:", err)
			os.Exit(123)
			return
		}

		enablePair(exMgr, pair)

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

func enablePair(exMgr iExchangeManager, pair currency.Pair) {
	exs, _ := exMgr.GetExchanges()
	// fmt.Println("loaded", len(exs), "exchanges")
	if len(exs) == 0 {
		panic(123)
	}

	// err = c.SetPairs(exchName, assetTypes[0], true, currency.Pairs{newPair})

	for _, ex := range exs {
		// fmt.Println("enabling", ex, pair)
		ex.SetPairs(currency.Pairs{pair}, asset.Spot, true)
	}
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
