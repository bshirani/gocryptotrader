package cancel

import (
	"gocryptotrader/eventtypes"
	"gocryptotrader/eventtypes/event"
	"gocryptotrader/exchange/order"
)

// Cancel is an event that details the events from placing an order
type Cancel struct {
	event.Base
	StrategyID string        `json:"strategy-id"`
	Order      *order.Detail `json:"-"`
}

// Event holds all functions required to handle a Cancel event
type Event interface {
	eventtypes.EventHandler
	GetOrder() *order.Detail
}
