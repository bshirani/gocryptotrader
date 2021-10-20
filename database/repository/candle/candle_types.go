package candle

import (
	"errors"
	"time"
)

var (
	errInvalidInput = errors.New("exchange, base, quote, asset, interval, start & end cannot be empty")
	errNoCandleData = errors.New("no candle data provided")
	// ErrNoCandleDataFound returns when no candle data is found
	ErrNoCandleDataFound = errors.New("no candle data found")
)

// Item generic candle holder for modelPSQL
type Item struct {
	ID         int64
	ExchangeID string
	Base       string
	Quote      string
	Interval   int64
	Asset      string
	Candles    []Candle
}

// Candle holds each interval
type Candle struct {
	Timestamp        time.Time
	Open             float64
	High             float64
	Low              float64
	Close            float64
	Volume           float64
	SourceJobID      string
	ValidationJobID  string
	ValidationIssues string
}

type PairCandleCountResponse struct {
	ExchangeID string    `boil:"exchange_id"`
	AssetType  string    `boil:"asset_type"`
	Date       time.Time `boil:"date"`
	Base       string    `boil:"base"`
	Quote      string    `boil:"quote"`
	Count      int64     `boil:"count"`
}
