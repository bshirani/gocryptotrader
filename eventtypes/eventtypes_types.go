package eventtypes

import (
	"errors"
	"fmt"
	"time"

	"gocryptotrader/currency"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/exchange/kline"
	"gocryptotrader/exchange/order"

	"github.com/shopspring/decimal"
)

// DataTypeToInt converts the config string value into an int
func DataTypeToInt(dataType string) (int64, error) {
	switch dataType {
	case CandleStr:
		return DataCandle, nil
	case TradeStr:
		return DataTrade, nil
	default:
		return 0, fmt.Errorf("unrecognised dataType '%v'", dataType)
	}
}

const (
	// DoNothing is an explicit signal for the backtester to not perform an action
	// based upon indicator results
	DoNothing order.Side = "DO NOTHING"
	// TransferredFunds is a status signal to do nothing
	TransferredFunds order.Side = "TRANSFERRED FUNDS"
	// CouldNotBuy is flagged when a BUY  signal is raised in the strategy/signal phase, but the
	// portfolio manager or exchange cannot place an order
	CouldNotBuy order.Side = "COULD NOT BUY"
	// CouldNotSell is flagged when a SELL  signal is raised in the strategy/signal phase, but the
	// portfolio manager or exchange cannot place an order
	CouldNotSell order.Side = "COULD NOT SELL"
	// MissingData is signalled during the strategy/signal phase when data has been identified as missing
	// No buy or sell events can occur
	MissingData order.Side = "MISSING DATA"
	// CandleStr is a config readable data type to tell the backtester to retrieve candle data
	CandleStr = "candle"
	// TradeStr is a config readable data type to tell the backtester to retrieve trade data
	TradeStr = "trade"
)

// DataCandle is an int64 representation of a candle data type
const (
	DataCandle = iota
	DataTrade
)

var (
	// ErrNilArguments is a common error response to highlight that nils were passed in
	// when they should not have been
	ErrNilArguments = errors.New("received nil argument(s)")
	// ErrNilEvent is a common error for whenever a nil event occurs when it shouldn't have
	ErrNilEvent = errors.New("nil event received")
	// ErrInvalidDataType occurs when an invalid data type is defined in the config
	ErrInvalidDataType = errors.New("invalid datatype received")
)

// EventHandler interface implements required GetTime() & Pair() return
type EventHandler interface {
	GetOffset() int64
	SetOffset(int64)
	IsEvent() bool
	GetTime() time.Time
	Pair() currency.Pair
	GetExchange() string
	GetInterval() kline.Interval
	GetAssetType() asset.Item
	GetReason() string
	AppendReason(string)
	GetStrategyID() int
	SetStrategyID(int)
}

// DataEventHandler interface used for loading and interacting with Data
type DataEventHandler interface {
	EventHandler
	ClosePrice() decimal.Decimal
	HighPrice() decimal.Decimal
	LowPrice() decimal.Decimal
	OpenPrice() decimal.Decimal
}

// Directioner dictates the side of an order
type Directioner interface {
	SetDirection(side order.Side)
	GetDirection() order.Side
}
