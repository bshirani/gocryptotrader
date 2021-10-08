package portfolio

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/backtester/common"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/exchange"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/compliance"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/holdings"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/positions"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/risk"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/settings"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventhandlers/portfolio/trades"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventtypes/event"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventtypes/fill"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventtypes/order"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventtypes/signal"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/engine"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	gctorder "github.com/thrasher-corp/gocryptotrader/exchanges/order"
	"github.com/thrasher-corp/gocryptotrader/log"
)

// Setup creates a portfolio manager instance and sets private fields
func Setup(bot engine.Engine, sh SizeHandler, r risk.Handler, riskFreeRate decimal.Decimal) (*Portfolio, error) {
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
	p.store.Positions = make(map[int64]*positions.Position)
	p.store.Positions[123] = &positions.Position{}
	p.bot = bot
	p.sizeManager = sh
	p.riskManager = r
	p.riskFreeRate = riskFreeRate

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
func (p *Portfolio) OnSignal(ev signal.Event, cs *exchange.Settings) (*order.Order, error) {

	switch ev.GetDecision() {
	case signal.Enter:
		fmt.Println("enter")
	case signal.Exit:
		fmt.Println("exit")
	case signal.DoNothing:
		return nil, nil
	default:
		return nil, errNoDecision
	}

	if ev == nil || cs == nil {
		return nil, common.ErrNilArguments
	}
	if p.sizeManager == nil {
		return nil, errSizeManagerUnset
	}
	if p.riskManager == nil {
		return nil, errRiskManagerUnset
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

	if ev.GetDirection() == common.DoNothing ||
		ev.GetDirection() == common.MissingData ||
		ev.GetDirection() == common.TransferredFunds ||
		ev.GetDirection() == "" {
		return nil, nil
	}

	o.Price = ev.GetPrice()
	o.OrderType = gctorder.Market
	o.BuyLimit = ev.GetBuyLimit()
	o.SellLimit = ev.GetSellLimit()
	// prevHolding := lookup.GetLatestHoldings()
	// var sizingFunds decimal.Decimal
	// if ev.GetDirection() == gctorder.Sell {
	// 	sizingFunds = prevHolding.BaseValue
	// }
	sizedOrder := p.sizeOrder(ev, cs, o, decimal.NewFromFloat(100))
	sizedOrder.Amount = decimal.NewFromFloat(0.001)

	p.recordTrade(ev)

	return p.evaluateOrder(ev, o, sizedOrder)
}

func (p *Portfolio) updatePosition(pos *positions.Position, amount decimal.Decimal) {
	pos.Amount = decimal.NewFromFloat(100.0)
}

// OnFill processes the event after an order has been placed by the exchange. Its purpose is to track holdings for future portfolio decisions.
func (p *Portfolio) OnFill(f fill.Event) (*fill.Fill, error) {
	if f == nil {
		return nil, common.ErrNilEvent
	}
	lookup := p.exchangeAssetPairSettings[f.GetExchange()][f.GetAssetType()][f.Pair()]
	if lookup == nil {
		return nil, fmt.Errorf("%w for %v %v %v", errNoPortfolioSettings, f.GetExchange(), f.GetAssetType(), f.Pair())
	}
	var err error

	// which strategy was filled?
	// what was the amount filled?
	// what was the direction of the fill?

	// create or update position
	for i, x := range p.store.Positions {
		if i == 123 {
			pos := p.store.Positions[i]
			pos.Amount = x.Amount.Add(decimal.NewFromFloat(10.0))
			if !pos.Amount.IsZero() {
				pos.Active = true
			}
		}
	}
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
	if direction == common.DoNothing ||
		direction == common.CouldNotBuy ||
		direction == common.CouldNotSell ||
		direction == common.MissingData ||
		direction == "" {
		fe, ok := f.(*fill.Fill)
		if !ok {
			return nil, fmt.Errorf("%w expected fill event", common.ErrInvalidDataType)
		}
		fe.ExchangeFee = decimal.Zero
		return fe, nil
	}

	fe, ok := f.(*fill.Fill)
	if !ok {
		return nil, fmt.Errorf("%w expected fill event", common.ErrInvalidDataType)
	}
	return fe, nil
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
func (p *Portfolio) UpdateHoldings(ev common.DataEventHandler) error {
	if ev == nil {
		return common.ErrNilEvent
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
func (p *Portfolio) UpdateTrades(ev common.DataEventHandler) {
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

func (p *Portfolio) GetPositionForStrategy(sid int64) *positions.Position {
	return p.store.Positions[sid]
	// return pos
}

// UpdatePositions updates the strategy's position for the data event
func (p *Portfolio) UpdatePositions(ev common.DataEventHandler) {
	if ev == nil {
		return
	}

	// for _, s := range p.GetStrategies() {
	// 	// update the strategies positions
	// 	fmt.Println("update position for ", s.Name)
	// }

	// portfolio has many strategies
	// we keep the position for each strategy
	for i, p := range p.store.Positions {
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
func (p *Portfolio) ViewHoldingAtTimePeriod(ev common.EventHandler) (*holdings.Holding, error) {
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
func (p *Portfolio) SetupCurrencySettingsMap(exch string, a asset.Item, cp currency.Pair) (*settings.Settings, error) {
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
		p.exchangeAssetPairSettings = make(map[string]map[asset.Item]map[currency.Pair]*settings.Settings)
	}
	if p.exchangeAssetPairSettings[exch] == nil {
		p.exchangeAssetPairSettings[exch] = make(map[asset.Item]map[currency.Pair]*settings.Settings)
	}
	if p.exchangeAssetPairSettings[exch][a] == nil {
		p.exchangeAssetPairSettings[exch][a] = make(map[currency.Pair]*settings.Settings)
	}
	if _, ok := p.exchangeAssetPairSettings[exch][a][cp]; !ok {
		p.exchangeAssetPairSettings[exch][a][cp] = &settings.Settings{}
	}

	return p.exchangeAssetPairSettings[exch][a][cp], nil
}

// // GetOpenTrades returns the latest holdings after being sorted by time
// func (p *Portfolio) GetOpenTrade() (trades.Trade, error) {
// 	if p.openTrade.EntryPrice.IsZero() {
// 		return trades.Trade{}, errNoOpenTrade
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

func (p *Portfolio) evaluateOrder(d common.Directioner, originalOrderSignal, sizedOrder *order.Order) (*order.Order, error) {
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
			originalOrderSignal.Direction = common.CouldNotBuy
		case gctorder.Sell:
			originalOrderSignal.Direction = common.CouldNotSell
		case common.CouldNotBuy, common.CouldNotSell:
		default:
			originalOrderSignal.Direction = common.DoNothing
		}
		d.SetDirection(originalOrderSignal.Direction)

		return originalOrderSignal, nil
	}

	return evaluatedOrder, nil
}

func (p *Portfolio) sizeOrder(d common.Directioner, cs *exchange.Settings, originalOrderSignal *order.Order, sizingFunds decimal.Decimal) *order.Order {
	sizedOrder, err := p.sizeManager.SizeOrder(originalOrderSignal, sizingFunds, cs)
	if err != nil {
		originalOrderSignal.AppendReason(err.Error())
		switch originalOrderSignal.Direction {
		case gctorder.Buy:
			originalOrderSignal.Direction = common.CouldNotBuy
		case gctorder.Sell:
			originalOrderSignal.Direction = common.CouldNotSell
		default:
			originalOrderSignal.Direction = common.DoNothing
		}
		d.SetDirection(originalOrderSignal.Direction)
		return originalOrderSignal
	}

	if sizedOrder.Amount.IsZero() {
		switch originalOrderSignal.Direction {
		case gctorder.Buy:
			originalOrderSignal.Direction = common.CouldNotBuy
		case gctorder.Sell:
			originalOrderSignal.Direction = common.CouldNotSell
		default:
			originalOrderSignal.Direction = common.DoNothing
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
		sizedOrder.Direction = common.DoNothing
		sizedOrder.AppendReason(err.Error())
	}
	return sizedOrder
}

// addComplianceSnapshot gets the previous snapshot of compliance events, updates with the latest fillevent
// then saves the snapshot to the c
func (p *Portfolio) addComplianceSnapshot(fillEvent fill.Event) error {
	if fillEvent == nil {
		return common.ErrNilEvent
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
