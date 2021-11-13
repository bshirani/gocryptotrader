package engine

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gocryptotrader/common"
	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/currency/coinmarketcap"
	"gocryptotrader/dispatch"
	"gocryptotrader/exchange"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/exchange/request"
	"gocryptotrader/exchange/trade"
	gctscript "gocryptotrader/gctscript/vm"
	"gocryptotrader/log"
	"gocryptotrader/portfolio/strategies"
	"gocryptotrader/portfolio/withdraw"
	"gocryptotrader/utils"

	"github.com/fatih/color"
)

// overarching type across this code base.
type Engine struct {
	CommunicationsManager   *CommunicationManager
	Config                  *config.Config
	StrategiesConfig        []*config.StrategySetting
	CurrencySettings        []*ExchangeAssetPairSettings
	DatabaseManager         *DatabaseConnectionManager
	DepositAddressManager   *DepositAddressManager
	ExchangeManager         *ExchangeManager
	OrderManager            OrderManagerHandler
	ServicesWG              sync.WaitGroup
	Settings                Settings
	TradeManager            *TradeManager
	WithdrawManager         *WithdrawManager
	apiServer               *apiServerManager
	connectionManager       *connectionManager
	currencyPairSyncer      *syncManager
	currencyStateManager    *CurrencyStateManager
	DataHistoryManager      *DataHistoryManager
	eventManager            *eventManager
	gctScriptManager        *gctscript.GctScriptManager
	ntpManager              *ntpManager
	portfolioManager        *portfolioManager
	Strategies              []strategies.Handler
	uptime                  time.Time
	watcher                 *Watcher
	websocketRoutineManager *websocketRoutineManager
}

// Bot is a happy global engine to allow various areas of the application
// to access its setup services and functions
var Bot *Engine

// New starts a new engine
func New() (*Engine, error) {
	newEngineMutex.Lock()
	defer newEngineMutex.Unlock()
	var b Engine
	b.Config = &config.Cfg

	err := b.Config.LoadConfig("", false)
	if err != nil {
		return nil, fmt.Errorf("failed to load config. Err: %s", err)
	}

	return &b, nil
}

// NewFromSettings starts a new engine based on supplied settings
func NewFromSettings(settings *Settings, flagSet map[string]bool) (*Engine, error) {
	newEngineMutex.Lock()
	defer newEngineMutex.Unlock()
	if settings == nil {
		return nil, errors.New("engine: settings is nil")
	}

	var b Engine
	var err error

	b.Config, err = loadConfigWithSettings(settings, flagSet)
	if err != nil {
		return nil, fmt.Errorf("failed to load config. Err: %s", err)
	}

	// log.Infoln(log.Global, "setting dry run to true for testing")

	if *b.Config.Logging.Enabled {
		log.SetupGlobalLogger()
		log.SetupSubLoggers(b.Config.Logging.SubLoggers)
		log.Debugln(log.Global, "Logger initialised.")
	}

	b.Settings.ConfigFile = settings.ConfigFile
	b.Settings.DataDir = b.Config.GetDataPath()
	b.Settings.CheckParamInteraction = settings.CheckParamInteraction

	err = utils.AdjustGoMaxProcs(settings.GoMaxProcs)
	if err != nil {
		return nil, fmt.Errorf("unable to adjust runtime GOMAXPROCS value. Err: %s", err)
	}

	b.gctScriptManager, err = gctscript.NewManager(&b.Config.GCTScript)
	if err != nil {
		return nil, fmt.Errorf("failed to create script manager. Err: %s", err)
	}

	b.ExchangeManager = SetupExchangeManager()

	validateSettings(&b, settings, flagSet)

	return &b, nil
}

// loadConfigWithSettings creates configuration based on the provided settings
func loadConfigWithSettings(settings *Settings, flagSet map[string]bool) (*config.Config, error) {
	filePath, err := config.GetAndMigrateDefaultPath(settings.ConfigFile)
	if err != nil {
		return nil, err
	}
	//logger hasn't been initizled yet
	// log.Info(log.Global, fmt.Sprintf("Loading config file %s..\n", filePath))
	fmt.Printf("Loading config file %s\n", filePath)

	conf := &config.Config{}
	err = conf.ReadConfigFromFile(filePath, true)
	if err != nil {
		return nil, fmt.Errorf(config.ErrFailureOpeningConfig, filePath, err)
	}

	fmt.Printf("Loading strategies file %s\n", settings.TradeConfigFile)
	ss, err := config.ReadStrategyConfigFromFile(settings.TradeConfigFile)
	if err != nil {
		return nil, fmt.Errorf(config.ErrFailureOpeningConfig, settings.TradeConfigFile, err)
	}
	conf.TradeManager.Strategies = ss
	// fmt.Println("loaded", len(ss), "strategies")
	// for _, s := range ss {
	// 	fmt.Println("pair", s.Pair)
	// }

	for _, ex := range conf.Exchanges {
		if ex.Enabled {
			// get list of pairs from strategy config
			for _, s := range ss {
				ex.CurrencyPairs.EnablePair(asset.Spot, s.Pair)
			}
			// fmt.Println("ENABLED EXXX", ex.Name, pairs)
		}
	}

	// Apply overrides from settings
	if flagSet["datadir"] {
		// warn if dryrun isn't enabled
		// if !settings.EnableDryRun {
		// 	// log.Warn(log.Global, "Command line argument '-datadir' induces dry run mode.")
		// 	fmt.Println("Command line argument '-datadir' induces dry run mode.")
		// }
		// settings.EnableDryRun = true
		conf.DataDirectory = settings.DataDir
	}

	return conf, conf.CheckConfig()
}

// validateSettings validates and sets all bot settings
func validateSettings(b *Engine, s *Settings, flagSet map[string]bool) {
	b.Settings = *s

	b.Settings.EnableDataHistoryManager = (flagSet["datahistory"] && b.Settings.EnableDatabaseManager) || b.Config.DataHistory.Enabled
	b.Settings.EnableTradeManager = (flagSet["trade"] && b.Settings.EnableTradeManager) || b.Config.TradeManager.Enabled
	b.Settings.EnableTrading = (flagSet["strategies"] && b.Settings.EnableTrading) || b.Config.TradeManager.Trading
	b.Settings.EnableOrderManager = flagSet["orders"] && b.Settings.EnableOrderManager || b.Config.OrderManager.Enabled
	b.Settings.EnableExchangeSyncManager = flagSet["sync"] && b.Settings.EnableExchangeSyncManager || b.Config.SyncManager.Enabled
	b.Settings.EnableDryRun = flagSet["dryrun"] && b.Settings.EnableDryRun || b.Config.DryRun
	b.Settings.EnableLiveMode = b.Config.LiveMode
	b.Settings.EnableProductionMode = flagSet["production"] && b.Settings.EnableProductionMode || b.Config.ProductionMode
	b.Settings.EnableCommsRelayer = b.Config.Communications.TelegramConfig.Enabled

	// if b.Settings.EnableTradeManager {
	// 	b.Settings.EnableDataHistoryManager = true
	// }

	b.Settings.EnableCurrencyStateManager = (flagSet["currencystatemanager"] &&
		b.Settings.EnableCurrencyStateManager) ||
		b.Config.CurrencyStateManager.Enabled != nil &&
			*b.Config.CurrencyStateManager.Enabled

	b.Settings.EnableWatcher = (flagSet["watcher"] &&
		b.Settings.EnableWatcher) ||
		b.Config.Watcher.Enabled != nil &&
			*b.Config.Watcher.Enabled

	b.Settings.EnableDataImporter = (flagSet["dataimporter"] &&
		b.Settings.EnableDataImporter) || b.Config.DataImporter.Enabled

	b.Settings.EnableGCTScriptManager = b.Settings.EnableGCTScriptManager &&
		(flagSet["gctscriptmanager"] || b.Config.GCTScript.Enabled)

	if b.Settings.EnableLiveMode && b.Config.ProductionMode {
		b.Settings.EnablePortfolioManager = true //flagSet["account"]
	}

	if b.Settings.EnablePortfolioManager &&
		b.Settings.PortfolioManagerDelay <= 0 {
		b.Settings.PortfolioManagerDelay = PortfolioSleepDelay
	}
	if flagSet["cleardb"] {
		if !b.Config.ProductionMode {
			b.Settings.EnableClearDB = true
		}
	}
	if flagSet["useml"] {
		b.Settings.EnableMachineLearning = true
	}

	if !flagSet["grpc"] {
		b.Settings.EnableGRPC = b.Config.RemoteControl.GRPC.Enabled
	}

	if !flagSet["grpcproxy"] {
		b.Settings.EnableGRPCProxy = b.Config.RemoteControl.GRPC.GRPCProxyEnabled
	}

	if !flagSet["websocketrpc"] {
		b.Settings.EnableWebsocketRPC = b.Config.RemoteControl.WebsocketRPC.Enabled
	}

	if flagSet["maxvirtualmachines"] {
		maxMachines := uint8(b.Settings.MaxVirtualMachines)
		b.gctScriptManager.MaxVirtualMachines = &maxMachines
	}

	if flagSet["withdrawcachesize"] {
		withdraw.CacheSize = b.Settings.WithdrawCacheSize
	}

	if b.Settings.EnableEventManager && b.Settings.EventManagerDelay <= 0 {
		b.Settings.EventManagerDelay = EventSleepDelay
	}

	// Checks if the flag values are different from the defaults
	if b.Settings.MaxHTTPRequestJobsLimit != int(request.DefaultMaxRequestJobs) &&
		b.Settings.MaxHTTPRequestJobsLimit > 0 {
		request.MaxRequestJobs = int32(b.Settings.MaxHTTPRequestJobsLimit)
	}

	if b.Settings.TradeBufferProcessingInterval != trade.DefaultProcessorIntervalTime {
		if b.Settings.TradeBufferProcessingInterval >= time.Second {
			trade.BufferProcessorIntervalTime = b.Settings.TradeBufferProcessingInterval
		} else {
			b.Settings.TradeBufferProcessingInterval = trade.DefaultProcessorIntervalTime
			// log.Warnf(log.Global, "-tradeprocessinginterval must be >= to 1 second, using default value of %v",
			// 	trade.DefaultProcessorIntervalTime)
		}
	}

	if b.Settings.RequestMaxRetryAttempts != request.DefaultMaxRetryAttempts &&
		b.Settings.RequestMaxRetryAttempts > 0 {
		request.MaxRetryAttempts = b.Settings.RequestMaxRetryAttempts
	}

	if b.Settings.HTTPTimeout <= 0 {
		b.Settings.HTTPTimeout = b.Config.GlobalHTTPTimeout
	}

	if b.Settings.GlobalHTTPTimeout <= 0 {
		b.Settings.GlobalHTTPTimeout = b.Config.GlobalHTTPTimeout
	}

	err := common.SetHTTPClientWithTimeout(b.Settings.GlobalHTTPTimeout)
	if err != nil {
		log.Errorf(log.Global,
			"Could not set new HTTP Client with timeout %s error: %v",
			b.Settings.GlobalHTTPTimeout,
			err)
	}

	if b.Settings.GlobalHTTPUserAgent != "" {
		err = common.SetHTTPUserAgent(b.Settings.GlobalHTTPUserAgent)
		if err != nil {
			log.Errorf(log.Global, "Could not set HTTP User Agent for %s error: %v",
				b.Settings.GlobalHTTPUserAgent,
				err)
		}
	}
}

func setColor(value bool) {
	if value {
		color.Set(color.FgGreen, color.Bold)
	} else {
		color.Set(color.FgRed, color.Bold)
	}
}

func elog(str string, args ...interface{}) {
	if len(args) > 0 {
		switch args[0].(type) {
		case bool:
			defer color.Unset()
			setColor(args[0].(bool))
		}
	}
	str = fmt.Sprintf("%s%s", str, "\n")
	fmt.Printf(str, args...)
}

// PrintSettings returns the engine settings
func PrintSettings(s *Settings) {
	if s.EnableProductionMode {
		elog("\t $$$$$$$$$$$$$$$$$$$$$$$$$$$$$$ REAL MONIES $$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$: %v", s.EnableProductionMode)
	}
	elog("\t live:%v", s.EnableLiveMode)
	elog("\t save_db: %v", !s.EnableDryRun)
	elog("\t trader: %v", s.EnableTradeManager)
	elog("\t trading: %v", s.EnableTrading)
	elog("\t sync: %v kline:%v ticker:%v trade:%v wsTimeout:%v", s.EnableExchangeSyncManager, s.EnableKlineSyncing, s.EnableTickerSyncing, s.EnableTradeSyncing, s.SyncTimeoutWebsocket)

	// elog("\t verbose: %v", s.Verbose)
	// elog("\t data history: %v", s.EnableDataHistoryManager)
	// elog("\t coinmarketcap analaysis: %v", s.EnableCoinmarketcapAnalysis)
	// elog("\t gPRC: %v", s.EnableGRPC)
	// elog("\t database: %v", s.EnableDatabaseManager)
	// elog("\t watcher: %v", s.EnableWatcher)
	// elog("\t comms relayer: %v", s.EnableCommsRelayer)
	// elog("\t event manager: %v", s.EnableEventManager)
	// elog("\t websocket RPC: %v", s.EnableWebsocketRPC)
	// elog("\t websocket routine: %v", s.EnableWebsocketRoutine)
	// elog("\t Enable orderbook syncing: %v\n", s.EnableOrderbookSyncing)
	// elog("\t TM Verbose: %v", s.TradeManager.Verbose)
	// elog("\t Enable all exchanges: %v", s.EnableAllExchanges)
	// elog("\t Enable all pairs: %v", s.EnableAllPairs)
	// elog("\t Enable portfolio manager: %v", s.EnablePortfolioManager)
	// elog("\t Enable currency state manager: %v", s.EnableCurrencyStateManager)
	// elog("\t Portfolio manager sleep delay: %v\n", s.PortfolioManagerDelay)
	// elog("\t Enable gRPC Proxy: %v", s.EnableGRPCProxy)
	// elog("\t Event manager sleep delay: %v", s.EventManagerDelay)
	// elog("\t Enable deposit address manager: %v\n", s.EnableDepositAddressManager)
	// elog("\t Enable NTP client: %v", s.EnableNTPClient)
	// elog("\t Enable dispatcher: %v", s.EnableDispatcher)
	// elog("\t Dispatch package max worker amount: %d", s.DispatchMaxWorkerAmount)
	// elog("\t Dispatch package jobs limit: %d", s.DispatchJobsLimit)
	// elog("\t Exchange sync continuously: %v\n", s.SyncContinuously)
	// elog("\t Exchange sync workers: %v\n", s.SyncWorkers)
	// elog("\t Exchange REST sync timeout: %v\n", s.SyncTimeoutREST)
	// elog("- FOREX SETTINGS:")
	// elog("\t Enable currency conveter: %v", s.EnableCurrencyConverter)
	// elog("\t Enable currency layer: %v", s.EnableCurrencyLayer)
	// elog("\t Enable fixer: %v", s.EnableFixer)
	// elog("\t Enable OpenExchangeRates: %v", s.EnableOpenExchangeRates)
	// elog("\t Enable ExchangeRateHost: %v", s.EnableExchangeRateHost)
	// elog("- EXCHANGE SETTINGS:")
	// elog("\t Enable exchange auto pair updates: %v", s.EnableExchangeAutoPairUpdates)
	// elog("\t Disable all exchange auto pair updates: %v", s.DisableExchangeAutoPairUpdates)
	// elog("\t Enable exchange websocket support: %v", s.EnableExchangeWebsocketSupport)
	// elog("\t Enable exchange verbose mode: %v", s.EnableExchangeVerbose)
	// elog("\t Enable exchange HTTP rate limiter: %v", s.EnableExchangeHTTPRateLimiter)
	// elog("\t Enable exchange HTTP debugging: %v", s.EnableExchangeHTTPDebugging)
	// elog("\t Max HTTP request jobs: %v", s.MaxHTTPRequestJobsLimit)
	// elog("\t HTTP request max retry attempts: %v", s.RequestMaxRetryAttempts)
	// elog("\t Trade buffer processing interval: %v", s.TradeBufferProcessingInterval)
	// elog("\t HTTP timeout: %v", s.HTTPTimeout)
	// elog("\t HTTP user agent: %v", s.HTTPUserAgent)
	// elog("- GCTSCRIPT SETTINGS: ")
	// elog("\t Enable GCTScript manager: %v", s.EnableGCTScriptManager)
	// elog("\t GCTScript max virtual machines: %v", s.MaxVirtualMachines)
	// elog("- WITHDRAW SETTINGS: ")
	// elog("\t Withdraw Cache size: %v", s.WithdrawCacheSize)
	// elog("- COMMON SETTINGS:")
	// elog("\t Global HTTP timeout: %v", s.GlobalHTTPTimeout)
	// elog("\t Global HTTP user agent: %v", s.GlobalHTTPUserAgent)
	// elog("\t Global HTTP proxy: %v", s.GlobalHTTPProxy)
}

// Start starts the engine
func (bot *Engine) Start() error {
	if bot == nil {
		return errors.New("engine instance is nil")
	}
	var err error
	newEngineMutex.Lock()
	defer newEngineMutex.Unlock()

	if bot.Settings.EnableDatabaseManager {
		bot.DatabaseManager, err = SetupDatabaseConnectionManager(&bot.Config.Database)
		if err != nil {
			log.Errorf(log.Global, "Database manager unable to setup: %v", err)
		} else {
			err = bot.DatabaseManager.Start(&bot.ServicesWG)
			if err != nil {
				log.Errorf(log.Global, "Database manager unable to start: %v", err)
			}
		}
	}

	if bot.Settings.EnableDispatcher {
		if err = dispatch.Start(bot.Settings.DispatchMaxWorkerAmount, bot.Settings.DispatchJobsLimit); err != nil {
			log.Errorf(log.DispatchMgr, "Dispatcher unable to start: %v", err)
		}
	}

	// Sets up internet connectivity monitor
	if bot.Settings.EnableConnectivityMonitor {
		bot.connectionManager, err = setupConnectionManager(&bot.Config.ConnectionMonitor)
		if err != nil {
			log.Errorf(log.Global, "Connection manager unable to setup: %v", err)
		} else {
			err = bot.connectionManager.Start()
			if err != nil {
				log.Errorf(log.Global, "Connection manager unable to start: %v", err)
			}
		}
	}

	if bot.Settings.EnableNTPClient {
		if bot.Config.NTPClient.Level == 0 {
			var responseMessage string
			responseMessage, err = bot.Config.SetNTPCheck(os.Stdin)
			if err != nil {
				return fmt.Errorf("unable to set NTP check: %w", err)
			}
			log.Info(log.TimeMgr, responseMessage)
		}
		bot.ntpManager, err = setupNTPManager(&bot.Config.NTPClient, *bot.Config.Logging.Enabled)
		if err != nil {
			log.Errorf(log.Global, "NTP manager unable to start: %s", err)
		}
	}

	bot.uptime = time.Now()
	log.Debugf(log.Global, "Using data dir: %s\n", bot.Settings.DataDir)
	if *bot.Config.Logging.Enabled && strings.Contains(bot.Config.Logging.Output, "file") {
		log.Debugf(log.Global, "Using log file: %s\n",
			filepath.Join(log.LogPath, bot.Config.Logging.LoggerFileConfig.FileName))
	}

	// log.Debugf(log.Global,
	// 	"Using %d out of %d logical processors for runtime performance\n",
	// 	runtime.GOMAXPROCS(-1), runtime.NumCPU())

	// enabledExchanges := bot.Config.CountEnabledExchanges()
	// if bot.Settings.EnableAllExchanges {
	// 	enabledExchanges = len(bot.Config.Exchanges)
	// }
	// log.Debugln(log.Global, "EXCHANGE COVERAGE")
	// log.Debugf(log.Global, "\t Available Exchanges: %d. Enabled Exchanges: %d.\n",
	// 	len(bot.Config.Exchanges), enabledExchanges)

	if bot.Settings.ExchangePurgeCredentials {
		log.Debugln(log.Global, "Purging exchange API credentials.")
		bot.Config.PurgeExchangeAPICredentials()
	}

	log.Debugln(log.Global, "Setting up exchanges..")
	err = bot.SetupExchanges()
	if err != nil {
		return err
	}

	if bot.Settings.EnableCommsRelayer {
		bot.CommunicationsManager, err = SetupCommunicationManager(&bot.Config.Communications)
		if err != nil {
			log.Errorf(log.Global, "Communications manager unable to setup: %s", err)
		} else {
			err = bot.CommunicationsManager.Start()
			if err != nil {
				log.Errorf(log.Global, "Communications manager unable to start: %s", err)
			}
		}
	}

	if bot.Settings.EnableCoinmarketcapAnalysis ||
		bot.Settings.EnableCurrencyConverter ||
		bot.Settings.EnableCurrencyLayer ||
		bot.Settings.EnableFixer ||
		bot.Settings.EnableOpenExchangeRates ||
		bot.Settings.EnableExchangeRateHost {
		os.Exit(123)
		err = currency.RunStorageUpdater(currency.BotOverrides{
			Coinmarketcap:       bot.Settings.EnableCoinmarketcapAnalysis,
			FxCurrencyConverter: bot.Settings.EnableCurrencyConverter,
			FxCurrencyLayer:     bot.Settings.EnableCurrencyLayer,
			FxFixer:             bot.Settings.EnableFixer,
			FxOpenExchangeRates: bot.Settings.EnableOpenExchangeRates,
			FxExchangeRateHost:  bot.Settings.EnableExchangeRateHost,
		},
			&currency.MainConfiguration{
				ForexProviders:         bot.Config.GetForexProviders(),
				CryptocurrencyProvider: coinmarketcap.Settings(bot.Config.Currency.CryptocurrencyProvider),
				Cryptocurrencies:       bot.Config.Currency.Cryptocurrencies,
				FiatDisplayCurrency:    bot.Config.Currency.FiatDisplayCurrency,
				CurrencyDelay:          bot.Config.Currency.CurrencyFileUpdateDuration,
				FxRateDelay:            bot.Config.Currency.ForeignExchangeUpdateDuration,
			},
			bot.Settings.DataDir)
		if err != nil {
			log.Errorf(log.Global, "ExchangeSettings updater system failed to start %s", err)
		}
	}

	if bot.Settings.EnableGRPC {
		go StartRPCServer(bot)
	}

	if bot.Settings.EnablePortfolioManager {
		if bot.portfolioManager == nil {
			bot.portfolioManager, err = setupPortfolioManager(
				bot.ExchangeManager,
				bot.Settings.PortfolioManagerDelay,
				&bot.Config.Portfolio,
				bot.CommunicationsManager)
			if err != nil {
				log.Errorf(log.Global, "portfolio manager unable to setup: %s", err)
			} else {
				err = bot.portfolioManager.Start(&bot.ServicesWG)
				if err != nil {
					log.Errorf(log.Global, "portfolio manager unable to start: %s", err)
				}
			}
		}
	}

	bot.WithdrawManager, err = SetupWithdrawManager(bot.ExchangeManager, bot.portfolioManager, true)
	if err != nil {
		return err
	}

	if bot.Settings.EnableWebsocketRPC {
		var filePath string
		filePath, err = config.GetAndMigrateDefaultPath(bot.Settings.ConfigFile)
		if err != nil {
			return err
		}
		bot.apiServer, err = setupAPIServerManager(&bot.Config.RemoteControl, &bot.Config.Profiler, bot.ExchangeManager, bot, bot.portfolioManager, filePath)
		if err != nil {
			log.Errorf(log.Global, "API Server unable to start: %s", err)
		} else {
			if bot.Settings.EnableWebsocketRPC {
				err = bot.apiServer.StartWebsocketServer()
				if err != nil {
					log.Errorf(log.Global, "could not start websocket API server: %s", err)
				}
			}
		}
	}

	if bot.Settings.EnableDepositAddressManager {
		bot.DepositAddressManager = SetupDepositAddressManager()
		go func() {
			err = bot.DepositAddressManager.Sync(bot.GetExchangeCryptocurrencyDepositAddresses())
			if err != nil {
				log.Errorf(log.Global, "Deposit address manager unable to setup: %s", err)
			}
		}()
	}

	if bot.Settings.EnableOrderManager && bot.Config.ProductionMode {
		log.Warn(log.TradeMgr, "Enabling REAL order manager")
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
			log.Errorf(log.Global, "Order manager unable to setup: %s", err)
		} else {
			err = bot.OrderManager.Start()
			if err != nil {
				log.Errorf(log.Global, "Order manager unable to start: %s", err)
			}
		}
	}

	if bot.Settings.EnableExchangeSyncManager {
		exchangeSyncCfg := &Config{
			SyncTicker:           bot.Config.SyncManager.TickerSync,
			SyncOrderbook:        bot.Config.SyncManager.OrderBookSync,
			SyncKlines:           bot.Config.SyncManager.KlineSync,
			SyncTrades:           bot.Config.SyncManager.TradeSync,
			SyncContinuously:     bot.Config.SyncManager.SyncContinuously,
			NumWorkers:           bot.Config.SyncManager.SyncWorkers,
			Verbose:              bot.Config.SyncManager.Verbose,
			SyncTimeoutREST:      bot.Config.SyncManager.SyncTimeoutREST,
			SyncTimeoutWebsocket: bot.Config.SyncManager.SyncTimeoutWebsocket,
		}

		bot.currencyPairSyncer, err = setupSyncManager(
			exchangeSyncCfg,
			bot.ExchangeManager,
			&bot.Config.RemoteControl,
			bot.Settings.EnableWebsocketRoutine)
		if err != nil {
			log.Errorf(log.Global, "Unable to initialise exchange currency pair syncer. Err: %s", err)
		} else {
			go func() {
				err = bot.currencyPairSyncer.Start()
				if err != nil {
					log.Errorf(log.Global, "failed to start exchange currency pair manager. Err: %s", err)
				}
			}()
		}
	}

	if bot.Settings.EnableEventManager {
		bot.eventManager, err = setupEventManager(bot.CommunicationsManager, bot.ExchangeManager, bot.Settings.EventManagerDelay, bot.Settings.Verbose)
		if err != nil {
			log.Errorf(log.Global, "Unable to initialise event manager. Err: %s", err)
		} else {
			err = bot.eventManager.Start()
			if err != nil {
				log.Errorf(log.Global, "failed to start event manager. Err: %s", err)
			}
		}
	}

	if bot.Settings.EnableWebsocketRoutine {
		bot.websocketRoutineManager, err = setupWebsocketRoutineManager(bot.ExchangeManager, bot.OrderManager, bot.currencyPairSyncer, &bot.Config.Currency, bot.Settings.Verbose)
		if err != nil {
			log.Errorf(log.Global, "Unable to initialise websocket routine manager. Err: %s", err)
		} else {
			err = bot.websocketRoutineManager.Start()
			if err != nil {
				log.Errorf(log.Global, "failed to start websocket routine manager. Err: %s", err)
			}
		}
	}

	if bot.Settings.EnableGCTScriptManager {
		bot.gctScriptManager, err = gctscript.NewManager(&bot.Config.GCTScript)
		if err != nil {
			log.Errorf(log.Global, "failed to create script manager. Err: %s", err)
		}
		if err = bot.gctScriptManager.Start(&bot.ServicesWG); err != nil {
			log.Errorf(log.Global, "GCTScript manager unable to start: %s", err)
		}
	}

	if bot.Settings.EnableCurrencyStateManager {
		bot.currencyStateManager, err = SetupCurrencyStateManager(
			bot.Config.CurrencyStateManager.Delay,
			bot.ExchangeManager)
		if err != nil {
			log.Errorf(log.Global,
				"%s unable to setup: %s",
				CurrencyStateManagementName,
				err)
		} else {
			err = bot.currencyStateManager.Start()
			if err != nil {
				log.Errorf(log.Global,
					"%s unable to start: %s",
					CurrencyStateManagementName,
					err)
			}
		}
	}

	err = bot.SetupExchangeSettings()
	if err != nil {
		fmt.Println("error setting up exchange settings", bot.Config, err)
		return err
	}
	if len(bot.CurrencySettings) < 1 {
		fmt.Println("no currencie settings, exiting")
		os.Exit(2)
	}

	if bot.Settings.EnableDataHistoryManager {
		if bot.DataHistoryManager == nil {
			bot.DataHistoryManager, err = SetupDataHistoryManager(bot, bot.ExchangeManager, bot.DatabaseManager, &bot.Config.DataHistory)
			if err != nil {
				log.Errorf(log.Global, "database history manager unable to setup: %s", err)
			} else {
				err = bot.DataHistoryManager.Start()
				if err != nil {
					log.Errorf(log.Global, "database history manager unable to start: %s", err)
				}
			}
		}
	}

	if bot.Settings.EnableTradeManager {
		if bot.TradeManager == nil {
			bot.TradeManager, err = NewTradeManager(bot)
			if err != nil {
				fmt.Printf("Could not setup trade manager from config. Error: %v.\n", err)
				os.Exit(1)
			} else {
				bot.TradeManager.Start()
			}
		}

		// handles data
		if bot.Settings.EnableWatcher {
			bot.watcher, err = SetupWatcher(
				bot.Config.Watcher.Delay,
				bot,
				bot.TradeManager)
			if err != nil {
				log.Errorf(log.Global,
					"%s unable to setup: %s",
					WatcherName,
					err)
			} else {
				err = bot.watcher.Start()
				if err != nil {
					log.Errorf(log.Global,
						"%s unable to start: %s",
						WatcherName,
						err)
				}
			}
		}
	}

	// catchup data history to database
	// can move this to trade manager setup
	// end check
	log.Debugf(log.Global, "Engine '%s' started.\n", bot.Config.Name)

	// handle script here

	if bot.Settings.EnableDataImporter {
		panic("HERE")
		dataImporter := SetupDataImporter(bot, &bot.Config.DataImporter)
		dataImporter.Run("kucoin")
	}

	if bot.Config.Script {
		go func() {
			// fmt.Println("is script, stopping")
			time.Sleep(time.Second)
			bot.Stop()
			os.Exit(0)
		}()
	}

	return nil
}

// Stop correctly shuts down engine saving configuration files
func (bot *Engine) Stop() {
	newEngineMutex.Lock()
	defer newEngineMutex.Unlock()

	log.Debugln(log.Global, "Engine shutting down..")

	if len(bot.portfolioManager.GetAddresses()) != 0 {
		bot.Config.Portfolio = *bot.portfolioManager.GetPortfolio()
	}

	if bot.gctScriptManager.IsRunning() {
		if err := bot.gctScriptManager.Stop(); err != nil {
			log.Errorf(log.Global, "GCTScript manager unable to stop. Error: %v", err)
		}
	}
	if bot.OrderManager != nil {
		if bot.OrderManager.IsRunning() {
			if err := bot.OrderManager.Stop(); err != nil {
				log.Errorf(log.Global, "Order manager unable to stop. Error: %v", err)
			}
		}
	}
	if bot.eventManager.IsRunning() {
		if err := bot.eventManager.Stop(); err != nil {
			log.Errorf(log.Global, "event manager unable to stop. Error: %v", err)
		}
	}
	if bot.ntpManager.IsRunning() {
		if err := bot.ntpManager.Stop(); err != nil {
			log.Errorf(log.Global, "NTP manager unable to stop. Error: %v", err)
		}
	}
	if bot.CommunicationsManager.IsRunning() {
		if err := bot.CommunicationsManager.Stop(); err != nil {
			log.Errorf(log.Global, "Communication manager unable to stop. Error: %v", err)
		}
	}
	if bot.portfolioManager.IsRunning() {
		if err := bot.portfolioManager.Stop(); err != nil {
			log.Errorf(log.Global, "Fund manager unable to stop. Error: %v", err)
		}
	}
	if bot.connectionManager.IsRunning() {
		if err := bot.connectionManager.Stop(); err != nil {
			log.Errorf(log.Global, "Connection manager unable to stop. Error: %v", err)
		}
	}
	if bot.apiServer.IsRESTServerRunning() {
		if err := bot.apiServer.StopRESTServer(); err != nil {
			log.Errorf(log.Global, "API Server unable to stop REST server. Error: %s", err)
		}
	}
	if bot.apiServer.IsWebsocketServerRunning() {
		if err := bot.apiServer.StopWebsocketServer(); err != nil {
			log.Errorf(log.Global, "API Server unable to stop websocket server. Error: %s", err)
		}
	}
	if bot.DataHistoryManager.IsRunning() {
		if err := bot.DataHistoryManager.Stop(); err != nil {
			log.Errorf(log.DataHistory, "data history manager unable to stop. Error: %v", err)
		}
	}
	if bot.watcher.IsRunning() {
		if err := bot.watcher.Stop(); err != nil {
			log.Errorf(log.Global,
				"watcher unable to stop. Error: %v",
				err)
		}
	}
	if bot.TradeManager.IsRunning() {
		if err := bot.TradeManager.Stop(); err != nil {
			log.Errorf(log.Global, "bt unable to stop. Error: %v", err)
		}
	}
	if bot.DatabaseManager.IsRunning() {
		if err := bot.DatabaseManager.Stop(); err != nil {
			log.Errorf(log.Global, "Database manager unable to stop. Error: %v", err)
		}
	}
	if dispatch.IsRunning() {
		if err := dispatch.Stop(); err != nil {
			log.Errorf(log.DispatchMgr, "Dispatch system unable to stop. Error: %v", err)
		}
	}
	if bot.websocketRoutineManager.IsRunning() {
		if err := bot.websocketRoutineManager.Stop(); err != nil {
			log.Errorf(log.Global, "websocket routine manager unable to stop. Error: %v", err)
		}
	}
	if bot.currencyStateManager.IsRunning() {
		if err := bot.currencyStateManager.Stop(); err != nil {
			log.Errorf(log.Global,
				"currency state manager unable to stop. Error: %v",
				err)
		}
	}

	if bot.Settings.EnableCoinmarketcapAnalysis ||
		bot.Settings.EnableCurrencyConverter ||
		bot.Settings.EnableCurrencyLayer ||
		bot.Settings.EnableFixer ||
		bot.Settings.EnableOpenExchangeRates ||
		bot.Settings.EnableExchangeRateHost {
		if err := currency.ShutdownStorageUpdater(); err != nil {
			log.Errorf(log.Global, "ExchangeSettings storage system. Error: %v", err)
		}
	}

	err := bot.Config.SaveConfigToFile(bot.Settings.ConfigFile)
	if err != nil {
		log.Errorln(log.Global, "Unable to save config.")
	} else {
		log.Debugln(log.Global, "Config file saved successfully.")
	}

	// Wait for services to gracefully shutdown
	bot.ServicesWG.Wait()
	err = log.CloseLogger()
	if err != nil {
		log.Errorln(log.TradeMgr, "Failed to close logger. Error: %v\n", err)
	}
}

// GetExchangeByName returns an exchange given an exchange name
func (bot *Engine) GetExchangeByName(exchName string) (exchange.IBotExchange, error) {
	return bot.ExchangeManager.GetExchangeByName(exchName)
}

// UnloadExchange unloads an exchange by name
func (bot *Engine) UnloadExchange(exchName string) error {
	exchCfg, err := bot.Config.GetExchangeConfig(exchName)
	if err != nil {
		return err
	}

	err = bot.ExchangeManager.RemoveExchange(exchName)
	if err != nil {
		return err
	}

	exchCfg.Enabled = false
	return nil
}

// GetExchanges retrieves the loaded exchanges
func (bot *Engine) GetExchanges() []exchange.IBotExchange {
	exch, err := bot.ExchangeManager.GetExchanges()
	if err != nil {
		log.Warnf(log.ExchangeSys, "Cannot get exchanges: %v", err)
		return []exchange.IBotExchange{}
	}
	return exch
}

// LoadExchange loads an exchange by name. Optional wait group can be added for
// external synchronization.
func (bot *Engine) LoadExchange(name string, wg *sync.WaitGroup) error {
	exch, err := bot.ExchangeManager.NewExchangeByName(name)
	if err != nil {
		return err
	}
	if exch.GetBase() == nil {
		return ErrExchangeFailedToLoad
	}

	var localWG sync.WaitGroup
	localWG.Add(1)
	go func() {
		exch.SetDefaults()
		localWG.Done()
	}()

	exchCfg, err := bot.Config.GetExchangeConfig(name)
	if err != nil {
		fmt.Println("returning error")
		return err
	}

	if bot.Settings.EnableAllPairs &&
		exchCfg.CurrencyPairs != nil {
		assets := exchCfg.CurrencyPairs.GetAssetTypes(false)
		for x := range assets {
			var pairs currency.Pairs
			pairs, err = exchCfg.CurrencyPairs.GetPairs(assets[x], false)
			if err != nil {
				return err
			}
			exchCfg.CurrencyPairs.StorePairs(assets[x], pairs, true)
		}
	}

	if bot.Settings.EnableExchangeVerbose {
		exchCfg.Verbose = true
	}
	if exchCfg.Features != nil {
		if bot.Settings.EnableExchangeWebsocketSupport &&
			exchCfg.Features.Supports.Websocket {
			exchCfg.Features.Enabled.Websocket = true
		}
		if bot.Settings.EnableExchangeAutoPairUpdates &&
			exchCfg.Features.Supports.RESTCapabilities.AutoPairUpdates {
			exchCfg.Features.Enabled.AutoPairUpdates = true
		}
		if bot.Settings.DisableExchangeAutoPairUpdates {
			if exchCfg.Features.Supports.RESTCapabilities.AutoPairUpdates {
				exchCfg.Features.Enabled.AutoPairUpdates = false
			}
		}
	}
	if bot.Settings.HTTPUserAgent != "" {
		exchCfg.HTTPUserAgent = bot.Settings.HTTPUserAgent
	}
	if bot.Settings.HTTPProxy != "" {
		exchCfg.ProxyAddress = bot.Settings.HTTPProxy
	}
	if bot.Settings.HTTPTimeout != exchange.DefaultHTTPTimeout {
		exchCfg.HTTPTimeout = bot.Settings.HTTPTimeout
	}
	if bot.Settings.EnableExchangeHTTPDebugging {
		exchCfg.HTTPDebugging = bot.Settings.EnableExchangeHTTPDebugging
	}

	localWG.Wait()
	if !bot.Settings.EnableExchangeHTTPRateLimiter {
		log.Warnf(log.ExchangeSys,
			"Loaded exchange %s rate limiting has been turned off.\n",
			exch.GetName(),
		)
		err = exch.DisableRateLimiter()
		if err != nil {
			log.Errorf(log.ExchangeSys,
				"Loaded exchange %s rate limiting cannot be turned off: %s.\n",
				exch.GetName(),
				err,
			)
		}
	}

	exchCfg.Enabled = true
	err = exch.Setup(exchCfg)
	if err != nil {
		exchCfg.Enabled = false
		fmt.Println("error", err)
		return err
	}

	// log.Infoln(log.Global, "validating credentials")
	bot.ExchangeManager.Add(exch)
	base := exch.GetBase()

	if base.API.AuthenticatedSupport ||
		base.API.AuthenticatedWebsocketSupport {
		assetTypes := base.GetAssetTypes(false)
		var useAsset asset.Item
		for a := range assetTypes {
			err = base.CurrencyPairs.IsAssetEnabled(assetTypes[a])
			if err != nil {
				continue
			}
			useAsset = assetTypes[a]
			break
		}
		err = exch.ValidateCredentials(context.TODO(), useAsset)
		if err != nil {
			log.Warnf(log.ExchangeSys,
				"%s: Cannot validate credentials, authenticated support has been disabled, Error: %s\n",
				base.Name,
				err)
			base.API.AuthenticatedSupport = false
			base.API.AuthenticatedWebsocketSupport = false
			exchCfg.API.AuthenticatedSupport = false
			exchCfg.API.AuthenticatedWebsocketSupport = false
		}
	}

	if wg != nil {
		// elog("starting exchange...")
		if bot.Settings.EnableLiveMode {
			exch.Start(wg)
		}
	} else {
		tempWG := sync.WaitGroup{}
		exch.Start(&tempWG)
		tempWG.Wait()
	}
	return nil
}

// SetupExchanges sets up the exchanges used by the Bot
func (bot *Engine) SetupExchanges() error {
	var wg sync.WaitGroup
	configs := bot.Config.GetAllExchangeConfigs()

	for x := range configs {
		if !configs[x].Enabled && !bot.Settings.EnableAllExchanges {
			// log.Debugf(log.ExchangeSys, "%s: Exchange support: Disabled\n", configs[x].Name)
			continue
		}
		wg.Add(1)

		// fmt.Println("enabling config", configs[x])

		// enablePair(exMgr, pair)
		// func enablePair(exMgr iExchangeManager, pair currency.Pair) {
		// 	exs, _ := exMgr.GetExchanges()
		// 	// fmt.Println("loaded", len(exs), "exchanges")
		// 	if len(exs) == 0 {
		// 		panic(123)
		// 	}
		//
		// 	// err = c.SetPairs(exchName, assetTypes[0], true, currency.Pairs{newPair})
		// 	for _, ex := range exs {
		// 		fmt.Println("~~~~~~~~~~enabling", pair)
		// 		ex.SetPairs(currency.Pairs{pair}, asset.Spot, true)
		// 	}
		// }
		//

		go func(c config.ExchangeConfig) {
			defer wg.Done()
			err := bot.LoadExchange(c.Name, &wg)
			if err != nil {
				log.Errorf(log.ExchangeSys, "LoadExchange %s failed: %s\n", c.Name, err)
				panic("failed to load exchange")
				return
			}
			log.Infof(log.ExchangeSys,
				"%s: Exchange support: Enabled (Authenticated API support: %s - Verbose mode: %s).\n",
				c.Name,
				common.IsEnabled(c.API.AuthenticatedSupport),
				common.IsEnabled(c.Verbose),
			)
		}(configs[x])
	}
	wg.Wait()
	if len(bot.GetExchanges()) == 0 {
		return ErrNoExchangesLoaded
	}
	return nil
}

// WaitForInitialCurrencySync allows for a routine to wait for the initial sync
// of the currency pair syncer management system.
func (bot *Engine) WaitForInitialCurrencySync() error {
	return bot.currencyPairSyncer.WaitForInitialSync()
}

func (bot *Engine) SetupExchangeSettings() error {
	for _, e := range bot.Config.GetEnabledExchanges() {
		for _, a := range []asset.Item{asset.Spot, asset.Futures} {
			enabledPairs, _ := bot.Config.GetEnabledPairs(e, a)
			// fmt.Println("enabled pairs", enabledPairs)
			for _, pair := range enabledPairs {
				fmt.Println("loading pair", pair, e)
				// _, pair, a, err := bot.loadExchangePairAssetBase(e, pair.Base.String(), pair.Quote.String(), asset.Spot)
				// if err != nil {
				// 	fmt.Println("error enabling pair", err)
				// 	return err
				// }
				bot.CurrencySettings = append(bot.CurrencySettings, &ExchangeAssetPairSettings{
					ExchangeName: strings.ToLower(e),
					CurrencyPair: pair,
					AssetType:    a,
				})
			}
		}
	}
	return nil
}

func (bot *Engine) loadExchangePairAssetBase(exch, base, quote, ass string) (exchange.IBotExchange, currency.Pair, asset.Item, error) {
	e, err := bot.GetExchangeByName(exch)
	if err != nil {
		fmt.Println(123, err)
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
	// 	log.Warnf(log.TradeMgr, "no credentials set for %v, this is theoretical only", exchangeBase.Name)
	// }

	fPair, err = exchangeBase.FormatExchangeCurrency(cp, a)
	if err != nil {
		return nil, currency.Pair{}, "", err
	}
	return e, fPair, a, nil
}

func (bot *Engine) GetAllCurrencySettings() ([]*ExchangeAssetPairSettings, error) {
	return bot.CurrencySettings, nil
}

// SetExchangeAssetCurrencySettings sets the settings for an exchange, asset, currency
func (bot *Engine) SetExchangeAssetCurrencySettings(exch string, a asset.Item, cp currency.Pair, c *ExchangeAssetPairSettings) {
	exch = strings.ToLower(exch)
	if c.ExchangeName == "" ||
		c.AssetType == "" ||
		c.CurrencyPair.IsEmpty() {
		return
	}

	for i := range bot.CurrencySettings {
		if bot.CurrencySettings[i].CurrencyPair == cp &&
			bot.CurrencySettings[i].AssetType == a &&
			exch == bot.CurrencySettings[i].ExchangeName {
			bot.CurrencySettings[i] = c
			return
		}
	}
	bot.CurrencySettings = append(bot.CurrencySettings, c)
}

// GetCurrencySettings returns the settings for an exchange, asset currency
func (bot *Engine) GetCurrencySettings(exch string, a asset.Item, cp currency.Pair) (*ExchangeAssetPairSettings, error) {
	exch = strings.ToLower(exch)
	for i := range bot.CurrencySettings {
		if bot.CurrencySettings[i].CurrencyPair.Equal(cp) {
			if bot.CurrencySettings[i].AssetType == a {
				if exch == bot.CurrencySettings[i].ExchangeName {
					return bot.CurrencySettings[i], nil
				}
			}
		}
	}
	return &ExchangeAssetPairSettings{}, fmt.Errorf("no currency settings found for %v %v %v", exch, a, cp)
}

// func GetBadSymbols() []string {
// 	file, err := os.Open("./bad_symbols.txt")
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(123)
// 	}
// 	defer file.Close()
//
// 	pairs := make([]string, 0)
// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		pairs = append(pairs, scanner.Text())
// 	}
// 	return pairs
// }
//
// func IsSymbolInList(p currency.Pair, badSymbols []string) bool {
// 	for _, bad := range badSymbols {
// 		if strings.EqualFold(p.Upper().String(), bad) {
// 			return true
// 		}
// 	}
// 	return false
// }
