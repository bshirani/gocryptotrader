package cancel

import (
	"gocryptotrader/exchange/order"
)

// GetOrder returns the order
func (c *Cancel) GetOrder() *order.Detail {
	return c.Order
}
