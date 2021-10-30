package factors

import (
	"fmt"
	"gocryptotrader/data"

	"github.com/shopspring/decimal"
	"gonum.org/v1/gonum/stat"
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
		Time:          kline.LastDate(),
		NLen:          curLength,
		Open:          open,
		High:          high,
		Low:           low,
		Range:         nrange,
		PctChange:     nclose.Sub(open).Div(nclose).Mul(decimal.NewFromFloat(100.0)),
		RangeDivClose: nrange.Div(nclose),
		Slope:         slope,
	}
}

func GetBaseStats(kline *IntervalDataFrame, n int, d data.Handler) *NCalculation {
	if len(kline.Close) <= n {
		return &NCalculation{}
	}
	highBars := kline.High[len(kline.High)-n+1 : len(kline.High)-1]
	lowBars := kline.Low[len(kline.Low)-n+1 : len(kline.Low)-1]
	if len(lowBars) != len(highBars) {
		fmt.Println("error not same amount of bars data")
	}
	high := highBars[0]
	for i := range highBars {
		h := highBars[i]
		if h.GreaterThan(high) {
			high = h
		}
	}
	low := lowBars[0]
	for i := range lowBars {
		l := lowBars[i]
		if l.LessThan(low) {
			low = l
		}
	}
	nopen := kline.Open[len(kline.Open)-n]
	nRange := high.Sub(low)
	nAgoClose := kline.Close[len(kline.Close)-n]
	bar := d.Latest()
	nclose := bar.ClosePrice()
	nPctChg := (nclose.Sub(nAgoClose)).Div(nclose).Mul(decimal.NewFromInt(100))
	nAgo := d.History()[len(d.History())-n]
	return &NCalculation{
		NLen:          n,
		Time:          bar.GetTime(),
		FirstTime:     nAgo.GetTime(),
		NAgoClose:     nAgoClose,
		Range:         nRange,
		PctChange:     nPctChg,
		Open:          nopen,
		High:          high,
		Low:           low,
		Close:         nclose,
		RangeDivClose: nRange.Div(nclose),
		Slope:         decimal.NewFromFloat(getSlope(kline, n)),
	}
}

func getSlope(kline *IntervalDataFrame, n int) float64 {
	var (
		xs      = make([]float64, n)
		ys      = make([]float64, n)
		weights []float64
	)

	for i := range xs {
		f, _ := kline.Close[len(kline.Time)-n+i].Float64()
		xs[i] = float64(i)
		ys[i] = f
		// fmt.Println(ys[i], xs[i])
	}
	origin := false
	_, beta := stat.LinearRegression(xs, ys, weights, origin)
	// r2 := stat.RSquared(xs, ys, weights, alpha, beta)
	// fmt.Printf("alpha=%.6f beta=%.6f R^2=%.6f n=%d\n", alpha, beta, r2, n)
	return beta
}
