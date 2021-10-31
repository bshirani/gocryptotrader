package analyze

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gocryptotrader/common/file"
	"gocryptotrader/config"
	"gocryptotrader/log"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/shopspring/decimal"
)

func (p *PortfolioAnalysis) CalculateProductionWeights() {
	p.Weights = &PortfolioWeights{}
	p.Weights.Strategies = make([]*config.StrategySetting, 0)

	// get the performance for this strategy
	for _, s := range p.Strategies {
		analysis := p.GetStrategyAnalysis(s)
		// fmt.Println("analysis",
		// 	analysis.Label,
		// 	analysis.NumTrades,
		// 	analysis.NetProfit,
		// )

		if analysis.Base.NetProfit == 0 {
			panic("net profit is zero")
		}
		ss := s.GetSettings()
		if analysis.Base.NetProfit > 0.0 {
			ss.Weight = decimal.NewFromFloat(1.0)
			p.Weights.Strategies = append(p.Weights.Strategies, ss)
		}
	}
}

func (w *PortfolioWeights) Save(fpath string) error {
	if fpath == "" {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Could not get working directory. Error: %v.\n", err)
			os.Exit(1)
		}
		fpath = fmt.Sprintf(
			"portfolio_analysis_%v.json",
			time.Now().Format("2006-01-02-15-04-05"))
		fpath = filepath.Join(wd, "results/pf", fpath)
		fmt.Println("saving", fpath)
	}
	writer, err := file.Writer(fpath)
	defer func() {
		if writer != nil {
			err = writer.Close()
			if err != nil {
				log.Error(log.Global, err)
			}
		}
	}()
	payload, err := json.MarshalIndent(w.Strategies, "", " ")
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, bytes.NewReader(payload))
	return err
}

func (p *PortfolioAnalysis) AnalyzePortfolio() {
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

func (p *PortfolioAnalysis) Save(fpath string) error {
	if fpath == "" {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Could not get working directory. Error: %v.\n", err)
			os.Exit(1)
		}
		fpath = fmt.Sprintf(
			"portfolio_analysis_%v.json",
			time.Now().Format("2006-01-02-15-04-05"))
		fpath = filepath.Join(wd, "results/pf", fpath)
		fmt.Println("saving", fpath)
	}
	writer, err := file.Writer(fpath)
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
