package engine

// OrderManagerName is an exported subsystem name
const FakeOrderManagerName = "fakeorders"

// OrderManager processes and stores orders across enabled exchanges
type FakeOrderManager struct {
	OrderQueue       EventHolder
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
