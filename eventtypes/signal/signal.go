package signal

import (
	"gocryptotrader/currency"
	"gocryptotrader/exchange/order"

	"github.com/shopspring/decimal"
)

// IsSignal returns whether the event is a signal type
func (s *Signal) IsSignal() bool {
	return true
}

// SetDirection sets the direction
func (s *Signal) SetDirection(st order.Side) {
	s.Direction = st
}

// GetDirection returns the direction
func (s *Signal) GetDirection() order.Side {
	return s.Direction
}

func (s *Signal) GetStrategyID() int {
	return s.Base.StrategyID
}

func (s *Signal) SetStrategyID(st int) {
	s.Base.StrategyID = st
}

// SetBuyLimit sets the buy limit
func (s *Signal) SetBuyLimit(f decimal.Decimal) {
	s.BuyLimit = f
}

// GetBuyLimit returns the buy limit
func (s *Signal) GetBuyLimit() decimal.Decimal {
	return s.BuyLimit
}

// SetSellLimit sets the sell limit
func (s *Signal) SetSellLimit(f decimal.Decimal) {
	s.SellLimit = f
}

// GetSellLimit returns the sell limit
func (s *Signal) GetSellLimit() decimal.Decimal {
	return s.SellLimit
}

// Pair returns the currency pair
func (s *Signal) Pair() currency.Pair {
	return s.CurrencyPair
}

// GetPrice returns the price
func (s *Signal) GetPrice() decimal.Decimal {
	return s.ClosePrice
}

// GetPrice returns the price
func (s *Signal) GetStopLossPrice() decimal.Decimal {
	return s.StopLossPrice
}

// GetPrice returns the price
func (s *Signal) SetStopLossPrice(stop decimal.Decimal) {
	s.StopLossPrice = stop
}

func (s *Signal) GetDecision() Decision {
	return s.Decision
}

func (s *Signal) SetDecision(d Decision) {
	s.Decision = d
}

// SetPrice sets the price
func (s *Signal) SetPrice(f decimal.Decimal) {
	s.ClosePrice = f
}

func (s *Signal) GetAmount() decimal.Decimal {
	return s.Amount
}

func (s *Signal) SetAmount(d decimal.Decimal) {
	s.Amount = d
}
