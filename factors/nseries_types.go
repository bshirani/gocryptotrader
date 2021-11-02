package factors

type NSeries struct {
	Open      Series
	OpenRel   Series
	Close     Series
	Low       Series
	LowRel    Series
	High      Series
	HighRel   Series
	Range     Series
	RangeRel  Series
	PctChange Series
	Slope     Series
	SlopeRel  Series
}
