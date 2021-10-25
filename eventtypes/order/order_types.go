package order

import (
	"gocryptotrader/eventtypes"
	"gocryptotrader/eventtypes/event"
	"gocryptotrader/eventtypes/signal"
	"gocryptotrader/exchange/order"

	"github.com/shopspring/decimal"
)

// Order contains all details for an order event
type Order struct {
	event.Base
	ID             string
	Direction      order.Side
	Decision       signal.Decision
	StrategyID     int
	Status         order.Status
	Price          decimal.Decimal
	StopLossPrice  decimal.Decimal
	StopPrice      decimal.Decimal
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
	GetAllocatedFunds() decimal.Decimal
	GetAmount() decimal.Decimal
	GetBuyLimit() decimal.Decimal
	GetDecision() signal.Decision
	GetExchangeFee() decimal.Decimal
	GetID() string
	GetPrice() decimal.Decimal
	GetSellLimit() decimal.Decimal
	GetStatus() order.Status
	GetStopLossPrice() decimal.Decimal
	IsLeveraged() bool
	IsOrder() bool
	SetAmount(decimal.Decimal)
	SetExchangeFee(decimal.Decimal)
	SetID(id string)
	SetPrice(decimal.Decimal)
	SetStopLossPrice(decimal.Decimal)
}
