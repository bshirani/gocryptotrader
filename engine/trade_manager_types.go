package engine

import (
	"errors"
	"gocryptotrader/backtester/statistics"
	"gocryptotrader/config"
	"gocryptotrader/data"
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
	Portfolio          PortfolioHandler
	Exchange           ExecutionHandler
	Statistic          statistics.Handler
	EventQueue         EventHolder
	Reports            report.Handler
	IsLive             bool
	Warmup             bool
	FactorEngine       *FactorEngine       // remove this
	DataHistoryManager *DataHistoryManager // remove this
	hasHandledEvent    bool
	shutdown           chan struct{}
	cfg                config.Config
	started            int32
}
