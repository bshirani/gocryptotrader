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

	fmt.Println("here")

	count := 0
	for _, cs := range cfg.TradeManager.Strategies {
		count += 1
		fmt.Println("cs", cs)
		strat, _ := strategies.LoadStrategyByName(cs.Capture)

		fmt.Println("load pair", cs.Pair)

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
		// fmt.Println("creating strategy for pair", baseStrategy.Capture, cps.CurrencyPair, cps.Side)
		slit = append(slit, strat)

	}

	// cpsS, _ := currencypairstrategy.All(cfg.LiveMode)
	//
	// for _, cps := range cpsS {
	//
	// 	// LOAD ONLY ACTIVE STRATEGIES
	//
	// 	// if !cps.Active {
	// 	// 	continue
	// 	// }
	// 	if cps.Weight.IsZero() && cfg.LiveMode {
	// 		fmt.Println("weightskip", cps.ID, cps.Weight)
	// 		continue
	// 	}
	//
	// 	var isWhitelisted bool
	// 	for _, name := range cfg.TradeManager.Strategies {
	// 		if strings.EqualFold(name, baseStrategy.Capture) {
	// 			isWhitelisted = true
	// 			break
	// 		}
	// 	}
	// 	if len(cfg.TradeManager.Strategies) == 0 && !strings.EqualFold(baseStrategy.Capture, "trenddev") {
	// 		isWhitelisted = true
	// 	}
	// 	if !isWhitelisted {
	// 		continue
	// 	}
	//
	// 	// filter out environment specific strategies
	// 	if cfg.ProductionMode && strings.EqualFold(baseStrategy.Capture, "trenddev") {
	// 		continue
	// 	}
	//
	// 	pairs, err := cfg.GetEnabledPairs("gateio", asset.Spot)
	// 	if err != nil {
	// 		fmt.Println("error getting pairs", err)
	// 	}
	// 	if !isActivePair(pairs, cps.CurrencyPair) {
	// 		continue
	// 	}
	//
	// 	strat, _ := strategies.LoadStrategyByName(baseStrategy.Capture)
	// 	fmt.Println("creating strategy", cps.ID, baseStrategy.Capture, cps.CurrencyPair, cps.Side)
	// 	strat.SetID(cps.ID)
	// 	strat.SetWeight(cps.Weight)
	// 	strat.SetNumID(cps.ID)
	// 	strat.SetPair(cps.CurrencyPair)
	// 	strat.SetDirection(cps.Side)
	// 	strat.SetDefaults()
	// 	// fmt.Println("creating strategy for pair", baseStrategy.Capture, cps.CurrencyPair, cps.Side)
	// 	slit = append(slit, strat)
	// }

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
		fmt.Println("ACTIVE STRATEGY", x.Name(), x.GetPair(), x.GetWeight())
	}
}
