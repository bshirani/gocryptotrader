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
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Strategy is an object representing the database table.
type Strategy struct {
	ID        int       `boil:"id" json:"id" toml:"id" yaml:"id"`
	Capture   string    `boil:"capture" json:"capture" toml:"capture" yaml:"capture"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *strategyR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L strategyL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var StrategyColumns = struct {
	ID        string
	Capture   string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "id",
	Capture:   "capture",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

var StrategyTableColumns = struct {
	ID        string
	Capture   string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "strategy.id",
	Capture:   "strategy.capture",
	CreatedAt: "strategy.created_at",
	UpdatedAt: "strategy.updated_at",
}

// Generated where

var StrategyWhere = struct {
	ID        whereHelperint
	Capture   whereHelperstring
	CreatedAt whereHelpertime_Time
	UpdatedAt whereHelpertime_Time
}{
	ID:        whereHelperint{field: "\"strategy\".\"id\""},
	Capture:   whereHelperstring{field: "\"strategy\".\"capture\""},
	CreatedAt: whereHelpertime_Time{field: "\"strategy\".\"created_at\""},
	UpdatedAt: whereHelpertime_Time{field: "\"strategy\".\"updated_at\""},
}

// StrategyRels is where relationship names are stored.
var StrategyRels = struct {
	CurrencyPairStrategies string
}{
	CurrencyPairStrategies: "CurrencyPairStrategies",
}

// strategyR is where relationships are stored.
type strategyR struct {
	CurrencyPairStrategies CurrencyPairStrategySlice `boil:"CurrencyPairStrategies" json:"CurrencyPairStrategies" toml:"CurrencyPairStrategies" yaml:"CurrencyPairStrategies"`
}

// NewStruct creates a new relationship struct
func (*strategyR) NewStruct() *strategyR {
	return &strategyR{}
}

// strategyL is where Load methods for each relationship are stored.
type strategyL struct{}

var (
	strategyAllColumns            = []string{"id", "capture", "created_at", "updated_at"}
	strategyColumnsWithoutDefault = []string{"capture"}
	strategyColumnsWithDefault    = []string{"id", "created_at", "updated_at"}
	strategyPrimaryKeyColumns     = []string{"id"}
)

type (
	// StrategySlice is an alias for a slice of pointers to Strategy.
	// This should almost always be used instead of []Strategy.
	StrategySlice []*Strategy
	// StrategyHook is the signature for custom Strategy hook methods
	StrategyHook func(context.Context, boil.ContextExecutor, *Strategy) error

	strategyQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	strategyType                 = reflect.TypeOf(&Strategy{})
	strategyMapping              = queries.MakeStructMapping(strategyType)
	strategyPrimaryKeyMapping, _ = queries.BindMapping(strategyType, strategyMapping, strategyPrimaryKeyColumns)
	strategyInsertCacheMut       sync.RWMutex
	strategyInsertCache          = make(map[string]insertCache)
	strategyUpdateCacheMut       sync.RWMutex
	strategyUpdateCache          = make(map[string]updateCache)
	strategyUpsertCacheMut       sync.RWMutex
	strategyUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var strategyBeforeInsertHooks []StrategyHook
var strategyBeforeUpdateHooks []StrategyHook
var strategyBeforeDeleteHooks []StrategyHook
var strategyBeforeUpsertHooks []StrategyHook

var strategyAfterInsertHooks []StrategyHook
var strategyAfterSelectHooks []StrategyHook
var strategyAfterUpdateHooks []StrategyHook
var strategyAfterDeleteHooks []StrategyHook
var strategyAfterUpsertHooks []StrategyHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Strategy) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range strategyBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Strategy) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range strategyBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Strategy) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range strategyBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Strategy) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range strategyBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Strategy) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range strategyAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Strategy) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range strategyAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Strategy) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range strategyAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Strategy) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range strategyAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Strategy) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range strategyAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddStrategyHook registers your hook function for all future operations.
func AddStrategyHook(hookPoint boil.HookPoint, strategyHook StrategyHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		strategyBeforeInsertHooks = append(strategyBeforeInsertHooks, strategyHook)
	case boil.BeforeUpdateHook:
		strategyBeforeUpdateHooks = append(strategyBeforeUpdateHooks, strategyHook)
	case boil.BeforeDeleteHook:
		strategyBeforeDeleteHooks = append(strategyBeforeDeleteHooks, strategyHook)
	case boil.BeforeUpsertHook:
		strategyBeforeUpsertHooks = append(strategyBeforeUpsertHooks, strategyHook)
	case boil.AfterInsertHook:
		strategyAfterInsertHooks = append(strategyAfterInsertHooks, strategyHook)
	case boil.AfterSelectHook:
		strategyAfterSelectHooks = append(strategyAfterSelectHooks, strategyHook)
	case boil.AfterUpdateHook:
		strategyAfterUpdateHooks = append(strategyAfterUpdateHooks, strategyHook)
	case boil.AfterDeleteHook:
		strategyAfterDeleteHooks = append(strategyAfterDeleteHooks, strategyHook)
	case boil.AfterUpsertHook:
		strategyAfterUpsertHooks = append(strategyAfterUpsertHooks, strategyHook)
	}
}

// One returns a single strategy record from the query.
func (q strategyQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Strategy, error) {
	o := &Strategy{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "postgres: failed to execute a one query for strategy")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Strategy records from the query.
func (q strategyQuery) All(ctx context.Context, exec boil.ContextExecutor) (StrategySlice, error) {
	var o []*Strategy

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "postgres: failed to assign all query results to Strategy slice")
	}

	if len(strategyAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Strategy records in the query.
func (q strategyQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: failed to count strategy rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q strategyQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "postgres: failed to check if strategy exists")
	}

	return count > 0, nil
}

// CurrencyPairStrategies retrieves all the currency_pair_strategy's CurrencyPairStrategies with an executor.
func (o *Strategy) CurrencyPairStrategies(mods ...qm.QueryMod) currencyPairStrategyQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"currency_pair_strategy\".\"strategy_id\"=?", o.ID),
	)

	query := CurrencyPairStrategies(queryMods...)
	queries.SetFrom(query.Query, "\"currency_pair_strategy\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"currency_pair_strategy\".*"})
	}

	return query
}

// LoadCurrencyPairStrategies allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (strategyL) LoadCurrencyPairStrategies(ctx context.Context, e boil.ContextExecutor, singular bool, maybeStrategy interface{}, mods queries.Applicator) error {
	var slice []*Strategy
	var object *Strategy

	if singular {
		object = maybeStrategy.(*Strategy)
	} else {
		slice = *maybeStrategy.(*[]*Strategy)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &strategyR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &strategyR{}
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
		qm.From(`currency_pair_strategy`),
		qm.WhereIn(`currency_pair_strategy.strategy_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load currency_pair_strategy")
	}

	var resultSlice []*CurrencyPairStrategy
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice currency_pair_strategy")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on currency_pair_strategy")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for currency_pair_strategy")
	}

	if len(currencyPairStrategyAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.CurrencyPairStrategies = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &currencyPairStrategyR{}
			}
			foreign.R.Strategy = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.StrategyID {
				local.R.CurrencyPairStrategies = append(local.R.CurrencyPairStrategies, foreign)
				if foreign.R == nil {
					foreign.R = &currencyPairStrategyR{}
				}
				foreign.R.Strategy = local
				break
			}
		}
	}

	return nil
}

// AddCurrencyPairStrategies adds the given related objects to the existing relationships
// of the strategy, optionally inserting them as new records.
// Appends related to o.R.CurrencyPairStrategies.
// Sets related.R.Strategy appropriately.
func (o *Strategy) AddCurrencyPairStrategies(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*CurrencyPairStrategy) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.StrategyID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"currency_pair_strategy\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"strategy_id"}),
				strmangle.WhereClause("\"", "\"", 2, currencyPairStrategyPrimaryKeyColumns),
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

			rel.StrategyID = o.ID
		}
	}

	if o.R == nil {
		o.R = &strategyR{
			CurrencyPairStrategies: related,
		}
	} else {
		o.R.CurrencyPairStrategies = append(o.R.CurrencyPairStrategies, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &currencyPairStrategyR{
				Strategy: o,
			}
		} else {
			rel.R.Strategy = o
		}
	}
	return nil
}

// Strategies retrieves all the records using an executor.
func Strategies(mods ...qm.QueryMod) strategyQuery {
	mods = append(mods, qm.From("\"strategy\""))
	return strategyQuery{NewQuery(mods...)}
}

// FindStrategy retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindStrategy(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*Strategy, error) {
	strategyObj := &Strategy{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"strategy\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, strategyObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "postgres: unable to select from strategy")
	}

	if err = strategyObj.doAfterSelectHooks(ctx, exec); err != nil {
		return strategyObj, err
	}

	return strategyObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Strategy) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("postgres: no strategy provided for insertion")
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

	nzDefaults := queries.NonZeroDefaultSet(strategyColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	strategyInsertCacheMut.RLock()
	cache, cached := strategyInsertCache[key]
	strategyInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			strategyAllColumns,
			strategyColumnsWithDefault,
			strategyColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(strategyType, strategyMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(strategyType, strategyMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"strategy\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"strategy\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "postgres: unable to insert into strategy")
	}

	if !cached {
		strategyInsertCacheMut.Lock()
		strategyInsertCache[key] = cache
		strategyInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Strategy.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Strategy) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	strategyUpdateCacheMut.RLock()
	cache, cached := strategyUpdateCache[key]
	strategyUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			strategyAllColumns,
			strategyPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("postgres: unable to update strategy, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"strategy\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, strategyPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(strategyType, strategyMapping, append(wl, strategyPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "postgres: unable to update strategy row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: failed to get rows affected by update for strategy")
	}

	if !cached {
		strategyUpdateCacheMut.Lock()
		strategyUpdateCache[key] = cache
		strategyUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q strategyQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to update all for strategy")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to retrieve rows affected for strategy")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o StrategySlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), strategyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"strategy\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, strategyPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to update all in strategy slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to retrieve rows affected all in update all strategy")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Strategy) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("postgres: no strategy provided for upsert")
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

	nzDefaults := queries.NonZeroDefaultSet(strategyColumnsWithDefault, o)

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

	strategyUpsertCacheMut.RLock()
	cache, cached := strategyUpsertCache[key]
	strategyUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			strategyAllColumns,
			strategyColumnsWithDefault,
			strategyColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			strategyAllColumns,
			strategyPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("postgres: unable to upsert strategy, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(strategyPrimaryKeyColumns))
			copy(conflict, strategyPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"strategy\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(strategyType, strategyMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(strategyType, strategyMapping, ret)
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
		return errors.Wrap(err, "postgres: unable to upsert strategy")
	}

	if !cached {
		strategyUpsertCacheMut.Lock()
		strategyUpsertCache[key] = cache
		strategyUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Strategy record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Strategy) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("postgres: no Strategy provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), strategyPrimaryKeyMapping)
	sql := "DELETE FROM \"strategy\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to delete from strategy")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: failed to get rows affected by delete for strategy")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q strategyQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("postgres: no strategyQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to delete all from strategy")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: failed to get rows affected by deleteall for strategy")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o StrategySlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(strategyBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), strategyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"strategy\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, strategyPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to delete all from strategy slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: failed to get rows affected by deleteall for strategy")
	}

	if len(strategyAfterDeleteHooks) != 0 {
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
func (o *Strategy) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindStrategy(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *StrategySlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := StrategySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), strategyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"strategy\".* FROM \"strategy\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, strategyPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "postgres: unable to reload all in StrategySlice")
	}

	*o = slice

	return nil
}

// StrategyExists checks if the Strategy row exists.
func StrategyExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"strategy\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "postgres: unable to check if strategy exists")
	}

	return exists, nil
}
