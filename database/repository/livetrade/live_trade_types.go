package livetrade

import (
	"errors"
	"gocryptotrader/exchange/order"
)

var (
	// errInvalidInput = errors.New("exchange, base, quote, asset, interval, start & end cannot be empty")
	// errNoCandleData = errors.New("no candle data provided")
	// ErrNoCandleDataFound returns when no candle data is found
	ErrTradeNotFound = errors.New("no trade found")
)

// type Status string
//
// // Order side types
// const (
// 	Pending Status = "PENDING"
// 	Open    Status = "OPEN"
// 	Closed  Status = "CLOSED"
// )

// Item generic candle holder for modelPSQL & modelSQLite
type Details struct {
	ID            int64
	Name          string
	EntryPrice    float64
	ExitPrice     float64
	StopLossPrice float64
	Status        order.Status
	StrategyID    string
	Pair          string
	EntryOrderID  string
}
