package factors

import (
	"gocryptotrader/data"

	"github.com/shopspring/decimal"
)

type Series []decimal.Decimal

func (s Series) Values() []decimal.Decimal {
	return s
}

func (s Series) ToFloats() []float64 {
	x := make([]float64, len(s))
	for i, si := range s {
		x[i], _ = si.Float64()
	}
	return x
}

func (s Series) Last(position int) decimal.Decimal {
	return s[len(s)-1-position]
}

func (s Series) LastValues(size int) []decimal.Decimal {
	if l := len(s); l > size {
		return s[l-size:]
	}
	return s
}

func (s Series) Crossover(ref Series) bool {
	return s.Last(0).GreaterThan(ref.Last(0)) && s.Last(1).LessThanOrEqual(ref.Last(1))
}

func (s Series) Crossunder(ref Series) bool {
	return s.Last(0).LessThanOrEqual(ref.Last(0)) && s.Last(1).GreaterThan(ref.Last(1))
}

func GetCurrentDateStats(kline *IntervalDataFrame, d data.Handler) *NCalculation {
	if len(kline.Close) <= 1 {
		return &NCalculation{}
	}
	// NLen:          n,
	// Time:          bar.GetTime(),
	// FirstTime:     nAgo.GetTime(),
	// NAgoClose:     nAgoClose,
	// Range:         nRange,
	// PctChange:     nPctChg,
	// Slope:         decimal.NewFromFloat(getSlope(kline, n)),
	high := kline.GetCurrentDateHigh()
	low := kline.GetCurrentDateLow()
	nrange := high.Sub(low)
	open := kline.GetCurrentDateOpen()
	nclose := d.Latest().ClosePrice()
	curLength := kline.GetCurrentDateLength()
	var slope decimal.Decimal

	if curLength > 1 {
		slope = decimal.NewFromFloat(getSlope(kline, curLength))
	} else {
		slope = decimal.NewFromFloat(0.0)
	}

	return &NCalculation{
		Time:      kline.LastDate(),
		NLen:      curLength,
		Open:      open,
		OpenRel:   open.Div(nclose),
		High:      high,
		HighRel:   high.Div(nclose),
		Low:       low,
		LowRel:    low.Div(nclose),
		Range:     nrange,
		PctChange: nclose.Sub(open).Div(nclose).Mul(decimal.NewFromFloat(100.0)),
		RangeRel:  nrange.Div(nclose),
		Slope:     slope,
		SlopeRel:  slope.Div(nclose),
	}
}
