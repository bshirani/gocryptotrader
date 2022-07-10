package config

import (
	"errors"
	"sync"
	"time"

	"gocryptotrader/communications/base"
	"gocryptotrader/currency"
	"gocryptotrader/database"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/exchange/kline"
	"gocryptotrader/exchange/order"
	"gocryptotrader/exchange/protocol"
	gctscript "gocryptotrader/gctscript/vm"
	"gocryptotrader/log"
	"gocryptotrader/portfolio"
	"gocryptotrader/portfolio/banking"

	"github.com/shopspring/decimal"
)

// Constants declared here are filename strings and test strings
const (
	FXProviderFixer                      = "fixer"
	EncryptedFile                        = "config.dat"
	File                                 = "config.json"
	TestFile                             = "../testdata/configtest.json"
	fileEncryptionPrompt                 = 0
	fileEncryptionEnabled                = 1
	fileEncryptionDisabled               = -1
	pairsLastUpdatedWarningThreshold     = 30 // 30 days
	defaultHTTPTimeout                   = time.Second * 15
	defaultWebsocketResponseCheckTimeout = time.Millisecond * 30
	defaultWebsocketResponseMaxLimit     = time.Second * 7
	defaultWebsocketOrderbookBufferLimit = 5
	defaultWebsocketTrafficTimeout       = time.Second * 30
	maxAuthFailures                      = 3
	defaultNTPAllowedDifference          = 50000000
	defaultNTPAllowedNegativeDifference  = 50000000
	DefaultAPIKey                        = "Key"
	DefaultAPISecret                     = "Secret"
	DefaultAPIClientID                   = "ClientID"
	defaultDataHistoryMonitorCheckTimer  = time.Minute
	defaultCurrencyStateManagerDelay     = time.Minute
	defaultMaxJobsPerCycle               = 50
	defaultDaysBack                      = 2
)

// Constants here hold some messages
const (
	ErrExchangeNameEmpty                       = "exchange #%d name is empty"
	ErrNoEnabledExchanges                      = "no exchanges enabled"
	ErrFailureOpeningConfig                    = "fatal error opening %s file. Error: %s"
	ErrCheckingConfigValues                    = "fatal error checking config values. Error: %s"
	WarningExchangeAuthAPIDefaultOrEmptyValues = "exchange %s authenticated API support disabled due to default/empty APIKey/Secret/ClientID values"
	WarningPairsLastUpdatedThresholdExceeded   = "exchange %s last manual update of available currency pairs has exceeded %d days. Manual update required!"
)

// Constants here define unset default values displayed in the config.json
// file
const (
	APIURLNonDefaultMessage              = "NON_DEFAULT_HTTP_LINK_TO_EXCHANGE_API"
	WebsocketURLNonDefaultMessage        = "NON_DEFAULT_HTTP_LINK_TO_WEBSOCKET_EXCHANGE_API"
	DefaultUnsetAPIKey                   = "Key"
	DefaultUnsetAPISecret                = "Secret"
	DefaultUnsetAccountPlan              = "accountPlan"
	DefaultForexProviderExchangeRatesAPI = "ExchangeRateHost"
)

// Variables here are used for configuration
var (
	Cfg                                 Config
	m                                   sync.Mutex
	ErrExchangeNotFound                 = errors.New("config.exchange not found")
	errBadDate                          = errors.New("start date >= end date, please check your config")
	errNoCurrencySettings               = errors.New("no currency settings set in the config")
	errBadInitialFunds                  = errors.New("initial funds set with invalid data, please check your config")
	errUnsetExchange                    = errors.New("exchange name unset for currency settings, please check your config")
	errUnsetAsset                       = errors.New("asset unset for currency settings, please check your config")
	errUnsetCurrency                    = errors.New("currency unset for currency settings, please check your config")
	errBadSlippageRates                 = errors.New("invalid slippage rates in currency settings, please check your config")
	errStartEndUnset                    = errors.New("data start and end dates are invalid, please check your config")
	errSimultaneousProcessingRequired   = errors.New("exchange level funding requires simultaneous processing, please check your config and view funding readme for details")
	errExchangeLevelFundingRequired     = errors.New("invalid config, funding details set while exchange level funding is disabled")
	errExchangeLevelFundingDataRequired = errors.New("invalid config, exchange level funding enabled with no funding data set")
	errSizeLessThanZero                 = errors.New("size less than zero")
	errMaxSizeMinSizeMismatch           = errors.New("maximum size must be greater to minimum size")
	errMinMaxEqual                      = errors.New("minimum and maximum limits cannot be equal")
)

// Config is the overarching object that holds all the information for
// prestart management of Portfolio, Communications, Webserver and Enabled
// Exchanges
type Config struct {
	ProductionMode bool `json:"productionMode"`
	LiveMode       bool `json:"liveMode"`
	DryRun         bool `json:"dryRun"`
	Script         bool `json:"script"`

	Nickname          string             `json:"nickname"`
	Goal              string             `json:"goal"`
	CurrencySettings  []CurrencySettings `json:"currencySettings"`
	DataSettings      DataSettings       `json:"dataSettings"`
	PortfolioSettings PortfolioSettings  `json:"portfolioSettings"`

	GoCryptoTraderConfigPath string `json:"gocryptotraderConfigPath"`

	StatisticSettings StatisticSettings  `json:"statisticSettings"`
	SyncManager       SyncManagerConfig  `json:"syncManager"`
	StrategyManager      StrategyManagerConfig `json:"tradeManager"`
	FactorEngine      FactorEngineConfig `json:"factorEngine"`
	OrderManager      OrderManagerConfig `json:"orderManager"`
	DataImporter      DataImporterConfig `json:"dataImporter"`

	Name                 string                    `json:"name"`
	DataDirectory        string                    `json:"dataDirectory"`
	EncryptConfig        int                       `json:"encryptConfig"`
	GlobalHTTPTimeout    time.Duration             `json:"globalHTTPTimeout"`
	Database             database.Config           `json:"database"`
	Logging              log.Config                `json:"logging"`
	ConnectionMonitor    ConnectionMonitorConfig   `json:"connectionMonitor"`
	DataHistory          DataHistoryManager        `json:"dataHistory"`
	CurrencyStateManager CurrencyStateManager      `json:"currencyStateManager"`
	Watcher              Watcher                   `json:"watcher"`
	Profiler             Profiler                  `json:"profiler"`
	NTPClient            NTPClientConfig           `json:"ntpclient"`
	GCTScript            gctscript.Config          `json:"gctscript"`
	Currency             CurrencyConfig            `json:"currencyConfig"`
	Communications       base.CommunicationsConfig `json:"communications"`
	RemoteControl        RemoteControlConfig       `json:"remoteControl"`
	Portfolio            portfolio.Base            `json:"portfolioAddresses"`
	Exchanges            []ExchangeConfig          `json:"exchanges"`
	BankAccounts         []banking.Account         `json:"bankAccounts"`

	// Deprecated config settings, will be removed at a future date
	Webserver           *WebserverConfig          `json:"webserver,omitempty"`
	CurrencyPairFormat  *CurrencyPairFormatConfig `json:"currencyPairFormat,omitempty"`
	FiatDisplayCurrency *currency.Code            `json:"fiatDispayCurrency,omitempty"`
	Cryptocurrencies    *currency.Currencies      `json:"cryptocurrencies,omitempty"`
	SMS                 *base.SMSGlobalConfig     `json:"smsGlobal,omitempty"`
	// encryption session values
	storedSalt []byte
	sessionDK  []byte
}

// DataHistoryManager holds all information required for the data history manager
type DataImporterConfig struct {
	Enabled       bool   `json:"enabled"`
	Verbose       bool   `json:"verbose"`
	DataDirectory string `json:"dataDirectory"`
	FromDate      string `json:"fromDate"`
	ToDate        string `json:"toDate"`
	NumWorkers    int    `json:"numWorkers"`
}

// DataHistoryManager holds all information required for the data history manager
type StrategyManagerConfig struct {
	Enabled    bool               `json:"enabled"`
	Verbose    bool               `json:"verbose"`
	Trading    bool               `json:"trading"`
	UseML      bool               `json:"useML"`
	Debug      bool               `json:"debug"`
	ClearDB    bool               `json:"clearDB"`
	Strategies []*StrategySetting `json:"strategies"`
}

type StrategySetting struct {
	Weight    decimal.Decimal `json:"weight"`
	Side      order.Side      `json:"side"`
	AssetType asset.Item      `json:"assetType"`
	Interval  kline.Interval  `json:"klineInterval"`
	Pair      currency.Pair   `json:"pair"`
	Capture   string          `json:"capture"`
	Debug     bool            `json:"debug"`
}

type FactorEngineConfig struct {
	Verbose bool `json:"verbose"`
}

type SyncManagerConfig struct {
	Verbose              bool          `json:"verbose"`
	Enabled              bool          `json:"enabled"`
	TickerSync           bool          `json:"tickerSync"`
	OrderBookSync        bool          `json:"orderBookSync"`
	KlineSync            bool          `json:"klineSync"`
	TradeSync            bool          `json:"tradeSync"`
	SyncContinuously     bool          `json:"syncContinuously"`
	SyncWorkers          int           `json:"syncWorkers"`
	SyncTimeoutREST      time.Duration `json:"syncTimeoutRest"`
	SyncTimeoutWebsocket time.Duration `json:"syncTimeoutWebsocket"`
}

type OrderManagerConfig struct {
	Verbose bool `json:"verbose"`
	Enabled bool `json:"enabled"`
}

// DataHistoryManager holds all information required for the data history manager
type DataHistoryManager struct {
	Enabled             bool          `json:"enabled"`
	Verbose             bool          `json:"verbose"`
	CheckInterval       time.Duration `json:"checkInterval"`
	MaxJobsPerCycle     int64         `json:"maxJobsPerCycle"`
	MaxResultInsertions int64         `json:"maxResultInsertions"`
	DaysBack            int64         `json:"daysBack"`
}

// CurrencyStateManager defines a set of configuration options for the currency
// state manager
type CurrencyStateManager struct {
	Enabled *bool         `json:"enabled"`
	Delay   time.Duration `json:"delay"`
}

type Watcher struct {
	Enabled *bool         `json:"enabled"`
	Delay   time.Duration `json:"delay"`
}

// ConnectionMonitorConfig defines the connection monitor variables to ensure
// that there is internet connectivity
type ConnectionMonitorConfig struct {
	DNSList          []string      `json:"preferredDNSList"`
	PublicDomainList []string      `json:"preferredDomainList"`
	CheckInterval    time.Duration `json:"checkInterval"`
}

// ExchangeConfig holds all the information needed for each enabled Exchange.
type ExchangeConfig struct {
	Name                          string                 `json:"name"`
	Enabled                       bool                   `json:"enabled"`
	Verbose                       bool                   `json:"verbose"`
	UseSandbox                    bool                   `json:"useSandbox,omitempty"`
	HTTPTimeout                   time.Duration          `json:"httpTimeout"`
	HTTPUserAgent                 string                 `json:"httpUserAgent,omitempty"`
	HTTPDebugging                 bool                   `json:"httpDebugging,omitempty"`
	WebsocketResponseCheckTimeout time.Duration          `json:"websocketResponseCheckTimeout"`
	WebsocketResponseMaxLimit     time.Duration          `json:"websocketResponseMaxLimit"`
	WebsocketTrafficTimeout       time.Duration          `json:"websocketTrafficTimeout"`
	ProxyAddress                  string                 `json:"proxyAddress,omitempty"`
	BaseCurrencies                currency.Currencies    `json:"baseCurrencies"`
	CurrencyPairs                 *currency.PairsManager `json:"currencyPairs"`
	API                           APIConfig              `json:"api"`
	Features                      *FeaturesConfig        `json:"features"`
	BankAccounts                  []banking.Account      `json:"bankAccounts,omitempty"`
	OrderbookConfig               `json:"orderbook"`

	// Deprecated settings which will be removed in a future update
	AvailablePairs                   *currency.Pairs      `json:"availablePairs,omitempty"`
	EnabledPairs                     *currency.Pairs      `json:"enabledPairs,omitempty"`
	AssetTypes                       *string              `json:"assetTypes,omitempty"`
	PairsLastUpdated                 *int64               `json:"pairsLastUpdated,omitempty"`
	ConfigCurrencyPairFormat         *currency.PairFormat `json:"configCurrencyPairFormat,omitempty"`
	RequestCurrencyPairFormat        *currency.PairFormat `json:"requestCurrencyPairFormat,omitempty"`
	AuthenticatedAPISupport          *bool                `json:"authenticatedApiSupport,omitempty"`
	AuthenticatedWebsocketAPISupport *bool                `json:"authenticatedWebsocketApiSupport,omitempty"`
	APIKey                           *string              `json:"apiKey,omitempty"`
	APISecret                        *string              `json:"apiSecret,omitempty"`
	APIAuthPEMKeySupport             *bool                `json:"apiAuthPemKeySupport,omitempty"`
	APIAuthPEMKey                    *string              `json:"apiAuthPemKey,omitempty"`
	APIURL                           *string              `json:"apiUrl,omitempty"`
	APIURLSecondary                  *string              `json:"apiUrlSecondary,omitempty"`
	ClientID                         *string              `json:"clientId,omitempty"`
	SupportsAutoPairUpdates          *bool                `json:"supportsAutoPairUpdates,omitempty"`
	Websocket                        *bool                `json:"websocket,omitempty"`
	WebsocketURL                     *string              `json:"websocketUrl,omitempty"`
}

// Profiler defines the profiler configuration to enable pprof
type Profiler struct {
	Enabled              bool `json:"enabled"`
	MutexProfileFraction int  `json:"mutexProfileFraction"`
}

// NTPClientConfig defines a network time protocol configuration to allow for
// positive and negative differences
type NTPClientConfig struct {
	Level                     int            `json:"enabled"`
	Pool                      []string       `json:"pool"`
	AllowedDifference         *time.Duration `json:"allowedDifference"`
	AllowedNegativeDifference *time.Duration `json:"allowedNegativeDifference"`
}

// GRPCConfig stores the gRPC settings
type GRPCConfig struct {
	Enabled                bool   `json:"enabled"`
	ListenAddress          string `json:"listenAddress"`
	GRPCProxyEnabled       bool   `json:"grpcProxyEnabled"`
	GRPCProxyListenAddress string `json:"grpcProxyListenAddress"`
	TimeInNanoSeconds      bool   `json:"timeInNanoSeconds"`
}

// DepcrecatedRPCConfig stores the deprecatedRPCConfig settings
type DepcrecatedRPCConfig struct {
	Enabled       bool   `json:"enabled"`
	ListenAddress string `json:"listenAddress"`
}

// WebsocketRPCConfig stores the websocket config info
type WebsocketRPCConfig struct {
	Enabled             bool   `json:"enabled"`
	ListenAddress       string `json:"listenAddress"`
	ConnectionLimit     int    `json:"connectionLimit"`
	MaxAuthFailures     int    `json:"maxAuthFailures"`
	AllowInsecureOrigin bool   `json:"allowInsecureOrigin"`
}

// RemoteControlConfig stores the RPC services config
type RemoteControlConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`

	GRPC          GRPCConfig           `json:"gRPC"`
	DeprecatedRPC DepcrecatedRPCConfig `json:"deprecatedRPC"`
	WebsocketRPC  WebsocketRPCConfig   `json:"websocketRPC"`
}

// WebserverConfig stores the old webserver config
type WebserverConfig struct {
	Enabled                      bool   `json:"enabled"`
	AdminUsername                string `json:"adminUsername"`
	AdminPassword                string `json:"adminPassword"`
	ListenAddress                string `json:"listenAddress"`
	WebsocketConnectionLimit     int    `json:"websocketConnectionLimit"`
	WebsocketMaxAuthFailures     int    `json:"websocketMaxAuthFailures"`
	WebsocketAllowInsecureOrigin bool   `json:"websocketAllowInsecureOrigin"`
}

// Post holds the bot configuration data
type Post struct {
	Data Config `json:"data"`
}

// CurrencyPairFormatConfig stores the users preferred currency pair display
type CurrencyPairFormatConfig struct {
	Uppercase bool   `json:"uppercase"`
	Delimiter string `json:"delimiter,omitempty"`
	Separator string `json:"separator,omitempty"`
	Index     string `json:"index,omitempty"`
}

// BankTransaction defines a related banking transaction
type BankTransaction struct {
	ReferenceNumber     string `json:"referenceNumber"`
	TransactionNumber   string `json:"transactionNumber"`
	PaymentInstructions string `json:"paymentInstructions"`
}

// CurrencyConfig holds all the information needed for currency related manipulation
type CurrencyConfig struct {
	ForexProviders                []currency.FXSettings     `json:"forexProviders"`
	CryptocurrencyProvider        CryptocurrencyProvider    `json:"cryptocurrencyProvider"`
	Cryptocurrencies              currency.Currencies       `json:"cryptocurrencies"`
	CurrencyPairFormat            *CurrencyPairFormatConfig `json:"currencyPairFormat"`
	FiatDisplayCurrency           currency.Code             `json:"fiatDisplayCurrency"`
	CurrencyFileUpdateDuration    time.Duration             `json:"currencyFileUpdateDuration"`
	ForeignExchangeUpdateDuration time.Duration             `json:"foreignExchangeUpdateDuration"`
}

// CryptocurrencyProvider defines coinmarketcap tools
type CryptocurrencyProvider struct {
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Verbose     bool   `json:"verbose"`
	APIkey      string `json:"apiKey"`
	AccountPlan string `json:"accountPlan"`
}

// FeaturesSupportedConfig stores the exchanges supported features
type FeaturesSupportedConfig struct {
	REST                  bool              `json:"restAPI"`
	RESTCapabilities      protocol.Features `json:"restCapabilities,omitempty"`
	Websocket             bool              `json:"websocketAPI"`
	WebsocketCapabilities protocol.Features `json:"websocketCapabilities,omitempty"`
}

// FeaturesEnabledConfig stores the exchanges enabled features
type FeaturesEnabledConfig struct {
	AutoPairUpdates bool `json:"autoPairUpdates"`
	Websocket       bool `json:"websocketAPI"`
	SaveTradeData   bool `json:"saveTradeData"`
}

// FeaturesConfig stores the exchanges supported and enabled features
type FeaturesConfig struct {
	Supports FeaturesSupportedConfig `json:"supports"`
	Enabled  FeaturesEnabledConfig   `json:"enabled"`
}

// APIEndpointsConfig stores the API endpoint addresses
type APIEndpointsConfig struct {
	URL          string `json:"url"`
	URLSecondary string `json:"urlSecondary"`
	WebsocketURL string `json:"websocketURL"`
}

// APICredentialsConfig stores the API credentials
type APICredentialsConfig struct {
	Key        string `json:"key,omitempty"`
	Secret     string `json:"secret,omitempty"`
	ClientID   string `json:"clientID,omitempty"`
	Subaccount string `json:"subaccount,omitempty"`
	PEMKey     string `json:"pemKey,omitempty"`
	OTPSecret  string `json:"otpSecret,omitempty"`
}

// APICredentialsValidatorConfig stores the API credentials validator settings
type APICredentialsValidatorConfig struct {
	// For Huobi (optional)
	RequiresPEM bool `json:"requiresPEM,omitempty"`

	RequiresKey                bool `json:"requiresKey,omitempty"`
	RequiresSecret             bool `json:"requiresSecret,omitempty"`
	RequiresClientID           bool `json:"requiresClientID,omitempty"`
	RequiresBase64DecodeSecret bool `json:"requiresBase64DecodeSecret,omitempty"`
}

// APIConfig stores the exchange API config
type APIConfig struct {
	AuthenticatedSupport          bool `json:"authenticatedSupport"`
	AuthenticatedWebsocketSupport bool `json:"authenticatedWebsocketApiSupport"`
	PEMKeySupport                 bool `json:"pemKeySupport,omitempty"`

	Credentials          APICredentialsConfig           `json:"credentials"`
	CredentialsValidator *APICredentialsValidatorConfig `json:"credentialsValidator,omitempty"`
	OldEndPoints         *APIEndpointsConfig            `json:"endpoints,omitempty"`
	Endpoints            map[string]string              `json:"urlEndpoints"`
}

// OrderbookConfig stores the orderbook configuration variables
type OrderbookConfig struct {
	VerificationBypass     bool `json:"verificationBypass"`
	WebsocketBufferLimit   int  `json:"websocketBufferLimit"`
	WebsocketBufferEnabled bool `json:"websocketBufferEnabled"`
}

// DataSettings is a container for each type of data retrieval setting.
// Only ONE can be populated per config
type DataSettings struct {
	Interval     time.Duration `json:"interval"`
	DataType     string        `json:"dataType"`
	APIData      *APIData      `json:"apiData,omitempty"`
	DatabaseData *DatabaseData `json:"databaseData,omitempty"`
	LiveData     *LiveData     `json:"liveData,omitempty"`
	CSVData      *CSVData      `json:"csvData,omitempty"`
}

// ExchangeLevelFunding allows the portfolio manager to access
// a shared pool. For example, The base currencies BTC and LTC can both
// access the same USDT funding to make purchasing decisions
// Similarly, when a BTC is sold, LTC can now utilise the increased funding
// Importantly, exchange level funding is all-inclusive, you cannot have it for only some uses
// It also is required to use SimultaneousSignalProcessing, otherwise the first currency processed
// will have dibs
type ExchangeLevelFunding struct {
	ExchangeName string          `json:"exchangeName"`
	Asset        string          `json:"asset"`
	Currency     string          `json:"currency"`
	InitialFunds decimal.Decimal `json:"initialFunds"`
	TransferFee  decimal.Decimal `json:"transferFee"`
}

// StatisticSettings adjusts ratios where
// proper data is currently lacking
type StatisticSettings struct {
	RiskFreeRate decimal.Decimal `json:"riskFreeRate"`
}

// PortfolioSettings act as a global protector for strategies
// these settings will override ExchangeSettings that go against it
// and assess the bigger picture
type PortfolioSettings struct {
	EnabledPairs      string   `json:"enabledPairs"`
	EnabledStrategies string   `json:"enabledStrategies"`
	Verbose           bool     `json:"verbose"`
	Leverage          Leverage `json:"leverage"`
	BuySide           MinMax   `json:"buySide"`
	SellSide          MinMax   `json:"sellSide"`
}

// Leverage rules are used to allow or limit the use of leverage in orders
// when supported
type Leverage struct {
	CanUseLeverage                 bool            `json:"canUseLeverage"`
	MaximumOrdersWithLeverageRatio decimal.Decimal `json:"maximumOrdersWithLeverageRatio"`
	MaximumLeverageRate            decimal.Decimal `json:"maximumLeverageRate"`
}

// MinMax are the rules which limit the placement of orders.
type MinMax struct {
	MinimumSize  decimal.Decimal `json:"minimumSize"` // will not place an order if under this amount
	MaximumSize  decimal.Decimal `json:"maximumSize"` // can only place an order up to this amount
	MaximumTotal decimal.Decimal `json:"maximumTotal"`
}

// CurrencySettings stores pair based variables
// It contains rules about the specific currency pair
// you wish to trade with
// Backtester will load the data of the currencies specified here
type CurrencySettings struct {
	ExchangeName string `json:"exchangeName"`
	Asset        string `json:"asset"`
	Base         string `json:"base"`
	Quote        string `json:"quote"`

	InitialBaseFunds   *decimal.Decimal `json:"initialBaseFunds,omitempty"`
	InitialQuoteFunds  *decimal.Decimal `json:"initialQuoteFunds,omitempty"`
	InitialLegacyFunds float64          `json:"initialFunds,omitempty"`

	Leverage Leverage `json:"leverage"`
	BuySide  MinMax   `json:"buySide"`
	SellSide MinMax   `json:"sellSide"`

	MinimumSlippagePercent decimal.Decimal `json:"minSlippagePercent"`
	MaximumSlippagePercent decimal.Decimal `json:"maxSlippagePercent"`

	MakerFee decimal.Decimal `json:"makerFeeOverride"`
	TakerFee decimal.Decimal `json:"takerFeeOverride"`

	MaximumHoldingsRatio decimal.Decimal `json:"maximumHoldingsRatio"`

	CanUseExchangeLimits          bool `json:"useExchangeOrderLimits"`
	SkipCandleVolumeFitting       bool `json:"skipCandleVolumeFitting"`
	ShowExchangeOrderLimitWarning bool `json:"-"`
}

// APIData defines all fields to configure API based data
type APIData struct {
	StartDate        time.Time `json:"startDate"`
	EndDate          time.Time `json:"endDate"`
	InclusiveEndDate bool      `json:"inclusiveEndDate"`
}

// CSVData defines all fields to configure CSV based data
type CSVData struct {
	FullPath string `json:"fullOath"`
}

// DatabaseData defines all fields to configure database based data
type DatabaseData struct {
	StartDate        time.Time        `json:"startDate"`
	EndDate          time.Time        `json:"endDate"`
	ConfigOverride   *database.Config `json:"configOverride"`
	InclusiveEndDate bool             `json:"inclusiveEndDate"`
}

// LiveData defines all fields to configure live data
type LiveData struct {
	APIKeyOverride        string `json:"apiKeyOverride"`
	APISecretOverride     string `json:"apiSecretOverride"`
	APIClientIDOverride   string `json:"apiClientIdOverride"`
	API2FAOverride        string `json:"api2faOverride"`
	APISubAccountOverride string `json:"apiSubSccountOverride"`
	RealOrders            bool   `json:"realOrders"`
}
