package analyze

import (
	"bytes"
	"encoding/json"
	"gocryptotrader/common/file"
	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/database/repository/livetrade"
	"gocryptotrader/exchange/order"
	"gocryptotrader/log"
	"gocryptotrader/portfolio/strategies"
	"io"
	"strings"

	"github.com/shopspring/decimal"
)

func (p *PortfolioAnalysis) AnalyzeStrategies() {
	p.analyzeGrouped()
	p.Report.Strategies = p.StrategiesAnalyses
}

func (p *PortfolioAnalysis) analyzeGrouped() {
	p.Strategies = make([]strategies.Handler, 0)
	for label, trades := range p.groupedTrades {
		strat := loadStrategyFromLabel(label)
		a := analyzeStrategy(strat, trades)
		p.StrategiesAnalyses = append(p.StrategiesAnalyses, a)
		p.GroupedSettings = append(p.GroupedSettings, strat.GetSettings())
		p.Strategies = append(p.Strategies, strat)
	}
}

func analyzeStrategy(s strategies.Handler, trades []*livetrade.Details) (a *StrategyAnalysis) {
	a = &StrategyAnalysis{}
	var predTrades []*livetrade.Details
	for _, t := range trades {
		if t.Prediction > 0 {
			predTrades = append(predTrades, t)
		}
	}

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
	a.NumTrades = len(trades)
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

	sumPl = 0.0
	sumProfits = 0.0
	sumLosses = 0.0
	winCount = 0
	lossCount = 0
	for _, t := range predTrades {
		pl, _ := t.ProfitLossQuote.Float64()
		pl = pl * t.Prediction
		sumPl += pl
		if pl > 0 {
			sumProfits += pl
			winCount += 1
		} else {
			sumLosses += pl
			lossCount += 1
		}
	}
	a.NumTradesPredicted = len(predTrades)
	a.NetProfitPredicted = decimal.NewFromFloat(sumPl)
	a.AveragePLPredicted = sumPl / float64(winCount+lossCount)
	a.AverageWinPredicted = sumProfits / float64(winCount)
	a.AverageLossPredicted = sumLosses / float64(lossCount)
	a.WinPercentagePredicted = float64(winCount) / float64(len(trades))

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

func SaveStrategiesConfigFile(outpath string, ss []config.StrategySetting) error {
	writer, err := file.Writer(outpath)
	defer func() {
		if writer != nil {
			err = writer.Close()
			if err != nil {
				log.Error(log.Global, err)
			}
		}
	}()
	payload, err := json.MarshalIndent(ss, "", " ")
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, bytes.NewReader(payload))
	return err
}
func (p *PortfolioAnalysis) GetStrategyAnalysis(s strategies.Handler) *StrategyAnalysis {
	for _, a := range p.StrategiesAnalyses {
		if strings.EqualFold(a.Label, s.GetLabel()) {
			return a
		}
	}
	panic("could not find strategy analysis")
	return nil
}

func loadStrategyFromTrade(t *livetrade.Details) strategies.Handler {
	s, _ := strategies.LoadStrategyByName("trend")
	s.SetName("trend")
	s.SetDirection(t.Side)
	s.SetPair(t.Pair)
	s.SetID(t.StrategyID)
	// fmt.Println("strategy label", s.GetLabel(), s.Name())
	return s
}

func loadStrategyFromLabel(label string) strategies.Handler {
	l := strings.Split(label, "@")
	name := l[0]
	symbol := l[1]
	dir := l[2]

	s, _ := strategies.LoadStrategyByName(name)
	s.SetName(name)
	s.SetDirection(order.Side(dir))
	pair, _ := currency.NewPairFromString(symbol)

	s.SetPair(pair)
	// s.SetID(t.StrategyID)
	// fmt.Println("strategy label", s.GetLabel(), s.Name())
	return s
}

func groupByStrategyID(trades []*livetrade.Details) (grouped map[string][]*livetrade.Details) {
	grouped = make(map[string][]*livetrade.Details)

	for _, lt := range trades {
		s := loadStrategyFromTrade(lt)
		grouped[s.GetLabel()] = append(grouped[s.GetLabel()], lt)
	}
	return grouped
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
// 			t.ProfitLossPoints = t.ExitPrice.Sub(t.EntryPrice)
// 		} else if t.Side == order.Sell {
// 			t.ProfitLossPoints = t.EntryPrice.Sub(t.ExitPrice)
// 		}
// 		netProfit = netProfit.Add(t.Amount.Mul(t.ProfitLossPoints))
// 	}
// 	return netProfit
// }
