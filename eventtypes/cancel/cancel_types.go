package cancel

import (
	"gocryptotrader/eventtypes"
	"gocryptotrader/eventtypes/event"
	"gocryptotrader/exchange/order"

	"github.com/shopspring/decimal"
)

// Cancel is an event that details the events from placing an order
type Cancel struct {
	event.Base
	Direction           order.Side      `json:"side"`
	StrategyID          string          `json:"strategy-id"`
	Amount              decimal.Decimal `json:"amount"`
	ClosePrice          decimal.Decimal `json:"close-price"`
	VolumeAdjustedPrice decimal.Decimal `json:"volume-adjusted-price"`
	PurchasePrice       decimal.Decimal `json:"purchase-price"`
	Total               decimal.Decimal `json:"total"`
	ExchangeFee         decimal.Decimal `json:"exchange-fee"`
	Slippage            decimal.Decimal `json:"slippage"`
	Order               *order.Detail   `json:"-"`
}

// Event holds all functions required to handle a Cancel event
type Event interface {
	eventtypes.EventHandler
	eventtypes.Directioner

	GetStrategyID() string
	SetStrategyID(string)
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
}
