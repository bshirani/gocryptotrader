package factors

import (
	"time"

	"github.com/shopspring/decimal"
)

type NCalculation struct {
	NLen      int
	Time      time.Time
	FirstTime time.Time
	NAgoClose decimal.Decimal
	Range     decimal.Decimal
	PctChange decimal.Decimal
	Open      decimal.Decimal
	OpenRel   decimal.Decimal
	Close     decimal.Decimal
	Low       decimal.Decimal
	LowRel    decimal.Decimal
	High      decimal.Decimal
	HighRel   decimal.Decimal
	RangeRel  decimal.Decimal
	Slope     decimal.Decimal
	SlopeRel  decimal.Decimal
}
