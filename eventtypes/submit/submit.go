package submit

func (s *Submit) GetInternalOrderID() string {
	return s.InternalOrderID
}

func (s *Submit) GetOrderID() string {
	return s.OrderID
}

func (s *Submit) GetIsOrderPlaced() bool {
	return s.IsOrderPlaced
}

func (s *Submit) GetStrategyID() string {
	return s.StrategyID
}
