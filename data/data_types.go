package data

import (
	"time"

	"github.com/shopspring/decimal"
	"gocryptotrader/currency"
	"gocryptotrader/eventtypes"
	"gocryptotrader/exchanges/asset"
)

// HandlerPerCurrency stores an event handler per exchange asset pair
type HandlerPerCurrency struct {
	data map[string]map[asset.Item]map[currency.Pair]Handler
}

// Holder interface dictates what a data holder is expected to do
type Holder interface {
	Setup()
	SetDataForCurrency(string, asset.Item, currency.Pair, Handler)
	GetAllData() map[string]map[asset.Item]map[currency.Pair]Handler
	GetDataForCurrency(string, asset.Item, currency.Pair) Handler
	Reset()
}

// Base is the base implementation of some interface functions
// where further specific functions are implmented in DataFromKline
type Base struct {
	latest eventtypes.DataEventHandler
	stream []eventtypes.DataEventHandler
	offset int
}

// Handler interface for Loading and Streaming data
type Handler interface {
	Loader
	Streamer
	Reset()
}

// Loader interface for Loading data into backtest supported format
type Loader interface {
	Load() error
}

// Streamer interface handles loading, parsing, distributing BackTest data
type Streamer interface {
	Next() eventtypes.DataEventHandler
	GetStream() []eventtypes.DataEventHandler
	History() []eventtypes.DataEventHandler
	Latest() eventtypes.DataEventHandler
	List() []eventtypes.DataEventHandler
	Offset() int

	StreamOpen() []decimal.Decimal
	StreamHigh() []decimal.Decimal
	StreamLow() []decimal.Decimal
	StreamClose() []decimal.Decimal
	StreamVol() []decimal.Decimal

	HasDataAtTime(time.Time) bool
}
