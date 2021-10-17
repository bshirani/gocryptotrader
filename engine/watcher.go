package engine

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"gocryptotrader/exchange"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/log"
)

const (
	WatcherName       = "watcher"
	DefaultWatchDelay = time.Second * 5
)

// Watcher watches the system
type Watcher struct {
	started  int32
	shutdown chan struct{}
	wg       sync.WaitGroup
	sleep    time.Duration
	bot      *Engine
	queue    EventHolder
}

// SetupWatcher applies configuration parameters before running
func SetupWatcher(interval time.Duration, bot *Engine) (*Watcher, error) {
	var c Watcher
	if interval <= 0 {
		log.Warnf(log.ExchangeSys,
			"Watcher interval is invalid, defaulting to: %s",
			DefaultWatchDelay)
		interval = DefaultWatchDelay
	}
	c.bot = bot
	c.queue = bot.TradeManager.EventQueue
	c.sleep = interval
	c.shutdown = make(chan struct{})
	return &c, nil
}

// Start runs the subsystem
func (c *Watcher) Start() error {
	log.Debugln(log.ExchangeSys, "Watcher starting...")
	if c == nil {
		return fmt.Errorf("%s %w", WatcherName, ErrNilSubsystem)
	}

	if !atomic.CompareAndSwapInt32(&c.started, 0, 1) {
		return fmt.Errorf("%s %w", WatcherName, ErrSubSystemAlreadyStarted)
	}
	c.wg.Add(1)
	go c.monitor()
	log.Debugln(log.ExchangeSys, "Watcher started.")
	return nil
}

// Stop stops the subsystem
func (c *Watcher) Stop() error {
	if c == nil {
		return fmt.Errorf("%s %w", WatcherName, ErrNilSubsystem)
	}
	if atomic.LoadInt32(&c.started) == 0 {
		return fmt.Errorf("%s %w", WatcherName, ErrSubSystemNotStarted)
	}

	log.Debugf(log.ExchangeSys, "Watcher %s", MsgSubSystemShuttingDown)
	close(c.shutdown)
	c.wg.Wait()
	c.shutdown = make(chan struct{})
	log.Debugf(log.ExchangeSys, "Watcher %s", MsgSubSystemShutdown)
	atomic.StoreInt32(&c.started, 0)
	return nil
}

// IsRunning safely checks whether the subsystem is running
func (c *Watcher) IsRunning() bool {
	if c == nil {
		return false
	}
	return atomic.LoadInt32(&c.started) == 1
}

func (c *Watcher) monitor() {
	// lup := make(map[string]map[asset.Item]map[currency.Pair]time.Time)
	//
	// defer c.wg.Done()
	// timer := time.NewTimer(0) // Prime firing of channel for initial sync.
	// for {
	// 	select {
	// 	case <-c.shutdown:
	// 		return
	// 	case <-timer.C:
	// 		t := time.Now()
	// 		thisMinute := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
	//
	// 		for _, cs := range c.bot.TradeManager.CurrencySettings {
	//
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
	// 					// c.queue.AppendEvent(signals[i])
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
	// 		timer.Reset(c.sleep)
	// 	}
	// }
}

func (c *Watcher) update(exch exchange.IBotExchange, wg *sync.WaitGroup, enabledAssets asset.Items) {
	defer wg.Done()
}
