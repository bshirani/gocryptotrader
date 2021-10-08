package positions

import "github.com/shopspring/decimal"

type Position struct {
	id     int64
	Active bool
	Amount decimal.Decimal
}
