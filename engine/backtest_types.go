package engine

import (
	"errors"

	"github.com/thrasher-corp/gocryptotrader/backtester/report"
	portfolio "github.com/thrasher-corp/gocryptotrader/bt_portfolio"
	"github.com/thrasher-corp/gocryptotrader/data"
	"github.com/thrasher-corp/gocryptotrader/engine"
	"github.com/thrasher-corp/gocryptotrader/factors"
	"github.com/thrasher-corp/gocryptotrader/statistics"
	"github.com/thrasher-corp/gocryptotrader/strategies"
)

var (
	errNilConfig           = errors.New("unable to setup backtester with nil config")
	errNilBot              = errors.New("unable to setup backtester without a loaded GoCryptoTrader bot")
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
	FactorEngine       *factors.Engine
	DataHistoryManager *engine.DataHistoryManager
	hasHandledEvent    bool
	shutdown           chan struct{}
	Datas              data.Holder
	Strategies         []strategies.Handler
	Portfolio          portfolio.Handler
	Exchange           ExecutionHandler
	Statistic          statistics.Handler
	EventQueue         EventHolder
	Reports            report.Handler
	catchup            bool
}
