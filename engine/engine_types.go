package engine

import (
	"sync"
	"time"
)

// Settings stores engine params
type Settings struct {
	ConfigFile            string
	TradeConfigFile       string
	DataDir               string
	MigrationDir          string
	LogFile               string
	GoMaxProcs            int
	CheckParamInteraction bool

	// Core Settings
	EnableProductionMode        bool
	EnableClearDB               bool
	EnableMachineLearning       bool
	EnableAllExchanges          bool
	EnableAllPairs              bool
	EnableCoinmarketcapAnalysis bool
	EnableCommsRelayer          bool
	EnableConnectivityMonitor   bool
	EnableCurrencyStateManager  bool
	EnableDataHistoryManager    bool
	EnableDatabaseManager       bool
	EnableDepositAddressManager bool
	EnableDryRun                bool
	EnableEventManager          bool
	EnableExchangeSyncManager   bool
	EnableGCTScriptManager      bool
	EnableGRPC                  bool
	EnableGRPCProxy             bool
	EnableLiveMode              bool
	EnableLogging               bool
	EnableNTPClient             bool
	EnableOrderManager          bool
	EnablePortfolioManager      bool
	EnableTradeManager          bool
	EnableTrading               bool
	EnableWatcher               bool
	EnableDataImporter          bool
	EnableWebsocketRPC          bool
	EnableWebsocketRoutine      bool
	EventManagerDelay           time.Duration
	PortfolioManagerDelay       time.Duration
	Verbose                     bool

	// Exchange syncer settings
	EnableTickerSyncing    bool
	EnableOrderbookSyncing bool
	EnableKlineSyncing     bool
	EnableTradeSyncing     bool
	SyncWorkers            int
	SyncContinuously       bool
	SyncTimeoutREST        time.Duration
	SyncTimeoutWebsocket   time.Duration

	// Forex settings
	EnableCurrencyConverter bool
	EnableCurrencyLayer     bool
	EnableFixer             bool
	EnableOpenExchangeRates bool
	EnableExchangeRateHost  bool

	// Exchange tuning settings
	EnableExchangeHTTPRateLimiter  bool
	EnableExchangeHTTPDebugging    bool
	EnableExchangeVerbose          bool
	ExchangePurgeCredentials       bool
	EnableExchangeAutoPairUpdates  bool
	DisableExchangeAutoPairUpdates bool
	EnableExchangeRESTSupport      bool
	EnableExchangeWebsocketSupport bool
	MaxHTTPRequestJobsLimit        int
	TradeBufferProcessingInterval  time.Duration
	RequestMaxRetryAttempts        int

	// Global HTTP related settings
	GlobalHTTPTimeout   time.Duration
	GlobalHTTPUserAgent string
	GlobalHTTPProxy     string

	// Exchange HTTP related settings
	HTTPTimeout   time.Duration
	HTTPUserAgent string
	HTTPProxy     string

	// Dispatch system settings
	EnableDispatcher        bool
	DispatchMaxWorkerAmount int
	DispatchJobsLimit       int

	// GCTscript settings
	MaxVirtualMachines uint

	// Withdraw settings
	WithdrawCacheSize uint64
}

const (
	// MsgStatusOK message to display when status is "OK"
	MsgStatusOK string = "ok"
	// MsgStatusSuccess message to display when status is successful
	MsgStatusSuccess string = "success"
	// MsgStatusError message to display when failure occurs
	MsgStatusError string = "error"
	grpcName       string = "grpc"
	grpcProxyName  string = "grpc_proxy"
)

// newConfigMutex only locks and unlocks on engine creation functions
// as engine modifies global files, this protects the main bot creation
// functions from interfering with each other
var newEngineMutex sync.Mutex
