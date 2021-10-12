package engine

import (
	"errors"

	"gocryptotrader/currency"
	"gocryptotrader/data"
	"gocryptotrader/eventtypes/fill"
	"gocryptotrader/eventtypes/order"
	"gocryptotrader/exchange/asset"
)

var (
	errDataMayBeIncorrect     = errors.New("data may be incorrect")
	errExceededPortfolioLimit = errors.New("exceeded portfolio limit")
	errNilCurrencySettings    = errors.New("received nil currency settings")
	// errInvalidDirection       = errors.New("received invalid order direction")
)

// ExecutionHandler interface dictates what functions are required to submit an order
type ExecutionHandler interface {
	SetExchangeAssetCurrencySettings(string, asset.Item, currency.Pair, *ExchangeAssetPairSettings)
	GetAllCurrencySettings() ([]ExchangeAssetPairSettings, error)
	GetCurrencySettings(string, asset.Item, currency.Pair) (ExchangeAssetPairSettings, error)
	ExecuteOrder(order.Event, data.Handler, *OrderManager) (*fill.Fill, error)
	Reset()
}

type OrderManagerHandler interface {
	ExecuteOrder(order.Event, data.Handler) (*fill.Fill, error)
}
