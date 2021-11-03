package factors

import (
	"time"

	"github.com/shopspring/decimal"
)

type NCalculation struct {
	NLen      int             `json:"n_len"`
	Time      time.Time       `json:"time"`
	FirstTime time.Time       `json:"first_time"`
	NAgoClose decimal.Decimal `json:"n_ago_close"`
	Range     decimal.Decimal `json:"range"`
	RangeRel  decimal.Decimal `json:"range_rel"`
	PctChange decimal.Decimal `json:"pct_change"`
	Open      decimal.Decimal `json:"open"`
	OpenRel   decimal.Decimal `json:"open_rel"`
	Close     decimal.Decimal `json:"close"`
	Low       decimal.Decimal `json:"low"`
	LowRel    decimal.Decimal `json:"low_rel"`
	High      decimal.Decimal `json:"high"`
	HighRel   decimal.Decimal `json:"high_rel"`
	Slope     decimal.Decimal `json:"slope"`
	SlopeRel  decimal.Decimal `json:"slope_rel"`
}
