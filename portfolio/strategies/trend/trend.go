package trend

import (
	"fmt"
	"os"

	"gocryptotrader/common"
	"gocryptotrader/data"
	"gocryptotrader/eventtypes"
	"gocryptotrader/eventtypes/signal"
	"gocryptotrader/exchange/order"
	"gocryptotrader/log"
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

func (s *Strategy) OnData(d data.Handler, p base.StrategyPortfolioHandler, fe base.FactorEngineHandler) (signal.Event, error) {
	if p.GetLiveMode() {
		s.Strategy.Debug = true
		log.Infoln(log.Global, "trend ONDATA", d.Latest().GetTime(), s.Strategy.GetDirection(), d.Latest().Pair(), len(d.History()), len(fe.Minute().Close))
	}
	if d == nil {
		return nil, eventtypes.ErrNilEvent
	}
	es, err := base.GetBaseData(d)
	if err != nil {
		return nil, err
	}

	// bar := fe.Minute().GetCurrentTime()
	// fmt.Println("straegy on bar", s.GetID(), s.GetDirection(), d.Latest().GetTime(), bar)

	// set defaults
	es.SetPrice(d.Latest().ClosePrice())
	es.SetAmount(decimal.NewFromFloat(1.0))

	// offset := d.Offset()
	// fmt.Println("trend.go offset", offset)

	// if offset <= int(60) {
	// 	es.SetDecision(signal.DoNothing)
	// 	es.AppendReason("Not enough data for signal generation")
	// 	return &es, nil
	// }

	if !d.HasDataAtTime(d.Latest().GetTime()) {
		es.SetDirection(eventtypes.MissingData)
		es.SetDecision(signal.DoNothing)
		es.AppendReason(fmt.Sprintf("missing data at %v, cannot perform any actions. trend", d.Latest().GetTime()))
		return &es, nil
	}

	currentTime := d.Latest().GetTime()
	orders := p.GetOpenOrdersForStrategy(s.GetID())
	trade := p.GetTradeForStrategy(s.GetID())

	// if trade != nil {
	// 	fmt.Println("trade amount", trade.Amount, "profit is", trade.ProfitLossPoints.Mul(trade.Amount))
	// }

	if trade == nil && len(orders) == 0 {
		m60Chg := fe.Minute().M60PctChange.Last(0)

		if s.Strategy.GetDirection() == order.Buy { // check for buy strategy
			if m60Chg.GreaterThan(decimal.NewFromInt(0)) {
				es.AppendReason(fmt.Sprintf("Strategy: m60Chg greater than zero (%v)", m60Chg))
				es.SetDecision(signal.Enter)
			} else {
				es.AppendReason(fmt.Sprintf("Strategy: m60Chg less than zero (%v)", m60Chg))
				es.SetDecision(signal.DoNothing)
			}

		} else if s.Strategy.GetDirection() == order.Sell { // check sell strategy
			if m60Chg.LessThan(decimal.NewFromInt(0)) {
				es.AppendReason(fmt.Sprintf("Strategy: m60Chg less than zero (%v)", m60Chg))
				es.SetDecision(signal.Enter)
			} else {
				es.AppendReason(fmt.Sprintf("Strategy: m60Chg greater than zero (%v)", m60Chg))
				es.SetDecision(signal.DoNothing)
			}
		}
	} else {
		minutesInTrade := int(currentTime.Sub(trade.EntryTime).Minutes())
		if minutesInTrade < -2 {

			fmt.Println("ERROR negative seconds in trade", currentTime, trade.EntryTime)
			reason := fmt.Sprintf("negative %d minutes in trade", minutesInTrade)
			es.AppendReason(reason)
			os.Exit(2)

		} else if minutesInTrade > 60 {
			// handle exit
			m60PctChg := fe.Minute().M60PctChange.Last(1)

			// CHECK EXIT BUY
			if s.Strategy.GetDirection() == order.Buy {
				if m60PctChg.LessThan(decimal.NewFromFloat(-1)) {
					es.SetDecision(signal.Exit)
					es.AppendReason(fmt.Sprintf("Strategy: t >. %d min and M60PctChange is negative.", minutesInTrade))
				} else {
					es.SetDecision(signal.DoNothing)
					es.AppendReason(fmt.Sprintf("Strategy: Stay in long. M60PctChange is positive. (%d).", minutesInTrade))
				}
			}

			// CHECK EXIT SELL
			if s.Strategy.GetDirection() == order.Sell {
				if m60PctChg.GreaterThan(decimal.NewFromFloat(1)) {
					es.SetDecision(signal.Exit)
					es.AppendReason(fmt.Sprintf("Strategy.go says: exiting t > (%d) min and M60PctChange is positive.", minutesInTrade))
				} else {
					es.SetDecision(signal.DoNothing)
					es.AppendReason(fmt.Sprintf("Strategy.go says: Stay in short. M60PctChange is negative. (%d).", minutesInTrade))
				}
			}

		} else {
			es.SetDecision(signal.DoNothing)
			es.AppendReason(fmt.Sprintf("trade started %d minutes ago.", minutesInTrade))
		}

	}

	if es.GetDecision() == "" {
		es.SetDecision(signal.DoNothing)
		es.SetDirection(eventtypes.DoNothing)
		es.AppendReason("null")
	}

	return &es, nil

	// fmt.Println("ALREADY IN TRADE")
	// if p.GetLiveMode() {
	// 	log.Debugln(log.TradeMgr, s.GetID(), "can trade")
	// }
	// fmt.Println(s.GetID(), "")

	// if trade.ProfitLossPoints.GreaterThan(decimal.NewFromFloat(10)) {
	// 	fmt.Println("trade profit greater than 10, exiting")
	// 	es.SetDecision(signal.Exit)
	// } else {
	// 	es.SetDecision(signal.DoNothing)
	// }
	// pos := p.GetPositionForStrategy(s.GetID())
	// if len(orders) > 0 {
	// 	fmt.Printf("%s has %d orders type %s status %s\n", s.GetID(), len(orders), orders[0].OrderType, orders[0].Status)
	// }

	// if trade != nil {
	// 	fmt.Println("trend.go trade", trade, s.GetID())
	// }

	// else {
	// 	es.SetDecision(signal.Exit)
	// }

	// else {
	// 	es.SetDecision(signal.Exit)
	// 	es.SetDirection(order.Sell)
	// }
	// 	// fmt.Println("check for entry")
	// 	// whats the current date
	// 	// fmt.Println("bar time: ", d.Latest().GetTime())
	// 	// get the current bar from the factor engine
	// 	// bar := fe.Minute().LatestClose()

	// if p.GetVerbose() {
	// 	m := fe.Minute()
	// 	daily := fe.Daily()
	// 	factors := fmt.Sprintf(
	// 		"%s,%d,%v,%v,%v,%v",
	// 		s.ID,
	// 		m.GetCurrentTime().Unix(),
	// 		m.GetCurrentDateOpen(),
	// 		m.GetCurrentDateHigh(),
	// 		m.GetCurrentDateLow(),
	// 		len(daily.Open))
	//
	// 	log.Debugf(log.TradeMgr, "%s trend factors %s", es.CurrencyPair, factors)
	// }

	// 	// fmt.Println("bar", m.LatestClose(), m.LastUpdate, d.Latest().GetTime())
	// 	// what was the open of the day
	// } else {
	// 	fmt.Println("check for exit")
	// }

	// no trade

	// fmt.Println(s.GetPosition())
	// fmt.Printf("%s@%v@%s now:%v pl:%v\n", t.Direction, t.EntryPrice, t.Timestamp, t.CurrentPrice, t.NetProfit)
}

// SupportsSimultaneousProcessing highlights whether the strategy can handle multiple currency calculation
// There is nothing actually stopping this strategy from considering multiple currencies at once
// but for demonstration purposes, this strategy does not
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
