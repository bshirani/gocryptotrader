package portfolio

import (
	"errors"
	"sync"

	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/backtester/common"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/exchange"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/compliance"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/factors"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/holdings"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/positions"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/risk"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/settings"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/strategies"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/trades"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventtypes/fill"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventtypes/order"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventtypes/signal"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/engine"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
)

var (
	errNoDecision           = errors.New("signal has no decision")
	errInvalidDirection     = errors.New("invalid direction")
	errRiskManagerUnset     = errors.New("risk manager unset")
	errStrategyIDUnset      = errors.New("strategy id unset")
	errAlreadyInTrade       = errors.New("already in trade")
	errSizeManagerUnset     = errors.New("size manager unset")
	errAssetUnset           = errors.New("asset unset")
	errCurrencyPairUnset    = errors.New("currency pair unset")
	errNoOpenTrade          = errors.New("no trade open")
	errExchangeUnset        = errors.New("exchange unset")
	errNegativeRiskFreeRate = errors.New("received negative risk free rate")
	errNoPortfolioSettings  = errors.New("no portfolio settings")
	errNoHoldings           = errors.New("no holdings found")
	errNoTrades             = errors.New("no trades found")
	errHoldingsNoTimestamp  = errors.New("holding with unset timestamp received")
	errTradesNoTimestamp    = errors.New("trade with unset timestamp received")
	errHoldingsAlreadySet   = errors.New("holding already set")
	errTradesAlreadySet     = errors.New("trade already set")
)

type store struct {
	m            sync.RWMutex
	positions    map[string]*positions.Position
	openTrade    map[string]*trades.Trade
	closedTrades map[string][]*trades.Trade
	wg           *sync.WaitGroup
}

// Portfolio stores all holdings and rules to assess orders, allowing the portfolio manager to
// modify, accept or reject strategy signals
type Portfolio struct {
	isLive                    bool
	riskFreeRate              decimal.Decimal
	sizeManager               SizeHandler
	riskManager               risk.Handler
	factorEngine              *factors.Engine
	bot                       engine.Engine
	strategies                []strategies.Handler
	store                     store
	exchangeAssetPairSettings map[string]map[asset.Item]map[currency.Pair]*settings.Settings
}

// Handler contains all functions expected to operate a portfolio manager
type Handler interface {
	OnSignal(signal.Event, *exchange.Settings) (*order.Order, error)
	OnFill(fill.Event) (*fill.Fill, error)

	ViewHoldingAtTimePeriod(common.EventHandler) (*holdings.Holding, error)
	setHoldingsForOffset(*holdings.Holding, bool) error
	UpdateHoldings(common.DataEventHandler) error
	GetPositionForStrategy(string) *positions.Position
	GetTradeForStrategy(string) *trades.Trade

	GetComplianceManager(string, asset.Item, currency.Pair) (*compliance.Manager, error)

	SetFee(string, asset.Item, currency.Pair, decimal.Decimal)
	GetFee(string, asset.Item, currency.Pair) decimal.Decimal

	Reset()
}

// SizeHandler is the interface to help size orders
type SizeHandler interface {
	SizeOrder(order.Event, decimal.Decimal, *exchange.Settings) (*order.Order, error)
}
