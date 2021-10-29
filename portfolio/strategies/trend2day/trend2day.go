package trend2day

import (
	"fmt"
	"os"

	"gocryptotrader/common"
	"gocryptotrader/data"
	"gocryptotrader/eventtypes"
	"gocryptotrader/eventtypes/signal"
	"gocryptotrader/exchange/order"
	"gocryptotrader/portfolio/strategies/base"

	"github.com/shopspring/decimal"
)

const (
	// Name is the strategy name
	Name        = "trend2day"
	description = `trend exploit`
)

type Strategy struct {
	base.Strategy
	rsiPeriod decimal.Decimal
}

func (s *Strategy) Name() string {
	return Name
}

func (s *Strategy) Description() string {
	return description
}

func (s *Strategy) OnData(d data.Handler, p base.StrategyPortfolioHandler, fe base.FactorEngineHandler) (signal.Event, error) {

	if d == nil {
		return nil, eventtypes.ErrNilEvent
	}
	es, err := base.GetBaseData(d)
	if err != nil {
		return nil, err
	}

	if s.Strategy.Debug {
		fmt.Println("trend on data")
	}

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

	orders := p.GetOpenOrdersForStrategy(s.GetLabel())
	trade := p.GetTradeForStrategy(s.GetLabel())

	// fmt.Println("trend.go has", len(orders), "orders", trade)

	if trade == nil && len(orders) == 0 {
		return s.checkEntry(es, p, d, fe)
	}

	if trade != nil {
		return s.checkExit(es, p, d, fe)
	}

	if es.GetDecision() == "" {
		es.SetDecision(signal.DoNothing)
		es.SetDirection(eventtypes.DoNothing)
		es.AppendReason("null")
	}

	return &es, nil
}

func (s *Strategy) SupportsSimultaneousProcessing() bool {
	return true
}

// OnSimultaneousSignals analyses multiple data points simultaneously, allowing flexibility
// in allowing a strategy to only place an order for X currency if Y currency's price is Z
func (s *Strategy) OnSimultaneousSignals(d []data.Handler, p base.StrategyPortfolioHandler, fe base.FactorEngineHandler) ([]signal.Event, error) {
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

func (s *Strategy) checkEntry(es signal.Signal, p base.StrategyPortfolioHandler, d data.Handler, fe base.FactorEngineHandler) (signal.Event, error) {
	n60Chg := fe.Kline().N60PctChange.Last(0)
	price := d.Latest().ClosePrice()

	if s.Strategy.GetDirection() == order.Buy { // check for buy strategy
		if n60Chg.GreaterThan(decimal.NewFromInt(0)) {
			es.AppendReason("Strategy: n60Chg greater than zero")
			es.SetDecision(signal.Enter)
			es.SetStopLossPrice(price.Mul(decimal.NewFromFloat(0.9)))
		} else {
			es.AppendReason("Strategy: n60Chg less than zero")
			es.SetDecision(signal.DoNothing)
		}

	} else if s.Strategy.GetDirection() == order.Sell { // check sell strategy
		if n60Chg.LessThan(decimal.NewFromInt(0)) {
			es.AppendReason("Strategy: n60Chg less than zero")
			es.SetDecision(signal.Enter)
			es.SetStopLossPrice(price.Mul(decimal.NewFromFloat(1.1)))
		} else {
			es.AppendReason("Strategy: n60Chg greater than zero")
			es.SetDecision(signal.DoNothing)
		}
	}
	return &es, nil
}

func (s *Strategy) checkExit(es signal.Signal, p base.StrategyPortfolioHandler, d data.Handler, fe base.FactorEngineHandler) (signal.Event, error) {

	// if trade.ProfitLossPoints.GreaterThan(decimal.NewFromFloat(10)) {
	// 	fmt.Println("trade profit greater than 10, exiting")

	currentTime := d.Latest().GetTime()
	trade := p.GetTradeForStrategy(s.GetLabel())
	minutesInTrade := int(currentTime.Sub(trade.EntryTime).Minutes())

	if minutesInTrade < -2 {
		fmt.Println("ERROR negative seconds in trade", currentTime, trade.EntryTime)
		reason := fmt.Sprintf("negative %d minutes in trade", minutesInTrade)
		es.AppendReason(reason)
		os.Exit(2)

	} else if minutesInTrade > 60 {
		// handle exit
		n60PctChg := fe.Kline().N60PctChange.Last(1)

		// fmt.Println("check exit", es.GetTime(), minutesInTrade, n60PctChg)

		// CHECK EXIT BUY
		if s.Strategy.GetDirection() == order.Buy {
			if n60PctChg.LessThan(decimal.NewFromFloat(0)) {
				es.SetDecision(signal.Exit)
				es.AppendReason(fmt.Sprintf("Strategy: t >. %d min and N60PctChange is negative.", minutesInTrade))
			} else {
				es.SetDecision(signal.DoNothing)
				es.AppendReason(fmt.Sprintf("Strategy: Stay in long. N60PctChange is positive. (%d).", minutesInTrade))
			}
		}

		// CHECK EXIT SELL
		if s.Strategy.GetDirection() == order.Sell {
			if n60PctChg.GreaterThan(decimal.NewFromFloat(1)) {
				es.SetDecision(signal.Exit)
				es.AppendReason(fmt.Sprintf("Strategy.go says: exiting t > (%d) min and N60PctChange is positive.", minutesInTrade))
			} else {
				es.SetDecision(signal.DoNothing)
				es.AppendReason(fmt.Sprintf("Strategy.go says: Stay in short. N60PctChange is negative. (%d).", minutesInTrade))
			}
		}

	} else {
		es.SetDecision(signal.DoNothing)
		es.AppendReason(fmt.Sprintf("trade started %d minutes ago.", minutesInTrade))
	}

	return &es, nil
}
