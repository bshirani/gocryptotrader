package instrument

import (
	"errors"
	"time"

	"gocryptotrader/common/cache"

	"github.com/shopspring/decimal"
)

var (
	exchangeCache = cache.New(30)
	// ErrNoExchangeFound is a basic predefined error
	ErrNoExchangeFound = errors.New("database.exchange not found")
)

type Details struct {
	ID                  int
	CMCID               int
	FirstHistoricalData time.Time
	LastHistoricalData  time.Time
	Name                string
	Slug                string
	Symbol              string
	Active              bool
	Status              bool
	// Base                currency.Code
	// Quote               currency.Code
	MarketCap decimal.Decimal
	DataFrom  time.Time
	UpdatedAt time.Time
	CreatedAt time.Time
}
