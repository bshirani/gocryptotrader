package accountlog

import (
	"time"
)

type Details struct {
	ID         int
	USDBalance float64
	XRPBalance float64
	BTCBalance float64
	OpenTrades int
	Timestamp  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
