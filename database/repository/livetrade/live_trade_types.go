package livetrade

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

// Item generic candle holder for modelPSQL
type Details struct {
	ID               int
	Side             order.Side
	EntryPrice       decimal.Decimal
	ExitPrice        decimal.Decimal
	Amount           decimal.Decimal
	EntryOrderID     string
	StrategyID       string
	Name             string
	EntryTime        time.Time
	ExitTime         time.Time
	StopLossPrice    decimal.Decimal
	Status           order.Status
	Pair             currency.Pair
	ProfitLossPoints decimal.Decimal
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
