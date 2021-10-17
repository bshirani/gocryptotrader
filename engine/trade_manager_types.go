package engine

import (
	"context"
	"errors"
	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/data"
	"gocryptotrader/eventtypes"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/exchange/order"
	gctorder "gocryptotrader/exchange/order"
	"gocryptotrader/portfolio/report"
	"gocryptotrader/portfolio/statistics"
	"gocryptotrader/portfolio/strategies"
	"sync"
	"time"
)

const tradeManagerName = "trade_manager"

var (
	errInvalidConfigAsset     = errors.New("invalid asset in config")
	errAmbiguousDataSource    = errors.New("ambiguous settings received. Only one data type can be set")
	errNoDataSource           = errors.New("no data settings set in config")
	errIntervalUnset          = errors.New("candle interval unset")
	errUnhandledDatatype      = errors.New("unhandled datatype")
	errLiveDataTimeout        = errors.New("no data returned in 5 minutes, shutting down")
	errNilData                = errors.New("nil data received")
	errNilExchange            = errors.New("nil exchange received")
	errDataMayBeIncorrect     = errors.New("data may be incorrect")
	errExceededPortfolioLimit = errors.New("exceeded portfolio limit")
	errNilCurrencySettings    = errors.New("received nil currency settings")
	// errInvalidDirection       = errors.New("received invalid order direction")
)

// TradeManager is the main holder of all strategy functionality
type TradeManager struct {
	bot                *Engine
	DataHistoryManager *DataHistoryManager
	Datas              data.Holder
	EventQueue         EventHolder
	Exchange           ExecutionHandler
	FactorEngines      map[string]map[asset.Item]map[currency.Pair]*FactorEngine
	OrderManager       OrderManagerHandler
	Portfolio          PortfolioHandler
	Reports            report.Handler
	Statistic          statistics.Handler
	Strategies         []strategies.Handler
	Warmup             bool
	cfg                config.Config
	hasHandledEvent    bool
	shutdown           chan struct{}
	started            int32
	syncManager        *syncManager
	verbose            bool
	wg                 sync.WaitGroup
}

// Holder contains the event queue for backtester processing
type Holder struct {
	Queue []eventtypes.EventHandler
}

// EventHolder interface details what is expected of an event holder to perform
type EventHolder interface {
	Reset()
	AppendEvent(eventtypes.EventHandler)
	NextEvent() eventtypes.EventHandler
}

// ExecutionHandler interface dictates what functions are required to submit an order
type ExecutionHandler interface {
	// SetExchangeAssetCurrencySettings(string, asset.Item, currency.Pair, *ExchangeAssetPairSettings)
	// GetAllCurrencySettings() ([]ExchangeAssetPairSettings, error)
	// GetCurrencySettings(string, asset.Item, currency.Pair) (ExchangeAssetPairSettings, error)

	GetOrdersSnapshot(order.Status) ([]order.Detail, time.Time)
	Submit(context.Context, *gctorder.Submit) (*OrderSubmitResponse, error)
	SetOnSubmit(func(*OrderSubmitResponse))
	SetOnFill(func(*OrderSubmitResponse))
	SetOnCancel(func(*OrderSubmitResponse))
	// Reset()
}
