package engine

import (
	"errors"
	"gocryptotrader/backtester/statistics"
	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/data"
	"gocryptotrader/eventtypes"
	"gocryptotrader/portfolio/report"
	"gocryptotrader/portfolio/strategies"
)

var (
	errInvalidConfigAsset  = errors.New("invalid asset in config")
	errAmbiguousDataSource = errors.New("ambiguous settings received. Only one data type can be set")
	errNoDataSource        = errors.New("no data settings set in config")
	errIntervalUnset       = errors.New("candle interval unset")
	errUnhandledDatatype   = errors.New("unhandled datatype")
	errLiveDataTimeout     = errors.New("no data returned in 5 minutes, shutting down")
	errNilData             = errors.New("nil data received")
	errNilExchange         = errors.New("nil exchange received")
)

// TradeManager is the main holder of all strategy functionality
type TradeManager struct {
	Bot                *Engine
	Datas              data.Holder
	Strategies         []strategies.Handler
	Portfolio          StrategyPortfolioHandler
	Exchange           ExecutionHandler
	Statistic          statistics.Handler
	EventQueue         EventHolder
	Reports            report.Handler
	Warmup             bool
	FactorEngines      map[currency.Pair]*FactorEngine
	DataHistoryManager *DataHistoryManager
	hasHandledEvent    bool
	shutdown           chan struct{}
	cfg                config.Config
	started            int32
	verbose            bool
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
