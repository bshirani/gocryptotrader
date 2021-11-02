package factors

func (c *NCalculation) CSVRow() []string {
	return []string{
		c.HighRel.String(),
		c.LowRel.String(),
		c.OpenRel.String(),
		c.PctChange.String(),
		c.SlopeRel.String(),
	}
}
