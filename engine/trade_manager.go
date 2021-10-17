package engine

// . "github.com/volatiletech/sqlboiler/v4/queries/qm"
import (
	"context"
	"errors"
	"fmt"
	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/data"
	"gocryptotrader/data/kline"
	"gocryptotrader/data/kline/database"
	"os"
	"path/filepath"
	"strings"
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
)

// Helper method for starting from live engine
// New returns a new TradeManager instance
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

// Reset TradeManager values to default
// this working will allow backtests to be run without shutting the engine on/off
func (tm *TradeManager) Reset() {
	tm.EventQueue.Reset()
	tm.Datas.Reset()
	tm.Portfolio.Reset()
	tm.Statistic.Reset()
	// tm.Exchange.Reset()
	// reset live trades here
	// tm.bot = nil
	// tm.FactorEngine = nil
}

// NewFromConfig takes a strategy config and configures a backtester variable to run
func NewTradeManagerFromConfig(cfg *config.Config, templatePath, output string, bot *Engine) (*TradeManager, error) {
	log.Debugln(log.TradeMgr, "TradeManager: Initializing... dry run", bot.Config.DryRun)
	if cfg == nil {
		return nil, errNilConfig
	}
	if bot == nil {
		return nil, errNilBot
	}
	tm := &TradeManager{
		shutdown: make(chan struct{}),
	}

	tm.cfg = *cfg
	tm.verbose = false
	tm.Warmup = bot.Config.LiveMode

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

	err = tm.setupBot(cfg)
	if err != nil {
		fmt.Println("error setting up bot")
		os.Exit(123)
	} else {
		fmt.Println("bot setup complete")
	}

	// initialize the data structure to hold the klines for each pair
	return tm, err
}

func (tm *TradeManager) setOrderManagerCallbacks() {
	// tm.bot.OrderManager.SetOnSubmit(tm.onSubmit)
	tm.OrderManager.SetOnFill(tm.onFill)
	tm.OrderManager.SetOnCancel(tm.onCancel)
}

// BACKTEST FUNCTIONALITY
// Run will iterate over loaded data events
// save them and then handle the event based on its type
func (tm *TradeManager) Run() error {
	tm.setOrderManagerCallbacks()
	log.Debugf(log.TradeMgr, "TradeManager Running. Warmup: %v\n", tm.Warmup)
dataLoadingIssue:
	for ev := tm.EventQueue.NextEvent(); ; ev = tm.EventQueue.NextEvent() {
		if ev == nil {
			dataHandlerMap := tm.Datas.GetAllData()
			for exchangeName, exchangeMap := range dataHandlerMap {
				for assetItem, assetMap := range exchangeMap {
					// var hasProcessedData bool
					for currencyPair, dataHandler := range assetMap {
						d := dataHandler.Next()
						if d == nil {
							fmt.Println("no data found for", currencyPair)
							if !tm.hasHandledEvent {
								log.Errorf(log.TradeMgr, "Unable to perform `Next` for %v %v %v", exchangeName, assetItem, currencyPair)
							}
							break dataLoadingIssue
						}
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
	tm.Warmup = false

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

// LIVE FUNCTIONALITY
func (tm *TradeManager) Start() error {
	tm.bot.WaitForInitialCurrencySync()
	tm.setOrderManagerCallbacks()

	// throw error if not live
	if !atomic.CompareAndSwapInt32(&tm.started, 0, 1) {
		return fmt.Errorf("backtester %w", ErrSubSystemAlreadyStarted)
	}

	// start trade manager
	log.Debugf(log.TradeMgr, "TradeManager  %s", MsgSubSystemStarting)
	tm.shutdown = make(chan struct{})

	// go tm.heartBeat()

	log.Debugln(log.TradeMgr, "Running Live")
	go tm.runLive()
	return nil
}

func (tm *TradeManager) runLive() error {
	processEventTicker := time.NewTicker(time.Second)
	for {
		select {
		case <-tm.shutdown:
			return nil
		case <-processEventTicker.C:
			for _, exchangeMap := range tm.Datas.GetAllData() { // for each exchange
				for _, assetMap := range exchangeMap { // asset
					for _, dataHandler := range assetMap { // coin
						d := dataHandler.Next()
						if d == nil {
							continue
						}
						tm.EventQueue.AppendEvent(d)
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

// Stop shuts down the live data loop
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
	// if tm.bot.DatabaseManager.IsRunning() {
	// 	tm.bot.DatabaseManager.Stop()
	// }

	for _, s := range tm.Strategies {
		s.Stop()
	}

	close(tm.shutdown)
	tm.bot.TradeManager = nil
	tm.wg.Wait()
	log.Debugln(log.TradeMgr, "Backtester Stopped.")
	return nil
}

// func (tm *TradeManager) loadOfflineDatas() error {
// 	fmt.Println("ok", len(tm.bot.CurrencySettings))
// 	// cfg := &tm.cfg
//
// 	tm.Datas.Setup()
//
// 	// LOAD DATA FOR EVERY PAIR
// 	for _, cs := range tm.bot.CurrencySettings {
// 		// exchangeName := strings.ToLower(exch.GetName())
// 		// klineData, err := tm.loadOfflineData(cfg, exch, pair, a)
//
// 		// err = resp.Load()
// 		// if err != nil {
// 		// 	fmt.Println("error", err)
// 		// 	continue
// 		// }
// 		// tm.Datas.SetDataForCurrency(cs.ExchangeName, cs.AssetType, cs.CurrencyPair, resp)
// 	}
// 	return nil
// }
//
// IsRunning returns if gctscript manager subsystem is started
func (b *TradeManager) IsRunning() bool {
	if b == nil {
		return false
	}
	return atomic.LoadInt32(&b.started) == 1
}

// this is only needed in backtest mode, except for when live runs the catchup process
// setupBot sets up a basic bot to retrieve exchange data
// as well as process orders
// setup order manager, exchange manager, database manager
func (tm *TradeManager) setupBot(strategyConfig *config.Config) error {
	var err error

	tm.Datas = &data.HandlerPerCurrency{}
	tm.Datas.Setup()

	if !tm.bot.Config.LiveMode {
		err = tm.startOfflineServices()
	}

	if err != nil {
		fmt.Println("failed to setup bot", err)
		return err
	}

	tm.initializeFactorEngines()

	if tm.bot.Settings.EnableTrading {
		fmt.Println("enable trading")
		tm.initializePortfolio(strategyConfig)
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
	// this is where we update the portfolio

	// if tm.Portfolio != nil {
	// 	tm.Portfolio.Update(ev)
	// }

	// update order manager
	// err := tm.updateStatsForDataEvent(ev)
	// if err != nil {
	// 	return err
	// }

	d := tm.Datas.GetDataForCurrency(ev.GetExchange(), ev.GetAssetType(), ev.Pair())
	cs, err := tm.bot.GetCurrencySettings(ev.GetExchange(), ev.GetAssetType(), ev.Pair())
	if cs == nil || err != nil {
		fmt.Println("error !!! FAIL getting cs", cs)
	}
	fe := tm.FactorEngines[cs]
	err = fe.OnBar(d)
	if err != nil {
		fmt.Printf("error updating factor engine for %v reason: %s", ev.Pair(), err)
	}

	// HANDLE warmup MODE
	// in warmup mode, we do not query the strategies
	if !tm.Warmup {
		tm.bot.OrderManager.Update()

		for _, strategy := range tm.Strategies {
			if strategy.GetPair() == ev.Pair() {
				if tm.bot.Config.LiveMode {
					fmt.Println("Updating strategy", strategy.GetID(), d.Latest().GetTime())
				}
				s, err := strategy.OnData(d, tm.Portfolio, fe)
				if err != nil {
					fmt.Println("error processing data event", err)
					return err
				}
				tm.EventQueue.AppendEvent(s)
			}
		}
	}

	// if err != nil {
	// 	if errors.Is(err, base.ErrTooMuchBadData) {
	// 		// too much bad data is a severe error and backtesting must cease
	// 		return err
	// 	}
	// 	log.Error(log.TradeMgr, err)
	// 	return nil
	// }
	// err = tm.Statistic.SetEventForOffset(ev)
	// if err != nil {
	// 	log.Error(log.TradeMgr, err)
	// }

	return nil
}

// processSimultaneousDataEvents determines what signal events are generated and appended
// to the event queue based on whether it is running a multi-currency consideration strategy order not
//
// for multi-currency-consideration it will pass all currency datas to the strategy for it to determine what
// currencies to act upon
//
// for non-multi-currency-consideration strategies, it will simply process every currency individually
// against the strategy and generate signals
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

// creates a fill event based on an order submit response
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

// new orders
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

// updateStatsForDataEvent makes various systems aware of price movements from
// data events
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

func loadDatabaseData(cfg *config.Config, name string, fPair currency.Pair, a asset.Item, dataType int64) (*kline.DataFromKline, error) {
	if cfg == nil || cfg.DataSettings.DatabaseData == nil {
		return nil, errors.New("nil config data received")
	}
	if cfg.DataSettings.Interval <= 0 {
		return nil, errIntervalUnset
	}

	return database.LoadData(
		time.Now(),
		time.Now().Add(-1*time.Minute),
		cfg.DataSettings.Interval,
		strings.ToLower(name),
		dataType,
		fPair,
		a)
}

// Reset returns struct to defaults
func (e *Holder) Reset() {
	e.Queue = nil
}

// AppendEvent adds and event to the queue
func (e *Holder) AppendEvent(i eventtypes.EventHandler) {
	e.Queue = append(e.Queue, i)
}

// NextEvent removes the current event and returns the next event in the queue
func (e *Holder) NextEvent() (i eventtypes.EventHandler) {
	if len(e.Queue) == 0 {
		return nil
	}

	i = e.Queue[0]
	e.Queue = e.Queue[1:]

	return i
}

// ExecuteOrder assesses the portfolio manager's order event and if it passes validation
// will send an order to the exchange/fake order manager to be stored and raise a fill event
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

func (tm *TradeManager) startOfflineServices() error {
	for i := range tm.bot.Config.CurrencySettings {
		err := tm.bot.LoadExchange(tm.bot.Config.CurrencySettings[i].ExchangeName, nil)
		if err != nil && !errors.Is(err, ErrExchangeAlreadyLoaded) {
			return err
		}
	}

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
		// for _, dir := range []gctorder.Side{gctorder.Buy, gctorder.Sell} {
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
			s.SetDirection(gctorder.Buy)

			// validate strategy
			if s.GetID() == "" {
				fmt.Println("no strategy id")
				os.Exit(2)
			}
			s.SetDefaults()
			slit = append(slit, s)
		}
		// }
	}
	tm.Strategies = slit
}

func (tm *TradeManager) initializePortfolio(strategyConfig *config.Config) error {
	// for i := range strategyConfig.Stratey
	tm.initializeStrategies(strategyConfig)

	// if tm.verbose {
	log.Infof(log.TradeMgr, "Running %d strategies\n", len(tm.Strategies))
	log.Infof(log.TradeMgr, "Running %d currencies\n", len(tm.bot.CurrencySettings))
	// log.Infof(log.TradeMgr, "Running %d exchanges\n", len(tm.bot.CurrencySettings))
	// }

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
	tm.FactorEngines = make(map[*ExchangeAssetPairSettings]*FactorEngine)
	for _, cs := range tm.bot.CurrencySettings {
		log.Debugln(log.TradeMgr, "creating fctor engine", cs.ExchangeName, cs.CurrencyPair)
		fe, _ := SetupFactorEngine(cs)
		tm.FactorEngines[cs] = fe
	}
	return nil
}
