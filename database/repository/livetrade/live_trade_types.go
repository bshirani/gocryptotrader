package livetrade

import (
	"errors"
	"gocryptotrader/currency"
	"gocryptotrader/exchange/order"

	"github.com/shopspring/decimal"
)

var (
	// errInvalidInput = errors.New("exchange, base, quote, asset, interval, start & end cannot be empty")
	// errNoCandleData = errors.New("no candle data provided")
	// ErrNoCandleDataFound returns when no candle data is found
	ErrTradeNotFound = errors.New("no trade found")
)

// Item generic candle holder for modelPSQL & modelSQLite
type Details struct {
	ID               int64
	Side             order.Side
	Name             string
	StrategyID       string
	EntryOrderID     string
	EntryPrice       decimal.Decimal
	ExitPrice        decimal.Decimal
	StopLossPrice    decimal.Decimal
	Status           order.Status
	Pair             currency.Pair
	ProfitLossPoints decimal.Decimal
}
