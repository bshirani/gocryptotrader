package statistics

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	gctcommon "gocryptotrader/common"
	"gocryptotrader/currency"
	"gocryptotrader/eventtypes"
	"gocryptotrader/eventtypes/fill"
	"gocryptotrader/eventtypes/order"
	"gocryptotrader/eventtypes/signal"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/log"
	"gocryptotrader/portfolio/compliance"
	"gocryptotrader/portfolio/holdings"
	"gocryptotrader/portfolio/statistics/strategystatistics"
)

// Reset returns the struct to defaults
func (s *Statistic) Reset() {
	*s = Statistic{}
}

// SetupEventForTime sets up the big map for to store important data at each time interval
func (s *Statistic) SetupEventForTime(ev eventtypes.DataEventHandler) error {
	if ev == nil {
		return eventtypes.ErrNilEvent
	}
	ex := ev.GetExchange()
	a := ev.GetAssetType()
	p := ev.Pair()
	s.setupMap(ex, a)
	lookup := s.StrategyStatistics[ex][a][p]
	if lookup == nil {
		// fmt.Println("lookup nil")
		lookup = &strategystatistics.StrategyStatistic{}
	}
	for i := range lookup.Events {
		if lookup.Events[i].DataEvent.GetTime().Equal(ev.GetTime()) &&
			lookup.Events[i].DataEvent.GetExchange() == ev.GetExchange() &&
			lookup.Events[i].DataEvent.GetAssetType() == ev.GetAssetType() &&
			lookup.Events[i].DataEvent.Pair().Equal(ev.Pair()) &&
			lookup.Events[i].DataEvent.GetOffset() == ev.GetOffset() {
			return ErrAlreadyProcessed
		}
	}
	lookup.Events = append(lookup.Events,
		strategystatistics.EventStore{
			DataEvent: ev,
		},
	)
	// fmt.Println("adding event to stats", ev.GetTime())
	s.StrategyStatistics[ex][a][p] = lookup

	return nil
}

func (s *Statistic) setupMap(ex string, a asset.Item) {
	if s.StrategyStatistics == nil {
		s.StrategyStatistics = make(map[string]map[asset.Item]map[currency.Pair]*strategystatistics.StrategyStatistic)
	}
	if s.StrategyStatistics[ex] == nil {
		s.StrategyStatistics[ex] = make(map[asset.Item]map[currency.Pair]*strategystatistics.StrategyStatistic)
	}
	if s.StrategyStatistics[ex][a] == nil {
		s.StrategyStatistics[ex][a] = make(map[currency.Pair]*strategystatistics.StrategyStatistic)
	}
}

// SetEventForOffset sets the event for the time period in the event
func (s *Statistic) SetEventForOffset(ev eventtypes.EventHandler) error {
	if ev == nil {
		return eventtypes.ErrNilEvent
	}
	if s.StrategyStatistics == nil {
		return errExchangeAssetPairStatsUnset
	}
	exch := ev.GetExchange()
	a := ev.GetAssetType()
	p := ev.Pair()
	offset := ev.GetOffset()
	lookup := s.StrategyStatistics[exch][a][p]
	if lookup == nil {
		return fmt.Errorf("%w for %v %v %v to set signal event", errStrategyStatisticsUnset, exch, a, p)
	}

	for i := len(lookup.Events) - 1; i >= 0; i-- {
		if lookup.Events[i].DataEvent.GetOffset() == offset {
			return applyEventAtOffset(ev, lookup, i)
		}
	}

	return nil
}

func applyEventAtOffset(ev eventtypes.EventHandler, lookup *strategystatistics.StrategyStatistic, i int) error {
	switch t := ev.(type) {
	case eventtypes.DataEventHandler:
		lookup.Events[i].DataEvent = t
	case signal.Event:
		lookup.Events[i].SignalEvent = t
	case order.Event:
		lookup.Events[i].OrderEvent = t
	case fill.Event:
		lookup.Events[i].FillEvent = t
	default:
		return fmt.Errorf("unknown event type received: %v", ev)
	}
	return nil
}

// AddHoldingsForTime adds all holdings to the statistics at the time period
func (s *Statistic) AddHoldingsForTime(h *holdings.Holding) error {
	// fmt.Println("add holding for time....................")
	if s.StrategyStatistics == nil {
		return errExchangeAssetPairStatsUnset
	}
	lookup := s.StrategyStatistics[h.Exchange][h.Asset][h.Pair]
	if lookup == nil {
		return fmt.Errorf("%w for %v %v %v to set holding event", errStrategyStatisticsUnset, h.Exchange, h.Asset, h.Pair)
	}
	for i := len(lookup.Events) - 1; i >= 0; i-- {
		if lookup.Events[i].DataEvent.GetOffset() == h.Offset {
			lookup.Events[i].Holdings = *h
			break
		}
	}
	return nil
}

// AddComplianceSnapshotForTime adds the compliance snapshot to the statistics at the time period
func (s *Statistic) AddComplianceSnapshotForTime(c compliance.Snapshot, e fill.Event) error {
	if e == nil {
		return eventtypes.ErrNilEvent
	}
	if s.StrategyStatistics == nil {
		return errExchangeAssetPairStatsUnset
	}
	exch := e.GetExchange()
	a := e.GetAssetType()
	p := e.Pair()
	lookup := s.StrategyStatistics[exch][a][p]
	if lookup == nil {
		return fmt.Errorf("%w for %v %v %v to set compliance snapshot", errStrategyStatisticsUnset, exch, a, p)
	}
	for i := len(lookup.Events) - 1; i >= 0; i-- {
		if lookup.Events[i].DataEvent.GetOffset() == e.GetOffset() {
			lookup.Events[i].Transactions = c
			break
		}
	}

	return nil
}

// CalculateAllResults calculates the statistics of all exchange asset pair holdings,
// orders, ratios and drawdowns
func (s *Statistic) CalculateAllResults() error {
	log.Info(log.StrategyMgr, "calculating backtesting results")
	s.PrintAllEventsChronologically()
	currCount := 0
	var finalResults []FinalResultsHolder
	var err error
	// var startDate, endDate time.Time
	for exchangeName, exchangeMap := range s.StrategyStatistics {
		for assetItem, assetMap := range exchangeMap {
			for pair, stats := range assetMap {
				currCount++
				// var f funding.IPairReader
				last := stats.Events[len(stats.Events)-1]
				// startDate = stats.Events[0].DataEvent.GetTime()
				// endDate = last.DataEvent.GetTime()
				// var event eventtypes.EventHandler
				// switch {
				// case last.FillEvent != nil:
				// 	event = last.FillEvent
				// case last.SignalEvent != nil:
				// 	event = last.SignalEvent
				// default:
				// 	event = last.DataEvent
				// }
				// f, err = funds.GetFundingForEvent(event)
				// if err != nil {
				// 	return err
				// }
				err = stats.CalculateResults()
				if err != nil {
					log.Error(log.StrategyMgr, err)
				}
				stats.PrintResults(exchangeName, assetItem, pair, false)
				stats.FinalHoldings = last.Holdings
				stats.InitialHoldings = stats.Events[0].Holdings
				stats.FinalOrders = last.Transactions
				s.AllStats = append(s.AllStats, *stats)

				finalResults = append(finalResults, FinalResultsHolder{
					Exchange:         exchangeName,
					Asset:            assetItem,
					Pair:             pair,
					MaxDrawdown:      stats.MaxDrawdown,
					MarketMovement:   stats.MarketMovement,
					StrategyMovement: stats.StrategyMovement,
				})
				s.TotalBuyOrders += stats.BuyOrders
				s.TotalSellOrders += stats.SellOrders
				if stats.ShowMissingDataWarning {
					s.WasAnyDataMissing = true
				}
			}
		}
	}
	// s.Funding = funds.GenerateReport(startDate, endDate)
	s.TotalOrders = s.TotalBuyOrders + s.TotalSellOrders
	if currCount > 1 {
		s.BiggestDrawdown = s.GetTheBiggestDrawdownAcrossCurrencies(finalResults)
		s.BestMarketMovement = s.GetBestMarketPerformer(finalResults)
		s.BestStrategyResults = s.GetBestStrategyPerformer(finalResults)
		s.PrintTotalResults(false)
	}

	return nil
}

// PrintTotalResults outputs all results to the CMD
func (s *Statistic) PrintTotalResults(isUsingExchangeLevelFunding bool) {
	log.Info(log.StrategyMgr, "------------------Strategy-----------------------------------")
	log.Infof(log.StrategyMgr, "Strategy Name: %v", s.StrategyName)
	log.Infof(log.StrategyMgr, "Strategy Nickname: %v", s.StrategyNickname)
	log.Infof(log.StrategyMgr, "Strategy Goal: %v\n\n", s.StrategyGoal)
	// log.Info(log.StrategyMgr, "------------------Funding------------------------------------")
	// for i := range s.Funding.Items {
	// 	log.Infof(log.StrategyMgr, "Exchange: %v", s.Funding.Items[i].Exchange)
	// 	log.Infof(log.StrategyMgr, "Asset: %v", s.Funding.Items[i].Asset)
	// 	log.Infof(log.StrategyMgr, "Currency: %v", s.Funding.Items[i].Currency)
	// 	if !s.Funding.Items[i].PairedWith.IsEmpty() {
	// 		log.Infof(log.StrategyMgr, "Paired with: %v", s.Funding.Items[i].PairedWith)
	// 	}
	// 	log.Infof(log.StrategyMgr, "Initial funds: %v", s.Funding.Items[i].InitialFunds)
	// 	log.Infof(log.StrategyMgr, "Initial funds in USD: $%v", s.Funding.Items[i].InitialFundsUSD)
	// 	log.Infof(log.StrategyMgr, "Final funds: %v", s.Funding.Items[i].FinalFunds)
	// 	log.Infof(log.StrategyMgr, "Final funds in USD: $%v", s.Funding.Items[i].FinalFundsUSD)
	// 	if s.Funding.Items[i].InitialFunds.IsZero() {
	// 		log.Info(log.StrategyMgr, "Difference: ∞%")
	// 	} else {
	// 		log.Infof(log.StrategyMgr, "Difference: %v%%", s.Funding.Items[i].Difference)
	// 	}
	// 	if s.Funding.Items[i].TransferFee.GreaterThan(decimal.Zero) {
	// 		log.Infof(log.StrategyMgr, "Transfer fee: %v", s.Funding.Items[i].TransferFee)
	// 	}
	// 	log.Info(log.StrategyMgr, "")
	// }
	// log.Infof(log.StrategyMgr, "Initial total funds in USD: $%v", s.Funding.InitialTotalUSD)
	// log.Infof(log.StrategyMgr, "Final total funds in USD: $%v", s.Funding.FinalTotalUSD)
	// log.Infof(log.StrategyMgr, "Difference: %v%%\n", s.Funding.Difference)

	log.Info(log.StrategyMgr, "------------------Total Results------------------------------")
	log.Info(log.StrategyMgr, "------------------Orders-------------------------------------")
	log.Infof(log.StrategyMgr, "Total buy orders: %v", s.TotalBuyOrders)
	log.Infof(log.StrategyMgr, "Total sell orders: %v", s.TotalSellOrders)
	log.Infof(log.StrategyMgr, "Total orders: %v\n\n", s.TotalOrders)

	if s.BiggestDrawdown != nil {
		log.Info(log.StrategyMgr, "------------------Biggest Drawdown-----------------------")
		log.Infof(log.StrategyMgr, "Exchange: %v Asset: %v Currency: %v", s.BiggestDrawdown.Exchange, s.BiggestDrawdown.Asset, s.BiggestDrawdown.Pair)
		log.Infof(log.StrategyMgr, "Highest Price: %v", s.BiggestDrawdown.MaxDrawdown.Highest.Price.Round(8))
		log.Infof(log.StrategyMgr, "Highest Price Time: %v", s.BiggestDrawdown.MaxDrawdown.Highest.Time)
		log.Infof(log.StrategyMgr, "Lowest Price: %v", s.BiggestDrawdown.MaxDrawdown.Lowest.Price.Round(8))
		log.Infof(log.StrategyMgr, "Lowest Price Time: %v", s.BiggestDrawdown.MaxDrawdown.Lowest.Time)
		log.Infof(log.StrategyMgr, "Calculated Drawdown: %v%%", s.BiggestDrawdown.MaxDrawdown.DrawdownPercent.Round(2))
		log.Infof(log.StrategyMgr, "Difference: %v", s.BiggestDrawdown.MaxDrawdown.Highest.Price.Sub(s.BiggestDrawdown.MaxDrawdown.Lowest.Price).Round(8))
		log.Infof(log.StrategyMgr, "Drawdown length: %v\n\n", s.BiggestDrawdown.MaxDrawdown.IntervalDuration)
	}
	if s.BestMarketMovement != nil && s.BestStrategyResults != nil {
		log.Info(log.StrategyMgr, "------------------Orders----------------------------------")
		log.Infof(log.StrategyMgr, "Best performing market movement: %v %v %v %v%%", s.BestMarketMovement.Exchange, s.BestMarketMovement.Asset, s.BestMarketMovement.Pair, s.BestMarketMovement.MarketMovement.Round(2))
		log.Infof(log.StrategyMgr, "Best performing strategy movement: %v %v %v %v%%\n\n", s.BestStrategyResults.Exchange, s.BestStrategyResults.Asset, s.BestStrategyResults.Pair, s.BestStrategyResults.StrategyMovement.Round(2))
	}
}

// GetBestMarketPerformer returns the best final market movement
func (s *Statistic) GetBestMarketPerformer(results []FinalResultsHolder) *FinalResultsHolder {
	result := &FinalResultsHolder{}
	for i := range results {
		if results[i].MarketMovement.GreaterThan(result.MarketMovement) || result.MarketMovement.IsZero() {
			result = &results[i]
			break
		}
	}

	return result
}

// GetBestStrategyPerformer returns the best performing strategy result
func (s *Statistic) GetBestStrategyPerformer(results []FinalResultsHolder) *FinalResultsHolder {
	result := &FinalResultsHolder{}
	for i := range results {
		if results[i].StrategyMovement.GreaterThan(result.StrategyMovement) || result.StrategyMovement.IsZero() {
			result = &results[i]
		}
	}

	return result
}

// GetTheBiggestDrawdownAcrossCurrencies returns the biggest drawdown across all currencies in a backtesting run
func (s *Statistic) GetTheBiggestDrawdownAcrossCurrencies(results []FinalResultsHolder) *FinalResultsHolder {
	result := &FinalResultsHolder{}
	for i := range results {
		if results[i].MaxDrawdown.DrawdownPercent.GreaterThan(result.MaxDrawdown.DrawdownPercent) || result.MaxDrawdown.DrawdownPercent.IsZero() {
			result = &results[i]
		}
	}

	return result
}

func addEventOutputToTime(events []eventOutputHolder, t time.Time, message string) []eventOutputHolder {
	for i := range events {
		if events[i].Time.Equal(t) {
			events[i].Events = append(events[i].Events, message)
			return events
		}
	}
	events = append(events, eventOutputHolder{Time: t, Events: []string{message}})
	return events
}

// PrintAllEventsChronologically outputs all event details in the CMD
// rather than separated by exchange, asset and currency pair, it's
// grouped by time to allow a clearer picture of events
func (s *Statistic) PrintAllEventsChronologically() {
	var results []eventOutputHolder
	log.Info(log.StrategyMgr, "------------------Events-------------------------------------")
	var errs gctcommon.Errors
	for exch, x := range s.StrategyStatistics {
		for a, y := range x {
			for pair, strategyStatistic := range y {
				for i := range strategyStatistic.Events {
					switch {
					case strategyStatistic.Events[i].FillEvent != nil:
						direction := strategyStatistic.Events[i].FillEvent.GetDirection()
						if direction == eventtypes.CouldNotBuy ||
							direction == eventtypes.CouldNotSell ||
							direction == eventtypes.DoNothing ||
							direction == eventtypes.MissingData ||
							direction == eventtypes.TransferredFunds ||
							direction == "" {
							results = addEventOutputToTime(results, strategyStatistic.Events[i].FillEvent.GetTime(),
								fmt.Sprintf("%v %v %v %v | Price: $%v - Direction: %v - Reason: %s",
									strategyStatistic.Events[i].FillEvent.GetTime().Format(gctcommon.SimpleTimeFormat),
									strategyStatistic.Events[i].FillEvent.GetExchange(),
									strategyStatistic.Events[i].FillEvent.GetAssetType(),
									strategyStatistic.Events[i].FillEvent.Pair(),
									strategyStatistic.Events[i].FillEvent.GetClosePrice().Round(8),
									strategyStatistic.Events[i].FillEvent.GetDirection(),
									strategyStatistic.Events[i].FillEvent.GetReason()))
						} else {
							results = addEventOutputToTime(results, strategyStatistic.Events[i].FillEvent.GetTime(),
								fmt.Sprintf("%v %v %v %v | Price: $%v - Amount: %v - Fee: $%v - Total: $%v - Direction %v - Reason: %s",
									strategyStatistic.Events[i].FillEvent.GetTime().Format(gctcommon.SimpleTimeFormat),
									strategyStatistic.Events[i].FillEvent.GetExchange(),
									strategyStatistic.Events[i].FillEvent.GetAssetType(),
									strategyStatistic.Events[i].FillEvent.Pair(),
									strategyStatistic.Events[i].FillEvent.GetPurchasePrice().Round(8),
									strategyStatistic.Events[i].FillEvent.GetAmount().Round(8),
									strategyStatistic.Events[i].FillEvent.GetExchangeFee().Round(8),
									strategyStatistic.Events[i].FillEvent.GetTotal().Round(8),
									strategyStatistic.Events[i].FillEvent.GetDirection(),
									strategyStatistic.Events[i].FillEvent.GetReason(),
								))
						}
					case strategyStatistic.Events[i].SignalEvent != nil:
						results = addEventOutputToTime(results, strategyStatistic.Events[i].SignalEvent.GetTime(),
							fmt.Sprintf("%v %v %v %v | Price: $%v - Reason: %v",
								strategyStatistic.Events[i].SignalEvent.GetTime().Format(gctcommon.SimpleTimeFormat),
								strategyStatistic.Events[i].SignalEvent.GetExchange(),
								strategyStatistic.Events[i].SignalEvent.GetAssetType(),
								strategyStatistic.Events[i].SignalEvent.Pair(),
								strategyStatistic.Events[i].SignalEvent.GetPrice().Round(8),
								strategyStatistic.Events[i].SignalEvent.GetReason()))
					case strategyStatistic.Events[i].DataEvent != nil:
						results = addEventOutputToTime(results, strategyStatistic.Events[i].DataEvent.GetTime(),
							fmt.Sprintf("%v %v %v %v | Price: $%v - Reason: %v",
								strategyStatistic.Events[i].DataEvent.GetTime().Format(gctcommon.SimpleTimeFormat),
								strategyStatistic.Events[i].DataEvent.GetExchange(),
								strategyStatistic.Events[i].DataEvent.GetAssetType(),
								strategyStatistic.Events[i].DataEvent.Pair(),
								strategyStatistic.Events[i].DataEvent.ClosePrice().Round(8),
								strategyStatistic.Events[i].DataEvent.GetReason()))
					default:
						errs = append(errs, fmt.Errorf("%v %v %v unexpected data received %+v", exch, a, pair, strategyStatistic.Events[i]))
					}
				}
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		b1 := results[i]
		b2 := results[j]
		return b1.Time.Before(b2.Time)
	})
	for i := range results {
		for j := range results[i].Events {
			log.Info(log.StrategyMgr, results[i].Events[j])
		}
	}
	if len(errs) > 0 {
		log.Info(log.StrategyMgr, "------------------Errors-------------------------------------")
		for i := range errs {
			log.Info(log.StrategyMgr, errs[i].Error())
		}
	}
}

// SetStrategyName sets the name for statistical identification
func (s *Statistic) SetStrategyName(name string) {
	s.StrategyName = name
}

// Serialise outputs the Statistic struct in json
func (s *Statistic) Serialise() (string, error) {
	resp, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return "", err
	}

	return string(resp), nil
}
