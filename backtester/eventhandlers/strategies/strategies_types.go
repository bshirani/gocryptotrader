package strategies

import (
	"github.com/thrasher-corp/gocryptotrader/backtester/data"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/strategies/base"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventtypes/signal"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
)

// Handler defines all functions required to run strategies against data events
type Handler interface {
	Direction() order.Side
	Name() string
	Stop()
	GetPosition() (base.Position, error)
	SetPosition(base.Position)
	Description() string
	OnData(data.Handler) (signal.Event, error)
	OnSimultaneousSignals([]data.Handler) ([]signal.Event, error)
	UsingSimultaneousProcessing() bool
	SupportsSimultaneousProcessing() bool
	SetSimultaneousProcessing(bool)
	SetDirection(order.Side)
	SetCustomSettings(map[string]interface{}) error
	SetDefaults()
}
