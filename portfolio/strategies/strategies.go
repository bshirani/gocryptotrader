package strategies

import (
	"fmt"
	"strings"

	"gocryptotrader/portfolio/strategies/base"
	"gocryptotrader/portfolio/strategies/trend"
	"gocryptotrader/portfolio/strategies/trend2day"
	"gocryptotrader/portfolio/strategies/trend3day"
	"gocryptotrader/portfolio/strategies/trenddev"
)

// LoadStrategyByName returns the strategy by its name
func LoadStrategyByName(name string) (Handler, error) {
	strats := GetStrategies()
	for i := range strats {

		if !strings.EqualFold(name, strats[i].Name()) {
			continue
		}
		// if useSimultaneousProcessing {
		// 	if !strats[i].SupportsSimultaneousProcessing() {
		// 		return nil, fmt.Errorf(
		// 			"strategy '%v' %w",
		// 			name,
		// 			base.ErrSimultaneousProcessingNotSupported)
		// 	}
		// 	strats[i].SetSimultaneousProcessing(useSimultaneousProcessing)
		// }
		return strats[i], nil
	}
	fmt.Printf("strategy '%v' %w", name, base.ErrStrategyNotFound)
	return nil, fmt.Errorf("strategy '%v' %w", name, base.ErrStrategyNotFound)
}

// GetStrategies returns a static list of set strategies
// they must be set in here for the backtester to recognise them
func GetStrategies() []Handler {
	x := []Handler{
		// new(dollarcostaverage.Strategy),
		new(trend.Strategy),
		new(trend2day.Strategy),
		new(trend3day.Strategy),
		new(trenddev.Strategy),
		// new(top2bottom2.Strategy),
	}
	// for i := range x {
	// 	fmt.Println(x[i])
	// 	x[i].SetDirection(order.Sell)
	// }
	return x
}
