package factors

import (
	"time"
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

	N10Low           Series
	N10High          Series
	N10Range         Series
	N10RangeDivClose Series
	N10PctChange     Series

	N20Low           Series
	N20High          Series
	N20Range         Series
	N20RangeDivClose Series
	N20PctChange     Series

	N60Low           Series
	N60High          Series
	N60Range         Series
	N60RangeDivClose Series
	N60PctChange     Series

	N100Low           Series
	N100High          Series
	N100Range         Series
	N100RangeDivClose Series
	N100PctChange     Series
}
