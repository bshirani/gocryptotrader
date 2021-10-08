package factors

import (
	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/backtester/data"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/factors/dataframe"
	"github.com/thrasher-corp/gocryptotrader/gct-ta/indicators"
)

// Setup creates a portfolio manager instance and sets private fields
func Setup() (*Engine, error) {
	f := &Engine{}

	// initialize minute and daily data series here
	// load data from cache here

	// values := [][]string{}
	// df := dataframe.LoadRecords(values)

	f.minute = &dataframe.DataFrame{}
	f.daily = &dataframe.DataFrame{}

	// f.minute = dataframe.New(
	// 	series.New([]string{"b", "a"}, series.String, "COL.1"),
	// 	series.New([]int{1, 2}, series.Int, "COL.2"),
	// 	series.New([]float64{3.0, 4.0}, series.Float, "COL.3"),
	// )
	// f.daily = dataframe.New(
	// 	series.New([]string{"b", "a"}, series.String, "COL.1"),
	// 	series.New([]int{1, 2}, series.Int, "COL.2"),
	// 	series.New([]float64{3.0, 4.0}, series.Float, "COL.3"),
	// )

	return f, nil
}

func (e *Engine) Start() {
}

func (f *Engine) Minute() *dataframe.DataFrame {
	return f.minute
}

func (f *Engine) Daily() *dataframe.DataFrame {
	return f.daily
}

func (f *Engine) OnBar(d data.Handler) {
	bar := d.Latest()
	// fmt.Println("bar", bar, f)
	f.minute.Close = append(f.minute.Close, bar.ClosePrice())
	f.minute.Open = append(f.minute.Open, bar.OpenPrice())
	f.minute.High = append(f.minute.High, bar.HighPrice())
	f.minute.Low = append(f.minute.Low, bar.LowPrice())
	// f.minute.Volume = append(f.minute.Volume, bar.GetVolume())
	f.minute.Time = append(f.minute.Time, bar.GetTime())
	f.minute.LastUpdate = bar.GetTime()
	// add bar to minute dataframe

	// dataRange := d.StreamClose()
	// var massagedData []float64
	// massagedData, _ = f.massageMissingData(dataRange, d.Latest().GetTime())
	// if err != nil {
	// 	return nil, err
	// }

	// calculate RSI
	rsi := indicators.RSI(f.minute.Close, 14)
	ma := indicators.MA(f.minute.Close, 14, indicators.Sma)
	latestRSIValue := decimal.NewFromFloat(rsi[len(rsi)-1])
	latestMAValue := decimal.NewFromFloat(ma[len(ma)-1])

	f.minute.RSI = append(f.minute.RSI, latestRSIValue)
	f.minute.MA = append(f.minute.MA, latestMAValue)

	// i := IndicatorValues{}
	// i.Timestamp = d.Latest().GetTime()
	// i.rsiValue = latestRSIValue
	// i.maValue = latestMAValue
	// s.indicatorValues = append(s.indicatorValues, i)

	// f.minute
}

// // massageMissingData will replace missing data with the previous candle's data
// // this will ensure that RSI can be calculated correctly
// // the decision to handle missing data occurs at the strategy level, not all strategies
// // may wish to modify data
// func (f *Engine) massageMissingData(data []decimal.Decimal, t time.Time) ([]float64, error) {
// 	var resp []float64
// 	var missingDataStreak int64
// 	for i := range data {
// 		if data[i].IsZero() && i > 14 {
// 			data[i] = data[i-1]
// 			missingDataStreak++
// 		} else {
// 			missingDataStreak = 0
// 		}
// 		if missingDataStreak >= 14 {
// 			return nil, fmt.Errorf("missing data exceeds RSI period length of %d at %s and will distort results. %w",
// 				14,
// 				t.Format(gctcommon.SimpleTimeFormat),
// 				base.ErrTooMuchBadData)
// 		}
// 		d, _ := data[i].Float64()
// 		resp = append(resp, d)
// 	}
// 	return resp, nil
// }
