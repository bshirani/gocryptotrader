package engine

import (
	"sync"
	"time"

	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/exchange/kline"
)

// syncBase stores information
type syncBase struct {
	IsUsingWebsocket bool
	IsUsingREST      bool
	IsProcessing     bool
	LastUpdated      time.Time
	HaveData         bool
	NumErrors        int
}

// currencyPairSyncAgent stores the sync agent info
type currencyPairSyncAgent struct {
	Created   time.Time
	Exchange  string
	AssetType asset.Item
	Pair      currency.Pair
	Ticker    syncBase
	Orderbook syncBase
	Trade     syncBase
	Kline     syncBase
}

// Config stores the currency pair config
type Config struct {
	SyncTicker           bool
	SyncOrderbook        bool
	SyncKlines           bool
	SyncTrades           bool
	SyncContinuously     bool
	SyncTimeoutREST      time.Duration
	SyncTimeoutWebsocket time.Duration
	NumWorkers           int
	Verbose              bool
}

// syncManager stores the exchange currency pair syncer object
type syncManager struct {
	initSyncCompleted              int32
	initSyncStarted                int32
	started                        int32
	delimiter                      string
	uppercase                      bool
	initSyncStartTime              time.Time
	fiatDisplayCurrency            currency.Code
	websocketRoutineManagerEnabled bool
	mux                            sync.Mutex
	initSyncWG                     sync.WaitGroup
	inService                      sync.WaitGroup
	candleSaver                    func(*kline.Item, bool) (uint64, error)
	heartBeatWg                    sync.WaitGroup
	shutdown                       chan struct{}

	currencyPairs            []currencyPairSyncAgent
	tickerBatchLastRequested map[string]time.Time

	remoteConfig    *config.RemoteControlConfig
	config          Config
	exchangeManager iExchangeManager
}
