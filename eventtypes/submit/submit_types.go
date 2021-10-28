package submit

import (
	"gocryptotrader/eventtypes"
	"gocryptotrader/eventtypes/event"
	"gocryptotrader/exchange/order"

	"github.com/shopspring/decimal"
)

// Submit is an event that details the events from placing an order
type Submit struct {
	event.Base
	InternalOrderID int
	InternalType    order.InternalOrderType
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
	GetInternalOrderType() order.InternalOrderType
	GetOrderID() string
	GetIsOrderPlaced() bool
	GetStopLossOrderID() int
	GetPrice() decimal.Decimal
	GetStrategyName() string
	SetStrategyName(string)
	GetStrategyID() int
	SetStrategyID(int)
}
