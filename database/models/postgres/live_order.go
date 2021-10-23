// Code generated by SQLBoiler 4.7.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// LiveOrder is an object representing the database table.
type LiveOrder struct {
	ID            string       `boil:"id" json:"id" toml:"id" yaml:"id"`
	Status        string       `boil:"status" json:"status" toml:"status" yaml:"status"`
	OrderType     string       `boil:"order_type" json:"order_type" toml:"order_type" yaml:"order_type"`
	Exchange      string       `boil:"exchange" json:"exchange" toml:"exchange" yaml:"exchange"`
	StrategyID    string       `boil:"strategy_id" json:"strategy_id" toml:"strategy_id" yaml:"strategy_id"`
	InternalID    string       `boil:"internal_id" json:"internal_id" toml:"internal_id" yaml:"internal_id"`
	Side          null.String  `boil:"side" json:"side,omitempty" toml:"side" yaml:"side,omitempty"`
	ClientOrderID null.String  `boil:"client_order_id" json:"client_order_id,omitempty" toml:"client_order_id" yaml:"client_order_id,omitempty"`
	Amount        null.Float32 `boil:"amount" json:"amount,omitempty" toml:"amount" yaml:"amount,omitempty"`
	Symbol        null.String  `boil:"symbol" json:"symbol,omitempty" toml:"symbol" yaml:"symbol,omitempty"`
	Price         null.Float32 `boil:"price" json:"price,omitempty" toml:"price" yaml:"price,omitempty"`
	Fee           null.Float32 `boil:"fee" json:"fee,omitempty" toml:"fee" yaml:"fee,omitempty"`
	Cost          null.Float32 `boil:"cost" json:"cost,omitempty" toml:"cost" yaml:"cost,omitempty"`
	FilledAt      null.Time    `boil:"filled_at" json:"filled_at,omitempty" toml:"filled_at" yaml:"filled_at,omitempty"`
	AssetType     null.Int     `boil:"asset_type" json:"asset_type,omitempty" toml:"asset_type" yaml:"asset_type,omitempty"`
	SubmittedAt   null.Time    `boil:"submitted_at" json:"submitted_at,omitempty" toml:"submitted_at" yaml:"submitted_at,omitempty"`
	CancelledAt   null.Time    `boil:"cancelled_at" json:"cancelled_at,omitempty" toml:"cancelled_at" yaml:"cancelled_at,omitempty"`
	CreatedAt     time.Time    `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt     time.Time    `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *liveOrderR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L liveOrderL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var LiveOrderColumns = struct {
	ID            string
	Status        string
	OrderType     string
	Exchange      string
	StrategyID    string
	InternalID    string
	Side          string
	ClientOrderID string
	Amount        string
	Symbol        string
	Price         string
	Fee           string
	Cost          string
	FilledAt      string
	AssetType     string
	SubmittedAt   string
	CancelledAt   string
	CreatedAt     string
	UpdatedAt     string
}{
	ID:            "id",
	Status:        "status",
	OrderType:     "order_type",
	Exchange:      "exchange",
	StrategyID:    "strategy_id",
	InternalID:    "internal_id",
	Side:          "side",
	ClientOrderID: "client_order_id",
	Amount:        "amount",
	Symbol:        "symbol",
	Price:         "price",
	Fee:           "fee",
	Cost:          "cost",
	FilledAt:      "filled_at",
	AssetType:     "asset_type",
	SubmittedAt:   "submitted_at",
	CancelledAt:   "cancelled_at",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

var LiveOrderTableColumns = struct {
	ID            string
	Status        string
	OrderType     string
	Exchange      string
	StrategyID    string
	InternalID    string
	Side          string
	ClientOrderID string
	Amount        string
	Symbol        string
	Price         string
	Fee           string
	Cost          string
	FilledAt      string
	AssetType     string
	SubmittedAt   string
	CancelledAt   string
	CreatedAt     string
	UpdatedAt     string
}{
	ID:            "live_order.id",
	Status:        "live_order.status",
	OrderType:     "live_order.order_type",
	Exchange:      "live_order.exchange",
	StrategyID:    "live_order.strategy_id",
	InternalID:    "live_order.internal_id",
	Side:          "live_order.side",
	ClientOrderID: "live_order.client_order_id",
	Amount:        "live_order.amount",
	Symbol:        "live_order.symbol",
	Price:         "live_order.price",
	Fee:           "live_order.fee",
	Cost:          "live_order.cost",
	FilledAt:      "live_order.filled_at",
	AssetType:     "live_order.asset_type",
	SubmittedAt:   "live_order.submitted_at",
	CancelledAt:   "live_order.cancelled_at",
	CreatedAt:     "live_order.created_at",
	UpdatedAt:     "live_order.updated_at",
}

// Generated where

type whereHelpernull_Float32 struct{ field string }

func (w whereHelpernull_Float32) EQ(x null.Float32) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Float32) NEQ(x null.Float32) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Float32) LT(x null.Float32) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Float32) LTE(x null.Float32) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Float32) GT(x null.Float32) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Float32) GTE(x null.Float32) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

func (w whereHelpernull_Float32) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Float32) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }

var LiveOrderWhere = struct {
	ID            whereHelperstring
	Status        whereHelperstring
	OrderType     whereHelperstring
	Exchange      whereHelperstring
	StrategyID    whereHelperstring
	InternalID    whereHelperstring
	Side          whereHelpernull_String
	ClientOrderID whereHelpernull_String
	Amount        whereHelpernull_Float32
	Symbol        whereHelpernull_String
	Price         whereHelpernull_Float32
	Fee           whereHelpernull_Float32
	Cost          whereHelpernull_Float32
	FilledAt      whereHelpernull_Time
	AssetType     whereHelpernull_Int
	SubmittedAt   whereHelpernull_Time
	CancelledAt   whereHelpernull_Time
	CreatedAt     whereHelpertime_Time
	UpdatedAt     whereHelpertime_Time
}{
	ID:            whereHelperstring{field: "\"live_order\".\"id\""},
	Status:        whereHelperstring{field: "\"live_order\".\"status\""},
	OrderType:     whereHelperstring{field: "\"live_order\".\"order_type\""},
	Exchange:      whereHelperstring{field: "\"live_order\".\"exchange\""},
	StrategyID:    whereHelperstring{field: "\"live_order\".\"strategy_id\""},
	InternalID:    whereHelperstring{field: "\"live_order\".\"internal_id\""},
	Side:          whereHelpernull_String{field: "\"live_order\".\"side\""},
	ClientOrderID: whereHelpernull_String{field: "\"live_order\".\"client_order_id\""},
	Amount:        whereHelpernull_Float32{field: "\"live_order\".\"amount\""},
	Symbol:        whereHelpernull_String{field: "\"live_order\".\"symbol\""},
	Price:         whereHelpernull_Float32{field: "\"live_order\".\"price\""},
	Fee:           whereHelpernull_Float32{field: "\"live_order\".\"fee\""},
	Cost:          whereHelpernull_Float32{field: "\"live_order\".\"cost\""},
	FilledAt:      whereHelpernull_Time{field: "\"live_order\".\"filled_at\""},
	AssetType:     whereHelpernull_Int{field: "\"live_order\".\"asset_type\""},
	SubmittedAt:   whereHelpernull_Time{field: "\"live_order\".\"submitted_at\""},
	CancelledAt:   whereHelpernull_Time{field: "\"live_order\".\"cancelled_at\""},
	CreatedAt:     whereHelpertime_Time{field: "\"live_order\".\"created_at\""},
	UpdatedAt:     whereHelpertime_Time{field: "\"live_order\".\"updated_at\""},
}

// LiveOrderRels is where relationship names are stored.
var LiveOrderRels = struct {
	EntryOrderLiveTrades string
}{
	EntryOrderLiveTrades: "EntryOrderLiveTrades",
}

// liveOrderR is where relationships are stored.
type liveOrderR struct {
	EntryOrderLiveTrades LiveTradeSlice `boil:"EntryOrderLiveTrades" json:"EntryOrderLiveTrades" toml:"EntryOrderLiveTrades" yaml:"EntryOrderLiveTrades"`
}

// NewStruct creates a new relationship struct
func (*liveOrderR) NewStruct() *liveOrderR {
	return &liveOrderR{}
}

// liveOrderL is where Load methods for each relationship are stored.
type liveOrderL struct{}

var (
	liveOrderAllColumns            = []string{"id", "status", "order_type", "exchange", "strategy_id", "internal_id", "side", "client_order_id", "amount", "symbol", "price", "fee", "cost", "filled_at", "asset_type", "submitted_at", "cancelled_at", "created_at", "updated_at"}
	liveOrderColumnsWithoutDefault = []string{"status", "order_type", "exchange", "strategy_id", "internal_id", "side", "client_order_id", "amount", "symbol", "price", "fee", "cost", "filled_at", "asset_type", "submitted_at", "cancelled_at"}
	liveOrderColumnsWithDefault    = []string{"id", "created_at", "updated_at"}
	liveOrderPrimaryKeyColumns     = []string{"id"}
)

type (
	// LiveOrderSlice is an alias for a slice of pointers to LiveOrder.
	// This should almost always be used instead of []LiveOrder.
	LiveOrderSlice []*LiveOrder
	// LiveOrderHook is the signature for custom LiveOrder hook methods
	LiveOrderHook func(context.Context, boil.ContextExecutor, *LiveOrder) error

	liveOrderQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	liveOrderType                 = reflect.TypeOf(&LiveOrder{})
	liveOrderMapping              = queries.MakeStructMapping(liveOrderType)
	liveOrderPrimaryKeyMapping, _ = queries.BindMapping(liveOrderType, liveOrderMapping, liveOrderPrimaryKeyColumns)
	liveOrderInsertCacheMut       sync.RWMutex
	liveOrderInsertCache          = make(map[string]insertCache)
	liveOrderUpdateCacheMut       sync.RWMutex
	liveOrderUpdateCache          = make(map[string]updateCache)
	liveOrderUpsertCacheMut       sync.RWMutex
	liveOrderUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var liveOrderBeforeInsertHooks []LiveOrderHook
var liveOrderBeforeUpdateHooks []LiveOrderHook
var liveOrderBeforeDeleteHooks []LiveOrderHook
var liveOrderBeforeUpsertHooks []LiveOrderHook

var liveOrderAfterInsertHooks []LiveOrderHook
var liveOrderAfterSelectHooks []LiveOrderHook
var liveOrderAfterUpdateHooks []LiveOrderHook
var liveOrderAfterDeleteHooks []LiveOrderHook
var liveOrderAfterUpsertHooks []LiveOrderHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *LiveOrder) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range liveOrderBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *LiveOrder) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range liveOrderBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *LiveOrder) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range liveOrderBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *LiveOrder) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range liveOrderBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *LiveOrder) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range liveOrderAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *LiveOrder) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range liveOrderAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *LiveOrder) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range liveOrderAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *LiveOrder) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range liveOrderAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *LiveOrder) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range liveOrderAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddLiveOrderHook registers your hook function for all future operations.
func AddLiveOrderHook(hookPoint boil.HookPoint, liveOrderHook LiveOrderHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		liveOrderBeforeInsertHooks = append(liveOrderBeforeInsertHooks, liveOrderHook)
	case boil.BeforeUpdateHook:
		liveOrderBeforeUpdateHooks = append(liveOrderBeforeUpdateHooks, liveOrderHook)
	case boil.BeforeDeleteHook:
		liveOrderBeforeDeleteHooks = append(liveOrderBeforeDeleteHooks, liveOrderHook)
	case boil.BeforeUpsertHook:
		liveOrderBeforeUpsertHooks = append(liveOrderBeforeUpsertHooks, liveOrderHook)
	case boil.AfterInsertHook:
		liveOrderAfterInsertHooks = append(liveOrderAfterInsertHooks, liveOrderHook)
	case boil.AfterSelectHook:
		liveOrderAfterSelectHooks = append(liveOrderAfterSelectHooks, liveOrderHook)
	case boil.AfterUpdateHook:
		liveOrderAfterUpdateHooks = append(liveOrderAfterUpdateHooks, liveOrderHook)
	case boil.AfterDeleteHook:
		liveOrderAfterDeleteHooks = append(liveOrderAfterDeleteHooks, liveOrderHook)
	case boil.AfterUpsertHook:
		liveOrderAfterUpsertHooks = append(liveOrderAfterUpsertHooks, liveOrderHook)
	}
}

// One returns a single liveOrder record from the query.
func (q liveOrderQuery) One(ctx context.Context, exec boil.ContextExecutor) (*LiveOrder, error) {
	o := &LiveOrder{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "postgres: failed to execute a one query for live_order")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all LiveOrder records from the query.
func (q liveOrderQuery) All(ctx context.Context, exec boil.ContextExecutor) (LiveOrderSlice, error) {
	var o []*LiveOrder

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "postgres: failed to assign all query results to LiveOrder slice")
	}

	if len(liveOrderAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all LiveOrder records in the query.
func (q liveOrderQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: failed to count live_order rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q liveOrderQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "postgres: failed to check if live_order exists")
	}

	return count > 0, nil
}

// EntryOrderLiveTrades retrieves all the live_trade's LiveTrades with an executor via entry_order_id column.
func (o *LiveOrder) EntryOrderLiveTrades(mods ...qm.QueryMod) liveTradeQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"live_trade\".\"entry_order_id\"=?", o.ID),
	)

	query := LiveTrades(queryMods...)
	queries.SetFrom(query.Query, "\"live_trade\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"live_trade\".*"})
	}

	return query
}

// LoadEntryOrderLiveTrades allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (liveOrderL) LoadEntryOrderLiveTrades(ctx context.Context, e boil.ContextExecutor, singular bool, maybeLiveOrder interface{}, mods queries.Applicator) error {
	var slice []*LiveOrder
	var object *LiveOrder

	if singular {
		object = maybeLiveOrder.(*LiveOrder)
	} else {
		slice = *maybeLiveOrder.(*[]*LiveOrder)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &liveOrderR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &liveOrderR{}
			}

			for _, a := range args {
				if a == obj.ID {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`live_trade`),
		qm.WhereIn(`live_trade.entry_order_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load live_trade")
	}

	var resultSlice []*LiveTrade
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice live_trade")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on live_trade")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for live_trade")
	}

	if len(liveTradeAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.EntryOrderLiveTrades = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &liveTradeR{}
			}
			foreign.R.EntryOrder = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.EntryOrderID {
				local.R.EntryOrderLiveTrades = append(local.R.EntryOrderLiveTrades, foreign)
				if foreign.R == nil {
					foreign.R = &liveTradeR{}
				}
				foreign.R.EntryOrder = local
				break
			}
		}
	}

	return nil
}

// AddEntryOrderLiveTrades adds the given related objects to the existing relationships
// of the live_order, optionally inserting them as new records.
// Appends related to o.R.EntryOrderLiveTrades.
// Sets related.R.EntryOrder appropriately.
func (o *LiveOrder) AddEntryOrderLiveTrades(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*LiveTrade) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.EntryOrderID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"live_trade\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"entry_order_id"}),
				strmangle.WhereClause("\"", "\"", 2, liveTradePrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.EntryOrderID = o.ID
		}
	}

	if o.R == nil {
		o.R = &liveOrderR{
			EntryOrderLiveTrades: related,
		}
	} else {
		o.R.EntryOrderLiveTrades = append(o.R.EntryOrderLiveTrades, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &liveTradeR{
				EntryOrder: o,
			}
		} else {
			rel.R.EntryOrder = o
		}
	}
	return nil
}

// LiveOrders retrieves all the records using an executor.
func LiveOrders(mods ...qm.QueryMod) liveOrderQuery {
	mods = append(mods, qm.From("\"live_order\""))
	return liveOrderQuery{NewQuery(mods...)}
}

// FindLiveOrder retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindLiveOrder(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*LiveOrder, error) {
	liveOrderObj := &LiveOrder{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"live_order\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, liveOrderObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "postgres: unable to select from live_order")
	}

	if err = liveOrderObj.doAfterSelectHooks(ctx, exec); err != nil {
		return liveOrderObj, err
	}

	return liveOrderObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *LiveOrder) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("postgres: no live_order provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		if o.UpdatedAt.IsZero() {
			o.UpdatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(liveOrderColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	liveOrderInsertCacheMut.RLock()
	cache, cached := liveOrderInsertCache[key]
	liveOrderInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			liveOrderAllColumns,
			liveOrderColumnsWithDefault,
			liveOrderColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(liveOrderType, liveOrderMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(liveOrderType, liveOrderMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"live_order\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"live_order\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "postgres: unable to insert into live_order")
	}

	if !cached {
		liveOrderInsertCacheMut.Lock()
		liveOrderInsertCache[key] = cache
		liveOrderInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the LiveOrder.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *LiveOrder) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	liveOrderUpdateCacheMut.RLock()
	cache, cached := liveOrderUpdateCache[key]
	liveOrderUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			liveOrderAllColumns,
			liveOrderPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("postgres: unable to update live_order, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"live_order\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, liveOrderPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(liveOrderType, liveOrderMapping, append(wl, liveOrderPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to update live_order row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: failed to get rows affected by update for live_order")
	}

	if !cached {
		liveOrderUpdateCacheMut.Lock()
		liveOrderUpdateCache[key] = cache
		liveOrderUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q liveOrderQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to update all for live_order")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to retrieve rows affected for live_order")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o LiveOrderSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("postgres: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), liveOrderPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"live_order\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, liveOrderPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to update all in liveOrder slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to retrieve rows affected all in update all liveOrder")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *LiveOrder) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("postgres: no live_order provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		o.UpdatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(liveOrderColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	liveOrderUpsertCacheMut.RLock()
	cache, cached := liveOrderUpsertCache[key]
	liveOrderUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			liveOrderAllColumns,
			liveOrderColumnsWithDefault,
			liveOrderColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			liveOrderAllColumns,
			liveOrderPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("postgres: unable to upsert live_order, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(liveOrderPrimaryKeyColumns))
			copy(conflict, liveOrderPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"live_order\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(liveOrderType, liveOrderMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(liveOrderType, liveOrderMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "postgres: unable to upsert live_order")
	}

	if !cached {
		liveOrderUpsertCacheMut.Lock()
		liveOrderUpsertCache[key] = cache
		liveOrderUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single LiveOrder record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *LiveOrder) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("postgres: no LiveOrder provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), liveOrderPrimaryKeyMapping)
	sql := "DELETE FROM \"live_order\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to delete from live_order")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: failed to get rows affected by delete for live_order")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q liveOrderQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("postgres: no liveOrderQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to delete all from live_order")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: failed to get rows affected by deleteall for live_order")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o LiveOrderSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(liveOrderBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), liveOrderPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"live_order\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, liveOrderPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to delete all from liveOrder slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: failed to get rows affected by deleteall for live_order")
	}

	if len(liveOrderAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *LiveOrder) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindLiveOrder(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *LiveOrderSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := LiveOrderSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), liveOrderPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"live_order\".* FROM \"live_order\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, liveOrderPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "postgres: unable to reload all in LiveOrderSlice")
	}

	*o = slice

	return nil
}

// LiveOrderExists checks if the LiveOrder row exists.
func LiveOrderExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"live_order\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "postgres: unable to check if live_order exists")
	}

	return exists, nil
}
