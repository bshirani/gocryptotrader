package engine

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"gocryptotrader/communications/base"
	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/database/repository/candle"
	"gocryptotrader/database/repository/liveorder"
	"gocryptotrader/database/repository/livetrade"
	"gocryptotrader/eventtypes"
	"gocryptotrader/eventtypes/cancel"
	"gocryptotrader/eventtypes/event"
	"gocryptotrader/eventtypes/fill"
	"gocryptotrader/eventtypes/order"
	"gocryptotrader/eventtypes/signal"
	"gocryptotrader/eventtypes/submit"
	"gocryptotrader/exchange"
	"gocryptotrader/exchange/asset"
	gctorder "gocryptotrader/exchange/order"
	"gocryptotrader/log"
	"gocryptotrader/portfolio/compliance"
	"gocryptotrader/portfolio/holdings"
	"gocryptotrader/portfolio/positions"
	"gocryptotrader/portfolio/risk"
	"gocryptotrader/portfolio/slippage"
	"gocryptotrader/portfolio/strategies"
	"gocryptotrader/portfolio/trades"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

// var (
// 	errStrategyNotFound = errors.New("Strategy Not found")
// )
// Setup creates a portfolio manager instance and sets private fields
func SetupPortfolio(st []strategies.Handler, bot *Engine, cfg *config.Config) (*Portfolio, error) {
	buyRule := config.MinMax{
		MinimumSize:  cfg.PortfolioSettings.BuySide.MinimumSize,
		MaximumSize:  cfg.PortfolioSettings.BuySide.MaximumSize,
		MaximumTotal: cfg.PortfolioSettings.BuySide.MaximumTotal,
	}
	sellRule := config.MinMax{
		MinimumSize:  cfg.PortfolioSettings.SellSide.MinimumSize,
		MaximumSize:  cfg.PortfolioSettings.SellSide.MaximumSize,
		MaximumTotal: cfg.PortfolioSettings.SellSide.MaximumTotal,
	}
	sizeManager := &Size{
		BuySide:  buyRule,
		SellSide: sellRule,
	}

	portfolioRisk := &risk.Risk{
		CurrencySettings: make(map[string]map[asset.Item]map[currency.Pair]*risk.CurrencySettings),
	}
	for i := range cfg.CurrencySettings {
		if portfolioRisk.CurrencySettings[cfg.CurrencySettings[i].ExchangeName] == nil {
			portfolioRisk.CurrencySettings[cfg.CurrencySettings[i].ExchangeName] = make(map[asset.Item]map[currency.Pair]*risk.CurrencySettings)
		}
		var a asset.Item
		a, err := asset.New(cfg.CurrencySettings[i].Asset)
		if err != nil {
			return nil, fmt.Errorf(
				"%w for %v %v %v. Err %v",
				errInvalidConfigAsset,
				cfg.CurrencySettings[i].ExchangeName,
				cfg.CurrencySettings[i].Asset,
				cfg.CurrencySettings[i].Base+cfg.CurrencySettings[i].Quote,
				err)
		}
		if portfolioRisk.CurrencySettings[cfg.CurrencySettings[i].ExchangeName][a] == nil {
			portfolioRisk.CurrencySettings[cfg.CurrencySettings[i].ExchangeName][a] = make(map[currency.Pair]*risk.CurrencySettings)
		}
		var curr currency.Pair
		var b, q currency.Code
		b = currency.NewCode(cfg.CurrencySettings[i].Base)
		q = currency.NewCode(cfg.CurrencySettings[i].Quote)
		curr = currency.NewPair(b, q)
		var exch exchange.IBotExchange
		exch, err = bot.ExchangeManager.GetExchangeByName(cfg.CurrencySettings[i].ExchangeName)
		if err != nil {
			return nil, err
		}
		exchBase := exch.GetBase()
		var requestFormat currency.PairFormat
		requestFormat, err = exchBase.GetPairFormat(a, true)
		if err != nil {
			return nil, fmt.Errorf("could not format currency %v, %w", curr, err)
		}
		curr = curr.Format(requestFormat.Delimiter, requestFormat.Uppercase)
		err = exchBase.CurrencyPairs.EnablePair(a, curr)
		if err != nil && !errors.Is(err, currency.ErrPairAlreadyEnabled) {
			return nil, fmt.Errorf(
				"could not enable currency %v %v %v. Err %w",
				cfg.CurrencySettings[i].ExchangeName,
				cfg.CurrencySettings[i].Asset,
				cfg.CurrencySettings[i].Base+cfg.CurrencySettings[i].Quote,
				err)
		}
		portfolioRisk.CurrencySettings[cfg.CurrencySettings[i].ExchangeName][a][curr] = &risk.CurrencySettings{
			MaximumOrdersWithLeverageRatio: cfg.CurrencySettings[i].Leverage.MaximumOrdersWithLeverageRatio,
			MaxLeverageRate:                cfg.CurrencySettings[i].Leverage.MaximumLeverageRate,
			MaximumHoldingRatio:            cfg.CurrencySettings[i].MaximumHoldingsRatio,
		}
		if cfg.CurrencySettings[i].MakerFee.GreaterThan(cfg.CurrencySettings[i].TakerFee) {
			log.Warnf(log.Portfolio, "maker fee '%v' should not exceed taker fee '%v'. Please review config",
				cfg.CurrencySettings[i].MakerFee,
				cfg.CurrencySettings[i].TakerFee)
		}
	}

	if sizeManager == nil {
		return nil, errSizeManagerUnset
	}
	riskFreeRate := cfg.StatisticSettings.RiskFreeRate
	if riskFreeRate.IsNegative() {
		return nil, errNegativeRiskFreeRate
	}
	if portfolioRisk == nil {
		return nil, errRiskManagerUnset
	}
	p := &Portfolio{}
	p.verbose = cfg.PortfolioSettings.Verbose
	// if !p.verbose {
	// 	fmt.Println("pf not verbose")
	// 	os.Exit(123)
	// }

	// create position for every strategy
	// create open trades array for every strategy
	// you need the strategy IDS here
	p.store.positions = make(map[int]*positions.Position)
	p.store.openTrade = make(map[int]*livetrade.Details)
	p.store.openOrders = make(map[int][]*liveorder.Details)
	p.store.closedOrders = make(map[int][]*liveorder.Details)
	p.store.closedTrades = make(map[int][]*livetrade.Details)

	p.orderManager = bot.OrderManager
	p.bot = bot
	p.sizeManager = sizeManager
	p.riskManager = portfolioRisk
	p.riskFreeRate = riskFreeRate
	p.Strategies = st

	// set initial opentrade/positions
	for _, s := range p.Strategies {
		p.store.positions[s.GetID()] = &positions.Position{}
		p.store.closedTrades[s.GetID()] = make([]*livetrade.Details, 0)
		p.store.openOrders[s.GetID()] = make([]*liveorder.Details, 0)
		s.SetWeight(decimal.NewFromFloat(1.5))
	}

	if !p.bot.Settings.EnableDryRun {
		activeTrades, _ := livetrade.Active()
		for _, t := range activeTrades {
			if t.Amount.IsZero() {
				panic("trade amount is zero")
			}
			// p.getStrategyTrade(t.StrategyID)
			// set open trade
			// set position
			p.store.openTrade[t.StrategyID] = &t
			pos := p.store.positions[t.StrategyID]
			pos.Active = true
		}
		activeOrders, _ := liveorder.Active()
		for _, t := range activeOrders {
			p.store.openOrders[t.StrategyID] = append(p.store.openOrders[t.StrategyID], &t)
		}

		log.Infof(log.Portfolio, "Started Portfolio w/ %d Strategies, %d Currencies. Loaded Trades %d Orders %d", len(st), len(p.bot.CurrencySettings), len(activeTrades), len(activeOrders))
	}

	for _, cs := range p.bot.CurrencySettings {
		_, err := p.SetupCurrencySettingsMap(cs.ExchangeName, cs.AssetType, cs.CurrencyPair)
		if err != nil {
			return nil, err
		}
	}

	go p.heartBeat()

	// p.SetupCurrencySettingsMap(exch string, a asset.Item, cp currency.Pair) (*PortfolioSettings, error) {

	return p, nil
}

// Reset returns the portfolio manager to its default state
func (p *Portfolio) Reset() {
	p.exchangeAssetPairSettings = nil
}

func (p *Portfolio) OnSubmit(ev submit.Event) {
	if p.debug {
		fmt.Println("portfolio.OnSubmit", ev.GetStrategyID(), "orderID", ev.GetOrderID())
	}
	var openOrder *liveorder.Details
	if len(p.store.openOrders[ev.GetStrategyID()]) == 0 {
		for i := range p.store.openOrders {
			fmt.Println(i)
		}
		panic(fmt.Sprintf("did not store open order for strategy %d", ev.GetStrategyID()))
	}

	for _, ord := range p.store.openOrders[ev.GetStrategyID()] {
		if ord.ID == ev.GetOrderID() {
			openOrder = ord
			break
		}
	}
	if openOrder == nil {
		fmt.Println("error !!!!!! no interal openOrder id")
		return
	}

	if ev.GetIsOrderPlaced() {
		p.completeOrder(ev)
	}

	openOrder.Status = gctorder.Closed

	if !p.bot.Settings.EnableDryRun {
		id, err := liveorder.Update(openOrder)
		if err != nil || id == 0 {
			fmt.Println("error saving to db")
			os.Exit(2)
		}
	}
}

func (p *Portfolio) OnCancel(cancel cancel.Event) {
	// fmt.Println("portfolio received cancelled order", cancel)
}

func (p *Portfolio) updateStrategyTrades(ev signal.Event) {
	// update the Profit and loss of the strategy
	// move this
	trade := p.GetTradeForStrategy(ev.GetStrategyID())
	// fmt.Println("got trade for strategy", trade)
	if trade != nil {
		if trade.Side == gctorder.Buy {
			// fmt.Println("current price", ev.GetPrice(), "trade price", trade.EntryPrice, ev.GetPrice().Sub(trade.EntryPrice))
			trade.ProfitLossPoints = ev.GetPrice().Sub(trade.EntryPrice)
		} else if trade.Side == gctorder.Sell {
			trade.ProfitLossPoints = trade.EntryPrice.Sub(ev.GetPrice())
		} else {
			fmt.Println("trade is not sell or buy")
			os.Exit(2)
		}
		// if p.bot.Config.LiveMode {
		// 	p.printTradeDetails(trade)
		// }
	}
}

// OnSignal receives the event from the strategy on whether it has signalled to buy, do nothing or sell
// on buy/sell, the portfolio manager will size the order and assess the risk of the order
// if successful, it will pass on an order.Order to be used by the exchange event handler to place an order based on
// the portfolio manager's recommendations
func (p *Portfolio) OnSignal(ev signal.Event, cs *ExchangeAssetPairSettings) (*order.Order, error) {
	// if p.GetLiveMode() {
	// 	fmt.Println("UPDATE STRATEGY TRADES", ev.GetStrategyID())
	// }
	// if p.verbose {

	s, _ := p.getStrategy(ev.GetStrategyID())
	// if ev.Pair() != s.GetPair() {
	// 	fmt.Println(ev.Pair(), s.GetPair(), s.Name())
	// 	panic("updating wrong strategy/pair")
	// }

	// }
	if ev == nil || cs == nil {
		return nil, eventtypes.ErrNilArguments
	}
	if p.sizeManager == nil {
		return nil, errSizeManagerUnset
	}
	if p.riskManager == nil {
		return nil, errRiskManagerUnset
	}
	if ev.GetStrategyID() == 0 {
		return nil, errStrategyIDUnset
	}
	p.lastUpdate = ev.GetTime()
	p.updateStrategyTrades(ev)

	// validate and prepare the event
	strategyDirection, err := p.getStrategyDirection(ev.GetStrategyID())
	if err != nil {
		fmt.Println("error getting strategy direction", err)
	}
	switch ev.GetDecision() {
	case signal.Enter:
		if strategyDirection == gctorder.Sell {
			// fmt.Println("STRATEGY DIRECTION SELL, ENTER SELL")
			// ev.Base.SetDirection(gctorder.Sell)
			ev.SetDirection(gctorder.Sell)
		} else if strategyDirection == gctorder.Buy {
			// fmt.Println("STRATEGY DIRECTION BUY, ENTER BUY")
			// ev.Base.SetDirection(gctorder.Buy)
			ev.SetDirection(gctorder.Buy)
		} else {
			panic("no valid strategy side")
		}
	case signal.Exit:
		if strategyDirection == gctorder.Sell {
			// ev.Base.SetDirection(gctorder.Buy)
			// fmt.Println("STRATEGY DIRECTION SELL, EXIT w BUY", ev.GetStrategyID())
			ev.SetDirection(gctorder.Buy)
		} else if strategyDirection == gctorder.Buy {
			// ev.Base.SetDirection(gctorder.Sell)
			// fmt.Println("STRATEGY DIRECTION BUY , EXIT w SELL", ev.GetStrategyID())
			ev.SetDirection(gctorder.Sell)
		} else {
			panic("no valid strategy side")
		}
	case signal.DoNothing:
		// ev.Base.SetDirection(signal.DoNothing)
		ev.SetDirection(eventtypes.DoNothing)
	default:
		return nil, errNoDecision
	}
	if ev.GetDirection() == "" {
		return nil, errInvalidDirection
	}
	activeTrades, _ := livetrade.Active()
	// activeOrders, _ := liveorder.Active()
	maxTradeCount := 3

	// validate new entry order
	if ev.GetDecision() == signal.Enter {
		// if len(activeOrders) >= maxTradeCount {
		// 	ev.SetDirection(eventtypes.DoNothing)
		// 	ev.SetDecision(signal.DoNothing)
		// 	ev.AppendReason(fmt.Sprintf("Mgr: NOGO. DoNothing. Has %d new orders", len(activeOrders))) } else
		if len(activeTrades) >= maxTradeCount {
			ev.SetDirection(eventtypes.DoNothing)
			ev.SetDecision(signal.DoNothing)
			ev.AppendReason("ENTRY_DENIED,active_trades|")
		} else {
			ev.AppendReason("Approve Entry")
		}
	}

	// get trade for strategy

	t := p.GetTradeForStrategy(ev.GetStrategyID())
	var tradeStatus string
	if t == nil {
		tradeStatus = "NOT_IN_TRADE"
	} else {
		tradeStatus = fmt.Sprintf("IN_TRADE PL:%v", t.ProfitLossPoints.Mul(t.Amount))
	}

	// logging
	if p.verbose {
		log.Debugf(
			log.Portfolio,
			"onsig name=%d-%s-%s decision=%s status=%s reason=%s time=%s ",
			ev.GetStrategyID(),
			s.GetDirection(),
			s.Name(),
			ev.GetDecision(),
			ev.Pair(),
			tradeStatus,
			ev.GetReason(),
			ev.GetTime(),
		)
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
			StrategyID:   ev.GetStrategyID(),
		},
		Amount:     ev.GetAmount(),
		StrategyID: ev.GetStrategyID(),
	}

	lo := liveorder.Details{
		Status:     gctorder.New,
		OrderType:  gctorder.Market,
		Exchange:   ev.GetExchange(),
		StrategyID: ev.GetStrategyID(),
	}

	lookup, _ := p.bot.GetCurrencySettings(ev.GetExchange(), ev.GetAssetType(), ev.Pair())
	if lookup == nil {
		return nil, fmt.Errorf("%w for %v %v %v",
			errNoPortfolioSettings,
			ev.GetExchange(),
			ev.GetAssetType(),
			ev.Pair())
	}

	// don't allow trade if already active
	// if pos.Active {
	// 	if (sdir == gctorder.Buy && ev.GetDirection() == gctorder.Buy) || (sdir == gctorder.Sell && ev.GetDirection() == gctorder.Sell) {
	// 		return nil, errAlreadyInTrade
	// 	}
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
	o = p.sizeOrder(ev, cs, o, decimal.NewFromFloat(1.1))
	o.Amount = decimal.NewFromFloat(0.0001)
	p.recordOrder(ev, lo, o)
	p.recordTrade(ev)
	return o, nil

	// var sizingFunds decimal.Decimal
	// if ev.GetDirection() == gctorder.Sell {
	// 	sizingFunds = funds.BaseAvailable()
	// } else {
	// 	sizingFunds = funds.QuoteAvailable()
	// }
	// sizedOrder := p.sizeOrder(ev, cs, o, sizingFunds, funds)
	// // Get the holding from the previous iteration, create it if it doesn't yet have a timestamp
	// lookup2 := p.exchangeAssetPairSettings[ev.GetExchange()][ev.GetAssetType()][ev.Pair()]
	// h := lookup2.GetHoldingsForTime(o.GetTime())
	// fmt.Println("holdings", h)
	// sizedOrder.Amount = ev.GetAmount()
	// fmt.Println("sized order", sizedOrder.Amount)
	// fmt.Println("PORTFOLIO", ev.GetDirection(), ev.GetStrategyID(), ev.GetReason())
	// return p.evaluateOrder(ev, o, sizedOrder)
}

// ADD LIVE ORDER TO PORTFOLIO STORE
func (p *Portfolio) recordOrder(ev signal.Event, lo liveorder.Details, o *order.Order) error {
	// store the order
	if !p.bot.Config.DryRun {
		id, err := liveorder.Insert(lo)
		if err != nil {
			log.Errorln(log.Portfolio, "Unable to store order in database", err)
			// ev.SetDirection(signal.DoNothing)
			// ev.AppendReason(fmt.Sprintf("unable to store in database. err: %s", err))
			panic(fmt.Sprintf("unable to store order in database", err))
			return fmt.Errorf("unable to store in database. %v", err)
		}
		// fmt.Println("order signal, inserted order to db", id)
		lo.ID = id
		o.ID = id
	} else {
		id, err := uuid.NewV4()
		if err != nil {
			log.Warnf(log.OrderMgr,
				"Order manager: Unable to generate UUID. Err: %s",
				err)
		}
		lo.ID = id.String()
		o.ID = id.String()
	}

	if p.debug {
		fmt.Println("adding order for strategy:", ev.GetStrategyID())
	}
	beforeLen := len(p.store.openOrders[ev.GetStrategyID()])
	p.store.openOrders[ev.GetStrategyID()] = append(p.store.openOrders[ev.GetStrategyID()], &lo)
	afterLen := len(p.store.openOrders[ev.GetStrategyID()])
	// fmt.Println("store now has", afterLen, "orders for", ev.GetStrategyID())

	if afterLen > 1 {
		panic(fmt.Sprintf("more than one open order for strategy: %d", ev.GetStrategyID()))
	}

	// verify open order exists
	if afterLen <= beforeLen {
		fmt.Println("ERROR did not add open order")
		return fmt.Errorf("did not return open order")
	}
	return nil
}

func (p *Portfolio) updatePosition(pos *positions.Position, amount decimal.Decimal) {
	pos.Amount = decimal.NewFromFloat(100.0)
}

func (p *Portfolio) GetOrderFromStore(orderid string) *gctorder.Detail {
	// fmt.Printf("getorderfromstore LOOKUP", orderid)
	var foundOrd *gctorder.Detail
	ords, _ := p.bot.OrderManager.GetOrdersSnapshot("")
	for _, ord := range ords {
		if !strings.EqualFold(ord.InternalOrderID, orderid) {
			continue
		}
		// fmt.Println("FOUND", ord.InternalOrderID)
		foundOrd = &ord
		break
	}
	if foundOrd == nil {
		panic("order not found in store")
	}
	if foundOrd.Price == 0 {
		fmt.Println("ERROR order has no price ")
	}

	if !strings.EqualFold(foundOrd.InternalOrderID, orderid) {
		// fmt.Println("FOUND ORDER internal:", foundOrd.InternalOrderID, foundOrd.ID)
		panic("baddd")
	}

	return foundOrd
}

func (p *Portfolio) recordEnterTrade(ev fill.Event) {
	s, _ := p.getStrategy(ev.GetStrategyID())
	// fmt.Println("STRATEGY DIR", s.GetDirection())
	// fmt.Println("EV DIR", ev.GetDirection())
	// fmt.Println("ORDER ID", ev.GetOrderID())

	if ev.GetDirection() != s.GetDirection() {
		str := fmt.Sprintf("%s %s for %s", ev.GetDirection(), s.GetDirection(), ev.GetStrategyID())
		panic(str)
	}
	foundOrd := p.GetOrderFromStore(ev.GetOrderID())
	// fmt.Println("found order", foundOrd.ID, foundOrd.InternalOrderID)
	stopLossPrice := decimal.NewFromFloat(foundOrd.Price).Mul(decimal.NewFromFloat(0.9))
	p.store.positions[ev.GetStrategyID()] = &positions.Position{Active: true}

	lt := livetrade.Details{
		Status:        gctorder.Open,
		StrategyID:    ev.GetStrategyID(),
		EntryTime:     ev.GetTime(),
		EntryOrderID:  foundOrd.InternalOrderID,
		EntryPrice:    decimal.NewFromFloat(foundOrd.Price),
		StopLossPrice: stopLossPrice,
		Side:          foundOrd.Side,
		Pair:          foundOrd.Pair,
		Amount:        decimal.NewFromFloat(foundOrd.Amount),
	}

	if lt.EntryPrice.IsZero() {
		panic("EntryPrice cannot be empty")
	}
	if lt.EntryTime.IsZero() {
		panic("EntryTime cannot be empty")
	}
	if lt.Amount.IsZero() {
		panic("Amount cannot be 0")
	}
	// else {
	// 	fmt.Println("creating trade, entrytime:", lt.EntryTime)
	// }

	if !p.bot.Settings.EnableDryRun {
		id, err := livetrade.Insert(lt)
		lt.ID = id
		if err != nil {
			fmt.Println("error inserting trade", err)
			os.Exit(2)
		}
	}
	if p.bot.Settings.EnableLiveMode {
		tradeMsg := fmt.Sprintf("created trade for s:%d %v %s %s", ev.GetStrategyID(), ev.Pair(), ev.GetAmount(), ev.GetDirection())
		log.Warnf(log.Portfolio, tradeMsg)

		// lt.EntryTime.UTC().AppendFormat(e.data, l.Timestamp)

		timestampFormat := " 15:04:05 UTC"
		s, _ := p.getStrategy(ev.GetStrategyID())
		notificationMsg := fmt.Sprintf(
			"ENTER TRADE: %d\n%s %v@%v %v\n%s",
			s.GetID(),
			lt.Side,
			lt.Amount,
			lt.EntryPrice,
			lt.EntryTime.Format(timestampFormat),
			ev.GetReason())

		// fmt.Print("notification message", notificationMsg)
		p.bot.CommunicationsManager.PushEvent(base.Event{
			Type:    "trade_open\n",
			Message: notificationMsg,
		})
	}

	p.store.openTrade[ev.GetStrategyID()] = &lt
}

func (p *Portfolio) recordExitTrade(f fill.Event, t *livetrade.Details) {

	if t.Status == gctorder.Open {
		t.Status = gctorder.Closed
		t.ExitTime = f.GetTime()
		t.ExitPrice = f.GetClosePrice()

		// duplicate code from updateStrategyTrades
		if t.Side == gctorder.Buy {
			// fmt.Println("current price", ev.GetPrice(), "trade price", trade.EntryPrice, ev.GetPrice().Sub(trade.EntryPrice))
			t.ProfitLossPoints = f.GetClosePrice().Sub(t.EntryPrice)
		} else if t.Side == gctorder.Sell {
			t.ProfitLossPoints = t.EntryPrice.Sub(f.GetClosePrice())
		} else {
			fmt.Println("trade is not sell or buy")
			os.Exit(2)
		}

		// msg := fmt.Sprintf("Order manager: Strategy=%s Exchange=%s submitted order ID=%v [Ours: %v] pair=%v price=%v amount=%v side=%v type=%v for time %v.",
		// 	newOrder.StrategyID,
		// 	newOrder.Exchange,
		// 	result.OrderID,
		// 	newOrder.ID,
		// 	newOrder.Pair,
		// 	newOrder.Price,
		// 	newOrder.Amount,
		// 	newOrder.Side,
		// 	newOrder.Type,
		// 	newOrder.Date)
		// log.Debugln(log.OrderMgr, msg)

		// get order for trade

	} else {
		fmt.Println("TRYING TO CLOSE  ALREADY CLOSED TRADE. TRADE IS NOT OPEN")
		os.Exit(1)
	}
	// Velse if t.Status == livetrade.Pending {
	// 	ot := *p.store.openTrade[f.GetStrategyID()]
	// 	ot.Status = livetrade.Open
	// 	p.store.openTrade[f.GetStrategyID()] = &ot
	// }

	if t.Amount.IsZero() {
		panic("trade amount is zero")
	}
	if t.ExitTime.IsZero() {
		fmt.Println("ftime", f.GetTime())
		panic("exit time is zero")
	}
	if t.ExitPrice.IsZero() {
		panic("exit price is zero")
	}
	if !p.bot.Settings.EnableDryRun {
		id, err := livetrade.Update(t)
		t.ID = int(id)
		if err != nil || id == 0 {
			fmt.Println("error saving to db")
			os.Exit(2)
		}
	}
	p.store.closedTrades[f.GetStrategyID()] = append(p.store.closedTrades[f.GetStrategyID()], t)
	p.store.openTrade[f.GetStrategyID()] = nil
	p.store.positions[f.GetStrategyID()] = &positions.Position{Active: false}

	if p.bot.Settings.EnableLiveMode {
		s, _ := p.getStrategy(f.GetStrategyID())

		timeFormat := "15:05:05"
		// fmt.Println(time.Now().Format("02-Jan-2006 15:04:05"))
		notificationMsg := fmt.Sprintf(
			"EXIT TRADE: %s\nEntry:%s %v@%v@%v\nExit:%v@%v\nReason: %s\nProfit: %v",
			s.GetID(),
			t.Side,
			t.Amount,
			t.EntryTime.Format(timeFormat),
			t.EntryPrice,
			t.ExitTime.Format(timeFormat),
			t.ExitPrice,
			f.GetReason(),
			t.ProfitLossPoints,
		)

		p.bot.CommunicationsManager.PushEvent(base.Event{
			Type:    "trade_close",
			Message: notificationMsg,
		})
	}
}

// OnFill processes the event after an order has been placed by the exchange. Its purpose is to track holdings for future portfolio decisions.
func (p *Portfolio) OnFill(f fill.Event) {
	if f.GetStrategyID() == 0 {
		fmt.Println("fill has no strategy ID")
		os.Exit(2)
	}
	if f.GetOrderID() == "" {
		fmt.Println("fill has no order ID")
		os.Exit(2)
	}

	// update trades and orders here
	t := p.store.openTrade[f.GetStrategyID()]
	if t == nil {
		// fmt.Println("PF ON fILL creating NEW TRADE")
		p.recordEnterTrade(f)

	} else if t.Status == gctorder.Open {
		// fmt.Println("PF ONFILL CLOSING TRADE")
		p.recordExitTrade(f, t)
	}

}

func (p *Portfolio) completeOrder(ev submit.Event) {
	// if p.verbose {
	// 	log.Infoln(log.Portfolio, "completing order", ev.GetStrategyID())
	// }
	// fmt.Println("COMPLETING ORDER for:", ev.GetStrategyID())
	// fmt.Println("open orders", len(p.store.openOrders[ev.GetStrategyID()]))
	order := p.store.openOrders[ev.GetStrategyID()][0]
	p.store.closedOrders[ev.GetStrategyID()] = append(p.store.closedOrders[ev.GetStrategyID()], order)
	p.store.openOrders[ev.GetStrategyID()] = make([]*liveorder.Details, 0)
	// fmt.Println(ev.GetStrategyID(), " now has ", len(p.store.closedOrders[ev.GetStrategyID()]))

	// p.store.positions[ev.GetStrategyID()] = &positions.Position{Active: false}
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

func (p *Portfolio) GetVerbose() bool {
	return p.verbose
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
		fmt.Println("updateholdings")
		panic("give me the stack")
		return fmt.Errorf("%w for %v %v %v",
			errNoPortfolioSettings,
			ev.GetExchange(),
			ev.GetAssetType(),
			ev.Pair())
	}
	h := lookup.GetLatestHoldings()
	if h.Timestamp.IsZero() {
		var err error
		h, err = holdings.Create(ev, decimal.NewFromFloat(7.0), p.riskFreeRate)
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
func (p *Portfolio) UpdateTrades(ev eventtypes.DataEventHandler) error {
	if ev == nil {
		return fmt.Errorf("UPDATETRADES no event")
	}

	_, ok := p.exchangeAssetPairSettings[ev.GetExchange()][ev.GetAssetType()][ev.Pair()]
	if !ok {
		return fmt.Errorf("%w for %v %v %v",
			errNoPortfolioSettings,
			ev.GetExchange(),
			ev.GetAssetType(),
			ev.Pair())
	}
	return nil
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
// 	return p.Strategies
// }

func (p *Portfolio) GetPositionForStrategy(sid int) *positions.Position {
	return p.store.positions[sid]
}

func (p *Portfolio) GetTradeForStrategy(sid int) *livetrade.Details {
	return p.store.openTrade[sid]
}

func (p *Portfolio) GetOpenOrdersForStrategy(sid int) []*liveorder.Details {
	if p.debug {
		fmt.Println("get open orders", sid, len(p.store.openOrders[sid]))
	}
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

	// pos := p.GetPositionForStrategy(p.Strategies[0].ID())
	// fmt.Println("position:", pos.Amount)
	// pos.Amount = decimal.NewFromFloat(123.0)
	// pos.Active = false

	// return nil
	_, ok := p.exchangeAssetPairSettings[ev.GetExchange()][ev.GetAssetType()][ev.Pair()]
	if !ok {
		return
		// return fmt.Errorf("%w for %v %v %v",
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

// func (p *Portfolio) printTradeDetails(t *livetrade.Details) {
// 	secondsInTrade := int64(p.lastUpdate.Sub(t.EntryTime).Seconds())
// 	log.Infof(log.Portfolio, "%s trade: pl:%v time:%d\n", t.StrategyID, t.ProfitLossPoints, secondsInTrade)
// 	return
// }

func (p *Portfolio) PrintPortfolioDetails() {
	// log.Infoln(log.Portfolio, "portfolio details", p.lastUpdate)
	// active, _ := livetrade.Active()
	// activeOrders, _ := liveorder.Active()
	// closed, _ := livetrade.Closed()

	// for _, t := range active {
	// 	p.printTradeDetails(&t)
	// }
	// log.Infof(log.Portfolio, "orders:%d open_trades:%d closed_trades:%d", len(activeOrders), len(active), len(closed))

	// get strategy last updated time
	// get factor engine last updated time for each pair

	// total PL past 24 hours
	// num trades past 24 hours
	// num orders past 24 hours

	// current positions and their stats

	// log.Infoln(log.Portfolio, "active strategies")
	// log.Infoln(log.Portfolio, "active pairs")
	return
}

func (p *Portfolio) Bot() *Engine {
	return p.bot
}

// getFees will return an exchange's fee rate from GCT's wrapper function
func getFees(ctx context.Context, exch exchange.IBotExchange, fPair currency.Pair) (makerFee, takerFee decimal.Decimal) {
	fTakerFee, err := exch.GetFeeByType(ctx,
		&exchange.FeeBuilder{FeeType: exchange.OfflineTradeFee,
			Pair:          fPair,
			IsMaker:       false,
			PurchasePrice: 1,
			Amount:        1,
		})
	if err != nil {
		log.Errorf(log.Portfolio, "Could not retrieve taker fee for %v. %v", exch.GetName(), err)
	}

	fMakerFee, err := exch.GetFeeByType(ctx,
		&exchange.FeeBuilder{
			FeeType:       exchange.OfflineTradeFee,
			Pair:          fPair,
			IsMaker:       true,
			PurchasePrice: 1,
			Amount:        1,
		})
	if err != nil {
		log.Errorf(log.Portfolio, "Could not retrieve maker fee for %v. %v", exch.GetName(), err)
	}

	return decimal.NewFromFloat(fMakerFee), decimal.NewFromFloat(fTakerFee)
}

func (p *Portfolio) sizeOfflineOrder(high, low, volume decimal.Decimal, cs *ExchangeAssetPairSettings, f *fill.Fill) (adjustedPrice, adjustedAmount decimal.Decimal, err error) {
	if cs == nil || f == nil {
		return decimal.Zero, decimal.Zero, eventtypes.ErrNilArguments
	}
	// provide history and estimate volatility
	slippageRate := slippage.EstimateSlippagePercentage(cs.MinimumSlippageRate, cs.MaximumSlippageRate)
	if cs.SkipCandleVolumeFitting {
		f.VolumeAdjustedPrice = f.ClosePrice
		adjustedAmount = f.Amount
	} else {
		f.VolumeAdjustedPrice, adjustedAmount = ensureOrderFitsWithinHLV(f.ClosePrice, f.Amount, high, low, volume)
		if !adjustedAmount.Equal(f.Amount) {
			f.AppendReason(fmt.Sprintf("Order size shrunk from %v to %v to fit candle", f.Amount, adjustedAmount))
		}
	}

	if adjustedAmount.LessThanOrEqual(decimal.Zero) && f.Amount.GreaterThan(decimal.Zero) {
		return decimal.Zero, decimal.Zero, fmt.Errorf("amount set to 0, %w", errDataMayBeIncorrect)
	}
	adjustedPrice = applySlippageToPrice(f.GetDirection(), f.GetVolumeAdjustedPrice(), slippageRate)

	f.Slippage = slippageRate.Mul(decimal.NewFromInt(100)).Sub(decimal.NewFromInt(100))
	f.ExchangeFee = calculateExchangeFee(adjustedPrice, adjustedAmount, cs.TakerFee)
	return adjustedPrice, adjustedAmount, nil
}

func applySlippageToPrice(direction gctorder.Side, price, slippageRate decimal.Decimal) decimal.Decimal {
	adjustedPrice := price
	if direction == gctorder.Buy {
		adjustedPrice = price.Add(price.Mul(decimal.NewFromInt(1).Sub(slippageRate)))
	} else if direction == gctorder.Sell {
		adjustedPrice = price.Mul(slippageRate)
	}
	return adjustedPrice
}

func ensureOrderFitsWithinHLV(slippagePrice, amount, high, low, volume decimal.Decimal) (adjustedPrice, adjustedAmount decimal.Decimal) {
	adjustedPrice = slippagePrice
	if adjustedPrice.LessThan(low) {
		adjustedPrice = low
	}
	if adjustedPrice.GreaterThan(high) {
		adjustedPrice = high
	}
	if volume.LessThanOrEqual(decimal.Zero) {
		return adjustedPrice, adjustedAmount
	}
	currentVolume := amount.Mul(adjustedPrice)
	if currentVolume.GreaterThan(volume) {
		// reduce the volume to not exceed the total volume of the candle
		// it is slightly less than the total to still allow for the illusion
		// that open high low close values are valid with the remaining volume
		// this is very opinionated
		currentVolume = volume.Mul(decimal.NewFromFloat(0.99999999))
	}
	// extract the amount from the adjusted volume
	adjustedAmount = currentVolume.Div(adjustedPrice)

	return adjustedPrice, adjustedAmount
}

func calculateExchangeFee(price, amount, fee decimal.Decimal) decimal.Decimal {
	return fee.Mul(price).Mul(amount)
}

func reduceAmountToFitPortfolioLimit(adjustedPrice, amount, sizedPortfolioTotal decimal.Decimal, side gctorder.Side) decimal.Decimal {
	// switch side {
	// case gctorder.Buy:
	// 	if adjustedPrice.Mul(amount).GreaterThan(sizedPortfolioTotal) {
	// 		// adjusted amounts exceeds portfolio manager's allowed funds
	// 		// the amount has to be reduced to equal the sizedPortfolioTotal
	// 		amount = sizedPortfolioTotal.Div(adjustedPrice)
	// 	}
	// case gctorder.Sell:
	// 	if amount.GreaterThan(sizedPortfolioTotal) {
	// 		amount = sizedPortfolioTotal
	// 	}
	// }
	return amount
}

func (p *Portfolio) heartBeat() {
	// tick := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-p.shutdown:
			return
			// case <-tick.C:
			// 	p.PrintPortfolioDetails()
		}
	}
	// time.Sleep(time.Second * 10)
	// // fmt.Println("........................HEARTBEAT")
	// exchanges, _ := p.bot.ExchangeManager.GetExchanges()
	// ex := exchanges[0]
	// fmt.Println("subscribing to ", p.bot.CurrencySettings[0])
	// pipe, err := ticker.SubscribeToExchangeTickers(ex.GetName())
	// if err != nil {
	// 	fmt.Println(".........error subscribing to ticker", err)
	// 	// wait and retry
	// }
	//
	// // defer func() {
	// // }()
	//
	// for {
	// 	select {
	// 	case <-p.shutdown:
	// 		pipeErr := pipe.Release()
	// 		if pipeErr != nil {
	// 			log.Error(log.DispatchMgr, pipeErr)
	// 		}
	// 		return
	// 	case data, ok := <-pipe.C:
	// 		if !ok {
	// 			fmt.Println("error dispatch system")
	// 			return
	// 		}
	// 		t := (*data.(*interface{})).(ticker.Price)
	// 		fmt.Println(t.Pair.String(), t.High, t.Low)
	// 	}
	// 	// err := stream.Send(&gctrpc.TickerResponse{
	// 	// 	Pair: &gctrpc.CurrencyPair{
	// 	// 		Base:      t.Pair.Base.String(),
	// 	// 		Quote:     t.Pair.Quote.String(),
	// 	// 		Delimiter: t.Pair.Delimiter},
	// 	// 	LastUpdated: s.unixTimestamp(t.LastUpdated),
	// 	// 	Last:        t.Last,
	// 	// 	High:        t.High,
	// 	// 	Low:         t.Low,
	// 	// 	Bid:         t.Bid,
	// 	// 	Ask:         t.Ask,
	// 	// 	Volume:      t.Volume,
	// 	// 	PriceAth:    t.PriceATH,
	// 	// })
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// }
	// fmt.Println("finished")
	// // if err != nil {
	// // 	return err
	// // }
	//
	// // p.wg.Add(1)
	// // tick := time.NewTicker(time.Second * 5)
	// // defer func() {
	// // 	tick.Stop()
	// // 	p.wg.Done()
	// // }()
	// // for {
	// // 	select {
	// // 	case <-p.shutdown:
	// // 		return
	// // 	case <-tick.C:
	// // 		exchanges, err := p.bot.ExchangeManager.GetExchanges()
	// // 		for _, ex := range exchanges {
	// // 			if err != nil {
	// // 				log.Infoln(log.Portfolio, "error getting tick", err)
	// // 			}
	// //
	// // 			for _, cp := range p.bot.CurrencySettings {
	// // 				tick, _ := ex.FetchTicker(context.Background(), cp.CurrencyPair, asset.Spot)
	// // 				t1 := time.Now()
	// // 				// ticker := m.currencyPairs[x].Ticker
	// // 				secondsAgo := int(t1.Sub(tick.LastUpdated).Seconds())
	// // 				if secondsAgo > 10 {
	// // 					log.Warnln(log.Portfolio, cp.CurrencyPair, tick.Last, secondsAgo)
	// // 				} else {
	// // 					log.Infoln(log.Portfolio, cp.CurrencyPair, tick.Last, secondsAgo)
	// // 				}
	// // 			}
	// // 		}
	// // 		// p.PrintTradingDetails()
	// // 	}
	// // }
}

func (p *Portfolio) PrintTradingDetails() {
	// fmt.Println("strategies running", len(p.Strategies))
	log.Infoln(log.Portfolio, len(p.Strategies), "strategies running")

	for _, cs := range p.bot.CurrencySettings {
		// fmt.Println("currency", cs)
		retCandle, _ := candle.Series(cs.ExchangeName,
			cs.CurrencyPair.Base.String(), cs.CurrencyPair.Quote.String(),
			60, cs.AssetType.String(), time.Now().Add(time.Minute*-5), time.Now())
		var lastCandle candle.Candle
		if len(retCandle.Candles) > 0 {
			lastCandle = retCandle.Candles[len(retCandle.Candles)-1]
		}
		secondsAgo := int(time.Now().Sub(lastCandle.Timestamp).Seconds())
		if secondsAgo > 60 {
			log.Infoln(log.Portfolio, cs.CurrencyPair, "last updated", secondsAgo, "seconds ago")
		}
		// else {
		// 	log.Debugln(log.StrategyMgr, cs.CurrencyPair, "last updated", secondsAgo, "seconds ago")
		// }
	}
}

func (p *Portfolio) getStrategyDirection(strategyID int) (gctorder.Side, error) {
	strategy, err := p.getStrategy(strategyID)
	return strategy.GetDirection(), err
}

func (p *Portfolio) getStrategy(strategyID int) (strategies.Handler, error) {
	for _, s := range p.Strategies {
		// fmt.Println("lookup strategy", strategyID)
		if s.GetID() == strategyID {
			// fmt.Println("returning", strategyID, s.Name(), s.GetDirection(), s.GetPair())
			return s, nil
		}
	}
	return nil, fmt.Errorf("strategy not found")
}
