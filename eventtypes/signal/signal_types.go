package signal

import (
	"gocryptotrader/eventtypes"
	"gocryptotrader/eventtypes/event"
	"gocryptotrader/exchange/order"

	"github.com/shopspring/decimal"
)

type Decision string

const (
	Enter     Decision = "ENTER"
	DoNothing Decision = "DO_NOTHING"
	Exit      Decision = "EXIT"
)

// Signal contains everything needed for a strategy to raise a signal event
type Signal struct {
	event.Base
	StopLossPrice decimal.Decimal
	Direction     order.Side
	Decision      Decision
	Amount        decimal.Decimal
	OpenPrice     decimal.Decimal
	HighPrice     decimal.Decimal
	LowPrice      decimal.Decimal
	ClosePrice    decimal.Decimal
	Volume        decimal.Decimal
	BuyLimit      decimal.Decimal
	SellLimit     decimal.Decimal
}

// Event handler is used for getting trade signal details
// Example Amount and Price of current candle tick
type Event interface {
	eventtypes.EventHandler
	eventtypes.Directioner

	GetPrice() decimal.Decimal
	IsSignal() bool
	GetAmount() decimal.Decimal
	SetAmount(decimal.Decimal)
	GetSellLimit() decimal.Decimal
	GetBuyLimit() decimal.Decimal
	GetDecision() Decision
	SetDecision(Decision)
	SetStopLossPrice(decimal.Decimal)
	GetStopLossPrice() decimal.Decimal
}
