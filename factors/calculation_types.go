package factors

import (
	"time"

	"github.com/shopspring/decimal"
)

type Calculation struct {
	Time time.Time
	Date time.Time

	High  decimal.Decimal
	Low   decimal.Decimal
	Close decimal.Decimal
	Open  decimal.Decimal

	CurrentDate *NCalculation // current date open/high/low/close/range, etc, updated every interval

	N10  *NCalculation
	N20  *NCalculation
	N30  *NCalculation
	N60  *NCalculation
	N100 *NCalculation
}
