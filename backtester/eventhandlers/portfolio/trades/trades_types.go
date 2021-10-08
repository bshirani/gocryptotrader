package trades

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
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
