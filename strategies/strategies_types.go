package strategies

import (
	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/data"
	"github.com/thrasher-corp/gocryptotrader/eventtypes/signal"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
	"github.com/thrasher-corp/gocryptotrader/strategies/base"
)

// Handler defines all functions required to run strategies against data events
type Handler interface {
	ID() string
	SetID(string)
	Direction() order.Side
	GetPair() currency.Pair
	SetPair(currency.Pair)
	Name() string
	Stop()
	Description() string
	OnData(data.Handler, base.PortfolioHandler, base.FactorEngineHandler) (signal.Event, error)
	OnSimultaneousSignals([]data.Handler, base.PortfolioHandler, base.FactorEngineHandler) ([]signal.Event, error)
	UsingSimultaneousProcessing() bool
	SupportsSimultaneousProcessing() bool
	SetSimultaneousProcessing(bool)
	SetDirection(order.Side)
	SetCustomSettings(map[string]interface{}) error
	SetDefaults()
	SetWeight(decimal.Decimal)
}
