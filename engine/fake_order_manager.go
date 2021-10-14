package engine

import (
	"context"
	"errors"
	"fmt"
	"gocryptotrader/communications/base"
	"gocryptotrader/exchange/order"
	"gocryptotrader/log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofrs/uuid"
)

// SetupOrderManager will boot up the OrderManager
func SetupFakeOrderManager(bot *Engine, exchangeManager iExchangeManager, communicationsManager iCommsManager, wg *sync.WaitGroup, verbose bool) (*FakeOrderManager, error) {
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
			bot:      bot,
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
// Kraken Order Status
// Enum: "pending" "open" "closed" "canceled" "expired"
func (m *FakeOrderManager) Submit(ctx context.Context, newOrder *order.Submit) (*OrderSubmitResponse, error) {
	if m.bot.Config.LiveMode {
		log.Debugln(log.FakeOrderMgr, "Order manager: Order Submitted", newOrder.Date, newOrder.StrategyID, newOrder.ID)
	}

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

	// all of the orders can be immediately submitted
	// meaning that we want on the reply of the broker
	// we could set this to happen in several goroutines instead and use a waitgroup to syncronize the results per strategy or currency
	// but that's not currently necessary since it's purely a live performance optimization that doesn't impact profits
	fakeSubmission := order.SubmitResponse{
		IsOrderPlaced: true,
		OrderID:       newOrder.ID,
		FullyMatched:  true,
	}

	resp, err := m.processSubmittedOrder(newOrder, fakeSubmission)
	if err != nil {
		log.Errorln(log.FakeOrderMgr, "error!!!!!!!!!!", err)
	}

	return resp, nil
}

// processSubmittedOrder adds a new order to the manager
func (m *FakeOrderManager) processSubmittedOrder(newOrder *order.Submit, result order.SubmitResponse) (*OrderSubmitResponse, error) {
	if !result.IsOrderPlaced {
		return nil, errors.New("order unable to be placed")
	}

	id, err := uuid.NewV4()
	if err != nil {
		log.Warnf(log.OrderMgr,
			"Order manager: Unable to generate UUID. Err: %s",
			err)
	}
	if newOrder.Date.IsZero() {
		newOrder.Date = time.Now()
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
	// log.Debugln(log.OrderMgr, msgInfo)

	m.orderStore.commsManager.PushEvent(base.Event{
		Type:    "order",
		Message: msgInfo,
	})
	status := order.New
	if result.FullyMatched {
		status = order.Filled
	}
	err = m.orderStore.add(&order.Detail{
		StrategyID:        newOrder.StrategyID,
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

	// // // custom on submit callback
	// if m.onSubmit != nil {
	// 	m.onSubmit(resp)
	// }

	return &OrderSubmitResponse{
		SubmitResponse: order.SubmitResponse{
			IsOrderPlaced: result.IsOrderPlaced,
			OrderID:       result.OrderID,
		},
		InternalOrderID: id.String(),
		StrategyID:      newOrder.StrategyID,
	}, nil
}
