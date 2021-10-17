package engine

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gocryptotrader/data/kline"
	"gocryptotrader/database/repository/candle"
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
	// lup := make(map[string]map[asset.Item]map[currency.Pair]time.Time)
	defer w.wg.Done()
	lup := make(map[ExchangeAssetPairSettings]time.Time)
	processEventTicker := time.NewTicker(time.Second)
	for {
		select {
		case <-w.shutdown:
			fmt.Println("received shutdown")
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
					exch, _, asset, err := w.bot.loadExchangePairAssetBase(
						cs.ExchangeName,
						cs.CurrencyPair.Base.String(),
						cs.CurrencyPair.Quote.String(),
						cs.AssetType.String())

					if exch == nil {
						fmt.Println("no exchange found")
					}

					trades, err := trade.GetTradesInRange(
						cs.ExchangeName,
						cs.AssetType.String(),
						cs.CurrencyPair.Base.String(),
						cs.CurrencyPair.Quote.String(),
						time.Now().Add(-time.Minute),
						time.Now())

					candles, _ := candle.Series(cs.ExchangeName,
						cs.CurrencyPair.Base.String(), cs.CurrencyPair.Quote.String(),
						60, cs.AssetType.String(), thisMinute, time.Now())

					if len(candles.Candles) > 0 {
						fmt.Println("candle range", thisMinute, cs.CurrencyPair, candles.Candles[0].Timestamp, candles.Candles[len(candles.Candles)-1].Timestamp)
					} else {
						fmt.Println("no candles", cs.CurrencyPair)
					}

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

						resp := &kline.DataFromKline{}
						resp.Item = gctkline.Item{
							Exchange: strings.ToLower(cs.ExchangeName),
							Pair:     cs.CurrencyPair,
							Interval: gctkline.OneMin,
							Asset:    asset,
						}
						trades[0].CurrencyPair = cs.CurrencyPair
						trades[0].Exchange = strings.ToLower(cs.ExchangeName)
						klineItem, err := trade.ConvertTradesToCandles(
							gctkline.Interval(gctkline.OneMin),
							trades...)
						if err != nil {
							log.Errorf(log.Watcher, "could not convert database trade data for %v %v %v, %v", cs.ExchangeName, cs.AssetType, cs.CurrencyPair, err)
						}

						// resp.Item.Candles = append(resp.Item.Candles, klineItem)
						if len(klineItem.Candles) > 0 {
							fmt.Println("last candle from trades", klineItem.Candles[len(klineItem.Candles)-1])
							klineItem.SortCandlesByTimestamp(true)
							resp.Item = klineItem
							resp.Load()

							// only if no already set
							if w.tradeManager.Datas.GetDataForCurrency(strings.ToLower(cs.ExchangeName), cs.AssetType, cs.CurrencyPair) == nil {
								fmt.Println("setting", &resp)
								w.tradeManager.Datas.SetDataForCurrency(strings.ToLower(cs.ExchangeName), cs.AssetType, cs.CurrencyPair, resp)
							}

						} else {
							fmt.Println("error n o candles")
							os.Exit(123)
						}

						// w.queue.AppendEvent(signals[i])
					}
					if err != nil {
						fmt.Println("error", err)
						continue
					}

					// for ev := w.EventQueue.NextEvent(); ; ev = w.EventQueue.NextEvent() {
					// 	w.EventQueue.AppendEvent(d)
					// }
				}
			}
		}
	}
	//

	// timer := time.NewTimer(0) // Prime firing of channel for initial sync.
	// for {
	// 	select {
	// 		thisMinute := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
	// 		for _, cs := range w.bot.TradeManager.CurrencySettings {
	// 			if lup[cs.ExchangeName] == nil {
	// 				lup[cs.ExchangeName] = make(map[asset.Item]map[currency.Pair]time.Time)
	// 			}
	// 			if lup[cs.ExchangeName][cs.AssetType] == nil {
	// 				lup[cs.ExchangeName][cs.AssetType] = make(map[currency.Pair]time.Time)
	// 			}
	// 			t1 := lup[cs.ExchangeName][cs.AssetType][cs.CurrencyPair]
	//
	// 			if t1 == thisMinute {
	// 				continue
	// 			} else {
	// 				// fmt.Println("updating", cs.CurrencyPair, t1, thisMinute)
	// 				// secondsPast := t1.Sub(t).Seconds()
	// 				trades, err := trade.GetTradesInRange(
	// 					cs.ExchangeName,
	// 					cs.AssetType.String(),
	// 					cs.CurrencyPair.Base.String(),
	// 					cs.CurrencyPair.Quote.String(),
	// 					time.Now().Add(-time.Minute),
	// 					time.Now())
	// 				if err != nil {
	// 					log.Errorf(log.Watcher, "could not retrieve database trade data for %v %v %v, %v", cs.ExchangeName, cs.AssetType, cs.CurrencyPair, err)
	// 				}
	// 				if len(trades) > 0 {
	// 					klineItem, err := trade.ConvertTradesToCandles(
	// 						gctkline.Interval(kline.OneMin),
	// 						trades...)
	// 					if err != nil {
	// 						log.Errorf(log.Watcher, "could not convert database trade data for %v %v %v, %v", cs.ExchangeName, cs.AssetType, cs.CurrencyPair, err)
	// 					}
	//
	// 					klineItem.SortCandlesByTimestamp(true)
	// 					fmt.Println("I have a bar for you")
	//
	// 					klineItem.Load()
	// 					// w.queue.AppendEvent(signals[i])
	// 					// startCandle := klineItem.Candles[0].Time
	// 					// lastCandle := klineItem.Candles[len(klineItem.Candles)-1].Time
	// 					// totalSeconds := startCandle.Sub(lastCandle).Seconds()
	// 					// fmt.Println(cs.CurrencyPair, len(klineItem.Candles), "candles", totalSeconds, "seconds", lastCandle, startCandle)
	// 					// if len(klineItem.Candles) > 0 {
	// 					// 	for _, candle := range klineItem.Candles {
	// 					// 		fmt.Println(candle.Time)
	// 					// 	}
	// 					// }
	// 				}
	//
	// 				lastMinuteUpdated[cs.ExchangeName][cs.AssetType][cs.CurrencyPair] = thisMinute
	// 			}
	// 		}
	//
	// 		timer.Reset(w.sleep)
	// 	}
	// }
}

func (w *Watcher) update(exch exchange.IBotExchange, wg *sync.WaitGroup, enabledAssets asset.Items) {
	defer wg.Done()
}
