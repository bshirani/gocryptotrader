package strategies

import (
	"gocryptotrader/currency"
	"gocryptotrader/data"
	"gocryptotrader/eventtypes/signal"
	"gocryptotrader/exchange/order"
	"gocryptotrader/portfolio/strategies/base"

	"github.com/shopspring/decimal"
)

// Handler defines all functions required to run strategies against data events
type Handler interface {
	SetID(int)
	GetID() int
	SetNumID(int)
	GetNumID() int
	GetDirection() order.Side
	GetPair() currency.Pair
	SetPair(currency.Pair)
	Name() string
	Stop()
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
