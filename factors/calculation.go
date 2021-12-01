package factors

import (
	"fmt"
	"gocryptotrader/data"
	"reflect"
	"strings"

	"github.com/shopspring/decimal"
	"gonum.org/v1/gonum/stat"
)

func (c *Calculation) ToStrings() (list []string) {
	for _, x := range c.N10.CSVRow() {
		list = append(list, x)
	}
	for _, x := range c.N20.CSVRow() {
		list = append(list, x)
	}
	for _, x := range c.N60.CSVRow() {
		list = append(list, x)
	}
	for _, x := range c.N100.CSVRow() {
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
		headers = append(headers, fmt.Sprintf("n20_%s", h))
		headers = append(headers, fmt.Sprintf("n60_%s", h))
		headers = append(headers, fmt.Sprintf("n100_%s", h))
	}
	return headers
}

func (c *Calculation) ToQueryParams() map[string]interface{} {
	return merge(
		c.N10.QueryParams("n10"),
		c.N20.QueryParams("n20"),
		c.N60.QueryParams("n60"),
		c.N100.QueryParams("n100"),
	)
}

func merge(ms ...map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for _, m := range ms {
		for k, v := range m {
			res[k] = v
		}
	}
	return res
}

func getAttr(obj interface{}, fieldName string) reflect.Value {
	pointToStruct := reflect.ValueOf(obj) // addressable
	curStruct := pointToStruct.Elem()
	if curStruct.Kind() != reflect.Struct {
		panic("not struct")
	}
	curField := curStruct.FieldByName(fieldName) // type: reflect.Value
	if !curField.IsValid() {
		panic("not found:" + fieldName)
	}
	return curField
}

func PrintFields(b interface{}) {
	val := reflect.ValueOf(b)
	for i := 0; i < val.Type().NumField(); i++ {
		t := val.Type().Field(i)
		fieldName := t.Name

		if jsonTag := t.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			// check for possible comma as in "...,omitempty"
			var commaIdx int
			if commaIdx = strings.Index(jsonTag, ","); commaIdx < 0 {
				commaIdx = len(jsonTag)
			}
			fieldName = jsonTag[:commaIdx]
		}
		fmt.Println(fieldName)
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
	nAgoTime := kline.Time[len(kline.Time)-n]
	return &NCalculation{
		NLen:      n,
		Time:      bar.GetTime(),
		FirstTime: nAgoTime,
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
