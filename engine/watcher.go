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
	WatcherName = "watcher"
)

// Watcher watches the system
type Watcher struct {
	started  int32
	shutdown chan struct{}
	wg       sync.WaitGroup
	iExchangeManager
	sleep time.Duration
}

// SetupWatcher applies configuration parameters before running
func SetupWatcher(interval time.Duration, em iExchangeManager) (*Watcher, error) {
	if em == nil {
		return nil, errNilExchangeManager
	}
	var c Watcher
	if interval <= 0 {
		log.Warnf(log.ExchangeSys,
			"Watcher interval is invalid, defaulting to: %s",
			DefaultStateManagerDelay)
		interval = DefaultStateManagerDelay
	}
	c.sleep = interval
	c.iExchangeManager = em
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
	defer c.wg.Done()
	timer := time.NewTimer(0) // Prime firing of channel for initial sync.
	for {
		select {
		case <-c.shutdown:
			return
		case <-timer.C:
			var wg sync.WaitGroup
			exchs, err := c.GetExchanges()
			if err != nil {
				log.Errorf(log.Global,
					"Watcher failed to get exchanges error: %v",
					err)
			}
			for x := range exchs {
				wg.Add(1)
				go c.update(exchs[x], &wg, exchs[x].GetAssetTypes(true))
			}
			wg.Wait() // This causes some variability in the timer due to
			// longest length of request time. Can do time.Ticker but don't
			// want routines to stack behind, this is more uniform.
			timer.Reset(c.sleep)
		}
	}
}

func (c *Watcher) update(exch exchange.IBotExchange, wg *sync.WaitGroup, enabledAssets asset.Items) {
	defer wg.Done()
}
