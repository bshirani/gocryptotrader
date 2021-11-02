package factors

import (
	"fmt"
	"gocryptotrader/data"

	"github.com/shopspring/decimal"
	"gonum.org/v1/gonum/stat"
)

func (c *Calculation) ToStrings() (list []string) {
	for _, x := range c.N10.CSVRow() {
		list = append(list, x)
	}
	return list
}

func (c *Calculation) CSVHeader() (headers []string) {
	strings := []string{
		"highrel",
		"lowrel",
		"openrel",
		"pctchg",
		"sloperel",
	}
	for _, h := range strings {
		headers = append(headers, fmt.Sprintf("n10_%s", h))
	}
	return headers
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
		NLen:      n,
		Time:      bar.GetTime(),
		FirstTime: nAgo.GetTime(),
		NAgoClose: nAgoClose,
		Range:     nRange,
		PctChange: nPctChg,
		Open:      nopen,
		OpenRel:   nopen.Div(nclose),
		High:      high,
		HighRel:   high.Div(nclose),
		Low:       low,
		LowRel:    low.Div(nclose),
		Close:     nclose,
		RangeRel:  nRange.Div(nclose),
		Slope:     decimal.NewFromFloat(getSlope(kline, n)),
		SlopeRel:  decimal.NewFromFloat(getSlope(kline, n)).Div(nclose),
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
