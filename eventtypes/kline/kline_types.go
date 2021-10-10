package kline

import (
	"github.com/shopspring/decimal"
	"gocryptotrader/eventtypes/event"
)

// Kline holds kline data and an event to be processed as
// a eventtypes.DataEventHandler type
type Kline struct {
	event.Base
	Open             decimal.Decimal
	Close            decimal.Decimal
	Low              decimal.Decimal
	High             decimal.Decimal
	Volume           decimal.Decimal
	ValidationIssues string
}
