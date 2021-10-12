package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gocryptotrader/currency"
	"gocryptotrader/eventtypes"
	exchange "gocryptotrader/exchanges"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/exchange/kline"
	"gocryptotrader/exchange/trade"
)

// LoadData retrieves data from a GoCryptoTrader exchange wrapper which calls the exchange's API
func LoadData(ctx context.Context, dataType int64, startDate, endDate time.Time, interval time.Duration, exch exchange.IBotExchange, fPair currency.Pair, a asset.Item) (*kline.Item, error) {
	var candles kline.Item
	var err error
	switch dataType {
	case eventtypes.DataCandle:
		candles, err = exch.GetHistoricCandlesExtended(ctx,
			fPair,
			a,
			startDate,
			endDate,
			kline.Interval(interval))
		if err != nil {
			return nil, fmt.Errorf("could not retrieve candle data for %v %v %v, %v", exch.GetName(), a, fPair, err)
		}
	case eventtypes.DataTrade:
		var trades []trade.Data
		trades, err = exch.GetHistoricTrades(ctx,
			fPair,
			a,
			startDate,
			endDate)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve trade data for %v %v %v, %v", exch.GetName(), a, fPair, err)
		}

		candles, err = trade.ConvertTradesToCandles(kline.Interval(interval), trades...)
		if err != nil {
			return nil, fmt.Errorf("could not convert trade data to candles for %v %v %v, %v", exch.GetName(), a, fPair, err)
		}
	default:
		return nil, fmt.Errorf("could not retrieve data for %v %v %v, %w", exch.GetName(), a, fPair, eventtypes.ErrInvalidDataType)
	}
	candles.Exchange = strings.ToLower(candles.Exchange)

	return &candles, nil
}
