package analyze

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gocryptotrader/common/file"
	"gocryptotrader/currency"
	"gocryptotrader/database/repository/livetrade"
	"gocryptotrader/exchange/order"
	"gocryptotrader/log"
	"gocryptotrader/portfolio/strategies"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

const (
	prodExchange     = "kraken"
	backtestExchange = "gateio"
)

func (p *PortfolioAnalysis) Analyze(filepath string) error {
	p.Report = &Report{}
	p.Report.Portfolio = &PortfolioReport{}
	lf := lastResult()
	fmt.Println("analyzing trades csv:", lf)
	trades, err := livetrade.LoadCSV(lf)
	enhanced := enhanceTrades(trades)
	p.trades = enhanced
	p.groupedTrades = groupByStrategyID(enhanced)
	p.loadAllStrategies()
	p.loadGroupedStrategies()
	p.analyzeGroupedStrategies()
	p.calculateReport()
	p.calculateProductionWeights()
	p.Report.Strategies = p.StrategiesAnalyses

	return err
}

func (p *PortfolioAnalysis) analyzeGroupedStrategies() {
}

func (p *PortfolioAnalysis) loadAllStrategies() {
	// get a list of all the strategy names
	// all := strategies.GetStrategies()

	names := []string{"trend", "trend2day", "trend3day"}
	symbols := []string{"ETH_USDT", "XBT_USDT", "DAI_USDT"}
	pairs := make([]currency.Pair, 0)
	for _, s := range symbols {
		pair, err := currency.NewPairFromString(s)
		if err != nil {
			fmt.Println("error hydrating pair", pair)
		}
		pairs = append(pairs, pair.Upper())
	}
	// pairs, _ := p.Config.GetEnabledPairs("gateio", asset.Spot)

	for _, name := range names {
		// for each direction
		for _, dir := range []order.Side{order.Buy, order.Sell} {
			for _, pair := range pairs {
				strat, _ := strategies.LoadStrategyByName(name)
				strat.SetDirection(dir)
				strat.SetPair(pair)
				strat.SetName(name)
				p.AllSettings = append(p.AllSettings, strat.GetSettings())
			}
		}
	}
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

func (p *PortfolioAnalysis) GetStrategyAnalysis(s strategies.Handler) *StrategyAnalysis {
	for _, a := range p.StrategiesAnalyses {
		if strings.EqualFold(a.Label, s.GetLabel()) {
			return a
		}
	}
	panic("could not find strategy analysis")
	return nil
}

func (p *PortfolioAnalysis) calculateReport() {
	// fmt.Println("pfloaded", len(trades), "trades from", len(grouped), "strategies")
	// p.Report.StrategiesAnalyses = make(map[strategies.Handler]*StrategyAnalysis)
	// for id := range grouped {

	for range p.AllSettings {
		// if the strategy is in the trades group

		// for i := range p.groupedTrades {
		// 	fmt.Println("check", i, ss.Capture, ss.Pair.Symbol, ss.Side)
		// }

		// p.Report.StrategiesAnalyses[id] = analyzeStrategy(id, grouped[id])
	}
	// }
	sumDurationMin := 0.0
	for _, lt := range p.trades {
		sumDurationMin += lt.DurationMinutes
	}
	p.Report.Portfolio.AverageDurationMin = sumDurationMin / float64(len(p.trades))
}

func (p *PortfolioAnalysis) PrintResults() {
	// for sid, sa := range p.Report.StrategiesAnalyses {
	// 	fmt.Println("strategy", sid, "num trades", sa.NumTrades)
	// }
}

// func (p *PortfolioAnalysis) WriteOutput() {
// 	fmt.Println("analyzing", len(p.StrategiesAnalyses), "strategies")
// 	for sid, sa := range p.StrategiesAnalyses {
// 		fmt.Println("strategy", sid, "num trades", sa.NumTrades)
// 	}
// }

func PrintTradeResults() {
	// for _, t := range trades {
	// 	fmt.Printf("enter=%v exit=%v enter=%v exit=%v profit=%v minutes=%d amount=%v stop=%v\n",
	// 		t.EntryTime.Format(common.SimpleTimeFormat),
	// 		t.ExitTime.Format(common.SimpleTimeFormat),
	// 		t.EntryPrice,
	// 		t.ExitPrice,
	// 		getProfit(t),
	// 		getDurationMin(t),
	// 		t.Amount,
	// 		t.StopLossPrice,
	// 	)
	// }
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
	l := strings.Split(label, ":")
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

func enhanceTrades(trades []*livetrade.Details) []*livetrade.Details {
	// create detailed trades
	// run preparation
	calculateDuration(trades)
	netProfitPoints(trades)
	netProfit(trades)
	return trades

	// enhance
	// for i := range trades {
	// 	enhanced = append(enhanced, &livetrade.Details{
	// 		EntryTime:       trades[i].EntryTime,
	// 		ExitTime:        trades[i].ExitTime,
	// 		EntryPrice:      trades[i].EntryPrice,
	// 		ExitPrice:       trades[i].ExitPrice,
	// 		Side:            trades[i].Side,
	// 		Amount:          trades[i].Amount,
	// 		StrategyID:      trades[i].StrategyID,
	// 		StopLossPrice:   trades[i].StopLossPrice,
	// 		TakeProfitPrice: trades[i].TakeProfitPrice,
	// 		Status:          trades[i].Status,
	// 		Pair:            trades[i].Pair,
	// 		CreatedAt:       trades[i].CreatedAt,
	// 		UpdatedAt:       trades[i].UpdatedAt,
	// 	})
	// }
	// return enhanced
}

func netProfitPoints(trades []*livetrade.Details) (netProfit decimal.Decimal) {
	for _, t := range trades {
		if t.Side == order.Buy {
			t.ProfitLossPoints = t.ExitPrice.Sub(t.EntryPrice)
		} else if t.Side == order.Sell {
			t.ProfitLossPoints = t.EntryPrice.Sub(t.ExitPrice)
		}
		netProfit = netProfit.Add(t.ProfitLossPoints)
	}
	return netProfit
}

func netProfit(trades []*livetrade.Details) (netProfit decimal.Decimal) {
	for _, t := range trades {
		if t.Side == order.Buy {
			t.ProfitLoss = t.ExitPrice.Sub(t.EntryPrice)
		} else if t.Side == order.Sell {
			t.ProfitLoss = t.EntryPrice.Sub(t.ExitPrice)
		}
		netProfit = netProfit.Add(t.Amount.Mul(t.ProfitLossPoints))
	}
	return netProfit
}

func calculateDuration(trades []*livetrade.Details) {
	for _, t := range trades {
		t.DurationMinutes = t.ExitTime.Sub(t.EntryTime).Minutes()
	}
}

func lastResult() string {
	// return os.MkdirAll(dir, 0770)
	wd, err := os.Getwd()
	dir := filepath.Join(wd, "results/bt")
	lf := lastFileInDir(dir)

	if err != nil {
		fmt.Println(err)
	}
	return filepath.Join(wd, "results/bt", lf)
}

func lastFileInDir(dir string) string {
	files, _ := ioutil.ReadDir(dir)
	var modTime time.Time
	var names []string
	for _, fi := range files {
		if fi.Mode().IsRegular() {
			if !fi.ModTime().Before(modTime) {
				if fi.ModTime().After(modTime) {
					modTime = fi.ModTime()
					names = names[:0]
				}
				names = append(names, fi.Name())
			}
		}
	}
	if len(names) == 0 {
		panic(fmt.Sprintf("could not find file in dir %s", dir))
		fmt.Println(modTime, names)
	}
	return names[len(names)-1]
}

func getProfit(trade livetrade.Details) decimal.Decimal {
	if trade.Side == order.Buy {
		return trade.ExitPrice.Sub(trade.EntryPrice)
	} else if trade.Side == order.Sell {
		return trade.EntryPrice.Sub(trade.ExitPrice)
	}
	return decimal.Decimal{}
}

func getDurationMin(trade livetrade.Details) int {
	return int(trade.ExitTime.Sub(trade.EntryTime).Minutes())
}

// func calculateMaxDrawdown(closePrices []eventtypes.DataEventHandler) Swing {
// 	var lowestPrice, highestPrice decimal.Decimal
// 	var lowestTime, highestTime time.Time
// 	var swings []Swing
// 	if len(closePrices) > 0 {
// 		lowestPrice = closePrices[0].LowPrice()
// 		highestPrice = closePrices[0].HighPrice()
// 		lowestTime = closePrices[0].GetTime()
// 		highestTime = closePrices[0].GetTime()
// 	}
// 	for i := range closePrices {
// 		currHigh := closePrices[i].HighPrice()
// 		currLow := closePrices[i].LowPrice()
// 		currTime := closePrices[i].GetTime()
// 		if lowestPrice.GreaterThan(currLow) && !currLow.IsZero() {
// 			lowestPrice = currLow
// 			lowestTime = currTime
// 		}
// 		if highestPrice.LessThan(currHigh) && highestPrice.IsPositive() {
// 			if lowestTime.Equal(highestTime) {
// 				// create distinction if the greatest drawdown occurs within the same candle
// 				lowestTime = lowestTime.Add((time.Hour * 23) + (time.Minute * 59) + (time.Second * 59))
// 			}
// 			intervals, err := gctkline.CalculateCandleDateRanges(highestTime, lowestTime, closePrices[i].GetInterval(), 0)
// 			if err != nil {
// 				log.Error(log.TradeMgr, err)
// 				continue
// 			}
// 			swings = append(swings, Swing{
// 				Highest: Iteration{
// 					Time:  highestTime,
// 					Price: highestPrice,
// 				},
// 				Lowest: Iteration{
// 					Time:  lowestTime,
// 					Price: lowestPrice,
// 				},
// 				DrawdownPercent:  lowestPrice.Sub(highestPrice).Div(highestPrice).Mul(decimal.NewFromInt(100)),
// 				IntervalDuration: int64(len(intervals.Ranges[0].Intervals)),
// 			})
// 			// reset the drawdown
// 			highestPrice = currHigh
// 			highestTime = currTime
// 			lowestPrice = currLow
// 			lowestTime = currTime
// 		}
// 	}
// 	if (len(swings) > 0 && swings[len(swings)-1].Lowest.Price != closePrices[len(closePrices)-1].LowPrice()) || swings == nil {
// 		// need to close out the final drawdown
// 		if lowestTime.Equal(highestTime) {
// 			// create distinction if the greatest drawdown occurs within the same candle
// 			lowestTime = lowestTime.Add((time.Hour * 23) + (time.Minute * 59) + (time.Second * 59))
// 		}
// 		intervals, err := gctkline.CalculateCandleDateRanges(highestTime, lowestTime, closePrices[0].GetInterval(), 0)
// 		if err != nil {
// 			log.Error(log.TradeMgr, err)
// 		}
// 		drawdownPercent := decimal.Zero
// 		if highestPrice.GreaterThan(decimal.Zero) {
// 			drawdownPercent = lowestPrice.Sub(highestPrice).Div(highestPrice).Mul(decimal.NewFromInt(100))
// 		}
// 		if lowestTime.Equal(highestTime) {
// 			// create distinction if the greatest drawdown occurs within the same candle
// 			lowestTime = lowestTime.Add((time.Hour * 23) + (time.Minute * 59) + (time.Second * 59))
// 		}
// 		swings = append(swings, Swing{
// 			Highest: Iteration{
// 				Time:  highestTime,
// 				Price: highestPrice,
// 			},
// 			Lowest: Iteration{
// 				Time:  lowestTime,
// 				Price: lowestPrice,
// 			},
// 			DrawdownPercent:  drawdownPercent,
// 			IntervalDuration: int64(len(intervals.Ranges[0].Intervals)),
// 		})
// 	}
//
// 	var maxDrawdown Swing
// 	if len(swings) > 0 {
// 		maxDrawdown = swings[0]
// 	}
// 	for i := range swings {
// 		if swings[i].DrawdownPercent.LessThan(maxDrawdown.DrawdownPercent) {
// 			// drawdowns are negative
// 			maxDrawdown = swings[i]
// 		}
// 	}
//
// 	return maxDrawdown
// }
//
// func (c *CurrencyStatistic) calculateHighestCommittedFunds() {
// 	for i := range c.Events {
// 		if c.Events[i].Holdings.BaseSize.Mul(c.Events[i].DataEvent.ClosePrice()).GreaterThan(c.HighestCommittedFunds.Value) {
// 			c.HighestCommittedFunds.Value = c.Events[i].Holdings.BaseSize.Mul(c.Events[i].DataEvent.ClosePrice())
// 			c.HighestCommittedFunds.Time = c.Events[i].Holdings.Timestamp
// 		}
// 	}
// }

func (p *PortfolioAnalysis) Save(filepath string) error {
	writer, err := file.Writer(filepath)
	defer func() {
		if writer != nil {
			err = writer.Close()
			if err != nil {
				log.Error(log.Global, err)
			}
		}
	}()
	payload, err := json.MarshalIndent(p.Report, "", " ")
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, bytes.NewReader(payload))
	return err
}

func (p *PortfolioAnalysis) SaveAllStrategiesConfigFile(outpath string) error {
	writer, err := file.Writer(outpath)
	defer func() {
		if writer != nil {
			err = writer.Close()
			if err != nil {
				log.Error(log.Global, err)
			}
		}
	}()
	payload, err := json.MarshalIndent(p.AllSettings, "", " ")
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, bytes.NewReader(payload))
	return err
}
