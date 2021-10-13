package order

import (
	"gocryptotrader/currency"
	"gocryptotrader/exchange/order"

	"github.com/shopspring/decimal"
)

// IsOrder returns whether the event is an order event
func (o *Order) IsOrder() bool {
	return true
}

// SetDirection sets the side of the order
func (o *Order) SetDirection(s order.Side) {
	o.Direction = s
}

// GetDirection returns the side of the order
func (o *Order) GetDirection() order.Side {
	return o.Direction
}

// SetAmount sets the amount
func (o *Order) SetAmount(i decimal.Decimal) {
	o.Amount = i
}

// GetAmount returns the amount
func (o *Order) GetAmount() decimal.Decimal {
	return o.Amount
}

// SetPrice sets the amount
func (o *Order) SetPrice(i decimal.Decimal) {
	o.Price = i
}

// GetPrice returns the amount
func (o *Order) GetPrice() decimal.Decimal {
	return o.Price
}

// SetExchangeFee sets the amount
func (o *Order) SetExchangeFee(i decimal.Decimal) {
	o.ExchangeFee = i
}

// GetExchangeFee returns the amount
func (o *Order) GetExchangeFee() decimal.Decimal {
	return o.ExchangeFee
}

// GetBuyLimit returns the buy limit
func (o *Order) GetBuyLimit() decimal.Decimal {
	return o.BuyLimit
}

// GetSellLimit returns the sell limit
func (o *Order) GetSellLimit() decimal.Decimal {
	return o.SellLimit
}

// Pair returns the currency pair
func (o *Order) Pair() currency.Pair {
	return o.CurrencyPair
}

// GetStatus returns order status
func (o *Order) GetStatus() order.Status {
	return o.Status
}

// SetID sets the order id
func (o *Order) SetID(id string) {
	o.ID = id
}

// GetID returns the ID
func (o *Order) GetID() string {
	return o.ID
}

// // SetID sets the order id
// func (o *Order) SetStrategyID(id string) {
// 	o.StrategyID = id
// }
//
// // GetID returns the ID
// func (o *Order) GetStrategyID() string {
// 	return o.StrategyID
// }

// // SetID sets the order id
// func (o *Order) SetStrategy(name string) {
// 	o.Strategy = name
// }
//
// func (o *Order) GetStrategy() string {
// 	return o.Strategy
// }

// IsLeveraged returns if it is leveraged
func (o *Order) IsLeveraged() bool {
	return o.Leverage.GreaterThan(decimal.NewFromFloat(1))
}

// GetLeverage returns leverage rate
func (o *Order) GetLeverage() decimal.Decimal {
	return o.Leverage
}

// SetLeverage sets leverage
func (o *Order) SetLeverage(l decimal.Decimal) {
	o.Leverage = l
}

// GetAllocatedFunds returns the amount of funds the portfolio manager
// has allocated to this potential position
func (o *Order) GetAllocatedFunds() decimal.Decimal {
	return o.AllocatedFunds
}
