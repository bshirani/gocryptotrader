package factors

import (
	"github.com/thrasher-corp/gocryptotrader/backtester/data"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/factors/dataframe"
)

type Engine struct {
	minute *dataframe.DataFrame
	daily  *dataframe.DataFrame
}

type Handler interface {
	Setup()
	Start()
	OnBar(data.Handler)
	Minute() *dataframe.DataFrame
	Daily() *dataframe.DataFrame
}
