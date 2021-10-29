package factors

import (
	"fmt"
	"gocryptotrader/data"
	"math/rand"

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
	nRange := high.Sub(low)
	nAgoClose := kline.Close[len(kline.Close)-n]
	bar := d.Latest()
	curClose := bar.ClosePrice()
	nPctChg := (curClose.Sub(nAgoClose)).Div(curClose).Mul(decimal.NewFromInt(100))
	return &NCalculation{
		NLen:          n,
		Time:          bar.GetTime(),
		NAgoClose:     nAgoClose,
		Range:         nRange,
		PctChange:     nPctChg,
		Close:         curClose,
		Low:           low,
		High:          high,
		RangeDivClose: nRange.Div(curClose),
		Slope:         decimal.NewFromFloat(getSlope(kline, n)),
	}
}

func getSlope(kline *IntervalDataFrame, n int) float64 {
	var (
		xs      = make([]float64, n)
		ys      = make([]float64, n)
		weights []float64
	)

	line := func(x float64) float64 {
		return 1 + 3*x
	}

	for i := range xs {
		xs[i] = float64(i)
		ys[i] = line(xs[i]) + 0.1*rand.NormFloat64()
	}
	origin := false
	alpha, _ := stat.LinearRegression(xs, ys, weights, origin)
	return alpha
}
