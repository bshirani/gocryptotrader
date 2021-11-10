package strategies

import (
	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/data"
	"gocryptotrader/database/repository/livetrade"
	"gocryptotrader/eventtypes/signal"
	"gocryptotrader/exchange/order"
	"gocryptotrader/portfolio/strategies/base"
	"time"

	"github.com/shopspring/decimal"
)

// Handler defines all functions required to run strategies against data events
type Handler interface {
	GetLabel() string
	GetPrediction(base.FactorEngineHandler, time.Time) float64
	SetDropFeatures()
	SetName(string)
	GetSettings() *config.StrategySetting
	SetID(int)
	GetID() int
	GetWeight() decimal.Decimal
	Learn(base.FactorEngineHandler, []*livetrade.Details) error
	SetNumID(int)
	GetNumID() int
	GetDirection() order.Side
	GetPair() currency.Pair
	// GetCurrencySettings() *ExchangeAssetPairSettings
	// SetCurrencySettings(*ExchangeAssetPairSettings)
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
