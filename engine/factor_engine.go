package engine

import (
	"fmt"
	"time"

	"gocryptotrader/common"
	"gocryptotrader/config"
	"gocryptotrader/data"
	"gocryptotrader/factors"

	"github.com/shopspring/decimal"
)

// initialize minute and daily data series here
// load data from cache here
func SetupFactorEngine(eap *ExchangeAssetPairSettings, cfg *config.FactorEngineConfig) (*FactorEngine, error) {
	f := &FactorEngine{}
	p := eap.CurrencyPair
	f.Verbose = cfg.Verbose

	f.Pair = p
	f.minute = &factors.MinuteDataFrame{}
	f.daily = &factors.DailyDataFrame{}

	return f, nil
}

func (f *FactorEngine) Minute() *factors.MinuteDataFrame {
	return f.minute
}

func (f *FactorEngine) Daily() *factors.DailyDataFrame {
	return f.daily
}

func (f *FactorEngine) OnBar(d data.Handler) error {
	// if f.Verbose {
	// 	log.Debugln(log.FactorEngine, "onbar", d.Latest().Pair(), d.Latest().GetTime(), d.Latest().ClosePrice())
	// }
	bar := d.Latest()
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
		// change date after checking for/creating new daily bar
		f.minute.Date = append(f.minute.Date, td)
		f.daily = f.createNewDailyBar(f.minute, f.daily)
	} else {
		f.minute.Date = append(f.minute.Date, td)
	}
	return nil

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

func (f *FactorEngine) createNewDailyBar(m *factors.MinuteDataFrame, d *factors.DailyDataFrame) *factors.DailyDataFrame {
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
func (f *FactorEngine) massageMissingData(data []decimal.Decimal, t time.Time) ([]float64, error) {
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

func (f *FactorEngine) warmup() error {
	// run the catchup process

	// err = f.bot.dataHistoryManager.Stop()
	// if err != nil {
	// 	return err
	// }

	// get latest bars for warmup
	// start := time.Now().Add(time.Minute * -10)
	// end := time.Now()
	//
	// fmt.Println("loading data for", pair.ExchangeName, pair.CurrencyPair)
	// dbData, err := database.LoadData(
	// 	start,
	// 	end,
	// 	time.Minute,
	// 	pair.ExchangeName,
	// 	eventtypes.DataCandle,
	// 	pair.CurrencyPair,
	// 	pair.AssetType)
	//
	// if err != nil {
	// 	fmt.Println("error loading db data", err)
	// }
	// dbData.Load()
	//
	// dbData.Item.RemoveDuplicates()
	// dbData.Item.SortCandlesByTimestamp(false)
	// dbData.RangeHolder, err = gctkline.CalculateCandleDateRanges(
	// 	start,
	// 	end,
	// 	gctkline.Interval(time.Minute),
	// 	0,
	// )
	//
	// f.Datas.SetDataForCurrency(
	// 	pair.ExchangeName,
	// 	pair.AssetType,
	// 	pair.CurrencyPair,
	// 	dbData)
	//
	// //
	// // validate the history is populated with current data
	// //
	// retCandle, _ := candle.Series(pair.ExchangeName,
	// 	pair.CurrencyPair.Base.String(), pair.CurrencyPair.Quote.String(),
	// 	int64(60), string(pair.AssetType), start, end)
	// lc := retCandle.Candles[len(retCandle.Candles)-1].Timestamp
	// t := time.Now()
	// t1 := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
	// t2 := time.Date(lc.Year(), lc.Month(), lc.Day(), lc.Hour(), lc.Minute(), 0, 0, t.Location())
	//
	// if t2 != t1 {
	// 	fmt.Println("sync time is off. History Catchup Failed.", t1, t2)
	// 	// os.Exit(1)
	// }
	//
	// if len(retCandle.Candles) == 0 {
	// 	fmt.Println("No candles returned, History catchup failed. Exiting.")
	// 	os.Exit(1)
	// }
	//
	// // precache the factor engines
	// log.Debugln(log.FactorEngine, "Warming up factor engines...")
	//
	// f.Run()
	//
	// //
	// // validate factor engines are cached
	// //
	// for _, fe := range f.FactorEngines {
	// 	log.Debugf(log.FactorEngine, "fe %v %v", fe.Pair, fe.Minute().LastDate())
	// }

	return nil
}
