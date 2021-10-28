package fill

import (
	"gocryptotrader/exchange/order"

	"github.com/shopspring/decimal"
)

// SetDirection sets the direction
func (f *Fill) SetDirection(s order.Side) {
	f.Direction = s
}

// GetDirection returns the direction
func (f *Fill) GetDirection() order.Side {
	return f.Direction
}

// SetAmount sets the amount
func (f *Fill) SetAmount(i decimal.Decimal) {
	f.Amount = i
}

// GetAmount returns the amount
func (f *Fill) GetAmount() decimal.Decimal {
	return f.Amount
}

// GetClosePrice returns the closing price
func (f *Fill) GetClosePrice() decimal.Decimal {
	return f.ClosePrice
}

// GetVolumeAdjustedPrice returns the volume adjusted price
func (f *Fill) GetVolumeAdjustedPrice() decimal.Decimal {
	return f.VolumeAdjustedPrice
}

// GetPurchasePrice returns the purchase price
func (f *Fill) GetPurchasePrice() decimal.Decimal {
	return f.PurchasePrice
}

// GetTotal returns the total cost
func (f *Fill) GetTotal() decimal.Decimal {
	return f.Total
}

// GetExchangeFee returns the exchange fee
func (f *Fill) GetExchangeFee() decimal.Decimal {
	return f.ExchangeFee
}

// SetExchangeFee sets the exchange fee
func (f *Fill) SetExchangeFee(fee decimal.Decimal) {
	f.ExchangeFee = fee
}

// GetOrder returns the order
func (f *Fill) GetOrder() *order.Detail {
	return f.Order
}

// GetOrder returns the order
func (f *Fill) GetOrderID() string {
	return f.OrderID
}

// GetOrder returns the order
func (f *Fill) GetInternalOrderID() int {
	return f.InternalOrderID
}

// GetOrder returns the order
func (f *Fill) GetInternalType() order.InternalOrderType {
	return f.InternalType
}

// GetOrder returns the order
func (f *Fill) GetStopLossOrderID() int {
	return f.StopLossOrderID
}

// GetSlippageRate returns the slippage rate
func (f *Fill) GetSlippageRate() decimal.Decimal {
	return f.Slippage
}

// GetStrategy returns the strategy
func (b *Fill) GetStrategyID() int {
	return b.StrategyID
}

// GetStrategy returns the strategy
func (b *Fill) SetStrategyID(s int) {
	b.StrategyID = s
}

func (o *Fill) GetStrategyName() string {
	return o.StrategyName
}

func (o *Fill) SetStrategyName(s string) {
	o.StrategyName = s
}
