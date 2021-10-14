package engine

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/database/repository/liveorder"
	"gocryptotrader/database/repository/livetrade"
	"gocryptotrader/eventtypes"
	"gocryptotrader/eventtypes/cancel"
	"gocryptotrader/eventtypes/event"
	"gocryptotrader/eventtypes/fill"
	"gocryptotrader/eventtypes/order"
	"gocryptotrader/eventtypes/signal"
	"gocryptotrader/eventtypes/submit"
	"gocryptotrader/exchange/asset"
	gctorder "gocryptotrader/exchange/order"
	"gocryptotrader/log"
	"gocryptotrader/portfolio/compliance"
	"gocryptotrader/portfolio/holdings"
	"gocryptotrader/portfolio/positions"
	"gocryptotrader/portfolio/risk"
	"gocryptotrader/portfolio/strategies"
	"gocryptotrader/portfolio/trades"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

// Setup creates a portfolio manager instance and sets private fields
func SetupPortfolio(st []strategies.Handler, bot *Engine, sh SizeHandler, r risk.Handler, riskFreeRate decimal.Decimal) (*Portfolio, error) {
	// log.Infof(log.TradeManager, "Setting up Portfolio")
	if sh == nil {
		return nil, errSizeManagerUnset
	}
	if riskFreeRate.IsNegative() {
		return nil, errNegativeRiskFreeRate
	}
	if r == nil {
		return nil, errRiskManagerUnset
	}
	p := &Portfolio{}

	// create position for every strategy
	// create open trades array for every strategy
	// you need the strategy IDS here
	p.store.positions = make(map[string]*positions.Position)
	p.store.openTrade = make(map[string]*livetrade.Details)
	p.store.openOrders = make(map[string][]*liveorder.Details)
	p.store.closedOrders = make(map[string][]*liveorder.Details)
	p.store.closedTrades = make(map[string][]*livetrade.Details)

	// load open trade from the database
	// only in live mode do we  do this, so portfolio must be live aware unless it has a callback
	// handle this in the trademanger
	// log.Infof(log.TradeManager, "there are %d trades running", livetrade.Count())

	// load all pending and open trades from the database into positions and trades

	// what does this do?
	// bot.Backtest.Datas.Setup()

	p.orderManager = bot.OrderManager
	p.bot = bot
	p.sizeManager = sh
	p.riskManager = r
	p.riskFreeRate = riskFreeRate
	p.strategies = st

	// set initial opentrade/positions
	for _, s := range p.strategies {
		p.store.positions[s.GetID()] = &positions.Position{}
		p.store.closedTrades[s.GetID()] = make([]*livetrade.Details, 0)
		p.store.openOrders[s.GetID()] = make([]*liveorder.Details, 0)
		s.SetWeight(decimal.NewFromFloat(1.5))

	}

	// load existing positions from database
	// only in live mode do we do this
	// should handle in trademanager
	activeTrades, _ := livetrade.Active()
	for _, t := range activeTrades {
		p.store.openTrade[t.StrategyID] = &t
		pos := p.store.positions[t.StrategyID]
		pos.Active = true
	}

	activeOrders, _ := liveorder.Active()
	for _, t := range activeOrders {
		p.store.openOrders[t.StrategyID] = append(p.store.openOrders[t.StrategyID], &t)
	}

	if !p.bot.Settings.EnableDryRun {
		log.Infoln(log.TradeManager, "(live mode) Loaded Trades", len(activeTrades))
		log.Infoln(log.TradeManager, "(live mode) Loaded Orders", len(activeOrders))
	}

	return p, nil
}

// Reset returns the portfolio manager to its default state
func (p *Portfolio) Reset() {
	p.exchangeAssetPairSettings = nil
}

func (p *Portfolio) OnSubmit(submit submit.Event) {
	// fmt.Println(submit.GetStrategyID(), "portfolio.OnSubmit", submit.GetInternalOrderID())
	// p.store.closedOrders[submit.GetStrategyID()] = append(p.store.closedOrders[submit.GetStrategyID()], ord)
	// p.store.openOrders[submit.GetStrategyID()] = make([]*liveorder.Details, 0)
}

func (p *Portfolio) OnCancel(cancel cancel.Event) {
	// fmt.Println("portfolio received cancelled order", cancel)
}

// OnSignal receives the event from the strategy on whether it has signalled to buy, do nothing or sell
// on buy/sell, the portfolio manager will size the order and assess the risk of the order
// if successful, it will pass on an order.Order to be used by the exchange event handler to place an order based on
// the portfolio manager's recommendations
func (p *Portfolio) OnSignal(ev signal.Event, cs *ExchangeAssetPairSettings) (*order.Order, error) {
	if ev == nil || cs == nil {
		return nil, eventtypes.ErrNilArguments
	}
	if p.sizeManager == nil {
		return nil, errSizeManagerUnset
	}
	if p.riskManager == nil {
		return nil, errRiskManagerUnset
	}
	if ev.GetStrategyID() == "" {
		return nil, errStrategyIDUnset
	}

	trade := p.GetTradeForStrategy(ev.GetStrategyID())
	if trade != nil {
		if trade.Side == gctorder.Buy {
			trade.ProfitLossPoints = ev.GetPrice().Sub(trade.EntryPrice)
		} else if trade.Side == gctorder.Sell {
			trade.ProfitLossPoints = trade.EntryPrice.Sub(ev.GetPrice())
		} else {
			fmt.Println("trade is not sell or buy")
			os.Exit(2)
		}

		if p.bot.Config.LiveMode {
			p.printTradeDetails(trade)
		}
	}

	id, _ := uuid.NewV4()

	ev.SetDirection(gctorder.Buy) // FIXME

	o := &order.Order{
		Base: event.Base{
			Offset:       ev.GetOffset(),
			Exchange:     ev.GetExchange(),
			Time:         ev.GetTime(),
			CurrencyPair: ev.Pair(),
			AssetType:    ev.GetAssetType(),
			Interval:     ev.GetInterval(),
			Reason:       ev.GetReason(),
			StrategyID:   ev.GetStrategyID(),
		},
		ID:         id.String(),
		Direction:  ev.GetDirection(),
		Amount:     ev.GetAmount(),
		StrategyID: ev.GetStrategyID(),
	}

	lo := liveorder.Details{
		Status:     gctorder.New,
		OrderType:  gctorder.Market,
		Exchange:   ev.GetExchange(),
		InternalID: id.String(),
		StrategyID: ev.GetStrategyID(),
		Side:       ev.GetDirection(),
	}

	switch ev.GetDecision() {
	case signal.Enter:
		lo.Side = gctorder.Buy
		if !p.bot.Settings.EnableDryRun {
			id, _ := liveorder.Insert(lo)
			lo.ID = id
		}

		p.store.openOrders[ev.GetStrategyID()] = append(p.store.openOrders[ev.GetStrategyID()], &lo)

	case signal.Exit:
		ev.SetDirection(gctorder.Sell) // FIXME
		lo.Side = gctorder.Sell
		p.store.openOrders[ev.GetStrategyID()] = append(p.store.openOrders[ev.GetStrategyID()], &lo)

	case signal.DoNothing:
		return nil, nil

	default:
		return nil, errNoDecision
	}

	if ev.GetDirection() == "" {
		return o, errInvalidDirection
	}

	lookup := p.exchangeAssetPairSettings[ev.GetExchange()][ev.GetAssetType()][ev.Pair()]
	if lookup == nil {
		return nil, fmt.Errorf("%w for %v %v %v",
			errNoPortfolioSettings,
			ev.GetExchange(),
			ev.GetAssetType(),
			ev.Pair())
	}

	// sdir := p.strategies[0].Direction()

	// if pos.Active {
	// if (sdir == gctorder.Buy && ev.GetDirection() == gctorder.Buy) || (sdir == gctorder.Sell && ev.GetDirection() == gctorder.Sell) {
	// 	return nil, errAlreadyInTrade
	// }
	// }

	if ev.GetDirection() == eventtypes.DoNothing ||
		ev.GetDirection() == eventtypes.MissingData ||
		ev.GetDirection() == eventtypes.TransferredFunds ||
		ev.GetDirection() == "" {
		return nil, nil
	}

	o.Price = ev.GetPrice()
	o.Direction = ev.GetDirection()
	o.OrderType = gctorder.Market
	o.BuyLimit = ev.GetBuyLimit()
	o.SellLimit = ev.GetSellLimit()
	o.StrategyID = ev.GetStrategyID()
	sizedOrder := p.sizeOrder(ev, cs, o, decimal.NewFromFloat(100))
	sizedOrder.Amount = ev.GetAmount()

	p.recordTrade(ev)

	fmt.Println("PORTFOLIO CREATING NEW ORDER FOR", ev.GetDirection(), ev.GetStrategyID(), ev.GetReason())

	return p.evaluateOrder(ev, o, sizedOrder)
}

func (p *Portfolio) updatePosition(pos *positions.Position, amount decimal.Decimal) {
	pos.Amount = decimal.NewFromFloat(100.0)
}

func (p *Portfolio) GetOrderFromStore(orderid string) *gctorder.Detail {
	var foundOrd *gctorder.Detail
	ords, _ := p.bot.OrderManager.GetOrdersSnapshot("")
	for _, ord := range ords {
		if ord.ID != orderid {
			continue
		}
		foundOrd = &ord
	}
	if foundOrd.ID == "" {
		fmt.Println("ERROR order has no ID")
	}

	if foundOrd.Price == 0 {
		fmt.Println("ERROR order has no price ")
	}

	return foundOrd
}

func (p *Portfolio) createTrade(ev fill.Event, order *liveorder.Details) {
	// fmt.Println("found order", foundOrd.ID)
	foundOrd := p.GetOrderFromStore(ev.GetOrderID())
	stopLossPrice := decimal.NewFromFloat(foundOrd.Price).Sub(decimal.NewFromFloat(20))

	lt := livetrade.Details{
		Status:        gctorder.Open,
		StrategyID:    ev.GetStrategyID(),
		EntryTime:     foundOrd.Date,
		EntryOrderID:  foundOrd.ID,
		EntryPrice:    decimal.NewFromFloat(foundOrd.Price),
		StopLossPrice: stopLossPrice,
		Side:          foundOrd.Side,
		Pair:          foundOrd.Pair,
	}

	if lt.EntryPrice.IsZero() {
		fmt.Println("1111")
		os.Exit(2)
	}

	if !p.bot.Settings.EnableDryRun {
		id, err := livetrade.Insert(lt)
		lt.ID = id
		if err != nil {
			fmt.Println("error inserting trade", err)
			os.Exit(2)
		}
	}

	p.store.openTrade[ev.GetStrategyID()] = &lt
}

func (p *Portfolio) closeTrade(f fill.Event, t *livetrade.Details) {
	if t.Status == gctorder.Open {
		t.Status = gctorder.Closed
		p.store.closedTrades[f.GetStrategyID()] = append(p.store.closedTrades[f.GetStrategyID()], t)
		p.store.openTrade[f.GetStrategyID()] = nil

		if !p.bot.Settings.EnableDryRun {
			id, err := livetrade.Update(t)
			t.ID = id
			if err != nil || id == 0 {
				fmt.Println("error saving to db")
				os.Exit(2)
			}
		}

	} else {
		fmt.Println("TRYING TO CLOSE  ALREADY CLOSED TRADE. TRADE IS NOT OPEN")
		os.Exit(1)
	}
	// Velse if t.Status == livetrade.Pending {
	// 	ot := *p.store.openTrade[f.GetStrategyID()]
	// 	ot.Status = livetrade.Open
	// 	p.store.openTrade[f.GetStrategyID()] = &ot
	// }
}

// OnFill processes the event after an order has been placed by the exchange. Its purpose is to track holdings for future portfolio decisions.
func (p *Portfolio) OnFill(f fill.Event) {
	if f.GetStrategyID() == "" {
		fmt.Println("fill has no strategy ID")
		os.Exit(2)
	}

	// fmt.Println("PF ONFILL", f.GetStrategyID())

	// // create or update position
	// for _, pos := range p.store.positions {
	// 	if f.GetDirection() == gctorder.Sell {
	// 		pos.Amount = pos.Amount.Sub(f.GetAmount())
	// 	} else if f.GetDirection() == gctorder.Buy {
	// 		pos.Amount = pos.Amount.Add(f.GetAmount())
	// 	}
	//
	// 	if !pos.Amount.IsZero() {
	// 		pos.Active = true
	// 	} else {
	// 		pos.Active = false
	// 	}
	// }

	// update the orders
	// fmt.Println(submit.GetStrategyID(), "portfolio.OnFill", submit.GetInternalOrderID())

	// fmt.Println(f.GetStrategyID(), "closing order")
	order := p.store.openOrders[f.GetStrategyID()][0]
	p.store.closedOrders[f.GetStrategyID()] = append(p.store.closedOrders[f.GetStrategyID()], order)
	p.store.openOrders[f.GetStrategyID()] = make([]*liveorder.Details, 0)
	// fmt.Println(f.GetStrategyID(), " now has ", len(p.store.closedOrders[f.GetStrategyID()]))

	// update trades and orders here
	t := p.store.openTrade[f.GetStrategyID()]
	if t == nil {
		// fmt.Println("NEW TRADE")
		p.createTrade(f, order)
	} else if t.Status == gctorder.Open {
		// fmt.Println("NEW CLOSE TRADE")
		p.closeTrade(f, t)
	}

	// if f == nil {
	// 	return nil, eventtypes.ErrNilEvent
	// }
	// lookup := p.exchangeAssetPairSettings[f.GetExchange()][f.GetAssetType()][f.Pair()]
	// if lookup == nil {
	// 	return nil, fmt.Errorf("%w for %v %v %v", errNoPortfolioSettings, f.GetExchange(), f.GetAssetType(), f.Pair())
	// }
	// var err error

	// which strategy was filled?
	// what was the amount filled?
	// what was the direction of the fill?
	// what is the direction of the strategy?

	// entryPrice, _ := ev.GetPrice().Float64()
	// stopLossPrice, _ := ev.GetPrice().Float64()

	// st := p.GetStrategy(f.GetStrategyID())
	// fmt.Println("strategy", st)

	// t.Status = trades.Open
	// whats the strategy id of this fill?

	// // Get the holding from the previous iteration, create it if it doesn't yet have a timestamp
	// h := lookup.GetHoldingsForTime(f.GetTime().Add(-f.GetInterval().Duration()))
	// if !h.Timestamp.IsZero() {
	// 	h.Update(f)
	// } else {
	// 	h = lookup.GetLatestHoldings()
	// 	if h.Timestamp.IsZero() {
	// 		h, err = holdings.Create(f, decimal.NewFromFloat(1000.0), p.riskFreeRate)
	// 		// if err != nil {
	// 		// 	return nil, err
	// 		// }
	// 	} else {
	// 		h.Update(f)
	// 	}
	// }
	// err = p.setHoldingsForOffset(&h, true)
	// if errors.Is(err, errNoHoldings) {
	// 	err = p.setHoldingsForOffset(&h, false)
	// }
	// if err != nil {
	// 	log.Error(log.TradeManager, err)
	// }

	// err = p.addComplianceSnapshot(f)
	// if err != nil {
	// 	log.Error(log.TradeManager, err)
	// }

	// direction := f.GetDirection()
	// if direction == eventtypes.DoNothing ||
	// 	direction == eventtypes.CouldNotBuy ||
	// 	direction == eventtypes.CouldNotSell ||
	// 	direction == eventtypes.MissingData ||
	// 	direction == "" {
	// 	// fe, ok := f.(*fill.Fill)
	// 	// if !ok {
	// 	// 	return nil, fmt.Errorf("%w expected fill event", eventtypes.ErrInvalidDataType)
	// 	// }
	// 	// fe.ExchangeFee = decimal.Zero
	// 	// return fe, nil
	// }

	// fe, ok := f.(*fill.Fill)
	// if !ok {
	// 	return nil, fmt.Errorf("%w expected fill event", eventtypes.ErrInvalidDataType)
	// }
	// return fe, nil
}

// func (p *Portfolio) GetStrategy(id string) *base.Strategy {
// 	for _, s := range p.strategies {
// 		if s.GetID() == id {
// 			return s
// 		}
// 	}
// 	return nil
// }

// GetComplianceManager returns the order snapshots for a given exchange, asset, pair
func (p *Portfolio) GetComplianceManager(exchangeName string, a asset.Item, cp currency.Pair) (*compliance.Manager, error) {
	lookup := p.exchangeAssetPairSettings[exchangeName][a][cp]
	if lookup == nil {
		return nil, fmt.Errorf("%w for %v %v %v could not retrieve compliance manager", errNoPortfolioSettings, exchangeName, a, cp)
	}
	return &lookup.ComplianceManager, nil
}

// SetFee sets the fee rate
func (p *Portfolio) SetFee(exch string, a asset.Item, cp currency.Pair, fee decimal.Decimal) {
	lookup := p.exchangeAssetPairSettings[exch][a][cp]
	lookup.Fee = fee
}

func (p *Portfolio) GetVerbose() bool {
	return p.verbose
}

func (p *Portfolio) SetVerbose(verbose bool) {
	p.verbose = verbose
}

// GetFee can panic for bad requests, but why are you getting things that don't exist?
func (p *Portfolio) GetFee(exchangeName string, a asset.Item, cp currency.Pair) decimal.Decimal {
	if p.exchangeAssetPairSettings == nil {
		return decimal.Zero
	}
	lookup := p.exchangeAssetPairSettings[exchangeName][a][cp]
	if lookup == nil {
		return decimal.Zero
	}
	return lookup.Fee
}

// UpdateHoldings updates the portfolio holdings for the data event
func (p *Portfolio) UpdateHoldings(ev eventtypes.DataEventHandler) error {
	if ev == nil {
		return eventtypes.ErrNilEvent
	}
	lookup, ok := p.exchangeAssetPairSettings[ev.GetExchange()][ev.GetAssetType()][ev.Pair()]
	if !ok {
		return fmt.Errorf("%w for %v %v %v",
			errNoPortfolioSettings,
			ev.GetExchange(),
			ev.GetAssetType(),
			ev.Pair())
	}
	h := lookup.GetLatestHoldings()
	if h.Timestamp.IsZero() {
		var err error
		h, err = holdings.Create(ev, decimal.NewFromFloat(1000.0), p.riskFreeRate)
		if err != nil {
			return err
		}
	}
	h.UpdateValue(ev)
	err := p.setHoldingsForOffset(&h, true)
	if errors.Is(err, errNoHoldings) {
		err = p.setHoldingsForOffset(&h, false)
	}
	return err
}

// UpdateTrades updates the portfolio trades for the data event
func (p *Portfolio) UpdateTrades(ev eventtypes.DataEventHandler) {
	if ev == nil {
		return
	}
	// return nil
	_, ok := p.exchangeAssetPairSettings[ev.GetExchange()][ev.GetAssetType()][ev.Pair()]
	if !ok {
		return
		// return fmt.Errorf("%w for %v %v %v",
		// 	errNoPortfolioSettings,
		// 	ev.GetExchange(),
		// 	ev.GetAssetType(),
		// 	ev.Pair())
	}
	// t, _ := p.GetOpenTrade()

	// if err != nil {
	// 	// fmt.Println("error", t)
	// 	// t, _ = trades.Create(ev)
	// }

	// p.openTrade = t
	// t, _ = p.GetOpenTrade()
	// p.openTrade.UpdateValue(ev)

	// if t.Strategy != nil {
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// }
	// err := p.setTradesForOffset(&t, true)
	// if errors.Is(err, errNoTrades) {
	// 	err = p.setTradesForOffset(&t, false)
	// }
	// return err
}

// func (p *Portfolio) GetStrategies() []strategies.Handler {
// 	return p.strategies
// }

func (p *Portfolio) GetPositionForStrategy(sid string) *positions.Position {
	return p.store.positions[sid]
}

func (p *Portfolio) GetTradeForStrategy(sid string) *livetrade.Details {
	return p.store.openTrade[sid]
}

func (p *Portfolio) GetOpenOrdersForStrategy(sid string) []*liveorder.Details {
	return p.store.openOrders[sid]
}

func (p *Portfolio) GetAllClosedTrades() []*livetrade.Details {
	count := 0
	for _, s := range p.store.closedTrades {
		for _, _ = range s {
			count += 1
		}
	}
	res := make([]*livetrade.Details, count)
	count = 0
	for _, s := range p.store.closedTrades {
		for _, t := range s {
			if t != nil {
				res[count] = t
				count += 1
			}
		}
	}
	return res
}

// UpdatePositions updates the strategy's position for the data event
func (p *Portfolio) UpdatePositions(ev eventtypes.DataEventHandler) {
	if ev == nil {
		return
	}

	// for _, s := range p.GetStrategies() {
	// 	// update the strategies positions
	// 	fmt.Println("update position for ", s.Name)
	// }

	// portfolio has many strategies
	// we keep the position for each strategy
	for i, p := range p.store.positions {
		fmt.Println(i, p)
	}

	// pos := p.GetPositionForStrategy(p.strategies[0].ID())
	// fmt.Println("position:", pos.Amount)
	// pos.Amount = decimal.NewFromFloat(123.0)
	// pos.Active = false

	// return nil
	_, ok := p.exchangeAssetPairSettings[ev.GetExchange()][ev.GetAssetType()][ev.Pair()]
	if !ok {
		return
		// return fmt.Errorf("%w for %v %v %v",
		// 	errNoPortfolioSettings,
		// 	ev.GetExchange(),
		// 	ev.GetAssetType(),
		// 	ev.Pair())
	}
	// t, _ := p.GetOpenTrade()

	// if err != nil {
	// 	// fmt.Println("error", t)
	// 	// t, _ = trades.Create(ev)
	// }

	// p.openTrade = t
	// t, _ = p.GetOpenTrade()
	// p.openTrade.UpdateValue(ev)
	//
	// if t.Strategy != nil {
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// }
	// err := p.setTradesForOffset(&t, true)
	// if errors.Is(err, errNoTrades) {
	// 	err = p.setTradesForOffset(&t, false)
	// }
	// return err
}

// GetLatestHoldingsForAllCurrencies will return the current holdings for all loaded currencies
// this is useful to assess the position of your entire portfolio in order to help with risk decisions
func (p *Portfolio) GetLatestHoldingsForAllCurrencies() []holdings.Holding {
	var resp []holdings.Holding
	for _, x := range p.exchangeAssetPairSettings {
		for _, y := range x {
			for _, z := range y {
				holds := z.GetLatestHoldings()
				if !holds.Timestamp.IsZero() {
					resp = append(resp, holds)
				}
			}
		}
	}
	return resp
}

// GetLatestTradesForAllCurrencies will return the current holdings for all loaded currencies
// this is useful to assess the position of your entire portfolio in order to help with risk decisions
func (p *Portfolio) GetLatestTradesForAllCurrencies() []holdings.Holding {
	var resp []holdings.Holding
	for _, x := range p.exchangeAssetPairSettings {
		for _, y := range x {
			for _, z := range y {
				holds := z.GetLatestHoldings()
				if !holds.Timestamp.IsZero() {
					resp = append(resp, holds)
				}
			}
		}
	}
	return resp
}

// ViewHoldingAtTimePeriod retrieves a snapshot of holdings at a specific time period,
// returning empty when not found
func (p *Portfolio) ViewHoldingAtTimePeriod(ev eventtypes.EventHandler) (*holdings.Holding, error) {
	exchangeAssetPairSettings := p.exchangeAssetPairSettings[ev.GetExchange()][ev.GetAssetType()][ev.Pair()]
	if exchangeAssetPairSettings == nil {
		return nil, fmt.Errorf("%w for %v %v %v", errNoHoldings, ev.GetExchange(), ev.GetAssetType(), ev.Pair())
	}

	for i := len(exchangeAssetPairSettings.HoldingsSnapshots) - 1; i >= 0; i-- {
		if ev.GetTime().Equal(exchangeAssetPairSettings.HoldingsSnapshots[i].Timestamp) {
			return &exchangeAssetPairSettings.HoldingsSnapshots[i], nil
		}
	}

	return nil, fmt.Errorf("%w for %v %v %v at %v", errNoHoldings, ev.GetExchange(), ev.GetAssetType(), ev.Pair(), ev.GetTime())
}

// SetupCurrencySettingsMap ensures a map is created and no panics happen
func (p *Portfolio) SetupCurrencySettingsMap(exch string, a asset.Item, cp currency.Pair) (*PortfolioSettings, error) {
	if exch == "" {
		return nil, errExchangeUnset
	}
	if a == "" {
		return nil, errAssetUnset
	}
	if cp.IsEmpty() {
		return nil, errCurrencyPairUnset
	}
	if p.exchangeAssetPairSettings == nil {
		p.exchangeAssetPairSettings = make(map[string]map[asset.Item]map[currency.Pair]*PortfolioSettings)
	}
	if p.exchangeAssetPairSettings[exch] == nil {
		p.exchangeAssetPairSettings[exch] = make(map[asset.Item]map[currency.Pair]*PortfolioSettings)
	}
	if p.exchangeAssetPairSettings[exch][a] == nil {
		p.exchangeAssetPairSettings[exch][a] = make(map[currency.Pair]*PortfolioSettings)
	}
	if _, ok := p.exchangeAssetPairSettings[exch][a][cp]; !ok {
		p.exchangeAssetPairSettings[exch][a][cp] = &PortfolioSettings{}
	}

	return p.exchangeAssetPairSettings[exch][a][cp], nil
}

// // GetOpenTrades returns the latest holdings after being sorted by time
// func (p *Portfolio) GetOpenTrade() (livetrade.Details, error) {
// 	if p.openTrade.EntryPrice.IsZero() {
// 		return livetrade.Details{}, errNoOpenTrade
// 	}
// 	return p.openTrade, nil
// }

func (p *Portfolio) GetLiveMode() bool {
	return p.bot.Config.LiveMode
}

// GetLatestHoldings returns the latest holdings after being sorted by time
func (e *PortfolioSettings) GetLatestHoldings() holdings.Holding {
	if len(e.HoldingsSnapshots) == 0 {
		return holdings.Holding{}
	}

	return e.HoldingsSnapshots[len(e.HoldingsSnapshots)-1]
}

// GetHoldingsForTime returns the holdings for a time period, or an empty holding if not found
func (e *PortfolioSettings) GetHoldingsForTime(t time.Time) holdings.Holding {
	if e.HoldingsSnapshots == nil {
		// no holdings yet
		return holdings.Holding{}
	}
	for i := len(e.HoldingsSnapshots) - 1; i >= 0; i-- {
		if e.HoldingsSnapshots[i].Timestamp.Equal(t) {
			return e.HoldingsSnapshots[i]
		}
	}
	return holdings.Holding{}
}

// Value returns the total value of the latest holdings
func (e *PortfolioSettings) Value() decimal.Decimal {
	latest := e.GetLatestHoldings()
	if latest.Timestamp.IsZero() {
		return decimal.Zero
	}
	return latest.TotalValue
}

func (p *Portfolio) recordTrade(ev signal.Event) {
	direction := ev.GetDirection()
	if direction == gctorder.Sell || direction == gctorder.Buy {
		t, _ := trades.Create(ev)
		t.Update(ev)
	}
}

func (p *Portfolio) evaluateOrder(d eventtypes.Directioner, originalOrderSignal, sizedOrder *order.Order) (*order.Order, error) {
	var evaluatedOrder *order.Order
	cm, err := p.GetComplianceManager(originalOrderSignal.GetExchange(), originalOrderSignal.GetAssetType(), originalOrderSignal.Pair())
	if err != nil {
		return nil, err
	}

	evaluatedOrder, err = p.riskManager.EvaluateOrder(sizedOrder, p.GetLatestHoldingsForAllCurrencies(), cm.GetLatestSnapshot())
	if err != nil {
		originalOrderSignal.AppendReason(err.Error())
		switch d.GetDirection() {
		case gctorder.Buy:
			originalOrderSignal.Direction = eventtypes.CouldNotBuy
		case gctorder.Sell:
			originalOrderSignal.Direction = eventtypes.CouldNotSell
		case eventtypes.CouldNotBuy, eventtypes.CouldNotSell:
		default:
			originalOrderSignal.Direction = eventtypes.DoNothing
		}
		d.SetDirection(originalOrderSignal.Direction)

		return originalOrderSignal, nil
	}

	return evaluatedOrder, nil
}

func (p *Portfolio) sizeOrder(d eventtypes.Directioner, cs *ExchangeAssetPairSettings, originalOrderSignal *order.Order, sizingFunds decimal.Decimal) *order.Order {
	sizedOrder, err := p.sizeManager.SizeOrder(originalOrderSignal, sizingFunds, cs)
	if err != nil {
		originalOrderSignal.AppendReason(err.Error())
		switch originalOrderSignal.Direction {
		case gctorder.Buy:
			originalOrderSignal.Direction = eventtypes.CouldNotBuy
		case gctorder.Sell:
			originalOrderSignal.Direction = eventtypes.CouldNotSell
		default:
			originalOrderSignal.Direction = eventtypes.DoNothing
		}
		d.SetDirection(originalOrderSignal.Direction)
		return originalOrderSignal
	}

	if sizedOrder.Amount.IsZero() {
		switch originalOrderSignal.Direction {
		case gctorder.Buy:
			originalOrderSignal.Direction = eventtypes.CouldNotBuy
		case gctorder.Sell:
			originalOrderSignal.Direction = eventtypes.CouldNotSell
		default:
			originalOrderSignal.Direction = eventtypes.DoNothing
		}
		d.SetDirection(originalOrderSignal.Direction)
		originalOrderSignal.AppendReason("sized order to 0")
	}
	if d.GetDirection() == gctorder.Sell {
		// err = funds.Reserve(sizedOrder.Amount, gctorder.Sell)
		sizedOrder.AllocatedFunds = sizedOrder.Amount
	} else {
		// err = funds.Reserve(sizedOrder.Amount.Mul(sizedOrder.Price), gctorder.Buy)
		sizedOrder.AllocatedFunds = sizedOrder.Amount.Mul(sizedOrder.Price)
	}
	if err != nil {
		sizedOrder.Direction = eventtypes.DoNothing
		sizedOrder.AppendReason(err.Error())
	}
	return sizedOrder
}

// addComplianceSnapshot gets the previous snapshot of compliance events, updates with the latest fillevent
// then saves the snapshot to the c
func (p *Portfolio) addComplianceSnapshot(fillEvent fill.Event) error {
	if fillEvent == nil {
		return eventtypes.ErrNilEvent
	}
	complianceManager, err := p.GetComplianceManager(fillEvent.GetExchange(), fillEvent.GetAssetType(), fillEvent.Pair())
	if err != nil {
		return err
	}
	prevSnap := complianceManager.GetLatestSnapshot()
	fo := fillEvent.GetOrder()
	if fo != nil {
		price := decimal.NewFromFloat(fo.Price)
		amount := decimal.NewFromFloat(fo.Amount)
		fee := decimal.NewFromFloat(fo.Fee)
		snapOrder := compliance.SnapshotOrder{
			ClosePrice:          fillEvent.GetClosePrice(),
			VolumeAdjustedPrice: fillEvent.GetVolumeAdjustedPrice(),
			SlippageRate:        fillEvent.GetSlippageRate(),
			Detail:              fo,
			CostBasis:           price.Mul(amount).Add(fee),
		}
		prevSnap.Orders = append(prevSnap.Orders, snapOrder)
	}
	return complianceManager.AddSnapshot(prevSnap.Orders, fillEvent.GetTime(), fillEvent.GetOffset(), false)
}

func (p *Portfolio) setHoldingsForOffset(h *holdings.Holding, overwriteExisting bool) error {
	if h.Timestamp.IsZero() {
		return errHoldingsNoTimestamp
	}
	lookup := p.exchangeAssetPairSettings[h.Exchange][h.Asset][h.Pair]
	if lookup == nil {
		var err error
		lookup, err = p.SetupCurrencySettingsMap(h.Exchange, h.Asset, h.Pair)
		if err != nil {
			return err
		}
	}
	if overwriteExisting && len(lookup.HoldingsSnapshots) == 0 {
		return errNoHoldings
	}
	for i := len(lookup.HoldingsSnapshots) - 1; i >= 0; i-- {
		if lookup.HoldingsSnapshots[i].Offset == h.Offset {
			if overwriteExisting {
				lookup.HoldingsSnapshots[i] = *h
				return nil
			}
			return errHoldingsAlreadySet
		}
	}
	if overwriteExisting {
		return fmt.Errorf("%w at %v", errNoHoldings, h.Timestamp)
	}

	lookup.HoldingsSnapshots = append(lookup.HoldingsSnapshots, *h)
	return nil
}

// verifyOrderWithinLimits conforms the amount to fall into the minimum size and maximum size limit after reduced
func verifyOrderWithinLimits(f *fill.Fill, limitReducedAmount decimal.Decimal, cs *ExchangeAssetPairSettings) error {
	if f == nil {
		return eventtypes.ErrNilEvent
	}
	if cs == nil {
		return errNilCurrencySettings
	}
	isBeyondLimit := false
	var minMax config.MinMax
	var direction gctorder.Side
	switch f.GetDirection() {
	case gctorder.Buy:
		minMax = cs.BuySide
		direction = eventtypes.CouldNotBuy
	case gctorder.Sell:
		minMax = cs.SellSide
		direction = eventtypes.CouldNotSell
	default:
		direction = f.GetDirection()
		f.SetDirection(eventtypes.DoNothing)
		return fmt.Errorf("%w: %v", errInvalidDirection, direction)
	}
	var minOrMax, belowExceed string
	var size decimal.Decimal
	if limitReducedAmount.LessThan(minMax.MinimumSize) && minMax.MinimumSize.GreaterThan(decimal.Zero) {
		isBeyondLimit = true
		belowExceed = "below"
		minOrMax = "minimum"
		size = minMax.MinimumSize
	}
	if limitReducedAmount.GreaterThan(minMax.MaximumSize) && minMax.MaximumSize.GreaterThan(decimal.Zero) {
		isBeyondLimit = true
		belowExceed = "exceeded"
		minOrMax = "maximum"
		size = minMax.MaximumSize
	}
	if isBeyondLimit {
		f.SetDirection(direction)
		e := fmt.Sprintf("Order size %v %s %s size %v", limitReducedAmount, belowExceed, minOrMax, size)
		f.AppendReason(e)
		return fmt.Errorf("%w %v", errExceededPortfolioLimit, e)
	}
	return nil
}

func (p *Portfolio) printTradeDetails(t *livetrade.Details) {
	// log.Infoln(log.TradeManager, "trade profit", t.ProfitLossPoints)
	return
}

func (p *Portfolio) PrintPortfolioDetails() {
	log.Infoln(log.TradeManager, "portfolio details")
	active, _ := livetrade.Active()
	closed, _ := livetrade.Closed()
	log.Infoln(log.TradeManager, "active trades", len(active))
	log.Infoln(log.TradeManager, "closed trades", len(closed))

	// get strategy last updated time
	// get factor engine last updated time for each pair

	// total PL past 24 hours
	// num trades past 24 hours
	// num orders past 24 hours

	// current positions and their stats

	// log.Infoln(log.TradeManager, "active strategies")
	// log.Infoln(log.TradeManager, "active pairs")
	return
}

func (p *Portfolio) Bot() *Engine {
	return p.bot
}
