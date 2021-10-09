package factors

import (
	"time"

	"github.com/shopspring/decimal"
)

type MinuteDataFrame struct {
	Close  Series
	Open   Series
	High   Series
	Low    Series
	Volume Series

	Time       []time.Time
	LastUpdate time.Time
	Date       []time.Time
	// Custom user metadata
	Metadata map[string]Series

	RSI             Series
	MA              Series
	Past24HourHigh  Series
	Past24HourLow   Series
	CurrentDateOpen Series
	CurrentDateLow  Series
	CurrentDateHigh Series
}

type MinuteDataFrameHandler interface {
	Last() Series
	LastDate() time.Time
	CurrentDate() time.Time
	CurrentDateHigh() decimal.Decimal
	CurrentDateLow() decimal.Decimal
	CurrentDateOpen() decimal.Decimal
}
