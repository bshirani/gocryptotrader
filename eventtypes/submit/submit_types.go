package submit

import (
	"gocryptotrader/eventtypes"
	"gocryptotrader/eventtypes/event"
)

// Submit is an event that details the events from placing an order
type Submit struct {
	event.Base
	InternalOrderID string
	OrderID         string
	StrategyID      int
	IsOrderPlaced   bool
}

// Event holds all functions required to handle a fill event
type Event interface {
	eventtypes.EventHandler
	GetInternalOrderID() string
	GetOrderID() int
	GetIsOrderPlaced() bool
}
