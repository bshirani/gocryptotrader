package submit

import (
	"gocryptotrader/exchange/order"

	"github.com/shopspring/decimal"
)

// SetDirection sets the direction
func (s *Submit) SetDirection(side order.Side) {
	s.Direction = side
}

// GetDirection returns the direction
func (s *Submit) GetDirection() order.Side {
	return s.Direction
}

// SetAmount sets the amount
func (s *Submit) SetAmount(i decimal.Decimal) {
	s.Amount = i
}

// GetAmount returns the amount
func (s *Submit) GetAmount() decimal.Decimal {
	return s.Amount
}

// GetClosePrice returns the closing price
func (s *Submit) GetClosePrice() decimal.Decimal {
	return s.ClosePrice
}

// GetVolumeAdjustedPrice returns the volume adjusted price
func (s *Submit) GetVolumeAdjustedPrice() decimal.Decimal {
	return s.VolumeAdjustedPrice
}

// GetPurchasePrice returns the purchase price
func (s *Submit) GetPurchasePrice() decimal.Decimal {
	return s.PurchasePrice
}

// GetTotal returns the total cost
func (s *Submit) GetTotal() decimal.Decimal {
	return s.Total
}

// GetExchangeFee returns the exchange fee
func (s *Submit) GetExchangeFee() decimal.Decimal {
	return s.ExchangeFee
}

// SetExchangeFee sets the exchange fee
func (s *Submit) SetExchangeFee(fee decimal.Decimal) {
	s.ExchangeFee = fee
}

// GetOrder returns the order
func (s *Submit) GetOrder() *order.Detail {
	return s.Order
}

// GetSlippageRate returns the slippage rate
func (s *Submit) GetSlippageRate() decimal.Decimal {
	return s.Slippage
}

// GetSlippageRate returns the slippage rate
func (s *Submit) GetStrategyID() string {
	return s.StrategyID
}

func (s *Submit) SetStrategyID(sid string) {
	s.StrategyID = sid
}
