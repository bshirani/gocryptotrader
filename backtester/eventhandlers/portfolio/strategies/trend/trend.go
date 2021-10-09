package trend

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/backtester/common"
	"github.com/thrasher-corp/gocryptotrader/backtester/data"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/factors"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/strategies/base"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventtypes/signal"
	gctcommon "github.com/thrasher-corp/gocryptotrader/common"
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

func (s *Strategy) OnData(d data.Handler, p base.PortfolioHandler, fe *factors.Engine) (signal.Event, error) {
	if d == nil {
		return nil, common.ErrNilEvent
	}
	es, err := s.GetBaseData(d)
	if err != nil {
		return nil, err
	}

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
		es.SetDirection(common.MissingData)
		es.SetDecision(signal.DoNothing)
		es.AppendReason(fmt.Sprintf("missing data at %v, cannot perform any actions. trend", d.Latest().GetTime()))
		return &es, nil
	}

	// pos := p.GetPositionForStrategy(s.Strategy.ID())
	// if !pos.Active {
	// 	// fmt.Println("check for entry")
	// 	// whats the current date
	// 	// fmt.Println("bar time: ", d.Latest().GetTime())
	// 	// get the current bar from the factor engine
	// 	// bar := fe.Minute().LatestClose()
	m := fe.Minute()
	factors := fmt.Sprintf(
		"%d,%v,%v,%v",
		m.GetCurrentTime().Unix(),
		m.GetCurrentDateOpen(),
		m.GetCurrentDateHigh(),
		m.GetCurrentDateLow())
	fmt.Println(factors)
	// 	// fmt.Println("bar", m.LatestClose(), m.LastUpdate, d.Latest().GetTime())
	// 	// what was the open of the day
	// } else {
	// 	fmt.Println("check for exit")
	// }

	// no trade
	if es.GetDecision() == "" {
		es.SetDecision(signal.DoNothing)
		es.SetDirection(common.DoNothing)
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
func (s *Strategy) OnSimultaneousSignals(d []data.Handler, p base.PortfolioHandler, fe *factors.Engine) ([]signal.Event, error) {
	var resp []signal.Event
	var errs gctcommon.Errors
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
