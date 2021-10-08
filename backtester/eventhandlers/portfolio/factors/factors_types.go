package factors

import (
	"time"

	"github.com/go-gota/gota/series"
	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/backtester/data"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/factors/dataframe"
)

type Engine struct {
	minute *dataframe.DataFrame
	daily  *dataframe.DataFrame
}

type DataFrameHandler interface {
	Last() series.Series
	LastDate() time.Time
	CurrentDate() time.Time
	CurrentDateHigh() decimal.Decimal
	CurrentDateLow() decimal.Decimal
}

type Handler interface {
	Setup()
	Start()
	OnBar(data.Handler)
	Minute() *DataFrameHandler
	Daily() *DataFrameHandler
}
