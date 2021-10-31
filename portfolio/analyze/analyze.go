package analyze

import (
	"gocryptotrader/common"
	"gocryptotrader/config"
	"gocryptotrader/database/repository/livetrade"
	"os"
	"path/filepath"
)

func SetupPortfolio(cfg *config.Config, filepath string) (pf *PortfolioAnalysis, err error) {
	return &PortfolioAnalysis{
		Config:   cfg,
		Filepath: filepath,
	}, nil
}

func (p *PortfolioAnalysis) Analyze() error {
	err := p.loadTradesFromFile(p.Filepath)
	p.AnalyzePortfolio()
	p.AnalyzeStrategies()
	p.CalculateProductionWeights()
	p.Save("")
	return err
}

func (p *PortfolioAnalysis) loadTradesFromFile(filepath string) error {
	p.Report = &Report{}
	p.Report.Portfolio = &PortfolioReport{}
	lf, err := getTradeFilePath(filepath)
	trades, err := livetrade.LoadJSON(lf)
	p.trades = trades
	p.groupedTrades = groupByStrategyID(trades)
	return err
}

func getTradeFilePath(path string) (string, error) {
	wd, _ := os.Getwd()
	var err error
	if path == "" {
		path, err = common.LastFileInDir("results/bt")
	}
	return filepath.Join(wd, "results/bt", path), err
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
