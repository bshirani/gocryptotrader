package analyze

import (
	"gocryptotrader/database/repository/livetrade"
	"gocryptotrader/portfolio/strategies"

	"github.com/shopspring/decimal"
)

func analyzeStrategy(strat strategies.Handler, trades []*livetrade.Details) (a *StrategyAnalysis) {
	a = &StrategyAnalysis{}
	// a.Trades = trades
	a.NumTrades = len(trades)

	sum := decimal.NewFromFloat(0.0)
	for _, t := range trades {
		sum = sum.Add(t.ProfitLoss)
	}
	a.NetProfit = sum
	a.Label = strat.GetLabel()
	return a
}

// func netProfitPoints(trades []*livetrade.Details) (netProfit decimal.Decimal) {
// 	for _, t := range trades {
// 		if t.Side == order.Buy {
// 			t.ProfitLossPoints = t.ExitPrice.Sub(t.EntryPrice)
// 		} else if t.Side == order.Sell {
// 			t.ProfitLossPoints = t.EntryPrice.Sub(t.ExitPrice)
// 		}
// 		netProfit = netProfit.Add(t.ProfitLossPoints)
// 	}
// 	return netProfit
// }
//
// func netProfit(trades []*livetrade.Details) (netProfit decimal.Decimal) {
// 	for _, t := range trades {
// 		if t.Side == order.Buy {
// 			t.ProfitLoss = t.ExitPrice.Sub(t.EntryPrice)
// 		} else if t.Side == order.Sell {
// 			t.ProfitLoss = t.EntryPrice.Sub(t.ExitPrice)
// 		}
// 		netProfit = netProfit.Add(t.Amount.Mul(t.ProfitLossPoints))
// 	}
// 	return netProfit
// }
