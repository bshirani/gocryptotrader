package factors

import "time"

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
