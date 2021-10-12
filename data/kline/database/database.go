package database

import (
	"fmt"
	"strings"
	"time"

	"gocryptotrader/currency"
	"gocryptotrader/data/kline"
	"gocryptotrader/eventtypes"
	"gocryptotrader/exchange/asset"
	gctkline "gocryptotrader/exchange/kline"
	"gocryptotrader/exchange/trade"
	"gocryptotrader/log"
)

// LoadData retrieves data from an existing database using GoCryptoTrader's database handling implementation
func LoadData(startDate, endDate time.Time, interval time.Duration, exchangeName string, dataType int64, fPair currency.Pair, a asset.Item) (*kline.DataFromKline, error) {
	resp := &kline.DataFromKline{}
	switch dataType {
	case eventtypes.DataCandle:
		klineItem, err := getCandleDatabaseData(
			startDate,
			endDate,
			interval,
			exchangeName,
			fPair,
			a)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve database candle data for %v %v %v, %v", exchangeName, a, fPair, err)
		}
		resp.Item = klineItem
		for i := range klineItem.Candles {
			if klineItem.Candles[i].ValidationIssues != "" {
				log.Warnf(log.TradeManager, "candle validation issue for %v %v %v: %v", klineItem.Exchange, klineItem.Asset, klineItem.Pair, klineItem.Candles[i].ValidationIssues)
			}
		}
	case eventtypes.DataTrade:
		trades, err := trade.GetTradesInRange(
			exchangeName,
			a.String(),
			fPair.Base.String(),
			fPair.Quote.String(),
			startDate,
			endDate)
		if err != nil {
			return nil, err
		}
		klineItem, err := trade.ConvertTradesToCandles(
			gctkline.Interval(interval),
			trades...)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve database trade data for %v %v %v, %v", exchangeName, a, fPair, err)
		}
		resp.Item = klineItem
	default:
		return nil, fmt.Errorf("could not retrieve database data for %v %v %v, %w", exchangeName, a, fPair, eventtypes.ErrInvalidDataType)
	}
	resp.Item.Exchange = strings.ToLower(resp.Item.Exchange)

	return resp, nil
}

func getCandleDatabaseData(startDate, endDate time.Time, interval time.Duration, exchangeName string, fPair currency.Pair, a asset.Item) (gctkline.Item, error) {
	return gctkline.LoadFromDatabase(
		exchangeName,
		fPair,
		a,
		gctkline.Interval(interval),
		startDate,
		endDate)
}
