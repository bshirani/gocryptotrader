package submit

import (
	"gocryptotrader/exchange/order"

	"github.com/shopspring/decimal"
)

func (s *Submit) GetInternalOrderID() int {
	return s.InternalOrderID
}

func (s *Submit) GetOrderID() string {
	return s.OrderID
}

func (s *Submit) GetIsOrderPlaced() bool {
	return s.IsOrderPlaced
}

func (s *Submit) GetStrategyID() int {
	return s.StrategyID
}

func (s *Submit) GetStopLossOrderID() int {
	return s.StopLossOrderID
}

func (s *Submit) GetPrice() decimal.Decimal {
	return decimal.NewFromFloat(s.Price)
}

func (o *Submit) GetStrategyName() string {
	return o.StrategyName
}

func (o *Submit) SetStrategyName(s string) {
	o.StrategyName = s
}

func (o *Submit) SetStrategyID(s int) {
	o.StrategyID = s
}

func (o *Submit) GetInternalType() order.InternalOrderType {
	return o.InternalType
}
