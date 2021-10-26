package engine

// . "github.com/volatiletech/sqlboiler/v4/queries/qm"
import (
	"context"
	"fmt"
	"gocryptotrader/common"
	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/data"
	datakline "gocryptotrader/data/kline"
	"gocryptotrader/data/kline/database"
	"gocryptotrader/database/repository/candle"
	"gocryptotrader/database/repository/datahistoryjob"
	"gocryptotrader/database/repository/liveorder"
	"gocryptotrader/database/repository/livetrade"
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
	"gocryptotrader/exchange/kline"

	gctorder "gocryptotrader/exchange/order"
	"gocryptotrader/log"

	gctlog "gocryptotrader/log"
	"gocryptotrader/portfolio/compliance"
	"gocryptotrader/portfolio/holdings"
	"gocryptotrader/portfolio/report"
	"gocryptotrader/portfolio/statistics"
	"gocryptotrader/portfolio/statistics/currencystatistics"
	"gocryptotrader/portfolio/tradereport"

	"github.com/fatih/color"
	"github.com/shopspring/decimal"
)

func NewTradeManager(bot *Engine) (*TradeManager, error) {
	configPath := bot.Settings.TradeConfigFile
	wd, err := os.Getwd()
	if configPath == "" {
		if bot.Config.LiveMode {
			configPath = filepath.Join(wd, "cmd/confs/prod.strat")
		} else {
			configPath = filepath.Join(wd, "cmd/confs/dev/backtest.strat")
		}
	} else {
		configPath = filepath.Join(wd, "cmd/confs/dev/strategy", fmt.Sprintf("%s.strat", configPath))
	}
	btcfg, err := config.ReadConfigFromFile(configPath)
	if err != nil {
		fmt.Println("error", configPath, err)
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
	log.Debugln(log.TradeMgr, "TradeManager: Initializing...")

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
	tm.verbose = bot.Config.TradeManager.Verbose
	tm.tradingEnabled = bot.Settings.EnableTrading
	tm.dryRun = bot.Settings.EnableDryRun
	tm.liveSimulationCfg = bot.Config.TradeManager.LiveSimulation
	tm.isSimulation = tm.liveSimulationCfg.Enabled
	tm.liveMode = bot.Config.LiveMode
	tm.debug = bot.Config.TradeManager.Debug

	if tm.isSimulation {
		tm.currentTime = tm.liveSimulationCfg.StartDate
	}
	// fmt.Println("tmconfig", cfg.TradeManager, cfg.TradeManager.Verbose, cfg.TradeManager.Trading, cfg.TradeManager.Enabled)

	stats := &statistics.Statistic{
		StrategyName:                "ok",
		StrategyNickname:            cfg.Nickname,
		StrategyDescription:         "ok",
		StrategyGoal:                cfg.Goal,
		ExchangeAssetPairStatistics: make(map[string]map[asset.Item]map[currency.Pair]*currencystatistics.CurrencyStatistic),
		RiskFreeRate:                cfg.StatisticSettings.RiskFreeRate,
	}
	tm.Statistic = stats

	tm.EventQueue = &Holder{}
	reports := &report.Data{
		Config:       cfg,
		TemplatePath: templatePath,
		OutputPath:   output,
	}
	tm.Reports = reports

	tradereports := &tradereport.Data{
		Config:       cfg,
		TemplatePath: templatePath,
		OutputPath:   output,
	}
	tm.TradeReports = tradereports
	reports.Statistics = stats
	tradereports.Statistics = stats

	tm.bot = bot
	var err error
	if bot.OrderManager == nil && bot.Settings.EnableOrderManager {
		// log.Warnln(log.TradeMgr, "!!!!!!!!!!!!!!!!!!!!!!!! Enabling REAL $$$$$$$$ order manager!!!!!!!!!!!!!!!!!!!!")
		bot.OrderManager, err = SetupOrderManager(
			bot.ExchangeManager,
			bot.CommunicationsManager,
			&bot.ServicesWG,
			bot.Config.OrderManager.Verbose,
			bot.Config.ProductionMode,
			bot.Config.LiveMode,
			bot.Config.DryRun,
		)

		if err != nil {
			log.Errorf(log.Global, "Fake Order manager unable to setup: %s", err)
		} else {
			err = bot.OrderManager.Start()

			if err != nil {
				log.Errorf(log.Global, "Fake Order manager unable to start: %s", err)
			}
		}
	}
	tm.syncManager = bot.currencyPairSyncer

	tm.Datas = &data.HandlerPerCurrency{}
	tm.Datas.Setup()
	if !tm.liveMode {
		// log.Debug(log.TradeMgr, "starting offline services")
		err = tm.startOfflineServices()
	}
	if err != nil {
		fmt.Println("failed to setup bot", err)
		// return err
	}

	if tm.tradingEnabled {
		tm.Strategies = SetupStrategies(cfg, tm.liveMode)

		if tm.bot.Settings.EnableClearDB {
			log.Warn(log.TradeMgr, "clearing DB")
			if tm.bot.Config.ProductionMode || tm.bot.Config.Database.ConnectionDetails.Database == "gct_prod" {
				// check database name to ensure we don't delete anything
				panic("trying to delete production")
			}
			err := liveorder.DeleteAll()
			err = livetrade.DeleteAll()
			if err != nil {
				fmt.Println("did not delete", err)
				os.Exit(123)
			}
		}

		p, err := SetupPortfolio(tm.Strategies, tm.bot, tm.bot.Config)
		if err != nil {
			return nil, fmt.Errorf("could not setup portfolio", err)
		}
		tm.Portfolio = p
	}

	if err != nil {
		fmt.Println("error setting up tm", err)
		os.Exit(123)
	}

	// fmt.Println("done setting up bot with", len(tm.bot.CurrencySettings), "currencies")
	// if len(tm.bot.CurrencySettings) < 1 {
	// 	log.Error(log.TradeMgr, "!!no currency settings")
	// 	os.Exit(123)
	// }

	return tm, err
}

func (tm *TradeManager) ExecuteOrder(o order.Event, data data.Handler, om ExecutionHandler) (submit.Event, error) {
	if tm.debug {
		log.Debugf(log.TradeMgr, "Executing order", o.GetDecision())
	}
	priceFloat, _ := o.GetPrice().Float64()
	a, _ := o.GetAmount().Float64()
	fee, _ := o.GetExchangeFee().Float64()

	var skipStop bool
	if o.GetDecision() == signal.Exit {
		skipStop = true
	} else if o.GetDecision() == "" {
		panic("order without decision")
	} else if o.GetDecision() == signal.Exit {
		fmt.Println("cancel the stop loss order too ey")
	} else {
		fmt.Println("decision", o.GetDecision())
	}
	if o.GetPrice().IsZero() {
		panic("order has no price")
	}

	stopLossPrice, _ := o.GetStopLossPrice().Float64()

	submission := &gctorder.Submit{
		Status:        gctorder.New,
		Price:         priceFloat,
		Amount:        a,
		Fee:           fee,
		Exchange:      o.GetExchange(),
		ID:            o.GetID(),
		Side:          o.GetDirection(),
		AssetType:     o.GetAssetType(),
		Date:          o.GetTime(),
		LastUpdated:   o.GetTime(),
		Pair:          o.Pair(),
		Type:          gctorder.Market,
		StrategyID:    o.GetStrategyID(),
		StopLossPrice: stopLossPrice,
	}

	if om == nil {
		panic("there is no order manager and trying to execute order")
	}

	var stopSide gctorder.Side
	if o.GetDirection() == gctorder.Buy {
		stopSide = gctorder.Sell
	} else if o.GetDirection() == gctorder.Sell {
		stopSide = gctorder.Buy
	}
	stopLossSubmission := &gctorder.Submit{
		Status:      gctorder.New,
		Price:       stopLossPrice,
		Amount:      a,
		Fee:         fee,
		Exchange:    o.GetExchange(),
		ID:          o.GetID(),
		Side:        stopSide,
		AssetType:   o.GetAssetType(),
		Date:        o.GetTime(),
		LastUpdated: o.GetTime(),
		Pair:        o.Pair(),
		Type:        gctorder.Stop,
		StrategyID:  o.GetStrategyID(),
	}

	var entryID, stopID int
	var err error
	if !tm.dryRun {
		entryID, err = liveorder.Insert(submission)
		if err != nil {
			fmt.Println("error inserted order", err)
			return nil, err
		}

		if !skipStop {
			stopID, err = liveorder.Insert(stopLossSubmission)
			if err != nil {
				fmt.Println("error inserted order", err)
				return nil, err
			}
		}
	} else {
		entryID = om.GenerateDryRunID()
		stopID = entryID + 1
	}

	// ords, _ := om.GetOrdersSnapshot("")
	// var internalOrderID int
	// for i := range ords {
	// 	fmt.Println("checking order id", ords[i].InternalOrderID, o.GetID())
	// 	if ords[i].ID != omr.InternalOrderID {
	// 		continue
	// 	}
	// 	ords[i].StrategyID = o.GetStrategyID()
	// 	ords[i].Date = o.GetTime()
	// 	ords[i].LastUpdated = o.GetTime()
	// 	ords[i].CloseTime = o.GetTime()
	// }

	submission.InternalOrderID = entryID
	stopLossSubmission.InternalOrderID = stopID

	omr, err := om.Submit(context.TODO(), submission)
	if err != nil {
		fmt.Println("tm: ERROR order manager submission", err, submission.Side, omr)
	}
	if omr.InternalOrderID == 0 {
		panic("no order id")
	}

	var stopLossOrderID int
	if !skipStop {
		somr, err := om.Submit(context.TODO(), stopLossSubmission)
		if err != nil {
			fmt.Println("tm: ERROR order manager submission", err, stopLossSubmission.Side, omr)
		}
		if somr.InternalOrderID == 0 {
			panic("no order id")
		}
		stopLossOrderID = somr.InternalOrderID
	}

	// fmt.Println("tm: order manager response", omr)
	// if order is placed, update the status of the order to Open
	// update order event order_id, status
	// add the submission to the store
	// we can call on submit manually here
	// if o.GetStrategyID() == "" {
	// 	return nil, fmt.Errorf("exchange: order has no strategyid")
	// }

	if omr.IsOrderPlaced && omr.Rate == 0 {
		panic("order placed without price/rate")
	}

	ev := &submit.Submit{
		Base: event.Base{
			Offset:       o.GetOffset(),
			Exchange:     o.GetExchange(),
			Time:         o.GetTime(),
			CurrencyPair: o.Pair(),
			AssetType:    o.GetAssetType(),
			Interval:     o.GetInterval(),
			Reason:       o.GetReason(),
			StrategyID:   o.GetStrategyID(),
		},
		InternalOrderID: omr.InternalOrderID,
		StopLossOrderID: stopLossOrderID,
		IsOrderPlaced:   omr.IsOrderPlaced,
		OrderID:         omr.OrderID,
		FullyMatched:    omr.FullyMatched,
		Price:           omr.Rate,
		StrategyID:      o.GetStrategyID(),
	}

	// if ev.GetInternalOrderID() == "" {
	// 	log.Errorln(log.TradeMgr, "error: order has no internal order id")
	// }

	return ev, nil
}

func (tm *TradeManager) Run() error {
	count := 0
	log.Debugf(log.TradeMgr, "TradeManager Running")
	if !tm.liveMode {
		err := tm.loadBacktestData()
		if err != nil {
			fmt.Println("error loadBacktestData:", err)
		}
		t1 := tm.cfg.DataSettings.DatabaseData.StartDate
		t2 := tm.cfg.DataSettings.DatabaseData.EndDate
		dayDuration := int(t2.Sub(t1).Minutes()) / 60 / 24
		if dayDuration > 30 {
			panic("more than 30 days")
		}
		log.Warnln(log.TradeMgr, "startdate:", t1)
		log.Warnln(log.TradeMgr, "enddate:", t2)
		log.Warnln(log.TradeMgr, "duration:", dayDuration, "days")
		log.Warnln(log.TradeMgr, "strategies:", len(tm.Strategies))
		for _, s := range tm.Strategies {
			log.Debugln(log.TradeMgr, s.GetPair(), s.Name(), s.GetDirection(), s.GetID())
		}
		// pairs, err := tm.bot.Config.GetEnabledPairs("gateio", asset.Spot)
		// for _, p := range pairs {
		// 	log.Debugln(log.TradeMgr, "Active Pair:", p)
		// }
	}
	// return nil
dataLoadingIssue:
	for ev := tm.EventQueue.NextEvent(); ; ev = tm.EventQueue.NextEvent() {
		// check for new day
		if ev == nil {
			// if !tm.liveMode && count%1000 == 0 {
			// 	fmt.Printf(".")
			// }

			dataHandlerMap := tm.Datas.GetAllData()
			for _, exchangeMap := range dataHandlerMap {
				for _, assetMap := range exchangeMap {
					// var hasProcessedData bool
					tm.hasHandledEvent = false
					for _, dataHandler := range assetMap {
						d := dataHandler.Next()
						if d == nil {
							// if !tm.hasHandledEvent {
							// 	log.Errorf(log.TradeMgr, "Unable to perform `Next` for %v %v %v", exchangeName, assetItem, currencyPair)
							// }
							break dataLoadingIssue
						}
						tm.hasHandledEvent = true
						count += 1
						// fmt.Println("data event", d)

						if !tm.bot.Config.ProductionMode {
							tm.bot.OrderManager.UpdateFakeOrders(d)
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

	if tm.liveMode {
		fmt.Println("done running", count, "data events")
	} else {
		livetrade.WriteCSV(tm.Portfolio.GetAllClosedTrades())
		livetrade.AnalyzeTrades("")
		log.Debugln(log.TradeMgr, "TradeManager Writing Config to File")
		tm.cfg.SaveConfigToFile("backtest_config_out.json")
	}

	return nil
}

func (tm *TradeManager) Start() error {
	if !atomic.CompareAndSwapInt32(&tm.started, 0, 1) {
		return fmt.Errorf("backtester %w", ErrSubSystemAlreadyStarted)
	}
	tm.shutdown = make(chan struct{})

	if tm.liveMode {
		printStrategies(tm.Strategies)
	}

	// go tm.heartBeat()
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

	log.Debugln(log.TradeMgr, "TradeManager Stopping...")

	if tm.bot.OrderManager != nil && tm.bot.OrderManager.IsRunning() {
		tm.bot.OrderManager.Stop()
	}
	for _, s := range tm.Strategies {
		s.Stop()
	}
	close(tm.shutdown)
	tm.bot.TradeManager = nil
	tm.wg.Wait()
	log.Debugln(log.TradeMgr, "TradeManager Stopped.")
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
	// fmt.Println(0)
	// var localWG sync.WaitGroup
	// localWG.Add(1)

	db := tm.bot.DatabaseManager.GetInstance()
	dhj, err := datahistoryjob.Setup(db)
	if err != nil {
		fmt.Println("error", err)
	}

	dhj.ClearJobs()

	log.Infoln(log.TradeMgr, "Catching up days...", tm.bot.dataHistoryManager.DaysBack)
	daysBack := make([]int, tm.bot.dataHistoryManager.DaysBack)

	for i := range daysBack {
		i += 1
		tm.bot.dataHistoryManager.CatchupDays(int64(i))

		for {
			active, err := dhj.CountActive()
			if err != nil {
				fmt.Println("error", err)
			}
			if active == 0 {
				fmt.Println("starting days back", i)
				break
			}
			time.Sleep(time.Second)
		}
	}

	// time.Sleep(time.Millisecond * 500)
	log.Infoln(log.TradeMgr, "Done with catchup")
	os.Exit(1123)
}

func (tm *TradeManager) waitForFactorEnginesWarmup() {
	// fmt.Println("warm up factor engines")
	tm.initializeFactorEngines()

	// load all candles for instrument
	for _, cs := range tm.bot.CurrencySettings {
		startDate := tm.GetCurrentTime().Add(time.Minute * -120)
		// candles, _ := CandleSeriesForSettings(cs, 60, startDate, tm.GetCurrentTime())
		dbData, err := database.LoadData(
			startDate,
			tm.GetCurrentTime(),
			time.Minute,
			cs.ExchangeName,
			0,
			cs.CurrencyPair,
			cs.AssetType)

		if err != nil {
			fmt.Errorf("error load db data", err)
		}

		if dbData != nil {
			tm.Datas.SetDataForCurrency(cs.ExchangeName, cs.AssetType, cs.CurrencyPair, dbData)
			dbData.Load()
		}
		// fmt.Println("dbData", dbData)
		// fmt.Println("dbDataitem", dbData.Item)
		// fmt.Println(cs.CurrencyPair, "loaded", len(dbData.Item.Candles), "candles")
	}

	tm.tradingEnabled = false
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
	log.Debugln(log.TradeMgr, "Waiting for initial currency sync...")
	tm.bot.WaitForInitialCurrencySync()
	log.Debugln(log.TradeMgr, "Finished Initial Currency Sync")

	// var processEventTicker time.Ticker
	// processEventTickerSim := time.NewTicker(time.Second)
	processEventTicker := time.NewTicker(time.Second * 5)
	if tm.bot.dataHistoryManager.IsRunning() {
		tm.waitForDataCatchup()
	}
	tm.waitForFactorEnginesWarmup()
	log.Infoln(log.TradeMgr, "Running Live!")

	tm.lastUpdateMin = make(map[*ExchangeAssetPairSettings]time.Time)

	for {
		select {
		case <-tm.shutdown:
			return nil
		case <-processEventTicker.C:
			err := tm.processLiveMinute()

			if err != nil {
				fmt.Println("error live min process", err)
			}
		}
	}
	return nil
}

func (tm *TradeManager) processLiveMinute() error {
	var thisMinute, lastMinute time.Time
	loc, _ := time.LoadLocation("UTC")
	t := tm.GetCurrentTime()
	thisMinute = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, loc)
	if !lastMinute.IsZero() && thisMinute != lastMinute {
		lastMinute = thisMinute
	}

	for _, cs := range tm.bot.CurrencySettings {
		if tm.lastUpdateMin[cs] != thisMinute {
			dbData, err := tm.loadLatestCandleFromDatabase(cs)
			if err != nil {
				fmt.Println("error loading latest candle", err)
				continue
			}
			fmt.Println("loaded latest candle")

			dataEvent := dbData.Next()
			for ; ; dataEvent = dbData.Next() {
				if dataEvent == nil {
					break
				}
			}
			dataEvent = dbData.Latest()

			if !common.IsSameMinute(thisMinute, dataEvent.GetTime()) {
				fmt.Println("skipping already seen bar", dataEvent.GetTime(), thisMinute)
				continue
			}

			if !dbData.HasDataAtTime(dataEvent.GetTime()) {
				log.Error(log.TradeMgr, "doesnt have data in range")
				os.Exit(123)
			}
			fmt.Println("updating with event", dataEvent)
			tm.lastUpdateMin[cs] = dataEvent.GetTime().UTC()
			tm.EventQueue.AppendEvent(dataEvent)
		}
	}

	err := tm.processEvents()
	if err != nil {
		log.Errorln(log.TradeMgr, "procesing events", err)
		return err
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
	// if tm.tradingEnabled {
	// fmt.Println("tm processing event at", ev.GetTime(), ev.Pair(), tm.GetCurrentTime())
	// }
	err := tm.updateStatsForDataEvent(ev)
	if err != nil {
		fmt.Println("error updating stats for data event")
		return err
	}

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

	// handle old events somehow
	// we need to ignore them to not trick the order manager or factor engine
	// while ensuring the factor engine has all the historical data
	// old events should not come through here

	if tm.tradingEnabled {
		if tm.bot.OrderManager != nil {
			tm.bot.OrderManager.Update()
		}

		if len(fe.Minute().M60Range) > 0 {
			if tm.liveMode {

				if tm.verbose {
					hrChg := fe.Minute().M60PctChange.Last(1).Round(2)

					if hrChg.GreaterThan(decimal.NewFromInt(0)) {
						color.Set(color.FgGreen, color.Bold)
					} else if hrChg.IsZero() {
						color.Set(color.FgWhite)
					} else if hrChg.LessThan(decimal.NewFromInt(0)) {
						color.Set(color.FgYellow, color.Bold)
					}
					defer color.Unset()

					log.Debugf(log.TradeMgr,
						"%2d:%2d %-12s %12v %7v%% %7v%% %12v %12v %12v",
						ev.GetTime().Hour(),
						ev.GetTime().Minute(),
						strings.ToUpper(ev.Pair().String()),
						fe.Minute().Close.Last(1),
						fe.Minute().M60RangeDivClose.Last(1).Mul(decimal.NewFromInt(100)).Round(2),
						hrChg,
						fe.Minute().M60Range.Last(1),
						fe.Minute().M60High.Last(1),
						fe.Minute().M60Low.Last(1))

				}

			}

			for _, strategy := range tm.Strategies {
				sp := strategy.GetPair()
				ep := ev.Pair()
				if arePairsEqual(sp, ep) {
					s, err := strategy.OnData(d, tm.Portfolio, fe)
					s.SetStrategyID(strategy.GetID())
					if err != nil {
						fmt.Println("error processing data event", err)
						return err
					}

					err = tm.Statistic.SetEventForOffset(s)
					if err != nil {
						log.Error(log.TradeMgr, err)
					}
					tm.EventQueue.AppendEvent(s)
				}
			}
		} else {
			if tm.bot.Settings.EnableLiveMode {
				fmt.Println("only have last", len(fe.Minute().M60Range), "m60 range bars, closes' length ", len(fe.Minute().Close))
			}
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
	// fmt.Println("process signal", ev.GetReason())
	cs, err := tm.bot.GetCurrencySettings(ev.GetExchange(), ev.GetAssetType(), ev.Pair())
	if err != nil {
		log.Error(log.TradeMgr, "error", err)
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
		if tm.debug {
			fmt.Println("tm received order from pf", o, err)
		}
		err = tm.Statistic.SetEventForOffset(o)
		tm.EventQueue.AppendEvent(o)
	}
}

func (tm *TradeManager) createFillEvent(ev submit.Event) {
	if ev.GetStrategyID() == 0 {
		panic("order submit response has no strategyID")
	}
	if !ev.GetIsOrderPlaced() {
		panic("trying filling an unsubmitted order")
	}
	if ev.GetTime().IsZero() {
		panic("event has no time")
	}
	if ev.GetInternalOrderID() == 0 {
		panic("event has no internal order id")
	}

	// fmt.Println("decision", ev.GetDecision(), "evdir", ev.GetDirection())

	// fmt.Println("CREATE FILL EVENT REFERENCING ORDER", ev.GetOrderID(), "for strategy", ev.GetStrategyID())
	// fmt.Println("returned order id:", o.ID, "internal", o.InternalOrderID)

	// if err != nil {
	// 	panic("error getting order from store")
	// }

	// validate the side here

	o := tm.Portfolio.GetOrderFromStore(ev.GetInternalOrderID())
	stopLossID := ev.GetStopLossOrderID()
	var stopLossPrice float64
	if stopLossID != 0 {
		so := tm.Portfolio.GetOrderFromStore(ev.GetStopLossOrderID())
		stopLossPrice = so.Price
		if stopLossPrice == 0 {
			panic("has 0 stop loss")
		}
		// else {
		// 	fmt.Println("retrieved secondary stop loss order ID:", ev.GetStopLossOrderID(), "at price", stopLossPrice)
		// }
	}

	if o.Amount == 0 {
		panic("order amount is 0")
	}
	if o.Price == 0 {
		panic("order price is 0 and filled")
	}

	e := &fill.Fill{
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
		OrderID:         ev.GetOrderID(),
		InternalOrderID: ev.GetInternalOrderID(),
		ClosePrice:      decimal.NewFromFloat(o.Price),
		PurchasePrice:   decimal.NewFromFloat(o.Price),
		StopLossPrice:   decimal.NewFromFloat(stopLossPrice),
		Direction:       o.Side,
		Amount:          decimal.NewFromFloat(o.Amount),
		StopLossOrderID: stopLossID,
		// Direction:  ev.GetDirection(),
		// Amount:     ev.GetAmount(),
	}
	tm.EventQueue.AppendEvent(e)
	// if o.InternalOrderID == "" {
	// 	fmt.Println("order submit response has internal order id")
	// 	os.Exit(2)
	// }
	// if o.SubmitResponse.OrderID == "" {
	// 	fmt.Println("order submit response has no order id")
	// 	os.Exit(2)
	// }
	// return &OrderSubmitResponse{
	// 	SubmitResponse: order.SubmitResponse{
	// 		IsOrderPlaced: result.IsOrderPlaced,
	// 		OrderID:       result.OrderID,
	// 	},
	// 	InternalOrderID: id.String(),
	// 	StrategyID:      newOrder.StrategyID,
	// }, nil
	// convert to submit event

}

func (tm *TradeManager) onCancel(o *OrderSubmitResponse) {
	// convert to submit event
	fmt.Println("onCancel", o)
	ev := &cancel.Cancel{}
	tm.EventQueue.AppendEvent(ev)
}

func (tm *TradeManager) processSubmitEvent(ev submit.Event) {
	if tm.debug {
		fmt.Println("processing submit event strategy:", ev.GetStrategyID())
	}
	if ev.GetStrategyID() == 0 {
		log.Error(log.TradeMgr, "submit event has no strategy ID")
		return
	}
	if ev.GetOrderID() == "" {
		log.Error(log.TradeMgr, "submit event has no order ID")
		panic("no order id")
		return
	}

	if ev.GetIsOrderPlaced() {
		// fmt.Println("creating fill", ev.GetStrategyID(), "orderid", ev.GetOrderID(), "internalID", ev.GetInternalOrderID())
		tm.createFillEvent(ev)
	}
}

func (tm *TradeManager) processCancelEvent(ev cancel.Event) {
	tm.Portfolio.OnCancel(ev)
}

func (tm *TradeManager) processFillEvent(ev fill.Event) {
	tm.Portfolio.OnFill(ev)
	// do it like this
	// t, err := bt.Portfolio.OnFill(ev, funds)
	// if err != nil {
	// 	log.Error(log.BackTester, err)
	// 	return
	// }

	err := tm.Statistic.SetEventForOffset(ev)
	if err != nil {
		log.Error(log.TradeMgr, err)
	}

	var holding *holdings.Holding
	holding, err = tm.Portfolio.ViewHoldingAtTimePeriod(ev)
	if err != nil {
		log.Error(log.TradeMgr, err)
	}

	err = tm.Statistic.AddHoldingsForTime(holding)
	if err != nil {
		log.Error(log.TradeMgr, err)
	}

	var cp *compliance.Manager
	cp, err = tm.Portfolio.GetComplianceManager(ev.GetExchange(), ev.GetAssetType(), ev.Pair())
	if err != nil {
		log.Error(log.TradeMgr, err)
	}

	snap := cp.GetLatestSnapshot()
	err = tm.Statistic.AddComplianceSnapshotForTime(snap, ev)
	if err != nil {
		log.Error(log.TradeMgr, err)
	}
}

func (tm *TradeManager) processOrderEvent(o order.Event) {
	if o.GetStrategyID() == 0 {
		log.Error(log.TradeMgr, "order event has no strategy ID")
	}
	// else {
	// 	// gctlog.Debugln(log.TradeMgr, "creating order for", o.GetStrategyID())
	// }
	d := tm.Datas.GetDataForCurrency(o.GetExchange(), o.GetAssetType(), o.Pair())

	// this blocks and returns a submission event
	submitEvent, err := tm.ExecuteOrder(o, d, tm.bot.OrderManager)

	if err != nil {
		log.Error(log.TradeMgr, err)
		return
	}

	if submitEvent.GetStrategyID() == 0 {
		log.Error(log.TradeMgr, "Not strategy ID in order event")
		return
	}

	tm.EventQueue.AppendEvent(submitEvent)
}

func (tm *TradeManager) updateStatsForDataEvent(ev eventtypes.DataEventHandler) error {
	// fmt.Println("update stats", ev.GetTime())
	// update statistics with the latest price
	err := tm.Statistic.SetupEventForTime(ev)
	if err != nil {
		if err == statistics.ErrAlreadyProcessed {
			return err
		}
		log.Error(log.TradeMgr, err)
	}
	// update portfolio manager with the latest price
	// fmt.Println("portfolio", tm.Portfolio)
	err = tm.Portfolio.UpdateHoldings(ev)
	if err != nil {
		log.Error(log.TradeMgr, err)
	}
	return nil
}

func (tm *TradeManager) startOfflineServices() error {
	// fmt.Println("TM start offline services")
	// for _, cs := range tm.cfg.CurrencySettings {
	// 	err := tm.bot.LoadExchange(cs.ExchangeName, nil)
	// 	if err != nil && !errors.Is(err, ErrExchangeAlreadyLoaded) {
	// 		return err
	// 	}
	// }

	err := tm.bot.SetupExchanges()
	if err != nil {
		return err
	}
	tm.bot.SetupExchangeSettings()

	// start fake order manager here since we don't start engine in backtest mode
	tm.bot.OrderManager, err = SetupOrderManager(
		tm.bot.ExchangeManager,
		tm.bot.CommunicationsManager,
		&tm.bot.ServicesWG,
		tm.bot.Config.OrderManager.Verbose,
		tm.bot.Config.ProductionMode,
		tm.liveMode,
		tm.bot.Settings.EnableDryRun,
	)
	if err != nil {
		gctlog.Errorf(gctlog.Global, "Order manager unable to setup: %s", err)
	} else {
		err = tm.bot.OrderManager.Start()
		if err != nil {
			gctlog.Errorf(gctlog.Global, "Order manager unable to start: %s", err)
		}
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

	tm.initializeFactorEngines()

	return err
}

func (tm *TradeManager) initializeFactorEngines() error {
	tm.FactorEngines = make(map[string]map[asset.Item]map[currency.Pair]*FactorEngine)
	tm.Datas.Setup()
	for _, cs := range tm.bot.CurrencySettings {
		if tm.FactorEngines[cs.ExchangeName] == nil {
			tm.FactorEngines[cs.ExchangeName] = make(map[asset.Item]map[currency.Pair]*FactorEngine)
		}
		if tm.FactorEngines[cs.ExchangeName][cs.AssetType] == nil {
			tm.FactorEngines[cs.ExchangeName][cs.AssetType] = make(map[currency.Pair]*FactorEngine)
		}
		fe, _ := SetupFactorEngine(cs, &tm.bot.Config.FactorEngine)
		tm.FactorEngines[cs.ExchangeName][cs.AssetType][cs.CurrencyPair] = fe

		var dbData *datakline.DataFromKline
		var err error
		if tm.bot.Settings.EnableLiveMode {
			if tm.debug {
				fmt.Println("initialize factor engine", cs.CurrencyPair)
			}
			// fmt.Println("get data for live", cs.CurrencyPair)
			dbData, err = database.LoadData(
				tm.GetCurrentTime().Add(time.Minute*-300),
				tm.GetCurrentTime().Add(time.Minute*-30),
				time.Minute,
				cs.ExchangeName,
				0,
				cs.CurrencyPair,
				cs.AssetType)
			if err != nil {
				fmt.Println("error initializeFactorEngines:", err)
			}

			if dbData == nil {
				return fmt.Errorf("no bars returned")
			}

			// fmt.Println(dbData.Item)
			barsRet := len(dbData.Item.Candles)
			// fmt.Println("returned", barsRet, "bars")
			if barsRet == 0 {
				if tm.isSimulation {
					panic("no bars returned")
				}
				return fmt.Errorf("no bars returned")
			}
			dbData.Load()
			tm.Datas.SetDataForCurrency(cs.ExchangeName, cs.AssetType, cs.CurrencyPair, dbData)
			if err != nil {
				fmt.Println("error loading db data", err)
				// create a data history request if there isn't one already
				os.Exit(123)
			}
		}
	}
	return nil
}

// Series returns candle data
func CandleSeriesForSettings(e *ExchangeAssetPairSettings, interval int64, start, end time.Time) (out candle.Item, err error) {
	return candle.Series(e.ExchangeName, e.CurrencyPair.Base.String(), e.CurrencyPair.Quote.String(), 60, e.AssetType.String(), start, end)
}

func (tm *TradeManager) loadBacktestData() (err error) {
	if len(tm.bot.CurrencySettings) == 0 {
		panic("no cs")
	}
	for _, eap := range tm.bot.CurrencySettings {
		e := eap.ExchangeName
		a := eap.AssetType
		p := eap.CurrencyPair
		// fmt.Println("loading data for", p)
		startDate := tm.cfg.DataSettings.DatabaseData.StartDate
		endDate := tm.cfg.DataSettings.DatabaseData.EndDate
		dbData, err := database.LoadData(
			startDate,
			endDate,
			time.Minute,
			e,
			0,
			p,
			a)
		if err != nil {
			fmt.Println("loaddata err", err)
			return err
		}

		tm.Datas.SetDataForCurrency(e, a, p, dbData)
		dbData.RangeHolder, err = kline.CalculateCandleDateRanges(
			startDate,
			endDate,
			kline.Interval(kline.OneMin),
			0)
		// fmt.Println("load data for currency", p)
		dbData.Load()
		tm.Reports.AddKlineItem(&dbData.Item)
		// tm.TradeReports.AddKlineItem(&dbData.Item)
		tm.originalCandles = dbData

		if err != nil {
			return fmt.Errorf("error creating range holder. error: %s", err)
		}

		for i := range dbData.RangeHolder.Ranges {
			for j := range dbData.RangeHolder.Ranges[i].Intervals {
				dbData.RangeHolder.Ranges[i].Intervals[j].HasData = true
			}
		}

		// fmt.Println("db data has", len(dbData.Item.Candles))
	}
	// fmt.Println("done loading tm data")

	return err
}

func (tm *TradeManager) loadLatestCandleFromDatabase(eap *ExchangeAssetPairSettings) (*datakline.DataFromKline, error) {
	e := eap.ExchangeName
	a := eap.AssetType
	p := eap.CurrencyPair
	thisMinute := common.ThisMinute()
	startTime := thisMinute.Add(time.Minute * -1)
	dbData, err := database.LoadData(
		startTime,
		thisMinute,
		time.Minute,
		e,
		0,
		p,
		a)
	if err != nil {
		return nil, err
	}

	// validate results
	lastCandle := dbData.Item.Candles[len(dbData.Item.Candles)-1]
	t1 := lastCandle.Time
	// sameTime := (t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day() && t1.Hour() == t2.Hour() && t1.Minute() == t2.Minute())
	if !common.IsSameMinute(t1, thisMinute) {
		// fmt.Println("don't have bar yet", lastCandle.Time, thisMinute)
		return nil, fmt.Errorf("don't have bar yet", lastCandle.Time, thisMinute)
	}
	// fmt.Println("db load data", p, thisMinute)
	dbData.Load()

	tm.Datas.SetDataForCurrency(e, a, p, dbData)
	// dbData.RangeHolder.SetHasDataFromCandles(dbData.Item.Candles)
	// fmt.Println("calculate data in range", startTime, thisMinute, "lasttime", lastCandle.Time)
	dbData.RangeHolder, err = kline.CalculateCandleDateRanges(
		startTime,
		tm.GetCurrentTime(),
		kline.Interval(kline.OneMin),
		0)

	if err != nil {
		return nil, fmt.Errorf("error creating range holder. error: %s", err)
	}

	for i := range dbData.RangeHolder.Ranges {
		for j := range dbData.RangeHolder.Ranges[i].Intervals {
			dbData.RangeHolder.Ranges[i].Intervals[j].HasData = true
		}
	}

	if !dbData.HasDataAtTime(lastCandle.Time) {
		fmt.Println("doesnt have data in range", lastCandle.Time)
		os.Exit(123)
	}

	return dbData, err
}

func (tm *TradeManager) GetCurrentTime() time.Time {
	if tm.isSimulation {
		return tm.currentTime
	}
	return time.Now().UTC()
}

func (tm *TradeManager) incrementMinute() {
	tm.currentTime = tm.currentTime.Add(time.Minute)
}

func arePairsEqual(p1, p2 currency.Pair) bool {
	return strings.EqualFold(p1.Quote.String(), p2.Quote.String()) && strings.EqualFold(p1.Base.String(), p2.Base.String())
}
