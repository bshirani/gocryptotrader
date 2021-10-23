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

// CurrencyPairStrategy is an object representing the database table.
type CurrencyPairStrategy struct {
	ID             int       `boil:"id" json:"id" toml:"id" yaml:"id"`
	CurrencyPairID int       `boil:"currency_pair_id" json:"currency_pair_id" toml:"currency_pair_id" yaml:"currency_pair_id"`
	StrategyID     int       `boil:"strategy_id" json:"strategy_id" toml:"strategy_id" yaml:"strategy_id"`
	Side           string    `boil:"side" json:"side" toml:"side" yaml:"side"`
	Active         bool      `boil:"active" json:"active" toml:"active" yaml:"active"`
	CreatedAt      time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt      time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *currencyPairStrategyR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L currencyPairStrategyL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var CurrencyPairStrategyColumns = struct {
	ID             string
	CurrencyPairID string
	StrategyID     string
	Side           string
	Active         string
	CreatedAt      string
	UpdatedAt      string
}{
	ID:             "id",
	CurrencyPairID: "currency_pair_id",
	StrategyID:     "strategy_id",
	Side:           "side",
	Active:         "active",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

var CurrencyPairStrategyTableColumns = struct {
	ID             string
	CurrencyPairID string
	StrategyID     string
	Side           string
	Active         string
	CreatedAt      string
	UpdatedAt      string
}{
	ID:             "currency_pair_strategy.id",
	CurrencyPairID: "currency_pair_strategy.currency_pair_id",
	StrategyID:     "currency_pair_strategy.strategy_id",
	Side:           "currency_pair_strategy.side",
	Active:         "currency_pair_strategy.active",
	CreatedAt:      "currency_pair_strategy.created_at",
	UpdatedAt:      "currency_pair_strategy.updated_at",
}

// Generated where

var CurrencyPairStrategyWhere = struct {
	ID             whereHelperint
	CurrencyPairID whereHelperint
	StrategyID     whereHelperint
	Side           whereHelperstring
	Active         whereHelperbool
	CreatedAt      whereHelpertime_Time
	UpdatedAt      whereHelpertime_Time
}{
	ID:             whereHelperint{field: "\"currency_pair_strategy\".\"id\""},
	CurrencyPairID: whereHelperint{field: "\"currency_pair_strategy\".\"currency_pair_id\""},
	StrategyID:     whereHelperint{field: "\"currency_pair_strategy\".\"strategy_id\""},
	Side:           whereHelperstring{field: "\"currency_pair_strategy\".\"side\""},
	Active:         whereHelperbool{field: "\"currency_pair_strategy\".\"active\""},
	CreatedAt:      whereHelpertime_Time{field: "\"currency_pair_strategy\".\"created_at\""},
	UpdatedAt:      whereHelpertime_Time{field: "\"currency_pair_strategy\".\"updated_at\""},
}

// CurrencyPairStrategyRels is where relationship names are stored.
var CurrencyPairStrategyRels = struct {
	CurrencyPair string
	Strategy     string
}{
	CurrencyPair: "CurrencyPair",
	Strategy:     "Strategy",
}

// currencyPairStrategyR is where relationships are stored.
type currencyPairStrategyR struct {
	CurrencyPair *Currency `boil:"CurrencyPair" json:"CurrencyPair" toml:"CurrencyPair" yaml:"CurrencyPair"`
	Strategy     *Strategy `boil:"Strategy" json:"Strategy" toml:"Strategy" yaml:"Strategy"`
}

// NewStruct creates a new relationship struct
func (*currencyPairStrategyR) NewStruct() *currencyPairStrategyR {
	return &currencyPairStrategyR{}
}

// currencyPairStrategyL is where Load methods for each relationship are stored.
type currencyPairStrategyL struct{}

var (
	currencyPairStrategyAllColumns            = []string{"id", "currency_pair_id", "strategy_id", "side", "active", "created_at", "updated_at"}
	currencyPairStrategyColumnsWithoutDefault = []string{"currency_pair_id", "strategy_id", "side"}
	currencyPairStrategyColumnsWithDefault    = []string{"id", "active", "created_at", "updated_at"}
	currencyPairStrategyPrimaryKeyColumns     = []string{"id"}
)

type (
	// CurrencyPairStrategySlice is an alias for a slice of pointers to CurrencyPairStrategy.
	// This should almost always be used instead of []CurrencyPairStrategy.
	CurrencyPairStrategySlice []*CurrencyPairStrategy
	// CurrencyPairStrategyHook is the signature for custom CurrencyPairStrategy hook methods
	CurrencyPairStrategyHook func(context.Context, boil.ContextExecutor, *CurrencyPairStrategy) error

	currencyPairStrategyQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	currencyPairStrategyType                 = reflect.TypeOf(&CurrencyPairStrategy{})
	currencyPairStrategyMapping              = queries.MakeStructMapping(currencyPairStrategyType)
	currencyPairStrategyPrimaryKeyMapping, _ = queries.BindMapping(currencyPairStrategyType, currencyPairStrategyMapping, currencyPairStrategyPrimaryKeyColumns)
	currencyPairStrategyInsertCacheMut       sync.RWMutex
	currencyPairStrategyInsertCache          = make(map[string]insertCache)
	currencyPairStrategyUpdateCacheMut       sync.RWMutex
	currencyPairStrategyUpdateCache          = make(map[string]updateCache)
	currencyPairStrategyUpsertCacheMut       sync.RWMutex
	currencyPairStrategyUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var currencyPairStrategyBeforeInsertHooks []CurrencyPairStrategyHook
var currencyPairStrategyBeforeUpdateHooks []CurrencyPairStrategyHook
var currencyPairStrategyBeforeDeleteHooks []CurrencyPairStrategyHook
var currencyPairStrategyBeforeUpsertHooks []CurrencyPairStrategyHook

var currencyPairStrategyAfterInsertHooks []CurrencyPairStrategyHook
var currencyPairStrategyAfterSelectHooks []CurrencyPairStrategyHook
var currencyPairStrategyAfterUpdateHooks []CurrencyPairStrategyHook
var currencyPairStrategyAfterDeleteHooks []CurrencyPairStrategyHook
var currencyPairStrategyAfterUpsertHooks []CurrencyPairStrategyHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *CurrencyPairStrategy) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range currencyPairStrategyBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *CurrencyPairStrategy) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range currencyPairStrategyBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *CurrencyPairStrategy) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range currencyPairStrategyBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *CurrencyPairStrategy) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range currencyPairStrategyBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *CurrencyPairStrategy) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range currencyPairStrategyAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *CurrencyPairStrategy) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range currencyPairStrategyAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *CurrencyPairStrategy) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range currencyPairStrategyAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *CurrencyPairStrategy) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range currencyPairStrategyAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *CurrencyPairStrategy) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range currencyPairStrategyAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddCurrencyPairStrategyHook registers your hook function for all future operations.
func AddCurrencyPairStrategyHook(hookPoint boil.HookPoint, currencyPairStrategyHook CurrencyPairStrategyHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		currencyPairStrategyBeforeInsertHooks = append(currencyPairStrategyBeforeInsertHooks, currencyPairStrategyHook)
	case boil.BeforeUpdateHook:
		currencyPairStrategyBeforeUpdateHooks = append(currencyPairStrategyBeforeUpdateHooks, currencyPairStrategyHook)
	case boil.BeforeDeleteHook:
		currencyPairStrategyBeforeDeleteHooks = append(currencyPairStrategyBeforeDeleteHooks, currencyPairStrategyHook)
	case boil.BeforeUpsertHook:
		currencyPairStrategyBeforeUpsertHooks = append(currencyPairStrategyBeforeUpsertHooks, currencyPairStrategyHook)
	case boil.AfterInsertHook:
		currencyPairStrategyAfterInsertHooks = append(currencyPairStrategyAfterInsertHooks, currencyPairStrategyHook)
	case boil.AfterSelectHook:
		currencyPairStrategyAfterSelectHooks = append(currencyPairStrategyAfterSelectHooks, currencyPairStrategyHook)
	case boil.AfterUpdateHook:
		currencyPairStrategyAfterUpdateHooks = append(currencyPairStrategyAfterUpdateHooks, currencyPairStrategyHook)
	case boil.AfterDeleteHook:
		currencyPairStrategyAfterDeleteHooks = append(currencyPairStrategyAfterDeleteHooks, currencyPairStrategyHook)
	case boil.AfterUpsertHook:
		currencyPairStrategyAfterUpsertHooks = append(currencyPairStrategyAfterUpsertHooks, currencyPairStrategyHook)
	}
}

// One returns a single currencyPairStrategy record from the query.
func (q currencyPairStrategyQuery) One(ctx context.Context, exec boil.ContextExecutor) (*CurrencyPairStrategy, error) {
	o := &CurrencyPairStrategy{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "postgres: failed to execute a one query for currency_pair_strategy")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all CurrencyPairStrategy records from the query.
func (q currencyPairStrategyQuery) All(ctx context.Context, exec boil.ContextExecutor) (CurrencyPairStrategySlice, error) {
	var o []*CurrencyPairStrategy

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "postgres: failed to assign all query results to CurrencyPairStrategy slice")
	}

	if len(currencyPairStrategyAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all CurrencyPairStrategy records in the query.
func (q currencyPairStrategyQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: failed to count currency_pair_strategy rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q currencyPairStrategyQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "postgres: failed to check if currency_pair_strategy exists")
	}

	return count > 0, nil
}

// CurrencyPair pointed to by the foreign key.
func (o *CurrencyPairStrategy) CurrencyPair(mods ...qm.QueryMod) currencyQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.CurrencyPairID),
	}

	queryMods = append(queryMods, mods...)

	query := Currencies(queryMods...)
	queries.SetFrom(query.Query, "\"currency\"")

	return query
}

// Strategy pointed to by the foreign key.
func (o *CurrencyPairStrategy) Strategy(mods ...qm.QueryMod) strategyQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.StrategyID),
	}

	queryMods = append(queryMods, mods...)

	query := Strategies(queryMods...)
	queries.SetFrom(query.Query, "\"strategy\"")

	return query
}

// LoadCurrencyPair allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (currencyPairStrategyL) LoadCurrencyPair(ctx context.Context, e boil.ContextExecutor, singular bool, maybeCurrencyPairStrategy interface{}, mods queries.Applicator) error {
	var slice []*CurrencyPairStrategy
	var object *CurrencyPairStrategy

	if singular {
		object = maybeCurrencyPairStrategy.(*CurrencyPairStrategy)
	} else {
		slice = *maybeCurrencyPairStrategy.(*[]*CurrencyPairStrategy)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &currencyPairStrategyR{}
		}
		args = append(args, object.CurrencyPairID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &currencyPairStrategyR{}
			}

			for _, a := range args {
				if a == obj.CurrencyPairID {
					continue Outer
				}
			}

			args = append(args, obj.CurrencyPairID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`currency`),
		qm.WhereIn(`currency.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Currency")
	}

	var resultSlice []*Currency
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Currency")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for currency")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for currency")
	}

	if len(currencyPairStrategyAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.CurrencyPair = foreign
		if foreign.R == nil {
			foreign.R = &currencyR{}
		}
		foreign.R.CurrencyPairCurrencyPairStrategies = append(foreign.R.CurrencyPairCurrencyPairStrategies, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.CurrencyPairID == foreign.ID {
				local.R.CurrencyPair = foreign
				if foreign.R == nil {
					foreign.R = &currencyR{}
				}
				foreign.R.CurrencyPairCurrencyPairStrategies = append(foreign.R.CurrencyPairCurrencyPairStrategies, local)
				break
			}
		}
	}

	return nil
}

// LoadStrategy allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (currencyPairStrategyL) LoadStrategy(ctx context.Context, e boil.ContextExecutor, singular bool, maybeCurrencyPairStrategy interface{}, mods queries.Applicator) error {
	var slice []*CurrencyPairStrategy
	var object *CurrencyPairStrategy

	if singular {
		object = maybeCurrencyPairStrategy.(*CurrencyPairStrategy)
	} else {
		slice = *maybeCurrencyPairStrategy.(*[]*CurrencyPairStrategy)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &currencyPairStrategyR{}
		}
		args = append(args, object.StrategyID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &currencyPairStrategyR{}
			}

			for _, a := range args {
				if a == obj.StrategyID {
					continue Outer
				}
			}

			args = append(args, obj.StrategyID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`strategy`),
		qm.WhereIn(`strategy.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Strategy")
	}

	var resultSlice []*Strategy
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Strategy")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for strategy")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for strategy")
	}

	if len(currencyPairStrategyAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Strategy = foreign
		if foreign.R == nil {
			foreign.R = &strategyR{}
		}
		foreign.R.CurrencyPairStrategies = append(foreign.R.CurrencyPairStrategies, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.StrategyID == foreign.ID {
				local.R.Strategy = foreign
				if foreign.R == nil {
					foreign.R = &strategyR{}
				}
				foreign.R.CurrencyPairStrategies = append(foreign.R.CurrencyPairStrategies, local)
				break
			}
		}
	}

	return nil
}

// SetCurrencyPair of the currencyPairStrategy to the related item.
// Sets o.R.CurrencyPair to related.
// Adds o to related.R.CurrencyPairCurrencyPairStrategies.
func (o *CurrencyPairStrategy) SetCurrencyPair(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Currency) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"currency_pair_strategy\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"currency_pair_id"}),
		strmangle.WhereClause("\"", "\"", 2, currencyPairStrategyPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.CurrencyPairID = related.ID
	if o.R == nil {
		o.R = &currencyPairStrategyR{
			CurrencyPair: related,
		}
	} else {
		o.R.CurrencyPair = related
	}

	if related.R == nil {
		related.R = &currencyR{
			CurrencyPairCurrencyPairStrategies: CurrencyPairStrategySlice{o},
		}
	} else {
		related.R.CurrencyPairCurrencyPairStrategies = append(related.R.CurrencyPairCurrencyPairStrategies, o)
	}

	return nil
}

// SetStrategy of the currencyPairStrategy to the related item.
// Sets o.R.Strategy to related.
// Adds o to related.R.CurrencyPairStrategies.
func (o *CurrencyPairStrategy) SetStrategy(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Strategy) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"currency_pair_strategy\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"strategy_id"}),
		strmangle.WhereClause("\"", "\"", 2, currencyPairStrategyPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.StrategyID = related.ID
	if o.R == nil {
		o.R = &currencyPairStrategyR{
			Strategy: related,
		}
	} else {
		o.R.Strategy = related
	}

	if related.R == nil {
		related.R = &strategyR{
			CurrencyPairStrategies: CurrencyPairStrategySlice{o},
		}
	} else {
		related.R.CurrencyPairStrategies = append(related.R.CurrencyPairStrategies, o)
	}

	return nil
}

// CurrencyPairStrategies retrieves all the records using an executor.
func CurrencyPairStrategies(mods ...qm.QueryMod) currencyPairStrategyQuery {
	mods = append(mods, qm.From("\"currency_pair_strategy\""))
	return currencyPairStrategyQuery{NewQuery(mods...)}
}

// FindCurrencyPairStrategy retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindCurrencyPairStrategy(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*CurrencyPairStrategy, error) {
	currencyPairStrategyObj := &CurrencyPairStrategy{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"currency_pair_strategy\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, currencyPairStrategyObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "postgres: unable to select from currency_pair_strategy")
	}

	if err = currencyPairStrategyObj.doAfterSelectHooks(ctx, exec); err != nil {
		return currencyPairStrategyObj, err
	}

	return currencyPairStrategyObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *CurrencyPairStrategy) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("postgres: no currency_pair_strategy provided for insertion")
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

	nzDefaults := queries.NonZeroDefaultSet(currencyPairStrategyColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	currencyPairStrategyInsertCacheMut.RLock()
	cache, cached := currencyPairStrategyInsertCache[key]
	currencyPairStrategyInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			currencyPairStrategyAllColumns,
			currencyPairStrategyColumnsWithDefault,
			currencyPairStrategyColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(currencyPairStrategyType, currencyPairStrategyMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(currencyPairStrategyType, currencyPairStrategyMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"currency_pair_strategy\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"currency_pair_strategy\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "postgres: unable to insert into currency_pair_strategy")
	}

	if !cached {
		currencyPairStrategyInsertCacheMut.Lock()
		currencyPairStrategyInsertCache[key] = cache
		currencyPairStrategyInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the CurrencyPairStrategy.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *CurrencyPairStrategy) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	currencyPairStrategyUpdateCacheMut.RLock()
	cache, cached := currencyPairStrategyUpdateCache[key]
	currencyPairStrategyUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			currencyPairStrategyAllColumns,
			currencyPairStrategyPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("postgres: unable to update currency_pair_strategy, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"currency_pair_strategy\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, currencyPairStrategyPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(currencyPairStrategyType, currencyPairStrategyMapping, append(wl, currencyPairStrategyPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "postgres: unable to update currency_pair_strategy row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: failed to get rows affected by update for currency_pair_strategy")
	}

	if !cached {
		currencyPairStrategyUpdateCacheMut.Lock()
		currencyPairStrategyUpdateCache[key] = cache
		currencyPairStrategyUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q currencyPairStrategyQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to update all for currency_pair_strategy")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to retrieve rows affected for currency_pair_strategy")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o CurrencyPairStrategySlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), currencyPairStrategyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"currency_pair_strategy\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, currencyPairStrategyPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to update all in currencyPairStrategy slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to retrieve rows affected all in update all currencyPairStrategy")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *CurrencyPairStrategy) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("postgres: no currency_pair_strategy provided for upsert")
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

	nzDefaults := queries.NonZeroDefaultSet(currencyPairStrategyColumnsWithDefault, o)

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

	currencyPairStrategyUpsertCacheMut.RLock()
	cache, cached := currencyPairStrategyUpsertCache[key]
	currencyPairStrategyUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			currencyPairStrategyAllColumns,
			currencyPairStrategyColumnsWithDefault,
			currencyPairStrategyColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			currencyPairStrategyAllColumns,
			currencyPairStrategyPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("postgres: unable to upsert currency_pair_strategy, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(currencyPairStrategyPrimaryKeyColumns))
			copy(conflict, currencyPairStrategyPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"currency_pair_strategy\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(currencyPairStrategyType, currencyPairStrategyMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(currencyPairStrategyType, currencyPairStrategyMapping, ret)
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
		return errors.Wrap(err, "postgres: unable to upsert currency_pair_strategy")
	}

	if !cached {
		currencyPairStrategyUpsertCacheMut.Lock()
		currencyPairStrategyUpsertCache[key] = cache
		currencyPairStrategyUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single CurrencyPairStrategy record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *CurrencyPairStrategy) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("postgres: no CurrencyPairStrategy provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), currencyPairStrategyPrimaryKeyMapping)
	sql := "DELETE FROM \"currency_pair_strategy\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to delete from currency_pair_strategy")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: failed to get rows affected by delete for currency_pair_strategy")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q currencyPairStrategyQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("postgres: no currencyPairStrategyQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to delete all from currency_pair_strategy")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: failed to get rows affected by deleteall for currency_pair_strategy")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o CurrencyPairStrategySlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(currencyPairStrategyBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), currencyPairStrategyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"currency_pair_strategy\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, currencyPairStrategyPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to delete all from currencyPairStrategy slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: failed to get rows affected by deleteall for currency_pair_strategy")
	}

	if len(currencyPairStrategyAfterDeleteHooks) != 0 {
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
func (o *CurrencyPairStrategy) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindCurrencyPairStrategy(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CurrencyPairStrategySlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := CurrencyPairStrategySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), currencyPairStrategyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"currency_pair_strategy\".* FROM \"currency_pair_strategy\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, currencyPairStrategyPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "postgres: unable to reload all in CurrencyPairStrategySlice")
	}

	*o = slice

	return nil
}

// CurrencyPairStrategyExists checks if the CurrencyPairStrategy row exists.
func CurrencyPairStrategyExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"currency_pair_strategy\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "postgres: unable to check if currency_pair_strategy exists")
	}

	return exists, nil
}
