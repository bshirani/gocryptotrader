package order

import (
	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/eventtypes"
	"github.com/thrasher-corp/gocryptotrader/eventtypes/event"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
)

// Order contains all details for an order event
type Order struct {
	event.Base
	ID             string
	Direction      order.Side
	StrategyID     string
	Status         order.Status
	Price          decimal.Decimal
	Amount         decimal.Decimal
	OrderType      order.Type
	Leverage       decimal.Decimal
	AllocatedFunds decimal.Decimal
	BuyLimit       decimal.Decimal
	SellLimit      decimal.Decimal
}

// Event inherits common event interfaces along with extra functions related to handling orders
type Event interface {
	eventtypes.EventHandler
	eventtypes.Directioner
	GetBuyLimit() decimal.Decimal
	GetSellLimit() decimal.Decimal
	SetAmount(decimal.Decimal)
	GetAmount() decimal.Decimal
	IsOrder() bool
	GetStatus() order.Status
	SetID(id string)
	GetID() string
	SetStrategyID(id string)
	GetStrategyID() string
	IsLeveraged() bool
	GetAllocatedFunds() decimal.Decimal
}
