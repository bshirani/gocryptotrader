// package dollarcostaverage
//
// import (
// 	"fmt"
//
// 	"gocryptotrader/common"
// 	"gocryptotrader/data"
// 	"gocryptotrader/factors"
// 	"gocryptotrader/strategies/base"
// 	"gocryptotrader/eventtypes/signal"
// 	gctcommon "gocryptotrader/common"
// 	"gocryptotrader/exchanges/order"
// )
//
// const (
// 	// Name is the strategy name
// 	Name        = "dollarcostaverage"
// 	description = `Dollar-cost averaging (DCA) is an investment strategy in which an investor divides up the total amount to be invested across periodic purchases of a target asset in an effort to reduce the impact of volatility on the overall purchase. The purchases occur regardless of the asset's price and at regular intervals. In effect, this strategy removes much of the detailed work of attempting to time the market in order to make purchases of equities at the best prices.`
// )
//
// // Strategy is an implementation of the Handler interface
// type Strategy struct {
// 	base.Strategy
// }
//
// // Name returns the name
// func (s *Strategy) Name() string {
// 	return Name
// }
//
// // Description provides a nice overview of the strategy
// // be it definition of terms or to highlight its purpose
// func (s *Strategy) Description() string {
// 	return description
// }
//
// // OnData handles a data event and returns what action the strategy believes should occur
// // For dollarcostaverage, this means returning a buy signal on every event
// func (s *Strategy) OnData(d data.Handler, _ base.StrategyPortfolioHandler, fe *factors.Engine) (signal.Event, error) {
// 	if d == nil {
// 		return nil, eventtypes.ErrNilEvent
// 	}
// 	es, err := base.GetBaseData(d)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if !d.HasDataAtTime(d.Latest().GetTime()) {
// 		es.SetDirection(common.MissingData)
// 		es.AppendReason(fmt.Sprintf("missing data at %v, cannot perform any actions", d.Latest().GetTime()))
// 		return &es, nil
// 	}
//
// 	es.SetPrice(d.Latest().ClosePrice())
// 	es.SetDirection(order.Buy)
// 	es.AppendReason("DCA purchases on every iteration")
// 	return &es, nil
// }
//
// // SupportsSimultaneousProcessing highlights whether the strategy can handle multiple currency calculation
// func (s *Strategy) SupportsSimultaneousProcessing() bool {
// 	return true
// }
//
// // OnSimultaneousSignals analyses multiple data points simultaneously, allowing flexibility
// // in allowing a strategy to only place an order for X currency if Y currency's price is Z
// // For dollarcostaverage, the strategy is always "buy", so it uses the OnData function
// func (s *Strategy) OnSimultaneousSignals(d []data.Handler, p base.StrategyPortfolioHandler, fe *factors.Engine) ([]signal.Event, error) {
// 	var resp []signal.Event
// 	var errs gctcommon.Errors
// 	for i := range d {
// 		sigEvent, err := s.OnData(d[i], p, fe)
// 		if err != nil {
// 			errs = append(errs, err)
// 		} else {
// 			resp = append(resp, sigEvent)
// 		}
// 	}
//
// 	if len(errs) > 0 {
// 		return nil, errs
// 	}
// 	return resp, nil
// }
//
// // SetCustomSettings not required for DCA
// func (s *Strategy) SetCustomSettings(_ map[string]interface{}) error {
// 	return base.ErrCustomSettingsUnsupported
// }
//
// // SetDefaults not required for DCA
// func (s *Strategy) SetDefaults() {}
