package factors

import (
	"github.com/thrasher-corp/gocryptotrader/data"
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
