package strategies

import (
	"fmt"
	"strings"

	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/strategies/base"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/strategies/dollarcostaverage"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/strategies/rsi"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/strategies/top2bottom2"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
)

// LoadStrategyByName returns the strategy by its name
func LoadStrategyByName(name string, direction order.Side, useSimultaneousProcessing bool) (Handler, error) {
	strats := GetStrategies()
	for i := range strats {
		strats[i].SetDirection(direction)

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
	return nil, fmt.Errorf("strategy '%v' %w", name, base.ErrStrategyNotFound)
}

// GetStrategies returns a static list of set strategies
// they must be set in here for the backtester to recognise them
func GetStrategies() []Handler {
	x := []Handler{
		new(dollarcostaverage.Strategy),
		new(rsi.Strategy),
		new(top2bottom2.Strategy),
	}
	// for i := range x {
	// 	fmt.Println(x[i])
	// 	x[i].SetDirection(order.Sell)
	// }
	return x
}
