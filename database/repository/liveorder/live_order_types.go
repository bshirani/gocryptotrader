package liveorder

import (
	"errors"
	"gocryptotrader/currency"
	"gocryptotrader/exchange/order"
	"time"

	"github.com/shopspring/decimal"
)

var (
	// errInvalidInput = errors.New("exchange, base, quote, asset, interval, start & end cannot be empty")
	// errNoCandleData = errors.New("no candle data provided")
	// ErrNoCandleDataFound returns when no candle data is found
	ErrTradeNotFound = errors.New("no trade found")
)

type Status string

// Item generic candle holder for modelPSQL
type Details struct {
	ID              int
	Side            order.Side
	Amount          decimal.Decimal
	Status          order.Status
	FilledAt        time.Time
	Pair            currency.Pair
	Price           decimal.Decimal
	StopLossPrice   decimal.Decimal
	TakeProfitPrice decimal.Decimal
	OrderType       order.Type
	Exchange        string
	InternalID      string
	StrategyName    string
	StrategyID      int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
