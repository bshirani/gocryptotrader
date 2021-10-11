package strategies

import (
	"github.com/shopspring/decimal"
	"gocryptotrader/currency"
	"gocryptotrader/data"
	"gocryptotrader/eventtypes/signal"
	"gocryptotrader/exchanges/order"
	"gocryptotrader/portfolio/strategies/base"
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
	OnData(data.Handler, base.StrategyPortfolioHandler, base.FactorEngineHandler) (signal.Event, error)
	OnSimultaneousSignals([]data.Handler, base.StrategyPortfolioHandler, base.FactorEngineHandler) ([]signal.Event, error)
	UsingSimultaneousProcessing() bool
	SupportsSimultaneousProcessing() bool
	SetSimultaneousProcessing(bool)
	SetDirection(order.Side)
	SetCustomSettings(map[string]interface{}) error
	SetDefaults()
	SetWeight(decimal.Decimal)
}
