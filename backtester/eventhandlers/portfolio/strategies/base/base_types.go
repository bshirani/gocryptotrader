package base

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/backtester/common"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/exchange"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/compliance"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/holdings"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/positions"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventtypes/fill"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventtypes/order"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventtypes/signal"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/database/repository/livetrade"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
)

var (
	// ErrCustomSettingsUnsupported used when custom settings are found in the start config when they shouldn't be
	ErrCustomSettingsUnsupported = errors.New("custom settings not supported")
	// ErrSimultaneousProcessingNotSupported used when strategy does not support simultaneous processing
	// but start config is set to use it
	ErrSimultaneousProcessingNotSupported = errors.New("does not support simultaneous processing and could not be loaded")
	// ErrStrategyNotFound used when strategy specified in start config does not exist
	ErrStrategyNotFound = errors.New("not found. Please ensure the strategy-settings field 'name' is spelled properly in your .start config")
	// ErrInvalidCustomSettings used when bad custom settings are found in the start config
	ErrInvalidCustomSettings = errors.New("invalid custom settings in config")
	// ErrTooMuchBadData used when there is too much missing data
	ErrTooMuchBadData = errors.New("backtesting cannot continue as there is too much invalid data. Please review your dataset")
)

// Handler contains all functions expected to operate a portfolio manager
type PortfolioHandler interface {
	OnSignal(signal.Event, *exchange.Settings) (*order.Order, error)
	OnFill(fill.Event) (*fill.Fill, error)

	ViewHoldingAtTimePeriod(common.EventHandler) (*holdings.Holding, error)
	UpdateHoldings(common.DataEventHandler) error
	GetPositionForStrategy(string) *positions.Position
	GetTradeForStrategy(string) *livetrade.Details

	GetComplianceManager(string, asset.Item, currency.Pair) (*compliance.Manager, error)

	SetFee(string, asset.Item, currency.Pair, decimal.Decimal)
	GetFee(string, asset.Item, currency.Pair) decimal.Decimal

	Reset()
}
