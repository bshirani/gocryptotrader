package analyze

import (
	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/database/repository/livetrade"
	"gocryptotrader/exchange/order"
	"gocryptotrader/portfolio/strategies"
	"strings"

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

func GenerateAllStrategies() (out []config.StrategySetting) {
	names := []string{"trend"}
	pairs := make([]currency.Pair, 0)
	symlist := "AAVEUSDT,AVAXUSDT,AXSUSDT,BATUSDT,CRVUSDT,DASHUSDT,DOGEUSDT,FILUSDT,FTMUSDT,ICPUSDT,KAVAUSDT,KNCUSDT,KSMUSDT,LITUSDT,LUNAUSDT,MATICUSDT,MKRUSDT,NEARUSDT,OMGUSDT,RENUSDT,SHIBUSDT,SNXUSDT,SOLUSDT,SUSHIUSDT,THETAUSDT,UNFIUSDT,XLMUSDT,XTZUSDT,YFIUSDT,ZRXUSDT,ADAUSDT,ATOMUSDT,BCHUSDT,BTCUSDT,DOTUSDT,EOSUSDT,ETCUSDT,ETHUSDT,LINKUSDT,LTCUSDT,TRXUSDT,UNIUSDT,XRPUSDT"
	symbols := strings.Split(symlist, ",")
	for _, s := range symbols {
		base := strings.Split(s, "USDT")[0]
		pair := currency.NewPairWithDelimiter(base, "USDT", "_")
		pairs = append(pairs, pair)
	}

	for _, name := range names {
		for _, dir := range []order.Side{order.Buy, order.Sell} {
			for _, pair := range pairs {
				strat, _ := strategies.LoadStrategyByName(name)
				strat.SetDirection(dir)
				strat.SetPair(pair)
				strat.SetName(name)
				out = append(out, *strat.GetSettings())
			}
		}
	}
	return out
}
