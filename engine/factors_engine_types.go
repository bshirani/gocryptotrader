package engine

import (
	"errors"

	"github.com/thrasher-corp/gocryptotrader/factors"
)

var (
	ErrTooMuchBadData = errors.New("backtesting cannot continue as there is too much invalid data. Please review your dataset")
)

type FactorEngine struct {
	minute *factors.MinuteDataFrame
	daily  *factors.DailyDataFrame
}
