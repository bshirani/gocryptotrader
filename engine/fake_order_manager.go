package engine

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gocryptotrader/common"
	"gocryptotrader/communications/base"
	"gocryptotrader/currency"
	"gocryptotrader/exchange"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/exchange/order"
	"gocryptotrader/log"

	"github.com/gofrs/uuid"
)

// SetupOrderManager will boot up the OrderManager
func SetupFakeOrderManager(exchangeManager iExchangeManager, communicationsManager iCommsManager, wg *sync.WaitGroup, verbose bool) (*FakeOrderManager, error) {
	// log.Debugln(log.FakeOrderMgr, "...")

	if exchangeManager == nil {
		return nil, errNilExchangeManager
	}
	if communicationsManager == nil {
		return nil, errNilCommunicationsManager
	}
	if wg == nil {
		return nil, errNilWaitGroup
	}

	return &FakeOrderManager{
		shutdown: make(chan struct{}),
		orderStore: store{
			Orders:          make(map[string][]*order.Detail),
			exchangeManager: exchangeManager,
			commsManager:    communicationsManager,
			wg:              wg,
		},
		verbose: verbose,
	}, nil
}

func (m *FakeOrderManager) SetOnSubmit(onSubmit func(*OrderSubmitResponse)) {
	m.onSubmit = onSubmit
}

func (m *FakeOrderManager) SetOnFill(onFill func(*OrderSubmitResponse)) {
	m.onFill = onFill
}

func (m *FakeOrderManager) SetOnCancel(onCancel func(*OrderSubmitResponse)) {
	m.onCancel = onCancel
}

// IsRunning safely checks whether the subsystem is running
func (m *FakeOrderManager) IsRunning() bool {
	if m == nil {
		return false
	}
	return atomic.LoadInt32(&m.started) == 1
}

func (m *FakeOrderManager) Start() error {
	if m == nil {
		return fmt.Errorf("fake order manager %w", ErrNilSubsystem)
	}
	if !atomic.CompareAndSwapInt32(&m.started, 0, 1) {
		return fmt.Errorf("fake order manager %w", ErrSubSystemAlreadyStarted)
	}
	log.Debugln(log.FakeOrderMgr, "fake order manager starting...")
	m.shutdown = make(chan struct{})
	go m.run()
	return nil
}

// Stop attempts to shutdown the subsystem
func (m *FakeOrderManager) Stop() error {
	if m == nil {
		return fmt.Errorf("fake order manager %w", ErrNilSubsystem)
	}
	if atomic.LoadInt32(&m.started) == 0 {
		return fmt.Errorf("fake order manager %w", ErrSubSystemNotStarted)
	}

	defer func() {
		atomic.CompareAndSwapInt32(&m.started, 1, 0)
	}()

	log.Debugln(log.FakeOrderMgr, "fake order manager shutting down...")
	close(m.shutdown)
	return nil
}

// gracefulShutdown cancels all orders (if enabled) before shutting down
func (m *FakeOrderManager) gracefulShutdown() {
	if m.cfg.CancelOrdersOnShutdown {
		log.Debugln(log.FakeOrderMgr, "fake order manager: Cancelling any open orders...")
		exchanges, err := m.orderStore.exchangeManager.GetExchanges()
		if err != nil {
			log.Errorf(log.FakeOrderMgr, "fake order manager cannot get exchanges: %v", err)
			return
		}
		m.CancelAllOrders(context.TODO(), exchanges)
	}
}

// run will periodically process orders
func (m *FakeOrderManager) run() {
	// log.Debugln(log.FakeOrderMgr, "fake order manager started.")
	m.processOrders()
	tick := time.NewTicker(orderManagerDelay)
	m.orderStore.wg.Add(1)
	defer func() {
		log.Debugln(log.FakeOrderMgr, "fake order manager shutdown.")
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
func (m *FakeOrderManager) CancelAllOrders(ctx context.Context, exchangeNames []exchange.IBotExchange) {
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
			log.Debugf(log.FakeOrderMgr, "fake order manager: Cancelling order(s) for exchange %s.", exchangeNames[i].GetName())
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
				log.Error(log.FakeOrderMgr, err)
			}
		}
	}
}

// Cancel will find the order in the FakeOrderManager, send a cancel request
// to the exchange and if successful, update the status of the order
func (m *FakeOrderManager) Cancel(ctx context.Context, cancel *order.Cancel) error {
	if m == nil {
		return fmt.Errorf("fake order manager %w", ErrNilSubsystem)
	}
	if atomic.LoadInt32(&m.started) == 0 {
		return fmt.Errorf("fake order manager %w", ErrSubSystemNotStarted)
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

	exch, err := m.orderStore.exchangeManager.GetExchangeByName(cancel.Exchange)
	if err != nil {
		return err
	}

	if cancel.AssetType.String() != "" && !exch.GetAssetTypes(false).Contains(cancel.AssetType) {
		err = errors.New("order asset type not supported by exchange")
		return err
	}

	log.Debugf(log.FakeOrderMgr, "fake order manager: Cancelling order ID %v [%+v]",
		cancel.ID, cancel)

	err = exch.CancelOrder(ctx, cancel)
	if err != nil {
		err = fmt.Errorf("%v - Failed to cancel order: %w", cancel.Exchange, err)
		return err
	}
	var od *order.Detail
	od, err = m.orderStore.getByExchangeAndID(cancel.Exchange, cancel.ID)
	if err != nil {
		err = fmt.Errorf("%v - Failed to retrieve order %v to update cancelled status: %w", cancel.Exchange, cancel.ID, err)
		return err
	}

	od.Status = order.Cancelled
	msg := fmt.Sprintf("fake order manager: Exchange %s order ID=%v cancelled.",
		od.Exchange, od.ID)
	log.Debugln(log.FakeOrderMgr, msg)
	m.orderStore.commsManager.PushEvent(base.Event{
		Type:    "order",
		Message: msg,
	})

	return nil
}

// GetOrderInfo calls the exchange's wrapper GetOrderInfo function
// and stores the result in the fake order manager
func (m *FakeOrderManager) GetOrderInfo(ctx context.Context, exchangeName, orderID string, cp currency.Pair, a asset.Item) (order.Detail, error) {
	if m == nil {
		return order.Detail{}, fmt.Errorf("fake order manager %w", ErrNilSubsystem)
	}
	if atomic.LoadInt32(&m.started) == 0 {
		return order.Detail{}, fmt.Errorf("fake order manager %w", ErrSubSystemNotStarted)
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
func (m *FakeOrderManager) validate(newOrder *order.Submit) error {
	if newOrder == nil {
		return errors.New("order cannot be nil")
	}

	if newOrder.Exchange == "" {
		return errors.New("order exchange name must be specified")
	}

	if err := newOrder.Validate(); err != nil {
		return fmt.Errorf("fake order manager: %w", err)
	}

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
func (m *FakeOrderManager) Modify(ctx context.Context, mod *order.Modify) (*order.ModifyResponse, error) {
	if m == nil {
		return nil, fmt.Errorf("fake order manager %w", ErrNilSubsystem)
	}
	if atomic.LoadInt32(&m.started) == 0 {
		return nil, fmt.Errorf("fake order manager %w", ErrSubSystemNotStarted)
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
			"fake order manager: Exchange %s order ID=%v: failed to modify",
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
		message = "fake order manager: Exchange %s order ID=%v: modified on exchange, but failed to modify locally"
	} else {
		message = "fake order manager: Exchange %s order ID=%v: modified successfully"
	}
	m.orderStore.commsManager.PushEvent(base.Event{
		Type:    "order",
		Message: fmt.Sprintf(message, mod.Exchange, res.ID),
	})
	return &order.ModifyResponse{OrderID: res.ID}, err
}

// Submit will take in an order struct, send it to the exchange and
// populate it in the FakeOrderManager if successful
func (m *FakeOrderManager) Submit(ctx context.Context, newOrder *order.Submit) (*OrderSubmitResponse, error) {
	log.Debugln(log.FakeOrderMgr, "Order manager: Order Submitted", newOrder.ID)

	if m == nil {
		return nil, fmt.Errorf("fake order manager %w", ErrNilSubsystem)
	}
	if atomic.LoadInt32(&m.started) == 0 {
		return nil, fmt.Errorf("fake order manager %w", ErrSubSystemNotStarted)
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
		return nil, fmt.Errorf("fake order manager: exchange %s unable to place order: %w",
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

	// we want to create a fake submission here
	// we only call SubmitOrder for real orders, we act like we are the exchange
	// result, err := exch.SubmitOrder(ctx, newOrder)
	// if err != nil {
	// 	log.Errorln(log.FakeOrderMgr, "error submitting order", err)
	// 	return nil, err
	// }

	fakeSubmission := order.SubmitResponse{
		IsOrderPlaced: true,
		OrderID:       newOrder.ID,
		FullyMatched:  true,
	}

	resp, err := m.processSubmittedOrder(newOrder, fakeSubmission)
	if err != nil {
		log.Errorln(log.FakeOrderMgr, "error", err)
	}

	if m.onSubmit != nil {
		m.onSubmit(resp)
	}
	return resp, nil
}

// GetOrdersSnapshot returns a snapshot of all orders in the orderstore. It optionally filters any orders that do not match the status
// but a status of "" or ANY will include all
// the time adds contexts for the when the snapshot is relevant for
func (m *FakeOrderManager) GetOrdersSnapshot(s order.Status) ([]order.Detail, time.Time) {
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
func (m *FakeOrderManager) GetOrdersFiltered(f *order.Filter) ([]order.Detail, error) {
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
func (m *FakeOrderManager) GetOrdersActive(f *order.Filter) ([]order.Detail, error) {
	if m == nil {
		return nil, fmt.Errorf("order manager %w", ErrNilSubsystem)
	}
	if atomic.LoadInt32(&m.started) == 0 {
		return nil, fmt.Errorf("order manager %w", ErrSubSystemNotStarted)
	}
	return m.orderStore.getActiveOrders(f), nil
}

// processSubmittedOrder adds a new order to the manager
func (m *FakeOrderManager) processSubmittedOrder(newOrder *order.Submit, result order.SubmitResponse) (*OrderSubmitResponse, error) {
	if !result.IsOrderPlaced {
		return nil, errors.New("order unable to be placed")
	}

	id, err := uuid.NewV4()
	if err != nil {
		log.Warnf(log.FakeOrderMgr,
			"Order manager: Unable to generate UUID. Err: %s",
			err)
	}
	if newOrder.Date.IsZero() {
		newOrder.Date = time.Now()
	}

	if newOrder.StrategyID == "" {
		return nil, errors.New("order must have a strategy")
	}

	// formatting for backtest only
	odate := newOrder.Date
	dtime := fmt.Sprintf("%d-%02d-%02d %d:%02d", odate.Year(), odate.Month(), odate.Day(), odate.Hour(), odate.Minute())
	msgInfo := fmt.Sprintf("%v %s %-10s %v %-5v %10v",
		newOrder.Pair,
		dtime,
		newOrder.StrategyID,
		newOrder.Pair,
		newOrder.Side,
		newOrder.Price)
	// fmt.Println(msgInfo)

	// msg := fmt.Sprintf("Order manager: Strategy=%s Exchange=%s submitted order ID=%v [Ours: %v] pair=%v price=%v amount=%v side=%v type=%v for time %v.",
	// 	newOrder.StrategyID,
	// 	newOrder.Exchange,
	// 	result.OrderID,
	// 	id.String(),
	// 	newOrder.Pair,
	// 	newOrder.Price,
	// 	newOrder.Amount,
	// 	newOrder.Side,
	// 	newOrder.Type,
	// 	newOrder.Date)
	// log.Debugln(log.FakeOrderMgr, msgInfo)

	m.orderStore.commsManager.PushEvent(base.Event{
		Type:    "order",
		Message: msgInfo,
	})
	status := order.New
	if result.FullyMatched {
		status = order.Filled
	}
	err = m.orderStore.add(&order.Detail{
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
		InternalOrderID:   id.String(),
		ID:                result.OrderID,
		AccountID:         newOrder.AccountID,
		ClientID:          newOrder.ClientID,
		ClientOrderID:     newOrder.ClientOrderID,
		WalletAddress:     newOrder.WalletAddress,
		Type:              newOrder.Type,
		Side:              newOrder.Side,
		Status:            status,
		AssetType:         newOrder.AssetType,
		Date:              time.Now(),
		LastUpdated:       time.Now(),
		Pair:              newOrder.Pair,
		Leverage:          newOrder.Leverage,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to add %v order %v to orderStore: %s", newOrder.Exchange, result.OrderID, err)
	}

	return &OrderSubmitResponse{
		SubmitResponse: order.SubmitResponse{
			IsOrderPlaced: result.IsOrderPlaced,
			OrderID:       result.OrderID,
		},
		InternalOrderID: id.String(),
	}, nil
}

// processOrders iterates over all exchange orders via API
// and adds them to the internal order store
func (m *FakeOrderManager) processOrders() {
	if !atomic.CompareAndSwapInt32(&m.processingOrders, 0, 1) {
		return
	}
	defer func() {
		atomic.StoreInt32(&m.processingOrders, 0)
	}()
	exchanges, err := m.orderStore.exchangeManager.GetExchanges()
	if err != nil {
		log.Errorf(log.FakeOrderMgr, "Order manager cannot get exchanges: %v", err)
		return
	}
	var wg sync.WaitGroup
	for i := range exchanges {
		if !exchanges[i].GetAuthenticatedAPISupport(exchange.RestAuthentication) {
			continue
		}
		log.Debugf(log.FakeOrderMgr,
			"Order manager: Processing orders for exchange %v.",
			exchanges[i].GetName())

		supportedAssets := exchanges[i].GetAssetTypes(true)
		for y := range supportedAssets {
			pairs, err := exchanges[i].GetEnabledPairs(supportedAssets[y])
			if err != nil {
				log.Errorf(log.FakeOrderMgr,
					"Order manager: Unable to get enabled pairs for %s and asset type %s: %s",
					exchanges[i].GetName(),
					supportedAssets[y],
					err)
				continue
			}

			if len(pairs) == 0 {
				if m.verbose {
					log.Debugf(log.FakeOrderMgr,
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
			requiresProcessing := make(map[string]bool, len(orders))
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
				log.Errorf(log.FakeOrderMgr,
					"Order manager: Unable to get active orders for %s and asset type %s: %s",
					exchanges[i].GetName(),
					supportedAssets[y],
					err)
				continue
			}
			if len(orders) == 0 && len(result) == 0 {
				continue
			}

			for z := range result {
				upsertResponse, err := m.UpsertOrder(&result[z])
				if err != nil {
					log.Error(log.FakeOrderMgr, err)
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

func (m *FakeOrderManager) processMatchingOrders(exch exchange.IBotExchange, orders []order.Detail, requiresProcessing map[string]bool, wg *sync.WaitGroup) {
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
				log.Error(log.FakeOrderMgr, err)
			}
		}
	}
}

// FetchAndUpdateExchangeOrder calls the exchange to upsert an order to the order store
func (m *FakeOrderManager) FetchAndUpdateExchangeOrder(exch exchange.IBotExchange, ord *order.Detail, assetType asset.Item) error {
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
func (m *FakeOrderManager) Exists(o *order.Detail) bool {
	if m == nil || atomic.LoadInt32(&m.started) == 0 {
		return false
	}

	return m.orderStore.exists(o)
}

// Add adds an order to the orderstore
func (m *FakeOrderManager) Add(o *order.Detail) error {
	if m == nil {
		return fmt.Errorf("order manager %w", ErrNilSubsystem)
	}
	if atomic.LoadInt32(&m.started) == 0 {
		return fmt.Errorf("order manager %w", ErrSubSystemNotStarted)
	}

	return m.orderStore.add(o)
}

// GetByExchangeAndID returns a copy of an order from an exchange if it matches the ID
func (m *FakeOrderManager) GetByExchangeAndID(exchangeName, id string) (*order.Detail, error) {
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
func (m *FakeOrderManager) UpdateExistingOrder(od *order.Detail) error {
	if m == nil {
		return fmt.Errorf("order manager %w", ErrNilSubsystem)
	}
	if atomic.LoadInt32(&m.started) == 0 {
		return fmt.Errorf("order manager %w", ErrSubSystemNotStarted)
	}
	return m.orderStore.updateExisting(od)
}

// UpsertOrder updates an existing order or adds a new one to the orderstore
func (m *FakeOrderManager) UpsertOrder(od *order.Detail) (resp *OrderUpsertResponse, err error) {
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
			log.Errorf(log.FakeOrderMgr, "UpsertOrder: produced nil order event message\n")
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
		log.Info(log.FakeOrderMgr, msg)
		return upsertResponse, nil
	}
	log.Debug(log.FakeOrderMgr, msg)
	return upsertResponse, nil
}
