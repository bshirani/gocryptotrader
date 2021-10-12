package cancel

import (
	"gocryptotrader/exchange/order"

	"github.com/shopspring/decimal"
)

// SetDirection sets the direction
func (c *Cancel) SetDirection(s order.Side) {
	c.Direction = s
}

// GetDirection returns the direction
func (c *Cancel) GetDirection() order.Side {
	return c.Direction
}

// SetAmount sets the amount
func (c *Cancel) SetAmount(i decimal.Decimal) {
	c.Amount = i
}

// GetAmount returns the amount
func (c *Cancel) GetAmount() decimal.Decimal {
	return c.Amount
}

// GetClosePrice returns the closing price
func (c *Cancel) GetClosePrice() decimal.Decimal {
	return c.ClosePrice
}

// GetVolumeAdjustedPrice returns the volume adjusted price
func (c *Cancel) GetVolumeAdjustedPrice() decimal.Decimal {
	return c.VolumeAdjustedPrice
}

// GetPurchasePrice returns the purchase price
func (c *Cancel) GetPurchasePrice() decimal.Decimal {
	return c.PurchasePrice
}

// GetTotal returns the total cost
func (c *Cancel) GetTotal() decimal.Decimal {
	return c.Total
}

// GetExchangeFee returns the exchange fee
func (c *Cancel) GetExchangeFee() decimal.Decimal {
	return c.ExchangeFee
}

// SetExchangeFee sets the exchange fee
func (c *Cancel) SetExchangeFee(fee decimal.Decimal) {
	c.ExchangeFee = fee
}

// GetOrder returns the order
func (c *Cancel) GetOrder() *order.Detail {
	return c.Order
}

// GetSlippageRate returns the slippage rate
func (c *Cancel) GetSlippageRate() decimal.Decimal {
	return c.Slippage
}

// GetSlippageRate returns the slippage rate
func (c *Cancel) GetStrategyID() string {
	return c.StrategyID
}

func (c *Cancel) SetStrategyID(sid string) {
	c.StrategyID = sid
}
