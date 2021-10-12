package engine

import (
	"context"
	"errors"

	"gocryptotrader/communications/base"
	"gocryptotrader/currency"
	"gocryptotrader/database"
	exchange "gocryptotrader/exchanges"
	"gocryptotrader/exchanges/asset"
	"gocryptotrader/exchanges/order"
	"gocryptotrader/exchanges/orderbook"
	"gocryptotrader/exchange/ticker"
	"gocryptotrader/portfolio"
)

const (
	// MsgSubSystemStarting message to return when subsystem is starting up
	MsgSubSystemStarting = "starting..."
	// MsgSubSystemStarted message to return when subsystem has started
	MsgSubSystemStarted = "started."
	// MsgSubSystemShuttingDown message to return when a subsystem is shutting down
	MsgSubSystemShuttingDown = "shutting down..."
	// MsgSubSystemShutdown message to return when a subsystem has shutdown
	MsgSubSystemShutdown = "shutdown."
)

var (
	// ErrSubSystemAlreadyStarted message to return when a subsystem is already started
	ErrSubSystemAlreadyStarted = errors.New("subsystem already started")
	// ErrSubSystemNotStarted message to return when subsystem not started
	ErrSubSystemNotStarted = errors.New("subsystem not started")
	// ErrNilSubsystem is returned when a subsystem hasn't had its Setup() func run
	ErrNilSubsystem                 = errors.New("subsystem not setup")
	errNilWaitGroup                 = errors.New("nil wait group received")
	errNilExchangeManager           = errors.New("cannot start with nil exchange manager")
	errNilDatabaseConnectionManager = errors.New("cannot start with nil database connection manager")
	errNilConfig                    = errors.New("received nil config")
)

// iExchangeManager limits exposure of accessible functions to exchange manager
// so that subsystems can use some functionality
type iExchangeManager interface {
	GetExchanges() ([]exchange.IBotExchange, error)
	GetExchangeByName(string) (exchange.IBotExchange, error)
}

// iCommsManager limits exposure of accessible functions to communication manager
type iCommsManager interface {
	PushEvent(evt base.Event)
}

// iOrderManager defines a limited scoped order manager
type iOrderManager interface {
	Exists(*order.Detail) bool
	Add(*order.Detail) error
	Cancel(context.Context, *order.Cancel) error
	GetByExchangeAndID(string, string) (*order.Detail, error)
	UpdateExistingOrder(*order.Detail) error
}

// iPortfolioManager limits exposure of accessible functions to portfolio manager
type iPortfolioManager interface {
	GetPortfolioSummary() portfolio.Summary
	IsWhiteListed(string) bool
	IsExchangeSupported(string, string) bool
}

// iBot limits exposure of accessible functions to engine bot
type iBot interface {
	SetupExchanges() error
}

// iCurrencyPairSyncer defines a limited scoped currency pair syncer
type iCurrencyPairSyncer interface {
	IsRunning() bool
	PrintTickerSummary(*ticker.Price, string, error)
	PrintOrderbookSummary(*orderbook.Base, string, error)
	Update(string, currency.Pair, asset.Item, int, error) error
}

// iDatabaseConnectionManager defines a limited scoped databaseConnectionManager
type iDatabaseConnectionManager interface {
	GetInstance() database.IDatabase
}
