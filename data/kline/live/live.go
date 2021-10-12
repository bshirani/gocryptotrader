package live

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gocryptotrader/currency"
	"gocryptotrader/eventtypes"
	exchange "gocryptotrader/exchanges"
	"gocryptotrader/exchanges/asset"
	"gocryptotrader/exchanges/kline"
	"gocryptotrader/exchange/trade"
)

// LoadData retrieves data from a GoCryptoTrader exchange wrapper which calls the exchange's API for the latest interval
func LoadData(ctx context.Context, exch exchange.IBotExchange, dataType int64, interval time.Duration, fPair currency.Pair, a asset.Item) (*kline.Item, error) {
	var candles kline.Item
	var err error
	switch dataType {
	case eventtypes.DataCandle:
		candles, err = exch.GetHistoricCandles(ctx,
			fPair,
			a,
			time.Now().Add(-interval*2), // multiplied by 2 to ensure the latest candle is always included
			time.Now(),
			kline.Interval(interval))
		if err != nil {
			return nil, fmt.Errorf("could not retrieve live candle data for %v %v %v, %v", exch.GetName(), a, fPair, err)
		}

		go func(candles kline.Item) {
			_, err = kline.StoreInDatabase(&candles, true)
			if err != nil {
				fmt.Println("problem saving", err)
			}
		}(candles)
	case eventtypes.DataTrade:
		var trades []trade.Data
		trades, err = exch.GetHistoricTrades(ctx,
			fPair,
			a,
			time.Now().Add(-interval*2), // multiplied by 2 to ensure the latest candle is always included
			time.Now())
		if err != nil {
			return nil, err
		}

		candles, err = trade.ConvertTradesToCandles(kline.Interval(interval), trades...)
		if err != nil {
			return nil, err
		}
		base := exch.GetBase()
		if len(candles.Candles) <= 1 && base.GetSupportedFeatures().RESTCapabilities.TradeHistory {
			trades, err = exch.GetHistoricTrades(ctx,
				fPair,
				a,
				time.Now().Add(-interval),
				time.Now())
			if err != nil {
				return nil, fmt.Errorf("could not retrieve live trade data for %v %v %v, %v", exch.GetName(), a, fPair, err)
			}

			candles, err = trade.ConvertTradesToCandles(kline.Interval(interval), trades...)
			if err != nil {
				return nil, fmt.Errorf("could not convert live trade data to candles for %v %v %v, %v", exch.GetName(), a, fPair, err)
			}
		}
	default:
		return nil, fmt.Errorf("could not retrieve live data for %v %v %v, %w", exch.GetName(), a, fPair, eventtypes.ErrInvalidDataType)
	}
	candles.Exchange = strings.ToLower(exch.GetName())
	return &candles, nil
}
