package engine

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"gocryptotrader/backtester/statistics"
	"gocryptotrader/backtester/statistics/currencystatistics"
	gctcommon "gocryptotrader/common"
	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/data"
	"gocryptotrader/data/kline"
	"gocryptotrader/data/kline/api"
	"gocryptotrader/data/kline/csv"
	"gocryptotrader/data/kline/database"
	"gocryptotrader/data/kline/live"
	gctdatabase "gocryptotrader/database"
	"gocryptotrader/eventtypes"
	"gocryptotrader/eventtypes/fill"
	"gocryptotrader/eventtypes/order"
	"gocryptotrader/eventtypes/signal"
	gctexchange "gocryptotrader/exchanges"
	"gocryptotrader/exchanges/asset"
	gctkline "gocryptotrader/exchanges/kline"
	gctorder "gocryptotrader/exchanges/order"
	"gocryptotrader/log"
	"gocryptotrader/portfolio/report"
	"gocryptotrader/portfolio/risk"
	"gocryptotrader/portfolio/slippage"
	"gocryptotrader/portfolio/strategies"
	"gocryptotrader/portfolio/strategies/base"

	"gocryptotrader/portfolio/compliance"

	"github.com/shopspring/decimal"
)

// New returns a new TradeManager instance
func NewTradeManager() *TradeManager {
	return &TradeManager{
		shutdown: make(chan struct{}),
	}
}

// Reset TradeManager values to default
func (tm *TradeManager) Reset() {
	tm.EventQueue.Reset()
	tm.Datas.Reset()
	tm.Portfolio.Reset()
	tm.Statistic.Reset()
	tm.Exchange.Reset()
	// reset live trades here
	tm.Bot = nil
	tm.FactorEngine = nil
}

// NewFromConfig takes a strategy config and configures a backtester variable to run
func NewTradeManagerFromConfig(cfg *config.Config, templatePath, output string, bot *Engine, live bool) (*TradeManager, error) {
	log.Infoln(log.TradeManager, "Backtest: Loading config...")
	if cfg == nil {
		return nil, errNilConfig
	}
	if bot == nil {
		return nil, errNilBot
	}
	tm := NewTradeManager()
	tm.cfg = *cfg
	tm.IsLive = live
	tm.Warmup = live

	tm.Datas = &data.HandlerPerCurrency{}
	tm.EventQueue = &Holder{}
	reports := &report.Data{
		Config:       cfg,
		TemplatePath: templatePath,
		OutputPath:   output,
	}
	tm.Reports = reports

	err := tm.setupBot(cfg, bot)
	if err != nil {
		return nil, err
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
	for i := range cfg.CurrencySettings {
		if portfolioRisk.CurrencySettings[cfg.CurrencySettings[i].ExchangeName] == nil {
			portfolioRisk.CurrencySettings[cfg.CurrencySettings[i].ExchangeName] = make(map[asset.Item]map[currency.Pair]*risk.CurrencySettings)
		}
		var a asset.Item
		a, err = asset.New(cfg.CurrencySettings[i].Asset)
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
		var exch gctexchange.IBotExchange
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
			log.Warnf(log.TradeManager, "maker fee '%v' should not exceed taker fee '%v'. Please review config",
				cfg.CurrencySettings[i].MakerFee,
				cfg.CurrencySettings[i].TakerFee)
		}
	}

	// LOAD ALL STRATEGIES HERE

	var slit []strategies.Handler

	s, err := strategies.LoadStrategyByName("trend", "SELL", false)
	s.SetID("trend_SELL")
	s.SetDefaults()
	slit = append(slit, s)

	s, err = strategies.LoadStrategyByName("trend", "BUY", false)
	s.SetDefaults()
	s.SetID("trend_BUY")
	slit = append(slit, s)

	s, err = strategies.LoadStrategyByName("trend2", "SELL", false)
	if err != nil {
		fmt.Println("error", err)
	}
	s.SetID("trend2_SELL")
	s.SetDefaults()
	slit = append(slit, s)

	s, err = strategies.LoadStrategyByName("trend2", "BUY", false)
	s.SetDefaults()
	s.SetID("trend2_BUY")
	slit = append(slit, s)

	tm.Strategies = slit

	log.Infof(log.TradeManager, "Loaded %d strategies\n", len(tm.Strategies))

	var p *Portfolio
	p, err = SetupPortfolio(tm.Strategies, bot, sizeManager, portfolioRisk, cfg.StatisticSettings.RiskFreeRate, tm.IsLive)
	if err != nil {
		return nil, err
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
	// reports.Statistics = stats

	e, err := tm.setupExchangeSettings(cfg)
	if err != nil {
		return nil, err
	}

	// load from configuration into datastructure
	// currencysettings returns the data from the config, exchangeassetpairsettings
	tm.Exchange = &e
	for i := range e.CurrencySettings {
		var lookup *PortfolioSettings
		lookup, err = p.SetupCurrencySettingsMap(e.CurrencySettings[i].ExchangeName, e.CurrencySettings[i].AssetType, e.CurrencySettings[i].CurrencyPair)
		if err != nil {
			return nil, err
		}
		lookup.Fee = e.CurrencySettings[i].TakerFee
		lookup.Leverage = e.CurrencySettings[i].Leverage
		lookup.BuySideSizing = e.CurrencySettings[i].BuySide
		lookup.SellSideSizing = e.CurrencySettings[i].SellSide
		lookup.ComplianceManager = compliance.Manager{
			Snapshots: []compliance.Snapshot{},
		}
	}
	tm.Portfolio = p

	// cfg.PrintSetting()

	fe, _ := SetupFactorEngine()
	tm.FactorEngine = fe

	return tm, nil
}

// Run will iterate over loaded data events
// save them and then handle the event based on its type
func (tm *TradeManager) Run() error {
	log.Info(log.TradeManager, "running trade manager")
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

// LIVE FUNCTIONALITY
func (tm *TradeManager) Start() error {
	// throw error if not live
	if !atomic.CompareAndSwapInt32(&tm.started, 0, 1) {
		return fmt.Errorf("backtester %w", ErrSubSystemAlreadyStarted)
	}
	log.Debugf(log.TradeManager, "TradeManager  %s", MsgSubSystemStarting)
	tm.shutdown = make(chan struct{})
	tm.FactorEngine.Start()
	// go tm.heartBeat()
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
			for ev := tm.EventQueue.NextEvent(); ; ev = tm.EventQueue.NextEvent() {
				if ev == nil {
					dataHandlerMap := tm.Datas.GetAllData()
					for exchangeName, exchangeMap := range dataHandlerMap {
						for assetItem, assetMap := range exchangeMap {
							for currencyPair, dataHandler := range assetMap {
								d := dataHandler.Next()
								if d == nil {
									if !tm.hasHandledEvent {
										log.Errorf(log.TradeManager, "Unable to perform `Next` for %v %v %v", exchangeName, assetItem, currencyPair)
									}
								}
								tm.EventQueue.AppendEvent(d)
							}
						}
					}
				}

				if ev != nil {
					err := tm.handleEvent(ev)
					if err != nil {
						return err
					}
				}
				if !tm.hasHandledEvent {
					tm.hasHandledEvent = true
				}
			}
		}
	}
	return nil
}

// Stop shuts down the live data loop
func (tm *TradeManager) Stop() error {
	log.Debugln(log.TradeManager, "stopping backtester")

	// if g == nil {
	// 	return fmt.Errorf("%s %w", caseName, ErrNilSubsystem)
	// }
	// if atomic.LoadInt32(&g.started) == 0 {
	// 	return fmt.Errorf("%s not running", caseName)
	// }
	// defer func() {
	// 	atomic.CompareAndSwapInt32(&g.started, 1, 0)
	// }()
	//
	// err := g.ShutdownAll()
	// if err != nil {
	// 	return err
	// }

	for _, s := range tm.Strategies {
		s.Stop()
	}

	tm.Bot.DatabaseManager.Stop()

	close(tm.shutdown)
	return nil
}

func (tm *TradeManager) setupExchangeSettings(cfg *config.Config) (Exchange, error) {
	log.Infoln(log.TradeManager, "setting exchange settings...")
	resp := Exchange{}

	for i := range cfg.CurrencySettings {
		exch, pair, a, err := tm.loadExchangePairAssetBase(
			cfg.CurrencySettings[i].ExchangeName,
			cfg.CurrencySettings[i].Base,
			cfg.CurrencySettings[i].Quote,
			cfg.CurrencySettings[i].Asset)
		if err != nil {
			return resp, err
		}

		exchangeName := strings.ToLower(exch.GetName())
		tm.Datas.Setup()
		klineData, err := tm.loadData(cfg, exch, pair, a)
		if err != nil {
			return resp, err
		}
		tm.Datas.SetDataForCurrency(exchangeName, a, pair, klineData)
		var makerFee, takerFee decimal.Decimal
		if cfg.CurrencySettings[i].MakerFee.GreaterThan(decimal.Zero) {
			makerFee = cfg.CurrencySettings[i].MakerFee
		}
		if cfg.CurrencySettings[i].TakerFee.GreaterThan(decimal.Zero) {
			takerFee = cfg.CurrencySettings[i].TakerFee
		}
		if makerFee.IsZero() || takerFee.IsZero() {
			var apiMakerFee, apiTakerFee decimal.Decimal
			apiMakerFee, apiTakerFee = getFees(context.TODO(), exch, pair)
			if makerFee.IsZero() {
				makerFee = apiMakerFee
			}
			if takerFee.IsZero() {
				takerFee = apiTakerFee
			}
		}

		if cfg.CurrencySettings[i].MaximumSlippagePercent.LessThan(decimal.Zero) {
			log.Warnf(log.TradeManager, "invalid maximum slippage percent '%v'. Slippage percent is defined as a number, eg '100.00', defaulting to '%v'",
				cfg.CurrencySettings[i].MaximumSlippagePercent,
				slippage.DefaultMaximumSlippagePercent)
			cfg.CurrencySettings[i].MaximumSlippagePercent = slippage.DefaultMaximumSlippagePercent
		}
		if cfg.CurrencySettings[i].MaximumSlippagePercent.IsZero() {
			cfg.CurrencySettings[i].MaximumSlippagePercent = slippage.DefaultMaximumSlippagePercent
		}
		if cfg.CurrencySettings[i].MinimumSlippagePercent.LessThan(decimal.Zero) {
			log.Warnf(log.TradeManager, "invalid minimum slippage percent '%v'. Slippage percent is defined as a number, eg '80.00', defaulting to '%v'",
				cfg.CurrencySettings[i].MinimumSlippagePercent,
				slippage.DefaultMinimumSlippagePercent)
			cfg.CurrencySettings[i].MinimumSlippagePercent = slippage.DefaultMinimumSlippagePercent
		}
		if cfg.CurrencySettings[i].MinimumSlippagePercent.IsZero() {
			cfg.CurrencySettings[i].MinimumSlippagePercent = slippage.DefaultMinimumSlippagePercent
		}
		if cfg.CurrencySettings[i].MaximumSlippagePercent.LessThan(cfg.CurrencySettings[i].MinimumSlippagePercent) {
			cfg.CurrencySettings[i].MaximumSlippagePercent = slippage.DefaultMaximumSlippagePercent
		}

		realOrders := false
		if cfg.DataSettings.LiveData != nil {
			realOrders = cfg.DataSettings.LiveData.RealOrders
		}

		buyRule := config.MinMax{
			MinimumSize:  cfg.CurrencySettings[i].BuySide.MinimumSize,
			MaximumSize:  cfg.CurrencySettings[i].BuySide.MaximumSize,
			MaximumTotal: cfg.CurrencySettings[i].BuySide.MaximumTotal,
		}
		sellRule := config.MinMax{
			MinimumSize:  cfg.CurrencySettings[i].SellSide.MinimumSize,
			MaximumSize:  cfg.CurrencySettings[i].SellSide.MaximumSize,
			MaximumTotal: cfg.CurrencySettings[i].SellSide.MaximumTotal,
		}

		limits, err := exch.GetOrderExecutionLimits(a, pair)
		if err != nil && !errors.Is(err, gctorder.ErrExchangeLimitNotLoaded) {
			return resp, err
		}

		if limits != nil {
			if !cfg.CurrencySettings[i].CanUseExchangeLimits {
				log.Warnf(log.TradeManager, "exchange %s order execution limits supported but disabled for %s %s, live results may differ",
					cfg.CurrencySettings[i].ExchangeName,
					pair,
					a)
				cfg.CurrencySettings[i].ShowExchangeOrderLimitWarning = true
			}
		}

		resp.CurrencySettings = append(resp.CurrencySettings, ExchangeAssetPairSettings{
			ExchangeName:        cfg.CurrencySettings[i].ExchangeName,
			MinimumSlippageRate: cfg.CurrencySettings[i].MinimumSlippagePercent,
			MaximumSlippageRate: cfg.CurrencySettings[i].MaximumSlippagePercent,
			CurrencyPair:        pair,
			AssetType:           a,
			ExchangeFee:         takerFee,
			MakerFee:            takerFee,
			TakerFee:            makerFee,
			UseRealOrders:       realOrders,
			BuySide:             buyRule,
			SellSide:            sellRule,
			Leverage: config.Leverage{
				CanUseLeverage:                 cfg.CurrencySettings[i].Leverage.CanUseLeverage,
				MaximumLeverageRate:            cfg.CurrencySettings[i].Leverage.MaximumLeverageRate,
				MaximumOrdersWithLeverageRatio: cfg.CurrencySettings[i].Leverage.MaximumOrdersWithLeverageRatio,
			},
			Limits:                  limits,
			SkipCandleVolumeFitting: cfg.CurrencySettings[i].SkipCandleVolumeFitting,
			CanUseExchangeLimits:    cfg.CurrencySettings[i].CanUseExchangeLimits,
		})
	}

	return resp, nil
}

func (tm *TradeManager) loadExchangePairAssetBase(exch, base, quote, ass string) (gctexchange.IBotExchange, currency.Pair, asset.Item, error) {
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
	if !exchangeBase.ValidateAPICredentials() {
		log.Warnf(log.TradeManager, "no credentials set for %v, this is theoretical only", exchangeBase.Name)
	}

	fPair, err = exchangeBase.FormatExchangeCurrency(cp, a)
	if err != nil {
		return nil, currency.Pair{}, "", err
	}
	return e, fPair, a, nil
}

// setupBot sets up a basic bot to retrieve exchange data
// as well as process orders
// setup order manager, exchange manager, database manager
func (tm *TradeManager) setupBot(cfg *config.Config, bot *Engine) error {
	var err error
	tm.Bot = bot
	tm.Bot.ExchangeManager = SetupExchangeManager()

	for i := range cfg.CurrencySettings {
		err = tm.Bot.LoadExchange(cfg.CurrencySettings[i].ExchangeName, nil)
		if err != nil && !errors.Is(err, ErrExchangeAlreadyLoaded) {
			return err
		}
	}

	// // if not live since we don't start the engine in backtest mode
	if !tm.IsLive {
		tm.Bot.DatabaseManager, err = SetupDatabaseConnectionManager(gctdatabase.DB.GetConfig())
		if err != nil {
			return err
		}

		err = tm.Bot.DatabaseManager.Start(&tm.Bot.ServicesWG)
		if err != nil {
			return err
		}

		// setup order manager so it can receive orders (and fill them?)
		// when live, the engine starts its own order manager so we do it here for backtesting
		if !tm.Bot.OrderManager.IsRunning() {
			tm.Bot.OrderManager, err = SetupOrderManager(
				tm.Bot.ExchangeManager,
				tm.Bot.CommunicationsManager,
				&tm.Bot.ServicesWG,
				bot.Settings.Verbose)
			if err != nil {
				return err
			}
			err = tm.Bot.OrderManager.Start()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// getFees will return an exchange's fee rate from GCT's wrapper function
func getFees(ctx context.Context, exch gctexchange.IBotExchange, fPair currency.Pair) (makerFee, takerFee decimal.Decimal) {
	fTakerFee, err := exch.GetFeeByType(ctx,
		&gctexchange.FeeBuilder{FeeType: gctexchange.OfflineTradeFee,
			Pair:          fPair,
			IsMaker:       false,
			PurchasePrice: 1,
			Amount:        1,
		})
	if err != nil {
		log.Errorf(log.TradeManager, "Could not retrieve taker fee for %v. %v", exch.GetName(), err)
	}

	fMakerFee, err := exch.GetFeeByType(ctx,
		&gctexchange.FeeBuilder{
			FeeType:       gctexchange.OfflineTradeFee,
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

func (tm *TradeManager) ReloadData() error {
	cfg := &tm.cfg
	for i := range cfg.CurrencySettings {
		exch, pair, a, err := tm.loadExchangePairAssetBase(
			cfg.CurrencySettings[i].ExchangeName,
			cfg.CurrencySettings[i].Base,
			cfg.CurrencySettings[i].Quote,
			cfg.CurrencySettings[i].Asset)
		if err != nil {
			return err
		}

		exchangeName := strings.ToLower(exch.GetName())
		tm.Datas.Setup()
		klineData, err := tm.loadData(cfg, exch, pair, a)
		if err != nil {
			return err
		}
		tm.Datas.SetDataForCurrency(exchangeName, a, pair, klineData)
	}
	return nil
}

// loadData will create kline data from the sources defined in start config files. It can exist from databases, csv or API endpoints
// it can also be generated from trade data which will be converted into kline data
func (tm *TradeManager) loadData(cfg *config.Config, exch gctexchange.IBotExchange, fPair currency.Pair, a asset.Item) (*kline.DataFromKline, error) {
	if exch == nil {
		return nil, ErrExchangeNotFound
	}
	b := exch.GetBase()
	if cfg.DataSettings.DatabaseData == nil &&
		cfg.DataSettings.LiveData == nil &&
		cfg.DataSettings.APIData == nil &&
		cfg.DataSettings.CSVData == nil {
		return nil, errNoDataSource
	}
	// if (cfg.DataSettings.APIData != nil && cfg.DataSettings.DatabaseData != nil) ||
	// 	(cfg.DataSettings.APIData != nil && cfg.DataSettings.LiveData != nil) ||
	// 	(cfg.DataSettings.APIData != nil && cfg.DataSettings.CSVData != nil) ||
	// 	(cfg.DataSettings.DatabaseData != nil && cfg.DataSettings.LiveData != nil) ||
	// 	(cfg.DataSettings.CSVData != nil && cfg.DataSettings.LiveData != nil) ||
	// 	(cfg.DataSettings.CSVData != nil && cfg.DataSettings.DatabaseData != nil) {
	// 	return nil, errAmbiguousDataSource
	// }

	dataType, err := eventtypes.DataTypeToInt(cfg.DataSettings.DataType)
	if err != nil {
		return nil, err
	}

	resp := &kline.DataFromKline{}
	switch {
	case cfg.DataSettings.CSVData != nil:
		log.Infof(log.TradeManager, "loading csv data for %v %v %v...\n", exch.GetName(), a, fPair)
		if cfg.DataSettings.Interval <= 0 {
			return nil, errIntervalUnset
		}
		resp, err = csv.LoadData(
			dataType,
			cfg.DataSettings.CSVData.FullPath,
			strings.ToLower(exch.GetName()),
			cfg.DataSettings.Interval,
			fPair,
			a)
		if err != nil {
			return nil, fmt.Errorf("%v. Please check your GoCryptoTrader configuration", err)
		}
		resp.Item.RemoveDuplicates()
		resp.Item.SortCandlesByTimestamp(false)
		resp.RangeHolder, err = gctkline.CalculateCandleDateRanges(
			resp.Item.Candles[0].Time,
			resp.Item.Candles[len(resp.Item.Candles)-1].Time.Add(cfg.DataSettings.Interval),
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
	// we want to do this when the live is warming up so it doesnt fail
	case tm.Warmup || !tm.IsLive:
		log.Infof(log.TradeManager, "loading db data for %v %v %v...\n", exch.GetName(), a, fPair)
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
		// tm.Bot.DatabaseManager, err = SetupDatabaseConnectionManager(gctdatabase.DB.GetConfig())
		// if err != nil {
		// 	return nil, err
		// }
		//
		// err = tm.Bot.DatabaseManager.Start(&tm.Bot.ServicesWG)
		// if err != nil {
		// 	return nil, err
		// }
		// defer func() {
		// 	stopErr := tm.Bot.DatabaseManager.Stop()
		// 	if stopErr != nil {
		// 		log.Error(log.TradeManager, stopErr)
		// 	}
		// }()
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
	case tm.IsLive && !tm.Warmup:
		log.Infof(log.TradeManager, "loading live data for %v %v %v...\n", exch.GetName(), a, fPair)
		if len(cfg.CurrencySettings) > 1 {
			return nil, errors.New("live data simulation only supports one currency")
		}

		err = loadLiveData(cfg, b)
		if err != nil {
			return nil, err
		}
		go tm.loadLiveDataLoop(
			resp,
			cfg,
			exch,
			fPair,
			a,
			dataType)
		return resp, nil
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

func loadAPIData(cfg *config.Config, exch gctexchange.IBotExchange, fPair currency.Pair, a asset.Item, resultLimit uint32, dataType int64) (*kline.DataFromKline, error) {
	if cfg.DataSettings.Interval <= 0 {
		return nil, errIntervalUnset
	}
	dates, err := gctkline.CalculateCandleDateRanges(
		cfg.DataSettings.APIData.StartDate,
		cfg.DataSettings.APIData.EndDate,
		gctkline.Interval(cfg.DataSettings.Interval),
		resultLimit)
	if err != nil {
		return nil, err
	}
	candles, err := api.LoadData(context.TODO(),
		dataType,
		cfg.DataSettings.APIData.StartDate,
		cfg.DataSettings.APIData.EndDate,
		cfg.DataSettings.Interval,
		exch,
		fPair,
		a)
	if err != nil {
		return nil, fmt.Errorf("%v. Please check your GoCryptoTrader configuration", err)
	}
	dates.SetHasDataFromCandles(candles.Candles)
	summary := dates.DataSummary(false)
	if len(summary) > 0 {
		log.Warnf(log.TradeManager, "%v", summary)
	}
	candles.FillMissingDataWithEmptyEntries(dates)
	candles.RemoveOutsideRange(cfg.DataSettings.APIData.StartDate, cfg.DataSettings.APIData.EndDate)
	return &kline.DataFromKline{
		Item:        *candles,
		RangeHolder: dates,
	}, nil
}

func loadLiveData(cfg *config.Config, base *gctexchange.Base) error {
	if cfg == nil || base == nil || cfg.DataSettings.LiveData == nil {
		return eventtypes.ErrNilArguments
	}
	if cfg.DataSettings.Interval <= 0 {
		return errIntervalUnset
	}

	if cfg.DataSettings.LiveData.APIKeyOverride != "" {
		base.API.Credentials.Key = cfg.DataSettings.LiveData.APIKeyOverride
	}
	if cfg.DataSettings.LiveData.APISecretOverride != "" {
		base.API.Credentials.Secret = cfg.DataSettings.LiveData.APISecretOverride
	}
	if cfg.DataSettings.LiveData.APIClientIDOverride != "" {
		base.API.Credentials.ClientID = cfg.DataSettings.LiveData.APIClientIDOverride
	}
	if cfg.DataSettings.LiveData.API2FAOverride != "" {
		base.API.Credentials.PEMKey = cfg.DataSettings.LiveData.API2FAOverride
	}
	if cfg.DataSettings.LiveData.APISubAccountOverride != "" {
		base.API.Credentials.Subaccount = cfg.DataSettings.LiveData.APISubAccountOverride
	}
	validated := base.ValidateAPICredentials()
	base.API.AuthenticatedSupport = validated
	if !validated && cfg.DataSettings.LiveData.RealOrders {
		log.Warn(log.TradeManager, "invalid API credentials set, real orders set to false")
		cfg.DataSettings.LiveData.RealOrders = false
	}
	return nil
}

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
	case fill.Event:
		// fmt.Println("fill event")
		tm.processFillEvent(eType)
	default:
		return fmt.Errorf("%w %v received, could not process",
			errUnhandledDatatype,
			ev)
	}

	return nil
}

func (tm *TradeManager) processSingleDataEvent(ev eventtypes.DataEventHandler) error {

	err := tm.updateStatsForDataEvent(ev)
	if err != nil {
		return err
	}

	d := tm.Datas.GetDataForCurrency(ev.GetExchange(), ev.GetAssetType(), ev.Pair())

	// update factor engine
	tm.FactorEngine.OnBar(d)

	// HANDLE warmup MODE
	if !tm.Warmup {
		var s signal.Event
		for _, strategy := range tm.Strategies {
			s, err = strategy.OnData(d, tm.Portfolio, tm.FactorEngine)
			tm.EventQueue.AppendEvent(s)
		}
	}

	if err != nil {
		if errors.Is(err, base.ErrTooMuchBadData) {
			// too much bad data is a severe error and backtesting must cease
			return err
		}
		log.Error(log.TradeManager, err)
		return nil
	}
	err = tm.Statistic.SetEventForOffset(ev)
	if err != nil {
		log.Error(log.TradeManager, err)
	}

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
	var dataEvents []data.Handler
	dataHandlerMap := tm.Datas.GetAllData()
	for _, exchangeMap := range dataHandlerMap {
		for _, assetMap := range exchangeMap {
			for _, dataHandler := range assetMap {
				latestData := dataHandler.Latest()
				err := tm.updateStatsForDataEvent(latestData)
				if err != nil && err == statistics.ErrAlreadyProcessed {
					continue
				}
				dataEvents = append(dataEvents, dataHandler)
			}
		}
	}
	signals, err := tm.Strategies[0].OnSimultaneousSignals(dataEvents, tm.Portfolio, tm.FactorEngine)
	if err != nil {
		if errors.Is(err, base.ErrTooMuchBadData) {
			// too much bad data is a severe error and backtesting must cease
			return err
		}
		log.Error(log.TradeManager, err)
		return nil
	}
	for i := range signals {
		err = tm.Statistic.SetEventForOffset(signals[i])
		if err != nil {
			log.Error(log.TradeManager, err)
		}
		tm.EventQueue.AppendEvent(signals[i])
	}
	return nil
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

// processSignalEvent receives an event from the strategy for processing under the portfolio
func (tm *TradeManager) processSignalEvent(ev signal.Event) {
	cs, err := tm.Exchange.GetCurrencySettings(ev.GetExchange(), ev.GetAssetType(), ev.Pair())
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

func (tm *TradeManager) processOrderEvent(ev order.Event) {
	d := tm.Datas.GetDataForCurrency(ev.GetExchange(), ev.GetAssetType(), ev.Pair())
	f, err := tm.Exchange.ExecuteOrder(ev, d, tm.Bot.OrderManager)
	if err != nil {
		if f == nil {
			log.Errorf(log.TradeManager, "fill event should always be returned, please fix, %v", err)
			return
		}
		log.Errorf(log.TradeManager, "%v %v %v %v", f.GetExchange(), f.GetAssetType(), f.Pair(), err)
	}
	err = tm.Statistic.SetEventForOffset(f)
	if err != nil {
		log.Error(log.TradeManager, err)
	}
	tm.EventQueue.AppendEvent(f)
}

func (tm *TradeManager) processFillEvent(ev fill.Event) {
	_, err := tm.Portfolio.OnFill(ev)
	if err != nil {
		log.Error(log.TradeManager, err)
		return
	}
}

// loadLiveDataLoop is an incomplete function to continuously retrieve exchange data on a loop
// from live. Its purpose is to be able to perform strategy analysis against current data
func (tm *TradeManager) loadLiveDataLoop(resp *kline.DataFromKline, cfg *config.Config, exch gctexchange.IBotExchange, fPair currency.Pair, a asset.Item, dataType int64) {
	startDate := time.Now().Add(-cfg.DataSettings.Interval * 2)
	dates, err := gctkline.CalculateCandleDateRanges(
		startDate,
		startDate.AddDate(1, 0, 0),
		gctkline.Interval(cfg.DataSettings.Interval),
		0)
	if err != nil {
		log.Errorf(log.TradeManager, "%v. Please check your GoCryptoTrader configuration", err)
		return
	}
	candles, err := live.LoadData(context.TODO(),
		exch,
		dataType,
		cfg.DataSettings.Interval,
		fPair,
		a)
	if err != nil {
		log.Errorf(log.TradeManager, "%v. Please check your GoCryptoTrader configuration", err)
		return
	}
	dates.SetHasDataFromCandles(candles.Candles)
	resp.RangeHolder = dates
	resp.Item = *candles

	loadNewDataTimer := time.NewTimer(time.Second * 5)
	for {
		select {
		case <-tm.shutdown:
			return
		case <-loadNewDataTimer.C:
			log.Infof(log.TradeManager, "fetching data for %v %v %v %v", exch.GetName(), a, fPair, cfg.DataSettings.Interval)
			loadNewDataTimer.Reset(time.Second * 15)
			err = tm.loadLiveData(resp, cfg, exch, fPair, a, dataType)
			if err != nil {
				log.Error(log.TradeManager, err)
				return
			}
		}
	}
}

func (tm *TradeManager) loadLiveData(resp *kline.DataFromKline, cfg *config.Config, exch gctexchange.IBotExchange, fPair currency.Pair, a asset.Item, dataType int64) error {
	if resp == nil {
		return errNilData
	}
	if cfg == nil {
		return errNilConfig
	}
	if exch == nil {
		return errNilExchange
	}
	candles, err := live.LoadData(context.TODO(),
		exch,
		dataType,
		cfg.DataSettings.Interval,
		fPair,
		a)
	if err != nil {
		return err
	}
	if len(candles.Candles) == 0 {
		return nil
	}
	resp.AppendResults(candles)
	tm.Reports.UpdateItem(&resp.Item)
	log.Info(log.TradeManager, "sleeping for 30 seconds before checking for new candle data")
	return nil
}

func (b *TradeManager) heartBeat() {
	for range time.Tick(time.Second * 15) {
		log.Info(log.TradeManager, "heartbeat")
	}
}

// IsRunning returns if gctscript manager subsystem is started
func (b *TradeManager) IsRunning() bool {
	if b == nil {
		return false
	}
	return atomic.LoadInt32(&b.started) == 1
}
