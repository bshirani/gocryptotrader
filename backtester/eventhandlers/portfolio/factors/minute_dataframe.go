package factors

import (
	"time"

	"github.com/shopspring/decimal"
)

func (d *MinuteDataFrame) LastClose() decimal.Decimal {
	return d.Close[len(d.Close)-1]
}

func (d *MinuteDataFrame) GetCurrentTime() time.Time {
	return d.Time[len(d.Time)-1]
}

func (d *MinuteDataFrame) GetCurrentDateHigh() decimal.Decimal {
	// get all bars for current date
	// start from end and check array until date changes
	// find max of those

	date := d.LastDate()
	var min decimal.Decimal
	var max decimal.Decimal

	for i := len(d.Close) - 1; i >= 0; i-- {
		if d.Date[i] == date {
			if d.Low[i].LessThan(min) {
				min = d.Low[i]
			}
			if d.High[i].GreaterThan(max) {
				max = d.High[i]
			}
		} else {
			break
		}
	}

	return max
}

func (d *MinuteDataFrame) GetCurrentDateLow() decimal.Decimal {
	// get all bars for current date
	// start from end and check array until date changes
	// find max of those

	date := d.LastDate()
	var min decimal.Decimal
	min = d.Low[0]

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

func (d *MinuteDataFrame) GetCurrentDateOpen() decimal.Decimal {
	date := d.LastDate()
	var open decimal.Decimal
	for i := len(d.Close) - 1; i >= 0; i-- {
		if d.Date[i] == date {
			continue
		} else {
			open = d.Open[i] // this is wrong, it should be the one before this
		}
	}
	return open
}

// func (d *MinuteDataFrame) GetDateOpen(date time.Time) decimal.Decimal {
// 	var open decimal.Decimal
// 	for i := len(d.Close) - 1; i >= 0; i-- {
// 		if d.Date[i] == date {
// 			continue
// 		} else {
// 			open = d.Open[i]
// 		}
// 	}
// 	return open
// }

func (d *MinuteDataFrame) LastTime() time.Time {
	return d.Time[len(d.Time)-1]
}

func (d *MinuteDataFrame) LastDate() time.Time {
	return d.Date[len(d.Date)-1]
}

func (d *MinuteDataFrame) Last() Series {
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
