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
	ID         int
	Side       order.Side
	Status     order.Status
	Pair       currency.Pair
	Price      decimal.Decimal
	StopPrice  decimal.Decimal
	OrderType  order.Type
	Exchange   string
	InternalID string
	StrategyID int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
