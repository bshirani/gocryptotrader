package engine

import (
	"context"
	"errors"
	"sync"
	"time"

	"gocryptotrader/currency"
	"gocryptotrader/exchange"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/exchange/order"
)

// OrderManagerName is an exported subsystem name
const OrderManagerName = "orders"
const FakeOrderManagerName = "fakeorders"

// vars for the fund manager package
var (
	orderManagerDelay = time.Second * 10
	// ErrOrdersAlreadyExists occurs when the order already exists in the manager
	ErrOrdersAlreadyExists = errors.New("order already exists")
	// ErrOrderNotFound occurs when an order is not found in the orderstore
	ErrOrderNotFound            = errors.New("order does not exist")
	errNilCommunicationsManager = errors.New("cannot start with nil communications manager")
	// ErrOrderIDCannotBeEmpty occurs when an order does not have an ID
	ErrOrderIDCannotBeEmpty = errors.New("orderID cannot be empty")
	errNilOrder             = errors.New("nil order received")
)

type orderManagerConfig struct {
	EnforceLimitConfig     bool
	AllowMarketOrders      bool
	CancelOrdersOnShutdown bool
	LimitAmount            float64
	AllowedPairs           currency.Pairs
	AllowedExchanges       []string
	OrderSubmissionRetries int64
	LiveMode               bool
}

// store holds all orders by exchange
type store struct {
	m               sync.RWMutex
	Orders          map[string][]*order.Detail
	commsManager    iCommsManager
	exchangeManager iExchangeManager
	wg              *sync.WaitGroup
}

// OrderManager processes and stores orders across enabled exchanges
type OrderManager struct {
	started          int32
	processingOrders int32
	shutdown         chan struct{}
	orderStore       store
	cfg              orderManagerConfig
	verbose          bool
	onSubmit         func(*OrderSubmitResponse)
	onFill           func(*OrderSubmitResponse)
	onCancel         func(*OrderSubmitResponse)
}

type FakeOrderManager struct {
	OrderManager
}

// OrderManager processes and stores orders across enabled exchanges
type RealOrderManager struct {
	OrderManager
}

// OrderSubmitResponse contains the order response along with an internal order ID
type OrderSubmitResponse struct {
	order.SubmitResponse
	InternalOrderID string
}

// OrderUpsertResponse contains a copy of the resulting order details and a bool
// indicating if the order details were inserted (true) or updated (false)
type OrderUpsertResponse struct {
	OrderDetails order.Detail
	IsNewOrder   bool
}

// SetupOrderManager(exchangeManager iExchangeManager, communicationsManager iCommsManager, wg *sync.WaitGroup, verbose bool) (*FakeOrderManager, error)
type OrderManagerHandler interface {
	IsRunning() bool
	Start() error
	Stop() error
	Update()

	Add(o *order.Detail) error
	Cancel(ctx context.Context, cancel *order.Cancel) error
	CancelAllOrders(ctx context.Context, exchangeNames []exchange.IBotExchange)
	Exists(o *order.Detail) bool
	FetchAndUpdateExchangeOrder(exch exchange.IBotExchange, ord *order.Detail, assetType asset.Item) error
	GetByExchangeAndID(exchangeName, id string) (*order.Detail, error)
	GetOrderInfo(ctx context.Context, exchangeName, orderID string, cp currency.Pair, a asset.Item) (order.Detail, error)
	GetOrdersActive(f *order.Filter) ([]order.Detail, error)
	GetOrdersFiltered(f *order.Filter) ([]order.Detail, error)
	GetOrdersSnapshot(s order.Status) ([]order.Detail, time.Time)
	Modify(ctx context.Context, mod *order.Modify) (*order.ModifyResponse, error)
	SetOnCancel(onCancel func(*OrderSubmitResponse))
	SetOnFill(onFill func(*OrderSubmitResponse))
	SetOnSubmit(onSubmit func(*OrderSubmitResponse))
	Submit(ctx context.Context, newOrder *order.Submit) (*OrderSubmitResponse, error)
	UpdateExistingOrder(od *order.Detail) error
	UpsertOrder(od *order.Detail) (resp *OrderUpsertResponse, err error)
}
