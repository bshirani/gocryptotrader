package engine

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gocryptotrader/data/kline/database"
	"gocryptotrader/exchange"
	"gocryptotrader/exchange/asset"
	gctkline "gocryptotrader/exchange/kline"
	"gocryptotrader/exchange/trade"
	"gocryptotrader/log"
)

const (
	WatcherName       = "watcher"
	DefaultWatchDelay = time.Second * 5
)

// Watcher watches the system
type Watcher struct {
	tradeManager *TradeManager
	started      int32
	shutdown     chan struct{}
	wg           sync.WaitGroup
	sleep        time.Duration
	bot          *Engine
	queue        EventHolder
}

// SetupWatcher applies configuration parameters before running
func SetupWatcher(interval time.Duration, bot *Engine, tradeManager *TradeManager) (*Watcher, error) {
	var w Watcher
	if interval <= 0 {
		log.Warnf(log.ExchangeSys,
			"Watcher interval is invalid, defaulting to: %s",
			DefaultWatchDelay)
		interval = DefaultWatchDelay
	}
	w.bot = bot
	w.tradeManager = tradeManager
	w.queue = bot.TradeManager.EventQueue
	w.sleep = interval
	w.shutdown = make(chan struct{})
	return &w, nil
}

// Start runs the subsystem
func (w *Watcher) Start() error {
	log.Debugln(log.ExchangeSys, "Watcher starting...")
	if w == nil {
		return fmt.Errorf("%s %w", WatcherName, ErrNilSubsystem)
	}

	if !atomic.CompareAndSwapInt32(&w.started, 0, 1) {
		return fmt.Errorf("%s %w", WatcherName, ErrSubSystemAlreadyStarted)
	}
	w.wg.Add(1)
	go w.monitor()
	log.Debugln(log.ExchangeSys, "Watcher started.")
	return nil
}

// Stop stops the subsystem
func (w *Watcher) Stop() error {
	fmt.Println("trying to stop")
	if w == nil {
		return fmt.Errorf("%s %w", WatcherName, ErrNilSubsystem)
	}
	if atomic.LoadInt32(&w.started) == 0 {
		return fmt.Errorf("%s %w", WatcherName, ErrSubSystemNotStarted)
	}

	log.Debugf(log.ExchangeSys, "Watcher %s", MsgSubSystemShuttingDown)
	close(w.shutdown)
	w.wg.Wait()
	w.shutdown = make(chan struct{})
	log.Debugf(log.ExchangeSys, "Watcher %s", MsgSubSystemShutdown)
	atomic.StoreInt32(&w.started, 0)
	return nil
}

// IsRunning safely checks whether the subsystem is running
func (w *Watcher) IsRunning() bool {
	if w == nil {
		return false
	}
	return atomic.LoadInt32(&w.started) == 1
}

func (w *Watcher) monitor() {
	defer w.wg.Done()
	lup := make(map[*ExchangeAssetPairSettings]time.Time)
	processEventTicker := time.NewTicker(time.Second)
	for {
		select {
		case <-w.shutdown:
			return
		case <-processEventTicker.C:
			t := time.Now()
			thisMinute := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
			// thisMinute := time.Now()
			for _, cs := range w.bot.CurrencySettings {
				t1 := lup[cs]

				if t1 == thisMinute { //skip if alrady updated this minute
					continue
				} else {
					trades, err := trade.GetTradesInRange(
						cs.ExchangeName,
						cs.AssetType.String(),
						cs.CurrencyPair.Base.String(),
						cs.CurrencyPair.Quote.String(),
						time.Now().Add(-time.Minute),
						time.Now())

					if err != nil {
						fmt.Println("unable to retrieve data from GoCryptoTrader database. Error: %v. Please ensure the database is setup correctly and has data before use", err)
						continue
					}

					if len(trades) > 0 {
						lastTradeTime := trades[len(trades)-1].Timestamp
						if lastTradeTime.Before(thisMinute) {
							continue
						}

						lup[cs] = thisMinute

						trades[0].CurrencyPair = cs.CurrencyPair
						trades[0].Exchange = strings.ToLower(cs.ExchangeName)
						klineItem, err := trade.ConvertTradesToCandles(
							gctkline.Interval(gctkline.OneMin),
							trades...)
						if err != nil {
							log.Errorf(log.Watcher, "could not convert database trade data for %v %v %v, %v", cs.ExchangeName, cs.AssetType, cs.CurrencyPair, err)
						}

						// store the candles in the database

						gctkline.StoreInDatabase(&klineItem, false)

						// func LoadData(startDate, endDate time.Time, interval time.Duration, exchangeName string, dataType int64, fPair currency.Pair, a asset.Item) (*kline.DataFromKline, error) {
						dbData, err := database.LoadData(
							thisMinute,
							thisMinute.Add(time.Minute*2),
							time.Minute,
							cs.ExchangeName,
							0,
							cs.CurrencyPair,
							cs.AssetType)

						if err != nil {
							fmt.Println("error load db data", err)
						}

						// saved the trades as candlesticks, now pull the new candlesticks and queue them
						// retCandle, _ := candle.Series(cs.ExchangeName, cs.CurrencyPair.Base.String(), cs.CurrencyPair.Quote.String(), 60, cs.AssetType.String(), time.Now().Add(time.Minute*-5), time.Now())
						// candles := retCandle.Candles
						// fmt.Println("last candle1", candles[len(candles)-1])
						// fmt.Println("last candle2", dbData.Item.Candles[len(dbData.Item.Candles)-1])

						// only if no already set
						// if w.tradeManager.Datas.GetDataForCurrency(cs.ExchangeName, cs.AssetType, cs.CurrencyPair) == nil {
						w.tradeManager.Datas.SetDataForCurrency(cs.ExchangeName, cs.AssetType, cs.CurrencyPair, dbData)
						// }

						dbData.Load()

					}
					if err != nil {
						fmt.Println("error", err)
						continue
					}
				}
			}
		}
	}
}

func (w *Watcher) update(exch exchange.IBotExchange, wg *sync.WaitGroup, enabledAssets asset.Items) {
	defer wg.Done()
}
