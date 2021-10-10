package engine

import (
	"errors"

	"github.com/shopspring/decimal"
	config "github.com/thrasher-corp/gocryptotrader/bt_config"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/data"
	"github.com/thrasher-corp/gocryptotrader/eventtypes/fill"
	"github.com/thrasher-corp/gocryptotrader/eventtypes/order"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	gctorder "github.com/thrasher-corp/gocryptotrader/exchanges/order"
)

var (
	errDataMayBeIncorrect     = errors.New("data may be incorrect")
	errExceededPortfolioLimit = errors.New("exceeded portfolio limit")
	errNilCurrencySettings    = errors.New("received nil currency settings")
	// errInvalidDirection       = errors.New("received invalid order direction")
)

// ExecutionHandler interface dictates what functions are required to submit an order
type ExecutionHandler interface {
	SetExchangeAssetCurrencySettings(string, asset.Item, currency.Pair, *PortfolioExchangeSettings)
	GetAllCurrencySettings() ([]PortfolioExchangeSettings, error)
	GetCurrencySettings(string, asset.Item, currency.Pair) (PortfolioExchangeSettings, error)
	ExecuteOrder(order.Event, data.Handler, *OrderManagerHandler) (*fill.Fill, error)
	Reset()
}

type OrderManagerHandler interface {
	ExecuteOrder(order.Event, data.Handler) (*fill.Fill, error)
}

// Exchange contains all the currency settings
type FakeExchange struct {
	CurrencySettings []PortfolioExchangeSettings
}

// Settings allow the eventhandler to size an order within the limitations set by the config file
type PortfolioExchangeSettings struct {
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
