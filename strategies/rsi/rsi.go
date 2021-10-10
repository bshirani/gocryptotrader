package rsi

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	gctcommon "github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/data"
	"github.com/thrasher-corp/gocryptotrader/eventtypes"
	"github.com/thrasher-corp/gocryptotrader/eventtypes/signal"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
	"github.com/thrasher-corp/gocryptotrader/factors"
	"github.com/thrasher-corp/gocryptotrader/strategies/base"
)

const (
	// Name is the strategy name
	Name         = "rsi"
	rsiPeriodKey = "rsi-period"
	rsiLowKey    = "rsi-low"
	rsiHighKey   = "rsi-high"
	description  = `The relative strength index is a technical indicator used in the analysis of financial markets. It is intended to chart the current and historical strength or weakness of a stock or market based on the closing prices of a recent trading period`
)

type IndicatorValues struct {
	Timestamp time.Time
	rsiValue  decimal.Decimal
	maValue   decimal.Decimal
}

// Strategy is an implementation of the Handler interface
type Strategy struct {
	base.Strategy
	rsiPeriod decimal.Decimal
	rsiLow    decimal.Decimal
	rsiHigh   decimal.Decimal
}

// Name returns the name of the strategy
func (s *Strategy) Name() string {
	return Name
}

// Description provides a nice overview of the strategy
// be it definition of terms or to highlight its purpose
func (s *Strategy) Description() string {
	return description
}

// OnData handles a data event and returns what action the strategy believes should occur
// For rsi, this means returning a buy signal when rsi is at or below a certain level, and a
// sell signal when it is at or above a certain level
func (s *Strategy) OnData(d data.Handler, p base.PortfolioHandler, fe *factors.Engine) (signal.Event, error) {

	// fmt.Println("checking strategy", s.Strategy.ID())
	// fmt.Printf("%s %s\n", d.Latest().GetTime(), d.Latest().ClosePrice())

	if d == nil {
		return nil, eventtypes.ErrNilEvent
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

	latestRSIValue := fe.Minute().RSI[len(fe.Minute().RSI)-1]
	fmt.Println("rsi ondata", latestRSIValue)

	if !d.HasDataAtTime(d.Latest().GetTime()) {
		es.SetDirection(eventtypes.MissingData)
		es.SetDecision(signal.DoNothing)
		es.AppendReason(fmt.Sprintf("missing data at %v, cannot perform any actions. RSI %v", d.Latest().GetTime(), latestRSIValue))
		return &es, nil
	}

	// fmt.Println("pair", d.Latest().Pair())
	// fmt.Println(s.Strategy.GetWeight(), m.LastUpdate, m.Close.Last(1))

	pos := p.GetPositionForStrategy(s.Strategy.ID())
	if !pos.Active {
		// check for entry
		if s.Direction() == order.Sell {
			if latestRSIValue.GreaterThanOrEqual(s.rsiHigh) {
				es.SetDecision(signal.Enter)
				es.SetDirection(order.Sell)
			}
		} else if s.Direction() == order.Buy {
			if latestRSIValue.LessThanOrEqual(s.rsiLow) {
				es.SetDecision(signal.Enter)
				es.SetDirection(order.Buy)
			}
		}
	} else {
		// check for exit
		if latestRSIValue.GreaterThanOrEqual(s.rsiHigh) {
			es.SetDecision(signal.Exit)
			es.SetDirection(order.Sell)
		}
	}

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

// SetCustomSettings allows a user to modify the RSI limits in their config
func (s *Strategy) SetCustomSettings(customSettings map[string]interface{}) error {
	for k, v := range customSettings {
		switch k {
		case rsiHighKey:
			rsiHigh, ok := v.(float64)
			if !ok || rsiHigh <= 0 {
				return fmt.Errorf("%w provided rsi-high value could not be parsed: %v", base.ErrInvalidCustomSettings, v)
			}
			s.rsiHigh = decimal.NewFromFloat(rsiHigh)
		case rsiLowKey:
			rsiLow, ok := v.(float64)
			if !ok || rsiLow <= 0 {
				return fmt.Errorf("%w provided rsi-low value could not be parsed: %v", base.ErrInvalidCustomSettings, v)
			}
			s.rsiLow = decimal.NewFromFloat(rsiLow)
		case rsiPeriodKey:
			rsiPeriod, ok := v.(float64)
			if !ok || rsiPeriod <= 0 {
				return fmt.Errorf("%w provided rsi-period value could not be parsed: %v", base.ErrInvalidCustomSettings, v)
			}
			s.rsiPeriod = decimal.NewFromFloat(rsiPeriod)
		default:
			return fmt.Errorf("%w unrecognised custom setting key %v with value %v. Cannot apply", base.ErrInvalidCustomSettings, k, v)
		}
	}

	return nil
}

// SetDefaults sets the custom settings to their default values
func (s *Strategy) SetDefaults() {
	s.rsiHigh = decimal.NewFromInt(70)
	s.rsiLow = decimal.NewFromInt(30)
	s.rsiPeriod = decimal.NewFromInt(14)
}
