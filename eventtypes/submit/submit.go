package submit

import "github.com/shopspring/decimal"

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
