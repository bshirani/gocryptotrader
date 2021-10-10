package engine

import (
	"errors"

	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/data"
	"github.com/thrasher-corp/gocryptotrader/eventtypes/fill"
	"github.com/thrasher-corp/gocryptotrader/eventtypes/order"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
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
