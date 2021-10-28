package analyze

import (
	"gocryptotrader/database/repository/livetrade"
	"gocryptotrader/portfolio/strategies"

	"github.com/shopspring/decimal"
)

func analyzeStrategy(s strategies.Handler, trades []*livetrade.Details) (a *StrategyAnalysis) {
	a = &StrategyAnalysis{}
	a.NumTrades = len(trades)
	sumPl := 0.0
	sumProfits := 0.0
	sumLosses := 0.0
	winCount := 0
	lossCount := 0
	for _, t := range trades {
		pl, _ := t.ProfitLossQuote.Float64()
		sumPl += pl
		if pl > 0 {
			sumProfits += pl
			winCount += 1
		} else {
			sumLosses += pl
			lossCount += 1
		}
	}
	a.NetProfit = decimal.NewFromFloat(sumPl)
	a.Label = s.GetLabel()
	a.Name = s.Name()
	a.Direction = s.GetDirection()
	a.Pair = s.GetPair()
	a.StartDate = trades[0].EntryTime
	a.EndDate = trades[len(trades)-1].ExitTime
	a.AveragePL = sumPl / float64(winCount+lossCount)
	a.AverageWin = sumProfits / float64(winCount)
	a.AverageLoss = sumLosses / float64(lossCount)
	a.WinPercentage = float64(winCount) / float64(len(trades))

	return a
}
