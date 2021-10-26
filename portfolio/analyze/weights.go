package analyze

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gocryptotrader/common/file"
	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/log"
	"io"
	"strings"
)

func (p *PortfolioAnalysis) calculateProductionWeights() {
	p.Weights = &PortfolioWeights{}
	p.Weights.Strategies = make([]*config.StrategySetting, 0)
	fmt.Println("there are", len(p.groupedTrades), "strategies")
	for i := range p.groupedTrades {
		fmt.Println("strategy", i)
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

func (p *PortfolioAnalysis) getPairForExchange(ex string, pair currency.Pair) currency.Pair {
	// pairs, err := p.Config.GetEnabledPairs(ex, asset.Spot)

	if strings.EqualFold(pair.Base.String(), "BTC") {
		pair.Base = currency.NewCode("XBT")
	}
	return pair

	// 	fmt.Println("pair", exp.Base, exp.Quote)
	// 	if strings.EqualFold(exp.Base.String(), "XBT") {
	// 		exp.Base = currency.NewCode("XBT")
	// 	}
	// }
	//
	// if err != nil {
	// 	fmt.Println("error", err)
	// }
}
