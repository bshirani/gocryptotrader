package engine

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"gocryptotrader/common"
	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/exchange/orderbook"
	"gocryptotrader/exchange/stats"
	"gocryptotrader/exchange/ticker"
	"gocryptotrader/log"
)

// const holds the sync item types
const (
	CandleSyncItemTicker = iota
	CandleSyncItemOrderbook
	CandleSyncItemTrade
	CandleSyncManagerName = "candle_syncer"
)

var (
	candleCreatedCounter = 0
	candleRemovedCounter = 0
	// DefaultCandleSyncerWorkers limits the number of sync workers
	DefaultCandleSyncerWorkers = 1
	// DefaultCandleSyncerTimeoutREST the default time to switch from REST to websocket protocols without a response
	DefaultCandleSyncerTimeoutREST = time.Second * 15
	// DefaultCandleSyncerTimeoutWebsocket the default time to switch from websocket to REST protocols without a response
	DefaultCandleSyncerTimeoutWebsocket = time.Minute
	errNoCandleSyncItemsEnabled         = errors.New("no sync items enabled")
	errUnknownCandleSyncItem            = errors.New("unknown sync item")
	errCandleSyncPairNotFound           = errors.New("exchange currency pair syncer not found")
)

// setupCandleSyncManager starts a new CandleSyncer
func setupCandleSyncManager(c *Config, exchangeManager iExchangeManager, remoteConfig *config.RemoteControlConfig, websocketRoutineManagerEnabled bool) (*candleSyncManager, error) {
	// if !c.SyncOrderbook && !c.SyncTicker && !c.SyncTrades {
	// 	return nil, errNoCandleSyncItemsEnabled
	// }
	if exchangeManager == nil {
		return nil, errNilExchangeManager
	}
	if remoteConfig == nil {
		return nil, errNilConfig
	}

	if c.NumWorkers <= 0 {
		c.NumWorkers = DefaultCandleSyncerWorkers
	}

	if c.SyncTimeoutREST <= time.Duration(0) {
		c.SyncTimeoutREST = DefaultCandleSyncerTimeoutREST
	}

	if c.SyncTimeoutWebsocket <= time.Duration(0) {
		c.SyncTimeoutWebsocket = DefaultCandleSyncerTimeoutWebsocket
	}

	s := &candleSyncManager{
		config:                         *c,
		remoteConfig:                   remoteConfig,
		exchangeManager:                exchangeManager,
		websocketRoutineManagerEnabled: websocketRoutineManagerEnabled,
	}

	s.tickerBatchLastRequested = make(map[string]time.Time)

	log.Debugf(log.CandleSyncMgr,
		"Exchange currency pair syncer config: continuous: %v ticker: %v"+
			" orderbook: %v trades: %v workers: %v verbose: %v timeout REST: %v"+
			" timeout Websocket: %v",
		s.config.SyncContinuously, s.config.SyncTicker, s.config.SyncOrderbook,
		s.config.SyncTrades, s.config.NumWorkers, s.config.Verbose, s.config.SyncTimeoutREST,
		s.config.SyncTimeoutWebsocket)
	s.inService.Add(1)
	return s, nil
}

// IsRunning safely checks whether the subsystem is running
func (m *candleSyncManager) IsRunning() bool {
	if m == nil {
		return false
	}
	return atomic.LoadInt32(&m.started) == 1
}

// Start runs the subsystem
func (m *candleSyncManager) Start() error {
	if m == nil {
		return fmt.Errorf("exchange CandleSyncer %w", ErrNilSubsystem)
	}
	if !atomic.CompareAndSwapInt32(&m.started, 0, 1) {
		return ErrSubSystemAlreadyStarted
	}
	m.initSyncWG.Add(1)
	m.inService.Done()
	log.Debugln(log.CandleSyncMgr, "Exchange CandleSyncer started.")
	exchanges, err := m.exchangeManager.GetExchanges()
	if err != nil {
		return err
	}
	for x := range exchanges {
		exchangeName := exchanges[x].GetName()
		supportsWebsocket := exchanges[x].SupportsWebsocket()
		supportsREST := exchanges[x].SupportsREST()

		if !supportsREST && !supportsWebsocket {
			log.Warnf(log.CandleSyncMgr,
				"Loaded exchange %s does not support REST or Websocket.",
				exchangeName)
			continue
		}

		var usingWebsocket bool
		var usingREST bool
		if m.websocketRoutineManagerEnabled &&
			supportsWebsocket &&
			exchanges[x].IsWebsocketEnabled() {
			usingWebsocket = true
		} else if supportsREST {
			usingREST = true
		}

		assetTypes := exchanges[x].GetAssetTypes(false)
		for y := range assetTypes {
			if exchanges[x].GetBase().CurrencyPairs.IsAssetEnabled(assetTypes[y]) != nil {
				log.Warnf(log.CandleSyncMgr,
					"%s asset type %s is disabled, fetching enabled pairs is paused",
					exchangeName,
					assetTypes[y])
				continue
			}

			wsAssetSupported := exchanges[x].IsAssetWebsocketSupported(assetTypes[y])
			if !wsAssetSupported {
				log.Warnf(log.CandleSyncMgr,
					"%s asset type %s websocket functionality is unsupported, REST fetching only.",
					exchangeName,
					assetTypes[y])
			}
			enabledPairs, err := exchanges[x].GetEnabledPairs(assetTypes[y])
			if err != nil {
				log.Errorf(log.CandleSyncMgr,
					"%s failed to get enabled pairs. Err: %s",
					exchangeName,
					err)
				continue
			}
			for i := range enabledPairs {
				if m.exists(exchangeName, enabledPairs[i], assetTypes[y]) {
					continue
				}

				c := &candleSyncAgent{
					AssetType: assetTypes[y],
					Exchange:  exchangeName,
					Pair:      enabledPairs[i],
				}
				sBase := syncBase{
					IsUsingREST:      usingREST || !wsAssetSupported,
					IsUsingWebsocket: usingWebsocket && wsAssetSupported,
				}
				if m.config.SyncTicker {
					c.Ticker = sBase
				}
				if m.config.SyncOrderbook {
					c.Orderbook = sBase
				}
				if m.config.SyncTrades {
					c.Trade = sBase
				}

				m.add(c)
			}
		}
	}

	if atomic.CompareAndSwapInt32(&m.initSyncStarted, 0, 1) {
		log.Debugf(log.CandleSyncMgr,
			"Exchange CandleSyncer initial sync started. %d items to process.",
			candleCreatedCounter)
		m.initSyncStartTime = time.Now()
	}

	go func() {
		m.initSyncWG.Wait()
		if atomic.CompareAndSwapInt32(&m.initSyncCompleted, 0, 1) {
			log.Debugf(log.CandleSyncMgr, "Exchange CandleSyncer initial sync is complete.")
			completedTime := time.Now()
			log.Debugf(log.CandleSyncMgr, "Exchange CandleSyncer initial sync took %v [%v sync items].",
				completedTime.Sub(m.initSyncStartTime), candleCreatedCounter)

			if !m.config.SyncContinuously {
				log.Debugln(log.CandleSyncMgr, "Exchange CandleSyncer stopping.")
				err := m.Stop()
				if err != nil {
					log.Error(log.CandleSyncMgr, err)
				}
				return
			}
		}
	}()

	if atomic.LoadInt32(&m.initSyncCompleted) == 1 && !m.config.SyncContinuously {
		return nil
	}

	for i := 0; i < m.config.NumWorkers; i++ {
		go m.worker()
	}
	m.initSyncWG.Done()
	return nil
}

// Stop shuts down the exchange currency pair syncer
func (m *candleSyncManager) Stop() error {
	if m == nil {
		return fmt.Errorf("exchange CandleSyncer %w", ErrNilSubsystem)
	}
	if !atomic.CompareAndSwapInt32(&m.started, 1, 0) {
		return fmt.Errorf("exchange CandleSyncer %w", ErrSubSystemNotStarted)
	}
	m.inService.Add(1)
	log.Debugln(log.CandleSyncMgr, "Exchange CandleSyncer stopped.")
	return nil
}

func (m *candleSyncManager) get(exchangeName string, p currency.Pair, a asset.Item) (*candleSyncAgent, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	for x := range m.currencyPairs {
		if m.currencyPairs[x].Exchange == exchangeName &&
			m.currencyPairs[x].Pair.Equal(p) &&
			m.currencyPairs[x].AssetType == a {
			return &m.currencyPairs[x], nil
		}
	}

	return nil, fmt.Errorf("%v %v %v %w", exchangeName, a, p, errCandleSyncPairNotFound)
}

func (m *candleSyncManager) exists(exchangeName string, p currency.Pair, a asset.Item) bool {
	m.mux.Lock()
	defer m.mux.Unlock()

	for x := range m.currencyPairs {
		if m.currencyPairs[x].Exchange == exchangeName &&
			m.currencyPairs[x].Pair.Equal(p) &&
			m.currencyPairs[x].AssetType == a {
			return true
		}
	}
	return false
}

func (m *candleSyncManager) add(c *candleSyncAgent) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if m.config.SyncTicker {
		if m.config.Verbose {
			log.Debugf(log.CandleSyncMgr,
				"%s: Added ticker sync item %v: using websocket: %v using REST: %v",
				c.Exchange, m.FormatCurrency(c.Pair).String(), c.Ticker.IsUsingWebsocket,
				c.Ticker.IsUsingREST)
		}
		if atomic.LoadInt32(&m.initSyncCompleted) != 1 {
			m.initSyncWG.Add(1)
			candleCreatedCounter++
		}
	}

	c.Created = time.Now()
	m.currencyPairs = append(m.currencyPairs, *c)
}

func (m *candleSyncManager) isProcessing(exchangeName string, p currency.Pair, a asset.Item, syncType int) bool {
	m.mux.Lock()
	defer m.mux.Unlock()

	for x := range m.currencyPairs {
		if m.currencyPairs[x].Exchange == exchangeName &&
			m.currencyPairs[x].Pair.Equal(p) &&
			m.currencyPairs[x].AssetType == a {
			switch syncType {
			case CandleSyncItemTicker:
				return m.currencyPairs[x].Ticker.IsProcessing
			case CandleSyncItemOrderbook:
				return m.currencyPairs[x].Orderbook.IsProcessing
			case CandleSyncItemTrade:
				return m.currencyPairs[x].Trade.IsProcessing
			}
		}
	}

	return false
}

func (m *candleSyncManager) setProcessing(exchangeName string, p currency.Pair, a asset.Item, syncType int, processing bool) {
	m.mux.Lock()
	defer m.mux.Unlock()

	for x := range m.currencyPairs {
		if m.currencyPairs[x].Exchange == exchangeName &&
			m.currencyPairs[x].Pair.Equal(p) &&
			m.currencyPairs[x].AssetType == a {
			switch syncType {
			case CandleSyncItemTicker:
				m.currencyPairs[x].Ticker.IsProcessing = processing
			case CandleSyncItemOrderbook:
				m.currencyPairs[x].Orderbook.IsProcessing = processing
			case CandleSyncItemTrade:
				m.currencyPairs[x].Trade.IsProcessing = processing
			}
		}
	}
}

// Update notifies the candleSyncManager to change the last updated time for a exchange asset pair
func (m *candleSyncManager) Update(exchangeName string, p currency.Pair, a asset.Item, syncType int, err error) error {
	if m == nil {
		return fmt.Errorf("exchange CandleSyncer %w", ErrNilSubsystem)
	}
	if atomic.LoadInt32(&m.started) == 0 {
		return fmt.Errorf("exchange CandleSyncer %w", ErrSubSystemNotStarted)
	}

	if atomic.LoadInt32(&m.initSyncStarted) != 1 {
		return nil
	}

	switch syncType {
	case CandleSyncItemOrderbook:
		if !m.config.SyncOrderbook {
			return nil
		}
	case CandleSyncItemTicker:
		if !m.config.SyncTicker {
			return nil
		}
	case CandleSyncItemTrade:
		if !m.config.SyncTrades {
			return nil
		}
	default:
		return fmt.Errorf("%v %w", syncType, errUnknownCandleSyncItem)
	}

	m.mux.Lock()
	defer m.mux.Unlock()

	for x := range m.currencyPairs {
		if m.currencyPairs[x].Exchange == exchangeName &&
			m.currencyPairs[x].Pair.Equal(p) &&
			m.currencyPairs[x].AssetType == a {
			switch syncType {
			case CandleSyncItemTicker:
				origHadData := m.currencyPairs[x].Ticker.HaveData
				m.currencyPairs[x].Ticker.LastUpdated = time.Now()
				if err != nil {
					m.currencyPairs[x].Ticker.NumErrors++
				}
				m.currencyPairs[x].Ticker.HaveData = true
				m.currencyPairs[x].Ticker.IsProcessing = false
				if atomic.LoadInt32(&m.initSyncCompleted) != 1 && !origHadData {
					candleRemovedCounter++
					log.Debugf(log.CandleSyncMgr, "%s ticker sync complete %v [%d/%d].",
						exchangeName,
						m.FormatCurrency(p).String(),
						candleRemovedCounter,
						candleCreatedCounter)
					m.initSyncWG.Done()
				}

			case CandleSyncItemOrderbook:
				origHadData := m.currencyPairs[x].Orderbook.HaveData
				m.currencyPairs[x].Orderbook.LastUpdated = time.Now()
				if err != nil {
					m.currencyPairs[x].Orderbook.NumErrors++
				}
				m.currencyPairs[x].Orderbook.HaveData = true
				m.currencyPairs[x].Orderbook.IsProcessing = false
				if atomic.LoadInt32(&m.initSyncCompleted) != 1 && !origHadData {
					candleRemovedCounter++
					log.Debugf(log.CandleSyncMgr, "%s orderbook sync complete %v [%d/%d].",
						exchangeName,
						m.FormatCurrency(p).String(),
						candleRemovedCounter,
						candleCreatedCounter)
					m.initSyncWG.Done()
				}

			case CandleSyncItemTrade:
				origHadData := m.currencyPairs[x].Trade.HaveData
				m.currencyPairs[x].Trade.LastUpdated = time.Now()
				if err != nil {
					m.currencyPairs[x].Trade.NumErrors++
				}
				m.currencyPairs[x].Trade.HaveData = true
				m.currencyPairs[x].Trade.IsProcessing = false
				if atomic.LoadInt32(&m.initSyncCompleted) != 1 && !origHadData {
					candleRemovedCounter++
					log.Debugf(log.CandleSyncMgr, "%s trade sync complete %v [%d/%d].",
						exchangeName,
						m.FormatCurrency(p).String(),
						candleRemovedCounter,
						candleCreatedCounter)
					m.initSyncWG.Done()
				}
			}
		}
	}
	return nil
}

func (m *candleSyncManager) worker() {
	cleanup := func() {
		log.Debugln(log.CandleSyncMgr,
			"Exchange CandleSyncer worker shutting down.")
	}
	defer cleanup()

	for atomic.LoadInt32(&m.started) != 0 {
		exchanges, err := m.exchangeManager.GetExchanges()
		if err != nil {
			log.Errorf(log.CandleSyncMgr, "CandleSync manager cannot get exchanges: %v", err)
		}
		for x := range exchanges {
			exchangeName := exchanges[x].GetName()
			supportsREST := exchanges[x].SupportsREST()
			supportsRESTTickerBatching := exchanges[x].SupportsRESTTickerBatchUpdates()
			var usingREST bool
			var usingWebsocket bool
			var switchedToRest bool
			if exchanges[x].SupportsWebsocket() && exchanges[x].IsWebsocketEnabled() {
				ws, err := exchanges[x].GetWebsocket()
				if err != nil {
					log.Errorf(log.CandleSyncMgr,
						"%s unable to get websocket pointer. Err: %s",
						exchangeName,
						err)
					usingREST = true
				}

				if ws.IsConnected() {
					usingWebsocket = true
				} else {
					usingREST = true
				}
			} else if supportsREST {
				usingREST = true
			}

			assetTypes := exchanges[x].GetAssetTypes(true)
			for y := range assetTypes {
				wsAssetSupported := exchanges[x].IsAssetWebsocketSupported(assetTypes[y])
				enabledPairs, err := exchanges[x].GetEnabledPairs(assetTypes[y])
				if err != nil {
					log.Errorf(log.CandleSyncMgr,
						"%s failed to get enabled pairs. Err: %s",
						exchangeName,
						err)
					continue
				}
				for i := range enabledPairs {
					if atomic.LoadInt32(&m.started) == 0 {
						return
					}

					c, err := m.get(exchangeName, enabledPairs[i], assetTypes[y])
					if err != nil {
						if err == errCandleSyncPairNotFound {
							c = &candleSyncAgent{
								AssetType: assetTypes[y],
								Exchange:  exchangeName,
								Pair:      enabledPairs[i],
							}

							sBase := syncBase{
								IsUsingREST:      usingREST || !wsAssetSupported,
								IsUsingWebsocket: usingWebsocket && wsAssetSupported,
							}

							if m.config.SyncTicker {
								c.Ticker = sBase
							}

							if m.config.SyncOrderbook {
								c.Orderbook = sBase
							}

							if m.config.SyncTrades {
								c.Trade = sBase
							}

							m.add(c)
						} else {
							log.Error(log.CandleSyncMgr, err)
							continue
						}
					}
					if switchedToRest && usingWebsocket {
						log.Warnf(log.CandleSyncMgr,
							"%s %s: Websocket re-enabled, switching from rest to websocket",
							c.Exchange, m.FormatCurrency(enabledPairs[i]).String())
						switchedToRest = false
					}

					if m.config.SyncTicker {
						if !m.isProcessing(exchangeName, c.Pair, c.AssetType, CandleSyncItemTicker) {
							if c.Ticker.LastUpdated.IsZero() ||
								(time.Since(c.Ticker.LastUpdated) > m.config.SyncTimeoutREST && c.Ticker.IsUsingREST) ||
								(time.Since(c.Ticker.LastUpdated) > m.config.SyncTimeoutWebsocket && c.Ticker.IsUsingWebsocket) {
								if c.Ticker.IsUsingWebsocket {
									if time.Since(c.Created) < m.config.SyncTimeoutWebsocket {
										continue
									}

									if supportsREST {
										m.setProcessing(c.Exchange, c.Pair, c.AssetType, CandleSyncItemTicker, true)
										c.Ticker.IsUsingWebsocket = false
										c.Ticker.IsUsingREST = true
										log.Warnf(log.CandleSyncMgr,
											"%s %s %s: No ticker update after %s, switching from websocket to rest",
											c.Exchange,
											m.FormatCurrency(enabledPairs[i]).String(),
											strings.ToUpper(c.AssetType.String()),
											m.config.SyncTimeoutWebsocket,
										)
										switchedToRest = true
										m.setProcessing(c.Exchange, c.Pair, c.AssetType, CandleSyncItemTicker, false)
									}
								}

								if c.Ticker.IsUsingREST {
									m.setProcessing(c.Exchange, c.Pair, c.AssetType, CandleSyncItemTicker, true)
									var result *ticker.Price
									var err error

									if supportsRESTTickerBatching {
										m.mux.Lock()
										batchLastDone, ok := m.tickerBatchLastRequested[exchangeName]
										if !ok {
											m.tickerBatchLastRequested[exchangeName] = time.Time{}
										}
										m.mux.Unlock()

										if batchLastDone.IsZero() || time.Since(batchLastDone) > m.config.SyncTimeoutREST {
											m.mux.Lock()
											if m.config.Verbose {
												log.Debugf(log.CandleSyncMgr, "Initialising %s REST ticker batching", exchangeName)
											}
											err = exchanges[x].UpdateTickers(context.TODO(), c.AssetType)
											if err == nil {
												result, err = exchanges[x].FetchTicker(context.TODO(), c.Pair, c.AssetType)
											}
											m.tickerBatchLastRequested[exchangeName] = time.Now()
											m.mux.Unlock()
										} else {
											if m.config.Verbose {
												log.Debugf(log.CandleSyncMgr, "%s Using recent batching cache", exchangeName)
											}
											result, err = exchanges[x].FetchTicker(context.TODO(),
												c.Pair,
												c.AssetType)
										}
									} else {
										result, err = exchanges[x].UpdateTicker(context.TODO(),
											c.Pair,
											c.AssetType)
									}
									m.PrintTickerSummary(result, "REST", err)
									if err == nil {
										if m.remoteConfig.WebsocketRPC.Enabled {
											relayWebsocketEvent(result, "ticker_update", c.AssetType.String(), exchangeName)
										}
									}
									updateErr := m.Update(c.Exchange, c.Pair, c.AssetType, CandleSyncItemTicker, err)
									if updateErr != nil {
										log.Error(log.CandleSyncMgr, updateErr)
									}
								}
							} else {
								time.Sleep(time.Millisecond * 50)
							}
						}
					}
				}
			}
		}
	}
}

// PrintTickerSummary outputs the ticker results
func (m *candleSyncManager) PrintTickerSummary(result *ticker.Price, protocol string, err error) {
	if m == nil || atomic.LoadInt32(&m.started) == 0 {
		return
	}
	if err != nil {
		if err == common.ErrNotYetImplemented {
			log.Warnf(log.CandleSyncMgr, "Failed to get %s ticker. Error: %s",
				protocol,
				err)
			return
		}
		log.Errorf(log.CandleSyncMgr, "Failed to get %s ticker. Error: %s",
			protocol,
			err)
		return
	}

	// ignoring error as not all tickers have volume populated and error is not actionable
	_ = stats.Add(result.ExchangeName, result.Pair, result.AssetType, result.Last, result.Volume)

	if result.Pair.Quote.IsFiatCurrency() &&
		result.Pair.Quote != m.fiatDisplayCurrency &&
		!m.fiatDisplayCurrency.IsEmpty() {
		origCurrency := result.Pair.Quote.Upper()
		log.Infof(log.Ticker, "%s %s %s %s: TICKER: Last %s Ask %s Bid %s High %s Low %s Volume %.8f",
			result.ExchangeName,
			protocol,
			m.FormatCurrency(result.Pair),
			strings.ToUpper(result.AssetType.String()),
			printConvertCurrencyFormat(origCurrency, result.Last, m.fiatDisplayCurrency),
			printConvertCurrencyFormat(origCurrency, result.Ask, m.fiatDisplayCurrency),
			printConvertCurrencyFormat(origCurrency, result.Bid, m.fiatDisplayCurrency),
			printConvertCurrencyFormat(origCurrency, result.High, m.fiatDisplayCurrency),
			printConvertCurrencyFormat(origCurrency, result.Low, m.fiatDisplayCurrency),
			result.Volume)
	} else {
		if result.Pair.Quote.IsFiatCurrency() &&
			result.Pair.Quote == m.fiatDisplayCurrency &&
			!m.fiatDisplayCurrency.IsEmpty() {
			log.Infof(log.Ticker, "%s %s %s %s: TICKER: Last %s Ask %s Bid %s High %s Low %s Volume %.8f",
				result.ExchangeName,
				protocol,
				m.FormatCurrency(result.Pair),
				strings.ToUpper(result.AssetType.String()),
				printCurrencyFormat(result.Last, m.fiatDisplayCurrency),
				printCurrencyFormat(result.Ask, m.fiatDisplayCurrency),
				printCurrencyFormat(result.Bid, m.fiatDisplayCurrency),
				printCurrencyFormat(result.High, m.fiatDisplayCurrency),
				printCurrencyFormat(result.Low, m.fiatDisplayCurrency),
				result.Volume)
		} else {
			log.Infof(log.Ticker, "%s %s %s %s: TICKER: Last %.8f Ask %.8f Bid %.8f High %.8f Low %.8f Volume %.8f",
				result.ExchangeName,
				protocol,
				m.FormatCurrency(result.Pair),
				strings.ToUpper(result.AssetType.String()),
				result.Last,
				result.Ask,
				result.Bid,
				result.High,
				result.Low,
				result.Volume)
		}
	}
}

// FormatCurrency is a method that formats and returns a currency pair
// based on the user currency display preferences
func (m *candleSyncManager) FormatCurrency(p currency.Pair) currency.Pair {
	if m == nil || atomic.LoadInt32(&m.started) == 0 {
		return p
	}
	return p.Format(m.delimiter, m.uppercase)
}

// const (
// 	book = "%s %s %s %s: ORDERBOOK: Bids len: %d Amount: %f %s. Total value: %s Asks len: %d Amount: %f %s. Total value: %s"
// )
//
// PrintOrderbookSummary outputs orderbook results
func (m *candleSyncManager) PrintOrderbookSummary(result *orderbook.Base, protocol string, err error) {
	if m == nil || atomic.LoadInt32(&m.started) == 0 {
		return
	}
	if err != nil {
		if result == nil {
			log.Errorf(log.OrderBook, "Failed to get %s orderbook. Error: %s",
				protocol,
				err)
			return
		}
		if err == common.ErrNotYetImplemented {
			log.Warnf(log.OrderBook, "Failed to get %s orderbook for %s %s %s. Error: %s",
				protocol,
				result.Exchange,
				result.Pair,
				result.Asset,
				err)
			return
		}
		log.Errorf(log.OrderBook, "Failed to get %s orderbook for %s %s %s. Error: %s",
			protocol,
			result.Exchange,
			result.Pair,
			result.Asset,
			err)
		return
	}

	bidsAmount, bidsValue := result.TotalBidsAmount()
	asksAmount, asksValue := result.TotalAsksAmount()

	var bidValueResult, askValueResult string
	switch {
	case result.Pair.Quote.IsFiatCurrency() && result.Pair.Quote != m.fiatDisplayCurrency && !m.fiatDisplayCurrency.IsEmpty():
		origCurrency := result.Pair.Quote.Upper()
		bidValueResult = printConvertCurrencyFormat(origCurrency, bidsValue, m.fiatDisplayCurrency)
		askValueResult = printConvertCurrencyFormat(origCurrency, asksValue, m.fiatDisplayCurrency)
	case result.Pair.Quote.IsFiatCurrency() && result.Pair.Quote == m.fiatDisplayCurrency && !m.fiatDisplayCurrency.IsEmpty():
		bidValueResult = printCurrencyFormat(bidsValue, m.fiatDisplayCurrency)
		askValueResult = printCurrencyFormat(asksValue, m.fiatDisplayCurrency)
	default:
		bidValueResult = strconv.FormatFloat(bidsValue, 'f', -1, 64)
		askValueResult = strconv.FormatFloat(asksValue, 'f', -1, 64)
	}

	log.Debugf(log.OrderBook, book,
		result.Exchange,
		protocol,
		m.FormatCurrency(result.Pair),
		strings.ToUpper(result.Asset.String()),
		len(result.Bids),
		bidsAmount,
		result.Pair.Base,
		bidValueResult,
		len(result.Asks),
		asksAmount,
		result.Pair.Base,
		askValueResult,
	)
}

// WaitForInitialCandleSync allows for a routine to wait for an initial sync to be
// completed without exposing the underlying type. This needs to be called in a
// separate routine.
func (m *candleSyncManager) WaitForInitialCandleSync() error {
	if m == nil {
		return fmt.Errorf("sync manager %w", ErrNilSubsystem)
	}

	m.inService.Wait()
	if atomic.LoadInt32(&m.started) == 0 {
		return fmt.Errorf("sync manager %w", ErrSubSystemNotStarted)
	}

	m.initSyncWG.Wait()
	return nil
}
