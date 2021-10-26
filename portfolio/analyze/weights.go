package analyze

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gocryptotrader/common/file"
	"gocryptotrader/config"
	"gocryptotrader/log"
	"io"

	"github.com/shopspring/decimal"
)

func (p *PortfolioAnalysis) calculateProductionWeights() {
	p.Weights = &PortfolioWeights{}
	p.Weights.Strategies = make([]*config.StrategySetting, 0)
	fmt.Println("here")

	// get the performance for this strategy
	for _, s := range p.Strategies {
		analysis := p.GetStrategyAnalysis(s)
		fmt.Println("analysis",
			analysis.Label,
			analysis.NumTrades,
			analysis.NetProfit,
		)

		if analysis.NetProfit.IsZero() {
			panic("net profit is zero")
		}
		ss := s.GetSettings()
		if analysis.NetProfit.GreaterThan(decimal.NewFromFloat(0.0)) {
			ss.Weight = decimal.NewFromFloat(1.0)
			p.Weights.Strategies = append(p.Weights.Strategies, ss)
		}
	}
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
