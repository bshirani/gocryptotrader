package trend

import (
	"fmt"

	"gocryptotrader/common"
	"gocryptotrader/data"
	"gocryptotrader/eventtypes"
	"gocryptotrader/eventtypes/signal"
	"gocryptotrader/portfolio/strategies/base"

	"github.com/shopspring/decimal"
)

const (
	// Name is the strategy name
	Name        = "trend"
	description = `trend exploit`
)

type Strategy struct {
	base.Strategy
	rsiPeriod decimal.Decimal
	rsiLow    decimal.Decimal
	rsiHigh   decimal.Decimal
}

func (s *Strategy) Name() string {
	return Name
}

func (s *Strategy) Description() string {
	return description
}

func (s *Strategy) OnData(d data.Handler, p base.PortfolioHandler, fe base.FactorEngineHandler) (signal.Event, error) {
	if d == nil {
		return nil, eventtypes.ErrNilEvent
	}
	es, err := base.GetBaseData(d)
	if err != nil {
		return nil, err
	}

	// bar := fe.Minute().GetCurrentTime()
	// fmt.Println("straegy on bar", d.Latest().GetTime(), bar)

	// set defaults
	es.SetStrategy(Name)
	es.SetStrategyID(s.ID())
	es.SetPrice(d.Latest().ClosePrice())
	es.SetAmount(decimal.NewFromFloat(1.0))

	offset := d.Offset()

	if offset <= int(s.rsiPeriod.IntPart()) {
		es.SetDecision(signal.DoNothing)
		es.AppendReason("Not enough data for signal generation")
		return &es, nil
	}

	if !d.HasDataAtTime(d.Latest().GetTime()) {
		es.SetDirection(eventtypes.MissingData)
		es.SetDecision(signal.DoNothing)
		es.AppendReason(fmt.Sprintf("missing data at %v, cannot perform any actions. trend", d.Latest().GetTime()))
		return &es, nil
	}

	pos := p.GetPositionForStrategy(s.Strategy.ID())
	if !pos.Active {
		es.SetDecision(signal.Enter)
	} else {
		es.SetDecision(signal.Exit)
	}
	// else {
	// 	es.SetDecision(signal.Exit)
	// 	es.SetDirection(order.Sell)
	// }
	// 	// fmt.Println("check for entry")
	// 	// whats the current date
	// 	// fmt.Println("bar time: ", d.Latest().GetTime())
	// 	// get the current bar from the factor engine
	// 	// bar := fe.Minute().LatestClose()

	m := fe.Minute()
	daily := fe.Daily()
	factors := fmt.Sprintf(
		"%s,%d,%v,%v,%v,%v",
		s.ID(),
		m.GetCurrentTime().Unix(),
		m.GetCurrentDateOpen(),
		m.GetCurrentDateHigh(),
		m.GetCurrentDateLow(),
		len(daily.Open))
	fmt.Println(factors)

	// 	// fmt.Println("bar", m.LatestClose(), m.LastUpdate, d.Latest().GetTime())
	// 	// what was the open of the day
	// } else {
	// 	fmt.Println("check for exit")
	// }

	// no trade
	if es.GetDecision() == "" {
		es.SetDecision(signal.DoNothing)
		es.SetDirection(eventtypes.DoNothing)
	}

	// fmt.Println(s.GetPosition())
	// fmt.Printf("%s@%v@%s now:%v pl:%v\n", t.Direction, t.EntryPrice, t.Timestamp, t.CurrentPrice, t.NetProfit)

	return &es, nil
}

// SupportsSimultaneousProcessing highlights whether the strategy can handle multiple currency calculation
// There is nothing actually stopping this strategy from considering multiple currencies at once
// but for demonstration purposes, this strategy does not
func (s *Strategy) SupportsSimultaneousProcessing() bool {
	return true
}

// OnSimultaneousSignals analyses multiple data points simultaneously, allowing flexibility
// in allowing a strategy to only place an order for X currency if Y currency's price is Z
func (s *Strategy) OnSimultaneousSignals(d []data.Handler, p base.PortfolioHandler, fe base.FactorEngineHandler) ([]signal.Event, error) {
	var resp []signal.Event
	var errs common.Errors
	for i := range d {
		sigEvent, err := s.OnData(d[i], p, fe)
		if err != nil {
			errs = append(errs, fmt.Errorf("%v %v %v %w", d[i].Latest().GetExchange(), d[i].Latest().GetAssetType(), d[i].Latest().Pair(), err))
		} else {
			resp = append(resp, sigEvent)
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}
	return resp, nil
}

// SetCustomSettings allows a user to modify the trend limits in their config
func (s *Strategy) SetCustomSettings(customSettings map[string]interface{}) error {
	return nil
}

// SetDefaults sets the custom settings to their default values
func (s *Strategy) SetDefaults() {
}
