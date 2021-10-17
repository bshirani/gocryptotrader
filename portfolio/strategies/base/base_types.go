package base

import (
	"errors"

	"gocryptotrader/data"
	"gocryptotrader/database/repository/liveorder"
	"gocryptotrader/database/repository/livetrade"
	"gocryptotrader/eventtypes"
	"gocryptotrader/factors"
	"gocryptotrader/portfolio/holdings"
	"gocryptotrader/portfolio/positions"
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
type StrategyPortfolioHandler interface {
	GetLiveMode() bool
	ViewHoldingAtTimePeriod(eventtypes.EventHandler) (*holdings.Holding, error)
	GetPositionForStrategy(string) *positions.Position
	GetOpenOrdersForStrategy(string) []*liveorder.Details
	GetVerbose() bool
	GetTradeForStrategy(string) *livetrade.Details
}

type FactorEngineHandler interface {
	OnBar(data.Handler) error
	Minute() *factors.MinuteDataFrame
	Daily() *factors.DailyDataFrame
}
