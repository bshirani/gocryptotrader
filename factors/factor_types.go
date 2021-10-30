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
	Open          Series
	Close         Series
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
	NLen          int
	Time          time.Time
	FirstTime     time.Time
	NAgoClose     decimal.Decimal
	Range         decimal.Decimal
	PctChange     decimal.Decimal
	Open          decimal.Decimal
	Close         decimal.Decimal
	Low           decimal.Decimal
	High          decimal.Decimal
	RangeDivClose decimal.Decimal
	Slope         decimal.Decimal
}

type Calculation struct {
	Time time.Time
	Date time.Time

	High   decimal.Decimal
	Low    decimal.Decimal
	Close  decimal.Decimal
	Open   decimal.Decimal
	Volume decimal.Decimal

	CurrentDate *NCalculation // current date open/high/low/close/range, etc, updated every interval

	N10  *NCalculation
	N20  *NCalculation
	N30  *NCalculation
	N60  *NCalculation
	N100 *NCalculation
}
