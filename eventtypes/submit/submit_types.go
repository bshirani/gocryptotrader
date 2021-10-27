package submit

import (
	"gocryptotrader/eventtypes"
	"gocryptotrader/eventtypes/event"

	"github.com/shopspring/decimal"
)

// Submit is an event that details the events from placing an order
type Submit struct {
	event.Base
	InternalOrderID int
	OrderID         string
	StrategyID      int
	StrategyName    string
	IsOrderPlaced   bool
	FullyMatched    bool
	StopLossOrderID int
	Price           float64
}

// Event holds all functions required to handle a fill event
type Event interface {
	eventtypes.EventHandler
	GetInternalOrderID() int
	GetOrderID() string
	GetIsOrderPlaced() bool
	GetStopLossOrderID() int
	GetPrice() decimal.Decimal
	GetStrategyName() string
	SetStrategyName(string)
	GetStrategyID() int
	SetStrategyID(int)
}
