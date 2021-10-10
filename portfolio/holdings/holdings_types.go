package holdings

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"gocryptotrader/currency"
	"gocryptotrader/exchanges/asset"
)

// ErrInitialFundsZero is an error when initial funds are zero or less
var ErrInitialFundsZero = errors.New("initial funds < 0")

// Holding contains pricing statistics for a given time
// for a given exchange asset pair
type Holding struct {
	Offset            int64
	Item              currency.Code
	Pair              currency.Pair
	Asset             asset.Item      `json:"asset"`
	Exchange          string          `json:"exchange"`
	Timestamp         time.Time       `json:"timestamp"`
	BaseInitialFunds  decimal.Decimal `json:"base-initial-funds"`
	BaseSize          decimal.Decimal `json:"base-size"`
	BaseValue         decimal.Decimal `json:"base-value"`
	QuoteInitialFunds decimal.Decimal `json:"quote-initial-funds"`
	TotalInitialValue decimal.Decimal `json:"total-initial-value"`
	QuoteSize         decimal.Decimal `json:"quote-size"`
	SoldAmount        decimal.Decimal `json:"sold-amount"`
	SoldValue         decimal.Decimal `json:"sold-value"`
	BoughtAmount      decimal.Decimal `json:"bought-amount"`
	BoughtValue       decimal.Decimal `json:"bought-value"`

	TotalValueDifference      decimal.Decimal
	ChangeInTotalValuePercent decimal.Decimal
	BoughtValueDifference     decimal.Decimal
	SoldValueDifference       decimal.Decimal
	PositionsValueDifference  decimal.Decimal

	TotalValue                   decimal.Decimal `json:"total-value"`
	TotalFees                    decimal.Decimal `json:"total-fees"`
	TotalValueLostToVolumeSizing decimal.Decimal `json:"total-value-lost-to-volume-sizing"`
	TotalValueLostToSlippage     decimal.Decimal `json:"total-value-lost-to-slippage"`
	TotalValueLost               decimal.Decimal `json:"total-value-lost"`

	RiskFreeRate decimal.Decimal `json:"risk-free-rate"`
}
