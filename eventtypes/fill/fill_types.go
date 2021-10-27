package fill

import (
	"gocryptotrader/eventtypes"
	"gocryptotrader/eventtypes/event"
	"gocryptotrader/exchange/order"

	"github.com/shopspring/decimal"
)

// Fill is an event that details the events from placing an order
type Fill struct {
	event.Base
	OrderID string

	InternalOrderID     int             `json:"internal-order-id"`
	Direction           order.Side      `json:"side"`
	Amount              decimal.Decimal `json:"amount"`
	ClosePrice          decimal.Decimal `json:"close-price"`
	VolumeAdjustedPrice decimal.Decimal `json:"volume-adjusted-price"`
	PurchasePrice       decimal.Decimal `json:"purchase-price"`
	Total               decimal.Decimal `json:"total"`
	StopLossPrice       decimal.Decimal `json:"total"`
	ExchangeFee         decimal.Decimal `json:"exchange-fee"`
	Slippage            decimal.Decimal `json:"slippage"`
	StrategyID          int
	StrategyName        string
	Order               *order.Detail `json:"-"`
	StopLossOrderID     int           `json:"stop-loss-order-id"`
}

// Event holds all functions required to handle a fill event
type Event interface {
	eventtypes.EventHandler
	eventtypes.Directioner

	GetOrderID() string
	SetAmount(decimal.Decimal)
	GetAmount() decimal.Decimal
	GetClosePrice() decimal.Decimal
	GetVolumeAdjustedPrice() decimal.Decimal
	GetSlippageRate() decimal.Decimal
	GetPurchasePrice() decimal.Decimal
	GetTotal() decimal.Decimal
	GetExchangeFee() decimal.Decimal
	SetExchangeFee(decimal.Decimal)
	GetOrder() *order.Detail
	GetInternalOrderID() int
	GetStopLossOrderID() int
	GetStrategyID() int
	SetStrategyID(int)
	GetStrategyName() string
	SetStrategyName(string)
}
