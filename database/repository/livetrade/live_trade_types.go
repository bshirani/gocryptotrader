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
	DurationMinutes  float64
	EntryOrderID     int
	StrategyName     string
	StrategyID       int
	Name             string
	EntryTime        time.Time
	ExitTime         time.Time
	StopLossPrice    decimal.Decimal
	TakeProfitPrice  decimal.Decimal
	Status           order.Status
	Pair             currency.Pair
	ProfitLossPoints decimal.Decimal
	ProfitLoss       decimal.Decimal
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
