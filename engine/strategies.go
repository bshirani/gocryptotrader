package engine

import (
	"fmt"
	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/database/repository/currencypairstrategy"
	"gocryptotrader/database/repository/strategy"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/portfolio/strategies"
	"strings"
)

func SetupStrategies(cfg *config.Config) (slit []strategies.Handler) {
	cpsS, _ := currencypairstrategy.All(cfg.LiveMode)

	for _, cps := range cpsS {
		if !cps.Active {
			continue
		}
		if cps.Weight.IsZero() && cfg.LiveMode {
			fmt.Println("skip", cps.ID, cps.Weight)
			continue
		}

		baseStrategy, _ := strategy.One(cps.StrategyID)
		if baseStrategy.TimeframeDays != 1 {
			continue
		}

		var isWhitelisted bool
		for _, name := range cfg.TradeManager.Strategies {
			if strings.EqualFold(name, baseStrategy.Capture) {
				isWhitelisted = true
				break
			}
		}
		if len(cfg.TradeManager.Strategies) == 0 && !strings.EqualFold(baseStrategy.Capture, "trenddev") {
			isWhitelisted = true
		}
		if !isWhitelisted {
			continue
		}

		pairs, err := cfg.GetEnabledPairs("gateio", asset.Spot)
		if err != nil {
			fmt.Println("error getting pairs", err)
		}
		if !isActivePair(pairs, cps.CurrencyPair) {
			continue
		}

		// fmt.Println("creating strategy for pair", cps.ID, baseStrategy.Capture, cps.CurrencyPair, cps.Side)
		strat, _ := strategies.LoadStrategyByName(baseStrategy.Capture)
		// fmt.Println("creating strategy", cps.ID, cps.CurrencyPair, cps.Side)
		strat.SetID(cps.ID)
		strat.SetNumID(cps.ID)
		strat.SetPair(cps.CurrencyPair)
		strat.SetDirection(cps.Side)
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
		if arePairsEqual(p, mypair) {
			return true
		}
	}
	return false
}
