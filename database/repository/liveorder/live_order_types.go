package liveorder

import (
	"errors"
)

var (
	// errInvalidInput = errors.New("exchange, base, quote, asset, interval, start & end cannot be empty")
	// errNoCandleData = errors.New("no candle data provided")
	// ErrNoCandleDataFound returns when no candle data is found
	ErrTradeNotFound = errors.New("no trade found")
)

type Status string

// Order side types
const (
	Pending Status = "PENDING"
	Open    Status = "OPEN"
	Closed  Status = "CLOSED"
)

// Item generic candle holder for modelPSQL & modelSQLite
type Details struct {
	ID         int64
	Status     Status
	Pair       string
	OrderType  string
	Exchange   string
	InternalID string
	StrategyID string
}
