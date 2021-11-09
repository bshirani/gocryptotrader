package livesignal

import (
	"time"
)

type Details struct {
	ID           int
	SignalTime   time.Time
	StrategyName string
	ValidUntil   time.Time
	Prediction   float64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
