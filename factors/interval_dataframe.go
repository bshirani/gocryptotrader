package factors

import (
	"time"

	"github.com/shopspring/decimal"
)

func (d *IntervalDataFrame) LastClose() decimal.Decimal {
	return d.Close[len(d.Close)-1]
}

func (d *IntervalDataFrame) GetCurrentTime() time.Time {
	return d.Time[len(d.Time)-1]
}

func (d *IntervalDataFrame) GetCurrentDateHigh() decimal.Decimal {
	date := d.LastDate()
	var max decimal.Decimal
	max = d.High[len(d.High)-1]

	for i := len(d.Close) - 1; i >= 0; i-- {
		if d.Date[i] == date {
			if d.High[i].GreaterThan(max) {
				max = d.High[i]
			}
		} else {
			break
		}
	}

	return max
}

func (d *IntervalDataFrame) GetCurrentDateLow() decimal.Decimal {
	date := d.LastDate()
	var min decimal.Decimal
	min = d.Low[len(d.Low)-1]

	for i := len(d.Close) - 1; i >= 0; i-- {
		if d.Date[i] == date {
			if d.Low[i].LessThan(min) {
				min = d.Low[i]
			}
		} else {
			break
		}
	}

	return min
}

func (d *IntervalDataFrame) GetCurrentDateOpen() decimal.Decimal {
	date := d.LastDate()
	var open decimal.Decimal
	for i := len(d.Close) - 1; i >= 0; i-- {
		if d.Date[i] == date {
			open = d.Open[i]
		} else {
			break
		}
	}
	return open
}

func (d *IntervalDataFrame) GetCurrentDateLength() int {
	date := d.LastDate()
	length := 0
	for i := len(d.Close) - 1; i >= 0; i-- {
		if d.Date[i] == date {
			length += 1
		} else {
			break
		}
	}
	return length
}

func (d *IntervalDataFrame) LastTime() time.Time {
	return d.Time[len(d.Time)-1]
}

func (d *IntervalDataFrame) LastDate() time.Time {
	return d.Date[len(d.Date)-1]
}

func (d *IntervalDataFrame) Last() Series {
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
