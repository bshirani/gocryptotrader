package engine

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gocryptotrader/common"
	"gocryptotrader/communications/base"
	"gocryptotrader/currency"
	"gocryptotrader/database/repository/liveorder"
	"gocryptotrader/eventtypes"
	"gocryptotrader/exchange"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/exchange/order"
	"gocryptotrader/log"

	"github.com/shopspring/decimal"
)

// SetupOrderManager will boot up the OrderManager
func SetupOrderManager(exchangeManager iExchangeManager, communicationsManager iCommsManager, wg *sync.WaitGroup, verbose bool, realOrders bool, liveMode bool, dryRun bool) (*OrderManager, error) {
	if exchangeManager == nil {
		return nil, errNilExchangeManager
	}
	if communicationsManager == nil {
		return nil, errNilCommunicationsManager
	}
	if wg == nil {
		return nil, errNilWaitGroup
	}

	// load orders from database
	// activeOrders, _ := liveorder.Active()
	// for _, t := range activeOrders {
	// 	p.store.openOrders[t.StrategyID] = append(p.store.openOrders[t.StrategyID], &t)
	// }

	return &OrderManager{
		shutdown: make(chan struct{}),
		orderStore: store{
			Orders:          make(map[string][]*order.Detail),
			exchangeManager: exchangeManager,
			commsManager:    communicationsManager,
			wg:              wg,
			dryRun:          dryRun,
		},
		realOrders:    realOrders,
		verbose:       verbose,
		liveMode:      liveMode,
		dryRun:        dryRun,
		currentCloses: make(map[string]map[asset.Item]map[currency.Pair]decimal.Decimal),
	}, nil
}

// IsRunning safely checks whether the subsystem is running
func (m *OrderManager) IsRunning() bool {
	if m == nil {
		return false
	}
	return atomic.LoadInt32(&m.started) == 1
}

// Start runs the subsystem
func (m *OrderManager) Start() error {
	if m == nil {
		return fmt.Errorf("order manager %w", ErrNilSubsystem)
	}
	if !atomic.CompareAndSwapInt32(&m.started, 0, 1) {
		return fmt.Errorf("order manager %w", ErrSubSystemAlreadyStarted)
	}
	log.Debugln(log.OrderMgr, "Order manager starting...")
	m.shutdown = make(chan struct{})
	if m.realOrders {
		go m.run()
	}
	return nil
}

func (m *OrderManager) Update() {
	// fmt.Println("updating order manager from trade manager")
}

// Stop attempts to shutdown the subsystem
func (m *OrderManager) Stop() error {
	if m == nil {
		return fmt.Errorf("order manager %w", ErrNilSubsystem)
	}
	if atomic.LoadInt32(&m.started) == 0 {
		return fmt.Errorf("order manager %w", ErrSubSystemNotStarted)
	}

	defer func() {
		atomic.CompareAndSwapInt32(&m.started, 1, 0)
	}()

	log.Debugln(log.OrderMgr, "Order manager shutting down...")
	close(m.shutdown)
	return nil
}

// gracefulShutdown cancels all orders (if enabled) before shutting down
func (m *OrderManager) gracefulShutdown() {
	if m.cfg.CancelOrdersOnShutdown {
		log.Debugln(log.OrderMgr, "Order manager: Cancelling any open orders...")
		exchanges, err := m.orderStore.exchangeManager.GetExchanges()
		if err != nil {
			log.Errorf(log.OrderMgr, "Order manager cannot get exchanges: %v", err)
			return
		}
		m.CancelAllOrders(context.TODO(), exchanges)
	}
}

// run will periodically process orders
func (m *OrderManager) UpdateFakeOrders(d eventtypes.DataEventHandler) error {
	if m.realOrders {
		panic("updating fake orders in production")
	}
	if m.currentCloses[d.GetExchange()] == nil {
		m.currentCloses[d.GetExchange()] = make(map[asset.Item]map[currency.Pair]decimal.Decimal)
	}
	if m.currentCloses[d.GetExchange()][d.GetAssetType()] == nil {
		m.currentCloses[d.GetExchange()][d.GetAssetType()] = make(map[currency.Pair]decimal.Decimal)
	}

	m.currentCloses[d.GetExchange()][d.GetAssetType()][d.Pair()] = d.ClosePrice()

	active, _ := m.GetOrdersActive(nil)
	for _, ao := range active {
		if ao.Pair != d.Pair() {
			continue
		}
		// handle stop orders only
		if ao.Type != order.Stop {
			continue
		}
		// if ao.Side == order.Sell {
		// 	fmt.Println("update fake",
		// 		len(active),
		// 		d.Pair(),
		// 		ao.Type,
		// 		ao.Side,
		// 		"pts2close=",
		// 		d.ClosePrice().Sub(decimal.NewFromFloat(ao.Price)),
		// 	)
		// }

		if (ao.Side == order.Sell && d.ClosePrice().LessThan(decimal.NewFromFloat(ao.Price))) ||
			(ao.Side == order.Buy && d.ClosePrice().GreaterThan(decimal.NewFromFloat(ao.Price))) {
			ao.Status = order.Filled
			ao.InternalType = order.InternalStopLoss
			_, err := m.orderStore.upsert(&ao)
			if err != nil {
				fmt.Println("error upserting stop order")
			}
			active, _ = m.GetOrdersActive(nil)
			m.onFill(ao, d)
		}
	}
	return nil
}

// run will periodically process orders
func (m *OrderManager) run() {
	log.Debugln(log.OrderMgr, "Order manager started.")
	m.processOrders()
	tick := time.NewTicker(orderManagerDelay)
	m.orderStore.wg.Add(1)
	defer func() {
		log.Debugln(log.OrderMgr, "Order manager shutdown.")
		tick.Stop()
		m.orderStore.wg.Done()
	}()

	for {
		select {
		case <-m.shutdown:
			m.gracefulShutdown()
			return
		case <-tick.C:
			go m.processOrders()
		}
	}
}

// CancelAllOrders iterates and cancels all orders for each exchange provided
func (m *OrderManager) CancelAllOrders(ctx context.Context, exchangeNames []exchange.IBotExchange) {
	if m == nil || atomic.LoadInt32(&m.started) == 0 {
		return
	}

	orders := m.orderStore.get()
	if orders == nil {
		return
	}

	for i := range exchangeNames {
		exchangeOrders, ok := orders[strings.ToLower(exchangeNames[i].GetName())]
		if !ok {
			continue
		}
		for j := range exchangeOrders {
			log.Debugf(log.OrderMgr, "Order manager: Cancelling order(s) for exchange %s.", exchangeNames[i].GetName())
			err := m.Cancel(ctx, &order.Cancel{
				Exchange:      exchangeOrders[j].Exchange,
				ID:            exchangeOrders[j].ID,
				AccountID:     exchangeOrders[j].AccountID,
				ClientID:      exchangeOrders[j].ClientID,
				WalletAddress: exchangeOrders[j].WalletAddress,
				Type:          exchangeOrders[j].Type,
				Side:          exchangeOrders[j].Side,
				Pair:          exchangeOrders[j].Pair,
				AssetType:     exchangeOrders[j].AssetType,
			})
			if err != nil {
				log.Error(log.OrderMgr, err)
			}
		}
	}
}

// Cancel will find the order in the OrderManager, send a cancel request
// to the exchange and if successful, update the status of the order
func (m *OrderManager) Cancel(ctx context.Context, cancel *order.Cancel) error {
	if m == nil {
		return fmt.Errorf("order manager %w", ErrNilSubsystem)
	}
	if atomic.LoadInt32(&m.started) == 0 {
		return fmt.Errorf("order manager %w", ErrSubSystemNotStarted)
	}
	var err error
	defer func() {
		if err != nil {
			m.orderStore.commsManager.PushEvent(base.Event{
				Type:    "order",
				Message: err.Error(),
			})
		}
	}()

	if cancel == nil {
		err = errors.New("order cancel param is nil")
		return err
	}
	if cancel.Exchange == "" {
		err = errors.New("order exchange name is empty")
		return err
	}
	if cancel.ID == "" {
		err = errors.New("order id is empty")
		return err
	}

	if m.realOrders {
		panic("trying to cancel for real")
		exch, err := m.orderStore.exchangeManager.GetExchangeByName(cancel.Exchange)
		if err != nil {
			return err
		}

		if cancel.AssetType.String() != "" && !exch.GetAssetTypes(false).Contains(cancel.AssetType) {
			err = errors.New("order asset type not supported by exchange")
			return err
		}

		log.Debugf(log.OrderMgr, "Order manager: Cancelling order ID %v [%+v]",
			cancel.ID, cancel)

		err = exch.CancelOrder(ctx, cancel)
		if err != nil {
			err = fmt.Errorf("%v - Failed to cancel order: %w", cancel.Exchange, err)
			return err
		}
	}

	if cancel.Date.IsZero() {
		panic("cannot save cancel without date")
	}
	var od *order.Detail
	od, err = m.orderStore.getByExchangeAndID(cancel.Exchange, cancel.ID)
	if err != nil {
		err = fmt.Errorf("%v - Failed to retrieve order %v to update cancelled status: %w", cancel.Exchange, cancel.ID, err)
		return err
	}

	od.Status = order.Cancelled
	od.CancelledAt = cancel.Date
	msg := fmt.Sprintf("Order manager: Exchange %s order ID=%v cancelled.",
		od.Exchange, od.ID)

	if m.verbose {
		log.Debugln(log.OrderMgr, msg)
	}

	if !m.dryRun {
		id, err := liveorder.Upsert(od)
		// fmt.Println("lookup order id", id, "oid", od.ID, "internalid", od.InternalOrderID)
		dblo, err := liveorder.OneByID(id)
		if dblo.ID == 0 || err != nil {
			fmt.Println("database not recorded as cancelled")
			panic(err)
		}
		if dblo.Status != order.Cancelled {
			panic(fmt.Sprintf("db status not cancelled %s", dblo.Status))
		}

		// ensure that order is cancelled in the databasek
		if err != nil {
			fmt.Println("error upserting cancelled order", err)
			panic(err)
		}
		// m.orderStore.commsManager.PushEvent(base.Event{
		// 	Type:    "order",
		// 	Message: msg,
		// })
	}
	return nil
}

// GetOrderInfo calls the exchange's wrapper GetOrderInfo function
// and stores the result in the order manager
func (m *OrderManager) GetOrderInfo(ctx context.Context, exchangeName, orderID string, cp currency.Pair, a asset.Item) (order.Detail, error) {
	if m == nil {
		return order.Detail{}, fmt.Errorf("order manager %w", ErrNilSubsystem)
	}
	if atomic.LoadInt32(&m.started) == 0 {
		return order.Detail{}, fmt.Errorf("order manager %w", ErrSubSystemNotStarted)
	}

	if orderID == "" {
		return order.Detail{}, ErrOrderIDCannotBeEmpty
	}

	exch, err := m.orderStore.exchangeManager.GetExchangeByName(exchangeName)
	if err != nil {
		return order.Detail{}, err
	}
	result, err := exch.GetOrderInfo(ctx, orderID, cp, a)
	if err != nil {
		return order.Detail{}, err
	}

	upsertResponse, err := m.orderStore.upsert(&result)
	if err != nil {
		return order.Detail{}, err
	}

	return upsertResponse.OrderDetails, nil
}

// validate ensures a submitted order is valid before adding to the manager
func (m *OrderManager) validate(newOrder *order.Submit) error {
	if newOrder == nil {
		return errors.New("order cannot be nil")
	}
	if newOrder.Exchange == "" {
		return errors.New("order exchange name must be specified")
	}
	if err := newOrder.Validate(); err != nil {
		return fmt.Errorf("order manager: %w", err)
	}
	if newOrder.StrategyID == 0 {
		panic("order without strategy")
	}
	if newOrder.StrategyName == "" {
		panic("order without strategy name")
	}
	if newOrder.InternalType == "" {
		panic("order without internal type")
	}
	// else {
	// 	fmt.Println("internal order type is ", newOrder.InternalType)
	// }

	if m.cfg.EnforceLimitConfig {
		if !m.cfg.AllowMarketOrders && newOrder.Type == order.Market {
			return errors.New("order market type is not allowed")
		}

		if m.cfg.LimitAmount > 0 && newOrder.Amount > m.cfg.LimitAmount {
			return errors.New("order limit exceeds allowed limit")
		}

		if len(m.cfg.AllowedExchanges) > 0 &&
			!common.StringDataCompareInsensitive(m.cfg.AllowedExchanges, newOrder.Exchange) {
			return errors.New("order exchange not found in allowed list")
		}

		if len(m.cfg.AllowedPairs) > 0 && !m.cfg.AllowedPairs.Contains(newOrder.Pair, true) {
			return errors.New("order pair not found in allowed list")
		}
	}
	return nil
}

// Modify depends on the order.Modify.ID and order.Modify.Exchange fields to uniquely
// identify an order to modify.
func (m *OrderManager) Modify(ctx context.Context, mod *order.Modify) (*order.ModifyResponse, error) {
	if m == nil {
		return nil, fmt.Errorf("order manager %w", ErrNilSubsystem)
	}
	if atomic.LoadInt32(&m.started) == 0 {
		return nil, fmt.Errorf("order manager %w", ErrSubSystemNotStarted)
	}

	// Fetch details from locally managed order store.
	det, err := m.orderStore.getByExchangeAndID(mod.Exchange, mod.ID)
	if det == nil || err != nil {
		return nil, fmt.Errorf("order does not exist: %w", err)
	}

	// Populate additional Modify fields as some of them are required by various
	// exchange implementations.
	mod.Pair = det.Pair                           // Used by Bithumb.
	mod.Side = det.Side                           // Used by Bithumb.
	mod.PostOnly = det.PostOnly                   // Used by Poloniex.
	mod.ImmediateOrCancel = det.ImmediateOrCancel // Used by Poloniex.

	// Following is just a precaution to not modify orders by mistake if exchange
	// implementations do not check fields of the Modify struct for zero values.
	if mod.Amount == 0 {
		mod.Amount = det.Amount
	}
	if mod.Price == 0 {
		mod.Price = det.Price
	}

	// Get exchange instance and submit order modification request.
	exch, err := m.orderStore.exchangeManager.GetExchangeByName(mod.Exchange)
	if err != nil {
		return nil, err
	}
	res, err := exch.ModifyOrder(ctx, mod)
	if err != nil {
		message := fmt.Sprintf(
			"Order manager: Exchange %s order ID=%v: failed to modify",
			mod.Exchange,
			mod.ID,
		)
		m.orderStore.commsManager.PushEvent(base.Event{
			Type:    "order",
			Message: message,
		})
		return nil, err
	}

	// If modification is successful, apply changes to local order store.
	//
	// XXX: This comes with a race condition, because [request -> changes] are not
	// atomic.
	err = m.orderStore.modifyExisting(mod.ID, &res)

	// Notify observers.
	var message string
	if err != nil {
		message = "Order manager: Exchange %s order ID=%v: modified on exchange, but failed to modify locally"
	} else {
		message = "Order manager: Exchange %s order ID=%v: modified successfully"
	}
	m.orderStore.commsManager.PushEvent(base.Event{
		Type:    "order",
		Message: fmt.Sprintf(message, mod.Exchange, res.ID),
	})
	return &order.ModifyResponse{OrderID: res.ID}, err
}

// Submit will take in an order struct, send it to the exchange and
// populate it in the OrderManager if successful
func (m *OrderManager) Submit(ctx context.Context, newOrder *order.Submit) (*OrderSubmitResponse, error) {
	// if m.liveMode {
	if m.debug {
		fmt.Println("submitting order type:", newOrder.Type)
	}

	if m.liveMode {
		log.Warnln(log.OrderMgr, "Order manager: Order", newOrder.Side, newOrder.Date, newOrder.StrategyID, newOrder.ID)
	}
	// }

	if m == nil {
		return nil, fmt.Errorf("order manager %w", ErrNilSubsystem)
	}
	if atomic.LoadInt32(&m.started) == 0 {
		return nil, fmt.Errorf("order manager %w", ErrSubSystemNotStarted)
	}

	err := m.validate(newOrder)
	if err != nil {
		return nil, err
	}
	exch, err := m.orderStore.exchangeManager.GetExchangeByName(newOrder.Exchange)
	if err != nil {
		return nil, err
	}

	// Checks for exchange min max limits for order amounts before order
	// execution can occur
	err = exch.CheckOrderExecutionLimits(newOrder.AssetType,
		newOrder.Pair,
		newOrder.Price,
		newOrder.Amount,
		newOrder.Type)
	if err != nil {
		return nil, fmt.Errorf("order manager: exchange %s unable to place order: %w",
			newOrder.Exchange,
			err)
	}

	// Determines if current trading activity is turned off by the exchange for
	// the currency pair
	err = exch.CanTradePair(newOrder.Pair, newOrder.AssetType)
	if err != nil {
		return nil, fmt.Errorf("order manager: exchange %s cannot trade pair %s %s: %w",
			newOrder.Exchange,
			newOrder.Pair,
			newOrder.AssetType,
			err)
	}

	// retrieve order from db
	// if dry run, skip and generate fake id
	var id int
	if !m.dryRun {
		// retrieve the already created order in the database
		// fail if you can't find it
		// in case system fails after creating orders
		// we can recover using the orders in the databse
		id = newOrder.InternalOrderID
		lo, err := liveorder.OneByID(id)
		if lo.ID != id || err != nil {
			panic(err)
		}
	} else {
		id = m.GenerateDryRunID()
	}
	newOrder.InternalOrderID = id

	if newOrder.InternalType == "" {
		panic("order has no internal type")
	}

	var result order.SubmitResponse

	var isOrderFilled bool

	// fmt.Println("order for strategy", newOrder.StrategyID)
	if m.realOrders {
		exch.GetBase().Verbose = true
		result, err = exch.SubmitOrder(ctx, newOrder)
		exch.GetBase().Verbose = false
		if err != nil {
			return nil, err
		}
	} else {
		if newOrder.Type == order.Stop {
			isOrderFilled = false
		} else {
			isOrderFilled = true
			currentPrice, _ := m.currentCloses[newOrder.Exchange][newOrder.AssetType][newOrder.Pair].Float64()
			newOrder.Price = currentPrice
		}
		if newOrder.Price == 0 {
			panic("did not get current price for fake order")
		}

		result = order.SubmitResponse{
			IsOrderPlaced:   true,
			Rate:            newOrder.Price,
			OrderID:         randString(12),
			FullyMatched:    isOrderFilled,
			InternalOrderID: id,
		}
	}

	return m.processSubmittedOrder(newOrder, result)

}

// GetOrdersSnapshot returns a snapshot of all orders in the orderstore. It optionally filters any orders that do not match the status
// but a status of "" or ANY will include all
// the time adds contexts for the when the snapshot is relevant for
func (m *OrderManager) GetOrdersSnapshot(s order.Status) ([]order.Detail, time.Time) {
	if m == nil || atomic.LoadInt32(&m.started) == 0 {
		return nil, time.Time{}
	}
	var os []order.Detail
	var latestUpdate time.Time
	for _, v := range m.orderStore.Orders {
		for i := range v {
			if s != v[i].Status &&
				s != order.AnyStatus &&
				s != "" {
				continue
			}
			if v[i].LastUpdated.After(latestUpdate) {
				latestUpdate = v[i].LastUpdated
			}
			os = append(os, *v[i])
		}
	}

	return os, latestUpdate
}

// GetOrdersFiltered returns a snapshot of all orders in the order store.
// Filtering is applied based on the order.Filter unless entries are empty
func (m *OrderManager) GetOrdersFiltered(f *order.Filter) ([]order.Detail, error) {
	if m == nil {
		return nil, fmt.Errorf("order manager %w", ErrNilSubsystem)
	}
	if f == nil {
		return nil, fmt.Errorf("order manager, GetOrdersFiltered: Filter is nil")
	}
	if atomic.LoadInt32(&m.started) == 0 {
		return nil, fmt.Errorf("order manager %w", ErrSubSystemNotStarted)
	}
	return m.orderStore.getFilteredOrders(f)
}

// GetOrdersActive returns a snapshot of all orders in the order store
// that have a status that indicates it's currently tradable
func (m *OrderManager) GetOrdersActive(f *order.Filter) ([]order.Detail, error) {
	if m == nil {
		return nil, fmt.Errorf("order manager %w", ErrNilSubsystem)
	}
	if atomic.LoadInt32(&m.started) == 0 {
		return nil, fmt.Errorf("order manager %w", ErrSubSystemNotStarted)
	}
	return m.orderStore.getActiveOrders(f), nil
}

// processSubmittedOrder adds a new order to the manager
func (m *OrderManager) processSubmittedOrder(newOrder *order.Submit, result order.SubmitResponse) (*OrderSubmitResponse, error) {
	if !result.IsOrderPlaced {
		return nil, errors.New("order unable to be placed")
	}

	if newOrder.Date.IsZero() {
		newOrder.Date = time.Now()
	}

	if newOrder.StrategyID == 0 {
		return nil, errors.New("order must have a strategy")
	}

	// formatting for backtest only
	// odate := newOrder.Date
	// dtime := fmt.Sprintf("%d-%02d-%02d %d:%02d", odate.Year(), odate.Month(), odate.Day(), odate.Hour(), odate.Minute())
	// fmt.Println(msgInfo)

	if m.verbose {
		msg := fmt.Sprintf("Order manager: Strategy=%d Exchange=%s submitted order ID=%v [Ours: %v] pair=%v price=%v amount=%v side=%v type=%v for time %v.",
			newOrder.StrategyID,
			newOrder.Exchange,
			result.OrderID,
			newOrder.ID,
			newOrder.Pair,
			newOrder.Price,
			newOrder.Amount,
			newOrder.Side,
			newOrder.Type,
			newOrder.Date)
		log.Debugln(log.OrderMgr, msg)
	}

	// msgInfo := fmt.Sprintf("%v %s %-10s %v %-5v %10v",
	// 	newOrder.Pair,
	// 	dtime,
	// 	newOrder.StrategyID,
	// 	newOrder.Pair,
	// 	newOrder.Side,
	// 	newOrder.Price)
	// m.orderStore.commsManager.PushEvent(base.Event{
	// 	Type:    "order",
	// 	Message: msgInfo,
	// })
	status := order.New
	var filledAt time.Time
	if result.FullyMatched {
		status = order.Filled
		filledAt = newOrder.Date
		if filledAt.IsZero() {
			panic("filled at cannot be empty")
		}
	} else {
		status = order.Active
	}
	if newOrder.Price == 0 {
		panic("new order doesnt have a price")
	}
	if newOrder.InternalType == "" {
		panic("no internal type")
	}
	if status != order.Active && status != order.Filled {
		panic("did not fill or submit the order")
	}

	err := m.orderStore.add(&order.Detail{
		InternalType:      newOrder.InternalType,
		Status:            status,
		FilledAt:          filledAt,
		ImmediateOrCancel: newOrder.ImmediateOrCancel,
		HiddenOrder:       newOrder.HiddenOrder,
		FillOrKill:        newOrder.FillOrKill,
		PostOnly:          newOrder.PostOnly,
		Price:             newOrder.Price,
		Amount:            newOrder.Amount,
		LimitPriceUpper:   newOrder.LimitPriceUpper,
		LimitPriceLower:   newOrder.LimitPriceLower,
		TriggerPrice:      newOrder.TriggerPrice,
		TargetAmount:      newOrder.TargetAmount,
		ExecutedAmount:    newOrder.ExecutedAmount,
		RemainingAmount:   newOrder.RemainingAmount,
		Fee:               newOrder.Fee,
		Exchange:          newOrder.Exchange,
		InternalOrderID:   result.InternalOrderID,
		ID:                result.OrderID,
		AccountID:         newOrder.AccountID,
		ClientID:          newOrder.ClientID,
		ClientOrderID:     newOrder.ClientOrderID,
		WalletAddress:     newOrder.WalletAddress,
		Type:              newOrder.Type,
		Side:              newOrder.Side,
		AssetType:         newOrder.AssetType,
		Date:              newOrder.Date,
		LastUpdated:       newOrder.Date,
		Pair:              newOrder.Pair,
		Leverage:          newOrder.Leverage,
		StopLossPrice:     newOrder.StopLossPrice,
		StrategyID:        newOrder.StrategyID,
		StrategyName:      newOrder.StrategyName,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to add %v order %v to orderStore: %s", newOrder.Exchange, result.OrderID, err)
	}

	if result.Rate == 0 {
		panic("order submit response without rate")
	}

	return &OrderSubmitResponse{
		SubmitResponse:  result,
		InternalOrderID: result.InternalOrderID,
		StrategyID:      newOrder.StrategyID,
	}, nil
}

// processOrders iterates over all exchange orders via API
// and adds them to the internal order store
func (m *OrderManager) processOrders() {
	if !atomic.CompareAndSwapInt32(&m.processingOrders, 0, 1) {
		return
	}
	defer func() {
		atomic.StoreInt32(&m.processingOrders, 0)
	}()

	exchanges, err := m.orderStore.exchangeManager.GetExchanges()
	if err != nil {
		log.Errorf(log.OrderMgr, "Order manager cannot get exchanges: %v", err)
		return
	}
	var wg sync.WaitGroup
	for i := range exchanges {
		if !exchanges[i].GetAuthenticatedAPISupport(exchange.RestAuthentication) {
			continue
		}
		if m.verbose {
			log.Debugf(log.OrderMgr,
				"Order manager: Processing orders for exchange %v.",
				exchanges[i].GetName())
		}

		supportedAssets := exchanges[i].GetAssetTypes(true)
		for y := range supportedAssets {
			pairs, err := exchanges[i].GetEnabledPairs(supportedAssets[y])
			if err != nil {
				log.Errorf(log.OrderMgr,
					"Order manager: Unable to get enabled pairs for %s and asset type %s: %s",
					exchanges[i].GetName(),
					supportedAssets[y],
					err)
				continue
			}

			if len(pairs) == 0 {
				if m.verbose {
					log.Debugf(log.OrderMgr,
						"Order manager: No pairs enabled for %s and asset type %s, skipping...",
						exchanges[i].GetName(),
						supportedAssets[y])
				}
				continue
			}

			filter := &order.Filter{
				Exchange: exchanges[i].GetName(),
			}
			orders := m.orderStore.getActiveOrders(filter)
			order.FilterOrdersByCurrencies(&orders, pairs)
			requiresProcessing := make(map[int]bool, len(orders))
			for x := range orders {
				requiresProcessing[orders[x].InternalOrderID] = true
			}

			req := order.GetOrdersRequest{
				Side:      order.AnySide,
				Type:      order.AnyType,
				Pairs:     pairs,
				AssetType: supportedAssets[y],
			}
			result, err := exchanges[i].GetActiveOrders(context.TODO(), &req)
			if err != nil {
				log.Errorf(log.OrderMgr,
					"Order manager: Unable to get active orders for %s and asset type %s: %s",
					exchanges[i].GetName(),
					supportedAssets[y],
					err)
				continue
			}
			if m.verbose {
				log.Infoln(log.OrderMgr, "open orders in store:", len(orders), "from broker:", len(result))
			}
			if len(orders) == 0 && len(result) == 0 {
				continue
			}

			for z := range result {
				upsertResponse, err := m.UpsertOrder(&result[z])
				if err != nil {
					log.Error(log.OrderMgr, err)
				}
				requiresProcessing[upsertResponse.OrderDetails.InternalOrderID] = false
			}
			if !exchanges[i].GetBase().GetSupportedFeatures().RESTCapabilities.GetOrder {
				continue
			}
			wg.Add(1)
			go m.processMatchingOrders(exchanges[i], orders, requiresProcessing, &wg)
		}
	}
	wg.Wait()
}

func (m *OrderManager) processMatchingOrders(exch exchange.IBotExchange, orders []order.Detail, requiresProcessing map[int]bool, wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	for x := range orders {
		if time.Since(orders[x].LastUpdated) < time.Minute {
			continue
		}
		if requiresProcessing[orders[x].InternalOrderID] {
			err := m.FetchAndUpdateExchangeOrder(exch, &orders[x], orders[x].AssetType)
			if err != nil {
				log.Error(log.OrderMgr, err)
			}
		}
	}
}

// FetchAndUpdateExchangeOrder calls the exchange to upsert an order to the order store
func (m *OrderManager) FetchAndUpdateExchangeOrder(exch exchange.IBotExchange, ord *order.Detail, assetType asset.Item) error {
	if ord == nil {
		return errors.New("order manager: Order is nil")
	}
	fetchedOrder, err := exch.GetOrderInfo(context.TODO(), ord.ID, ord.Pair, assetType)
	if err != nil {
		ord.Status = order.UnknownStatus
		return err
	}
	fetchedOrder.LastUpdated = time.Now()
	_, err = m.UpsertOrder(&fetchedOrder)
	return err
}

// Exists checks whether an order exists in the order store
func (m *OrderManager) Exists(o *order.Detail) bool {
	if m == nil || atomic.LoadInt32(&m.started) == 0 {
		return false
	}

	return m.orderStore.exists(o)
}

// Add adds an order to the orderstore
func (m *OrderManager) Add(o *order.Detail) error {
	if m == nil {
		return fmt.Errorf("order manager %w", ErrNilSubsystem)
	}
	if atomic.LoadInt32(&m.started) == 0 {
		return fmt.Errorf("order manager %w", ErrSubSystemNotStarted)
	}

	return m.orderStore.add(o)
}

// GetByExchangeAndID returns a copy of an order from an exchange if it matches the ID
func (m *OrderManager) GetByExchangeAndID(exchangeName, id string) (*order.Detail, error) {
	if m == nil {
		return nil, fmt.Errorf("order manager %w", ErrNilSubsystem)
	}
	if atomic.LoadInt32(&m.started) == 0 {
		return nil, fmt.Errorf("order manager %w", ErrSubSystemNotStarted)
	}

	o, err := m.orderStore.getByExchangeAndID(exchangeName, id)
	if err != nil {
		return nil, err
	}
	var cpy order.Detail
	cpy.UpdateOrderFromDetail(o)
	return &cpy, nil
}

// UpdateExistingOrder will update an existing order in the orderstore
func (m *OrderManager) UpdateExistingOrder(od *order.Detail) error {
	if m == nil {
		return fmt.Errorf("order manager %w", ErrNilSubsystem)
	}
	if atomic.LoadInt32(&m.started) == 0 {
		return fmt.Errorf("order manager %w", ErrSubSystemNotStarted)
	}
	return m.orderStore.updateExisting(od)
}

// UpsertOrder updates an existing order or adds a new one to the orderstore
func (m *OrderManager) UpsertOrder(od *order.Detail) (resp *OrderUpsertResponse, err error) {
	if m == nil {
		return nil, fmt.Errorf("order manager %w", ErrNilSubsystem)
	}
	if atomic.LoadInt32(&m.started) == 0 {
		return nil, fmt.Errorf("order manager %w", ErrSubSystemNotStarted)
	}
	if od == nil {
		return nil, errNilOrder
	}
	var msg string
	defer func(message *string) {
		if message == nil {
			log.Errorf(log.OrderMgr, "UpsertOrder: produced nil order event message\n")
			return
		}
		m.orderStore.commsManager.PushEvent(base.Event{
			Type:    "order",
			Message: *message,
		})
	}(&msg)

	upsertResponse, err := m.orderStore.upsert(od)
	if err != nil {
		msg = fmt.Sprintf(
			"Order manager: Exchange %s unable to upsert order ID=%v internal ID=%v pair=%v price=%.8f amount=%.8f side=%v type=%v status=%v: %s",
			od.Exchange, od.ID, od.InternalOrderID, od.Pair, od.Price, od.Amount, od.Side, od.Type, od.Status, err)
		return nil, err
	}

	status := "updated"
	if upsertResponse.IsNewOrder {
		status = "added"
	}
	msg = fmt.Sprintf("Order manager: Exchange !!! %s %s order ID=%v internal ID=%v pair=%v price=%.8f amount=%.8f side=%v type=%v status=%v strategy=%s.",
		upsertResponse.OrderDetails.Exchange, status, upsertResponse.OrderDetails.ID, upsertResponse.OrderDetails.InternalOrderID,
		upsertResponse.OrderDetails.Pair, upsertResponse.OrderDetails.Price, upsertResponse.OrderDetails.Amount,
		upsertResponse.OrderDetails.Side, upsertResponse.OrderDetails.Type, upsertResponse.OrderDetails.Status, upsertResponse.OrderDetails.Strategy)
	if upsertResponse.IsNewOrder {
		log.Info(log.OrderMgr, msg)
		return upsertResponse, nil
	}
	log.Debug(log.OrderMgr, msg)
	return upsertResponse, nil
}

// get returns all orders for all exchanges
// should not be exported as it can have large impact if used improperly
func (s *store) get() map[string][]*order.Detail {
	s.m.Lock()
	orders := s.Orders
	s.m.Unlock()
	return orders
}

// getByExchangeAndID returns a specific order by exchange and id
func (s *store) getByExchangeAndID(exchange, id string) (*order.Detail, error) {
	s.m.Lock()
	defer s.m.Unlock()
	r, ok := s.Orders[strings.ToLower(exchange)]
	if !ok {
		return nil, ErrExchangeNotFound
	}

	for x := range r {
		if r[x].ID == id {
			return r[x], nil
		}
	}
	return nil, ErrOrderNotFound
}

// updateExisting checks if an order exists in the orderstore
// and then updates it
func (s *store) updateExisting(od *order.Detail) error {
	s.m.Lock()
	defer s.m.Unlock()
	r, ok := s.Orders[strings.ToLower(od.Exchange)]
	if !ok {
		return ErrExchangeNotFound
	}
	for x := range r {
		if r[x].ID == od.ID {
			r[x].UpdateOrderFromDetail(od)
			return nil
		}
	}

	return ErrOrderNotFound
}

// modifyExisting depends on mod.Exchange and given ID to uniquely identify an order and
// modify it.
func (s *store) modifyExisting(id string, mod *order.Modify) error {
	s.m.Lock()
	defer s.m.Unlock()
	r, ok := s.Orders[strings.ToLower(mod.Exchange)]
	if !ok {
		return ErrExchangeNotFound
	}
	for x := range r {
		if r[x].ID == id {
			r[x].UpdateOrderFromModify(mod)
			return nil
		}
	}
	return ErrOrderNotFound
}

// upsert (1) checks if such an exchange exists in the exchangeManager, (2) checks if
// order exists and updates/creates it.
func (s *store) upsert(od *order.Detail) (resp *OrderUpsertResponse, err error) {
	// fmt.Println("UPSERTING ORDERRRRRRRRRRRRRRRRRRRRR")
	if od == nil {
		return nil, errNilOrder
	}
	lName := strings.ToLower(od.Exchange)
	_, err = s.exchangeManager.GetExchangeByName(lName)
	if err != nil {
		return nil, err
	}
	s.m.Lock()
	defer s.m.Unlock()
	r, ok := s.Orders[lName]
	if !ok {
		s.Orders[lName] = []*order.Detail{od}
		resp = &OrderUpsertResponse{
			OrderDetails: od.Copy(),
			IsNewOrder:   true,
		}
		return resp, nil
	}
	for x := range r {
		if r[x].ID == od.ID {
			r[x].UpdateOrderFromDetail(od)
			if !s.dryRun {
				liveorder.Upsert(r[x])
			}
			resp = &OrderUpsertResponse{
				OrderDetails: r[x].Copy(),
				IsNewOrder:   false,
			}
			return resp, nil
		}
	}
	// Untracked websocket orders will not have internalIDs yet
	s.Orders[lName] = append(s.Orders[lName], od)
	resp = &OrderUpsertResponse{
		OrderDetails: od.Copy(),
		IsNewOrder:   true,
	}
	return resp, nil
}

// getByExchange returns orders by exchange
func (s *store) getByExchange(exchange string) ([]*order.Detail, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	r, ok := s.Orders[strings.ToLower(exchange)]
	if !ok {
		return nil, ErrExchangeNotFound
	}
	return r, nil
}

// getByInternalOrderID will search all orders for our internal orderID
// and return the order
func (s *store) getByInternalOrderID(internalOrderID int) (*order.Detail, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	for _, v := range s.Orders {
		for x := range v {
			if v[x].InternalOrderID == internalOrderID {
				return v[x], nil
			}
		}
	}
	return nil, ErrOrderNotFound
}

// exists verifies if the orderstore contains the provided order
func (s *store) exists(det *order.Detail) bool {
	if det == nil {
		return false
	}
	s.m.RLock()
	defer s.m.RUnlock()
	r, ok := s.Orders[strings.ToLower(det.Exchange)]
	if !ok {
		return false
	}

	for x := range r {
		if r[x].ID == det.ID {
			return true
		}
	}
	return false
}

func (m *OrderManager) SetOnFill(onFill func(order.Detail, eventtypes.DataEventHandler)) {
	m.onFill = onFill
}

// Add Adds an order to the orderStore for tracking the lifecycle
func (s *store) add(det *order.Detail) error {
	if det == nil {
		return errors.New("order store: Order is nil")
	}
	_, err := s.exchangeManager.GetExchangeByName(det.Exchange)
	if err != nil {
		return err
	}
	if s.exists(det) {
		return ErrOrdersAlreadyExists
	}
	s.m.Lock()
	defer s.m.Unlock()

	// fmt.Println(
	// 	"add order to store id:",
	// 	det.ID,
	// 	"internal:",
	// 	det.InternalOrderID)

	if !s.dryRun {
		// fmt.Println("status", det.Status, "filedat", det.FilledAt)
		_, err = liveorder.Upsert(det)
		if err != nil {
			errMsg := fmt.Sprintf("error upserting order", err)
			panic(errMsg)
			return err
		}
	}

	orders := s.Orders[strings.ToLower(det.Exchange)]
	orders = append(orders, det)
	s.Orders[strings.ToLower(det.Exchange)] = orders

	return nil
}

// SAVE THE TRADE TO THE DATABASE HERE?
// m.recordOrder(newOrder)
// // ADD LIVE ORDER TO PORTFOLIO STORE
// func (m *OrderManager) recordOrder(o *order.Submit) error {
// 	fmt.Println("insert side", o.Side)
//
// 	// store the order
// 	lo := liveorder.Details{
// 		Status:       order.New,
// 		OrderType:    order.Market,
// 		Exchange:     o.Exchange,
// 		StrategyID:   o.StrategyID,
// 		StrategyName: o.StrategyName,
// 		Pair:         o.Pair,
// 		Side:         o.Side,
// 	}
//
// 	var openOrder *liveorder.Details
// 	if len(p.store.openOrders[ev.GetStrategyID()]) == 0 {
// 		for i := range p.store.openOrders {
// 			fmt.Println(i)
// 		}
// 		panic(fmt.Sprintf("did not store open order for strategy %d", ev.GetStrategyID()))
// 	}
// 	for _, ord := range p.store.openOrders[ev.GetStrategyID()] {
// 		if ord.ID == ev.GetInternalOrderID() {
// 			openOrder = ord
// 			break
// 		}
// 	}
// 	if openOrder == nil {
// 		fmt.Println("error !!!!!! no interal openOrder id")
// 		return
// 	}
//
// 	if ev.GetIsOrderPlaced() {
// 		m.completeOrder(ev)
// 	}
// 	openOrder.Status = gctorder.Closed
// 	if !p.bot.Settings.EnableDryRun {
// 		id, err := liveorder.Update(openOrder)
// 		if err != nil || id == 0 {
// 			fmt.Println("error saving to db")
// 			os.Exit(2)
// 		}
// 	}
// 	if !m.dryRun {
// 		id, err := liveorder.Insert(lo)
// 		if err != nil {
// 			log.Errorln(log.Portfolio, "Unable to store order in database", err)
// 			// ev.SetDirection(signal.DoNothing)
// 			// ev.AppendReason(fmt.Sprintf("unable to store in database. err: %s", err))
// 			panic(fmt.Sprintf("unable to store order in database", err))
// 			return fmt.Errorf("unable to store in database. %v", err)
// 		}
// 		lo.ID = id
// 		o.InternalOrderID = lo.ID
// 	} else {
// 		// generate random lo ID to keep track
// 		// or lookup using another way (side/name/pair)
// 		o.ID = fmt.Sprintf("%s-%s-%s", lo.StrategyName, lo.Pair, lo.Side)
// 		o.InternalOrderID = 1234
// 	}
//
// 	if m.debug {
// 		fmt.Println("portfolio.OnSubmit", o.StrategyName, "orderID", o.ID, "internal", o.InternalOrderID)
// 	}
//
// 	// validate write
// 	// if m.debug {
// 	// 	fmt.Println("adding order for strategy:", o.StrategyName)
// 	// }
// 	// beforeLen := len(p.store.openOrders[ev.GetStrategyID()])
// 	// p.store.openOrders[ev.GetStrategyID()] = append(p.store.openOrders[ev.GetStrategyID()], &lo)
// 	// afterLen := len(p.store.openOrders[ev.GetStrategyID()])
// 	// // fmt.Println("store now has", afterLen, "orders for", ev.GetStrategyID())
// 	//
// 	// if afterLen > 1 {
// 	// 	panic(fmt.Sprintf("more than one open order for strategy: %d", ev.GetStrategyID()))
// 	// }
// 	//
// 	// // verify open order exists
// 	// if afterLen <= beforeLen {
// 	// 	fmt.Println("ERROR did not add open order")
// 	// 	return fmt.Errorf("did not return open order")
// 	// }
// 	return nil
// }

// // update the store with the submission ID
// ords, _ := om.GetOrdersSnapshot("")
// var internalOrderID int
// for i := range ords {
// 	fmt.Println("checking order id", ords[i].InternalOrderID, o.GetID())
// 	if ords[i].ID != omr.InternalOrderID {
// 		continue
// 	}
// 	ords[i].StrategyID = o.GetStrategyID()
// 	ords[i].Date = o.GetTime()
// 	ords[i].LastUpdated = o.GetTime()
// 	ords[i].CloseTime = o.GetTime()
// }

// id, err := uuid.NewV4()
// if err != nil {
// 	log.Warnf(log.OrderMgr,
// 		"Order manager: Unable to generate UUID. Err: %s",
// 		err)
// }

// getFilteredOrders returns a filtered copy of the orders
func (s *store) getFilteredOrders(f *order.Filter) ([]order.Detail, error) {
	if f == nil {
		return nil, errors.New("filter is nil")
	}
	s.m.RLock()
	defer s.m.RUnlock()

	var os []order.Detail
	// optimization if Exchange is filtered
	if f.Exchange != "" {
		if e, ok := s.Orders[strings.ToLower(f.Exchange)]; ok {
			for i := range e {
				if !e[i].MatchFilter(f) {
					continue
				}
				os = append(os, e[i].Copy())
			}
		}
	} else {
		for _, e := range s.Orders {
			for i := range e {
				if !e[i].MatchFilter(f) {
					continue
				}
				os = append(os, e[i].Copy())
			}
		}
	}
	return os, nil
}

// getActiveOrders returns copy of the orders that are active
func (s *store) getActiveOrders(f *order.Filter) []order.Detail {
	s.m.RLock()
	defer s.m.RUnlock()

	var orders []order.Detail
	switch {
	case f == nil:
		for _, e := range s.Orders {
			for i := range e {
				if !e[i].IsActive() {
					continue
				}
				orders = append(orders, e[i].Copy())
			}
		}
	case f.Exchange != "":
		// optimization if Exchange is filtered
		if e, ok := s.Orders[strings.ToLower(f.Exchange)]; ok {
			for i := range e {
				if !e[i].IsActive() || !e[i].MatchFilter(f) {
					continue
				}
				orders = append(orders, e[i].Copy())
			}
		}
	default:
		for _, e := range s.Orders {
			for i := range e {
				if !e[i].IsActive() || !e[i].MatchFilter(f) {
					continue
				}
				orders = append(orders, e[i].Copy())
			}
		}
	}

	return orders
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// func (p *Portfolio) completeOrder(ev submit.Event) {
// 	// if p.verbose {
// 	// 	log.Infoln(log.Portfolio, "completing order", ev.GetStrategyID())
// 	// }
// 	// fmt.Println("COMPLETING ORDER for:", ev.GetStrategyID())
// 	// fmt.Println("open orders", len(p.store.openOrders[ev.GetStrategyID()]))
// 	order := p.store.openOrders[ev.GetStrategyID()][0]
// 	p.store.closedOrders[ev.GetStrategyID()] = append(p.store.closedOrders[ev.GetStrategyID()], order)
// 	p.store.openOrders[ev.GetStrategyID()] = make([]*liveorder.Details, 0)
// 	// fmt.Println(ev.GetStrategyID(), " now has ", len(p.store.closedOrders[ev.GetStrategyID()]))
//
// 	// p.store.positions[ev.GetStrategyID()] = &positions.Position{Active: false}
// 	// // create or update position
// 	// for _, pos := range p.store.positions {
// 	// 	if f.GetDirection() == gctorder.Sell {
// 	// 		pos.Amount = pos.Amount.Sub(f.GetAmount())
// 	// 	} else if f.GetDirection() == gctorder.Buy {
// 	// 		pos.Amount = pos.Amount.Add(f.GetAmount())
// 	// 	}
// 	//
// 	// 	if !pos.Amount.IsZero() {
// 	// 		pos.Active = true
// 	// 	} else {
// 	// 		pos.Active = false
// 	// 	}
// 	// }
// }

func (m *OrderManager) GenerateDryRunID() int {
	orders, _ := m.GetOrdersSnapshot("")
	return len(orders) + 1
}
