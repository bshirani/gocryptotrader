package engine

import (
	"errors"
	"sync"

	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/database/repository/livetrade"
	"gocryptotrader/eventtypes"
	"gocryptotrader/eventtypes/fill"
	"gocryptotrader/eventtypes/order"
	"gocryptotrader/eventtypes/signal"
	"gocryptotrader/exchange/asset"
	gctorder "gocryptotrader/exchange/order"
	"gocryptotrader/portfolio/compliance"
	"gocryptotrader/portfolio/holdings"
	"gocryptotrader/portfolio/positions"
	"gocryptotrader/portfolio/risk"
	"gocryptotrader/portfolio/strategies"

	"github.com/shopspring/decimal"
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
	factorEngine              *FactorEngine
	bot                       *Engine
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
type StrategyPortfolioHandler interface {
	OnSignal(signal.Event, *ExchangeAssetPairSettings) (*order.Order, error)
	OnFill(fill.Event) (*fill.Fill, error)
	GetAllClosedTrades() []*livetrade.Details

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
