package dataframe

import (
	"time"

	"github.com/shopspring/decimal"
)

func (d *DataFrame) LastClose() decimal.Decimal {
	return d.Close[len(d.Close)-1]
}

func (d *DataFrame) LastTime() time.Time {
	return d.Time[len(d.Time)-1]
}

func (d *DataFrame) LastDate() time.Time {
	return d.Date[len(d.Date)-1]
}

func (d *DataFrame) Last() Series {
	res := make([]decimal.Decimal, 5)
	res[0] = decimal.NewFromInt(d.LastTime().Unix())
	res[1] = decimal.NewFromInt(d.LastDate().Unix())
	return res
	// res[0] = decimal.NewFromInt(d.LastTime().Unix())
	// return res

	// res
	// for i,x := range columns {
	// 	res[i] =
	// }
}
