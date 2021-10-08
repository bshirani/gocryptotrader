package strategies

import (
	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/backtester/data"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/strategies/base"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventtypes/signal"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
)

// Handler defines all functions required to run strategies against data events
type Handler interface {
	ID() string
	SetID(string)
	Direction() order.Side
	Name() string
	Stop()
	Description() string
	OnData(data.Handler, base.PortfolioHandler) (signal.Event, error)
	OnSimultaneousSignals([]data.Handler, base.PortfolioHandler) ([]signal.Event, error)
	UsingSimultaneousProcessing() bool
	SupportsSimultaneousProcessing() bool
	SetSimultaneousProcessing(bool)
	SetDirection(order.Side)
	SetCustomSettings(map[string]interface{}) error
	SetDefaults()
	SetWeight(decimal.Decimal)
}
