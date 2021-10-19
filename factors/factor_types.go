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

	M60Low           Series
	M60High          Series
	M60Range         Series
	M60RangeDivClose Series
	M60PctChange     Series
}
