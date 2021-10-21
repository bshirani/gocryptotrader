package engine

import (
	"fmt"
	"time"

	"gocryptotrader/common"
	"gocryptotrader/config"
	"gocryptotrader/data"
	"gocryptotrader/factors"
	"gocryptotrader/log"

	"github.com/shopspring/decimal"
)

// initialize minute and daily data series here
// load data from cache here
func SetupFactorEngine(cs *ExchangeAssetPairSettings, cfg *config.FactorEngineConfig) (*FactorEngine, error) {
	f := &FactorEngine{}
	p := cs.CurrencyPair
	f.Verbose = cfg.Verbose

	f.Pair = p
	f.minute = &factors.MinuteDataFrame{}
	f.daily = &factors.DailyDataFrame{}

	// warmup the factor engine

	// load candles from db and convert to events

	// dbData, err := database.LoadData(
	// 	time.Now().Add(time.Minute*-20),
	// 	time.Now(),
	// 	time.Minute,
	// 	cs.ExchangeName,
	// 	0,
	// 	cs.CurrencyPair,
	// 	cs.AssetType)
	//
	// if err != nil {
	// 	fmt.Println("error factor engine warmup", err)
	// }
	// datas := &data.HandlerPerCurrency{}
	// datas.Setup()
	// datas.SetDataForCurrency(cs.ExchangeName, cs.AssetType, cs.CurrencyPair, dbData)
	// dbData.Load()
	// dataHandler := datas.GetDataForCurrency(cs.ExchangeName, cs.AssetType, cs.CurrencyPair)
	//
	// for ev := dataHandler.Next(); ev != nil; dataHandler.Next() {
	// 	f.OnBar(dataHandler)
	// }

	return f, nil
}

func (f *FactorEngine) Minute() *factors.MinuteDataFrame {
	return f.minute
}

func (f *FactorEngine) Daily() *factors.DailyDataFrame {
	return f.daily
}

func (f *FactorEngine) OnBar(d data.Handler) error {
	bar := d.Latest()

	if len(f.minute.Close) > 60 {
		// how much has moved in past hour
		highBars := f.minute.High[len(f.minute.High)-61 : len(f.minute.High)-1]
		lowBars := f.minute.Low[len(f.minute.Low)-61 : len(f.minute.Low)-1]

		if len(lowBars) != len(highBars) {
			fmt.Println("error not same amount of bars data")
		}
		// fmt.Println("have", len(highBars), "bars")

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
		hrRange := high.Sub(low)
		hrRangeRelClose := hrRange.Div(bar.ClosePrice())
		hrRangeRelClose = hrRangeRelClose.Mul(decimal.NewFromInt(100))
		hrAgoClose := f.minute.Close[len(f.minute.Close)-60]
		curClose := bar.ClosePrice()
		hrPctChg := (curClose.Sub(hrAgoClose)).Div(curClose).Mul(decimal.NewFromInt(100))

		f.minute.M60Low = append(f.minute.M60Low, low)
		f.minute.M60High = append(f.minute.M60High, high)
		f.minute.M60Range = append(f.minute.M60Range, hrRange)
		f.minute.M60RangeDivClose = append(f.minute.M60RangeDivClose, hrRange.Div(bar.ClosePrice()))
		f.minute.M60PctChange = append(f.minute.M60PctChange, hrPctChg)

		if f.Verbose {
			f.PrintLast(d)
		}

	} else {
		if f.Verbose {
			log.Debugln(log.FactorEngine, "onbar", bar.Pair(), bar.GetTime(), bar.ClosePrice())
		}
	}

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
	// fmt.Println("close length", d.Latest().Pair(), len(f.minute.Close))
	// fmt.Println("history first", d.History()[0])
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

func (f *FactorEngine) PrintLast(d data.Handler) {
	if len(f.Minute().Close) > 60 {
		hrRangeRelClose := f.minute.M60RangeDivClose[len(f.minute.M60RangeDivClose)-1]
		hrRange := f.minute.M60Range[len(f.minute.M60Range)-1]
		lt := d.Latest()

		if hrRangeRelClose.GreaterThan(decimal.NewFromInt(1)) {
			log.Infof(log.FactorEngine, "%s %s %s %v %v%%", lt.GetTime(), lt.Pair(), "60m range", hrRange, hrRangeRelClose.Round(2))
		} else if !hrRange.IsZero() {
			log.Debugf(log.FactorEngine, "%s %s %s %v %v%%", lt.GetTime(), lt.Pair(), "60m range", hrRange, hrRangeRelClose.Round(2))
		} else {
			log.Errorf(log.FactorEngine, "ZERO %s %s close:%v high:%v low:%v range:%v", lt.GetTime(), lt.Pair(), lt.ClosePrice(), lt.HighPrice(), lt.LowPrice(), hrRange)
		}
	} else {
		log.Debugln(log.FactorEngine, "HH onbar price change", d.Latest().Pair(), d.Latest().GetTime(), d.Latest().ClosePrice())
	}
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
