package livetrade

import (
	"errors"
)

var (
	// errInvalidInput = errors.New("exchange, base, quote, asset, interval, start & end cannot be empty")
	// errNoCandleData = errors.New("no candle data provided")
	// ErrNoCandleDataFound returns when no candle data is found
	ErrTradeNotFound = errors.New("no trade found")
)

// Item generic candle holder for modelPSQL & modelSQLite
type Details struct {
	ID         string
	Name       string
	EntryPrice float64
	ExitPrice  float64
	StopPrice  float64
}
