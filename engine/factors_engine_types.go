package engine

import (
	"errors"

	"gocryptotrader/currency"
	"gocryptotrader/factors"
)

var (
	ErrTooMuchBadData = errors.New("backtesting cannot continue as there is too much invalid data. Please review your dataset")
)

type FactorEngine struct {
	Pair    currency.Pair
	kline   *factors.MinuteDataFrame
	daily   *factors.DailyDataFrame
	Verbose bool
}
