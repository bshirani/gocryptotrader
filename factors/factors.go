package factors

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/data"
)

// initialize minute and daily data series here
// load data from cache here
func Setup() (*Engine, error) {
	f := &Engine{}

	f.minute = &MinuteDataFrame{}
	f.daily = &DailyDataFrame{}

	return f, nil
}

func (e *Engine) Start() {
}

func (f *Engine) Minute() *MinuteDataFrame {
	return f.minute
}

func (f *Engine) Daily() *DailyDataFrame {
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

	t := bar.GetTime()
	td := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, t.Nanosecond(), t.Location())

	// logic to create a new daily dataframe
	if len(d.History()) > 1 && td != f.minute.LastDate() {
		f.minute.Date = append(f.minute.Date, td)
		f.daily = f.createNewDailyBar(f.minute, f.daily)
	} else {
		f.minute.Date = append(f.minute.Date, td)
	}

	// change date after checking for/creating new daily bar

	// dataRange := d.StreamClose()
	// var massagedData []float64
	// massagedData, _ = f.massageMissingData(dataRange, d.Latest().GetTime())
	// if err != nil {
	// 	return nil, err
	// }

	// calculate RSI
	// fmt.Println(len(dataRange))
	// rsi := indicators.RSI(f.minute.Close.ToFloats(), 14)
	// latestRSIValue := decimal.NewFromFloat(rsi[len(rsi)-1])
	// f.minute.RSI = append(f.minute.RSI, latestRSIValue)

	// calcMA
	// ma := indicators.MA(f.minute.Close.ToFloats(), 15, indicators.Sma)
	// latestMAValue := decimal.NewFromFloat(ma[len(ma)-1])
	// f.minute.MA = append(f.minute.MA, latestMAValue)

	// i := IndicatorValues{}
	// i.Timestamp = d.Latest().GetTime()
	// i.rsiValue = latestRSIValue
	// i.maValue = latestMAValue
	// s.indicatorValues = append(s.indicatorValues, i)

	// f.minute
}

func (f *Engine) createNewDailyBar(m *MinuteDataFrame, d *DailyDataFrame) *DailyDataFrame {
	// d.Open = append(a.Date, decimal.NewFromFloat(421.0))

	// newDate := m.Date[len(m.Date)-1]
	// var ydayDate time.Time
	//
	// for i := len(m.Close) - 1; i >= 0; i-- {
	// 	if newDate != m.Date[i] {
	// 		ydayDate = m.Date[i]
	// 		break
	// 	}
	// }

	// calculate open here
	d.Open = append(d.Open, decimal.NewFromFloat(421.0))

	// calculate range here
	d.Range = append(d.Range, decimal.NewFromFloat(421.0))

	// fmt.Println("NEW DATE", f.minute.LastDate())
	// // get high/open/low/close
	// f.daily.Open = append(f.daily.Range,.GetOpenPrice()
	// f.daily.Range = append(f.daily.Range, decimal.NewFromFloat(1.0))
	return d
	// update daily dataframe
	// f.daily = append(f.daily,
}

// massageMissingData will replace missing data with the previous candle's data
// this will ensure that RSI can be calculated correctly
// the decision to handle missing data occurs at the strategy level, not all strategies
// may wish to modify data
func (f *Engine) massageMissingData(data []decimal.Decimal, t time.Time) ([]float64, error) {
	var resp []float64
	var missingDataStreak int64
	for i := range data {
		if data[i].IsZero() && i > 14 {
			data[i] = data[i-1]
			missingDataStreak++
		} else {
			missingDataStreak = 0
		}
		if missingDataStreak >= 14 {
			return nil, fmt.Errorf("missing data exceeds RSI period length of %d at %s and will distort results. %w",
				14,
				t.Format(common.SimpleTimeFormat),
				ErrTooMuchBadData)
		}
		d, _ := data[i].Float64()
		resp = append(resp, d)
	}
	return resp, nil
}
