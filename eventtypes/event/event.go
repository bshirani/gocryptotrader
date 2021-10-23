package event

import (
	"time"

	"gocryptotrader/currency"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/exchange/kline"
)

// GetOffset returns the offset
func (b *Base) GetOffset() int64 {
	return b.Offset
}

// SetOffset sets the offset
func (b *Base) SetOffset(o int64) {
	b.Offset = o
}

// IsEvent returns whether the event is an event
func (b *Base) IsEvent() bool {
	return true
}

// GetTime returns the time
func (b *Base) GetTime() time.Time {
	return b.Time
}

// Pair returns the currency pair
func (b *Base) Pair() currency.Pair {
	return b.CurrencyPair
}

// GetExchange returns the exchange
func (b *Base) GetExchange() string {
	return b.Exchange
}

// GetAssetType returns the asset type
func (b *Base) GetAssetType() asset.Item {
	return b.AssetType
}

// GetStrategy returns the strategy
func (b *Base) GetStrategyID() int {
	return b.StrategyID
}

// GetStrategy returns the strategy
func (b *Base) SetStrategyID(s int) {
	b.StrategyID = s
}

// GetInterval returns the interval
func (b *Base) GetInterval() kline.Interval {
	return b.Interval
}

// AppendReason adds reasoning for a decision being made
func (b *Base) AppendReason(y string) {
	if b.Reason == "" {
		b.Reason = y
	} else {
		b.Reason = y + ". " + b.Reason
	}
}

// GetReason returns the why
func (b *Base) GetReason() string {
	return b.Reason
}
