package kline

import "github.com/shopspring/decimal"

// // ClosePrice returns the closing price of a kline
// func (k *Kline) GetStrategyID() int {
// 	return 0
// }
//
// func (k *Kline) SetStrategyID(s int) {
// 	return
// }

// ClosePrice returns the closing price of a kline
func (k *Kline) ClosePrice() decimal.Decimal {
	return k.Close
}

// HighPrice returns the high price of a kline
func (k *Kline) HighPrice() decimal.Decimal {
	return k.High
}

// LowPrice returns the low price of a kline
func (k *Kline) LowPrice() decimal.Decimal {
	return k.Low
}

// OpenPrice returns the open price of a kline
func (k *Kline) OpenPrice() decimal.Decimal {
	return k.Open
}
