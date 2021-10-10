package engine

import (
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/compliance"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/database/repository/livetrade"
	"github.com/thrasher-corp/gocryptotrader/eventtypes"
	"github.com/thrasher-corp/gocryptotrader/eventtypes/event"
	"github.com/thrasher-corp/gocryptotrader/eventtypes/fill"
	"github.com/thrasher-corp/gocryptotrader/eventtypes/order"
	"github.com/thrasher-corp/gocryptotrader/eventtypes/signal"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	gctorder "github.com/thrasher-corp/gocryptotrader/exchanges/order"
	"github.com/thrasher-corp/gocryptotrader/log"
	"github.com/thrasher-corp/gocryptotrader/portfolio/holdings"
	"github.com/thrasher-corp/gocryptotrader/portfolio/positions"
	"github.com/thrasher-corp/gocryptotrader/portfolio/risk"
	"github.com/thrasher-corp/gocryptotrader/portfolio/trades"
	"github.com/thrasher-corp/gocryptotrader/strategies"
)

// Setup creates a portfolio manager instance and sets private fields
func SetupPortfolio(st []strategies.Handler, bot Engine, sh SizeHandler, r risk.Handler, riskFreeRate decimal.Decimal) (*Portfolio, error) {
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
	p.store.closedTrades = make(map[string][]*livetrade.Details)

	// load open trade from the database
	log.Infof(log.BackTester, "there are %d trades running", livetrade.Count())

	// load all pending and open trades from the database into positions and trades

	// what does this do?
	// bt.Datas.Setup()

	p.bot = bot
	p.sizeManager = sh
	p.riskManager = r
	p.riskFreeRate = riskFreeRate
	p.strategies = st

	// set initial opentrade/positions
	for _, s := range p.strategies {
		s.SetID(fmt.Sprintf("%s_%s", s.Name(), s.Direction()))
		p.store.positions[s.ID()] = &positions.Position{}
		p.store.closedTrades[s.ID()] = make([]*livetrade.Details, 10)
		s.SetWeight(decimal.NewFromFloat(1.5))
	}

	// load existing positions from database
	activeTrades, _ := livetrade.Active()
	for _, t := range activeTrades {
		p.store.openTrade[t.StrategyID] = &t
		pos := p.store.positions[t.StrategyID]
		pos.Active = true
	}

	return p, nil
}

// Reset returns the portfolio manager to its default state
func (p *Portfolio) Reset() {
	p.exchangeAssetPairSettings = nil
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
	switch ev.GetDecision() {
	case signal.Enter:
		if p.bot.IsLive {
			err := livetrade.Insert(livetrade.Details{
				EntryPrice:    123.0,
				StopLossPrice: 123.0,
				Status:        "PENDING",
				StrategyID:    ev.GetStrategyID(),
			})
			if err != nil {
				return nil, errStrategyIDUnset
			}
		}

		p.store.openTrade[ev.GetStrategyID()] = &livetrade.Details{Status: livetrade.Pending}

	case signal.Exit:
		// fmt.Println("exit")
	case signal.DoNothing:
		// fmt.Println("portfolio: do nothing", ev.GetReason())
		return nil, nil
	default:
		fmt.Println("no decision for", ev.GetStrategyID())
		return nil, errNoDecision
	}

	o := &order.Order{
		Base: event.Base{
			Offset:       ev.GetOffset(),
			Exchange:     ev.GetExchange(),
			Time:         ev.GetTime(),
			CurrencyPair: ev.Pair(),
			AssetType:    ev.GetAssetType(),
			Interval:     ev.GetInterval(),
			Reason:       ev.GetReason(),
		},
		Direction: ev.GetDirection(),
		Amount:    ev.GetAmount(),
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
	o.OrderType = gctorder.Market
	o.BuyLimit = ev.GetBuyLimit()
	o.SellLimit = ev.GetSellLimit()
	o.StrategyID = ev.GetStrategyID()
	sizedOrder := p.sizeOrder(ev, cs, o, decimal.NewFromFloat(100))
	sizedOrder.Amount = ev.GetAmount()

	p.recordTrade(ev)

	return p.evaluateOrder(ev, o, sizedOrder)
}

func (p *Portfolio) updatePosition(pos *positions.Position, amount decimal.Decimal) {
	pos.Amount = decimal.NewFromFloat(100.0)
}

// OnFill processes the event after an order has been placed by the exchange. Its purpose is to track holdings for future portfolio decisions.
func (p *Portfolio) OnFill(f fill.Event) (*fill.Fill, error) {
	if f == nil {
		return nil, eventtypes.ErrNilEvent
	}
	lookup := p.exchangeAssetPairSettings[f.GetExchange()][f.GetAssetType()][f.Pair()]
	if lookup == nil {
		return nil, fmt.Errorf("%w for %v %v %v", errNoPortfolioSettings, f.GetExchange(), f.GetAssetType(), f.Pair())
	}
	var err error

	// which strategy was filled?
	// what was the amount filled?
	// what was the direction of the fill?
	// what is the direction of the strategy?

	// create or update position
	for _, pos := range p.store.positions {
		if f.GetDirection() == gctorder.Sell {
			pos.Amount = pos.Amount.Sub(f.GetAmount())
		} else if f.GetDirection() == gctorder.Buy {
			pos.Amount = pos.Amount.Add(f.GetAmount())
		}

		if !pos.Amount.IsZero() {
			pos.Active = true
		} else {
			pos.Active = false
		}
	}

	// st := p.GetStrategy(f.GetStrategyID())
	// fmt.Println("strategy", st)

	t := p.store.openTrade[f.GetStrategyID()]

	if t == nil {
		p.store.openTrade[f.GetStrategyID()] = &livetrade.Details{Status: livetrade.Open}
	} else {
		// fmt.Println("open trade", t)
		t.Status = livetrade.Closed
		p.store.closedTrades[f.GetStrategyID()] = append(p.store.closedTrades[f.GetStrategyID()], t)
		p.store.openTrade[f.GetStrategyID()] = nil
	}

	// t.Status = trades.Open
	// whats the strategy id of this fill?

	// Get the holding from the previous iteration, create it if it doesn't yet have a timestamp
	h := lookup.GetHoldingsForTime(f.GetTime().Add(-f.GetInterval().Duration()))
	if !h.Timestamp.IsZero() {
		h.Update(f)
	} else {
		h = lookup.GetLatestHoldings()
		if h.Timestamp.IsZero() {
			h, err = holdings.Create(f, decimal.NewFromFloat(1000.0), p.riskFreeRate)
			if err != nil {
				return nil, err
			}
		} else {
			h.Update(f)
		}
	}
	err = p.setHoldingsForOffset(&h, true)
	if errors.Is(err, errNoHoldings) {
		err = p.setHoldingsForOffset(&h, false)
	}
	if err != nil {
		log.Error(log.BackTester, err)
	}

	err = p.addComplianceSnapshot(f)
	if err != nil {
		log.Error(log.BackTester, err)
	}

	direction := f.GetDirection()
	if direction == eventtypes.DoNothing ||
		direction == eventtypes.CouldNotBuy ||
		direction == eventtypes.CouldNotSell ||
		direction == eventtypes.MissingData ||
		direction == "" {
		fe, ok := f.(*fill.Fill)
		if !ok {
			return nil, fmt.Errorf("%w expected fill event", eventtypes.ErrInvalidDataType)
		}
		fe.ExchangeFee = decimal.Zero
		return fe, nil
	}

	fe, ok := f.(*fill.Fill)
	if !ok {
		return nil, fmt.Errorf("%w expected fill event", eventtypes.ErrInvalidDataType)
	}
	return fe, nil
}

func (p *Portfolio) GetStrategy(id string) *strategies.Handler {
	for _, s := range p.strategies {
		fmt.Println("checking", s.ID(), id)
		if s.ID() == id {
			return &s
		}
	}
	return nil
}

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
