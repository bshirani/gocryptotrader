package trades

import (
	"time"

	"github.com/shopspring/decimal"
	"gocryptotrader/currency"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/exchange/order"
)

type Status string

// Order side types
const (
	Pending Status = "PENDING"
	Open    Status = "OPEN"
	Closed  Status = "CLOSED"
)

// Trade contains trade data for a given time
// for a given exchange asset pair
type Trade struct {
	Status       Status
	Offset       int64
	Item         currency.Code
	Pair         currency.Pair
	EntryPrice   decimal.Decimal
	Direction    order.Side `json:"side"`
	Asset        asset.Item `json:"asset"`
	Exchange     string     `json:"exchange"`
	Timestamp    time.Time  `json:"timestamp"`
	NetProfit    decimal.Decimal
	CurrentPrice decimal.Decimal
}
