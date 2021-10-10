package base

import (
	"errors"

	"github.com/thrasher-corp/gocryptotrader/data"
	"github.com/thrasher-corp/gocryptotrader/database/repository/livetrade"
	"github.com/thrasher-corp/gocryptotrader/eventtypes"
	"github.com/thrasher-corp/gocryptotrader/factors"
	"github.com/thrasher-corp/gocryptotrader/portfolio/holdings"
	"github.com/thrasher-corp/gocryptotrader/portfolio/positions"
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
	ViewHoldingAtTimePeriod(eventtypes.EventHandler) (*holdings.Holding, error)
	GetPositionForStrategy(string) *positions.Position
	GetTradeForStrategy(string) *livetrade.Details
}

type FactorEngineHandler interface {
	Start()
	OnBar(data.Handler)
	Minute() *factors.MinuteDataFrame
	Daily() *factors.DailyDataFrame
}
