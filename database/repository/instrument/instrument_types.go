package instrument

import (
	"errors"
	"time"

	"gocryptotrader/common/cache"
	"gocryptotrader/currency"

	"github.com/shopspring/decimal"
)

var (
	exchangeCache = cache.New(30)
	// ErrNoExchangeFound is a basic predefined error
	ErrNoExchangeFound = errors.New("database.exchange not found")
)

type Details struct {
	ID        int
	Base      currency.Code
	Quote     currency.Code
	MarketCap decimal.Decimal
	DataFrom  time.Time
}
