package event

import (
	"time"

	"gocryptotrader/currency"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/exchange/kline"
)

// Base is the underlying event across all actions that occur for the backtester
// Data, fill, order events all contain the base event and store important and
// consistent information
type Base struct {
	Offset       int64          `json:"-"`
	Exchange     string         `json:"exchange"`
	Time         time.Time      `json:"timestamp"`
	Interval     kline.Interval `json:"interval-size"`
	CurrencyPair currency.Pair  `json:"pair"`
	AssetType    asset.Item     `json:"asset"`
	Reason       string         `json:"reason"`
}
