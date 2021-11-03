package factors

import "fmt"

func (c *NCalculation) CSVRow() []string {
	return []string{
		c.HighRel.String(),
		c.LowRel.String(),
		c.OpenRel.String(),
		c.PctChange.String(),
		c.SlopeRel.String(),
	}
}

func (c *NCalculation) QueryParams(prefix string) map[string]interface{} {
	params := make(map[string]interface{})
	params[fmt.Sprintf("%s_sloperel", prefix)], _ = c.SlopeRel.Float64()
	params[fmt.Sprintf("%s_rangerel", prefix)], _ = c.RangeRel.Float64()
	params[fmt.Sprintf("%s_pctchg", prefix)], _ = c.PctChange.Float64()
	params[fmt.Sprintf("%s_highrel", prefix)], _ = c.HighRel.Float64()
	params[fmt.Sprintf("%s_lowrel", prefix)], _ = c.LowRel.Float64()
	params[fmt.Sprintf("%s_openrel", prefix)], _ = c.OpenRel.Float64()
	return params
}
