package engine

import (
	"errors"
	"time"

	"gocryptotrader/currency"
	"gocryptotrader/factors"

	"github.com/shopspring/decimal"
)

var (
	ErrTooMuchBadData = errors.New("backtesting cannot continue as there is too much invalid data. Please review your dataset")
)

type FactorEngine struct {
	Pair    currency.Pair
	kline   *factors.IntervalDataFrame
	daily   *factors.DailyDataFrame
	Verbose bool
}

type FactorCalculation struct {
	LastTime      time.Time
	NLen          int
	Range         decimal.Decimal
	High          decimal.Decimal
	Low           decimal.Decimal
	RangeRelClose decimal.Decimal
	NAgoClose     decimal.Decimal
	CurrentClose  decimal.Decimal
	PercentChange decimal.Decimal
}
