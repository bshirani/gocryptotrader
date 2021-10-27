package analyze

import (
	"fmt"
	"gocryptotrader/database/repository/livetrade"
	"gocryptotrader/portfolio/strategies"

	"github.com/shopspring/decimal"
)

func analyzeStrategy(strat strategies.Handler, trades []*livetrade.Details) (a *StrategyAnalysis) {
	a = &StrategyAnalysis{}
	a.NumTrades = len(trades)
	sum := decimal.NewFromFloat(0.0)
	for _, t := range trades {
		sum = sum.Add(t.ProfitLoss)
	}
	fmt.Println("setting net profit", sum)
	a.NetProfit = sum
	a.Label = strat.GetLabel()
	return a
}
