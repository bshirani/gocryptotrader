package signal

import (
	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/backtester/common"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventtypes/event"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
)

type Decision string

const (
	Enter     Decision = "ENTER"
	DoNothing Decision = "DO_NOTHING"
	Exit      Decision = "EXIT"
)

// Event handler is used for getting trade signal details
// Example Amount and Price of current candle tick
type Event interface {
	common.EventHandler
	common.Directioner

	GetPrice() decimal.Decimal
	IsSignal() bool
	GetAmount() decimal.Decimal
	SetAmount(decimal.Decimal)
	GetSellLimit() decimal.Decimal
	GetBuyLimit() decimal.Decimal
	GetDecision() Decision
	GetStrategyID() string
	SetStrategyID(string)
	SetDecision(Decision)
}

// Signal contains everything needed for a strategy to raise a signal event
type Signal struct {
	event.Base
	StrategyID string
	Direction  order.Side
	Decision   Decision
	Amount     decimal.Decimal
	OpenPrice  decimal.Decimal
	HighPrice  decimal.Decimal
	LowPrice   decimal.Decimal
	ClosePrice decimal.Decimal
	Volume     decimal.Decimal
	BuyLimit   decimal.Decimal
	SellLimit  decimal.Decimal
}
