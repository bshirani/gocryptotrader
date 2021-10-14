package liveorder

import (
	"errors"
	"gocryptotrader/currency"
	"gocryptotrader/exchange/order"
)

var (
	// errInvalidInput = errors.New("exchange, base, quote, asset, interval, start & end cannot be empty")
	// errNoCandleData = errors.New("no candle data provided")
	// ErrNoCandleDataFound returns when no candle data is found
	ErrTradeNotFound = errors.New("no trade found")
)

type Status string

// Item generic candle holder for modelPSQL & modelSQLite
type Details struct {
	ID         int64
	Side       order.Side
	Status     order.Status
	Pair       currency.Pair
	OrderType  order.Type
	Exchange   string
	InternalID string
	StrategyID string
}
