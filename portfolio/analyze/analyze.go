package analyze

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gocryptotrader/common/file"
	"gocryptotrader/config"
	"gocryptotrader/database/repository/livetrade"
	"gocryptotrader/exchange/order"
	"gocryptotrader/log"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/shopspring/decimal"
)

func (p *PortfolioAnalysis) Analyze(filepath string) error {
	p.Report = &PortfolioReport{}
	lf := lastResult()
	trades, err := livetrade.LoadCSV(lf)
	enhanced := enhanceTrades(trades)
	grouped := groupByStrategyID(enhanced)
	fmt.Println("pfloaded", len(trades), "trades from", len(grouped), "strategies")
	p.Report.StrategiesAnalyses = make(map[int]*StrategyAnalysis)
	for id := range grouped {
		p.Report.StrategiesAnalyses[id] = analyzeStrategy(id, grouped[id])
	}

	sumDurationMin := 0.0
	for _, lt := range trades {
		sumDurationMin += lt.DurationMinutes
	}
	p.Report.AverageDurationMin = sumDurationMin / float64(len(trades))

	p.calculateProductionWeights()

	return err
}

func (p *PortfolioAnalysis) calculateProductionWeights() {
	p.Weights = &PortfolioWeights{}
	p.Weights.Strategies = make([]*config.StrategySetting, 0)
	// write all strategies for now
	for _, s := range p.Strategies {
		ss := &config.StrategySetting{
			Weight:  decimal.NewFromFloat(1.5),
			Side:    s.GetDirection(),
			Capture: "trend",
		}
		p.Weights.Strategies = append(p.Weights.Strategies, ss)
	}
}

func (p *PortfolioAnalysis) PrintResults() {
	for sid, sa := range p.Report.StrategiesAnalyses {
		fmt.Println("strategy", sid, "num trades", sa.NumTrades)
	}
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

func groupByStrategyID(trades []*livetrade.Details) (grouped map[int][]*livetrade.Details) {
	grouped = make(map[int][]*livetrade.Details)
	for _, lt := range trades {
		if grouped[lt.StrategyID] == nil {
			grouped[lt.StrategyID] = make([]*livetrade.Details, 0)
		}
		if lt.StrategyID == 0 {
			fmt.Println(lt)
			panic("livetrade has no strategy id")
		}
		grouped[lt.StrategyID] = append(grouped[lt.StrategyID], lt)
	}
	return grouped
}

func enhanceTrades(trades []*livetrade.Details) []*livetrade.Details {
	// create detailed trades
	// run preparation
	calculateDuration(trades)
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

func analyzeStrategy(id int, trades []*livetrade.Details) (a *StrategyAnalysis) {
	a = &StrategyAnalysis{}
	// a.Trades = trades
	a.NumTrades = len(trades)
	return a
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
			t.ProfitLossPoints = t.ExitPrice.Sub(t.EntryPrice)
		} else if t.Side == order.Sell {
			t.ProfitLossPoints = t.EntryPrice.Sub(t.ExitPrice)
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
	dir := filepath.Join(wd, "../backtest/results")
	lf := lastFileInDir(dir)

	if err != nil {
		fmt.Println(err)
	}
	return filepath.Join(wd, "../backtest/results", lf)
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
	// if len(names) > 0 {
	// 	fmt.Println(modTime, names)
	// }
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
	payload, err := json.MarshalIndent(p, "", " ")
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, bytes.NewReader(payload))
	return err
}

func (w *PortfolioWeights) Save(filepath string) error {
	writer, err := file.Writer(filepath)
	defer func() {
		if writer != nil {
			err = writer.Close()
			if err != nil {
				log.Error(log.Global, err)
			}
		}
	}()
	payload, err := json.MarshalIndent(w, "", " ")
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, bytes.NewReader(payload))
	return err
}
