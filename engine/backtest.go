package engine

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	config "github.com/thrasher-corp/gocryptotrader/bt_config"
	gctcommon "github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/compliance"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/data"
	"github.com/thrasher-corp/gocryptotrader/data/kline"
	"github.com/thrasher-corp/gocryptotrader/data/kline/api"
	"github.com/thrasher-corp/gocryptotrader/data/kline/csv"
	"github.com/thrasher-corp/gocryptotrader/data/kline/database"
	"github.com/thrasher-corp/gocryptotrader/data/kline/live"
	gctdatabase "github.com/thrasher-corp/gocryptotrader/database"
	"github.com/thrasher-corp/gocryptotrader/eventtypes"
	"github.com/thrasher-corp/gocryptotrader/eventtypes/fill"
	"github.com/thrasher-corp/gocryptotrader/eventtypes/order"
	"github.com/thrasher-corp/gocryptotrader/eventtypes/signal"
	gctexchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	gctkline "github.com/thrasher-corp/gocryptotrader/exchanges/kline"
	gctorder "github.com/thrasher-corp/gocryptotrader/exchanges/order"
	"github.com/thrasher-corp/gocryptotrader/factors"
	"github.com/thrasher-corp/gocryptotrader/log"
	"github.com/thrasher-corp/gocryptotrader/report"
	"github.com/thrasher-corp/gocryptotrader/risk"
	"github.com/thrasher-corp/gocryptotrader/slippage"
	"github.com/thrasher-corp/gocryptotrader/statistics"
	"github.com/thrasher-corp/gocryptotrader/strategies"
	"github.com/thrasher-corp/gocryptotrader/strategies/base"
)

// New returns a new BackTest instance
func NewBacktest() *BackTest {
	return &BackTest{
		shutdown: make(chan struct{}),
	}
}

// Reset BackTest values to default
func (bt *BackTest) Reset() {
	bt.EventQueue.Reset()
	bt.Datas.Reset()
	bt.Portfolio.Reset()
	// bt.Statistic.Reset()
	bt.Exchange.Reset()
	bt.Bot = nil
	bt.FactorEngine = nil
}

// NewFromConfig takes a strategy config and configures a backtester variable to run
func NewFromConfig(cfg *config.Config, templatePath, output string, bot *Engine) (*BackTest, error) {
	log.Infoln(log.BackTester, "bt:loading config...")
	if cfg == nil {
		return nil, errNilConfig
	}
	if bot == nil {
		return nil, errNilBot
	}
	bt := NewBacktest()
	bt.Datas = &data.HandlerPerCurrency{}
	bt.EventQueue = &Holder{}
	reports := &report.Data{
		Config:       cfg,
		TemplatePath: templatePath,
		OutputPath:   output,
	}
	bt.Reports = reports

	err := bt.setupBot(cfg, bot)
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

	useExchangeLevelFunding := cfg.StrategySettings.UseExchangeLevelFunding
	if useExchangeLevelFunding {
		for i := range cfg.StrategySettings.ExchangeLevelFunding {
			_, err = asset.New(cfg.StrategySettings.ExchangeLevelFunding[i].Asset)
			if err != nil {
				return nil, err
			}
		}
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
			log.Warnf(log.BackTester, "maker fee '%v' should not exceed taker fee '%v'. Please review config",
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

	bt.Strategies = slit

	var p *Portfolio
	p, err = SetupPortfolio(bt.Strategies, *bot, sizeManager, portfolioRisk, cfg.StatisticSettings.RiskFreeRate)
	if err != nil {
		return nil, err
	}

	if cfg.StrategySettings.CustomSettings != nil {
		err = bt.Strategies[0].SetCustomSettings(cfg.StrategySettings.CustomSettings)
		if err != nil && !errors.Is(err, base.ErrCustomSettingsUnsupported) {
			return nil, err
		}
	}
	bt.FactorEngine, err = factors.Setup()
	if err != nil {
		return nil, err
	}
	// stats := &statistics.Statistic{
	// 	StrategyName:                s.Name(),
	// 	StrategyNickname:            cfg.Nickname,
	// 	StrategyDescription:         bt.Strategy.Description(),
	// 	StrategyGoal:                cfg.Goal,
	// 	ExchangeAssetPairStatistics: make(map[string]map[asset.Item]map[currency.Pair]*currencystatistics.CurrencyStatistic),
	// 	RiskFreeRate:                cfg.StatisticSettings.RiskFreeRate,
	// }
	// bt.Statistic = stats
	// reports.Statistics = stats

	e, err := bt.setupExchangeSettings(cfg)
	if err != nil {
		return nil, err
	}

	bt.Exchange = &e
	for i := range e.CurrencySettings {
		var lookup *PortfolioExchangeSettings
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

		dataType, _ := eventtypes.DataTypeToInt(cfg.DataSettings.DataType)
		resp, _ := database.LoadData(
			cfg.DataSettings.DatabaseData.StartDate,
			cfg.DataSettings.DatabaseData.EndDate,
			cfg.DataSettings.Interval,
			strings.ToLower(e.CurrencySettings[i].ExchangeName),
			dataType,
			e.CurrencySettings[i].CurrencyPair,
			e.CurrencySettings[i].AssetType)

		lastBar := resp.Item.Candles[len(resp.Item.Candles)-1]

		// how much do we have to catchup?

		// what's the current date?
		fmt.Println("catchup range", time.Now(), lastBar.Time)

		// perform the catchup
		// request the data missing from the broker

		// ask the data history manager to perform the catchup

		// warm the factor engine

		// then start the damn thing
		// strategies are stateless, and the factor engine is warmed up outside of the backtest
		// so we do not need any dataevents to trigger during catchup
		// we only need to make sure the factor engine is updated to the current time/last bar time
		// so that it can respond to queries

		// resp.Item.RemoveDuplicates()
		// resp.Item.SortCandlesByTimestamp(false)
		// resp.RangeHolder, err = gctkline.CalculateCandleDateRanges(
		// 	cfg.DataSettings.DatabaseData.StartDate,
		// 	cfg.DataSettings.DatabaseData.EndDate,
		// 	gctkline.Interval(cfg.DataSettings.Interval),
		// 	0,
		// )
		// if err != nil {
		// 	return nil, err
		// }
		// resp.RangeHolder.SetHasDataFromCandles(resp.Item.Candles)
		// summary := resp.RangeHolder.DataSummary(false)
		// if len(summary) > 0 {
		// 	log.Warnf(log.BackTester, "%v", summary)
		// }

	}
	bt.Portfolio = p

	// cfg.PrintSetting()

	return bt, nil
}

// Run will iterate over loaded data events
// save them and then handle the event based on its type
func (bt *BackTest) Run() error {
	log.Info(log.BackTester, "running backtester against pre-defined data")
dataLoadingIssue:
	for ev := bt.EventQueue.NextEvent(); ; ev = bt.EventQueue.NextEvent() {
		if ev == nil {
			dataHandlerMap := bt.Datas.GetAllData()
			for exchangeName, exchangeMap := range dataHandlerMap {
				for assetItem, assetMap := range exchangeMap {
					// var hasProcessedData bool
					for currencyPair, dataHandler := range assetMap {
						d := dataHandler.Next()
						if d == nil {
							if !bt.hasHandledEvent {
								log.Errorf(log.BackTester, "Unable to perform `Next` for %v %v %v", exchangeName, assetItem, currencyPair)
							}
							break dataLoadingIssue
						}
						bt.EventQueue.AppendEvent(d)
						// hasProcessedData = true
					}
				}
			}
		}
		if ev != nil {
			err := bt.handleEvent(ev)
			if err != nil {
				return err
			}
		}
		if !bt.hasHandledEvent {
			bt.hasHandledEvent = true
		}
	}

	return nil
}

// uses a single database
// load the data from the database
// do catchup before starting
// trades/activeOrderse are always loaded from the database in case system shut down with trades/orders
// data is always written to the database in the background after receiving it
// this will not need more than 60 days of data at the most to operate
// this should be able to handle 50-100+ instruments without a problem
// this database will not save the factors, unlike in backtests
// if the system is shut down, it will have to recalculate the factors from scratch from the data after catchup process is completed
// the first step is to write the data in the database as it comes in
func (bt *BackTest) RunLive() error {
	bt.FactorEngine.Start()
	bt.liveDataCatchup()

	timeoutTimer := time.NewTimer(time.Minute * 1)

	// a frequent timer so that when a new candle is released by an exchange
	// that it can be processed quickly
	processEventTicker := time.NewTicker(time.Second)
	doneARun := false
	for {
		select {
		case <-bt.shutdown:
			return nil
		case <-timeoutTimer.C:
			return errLiveDataTimeout
		case <-processEventTicker.C:
			for ev := bt.EventQueue.NextEvent(); ; ev = bt.EventQueue.NextEvent() {
				if ev == nil {
					dataHandlerMap := bt.Datas.GetAllData()
					for exchangeName, exchangeMap := range dataHandlerMap {
						for assetItem, assetMap := range exchangeMap {
							// var hasProcessedData bool
							for currencyPair, dataHandler := range assetMap {
								d := dataHandler.Next()
								if d == nil {
									if !bt.hasHandledEvent {
										log.Errorf(log.BackTester, "Unable to perform `Next` for %v %v %v", exchangeName, assetItem, currencyPair)
									}
									// break dataLoadingIssue
								}

								bt.EventQueue.AppendEvent(d)
								// hasProcessedData = true
							}
						}
					}
				}
				// if ev == nil {
				// 	// as live only supports singular currency, just get the proper reference manually
				// 	var d data.Handler
				// 	dd := bt.Datas.GetAllData()
				//
				// 	for k1, v1 := range dd {
				// 		for k2, v2 := range v1 {
				// 			for k3 := range v2 {
				// 				d = dd[k1][k2][k3]
				// 			}
				// 		}
				// 	}
				// 	de := d.Next()
				// 	if de == nil {
				// 		break
				// 	}
				//
				// 	bt.EventQueue.AppendEvent(de)
				// 	doneARun = true
				// 	continue
				// }

				if ev != nil {
					err := bt.handleEvent(ev)
					if err != nil {
						return err
					}
				}
				if !bt.hasHandledEvent {
					bt.hasHandledEvent = true
				}
			}
			if doneARun {
				timeoutTimer = time.NewTimer(time.Minute * 5)
			}
		}
	}
}

// Stop shuts down the live data loop
func (bt *BackTest) Stop() {
	for _, s := range bt.Strategies {
		s.Stop()
	}

	defer func() {
		stopErr := bt.Bot.DatabaseManager.Stop()
		if stopErr != nil {
			log.Error(log.BackTester, stopErr)
		}
	}()

	close(bt.shutdown)
}

func (bt *BackTest) liveDataCatchup() error {

	return nil
}

func (bt *BackTest) setupExchangeSettings(cfg *config.Config) (FakeExchange, error) {
	log.Infoln(log.BackTester, "setting exchange settings...")
	resp := FakeExchange{}

	for i := range cfg.CurrencySettings {
		exch, pair, a, err := bt.loadExchangePairAssetBase(
			cfg.CurrencySettings[i].ExchangeName,
			cfg.CurrencySettings[i].Base,
			cfg.CurrencySettings[i].Quote,
			cfg.CurrencySettings[i].Asset)
		if err != nil {
			return resp, err
		}

		exchangeName := strings.ToLower(exch.GetName())
		bt.Datas.Setup()
		klineData, err := bt.loadData(cfg, exch, pair, a)
		if err != nil {
			return resp, err
		}
		bt.Datas.SetDataForCurrency(exchangeName, a, pair, klineData)
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
			log.Warnf(log.BackTester, "invalid maximum slippage percent '%v'. Slippage percent is defined as a number, eg '100.00', defaulting to '%v'",
				cfg.CurrencySettings[i].MaximumSlippagePercent,
				slippage.DefaultMaximumSlippagePercent)
			cfg.CurrencySettings[i].MaximumSlippagePercent = slippage.DefaultMaximumSlippagePercent
		}
		if cfg.CurrencySettings[i].MaximumSlippagePercent.IsZero() {
			cfg.CurrencySettings[i].MaximumSlippagePercent = slippage.DefaultMaximumSlippagePercent
		}
		if cfg.CurrencySettings[i].MinimumSlippagePercent.LessThan(decimal.Zero) {
			log.Warnf(log.BackTester, "invalid minimum slippage percent '%v'. Slippage percent is defined as a number, eg '80.00', defaulting to '%v'",
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
				log.Warnf(log.BackTester, "exchange %s order execution limits supported but disabled for %s %s, live results may differ",
					cfg.CurrencySettings[i].ExchangeName,
					pair,
					a)
				cfg.CurrencySettings[i].ShowExchangeOrderLimitWarning = true
			}
		}

		resp.CurrencySettings = append(resp.CurrencySettings, PortfolioExchangeSettings{
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

func (bt *BackTest) loadExchangePairAssetBase(exch, base, quote, ass string) (gctexchange.IBotExchange, currency.Pair, asset.Item, error) {
	e, err := bt.Bot.GetExchangeByName(exch)
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
		log.Warnf(log.BackTester, "no credentials set for %v, this is theoretical only", exchangeBase.Name)
	}

	fPair, err = exchangeBase.FormatExchangeCurrency(cp, a)
	if err != nil {
		return nil, currency.Pair{}, "", err
	}
	return e, fPair, a, nil
}

// setupBot sets up a basic bot to retrieve exchange data
// as well as process orders
// setup order manager and exchange manager
func (bt *BackTest) setupBot(cfg *config.Config, bot *Engine) error {
	var err error
	bt.Bot = bot
	bt.Bot.ExchangeManager = SetupExchangeManager()

	for i := range cfg.CurrencySettings {
		err = bt.Bot.LoadExchange(cfg.CurrencySettings[i].ExchangeName, nil)
		if err != nil && !errors.Is(err, ErrExchangeAlreadyLoaded) {
			return err
		}
	}

	bt.Bot.DatabaseManager, err = SetupDatabaseConnectionManager(gctdatabase.DB.GetConfig())
	if err != nil {
		return err
	}

	err = bt.Bot.DatabaseManager.Start(&bt.Bot.ServicesWG)
	if err != nil {
		return err
	}

	if cfg.IsLive && !bt.Bot.CommunicationsManager.IsRunning() {
		communicationsConfig := bot.Config.GetCommunicationsConfig()
		bot.CommunicationsManager, err = SetupCommunicationManager(&communicationsConfig)
		if err != nil {
			return err
		}
		bot.CommunicationsManager.Start()
	}

	if !bt.Bot.OrderManager.IsRunning() {
		bt.Bot.OrderManager, err = SetupOrderManager(
			bt.Bot.ExchangeManager,
			bt.Bot.CommunicationsManager,
			&bt.Bot.ServicesWG,
			bot.Settings.Verbose)
		if err != nil {
			return err
		}
		err = bt.Bot.OrderManager.Start()
		if err != nil {
			return err
		}
	}

	// TODO
	// DataCatchup()
	// FactorEngineWarmup()
	// EnsureSyncronized()

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
		log.Errorf(log.BackTester, "Could not retrieve taker fee for %v. %v", exch.GetName(), err)
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
		log.Errorf(log.BackTester, "Could not retrieve maker fee for %v. %v", exch.GetName(), err)
	}

	return decimal.NewFromFloat(fMakerFee), decimal.NewFromFloat(fTakerFee)
}

// loadData will create kline data from the sources defined in start config files. It can exist from databases, csv or API endpoints
// it can also be generated from trade data which will be converted into kline data
func (bt *BackTest) loadData(cfg *config.Config, exch gctexchange.IBotExchange, fPair currency.Pair, a asset.Item) (*kline.DataFromKline, error) {
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
		log.Infof(log.BackTester, "loading csv data for %v %v %v...\n", exch.GetName(), a, fPair)
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
			log.Warnf(log.BackTester, "%v", summary)
		}
	case !cfg.IsLive && cfg.DataSettings.DatabaseData != nil:
		log.Infof(log.BackTester, "loading db data for %v %v %v...\n", exch.GetName(), a, fPair)
		if cfg.DataSettings.DatabaseData.InclusiveEndDate {
			cfg.DataSettings.DatabaseData.EndDate = cfg.DataSettings.DatabaseData.EndDate.Add(cfg.DataSettings.Interval)
		}
		if cfg.DataSettings.DatabaseData.ConfigOverride != nil {
			bt.Bot.Config.Database = *cfg.DataSettings.DatabaseData.ConfigOverride
			gctdatabase.DB.DataPath = filepath.Join(gctcommon.GetDefaultDataDir(runtime.GOOS), "database")
			err = gctdatabase.DB.SetConfig(cfg.DataSettings.DatabaseData.ConfigOverride)
			if err != nil {
				return nil, err
			}
		}
		bt.Bot.DatabaseManager, err = SetupDatabaseConnectionManager(gctdatabase.DB.GetConfig())
		if err != nil {
			return nil, err
		}

		err = bt.Bot.DatabaseManager.Start(&bt.Bot.ServicesWG)
		if err != nil {
			return nil, err
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
			log.Warnf(log.BackTester, "%v", summary)
		}
	case cfg.DataSettings.APIData != nil:
		log.Infof(log.BackTester, "loading api data for %v %v %v...\n", exch.GetName(), a, fPair)
		if cfg.DataSettings.APIData.InclusiveEndDate {
			cfg.DataSettings.APIData.EndDate = cfg.DataSettings.APIData.EndDate.Add(cfg.DataSettings.Interval)
		}
		resp, err = loadAPIData(
			cfg,
			exch,
			fPair,
			a,
			b.Features.Enabled.Kline.ResultLimit,
			dataType)
		if err != nil {
			return resp, err
		}
	case cfg.IsLive && cfg.DataSettings.LiveData != nil:
		if len(cfg.CurrencySettings) > 1 {
			return nil, errors.New("live data simulation only supports one currency")
		}

		err = loadLiveData(cfg, b)
		if err != nil {
			return nil, err
		}
		go bt.loadLiveDataLoop(
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
	bt.Reports.AddKlineItem(&resp.Item)
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
		log.Warnf(log.BackTester, "%v", summary)
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
		log.Warn(log.BackTester, "invalid API credentials set, real orders set to false")
		cfg.DataSettings.LiveData.RealOrders = false
	}
	return nil
}

// handleEvent is the main processor of data for the backtester
// after data has been loaded and Run has appended a data event to the queue,
// handle event will process events and add further events to the queue if they
// are required
func (bt *BackTest) handleEvent(ev eventtypes.EventHandler) error {
	switch eType := ev.(type) {
	case eventtypes.DataEventHandler:
		return bt.processSingleDataEvent(eType)
	case signal.Event:
		bt.processSignalEvent(eType)
	case order.Event:
		bt.processOrderEvent(eType)
	case fill.Event:
		bt.processFillEvent(eType)
	default:
		return fmt.Errorf("%w %v received, could not process",
			errUnhandledDatatype,
			ev)
	}

	return nil
}

func (bt *BackTest) processSingleDataEvent(ev eventtypes.DataEventHandler) error {

	err := bt.updateStatsForDataEvent(ev)
	if err != nil {
		return err
	}

	d := bt.Datas.GetDataForCurrency(ev.GetExchange(), ev.GetAssetType(), ev.Pair())

	// update factor engine
	// fmt.Println(len(d.History()))
	bt.FactorEngine.OnBar(d)

	// HANDLE CATCHUP MODE
	if !bt.catchup {
		var s signal.Event
		for _, strategy := range bt.Strategies {
			s, err = strategy.OnData(d, bt.Portfolio, bt.FactorEngine)
			bt.EventQueue.AppendEvent(s)
		}
	}

	// if err != nil {
	// 	if errors.Is(err, base.ErrTooMuchBadData) {
	// 		// too much bad data is a severe error and backtesting must cease
	// 		return err
	// 	}
	// 	log.Error(log.BackTester, err)
	// 	return nil
	// }
	// err = bt.Statistic.SetEventForOffset(s)
	// if err != nil {
	// 	log.Error(log.BackTester, err)
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
func (bt *BackTest) processSimultaneousDataEvents() error {
	var dataEvents []data.Handler
	dataHandlerMap := bt.Datas.GetAllData()
	for _, exchangeMap := range dataHandlerMap {
		for _, assetMap := range exchangeMap {
			for _, dataHandler := range assetMap {
				latestData := dataHandler.Latest()
				err := bt.updateStatsForDataEvent(latestData)
				if err != nil && err == statistics.ErrAlreadyProcessed {
					continue
				}
				dataEvents = append(dataEvents, dataHandler)
			}
		}
	}
	signals, err := bt.Strategies[0].OnSimultaneousSignals(dataEvents, bt.Portfolio, bt.FactorEngine)
	if err != nil {
		if errors.Is(err, base.ErrTooMuchBadData) {
			// too much bad data is a severe error and backtesting must cease
			return err
		}
		log.Error(log.BackTester, err)
		return nil
	}
	for i := range signals {
		err = bt.Statistic.SetEventForOffset(signals[i])
		if err != nil {
			log.Error(log.BackTester, err)
		}
		bt.EventQueue.AppendEvent(signals[i])
	}
	return nil
}

// updateStatsForDataEvent makes various systems aware of price movements from
// data events
func (bt *BackTest) updateStatsForDataEvent(ev eventtypes.DataEventHandler) error {
	// // update statistics with the latest price
	// err := bt.Statistic.SetupEventForTime(ev)
	// if err != nil {
	// 	if err == statistics.ErrAlreadyProcessed {
	// 		return err
	// 	}
	// 	log.Error(log.BackTester, err)
	// }
	// update portfolio manager with the latest price
	err := bt.Portfolio.UpdateHoldings(ev)
	if err != nil {
		log.Error(log.BackTester, err)
	}
	return nil
}

// processSignalEvent receives an event from the strategy for processing under the portfolio
func (bt *BackTest) processSignalEvent(ev signal.Event) {
	cs, err := bt.Exchange.GetCurrencySettings(ev.GetExchange(), ev.GetAssetType(), ev.Pair())
	if err != nil {
		log.Error(log.BackTester, err)
		return
	}
	var o *order.Order
	o, err = bt.Portfolio.OnSignal(ev, &cs)
	if err != nil {
		log.Error(log.BackTester, err)
		return
	}
	if err != nil {
		log.Error(log.BackTester, err)
	}

	if o != nil {
		// err = bt.Statistic.SetEventForOffset(o)
		bt.EventQueue.AppendEvent(o)
	}
}

func (bt *BackTest) processOrderEvent(ev order.Event) {
	d := bt.Datas.GetDataForCurrency(ev.GetExchange(), ev.GetAssetType(), ev.Pair())
	f, err := bt.Exchange.ExecuteOrder(ev, d, bt.Bot)
	if err != nil {
		if f == nil {
			log.Errorf(log.BackTester, "fill event should always be returned, please fix, %v", err)
			return
		}
		log.Errorf(log.BackTester, "%v %v %v %v", f.GetExchange(), f.GetAssetType(), f.Pair(), err)
	}
	// err = bt.Statistic.SetEventForOffset(f)
	// if err != nil {
	// 	log.Error(log.BackTester, err)
	// }
	bt.EventQueue.AppendEvent(f)
}

func (bt *BackTest) processFillEvent(ev fill.Event) {
	_, err := bt.Portfolio.OnFill(ev)
	if err != nil {
		log.Error(log.BackTester, err)
		return
	}

	// err = bt.Statistic.SetEventForOffset(t)
	// if err != nil {
	// 	log.Error(log.BackTester, err)
	// }

	// var holding *holdings.Holding
	// holding, err = bt.Portfolio.ViewHoldingAtTimePeriod(ev)
	// if err != nil {
	// 	log.Error(log.BackTester, err)
	// }

	// err = bt.Statistic.AddHoldingsForTime(holding)
	// if err != nil {
	// 	log.Error(log.BackTester, err)
	// }

	// var cp *compliance.Manager
	// cp, err = bt.Portfolio.GetComplianceManager(ev.GetExchange(), ev.GetAssetType(), ev.Pair())
	// if err != nil {
	// 	log.Error(log.BackTester, err)
	// }

	// snap := cp.GetLatestSnapshot()
	// err = bt.Statistic.AddComplianceSnapshotForTime(snap, ev)
	// if err != nil {
	// 	log.Error(log.BackTester, err)
	// }
}

// loadLiveDataLoop is an incomplete function to continuously retrieve exchange data on a loop
// from live. Its purpose is to be able to perform strategy analysis against current data
func (bt *BackTest) loadLiveDataLoop(resp *kline.DataFromKline, cfg *config.Config, exch gctexchange.IBotExchange, fPair currency.Pair, a asset.Item, dataType int64) {
	startDate := time.Now().Add(-cfg.DataSettings.Interval * 2)
	dates, err := gctkline.CalculateCandleDateRanges(
		startDate,
		startDate.AddDate(1, 0, 0),
		gctkline.Interval(cfg.DataSettings.Interval),
		0)
	if err != nil {
		log.Errorf(log.BackTester, "%v. Please check your GoCryptoTrader configuration", err)
		return
	}
	candles, err := live.LoadData(context.TODO(),
		exch,
		dataType,
		cfg.DataSettings.Interval,
		fPair,
		a)
	if err != nil {
		log.Errorf(log.BackTester, "%v. Please check your GoCryptoTrader configuration", err)
		return
	}
	dates.SetHasDataFromCandles(candles.Candles)
	resp.RangeHolder = dates
	resp.Item = *candles

	loadNewDataTimer := time.NewTimer(time.Second * 5)
	for {
		select {
		case <-bt.shutdown:
			return
		case <-loadNewDataTimer.C:
			log.Infof(log.BackTester, "fetching data for %v %v %v %v", exch.GetName(), a, fPair, cfg.DataSettings.Interval)
			loadNewDataTimer.Reset(time.Second * 15)
			err = bt.loadLiveData(resp, cfg, exch, fPair, a, dataType)
			if err != nil {
				log.Error(log.BackTester, err)
				return
			}
		}
	}
}

func (bt *BackTest) loadLiveData(resp *kline.DataFromKline, cfg *config.Config, exch gctexchange.IBotExchange, fPair currency.Pair, a asset.Item, dataType int64) error {
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
	bt.Reports.UpdateItem(&resp.Item)
	log.Info(log.BackTester, "sleeping for 30 seconds before checking for new candle data")
	return nil
}
