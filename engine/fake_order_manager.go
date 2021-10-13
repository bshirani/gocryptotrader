package engine

import (
	"context"
	"fmt"
	"gocryptotrader/exchange/order"
	"gocryptotrader/log"
	"sync"
	"sync/atomic"
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
		OrderManager{
			shutdown: make(chan struct{}),
			orderStore: store{
				Orders:          make(map[string][]*order.Detail),
				exchangeManager: exchangeManager,
				commsManager:    communicationsManager,
				wg:              wg,
			},
			verbose: verbose,
		},
	}, nil
}

// Submit will take in an order struct, send it to the exchange and
// populate it in the FakeOrderManager if successful
func (m *FakeOrderManager) Submit(ctx context.Context, newOrder *order.Submit) (*OrderSubmitResponse, error) {
	log.Debugln(log.FakeOrderMgr, "Order manager: Order Submitted", newOrder.Date, newOrder.StrategyID, newOrder.ID)

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

// processOrders iterates over all exchange orders via API
// and adds them to the internal order store
func (m *FakeOrderManager) processOrders() {}
