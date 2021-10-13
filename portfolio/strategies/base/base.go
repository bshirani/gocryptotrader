package base

import (
	"gocryptotrader/currency"
	"gocryptotrader/data"
	"gocryptotrader/eventtypes"
	"gocryptotrader/eventtypes/event"
	"gocryptotrader/eventtypes/signal"
	"gocryptotrader/exchange/order"

	"github.com/shopspring/decimal"
)

// Strategy is base implementation of the Handler interface
type Strategy struct {
	Name                      string
	ID                        string
	pair                      currency.Pair
	weight                    decimal.Decimal
	direction                 order.Side
	useSimultaneousProcessing bool
	usingExchangeLevelFunding bool
}

func (s *Strategy) GetName() string {
	return s.Name
}

func (s *Strategy) GetID() string {
	return s.ID
}

// GetBaseData returns the non-interface version of the Handler
func GetBaseData(d data.Handler) (signal.Signal, error) {
	if d == nil {
		return signal.Signal{}, eventtypes.ErrNilArguments
	}
	latest := d.Latest()
	if latest == nil {
		return signal.Signal{}, eventtypes.ErrNilEvent
	}
	return signal.Signal{
		Base: event.Base{
			Offset:       latest.GetOffset(),
			Exchange:     latest.GetExchange(),
			Time:         latest.GetTime(),
			CurrencyPair: latest.Pair(),
			AssetType:    latest.GetAssetType(),
			Interval:     latest.GetInterval(),
			Reason:       latest.GetReason(),
		},
		ClosePrice: latest.ClosePrice(),
		HighPrice:  latest.HighPrice(),
		OpenPrice:  latest.OpenPrice(),
		LowPrice:   latest.LowPrice(),
	}, nil
}

func (s *Strategy) SetID(id string) {
	s.ID = id
}

func (s *Strategy) GetPair() currency.Pair {
	return s.pair
}

func (s *Strategy) SetPair(p currency.Pair) {
	s.pair = p
}

func (s *Strategy) SetWeight(d decimal.Decimal) {
	s.weight = d
}

func (s *Strategy) GetWeight() decimal.Decimal {
	return s.weight
}

// UsingSimultaneousProcessing returns whether multiple currencies can be assessed in one go
func (s *Strategy) UsingSimultaneousProcessing() bool {
	return s.useSimultaneousProcessing
}

// SetSimultaneousProcessing sets whether multiple currencies can be assessed in one go
func (s *Strategy) SetSimultaneousProcessing(b bool) {
	s.useSimultaneousProcessing = b
}

// UsingExchangeLevelFunding returns whether funding is based on currency pairs or individual currencies at the exchange level
func (s *Strategy) UsingExchangeLevelFunding() bool {
	return s.usingExchangeLevelFunding
}

// SetExchangeLevelFunding sets whether funding is based on currency pairs or individual currencies at the exchange level
func (s *Strategy) SetExchangeLevelFunding(b bool) {
	s.usingExchangeLevelFunding = b
}

func (s *Strategy) Direction() order.Side {
	return s.direction
}

func (s *Strategy) SetDirection(direction order.Side) {
	s.direction = direction
}

func (s *Strategy) Stop() {
	// fmt.Println("num trades:", len(p.ClosedTrades))
	return
	// for i := range s.indicatorValues {
	// 	x := s.indicatorValues[i]
	// 	fmt.Printf("%d,%s,%s\n", x.Timestamp.Unix(), x.rsiValue, x.maValue)
	// }
}
