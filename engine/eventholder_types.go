package engine

import "github.com/thrasher-corp/gocryptotrader/eventtypes"

// Holder contains the event queue for backtester processing
type Holder struct {
	Queue []eventtypes.EventHandler
}

// EventHolder interface details what is expected of an event holder to perform
type EventHolder interface {
	Reset()
	AppendEvent(eventtypes.EventHandler)
	NextEvent() eventtypes.EventHandler
}
