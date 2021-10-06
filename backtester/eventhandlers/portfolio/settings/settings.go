package settings

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/holdings"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/trades"
)

// GetLatestHoldings returns the latest holdings after being sorted by time
func (e *Settings) GetLatestHoldings() holdings.Holding {
	if len(e.HoldingsSnapshots) == 0 {
		return holdings.Holding{}
	}

	return e.HoldingsSnapshots[len(e.HoldingsSnapshots)-1]
}

// GetOpenTrades returns the latest holdings after being sorted by time
func (e *Settings) GetOpenTrades() map[string]trades.Trade {
	if len(e.TradesMap) == 0 {
		e.TradesMap = make(map[string]trades.Trade)
	}

	return e.TradesMap
}

// GetHoldingsForTime returns the holdings for a time period, or an empty holding if not found
func (e *Settings) GetHoldingsForTime(t time.Time) holdings.Holding {
	if e.HoldingsSnapshots == nil {
		// no holdings yet
		return holdings.Holding{}
	}
	for i := len(e.HoldingsSnapshots) - 1; i >= 0; i-- {
		if e.HoldingsSnapshots[i].Timestamp.Equal(t) {
			return e.HoldingsSnapshots[i]
		}
	}
	return holdings.Holding{}
}

func (e *Settings) GetTradeForStrategy(s string) trades.Trade {
	if len(e.TradesMap) == 0 {
		e.TradesMap = make(map[string]trades.Trade)
	}

	// if e.TradesMap == nil {
	// 	return trades.Trade{}
	// }

	return e.TradesMap[s]
}

// Value returns the total value of the latest holdings
func (e *Settings) Value() decimal.Decimal {
	latest := e.GetLatestHoldings()
	if latest.Timestamp.IsZero() {
		return decimal.Zero
	}
	return latest.TotalValue
}
