# GoCryptoTrader OHLCV support

## Wrapper Methods

Candle retrieval is handled by two methods


GetHistoricCandles which makes a single request to the exchange and follows all exchange limitations
```go
func (b *base) GetHistoricCandles(ctx context.Context, pair currency.Pair, a asset.Item, start, end time.Time, interval kline.Interval) (kline.Item, error) {
	return kline.Item{}, common.ErrFunctionNotSupported
}
```

GetHistoricCandlesExtended that will make multiple requests to an exchange if the requested periods are outside exchange limits
```go
func (b *base) GetHistoricCandlesExtended(ctx context.Context, pair currency.Pair, a asset.Item, start, end time.Time, interval kline.Interval) (kline.Item, error) {
	return kline.Item{}, common.ErrFunctionNotSupported
}
```

both methods return kline.Item{}

```go
// Item holds all the relevant information for internal kline elements
type Item struct {
	Exchange string
	Pair     currency.Pair
	Asset    asset.Item
	Interval Interval
	Candles  []Candle
}

// Candle holds historic rate information.
type Candle struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}
```

### DBSeed helper

A helper tool [cmd/dbseed](../cmd/dbseed/README.md) has been created for assisting with candle data migration

## Exchange status
| Exchange       | Supported   |
|----------------|-------------|
| Binance        | Y           |
| Bitfinex       | Y           |
| Bitflyer       |             |
| Bithumb        | Y           |
| Bitmex         |             |
| Bitstamp       | Y           |
| BTC Markets    | Y           |
| Bittrex        |             |
| BTSE           | Y           |
| Coinbase Pro   | Y           |
| Coinbene       | Y           |
| Coinut         |             |
| Exmo           |             |
| GateIO         | Y           |
| Gemini         |             |
| HitBTC         | Y           |
| Huobi          | Y           |
| FTX            | Y           |
| itBIT          |             |
| Kraken         | Y           |
| lBank          | Y           |
| Localbitcoins  |             |
| Okcoin         | Y           |
| Okex           | Y           |
| Poloniex       | Y           |
| Yobit          |            |
| ZB             | Y           |
