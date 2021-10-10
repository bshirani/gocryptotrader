package engine

import (
	"errors"
	"sync"

	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/compliance"
	"github.com/thrasher-corp/gocryptotrader/config"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/database/repository/livetrade"
	"github.com/thrasher-corp/gocryptotrader/eventtypes"
	"github.com/thrasher-corp/gocryptotrader/eventtypes/fill"
	"github.com/thrasher-corp/gocryptotrader/eventtypes/order"
	"github.com/thrasher-corp/gocryptotrader/eventtypes/signal"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	gctorder "github.com/thrasher-corp/gocryptotrader/exchanges/order"
	"github.com/thrasher-corp/gocryptotrader/factors"
	"github.com/thrasher-corp/gocryptotrader/portfolio/holdings"
	"github.com/thrasher-corp/gocryptotrader/portfolio/positions"
	"github.com/thrasher-corp/gocryptotrader/portfolio/risk"
	"github.com/thrasher-corp/gocryptotrader/strategies"
)

var (
	errNoDecision           = errors.New("signal has no decision")
	errInvalidDirection     = errors.New("invalid direction")
	errRiskManagerUnset     = errors.New("risk manager unset")
	errStrategyIDUnset      = errors.New("strategy id unset")
	errAlreadyInTrade       = errors.New("already in trade")
	errSizeManagerUnset     = errors.New("size manager unset")
	errAssetUnset           = errors.New("asset unset")
	errNoOpenTrade          = errors.New("no trade open")
	errExchangeUnset        = errors.New("exchange unset")
	errNegativeRiskFreeRate = errors.New("received negative risk free rate")
	errNoPortfolioSettings  = errors.New("no portfolio settings")
	errNoHoldings           = errors.New("no holdings found")
	errHoldingsNoTimestamp  = errors.New("holding with unset timestamp received")
	errTradesNoTimestamp    = errors.New("trade with unset timestamp received")
	errHoldingsAlreadySet   = errors.New("holding already set")
	errTradesAlreadySet     = errors.New("trade already set")
)

type portfolioStore struct {
	m            sync.RWMutex
	positions    map[string]*positions.Position
	openTrade    map[string]*livetrade.Details
	closedTrades map[string][]*livetrade.Details
	wg           *sync.WaitGroup
}

// Portfolio stores all holdings and rules to assess orders, allowing the portfolio manager to
// modify, accept or reject strategy signals
type Portfolio struct {
	riskFreeRate              decimal.Decimal
	sizeManager               SizeHandler
	riskManager               risk.Handler
	factorEngine              *factors.Engine
	bot                       Engine
	strategies                []strategies.Handler
	store                     portfolioStore
	exchangeAssetPairSettings map[string]map[asset.Item]map[currency.Pair]*PortfolioSettings
}

// Settings holds all important information for the portfolio manager
// to assess purchasing decisions
type PortfolioSettings struct {
	Fee               decimal.Decimal
	BuySideSizing     config.MinMax
	SellSideSizing    config.MinMax
	Leverage          config.Leverage
	HoldingsSnapshots []holdings.Holding
	ComplianceManager compliance.Manager
}

// Settings allow the eventhandler to size an order within the limitations set by the config file
type ExchangeAssetPairSettings struct {
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

// Exchange contains all the currency settings
type Exchange struct {
	CurrencySettings []ExchangeAssetPairSettings
}

// Handler contains all functions expected to operate a portfolio manager
type PortfolioHandler interface {
	OnSignal(signal.Event, *ExchangeAssetPairSettings) (*order.Order, error)
	OnFill(fill.Event) (*fill.Fill, error)

	ViewHoldingAtTimePeriod(eventtypes.EventHandler) (*holdings.Holding, error)
	setHoldingsForOffset(*holdings.Holding, bool) error
	UpdateHoldings(eventtypes.DataEventHandler) error
	GetTradeForStrategy(string) *livetrade.Details
	GetPositionForStrategy(string) *positions.Position

	SetFee(string, asset.Item, currency.Pair, decimal.Decimal)
	GetFee(string, asset.Item, currency.Pair) decimal.Decimal

	Reset()
}

// SizeHandler is the interface to help size orders
type SizeHandler interface {
	SizeOrder(order.Event, decimal.Decimal, *ExchangeAssetPairSettings) (*order.Order, error)
}
