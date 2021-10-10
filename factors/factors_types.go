package factors

import (
	"errors"

	"github.com/thrasher-corp/gocryptotrader/data"
)

var (
	ErrTooMuchBadData = errors.New("backtesting cannot continue as there is too much invalid data. Please review your dataset")
)

type Engine struct {
	minute *MinuteDataFrame
	daily  *DailyDataFrame
}

type Handler interface {
	Setup()
	Start()
	OnBar(data.Handler)
	Minute() *MinuteDataFrameHandler
	Daily() *DailyDataFrameHandler
}
