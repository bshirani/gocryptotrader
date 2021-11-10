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

	InternalOrderID     int                     `json:"internalOrderId"`
	InternalType        order.InternalOrderType `json:"internalType"`
	Direction           order.Side              `json:"side"`
	Amount              decimal.Decimal         `json:"amount"`
	ClosePrice          decimal.Decimal         `json:"closePrice"`
	VolumeAdjustedPrice decimal.Decimal         `json:"volumeAdjustedPrice"`
	PurchasePrice       decimal.Decimal         `json:"purchasePrice"`
	Total               decimal.Decimal         `json:"total"`
	StopLossPrice       decimal.Decimal         `json:"total"`
	ExchangeFee         decimal.Decimal         `json:"exchangeFee"`
	Slippage            decimal.Decimal         `json:"slippage"`
	StrategyID          int                     `json:"strategyID"`
	StrategyName        string                  `json:"strategyName"`
	Order               *order.Detail           `json:"-"`
	StopLossOrderID     int                     `json:"stopLossOrderId"`
	Prediction          float64                 `json:"prediction"`
}

// Event holds all functions required to handle a fill event
type Event interface {
	eventtypes.EventHandler
	eventtypes.Directioner

	GetOrderID() string
	SetAmount(decimal.Decimal)
	GetAmount() decimal.Decimal
	GetPrediction() float64
	SetPrediction(float64)
	GetClosePrice() decimal.Decimal
	GetVolumeAdjustedPrice() decimal.Decimal
	GetSlippageRate() decimal.Decimal
	GetPurchasePrice() decimal.Decimal
	GetTotal() decimal.Decimal
	GetExchangeFee() decimal.Decimal
	SetExchangeFee(decimal.Decimal)
	GetOrder() *order.Detail
	GetInternalOrderID() int
	GetInternalType() order.InternalOrderType
	GetStopLossOrderID() int
	GetStrategyID() int
	SetStrategyID(int)
	GetStrategyName() string
	SetStrategyName(string)
}
