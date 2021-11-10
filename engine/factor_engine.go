package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"gocryptotrader/common"
	"gocryptotrader/common/file"
	"gocryptotrader/config"
	"gocryptotrader/data"
	"gocryptotrader/database/repository/livetrade"
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
	f.kline = &factors.IntervalDataFrame{
		N10:  &factors.NSeries{},
		N20:  &factors.NSeries{},
		N30:  &factors.NSeries{},
		N60:  &factors.NSeries{},
		N100: &factors.NSeries{},
	}
	f.daily = &factors.DailyDataFrame{}
	f.calcs = make([]*factors.Calculation, 0)

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

func (f *FactorEngine) Kline() *factors.IntervalDataFrame {
	return f.kline
}

func (f *FactorEngine) Daily() *factors.DailyDataFrame {
	return f.daily
}

func (f *FactorEngine) GetCalculations() []*factors.Calculation {
	return f.calcs
}

func (f *FactorEngine) getFactorCalculations(d data.Handler) *factors.Calculation {
	curDate := factors.GetCurrentDateStats(f.kline, d)
	n10 := factors.GetBaseStats(f.kline, 10, d)
	n20 := factors.GetBaseStats(f.kline, 20, d)
	n30 := factors.GetBaseStats(f.kline, 30, d)
	n60 := factors.GetBaseStats(f.kline, 60, d)
	n100 := factors.GetBaseStats(f.kline, 100, d)
	bar := d.Latest()

	return &factors.Calculation{
		CurrentDate: curDate,
		N10:         n10,
		N20:         n20,
		N30:         n30,
		N60:         n60,
		N100:        n100,
		Close:       bar.ClosePrice(),
		Open:        bar.OpenPrice(),
		High:        bar.HighPrice(),
		Low:         bar.LowPrice(),
		Time:        bar.GetTime(),
	}
}

func (f *FactorEngine) Last() *factors.Calculation {
	return f.calcs[len(f.calcs)-1]
}

func appendNSeries(fs *factors.NSeries, fc *factors.NCalculation) {
	fs.Open = append(fs.Open, fc.Open)
	fs.OpenRel = append(fs.OpenRel, fc.OpenRel)
	fs.Close = append(fs.Close, fc.Close)
	fs.Low = append(fs.Low, fc.Low)
	fs.LowRel = append(fs.LowRel, fc.LowRel)
	fs.High = append(fs.High, fc.High)
	fs.HighRel = append(fs.HighRel, fc.HighRel)
	fs.Range = append(fs.Range, fc.Range)
	fs.RangeRel = append(fs.RangeRel, fc.RangeRel)
	fs.PctChange = append(fs.PctChange, fc.PctChange)
	fs.Slope = append(fs.Slope, fc.Slope)
	fs.SlopeRel = append(fs.SlopeRel, fc.SlopeRel)
}

func (f *FactorEngine) ToQueryParams() map[string]interface{} {
	return f.calcs[len(f.calcs)-1].ToQueryParams()
}

func (f *FactorEngine) OnBar(d data.Handler) error {
	bar := d.Latest()
	t := bar.GetTime()
	td := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, t.Nanosecond(), t.Location())

	if len(d.History()) > 1 && td != f.kline.LastDate() {
		f.kline.Date = append(f.kline.Date, td)
		// f.daily = f.createNewDailyBar(f.kline, f.daily)
	} else {
		f.kline.Date = append(f.kline.Date, td)
	}

	fc := f.getFactorCalculations(d)
	f.calcs = append(f.calcs, fc)

	f.kline.Close = append(f.kline.Close, fc.Close)
	f.kline.Open = append(f.kline.Open, bar.OpenPrice())
	f.kline.High = append(f.kline.High, bar.HighPrice())
	f.kline.Low = append(f.kline.Low, bar.LowPrice())
	f.kline.Time = append(f.kline.Time, bar.GetTime())
	f.kline.LastUpdate = bar.GetTime()

	switch l := len(f.calcs); {
	case l >= 10:
		appendNSeries(f.kline.N10, fc.N10)
		fallthrough
	case l >= 20:
		appendNSeries(f.kline.N20, fc.N20)
		fallthrough
	case l >= 30:
		appendNSeries(f.kline.N30, fc.N30)
		fallthrough
	case l >= 60:
		appendNSeries(f.kline.N60, fc.N60)
		fallthrough
	case l >= 100:
		appendNSeries(f.kline.N100, fc.N100)
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

// func (f *FactorEngine) PrintLast(d data.Handler) {
// 	if len(f.Kline().Close) > 60 {
// 		nRangeRelClose := f.kline.N60RangeRel[len(f.kline.N60RangeRel)-1]
// 		hrRange := f.kline.N60Range[len(f.kline.N60Range)-1]
// 		lt := d.Latest()
//
// 		if nRangeRelClose.GreaterThan(decimal.NewFromInt(1)) {
// 			log.Infof(log.FactorEngine, "%s %s %s %v %v%%", lt.GetTime(), lt.Pair(), "60m range", hrRange, nRangeRelClose.Round(2))
// 		} else if !hrRange.IsZero() {
// 			log.Debugf(log.FactorEngine, "%s %s %s %v %v%%", lt.GetTime(), lt.Pair(), "60m range", hrRange, nRangeRelClose.Round(2))
// 		} else {
// 			log.Errorf(log.FactorEngine, "ZERO %s %s close:%v high:%v low:%v range:%v", lt.GetTime(), lt.Pair(), lt.ClosePrice(), lt.HighPrice(), lt.LowPrice(), hrRange)
// 		}
// 	} else {
// 		log.Debugln(log.FactorEngine, "HH onbar price change", d.Latest().Pair(), d.Latest().GetTime(), d.Latest().ClosePrice())
// 	}
// }

func (f *FactorEngine) createNewDailyBar(m *factors.IntervalDataFrame, d *factors.DailyDataFrame) *factors.DailyDataFrame {
	fmt.Println("creating daily bar!!")

	// whats the open of today

	// d.Open = append(a.Date, decimal.NewFromFloat(421.0))

	curDate := m.Date[len(m.Date)-1]
	var ydayDate time.Time

	for i := len(m.Close) - 1; i >= 0; i-- {
		if curDate != m.Date[i] {
			ydayDate = m.Date[i]
			break
		}
	}
	fmt.Println("ydaydate", ydayDate, "curdate", curDate)

	var open decimal.Decimal
	// calculate open here
	for i := len(m.Open) - 1; i >= 0; i-- {
		if m.Date[i] == curDate {
			continue
		}

		if m.Date[i] == ydayDate {
			open = m.Open[i]
		} else {
			break
		}
	}
	fmt.Println("appending open", open)
	d.Open = append(d.Open, open)

	// calculate range here
	d.Range = append(d.Range, decimal.NewFromFloat(421.0))
	os.Exit(123)

	// fmt.Println("NEW DATE", f.kline.LastDate())
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
	// 	log.Debugf(log.FactorEngine, "fe %v %v", fe.Pair, fe.Kline().LastDate())
	// }

	return nil
}

func (f *FactorEngine) WriteJSON(filepath string) error {
	writer, err := file.Writer(filepath)
	defer func() {
		if writer != nil {
			err = writer.Close()
			if err != nil {
				log.Error(log.Global, err)
			}
		}
	}()
	payload, err := json.MarshalIndent(&f.calcs, "", " ")
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, bytes.NewReader(payload))
	return err
}

type TestRow struct {
	x1 float64
	x2 float64
	y  float64
}

func (f *FactorEngine) WriteCSV(w io.Writer, t []*livetrade.Details) {
	factors.WriteCSV(w, f.calcs, t)
}

// func (f *FactorEngine) WriteCSV(fpath string) error {
// 	// buff := &bytes.Buffer{}
//
// 	structs := []TestRow{
// 		TestRow{x1: 0.20, x2: 0.1, y: 1.0},
// 		TestRow{x1: 0.3, x2: 0.52, y: 1.0},
// 		TestRow{x1: 0.22, x2: 0.01, y: 0.0},
// 	}
//
// 	fmt.Println("writing", fpath, len(f.calcs))
// 	writer, err := file.Writer(fpath)
// 	w := struct2csv.NewWriter(writer)
// 	err = w.WriteStructs(&structs)
// 	// _, err = io.Copy(w, buff)
// 	return err
// }
