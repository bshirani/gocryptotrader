package engine

import (
	"fmt"
	"gocryptotrader/currency"
	"gocryptotrader/database/repository/candle"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/exchange/kline"
	"gocryptotrader/log"
	"time"
)

func (m *DataHistoryManager) CatchupDays(daysBack int64) error {
	t := time.Now().UTC()
	dayTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	startDate := dayTime.AddDate(0, 0, int(daysBack)*-1)
	for _, p := range m.bot.CurrencySettings {
		// if p.CurrencyPair.Base.String() != "AAVE" {
		// 	continue
		// }
		// for x := startDate; x.Before(dayTime); x = x.AddDate(0, 0, 1) {
		t1 := startDate
		t2 := startDate.AddDate(0, 0, 1)
		candles, _ := candle.Series(p.ExchangeName, p.CurrencyPair.Base.String(), p.CurrencyPair.Quote.String(), 60, p.AssetType.String(), t1, t2)
		if len(candles.Candles) > 1400 {
			// fmt.Printf("%d-%d:%d, ", x.Month(), x.Day(), len(candles.Candles))
			continue
		}
		log.Warnln(log.DataHistory, "Data history manager Syncing Days. Only have:", len(candles.Candles), "bars", p.ExchangeName, p.CurrencyPair, t1, t2)
		m.createCatchupJob(p.ExchangeName, p.AssetType, p.CurrencyPair, t1, t2)
		// time.Sleep(time.Second)
		// }
	}

	// var activePair bool
	// counts, _ := candle.Counts()
	// t := time.Now().UTC()
	// dayTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	//
	// for _, co := range counts {
	// 	p, _ := currency.NewPairFromString(fmt.Sprintf("%s_%s", co.Base, co.Quote))
	// 	for _, cs := range m.bot.CurrencySettings {
	// 		if strings.EqualFold(p.String(), cs.CurrencyPair.String()) {
	// 			// fmt.Println("active", cs.CurrencyPair)
	// 			activePair = true
	// 		}
	// 	}
	// 	if !activePair {
	// 		continue
	// 	}
	//
	// 	if co.Date.Before(time.Now().AddDate(0, 0, -29)) {
	// 		// fmt.Println("skipping")
	// 		continue
	// 	}
	//
	// 	isToday := co.Date.Year() == dayTime.Year() && co.Date.Month() == dayTime.Month() && co.Date.Day() == dayTime.Day()
	// 	if isToday {
	// 		// fmt.Println("is today", co.Count)
	// 		continue
	// 	}
	//
	// 	if co.Count < 1400 {
	// 		t1 := co.Date
	// 		t2 := co.Date.AddDate(0, 0, 1)
	// 		uid, _ := uuid.FromString(co.ExchangeID)
	// 		e, _ := exchangesql.OneByUUID(uid)
	// 		a, _ := asset.New(co.AssetType)
	// 		fmt.Println(p, "ONLY HAS", co.Count, t1, t2, e, a)
	// 		// m.createCatchupJob(e.Name, a, p, t1, t2)
	// 	}
	// 	// else {
	// 	// 	fmt.Printf(".")
	// 	// }
	// 	activePair = false
	// }

	// if m.verbose {
	// 	log.Debugln(log.DataHistory, "catchup today")
	// }
	//
	// for _, p := range m.bot.CurrencySettings {
	// 	t1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).UTC()
	// 	t2 := time.Now().UTC()
	// 	minPast := int(t2.Sub(t1).Minutes())
	// 	candles, _ := candle.Series(p.ExchangeName, p.CurrencyPair.Base.String(), p.CurrencyPair.Quote.String(), 60, p.AssetType.String(), t1, t2)
	// 	missing := minPast - len(candles.Candles)
	// 	if missing < 60 {
	// 		continue
	// 	}
	// 	log.Warnf(log.DataHistory, "Data history manager Syncing More Than 60 minutes of Data %v", MsgSubSystemStarted)
	// 	m.createCatchupJob(p.ExchangeName, p.AssetType, p.CurrencyPair, t1, t2)
	// }

	return nil
}

func (m *DataHistoryManager) createCatchupJob(exchangeName string, a asset.Item, c currency.Pair, start, end time.Time) error {
	startFmt := fmt.Sprintf("%d-%02d-%02d", start.Year(), start.Month(), start.Day())
	endFmt := fmt.Sprintf("%d-%02d-%02d", end.Year(), end.Month(), end.Day())
	name := fmt.Sprintf("%v-%s-%s--%d-catchup", c, startFmt, endFmt, time.Now().Unix())

	dataType := dataHistoryCandleDataType
	// dataType := dataHistoryConvertTradesDataType
	// dataType := dataHistoryTradeDataType

	job := DataHistoryJob{
		Nickname:               name,
		Exchange:               exchangeName,
		Asset:                  a,
		Pair:                   c,
		StartDate:              start,
		EndDate:                end,
		Interval:               kline.Interval(60000000000),
		RunBatchLimit:          10,
		RequestSizeLimit:       999,
		DataType:               dataType,
		MaxRetryAttempts:       1,
		Status:                 dataHistoryStatusActive,
		OverwriteExistingData:  false,
		ConversionInterval:     60000000000,
		DecimalPlaceComparison: 3,
	}
	return m.UpsertJob(&job, true)
}
