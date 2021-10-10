package exchange

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/backtester/config"
	"github.com/thrasher-corp/gocryptotrader/backtester/data"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventtypes/fill"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventtypes/order"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	gctorder "github.com/thrasher-corp/gocryptotrader/exchanges/order"
)

var (
	errDataMayBeIncorrect     = errors.New("data may be incorrect")
	errExceededPortfolioLimit = errors.New("exceeded portfolio limit")
	errNilCurrencySettings    = errors.New("received nil currency settings")
	errInvalidDirection       = errors.New("received invalid order direction")
)

// ExecutionHandler interface dictates what functions are required to submit an order
type ExecutionHandler interface {
	SetExchangeAssetCurrencySettings(string, asset.Item, currency.Pair, *Settings)
	GetAllCurrencySettings() ([]Settings, error)
	GetCurrencySettings(string, asset.Item, currency.Pair) (Settings, error)
	ExecuteOrder(order.Event, data.Handler, OrderManagerHandler) (*fill.Fill, error)
	Reset()
}

type OrderManagerHandler interface {
	ExecuteOrder(order.Event, data.Handler) (*fill.Fill, error)
}

// Exchange contains all the currency settings
type Exchange struct {
	CurrencySettings []Settings
}

// Settings allow the eventhandler to size an order within the limitations set by the config file
type Settings struct {
	ExchangeName  string
	UseRealOrders bool

	CurrencyPair currency.Pair
	AssetType    asset.Item

	ExchangeFee decimal.Decimal
	MakerFee    decimal.Decimal
	TakerFee    decimal.Decimal

	BuySide  config.MinMax
	SellSide config.MinMax

	Leverage config.Leverage

	MinimumSlippageRate decimal.Decimal
	MaximumSlippageRate decimal.Decimal

	Limits                  *gctorder.Limits
	CanUseExchangeLimits    bool
	SkipCandleVolumeFitting bool
}
