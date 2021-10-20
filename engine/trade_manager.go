package engine

// . "github.com/volatiletech/sqlboiler/v4/queries/qm"
import (
	"context"
	"errors"
	"fmt"
	"gocryptotrader/common"
	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/data"
	"gocryptotrader/data/kline/database"
	"gocryptotrader/database/repository/candle"
	"gocryptotrader/database/repository/datahistoryjob"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	gctdatabase "gocryptotrader/database"
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

	gctlog "gocryptotrader/log"
	"gocryptotrader/portfolio/report"
	"gocryptotrader/portfolio/statistics"
	"gocryptotrader/portfolio/statistics/currencystatistics"
	"gocryptotrader/portfolio/strategies"

	"github.com/fatih/color"
	"github.com/shopspring/decimal"
)

func NewTradeManager(bot *Engine) (*TradeManager, error) {
	wd, err := os.Getwd()
	configPath := filepath.Join(wd, "backtester", "config", "trend.strat")
	btcfg, err := config.ReadConfigFromFile(configPath)
	if err != nil {
		fmt.Println("error", err)
		return nil, err
	}
	return NewTradeManagerFromConfig(btcfg, "xx", "xx", bot)
}

func (tm *TradeManager) Reset() {
	tm.EventQueue.Reset()
	tm.Datas.Reset()
	tm.Portfolio.Reset()
	tm.Statistic.Reset()
	// tm.Exchange.Reset()
	// reset live trades here
	// tm.bot = nil
}

func NewTradeManagerFromConfig(cfg *config.Config, templatePath, output string, bot *Engine) (*TradeManager, error) {
	log.Debugln(log.TradeMgr, "TradeManager: Initializing... dry run", bot.Config.DryRun)
	if cfg == nil {
		return nil, errNilConfig
	}
	if bot == nil {
		return nil, errNilBot
	}
	if len(cfg.CurrencySettings) < 1 {
		fmt.Println("no cs")
		os.Exit(123)
	}
	tm := &TradeManager{
		shutdown: make(chan struct{}),
	}

	tm.cfg = *cfg
	tm.verbose = cfg.TradeManager.Verbose

	tm.EventQueue = &Holder{}
	reports := &report.Data{
		Config:       cfg,
		TemplatePath: templatePath,
		OutputPath:   output,
	}
	tm.Reports = reports
	tm.bot = bot
	var err error
	if bot.OrderManager == nil {
		bot.FakeOrderManager, err = SetupFakeOrderManager(
			bot,
			bot.ExchangeManager,
			bot.CommunicationsManager,
			&bot.ServicesWG,
			bot.Settings.Verbose)
		if err != nil {
			log.Errorf(log.Global, "Fake Order manager unable to setup: %s", err)
		} else {
			err = bot.FakeOrderManager.Start()
			bot.OrderManager = bot.FakeOrderManager

			if err != nil {
				log.Errorf(log.Global, "Fake Order manager unable to start: %s", err)
			}
		}
	}
	tm.OrderManager = bot.OrderManager
	tm.syncManager = bot.currencyPairSyncer

	tm.Datas = &data.HandlerPerCurrency{}
	tm.Datas.Setup()
	if !tm.bot.Config.LiveMode {
		err = tm.startOfflineServices()
	}
	if err != nil {
		fmt.Println("failed to setup bot", err)
		// return err
	}

	if tm.bot.Settings.EnableTrading {
		tm.initializePortfolio(cfg)
	}
	if err != nil {
		fmt.Println("error setting up tm", err)
		os.Exit(123)
	}

	// fmt.Println("done setting up bot with", len(tm.bot.CurrencySettings), "currencies")
	if len(tm.bot.CurrencySettings) < 1 {
		fmt.Println("no currency settings")
		os.Exit(123)
	}

	tm.setOrderManagerCallbacks()

	return tm, err
}

func (tm *TradeManager) ExecuteOrder(o order.Event, data data.Handler, om ExecutionHandler) (submit.Event, error) {
	// u, _ := uuid.NewV4()
	// var orderID string
	priceFloat, _ := o.GetPrice().Float64()
	a, _ := o.GetAmount().Float64()
	fee, _ := o.GetExchangeFee().Float64()

	submission := &gctorder.Submit{
		Price:       priceFloat,
		Amount:      a,
		Fee:         fee,
		Exchange:    o.GetExchange(),
		ID:          o.GetID(),
		Side:        o.GetDirection(),
		AssetType:   o.GetAssetType(),
		Date:        o.GetTime(),
		LastUpdated: o.GetTime(),
		Pair:        o.Pair(),
		Type:        gctorder.Market,
		StrategyID:  o.GetStrategyID(),
	}

	omr, err := om.Submit(context.TODO(), submission)
	if err != nil {
		fmt.Println("tm: ERROR order manager submission", err, submission.Side, omr)
	}

	// fmt.Println("tm: order manager response", omr)

	// if order is placed, update the status of the order to Open

	// update order event order_id, status

	// add the submission to the store

	// we can call on submit manually here

	// if o.GetStrategyID() == "" {
	// 	return nil, fmt.Errorf("exchange: order has no strategyid")
	// }

	// update the store with the submission ID
	ords, _ := om.GetOrdersSnapshot("")
	var internalOrderID string
	for i := range ords {
		if ords[i].ID != o.GetID() {
			continue
		}
		internalOrderID = ords[i].InternalOrderID
		internalOrderID = ords[i].InternalOrderID
		ords[i].StrategyID = o.GetStrategyID()
		ords[i].Date = o.GetTime()
		ords[i].LastUpdated = o.GetTime()
		ords[i].CloseTime = o.GetTime()
	}

	ev := &submit.Submit{
		IsOrderPlaced:   omr.IsOrderPlaced,
		InternalOrderID: internalOrderID,
		StrategyID:      o.GetStrategyID(),
	} // transform into submit event

	if ev.GetInternalOrderID() == "" {
		log.Errorln(log.TradeMgr, "error: order has no internal order id")
	}

	if ev.IsOrderPlaced {
		// fmt.Println("TM ORDERPLACED, create fill event")
		tm.onFill(omr)
	} else {
		fmt.Println("TM ERROR: ORDERPLACED NOT")
	}

	return ev, nil
}

// load data if necessary
// in backtest mode...
// in live mode...

func (tm *TradeManager) Run() error {
	// log.Debugf(log.TradeMgr, "TradeManager Running. Warmup: %v\n", warmup)
	// for _, cs := range tm.bot.CurrencySettings {
	// 	dbData, err := database.LoadData(
	// 		time.Now().Add(time.Minute-30),
	// 		time.Now(),
	// 		time.Minute,
	// 		cs.ExchangeName,
	// 		0,
	// 		cs.CurrencyPair,
	// 		cs.AssetType)
	//
	// 	if err != nil {
	// 		fmt.Println("error loading db data", err)
	// 		// create a data history request if there isn't one already
	// 		os.Exit(123)
	// 	} else {
	// 		fmt.Println("loaded data for", cs.CurrencyPair)
	// 	}
	//
	// 	tm.Datas.SetDataForCurrency(cs.ExchangeName, cs.AssetType, cs.CurrencyPair, dbData)
	// 	dbData.Load()
	// }
dataLoadingIssue:
	for ev := tm.EventQueue.NextEvent(); ; ev = tm.EventQueue.NextEvent() {
		if ev == nil {
			dataHandlerMap := tm.Datas.GetAllData()
			for _, exchangeMap := range dataHandlerMap {
				for _, assetMap := range exchangeMap {
					// var hasProcessedData bool
					tm.hasHandledEvent = false
					for _, dataHandler := range assetMap {
						d := dataHandler.Next()
						if d == nil {
							// log.Errorf(log.TradeMgr, "No data found for %v", currencyPair)
							// if !tm.hasHandledEvent {
							// 	log.Errorf(log.TradeMgr, "Unable to perform `Next` for %v %v %v", exchangeName, assetItem, currencyPair)
							// }
							break dataLoadingIssue
						}
						tm.hasHandledEvent = true
						tm.EventQueue.AppendEvent(d)
					}
				}
			}
		}
		if ev != nil {
			err := tm.handleEvent(ev)
			if err != nil {
				fmt.Println("error handling event", err)
				return err
			}
		}
		if !tm.hasHandledEvent {
			tm.hasHandledEvent = true
		}
	}

	return nil
}

func (tm *TradeManager) Start() error {
	tm.bot.WaitForInitialCurrencySync()
	if !atomic.CompareAndSwapInt32(&tm.started, 0, 1) {
		return fmt.Errorf("backtester %w", ErrSubSystemAlreadyStarted)
	}
	log.Debugf(log.TradeMgr, "TradeManager  %s", MsgSubSystemStarting)
	tm.shutdown = make(chan struct{})
	// go tm.heartBeat()
	log.Debugln(log.TradeMgr, "Running Live")
	go tm.runLive()
	return nil
}

func (tm *TradeManager) Stop() error {
	if tm == nil {
		return ErrNilSubsystem
	}
	if !atomic.CompareAndSwapInt32(&tm.started, 1, 0) {
		return ErrSubSystemNotStarted
	}

	log.Debugln(log.TradeMgr, "Backtester Stopping...")
	if tm.bot.OrderManager.IsRunning() {
		tm.bot.OrderManager.Stop()
	}
	for _, s := range tm.Strategies {
		s.Stop()
	}
	close(tm.shutdown)
	tm.bot.TradeManager = nil
	tm.wg.Wait()
	log.Debugln(log.TradeMgr, "Backtester Stopped.")
	return nil
}

func (b *TradeManager) IsRunning() bool {
	if b == nil {
		return false
	}
	return atomic.LoadInt32(&b.started) == 1
}

func (e *Holder) Reset() {
	e.Queue = nil
}

func (e *Holder) AppendEvent(i eventtypes.EventHandler) {
	e.Queue = append(e.Queue, i)
}

func (e *Holder) NextEvent() (i eventtypes.EventHandler) {
	if len(e.Queue) == 0 {
		return nil
	}

	i = e.Queue[0]
	e.Queue = e.Queue[1:]

	return i
}

func (tm *TradeManager) waitForDataCatchup() {
	var localWG sync.WaitGroup
	localWG.Add(1)
	// localWG.Add(1)

	db := tm.bot.DatabaseManager.GetInstance()
	dhj, err := datahistoryjob.Setup(db)
	if err != nil {
		fmt.Println("error", err)
	}

	dhj.ClearJobs()

	if tm.bot.dataHistoryManager.IsRunning() {
		// names, err := tm.bot.dataHistoryManager.CatchupDays(func() { localWG.Done() })
		tm.bot.dataHistoryManager.CatchupToday(func() { fmt.Println("DONE"); localWG.Done() })
	}

	log.Infoln(log.TradeMgr, "Waiting for data catchup...")
	localWG.Wait()

	for {
		// count jobs running
		active, err := dhj.CountActive()
		if err != nil {
			fmt.Println("error", err)
		}
		if active == 0 {
			break
		}
		fmt.Println("jobs still running")
		time.Sleep(time.Second)
	}

	// fmt.Println("validating...")
	// rTotal := make(map[*ExchangeAssetPairSettings]int)
	// for _, p := range tm.bot.CurrencySettings {
	// 	t := time.Now()
	// 	nowTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	// 	startDate := nowTime.AddDate(0, -1, 0)
	// 	rTotal[p] = 0
	//
	// 	candles, _ := candle.Series(p.ExchangeName, p.CurrencyPair.Base.String(), p.CurrencyPair.Quote.String(), 60, p.AssetType.String(), startDate, time.Now())
	// 	rTotal[p] += len(candles.Candles)
	// 	fmt.Println(p.CurrencyPair, "has", rTotal[p], "bars")
	// }
}

// ensure that we're synced before moving on
// r := make(map[*ExchangeAssetPairSettings]map[time.Time]int)
// for x := startDate; x.Before(nowTime); x = x.AddDate(0, 0, 1) {
// 	if r[p] == nil {
// 		r[p] = make(map[time.Time]int)
// 	}
// 	candles, _ := candle.Series(p.ExchangeName, p.CurrencyPair.Base.String(), p.CurrencyPair.Quote.String(), 60, p.AssetType.String(), time.Now().Add(time.Minute*-5), time.Now())
// 	r[p][x] = len(candles.Candles)
// 	rTotal[p] += len(candles.Candles)
// 	fmt.Println("adding", len(candles.Candles), rTotal[p])
// }

func (tm *TradeManager) waitForFactorEnginesWarmup() {
	// var localWG sync.WaitGroup
	// localWG.Add(1)
	tm.initializeFactorEngines()

	// load all candles for instrument

	for _, cs := range tm.bot.CurrencySettings {
		startDate := time.Now().Add(time.Minute * -120)
		// candles, _ := CandleSeriesForSettings(cs, 60, startDate, time.Now())
		dbData, err := database.LoadData(
			startDate,
			time.Now(),
			time.Minute,
			cs.ExchangeName,
			0,
			cs.CurrencyPair,
			cs.AssetType)
		if err != nil {
			fmt.Println("error load db data", err)
		}
		// fmt.Println(cs.CurrencyPair, "loaded", len(dbData.Item.Candles), "candles")
		tm.Datas.SetDataForCurrency(cs.ExchangeName, cs.AssetType, cs.CurrencyPair, dbData)
		dbData.Load()
	}

	tm.Run()
	tm.tradingEnabled = true

	// dbm := tm.bot.DatabaseManager.GetInstance()
	// if err != nil {
	// 	fmt.Println("error", err)
	// }
	//
	// for {
	// 	// count jobs running
	// 	active, err := db.CountActive()
	// 	if err != nil {
	// 		fmt.Println("error", err)
	// 	}
	// 	if active == 0 {
	// 		break
	// 	}
	// 	time.Sleep(time.Second)
	// }

	// localWG.Wait()
}

func (tm *TradeManager) runLive() error {
	processEventTicker := time.NewTicker(time.Second * 5)
	tm.waitForDataCatchup()
	tm.waitForFactorEnginesWarmup()
	log.Infoln(log.TradeMgr, "Running Live!")

	lup := make(map[*ExchangeAssetPairSettings]time.Time)

	var thisMinute, lastMinute time.Time
	loc, _ := time.LoadLocation("GMT")

	for {
		select {
		case <-tm.shutdown:
			return nil
		case <-processEventTicker.C:
			t := time.Now().UTC()
			thisMinute = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, loc)
			if thisMinute != lastMinute {
				lup = make(map[*ExchangeAssetPairSettings]time.Time)
				lastMinute = thisMinute
			}

			for _, cs := range tm.bot.CurrencySettings {
				if lup[cs] != thisMinute {
					dbData, err := database.LoadData(
						thisMinute.Add(time.Minute*-1),
						time.Now(),
						time.Minute,
						cs.ExchangeName,
						0,
						cs.CurrencyPair,
						cs.AssetType)

					if err != nil {
						// fmt.Println(err)
						continue
					}
					lastCandle := dbData.Item.Candles[len(dbData.Item.Candles)-1]
					t1 := lastCandle.Time
					t2 := thisMinute
					sameTime := (t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day() && t1.Hour() == t2.Hour() && t1.Minute() == t2.Minute())

					if !sameTime {
						// fmt.Println("don't have bar yet", lastCandle.Time, thisMinute)
						continue
					}
					dbData.Load()

					for dataEvent := dbData.Next(); ; dataEvent = dbData.Next() {
						if dataEvent == nil {
							dataEvent := dbData.Latest()
							tm.EventQueue.AppendEvent(dataEvent)
							lup[cs] = thisMinute
							// fmt.Println(cs.CurrencyPair, "seen", thisMinute, "loadedbars", t1, len(dbData.Item.Candles))
							// fmt.Println("sending event", dataEvent.GetTime())
							tm.Datas.SetDataForCurrency(cs.ExchangeName, cs.AssetType, cs.CurrencyPair, dbData)
							break
						}
					}

				}
			}

			err := tm.processEvents()
			if err != nil {
				log.Errorln(log.TradeMgr, "procesing events", err)
			}
		}
	}
	return nil
}

func (tm *TradeManager) handleEvent(ev eventtypes.EventHandler) error {
	switch eType := ev.(type) {
	case eventtypes.DataEventHandler:
		return tm.processSingleDataEvent(eType)
	case signal.Event:
		tm.processSignalEvent(eType)
	case order.Event:
		tm.processOrderEvent(eType)
	case submit.Event:
		tm.processSubmitEvent(eType)
	case fill.Event:
		tm.processFillEvent(eType)
	default:
		return fmt.Errorf("%w %v received, could not process",
			errUnhandledDatatype,
			ev)
	}

	return nil
}

func (tm *TradeManager) processSingleDataEvent(ev eventtypes.DataEventHandler) error {
	d := tm.Datas.GetDataForCurrency(ev.GetExchange(), ev.GetAssetType(), ev.Pair())
	cs, err := tm.bot.GetCurrencySettings(ev.GetExchange(), ev.GetAssetType(), ev.Pair())
	if cs == nil || err != nil {
		fmt.Println("error !!! FAIL getting cs", cs)
	}
	fe := tm.FactorEngines[cs.ExchangeName][cs.AssetType][cs.CurrencyPair]
	err = fe.OnBar(d)
	if err != nil {
		fmt.Printf("error updating factor engine for %v reason: %s", ev.Pair(), err)
	}

	if tm.tradingEnabled {
		tm.bot.OrderManager.Update()

		if len(fe.Minute().M60Range) > 0 {
			if tm.bot.Config.LiveMode {
				if tm.verbose {
					hrChg := fe.Minute().M60PctChange.Last(1).Round(2)

					if hrChg.GreaterThan(decimal.NewFromInt(0)) {
						color.Set(color.FgGreen, color.Bold)
					} else if hrChg.IsZero() {
						color.Set(color.FgWhite)
					} else if hrChg.LessThan(decimal.NewFromInt(0)) {
						color.Set(color.FgCyan, color.Bold)
					}
					defer color.Unset()

					log.Debugf(log.TradeMgr,
						"%s %s rRng:%v%% hrPctChg:%v%% close:%v hrRng:%v hrH: %v hrL: %v ",
						ev.GetTime().Format(common.SimpleTimeFormat),
						strings.ToUpper(ev.Pair().String()),
						fe.Minute().M60RangeDivClose.Last(1).Mul(decimal.NewFromInt(100)).Round(2),
						hrChg,
						fe.Minute().Close.Last(1),
						fe.Minute().M60Range.Last(1),
						fe.Minute().M60High.Last(1),
						fe.Minute().M60Low.Last(1))

				}
			}

			for _, strategy := range tm.Strategies {
				if strategy.GetPair() == ev.Pair() {
					s, err := strategy.OnData(d, tm.Portfolio, fe)
					if err != nil {
						fmt.Println("error processing data event", err)
						return err
					}
					tm.EventQueue.AppendEvent(s)
				}
			}
		} else {
			fmt.Println("only have", len(fe.Minute().M60Range), "close", len(fe.Minute().Close))
		}
	}

	return nil
}

func (tm *TradeManager) processEvents() error {
	for ev := tm.EventQueue.NextEvent(); ; ev = tm.EventQueue.NextEvent() {
		if ev != nil {
			err := tm.handleEvent(ev)
			if err != nil {
				return err
			}
		} else {
			return nil
		}

		if !tm.hasHandledEvent {
			tm.hasHandledEvent = true
		}
	}
}

func (tm *TradeManager) processSimultaneousDataEvents() error {
	return nil
}

func (tm *TradeManager) processSignalEvent(ev signal.Event) {
	cs, err := tm.bot.GetCurrencySettings(ev.GetExchange(), ev.GetAssetType(), ev.Pair())
	if err != nil {
		log.Error(log.TradeMgr, err)
		return
	}
	var o *order.Order
	o, err = tm.Portfolio.OnSignal(ev, cs)
	if err != nil {
		log.Error(log.TradeMgr, err)
		return
	}
	if err != nil {
		log.Error(log.TradeMgr, err)
	}

	if o != nil {
		err = tm.Statistic.SetEventForOffset(o)
		tm.EventQueue.AppendEvent(o)
	}
}

func (tm *TradeManager) onFill(o *OrderSubmitResponse) {
	if o.StrategyID == "" {
		fmt.Println("order submit response has no strategyID")
		os.Exit(2)
	}
	// return &OrderSubmitResponse{
	// 	SubmitResponse: order.SubmitResponse{
	// 		IsOrderPlaced: result.IsOrderPlaced,
	// 		OrderID:       result.OrderID,
	// 	},
	// 	InternalOrderID: id.String(),
	// 	StrategyID:      newOrder.StrategyID,
	// }, nil
	// convert to submit event

	ev := &fill.Fill{
		Base: event.Base{
			StrategyID: o.StrategyID,
		},
		OrderID: o.SubmitResponse.OrderID,
	}
	// fmt.Println("tmonfill creating fill event for:", ev.GetStrategyID())
	if ev.GetStrategyID() == "" {
		fmt.Println("noooooo")
		os.Exit(2)
	}
	tm.EventQueue.AppendEvent(ev)
}

func (tm *TradeManager) onCancel(o *OrderSubmitResponse) {
	// convert to submit event
	fmt.Println("onCancel", o)
	ev := &cancel.Cancel{}
	tm.EventQueue.AppendEvent(ev)
}

func (tm *TradeManager) processSubmitEvent(ev submit.Event) {
	// fmt.Println("processing submit event", ev.GetStrategyID())
	if ev.GetStrategyID() == "" {
		log.Error(log.TradeMgr, "submit event has no strategy ID")
		return
	}
	// convert order submit response to submit.Event here
	tm.Portfolio.OnSubmit(ev)
}

func (tm *TradeManager) processCancelEvent(ev cancel.Event) {
	tm.Portfolio.OnCancel(ev)
}

func (tm *TradeManager) processFillEvent(ev fill.Event) {
	tm.Portfolio.OnFill(ev)
}

func (tm *TradeManager) processOrderEvent(o order.Event) {
	if o.GetStrategyID() == "" {
		log.Error(log.TradeMgr, "order event has no strategy ID")
	}
	// else {
	// 	// gctlog.Debugln(log.TradeMgr, "creating order for", o.GetStrategyID())
	// }
	d := tm.Datas.GetDataForCurrency(o.GetExchange(), o.GetAssetType(), o.Pair())
	// this blocks and returns a submission event
	submitEvent, err := tm.ExecuteOrder(o, d, tm.bot.FakeOrderManager)

	// call on submit here

	if err != nil {
		log.Error(log.TradeMgr, err)
		return
	}

	if submitEvent.GetStrategyID() == "" {
		log.Error(log.TradeMgr, "Not strategy ID in order event")
		return
	}

	tm.EventQueue.AppendEvent(submitEvent)
}

func (tm *TradeManager) updateStatsForDataEvent(ev eventtypes.DataEventHandler) error {
	// update statistics with the latest price
	err := tm.Statistic.SetupEventForTime(ev)
	if err != nil {
		if err == statistics.ErrAlreadyProcessed {
			return err
		}
		log.Error(log.TradeMgr, err)
	}
	// update portfolio manager with the latest price
	err = tm.Portfolio.UpdateHoldings(ev)
	if err != nil {
		log.Error(log.TradeMgr, err)
	}
	return nil
}

func (tm *TradeManager) startOfflineServices() error {
	for _, cs := range tm.cfg.CurrencySettings {
		err := tm.bot.LoadExchange(cs.ExchangeName, nil)
		if err != nil && !errors.Is(err, ErrExchangeAlreadyLoaded) {
			return err
		}
	}

	tm.bot.setupExchangeSettings()

	// start fake order manager here since we don't start engine in live mode
	var err error
	tm.bot.FakeOrderManager, err = SetupFakeOrderManager(
		tm.bot,
		tm.bot.ExchangeManager,
		tm.bot.CommunicationsManager,
		&tm.bot.ServicesWG,
		tm.bot.Settings.Verbose,
	)
	if err != nil {
		gctlog.Errorf(gctlog.Global, "Fake Order manager unable to setup: %s", err)
	} else {
		err = tm.bot.FakeOrderManager.Start()
		if err != nil {
			gctlog.Errorf(gctlog.Global, "Fake Order manager unable to start: %s", err)
		}
		tm.bot.OrderManager = tm.bot.FakeOrderManager
	}

	tm.bot.DatabaseManager, err = SetupDatabaseConnectionManager(gctdatabase.DB.GetConfig())
	if err != nil {
		return err
	} else {
		err = tm.bot.DatabaseManager.Start(&tm.bot.ServicesWG)
		if err != nil {
			gctlog.Errorf(gctlog.Global, "Database manager unable to start: %v", err)
		}
	}
	return err
}

func (tm *TradeManager) initializeStrategies(cfg *config.Config) {
	var slit []strategies.Handler
	var s strategies.Handler
	count := 0
	for _, strat := range cfg.StrategiesSettings {
		for _, dir := range []gctorder.Side{gctorder.Buy, gctorder.Sell} {
			for _, c := range tm.bot.CurrencySettings {
				// fmt.Println("c", c)
				// _, pair, _, _ := tm.loadExchangePairAssetBase(c.ExchangeName, c.Base, c.Quote, c.Asset)
				s, _ = strategies.LoadStrategyByName(strat.Name)

				// tm.SetExchangeAssetCurrencySettings(exch, a, cp , c *ExchangeAssetPairSettings) {
				count += 1

				// fmt.Println("type of s", reflect.New(reflect.TypeOf(s)))
				// fmt.Println("type", reflect.New(reflect.ValueOf(s).Elem().Type()).Interface())
				// fmt.Println("valueof", reflect.New(reflect.ValueOf(s).Elem().Type()))
				// strategy = reflect.New(reflect.ValueOf(s).Elem().Type()).Interface().(strategy.Handler)
				// fmt.Println("loaded", strategy)

				id := fmt.Sprintf("%d_%s_%s_%v", count, s.Name(), string(gctorder.Buy), c.CurrencyPair)
				s.SetID(id)
				s.SetNumID(count)
				s.SetPair(c.CurrencyPair)
				s.SetDirection(dir)

				// validate strategy
				if s.GetID() == "" {
					fmt.Println("no strategy id")
					os.Exit(2)
				}
				s.SetDefaults()
				slit = append(slit, s)
			}
		}
	}
	tm.Strategies = slit
}

func (tm *TradeManager) initializePortfolio(strategyConfig *config.Config) error {
	// for i := range strategyConfig.Strategy
	tm.initializeStrategies(strategyConfig)

	if tm.verbose {
		log.Infof(log.TradeMgr, "Running %d strategies %d currencies", len(tm.Strategies), len(tm.bot.CurrencySettings))
	}

	// setup portfolio with strategies
	var p *Portfolio
	p, err := SetupPortfolio(tm.Strategies, tm.bot, strategyConfig)
	p.SetVerbose(false)
	if err != nil {
		return err
	}

	stats := &statistics.Statistic{
		StrategyName:                "ok",
		StrategyNickname:            strategyConfig.Nickname,
		StrategyDescription:         "ok",
		StrategyGoal:                strategyConfig.Goal,
		ExchangeAssetPairStatistics: make(map[string]map[asset.Item]map[currency.Pair]*currencystatistics.CurrencyStatistic),
		RiskFreeRate:                strategyConfig.StatisticSettings.RiskFreeRate,
	}
	tm.Statistic = stats

	// load from configuration into datastructure
	// currencysettings returns the data from the config, exchangeassetpairsettings
	// tm.Exchange = &e

	log.Infoln(log.TradeMgr, "Loaded", len(tm.bot.CurrencySettings), "currencies")

	tm.Portfolio = p
	return err
}

func (tm *TradeManager) initializeFactorEngines() error {
	tm.FactorEngines = make(map[string]map[asset.Item]map[currency.Pair]*FactorEngine)
	for _, cs := range tm.bot.CurrencySettings {
		if tm.FactorEngines[cs.ExchangeName] == nil {
			tm.FactorEngines[cs.ExchangeName] = make(map[asset.Item]map[currency.Pair]*FactorEngine)
		}
		if tm.FactorEngines[cs.ExchangeName][cs.AssetType] == nil {
			tm.FactorEngines[cs.ExchangeName][cs.AssetType] = make(map[currency.Pair]*FactorEngine)
		}
		fe, _ := SetupFactorEngine(cs, &tm.bot.Config.FactorEngine)
		tm.FactorEngines[cs.ExchangeName][cs.AssetType][cs.CurrencyPair] = fe
	}
	return nil
}
func (tm *TradeManager) setOrderManagerCallbacks() {
	// tm.bot.OrderManager.SetOnSubmit(tm.onSubmit)
	tm.OrderManager.SetOnFill(tm.onFill)
	tm.OrderManager.SetOnCancel(tm.onCancel)
}

// Series returns candle data
func CandleSeriesForSettings(e *ExchangeAssetPairSettings, interval int64, start, end time.Time) (out candle.Item, err error) {
	return candle.Series(e.ExchangeName, e.CurrencyPair.Base.String(), e.CurrencyPair.Quote.String(), 60, e.AssetType.String(), start, end)
}
