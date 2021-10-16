package engine

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	gctcommon "gocryptotrader/common"
	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/data"
	"gocryptotrader/data/kline"
	"gocryptotrader/data/kline/database"
	gctdatabase "gocryptotrader/database"
	"gocryptotrader/database/repository/candle"
	"gocryptotrader/eventtypes"
	"gocryptotrader/eventtypes/cancel"
	"gocryptotrader/eventtypes/event"
	"gocryptotrader/eventtypes/fill"
	"gocryptotrader/eventtypes/order"
	"gocryptotrader/eventtypes/signal"
	"gocryptotrader/eventtypes/submit"
	"gocryptotrader/exchange"
	"gocryptotrader/exchange/asset"
	gctkline "gocryptotrader/exchange/kline"
	gctorder "gocryptotrader/exchange/order"
	"gocryptotrader/exchange/ticker"
	"gocryptotrader/log"
	gctlog "gocryptotrader/log"
	"gocryptotrader/portfolio/report"
	"gocryptotrader/portfolio/risk"
	"gocryptotrader/portfolio/slippage"
	"gocryptotrader/portfolio/statistics"
	"gocryptotrader/portfolio/statistics/currencystatistics"
	"gocryptotrader/portfolio/strategies"

	"github.com/shopspring/decimal"
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
	// tm.Bot = nil
	// tm.FactorEngine = nil
}

// NewFromConfig takes a strategy config and configures a backtester variable to run
func NewTradeManagerFromConfig(cfg *config.Config, templatePath, output string, bot *Engine) (*TradeManager, error) {
	log.Debugln(log.TradeManager, "TradeManager: Initializing... dry run", bot.Config.DryRun)
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
	tm.Bot = bot
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

	err = tm.setupBot(cfg)

	// initialize the data structure to hold the klines for each pair
	return tm, err
}

func (tm *TradeManager) setOrderManagerCallbacks() {
	// tm.Bot.OrderManager.SetOnSubmit(tm.onSubmit)
	tm.OrderManager.SetOnFill(tm.onFill)
	tm.OrderManager.SetOnCancel(tm.onCancel)
}

// BACKTEST FUNCTIONALITY
// Run will iterate over loaded data events
// save them and then handle the event based on its type
func (tm *TradeManager) Run() error {
	tm.setOrderManagerCallbacks()
	log.Debugf(log.TradeManager, "TradeManager Running. Warmup: %v\n", tm.Warmup)
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
								log.Errorf(log.TradeManager, "Unable to perform `Next` for %v %v %v", exchangeName, assetItem, currencyPair)
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
	tm.setOrderManagerCallbacks()

	tm.Warmup = false
	// tm.warmup()

	// throw error if not live
	if !atomic.CompareAndSwapInt32(&tm.started, 0, 1) {
		return fmt.Errorf("backtester %w", ErrSubSystemAlreadyStarted)
	}

	// start trade manager
	log.Debugf(log.TradeManager, "TradeManager  %s", MsgSubSystemStarting)
	tm.shutdown = make(chan struct{})

	go tm.heartBeat()

	// create data subscriptions

	log.Debugln(log.TradeManager, "Running Live")
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
			for _, exchangeMap := range tm.Datas.GetAllData() {
				for _, assetMap := range exchangeMap {
					for _, dataHandler := range assetMap {
						d := dataHandler.Next()
						if d == nil {
							// if !tm.hasHandledEvent {
							// 	log.Errorf(log.TradeManager, "Unable to perform `Next` for %v %v %v", exchangeName, assetItem, currencyPair)
							// }
							return nil
						}
						tm.EventQueue.AppendEvent(d)
					}
				}
			}
			// if err != nil {
			// 	log.Errorln(log.TradeManager, "error loading data events", err)
			// }
			err := tm.processEvents()
			if err != nil {
				log.Errorln(log.TradeManager, "procesing events", err)
			}
		}
	}
	return nil
}

func (tm *TradeManager) loadDataEvents() error {
	return nil
}

func (tm *TradeManager) warmup() error {

	// run the catchup process

	if tm.verbose {
		log.Infoln(log.TradeManager, "Running catchup processes")
	}
	_, err := tm.Bot.dataHistoryManager.Catchup(tm.GetAllCurrencySettings())
	if err != nil {
		log.Infoln(log.TradeManager, "history catchup failed")
		os.Exit(1)
	}
	err = tm.Bot.dataHistoryManager.RunJobs()
	if err != nil {
		return err
	}
	err = tm.Bot.dataHistoryManager.Stop()
	if err != nil {
		return err
	}

	// get latest bars for warmup
	cs, err := tm.GetAllCurrencySettings()
	if err != nil {
		return err
	}

	for _, pair := range cs {
		start := time.Now().Add(time.Minute * -10)
		end := time.Now()

		fmt.Println("loading data for", pair.ExchangeName, pair.CurrencyPair)
		dbData, err := database.LoadData(
			start,
			end,
			time.Minute,
			pair.ExchangeName,
			eventtypes.DataCandle,
			pair.CurrencyPair,
			pair.AssetType)

		if err != nil {
			fmt.Println("error loading db data", err)
		}
		dbData.Load()

		dbData.Item.RemoveDuplicates()
		dbData.Item.SortCandlesByTimestamp(false)
		dbData.RangeHolder, err = gctkline.CalculateCandleDateRanges(
			start,
			end,
			gctkline.Interval(time.Minute),
			0,
		)

		tm.Datas.SetDataForCurrency(
			pair.ExchangeName,
			pair.AssetType,
			pair.CurrencyPair,
			dbData)

		//
		// validate the history is populated with current data
		//
		retCandle, _ := candle.Series(pair.ExchangeName,
			pair.CurrencyPair.Base.String(), pair.CurrencyPair.Quote.String(),
			int64(60), string(pair.AssetType), start, end)
		lc := retCandle.Candles[len(retCandle.Candles)-1].Timestamp
		t := time.Now()
		t1 := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
		t2 := time.Date(lc.Year(), lc.Month(), lc.Day(), lc.Hour(), lc.Minute(), 0, 0, t.Location())

		if t2 != t1 {
			fmt.Println("sync time is off. History Catchup Failed.", t1, t2)
			// os.Exit(1)
		}

		if len(retCandle.Candles) == 0 {
			fmt.Println("No candles returned, History catchup failed. Exiting.")
			os.Exit(1)
		}
	}

	// precache the factor engines
	log.Debugln(log.TradeManager, "Warming up factor engines...")

	tm.Run()

	//
	// validate factor engines are cached
	//
	for _, ex := range tm.FactorEngines {
		for _, a := range ex {
			for _, fe := range a {
				log.Debugf(log.TradeManager, "fe %v %v", fe.Pair, fe.Minute().LastDate())
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

	log.Debugln(log.TradeManager, "Backtester Stopping...")

	if tm.Bot.OrderManager.IsRunning() {
		tm.Bot.OrderManager.Stop()
	}
	// if tm.Bot.DatabaseManager.IsRunning() {
	// 	tm.Bot.DatabaseManager.Stop()
	// }

	for _, s := range tm.Strategies {
		s.Stop()
	}

	close(tm.shutdown)
	tm.Bot.TradeManager = nil
	tm.wg.Wait()
	log.Debugln(log.TradeManager, "Backtester Stopped.")
	return nil
}

func (tm *TradeManager) loadOfflineDatas() error {
	cfg := &tm.cfg

	tm.Datas.Setup()

	// LOAD DATA FOR EVERY PAIR
	for _, cs := range tm.CurrencySettings {
		// exchangeName := strings.ToLower(exch.GetName())
		// klineData, err := tm.loadOfflineData(cfg, exch, pair, a)
		exch, pair, a, err := tm.loadExchangePairAssetBase(
			cs.ExchangeName,
			cs.CurrencyPair.Base.String(),
			cs.CurrencyPair.Quote.String(),
			cs.AssetType.String())
		klineData, err := tm.loadOfflineData(cfg, exch, pair, a)
		if err != nil {
			return err
		}
		tm.Datas.SetDataForCurrency(cs.ExchangeName, cs.AssetType, cs.CurrencyPair, klineData)
	}
	return nil
}

// IsRunning returns if gctscript manager subsystem is started
func (b *TradeManager) IsRunning() bool {
	if b == nil {
		return false
	}
	return atomic.LoadInt32(&b.started) == 1
}

func (tm *TradeManager) setupExchangeSettings(cfg *config.Config) error {
	for _, e := range tm.Bot.Config.GetEnabledExchanges() {
		enabledPairs, _ := tm.Bot.Config.GetEnabledPairs(e, asset.Spot)
		for _, pair := range enabledPairs {
			// fmt.Println("enabledpairs", e, pair)
			_, pair, a, err := tm.loadExchangePairAssetBase(e, pair.Base.String(), pair.Quote.String(), "spot")

			// log.Debugln(log.TradeManager, "setting exchange settings...", pair, a)
			if err != nil {
				return err
			}

			tm.CurrencySettings = append(tm.CurrencySettings, ExchangeAssetPairSettings{
				ExchangeName: e,
				CurrencyPair: pair,
				AssetType:    a,
			})
		}
	}
	return nil
}

func (tm *TradeManager) loadExchangePairAssetBase(exch, base, quote, ass string) (exchange.IBotExchange, currency.Pair, asset.Item, error) {
	e, err := tm.Bot.GetExchangeByName(exch)
	if err != nil {
		return nil, currency.Pair{}, "", err
	}

	var cp, fPair currency.Pair
	cp, err = currency.NewPairFromStrings(base, quote)
	if err != nil {
		return nil, currency.Pair{}, "", err
	}

	var a asset.Item
	a, err = asset.New(ass)
	if err != nil {
		return nil, currency.Pair{}, "", err
	}

	exchangeBase := e.GetBase()
	// if !exchangeBase.ValidateAPICredentials() {
	// 	log.Warnf(log.TradeManager, "no credentials set for %v, this is theoretical only", exchangeBase.Name)
	// }

	fPair, err = exchangeBase.FormatExchangeCurrency(cp, a)
	if err != nil {
		return nil, currency.Pair{}, "", err
	}
	return e, fPair, a, nil
}

// this is only needed in backtest mode, except for when live runs the catchup process
// setupBot sets up a basic bot to retrieve exchange data
// as well as process orders
// setup order manager, exchange manager, database manager
func (tm *TradeManager) setupBot(cfg *config.Config) error {
	var err error

	tm.Datas = &data.HandlerPerCurrency{}
	tm.Datas.Setup()

	// this already run in engine.newfromsettings
	// tm.Bot.ExchangeManager = SetupExchangeManager()

	if !tm.Bot.Config.LiveMode {
		for i := range cfg.CurrencySettings {
			err = tm.Bot.LoadExchange(cfg.CurrencySettings[i].ExchangeName, nil)
			if err != nil && !errors.Is(err, ErrExchangeAlreadyLoaded) {
				return err
			}
		}

		// start fake order manager here since we don't start engine in live mode
		tm.Bot.FakeOrderManager, err = SetupFakeOrderManager(
			tm.Bot,
			tm.Bot.ExchangeManager,
			tm.Bot.CommunicationsManager,
			&tm.Bot.ServicesWG,
			tm.Bot.Settings.Verbose,
		)
		if err != nil {
			gctlog.Errorf(gctlog.Global, "Fake Order manager unable to setup: %s", err)
		} else {
			err = tm.Bot.FakeOrderManager.Start()
			if err != nil {
				gctlog.Errorf(gctlog.Global, "Fake Order manager unable to start: %s", err)
			}
			tm.Bot.OrderManager = tm.Bot.FakeOrderManager
		}

		tm.Bot.DatabaseManager, err = SetupDatabaseConnectionManager(gctdatabase.DB.GetConfig())
		if err != nil {
			return err
		} else {
			err = tm.Bot.DatabaseManager.Start(&tm.Bot.ServicesWG)
			if err != nil {
				gctlog.Errorf(gctlog.Global, "Database manager unable to start: %v", err)
			}
		}
	}

	err = tm.setupExchangeSettings(cfg)
	if err != nil {
		fmt.Println("error setting up exchange settings", cfg, err)
		return err
	}

	if err != nil {
		return err
	}

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

	// for i := range cfg.Stratey

	var slit []strategies.Handler
	var s strategies.Handler

	count := 0
	for _, strat := range cfg.StrategiesSettings {
		// for _, dir := range []gctorder.Side{gctorder.Buy, gctorder.Sell} {
		for _, c := range tm.CurrencySettings {
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

	// if tm.verbose {
	log.Infof(log.TradeManager, "Running %d strategies\n", len(tm.Strategies))
	log.Infof(log.TradeManager, "Running %d currencies\n", len(tm.CurrencySettings))
	// log.Infof(log.TradeManager, "Running %d exchanges\n", len(tm.CurrencySettings))
	// }

	// setup portfolio with strategies
	var p *Portfolio
	p, err = SetupPortfolio(tm.Strategies, tm.Bot, sizeManager, portfolioRisk, cfg.StatisticSettings.RiskFreeRate)
	p.SetVerbose(false)
	if err != nil {
		return err
	}

	stats := &statistics.Statistic{
		StrategyName:                "ok",
		StrategyNickname:            cfg.Nickname,
		StrategyDescription:         "ok",
		StrategyGoal:                cfg.Goal,
		ExchangeAssetPairStatistics: make(map[string]map[asset.Item]map[currency.Pair]*currencystatistics.CurrencyStatistic),
		RiskFreeRate:                cfg.StatisticSettings.RiskFreeRate,
	}
	tm.Statistic = stats

	// load from configuration into datastructure
	// currencysettings returns the data from the config, exchangeassetpairsettings
	// tm.Exchange = &e

	log.Infoln(log.TradeManager, "Loaded", len(tm.CurrencySettings), "currencies")

	tm.Portfolio = p
	allCS, _ := tm.GetAllCurrencySettings()
	// tm.FactorEngines = make(map[currency.Pair]*FactorEngine)
	tm.FactorEngines = make(map[string]map[asset.Item]map[currency.Pair]*FactorEngine)
	for _, x := range allCS {
		fe, _ := SetupFactorEngine(x.CurrencyPair)
		ex := strings.ToLower(x.ExchangeName)

		if tm.FactorEngines[ex] == nil {
			// fmt.Println("setup exchanges factors engine", ex)
			tm.FactorEngines[ex] = make(map[asset.Item]map[currency.Pair]*FactorEngine)
		}
		if tm.FactorEngines[ex][x.AssetType] == nil {
			// fmt.Println("setup exchanges's asset's factor engines", ex, "spot")
			tm.FactorEngines[ex][x.AssetType] = make(map[currency.Pair]*FactorEngine)
		}

		// fmt.Println("setup factor engine for pair", ex, x.CurrencyPair)
		tm.FactorEngines[ex][x.AssetType][x.CurrencyPair] = fe
	}

	// cfg.PrintSetting()

	err = tm.loadOfflineDatas()
	if err != nil {
		log.Errorf(log.TradeManager, "error loading datas", err)
		return nil
	}

	return nil
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
		log.Errorf(log.TradeManager, "Could not retrieve taker fee for %v. %v", exch.GetName(), err)
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
		log.Errorf(log.TradeManager, "Could not retrieve maker fee for %v. %v", exch.GetName(), err)
	}

	return decimal.NewFromFloat(fMakerFee), decimal.NewFromFloat(fTakerFee)
}

// loadOfflineData will create kline data from the sources defined in start config files. It can exist from databases, csv or API endpoints
// it can also be generated from trade data which will be converted into kline data
func (tm *TradeManager) loadOfflineData(cfg *config.Config, exch exchange.IBotExchange, fPair currency.Pair, a asset.Item) (*kline.DataFromKline, error) {
	if exch == nil {
		return nil, ErrExchangeNotFound
	}
	b := exch.GetBase()

	dataType, err := eventtypes.DataTypeToInt(cfg.DataSettings.DataType)
	if err != nil {
		return nil, err
	}

	resp := &kline.DataFromKline{}
	// log.Infof(log.TradeManager, "loading db data for %v %v %v...\n", exch.GetName(), a, fPair)
	if cfg.DataSettings.DatabaseData.InclusiveEndDate {
		cfg.DataSettings.DatabaseData.EndDate = cfg.DataSettings.DatabaseData.EndDate.Add(cfg.DataSettings.Interval)
	}
	if cfg.DataSettings.DatabaseData.ConfigOverride != nil {
		tm.Bot.Config.Database = *cfg.DataSettings.DatabaseData.ConfigOverride
		gctdatabase.DB.DataPath = filepath.Join(gctcommon.GetDefaultDataDir(runtime.GOOS), "database")
		err = gctdatabase.DB.SetConfig(cfg.DataSettings.DatabaseData.ConfigOverride)
		if err != nil {
			return nil, err
		}
	}
	resp, err = loadDatabaseData(cfg, exch.GetName(), fPair, a, dataType)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from GoCryptoTrader database. Error: %v. Please ensure the database is setup correctly and has data before use", err)
	}

	resp.Item.RemoveDuplicates()
	resp.Item.SortCandlesByTimestamp(false)
	resp.RangeHolder, err = gctkline.CalculateCandleDateRanges(
		cfg.DataSettings.DatabaseData.StartDate,
		cfg.DataSettings.DatabaseData.EndDate,
		gctkline.Interval(cfg.DataSettings.Interval),
		0,
	)
	if err != nil {
		return nil, err
	}
	resp.RangeHolder.SetHasDataFromCandles(resp.Item.Candles)
	summary := resp.RangeHolder.DataSummary(false)
	if len(summary) > 0 {
		log.Warnf(log.TradeManager, "%v", summary)
	}
	if resp == nil {
		return nil, fmt.Errorf("processing error, response returned nil")
	}

	err = b.ValidateKline(fPair, a, resp.Item.Interval)
	if err != nil {
		if dataType != eventtypes.DataTrade || !strings.EqualFold(err.Error(), "interval not supported") {
			return nil, err
		}
	}

	err = resp.Load()
	if err != nil {
		return nil, err
	}
	tm.Reports.AddKlineItem(&resp.Item)
	return resp, nil
}

// -----------------------------------
// DATA PROCESSING
// -----------------------------------
// handleEvent is the main processor of data for the backtester
// after data has been loaded and Run has appended a data event to the queue,
// handle event will process events and add further events to the queue if they
// are required
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
	// err := tm.updateStatsForDataEvent(ev)
	// if err != nil {
	// 	return err
	// }

	d := tm.Datas.GetDataForCurrency(ev.GetExchange(), ev.GetAssetType(), ev.Pair())

	// update factor engine
	// if tm.Bot.Config.LiveMode {
	// 	fmt.Println("factor on bar update", ev.Pair(), d.Latest().GetTime(), len(tm.FactorEngines[ev.Pair()].Minute().Close))
	// }
	// exfe := tm.FactorEngines[ev.GetExchange()]
	// exassetfe := tm.FactorEngines[ev.GetExchange()][ev.GetAssetType()]
	// fmt.Println(len(tm.FactorEngines))
	// for i := range exassetfe {
	// 	fmt.Printf("%s..", i)
	// }
	// fmt.Println(ev.GetExchange(), len(exfe), len(exassetfe), ev.Pair(), fe)
	fe := tm.FactorEngines[ev.GetExchange()][ev.GetAssetType()][ev.Pair()]
	fe.OnBar(d)

	// HANDLE warmup MODE
	// in warmup mode, we do not query the strategies
	if !tm.Warmup {
		tm.Bot.OrderManager.Update()

		// var s signal.Event
		// for _, strategy := range tm.Strategies {
		// 	if strategy.GetPair() == ev.Pair() {
		// 		if tm.Bot.Config.LiveMode {
		// 			fmt.Println("Updating strategy", strategy.GetID(), d.Latest().GetTime())
		// 		}
		// 		s, err = strategy.OnData(d, tm.Portfolio, fe)
		// 		tm.EventQueue.AppendEvent(s)
		// 	}
		// }
	}

	// if err != nil {
	// 	if errors.Is(err, base.ErrTooMuchBadData) {
	// 		// too much bad data is a severe error and backtesting must cease
	// 		return err
	// 	}
	// 	log.Error(log.TradeManager, err)
	// 	return nil
	// }
	// err = tm.Statistic.SetEventForOffset(ev)
	// if err != nil {
	// 	log.Error(log.TradeManager, err)
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
	// var dataEvents []data.Handler
	// dataHandlerMap := tm.Datas.GetAllData()
	// for _, exchangeMap := range dataHandlerMap {
	// 	for _, assetMap := range exchangeMap {
	// 		for _, dataHandler := range assetMap {
	// 			latestData := dataHandler.Latest()
	// 			err := tm.updateStatsForDataEvent(latestData)
	// 			if err != nil && err == statistics.ErrAlreadyProcessed {
	// 				continue
	// 			}
	// 			dataEvents = append(dataEvents, dataHandler)
	// 		}
	// 	}
	// }
	// d, _ := base.GetBaseData(dataEvents)
	// // NOTE DISABLED
	// // signals, err := tm.Strategies[0].OnSimultaneousSignals(dataEvents, tm.Portfolio, tm.FactorEngines[d.CurrencyPair()])
	// if err != nil {
	// 	if errors.Is(err, base.ErrTooMuchBadData) {
	// 		// too much bad data is a severe error and backtesting must cease
	// 		return err
	// 	}
	// 	log.Error(log.TradeManager, err)
	// 	return nil
	// }
	// for i := range signals {
	// 	err = tm.Statistic.SetEventForOffset(signals[i])
	// 	if err != nil {
	// 		log.Error(log.TradeManager, err)
	// 	}
	// 	tm.EventQueue.AppendEvent(signals[i])
	// }
	return nil
}

//
// Event Processors
//
// processSignalEvent receives an event from the strategy for processing under the portfolio
func (tm *TradeManager) processSignalEvent(ev signal.Event) {
	cs, err := tm.GetCurrencySettings(ev.GetExchange(), ev.GetAssetType(), ev.Pair())
	if err != nil {
		log.Error(log.TradeManager, err)
		return
	}
	var o *order.Order
	o, err = tm.Portfolio.OnSignal(ev, &cs)
	if err != nil {
		log.Error(log.TradeManager, err)
		return
	}
	if err != nil {
		log.Error(log.TradeManager, err)
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
		log.Error(log.TradeManager, "submit event has no strategy ID")
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
		log.Error(log.TradeManager, "order event has no strategy ID")
	}
	// else {
	// 	// gctlog.Debugln(log.TradeManager, "creating order for", o.GetStrategyID())
	// }
	d := tm.Datas.GetDataForCurrency(o.GetExchange(), o.GetAssetType(), o.Pair())
	// this blocks and returns a submission event
	submitEvent, err := tm.ExecuteOrder(o, d, tm.Bot.FakeOrderManager)

	// call on submit here

	if err != nil {
		log.Error(log.TradeManager, err)
		return
	}

	if submitEvent.GetStrategyID() == "" {
		log.Error(log.TradeManager, "Not strategy ID in order event")
		return
	}

	tm.EventQueue.AppendEvent(submitEvent)
}

func (tm *TradeManager) heartBeat() {
	time.Sleep(time.Second * 10)
	fmt.Println("........................HEARTBEAT")
	exchanges, _ := tm.Bot.ExchangeManager.GetExchanges()
	ex := exchanges[0]
	fmt.Println("subscribing to ", tm.CurrencySettings[0])
	pipe, err := ticker.SubscribeTicker(ex.GetName(), tm.CurrencySettings[0].CurrencyPair, asset.Spot)
	if err != nil {
		fmt.Println(".........error subscribing to ticker", err)
		// wait and retry
	}

	// defer func() {
	// }()

	for {
		select {
		case <-tm.shutdown:
			pipeErr := pipe.Release()
			if pipeErr != nil {
				log.Error(log.DispatchMgr, pipeErr)
			}
			return
		case data, ok := <-pipe.C:
			if !ok {
				fmt.Println("error dispatch system")
				return
			}
			t := (*data.(*interface{})).(ticker.Price)
			fmt.Println("received data", t)
		}
		// err := stream.Send(&gctrpc.TickerResponse{
		// 	Pair: &gctrpc.CurrencyPair{
		// 		Base:      t.Pair.Base.String(),
		// 		Quote:     t.Pair.Quote.String(),
		// 		Delimiter: t.Pair.Delimiter},
		// 	LastUpdated: s.unixTimestamp(t.LastUpdated),
		// 	Last:        t.Last,
		// 	High:        t.High,
		// 	Low:         t.Low,
		// 	Bid:         t.Bid,
		// 	Ask:         t.Ask,
		// 	Volume:      t.Volume,
		// 	PriceAth:    t.PriceATH,
		// })
		// if err != nil {
		// 	return err
		// }
	}
	fmt.Println("finished")
	// if err != nil {
	// 	return err
	// }

	// tm.wg.Add(1)
	// tick := time.NewTicker(time.Second * 5)
	// defer func() {
	// 	tick.Stop()
	// 	tm.wg.Done()
	// }()
	// for {
	// 	select {
	// 	case <-tm.shutdown:
	// 		return
	// 	case <-tick.C:
	// 		exchanges, err := tm.Bot.ExchangeManager.GetExchanges()
	// 		for _, ex := range exchanges {
	// 			if err != nil {
	// 				log.Infoln(log.TradeManager, "error getting tick", err)
	// 			}
	//
	// 			for _, cp := range tm.CurrencySettings {
	// 				tick, _ := ex.FetchTicker(context.Background(), cp.CurrencyPair, asset.Spot)
	// 				t1 := time.Now()
	// 				// ticker := m.currencyPairs[x].Ticker
	// 				secondsAgo := int(t1.Sub(tick.LastUpdated).Seconds())
	// 				if secondsAgo > 10 {
	// 					log.Warnln(log.TradeManager, cp.CurrencyPair, tick.Last, secondsAgo)
	// 				} else {
	// 					log.Infoln(log.TradeManager, cp.CurrencyPair, tick.Last, secondsAgo)
	// 				}
	// 			}
	// 		}
	// 		// tm.PrintTradingDetails()
	// 	}
	// }
}

func (tm *TradeManager) PrintTradingDetails() {
	// fmt.Println("strategies running", len(tm.Strategies))
	log.Infoln(log.TradeManager, len(tm.Strategies), "strategies running")

	for _, cs := range tm.CurrencySettings {
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
			log.Infoln(log.TradeManager, cs.CurrencyPair, "last updated", secondsAgo, "seconds ago")
		}
		// else {
		// 	log.Debugln(log.TradeManager, cs.CurrencyPair, "last updated", secondsAgo, "seconds ago")
		// }
	}
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
		log.Error(log.TradeManager, err)
	}
	// update portfolio manager with the latest price
	err = tm.Portfolio.UpdateHoldings(ev)
	if err != nil {
		log.Error(log.TradeManager, err)
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
		cfg.DataSettings.DatabaseData.StartDate,
		cfg.DataSettings.DatabaseData.EndDate,
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

func (tm *TradeManager) GetAllCurrencySettings() ([]ExchangeAssetPairSettings, error) {
	return tm.CurrencySettings, nil
}

// SetExchangeAssetCurrencySettings sets the settings for an exchange, asset, currency
func (tm *TradeManager) SetExchangeAssetCurrencySettings(exch string, a asset.Item, cp currency.Pair, c *ExchangeAssetPairSettings) {
	if c.ExchangeName == "" ||
		c.AssetType == "" ||
		c.CurrencyPair.IsEmpty() {
		return
	}

	for i := range tm.CurrencySettings {
		if tm.CurrencySettings[i].CurrencyPair == cp &&
			tm.CurrencySettings[i].AssetType == a &&
			exch == tm.CurrencySettings[i].ExchangeName {
			tm.CurrencySettings[i] = *c
			return
		}
	}
	tm.CurrencySettings = append(tm.CurrencySettings, *c)
}

// GetCurrencySettings returns the settings for an exchange, asset currency
func (tm *TradeManager) GetCurrencySettings(exch string, a asset.Item, cp currency.Pair) (ExchangeAssetPairSettings, error) {
	for i := range tm.CurrencySettings {
		if tm.CurrencySettings[i].CurrencyPair.Equal(cp) {
			if tm.CurrencySettings[i].AssetType == a {
				if exch == tm.CurrencySettings[i].ExchangeName {
					return tm.CurrencySettings[i], nil
				}
			}
		}
	}
	return ExchangeAssetPairSettings{}, fmt.Errorf("no currency settings found for %v %v %v", exch, a, cp)
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
		fmt.Println("error order has no internal order id")
	}

	if ev.IsOrderPlaced {
		// fmt.Println("TM ORDERPLACED, create fill event")
		tm.onFill(omr)
	} else {
		fmt.Println("TM ERROR: ORDERPLACED NOT")
	}

	return ev, nil
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
