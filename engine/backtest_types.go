package engine

import (
	"errors"

	"github.com/thrasher-corp/gocryptotrader/backtester/statistics"
	"github.com/thrasher-corp/gocryptotrader/data"
	"github.com/thrasher-corp/gocryptotrader/report"
	"github.com/thrasher-corp/gocryptotrader/strategies"
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

// BackTest is the main holder of all backtesting functionality
type BackTest struct {
	Bot                *Engine
	FactorEngine       *FactorEngine
	DataHistoryManager *DataHistoryManager
	hasHandledEvent    bool
	shutdown           chan struct{}
	Datas              data.Holder
	Strategies         []strategies.Handler
	Portfolio          PortfolioHandler
	Exchange           ExecutionHandler
	Statistic          statistics.Handler
	EventQueue         EventHolder
	Reports            report.Handler
	catchup            bool
}
