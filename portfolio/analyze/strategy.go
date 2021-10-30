package analyze

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	p.loadGroupedStrategies()
	fmt.Println("saving strategies", len(p.StrategiesAnalyses))
	p.Report.Strategies = p.StrategiesAnalyses
}

func (p *PortfolioAnalysis) loadGroupedStrategies() {
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
