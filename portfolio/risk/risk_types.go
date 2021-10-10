package risk

import (
	"errors"

	"github.com/shopspring/decimal"
	"gocryptotrader/currency"
	"gocryptotrader/eventtypes/order"
	"gocryptotrader/exchanges/asset"
	"gocryptotrader/portfolio/compliance"
	"gocryptotrader/portfolio/holdings"
)

var (
	errNoCurrencySettings       = errors.New("lacking currency settings, cannot evaluate order")
	errLeverageNotAllowed       = errors.New("order is using leverage when leverage is not enabled in config")
	errCannotPlaceLeverageOrder = errors.New("cannot place leveraged order")
)

// Handler defines what is expected to be able to assess risk of an order
type Handler interface {
	EvaluateOrder(order.Event, []holdings.Holding, compliance.Snapshot) (*order.Order, error)
}

// Risk contains all currency settings in order to evaluate potential orders
type Risk struct {
	CurrencySettings map[string]map[asset.Item]map[currency.Pair]*CurrencySettings
	CanUseLeverage   bool
	MaximumLeverage  decimal.Decimal
}

// CurrencySettings contains relevant limits to assess risk
type CurrencySettings struct {
	MaximumOrdersWithLeverageRatio decimal.Decimal
	MaxLeverageRate                decimal.Decimal
	MaximumHoldingRatio            decimal.Decimal
}
