package order

import (
	"gocryptotrader/eventtypes"
	"gocryptotrader/eventtypes/event"
	"gocryptotrader/exchange/order"

	"github.com/shopspring/decimal"
)

// Order contains all details for an order event
type Order struct {
	event.Base
	ID             string
	Direction      order.Side
	StrategyID     int
	Status         order.Status
	Price          decimal.Decimal
	ExchangeFee    decimal.Decimal
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
	SetPrice(decimal.Decimal)
	GetPrice() decimal.Decimal
	SetExchangeFee(decimal.Decimal)
	GetExchangeFee() decimal.Decimal
	SetAmount(decimal.Decimal)
	GetAmount() decimal.Decimal
	IsOrder() bool
	GetStatus() order.Status
	SetID(id string)
	GetID() string
	IsLeveraged() bool
	GetAllocatedFunds() decimal.Decimal
}
