package factors

import (
	"time"

	"github.com/shopspring/decimal"
)

type DailyDataFrame struct {
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

	RSI   Series
	MA    Series
	Range Series
}

type NSeries struct {
	Low           Series
	High          Series
	Range         Series
	RangeDivClose Series
	PctChange     Series
}

type IntervalDataFrame struct {
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

	N10  *NSeries
	N20  *NSeries
	N30  *NSeries
	N60  *NSeries
	N100 *NSeries
}

type NCalculation struct {
	NLen int
	Time time.Time
	// FirstTime     time.Time
	NAgoClose     decimal.Decimal
	Range         decimal.Decimal
	PctChange     decimal.Decimal
	Close         decimal.Decimal
	Low           decimal.Decimal
	High          decimal.Decimal
	RangeDivClose decimal.Decimal
}

type Calculation struct {
	Time time.Time
	Date time.Time

	High            decimal.Decimal
	Low             decimal.Decimal
	Close           decimal.Decimal
	Open            decimal.Decimal
	Volume          decimal.Decimal
	RSI             decimal.Decimal
	MA              decimal.Decimal
	Past24HourHigh  decimal.Decimal
	Past24HourLow   decimal.Decimal
	CurrentDateOpen decimal.Decimal
	CurrentDateLow  decimal.Decimal
	CurrentDateHigh decimal.Decimal

	N10  *NCalculation
	N20  *NCalculation
	N30  *NCalculation
	N60  *NCalculation
	N100 *NCalculation
}
