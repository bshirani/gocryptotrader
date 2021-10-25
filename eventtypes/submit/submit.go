package submit

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
