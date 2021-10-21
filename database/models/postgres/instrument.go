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
	"github.com/volatiletech/sqlboiler/v4/types"
	"github.com/volatiletech/strmangle"
)

// Instrument is an object representing the database table.
type Instrument struct {
	ID                  int               `boil:"id" json:"id" toml:"id" yaml:"id"`
	Symbol              string            `boil:"symbol" json:"symbol" toml:"symbol" yaml:"symbol"`
	CMCID               int               `boil:"cmc_id" json:"cmc_id" toml:"cmc_id" yaml:"cmc_id"`
	Name                string            `boil:"name" json:"name" toml:"name" yaml:"name"`
	Slug                string            `boil:"slug" json:"slug" toml:"slug" yaml:"slug"`
	FirstHistoricalData time.Time         `boil:"first_historical_data" json:"first_historical_data" toml:"first_historical_data" yaml:"first_historical_data"`
	LastHistoricalData  time.Time         `boil:"last_historical_data" json:"last_historical_data" toml:"last_historical_data" yaml:"last_historical_data"`
	MarketCap           types.NullDecimal `boil:"market_cap" json:"market_cap,omitempty" toml:"market_cap" yaml:"market_cap,omitempty"`
	ListingStatus       string            `boil:"listing_status" json:"listing_status" toml:"listing_status" yaml:"listing_status"`
	Active              bool              `boil:"active" json:"active" toml:"active" yaml:"active"`
	Status              bool              `boil:"status" json:"status" toml:"status" yaml:"status"`
	CreatedAt           time.Time         `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt           time.Time         `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *instrumentR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L instrumentL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var InstrumentColumns = struct {
	ID                  string
	Symbol              string
	CMCID               string
	Name                string
	Slug                string
	FirstHistoricalData string
	LastHistoricalData  string
	MarketCap           string
	ListingStatus       string
	Active              string
	Status              string
	CreatedAt           string
	UpdatedAt           string
}{
	ID:                  "id",
	Symbol:              "symbol",
	CMCID:               "cmc_id",
	Name:                "name",
	Slug:                "slug",
	FirstHistoricalData: "first_historical_data",
	LastHistoricalData:  "last_historical_data",
	MarketCap:           "market_cap",
	ListingStatus:       "listing_status",
	Active:              "active",
	Status:              "status",
	CreatedAt:           "created_at",
	UpdatedAt:           "updated_at",
}

var InstrumentTableColumns = struct {
	ID                  string
	Symbol              string
	CMCID               string
	Name                string
	Slug                string
	FirstHistoricalData string
	LastHistoricalData  string
	MarketCap           string
	ListingStatus       string
	Active              string
	Status              string
	CreatedAt           string
	UpdatedAt           string
}{
	ID:                  "instrument.id",
	Symbol:              "instrument.symbol",
	CMCID:               "instrument.cmc_id",
	Name:                "instrument.name",
	Slug:                "instrument.slug",
	FirstHistoricalData: "instrument.first_historical_data",
	LastHistoricalData:  "instrument.last_historical_data",
	MarketCap:           "instrument.market_cap",
	ListingStatus:       "instrument.listing_status",
	Active:              "instrument.active",
	Status:              "instrument.status",
	CreatedAt:           "instrument.created_at",
	UpdatedAt:           "instrument.updated_at",
}

// Generated where

type whereHelpertypes_NullDecimal struct{ field string }

func (w whereHelpertypes_NullDecimal) EQ(x types.NullDecimal) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpertypes_NullDecimal) NEQ(x types.NullDecimal) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpertypes_NullDecimal) LT(x types.NullDecimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertypes_NullDecimal) LTE(x types.NullDecimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertypes_NullDecimal) GT(x types.NullDecimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertypes_NullDecimal) GTE(x types.NullDecimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

func (w whereHelpertypes_NullDecimal) IsNull() qm.QueryMod { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpertypes_NullDecimal) IsNotNull() qm.QueryMod {
	return qmhelper.WhereIsNotNull(w.field)
}

var InstrumentWhere = struct {
	ID                  whereHelperint
	Symbol              whereHelperstring
	CMCID               whereHelperint
	Name                whereHelperstring
	Slug                whereHelperstring
	FirstHistoricalData whereHelpertime_Time
	LastHistoricalData  whereHelpertime_Time
	MarketCap           whereHelpertypes_NullDecimal
	ListingStatus       whereHelperstring
	Active              whereHelperbool
	Status              whereHelperbool
	CreatedAt           whereHelpertime_Time
	UpdatedAt           whereHelpertime_Time
}{
	ID:                  whereHelperint{field: "\"instrument\".\"id\""},
	Symbol:              whereHelperstring{field: "\"instrument\".\"symbol\""},
	CMCID:               whereHelperint{field: "\"instrument\".\"cmc_id\""},
	Name:                whereHelperstring{field: "\"instrument\".\"name\""},
	Slug:                whereHelperstring{field: "\"instrument\".\"slug\""},
	FirstHistoricalData: whereHelpertime_Time{field: "\"instrument\".\"first_historical_data\""},
	LastHistoricalData:  whereHelpertime_Time{field: "\"instrument\".\"last_historical_data\""},
	MarketCap:           whereHelpertypes_NullDecimal{field: "\"instrument\".\"market_cap\""},
	ListingStatus:       whereHelperstring{field: "\"instrument\".\"listing_status\""},
	Active:              whereHelperbool{field: "\"instrument\".\"active\""},
	Status:              whereHelperbool{field: "\"instrument\".\"status\""},
	CreatedAt:           whereHelpertime_Time{field: "\"instrument\".\"created_at\""},
	UpdatedAt:           whereHelpertime_Time{field: "\"instrument\".\"updated_at\""},
}

// InstrumentRels is where relationship names are stored.
var InstrumentRels = struct {
}{}

// instrumentR is where relationships are stored.
type instrumentR struct {
}

// NewStruct creates a new relationship struct
func (*instrumentR) NewStruct() *instrumentR {
	return &instrumentR{}
}

// instrumentL is where Load methods for each relationship are stored.
type instrumentL struct{}

var (
	instrumentAllColumns            = []string{"id", "symbol", "cmc_id", "name", "slug", "first_historical_data", "last_historical_data", "market_cap", "listing_status", "active", "status", "created_at", "updated_at"}
	instrumentColumnsWithoutDefault = []string{"symbol", "cmc_id", "name", "slug", "first_historical_data", "last_historical_data", "market_cap", "listing_status", "active", "status"}
	instrumentColumnsWithDefault    = []string{"id", "created_at", "updated_at"}
	instrumentPrimaryKeyColumns     = []string{"id"}
)

type (
	// InstrumentSlice is an alias for a slice of pointers to Instrument.
	// This should almost always be used instead of []Instrument.
	InstrumentSlice []*Instrument
	// InstrumentHook is the signature for custom Instrument hook methods
	InstrumentHook func(context.Context, boil.ContextExecutor, *Instrument) error

	instrumentQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	instrumentType                 = reflect.TypeOf(&Instrument{})
	instrumentMapping              = queries.MakeStructMapping(instrumentType)
	instrumentPrimaryKeyMapping, _ = queries.BindMapping(instrumentType, instrumentMapping, instrumentPrimaryKeyColumns)
	instrumentInsertCacheMut       sync.RWMutex
	instrumentInsertCache          = make(map[string]insertCache)
	instrumentUpdateCacheMut       sync.RWMutex
	instrumentUpdateCache          = make(map[string]updateCache)
	instrumentUpsertCacheMut       sync.RWMutex
	instrumentUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var instrumentBeforeInsertHooks []InstrumentHook
var instrumentBeforeUpdateHooks []InstrumentHook
var instrumentBeforeDeleteHooks []InstrumentHook
var instrumentBeforeUpsertHooks []InstrumentHook

var instrumentAfterInsertHooks []InstrumentHook
var instrumentAfterSelectHooks []InstrumentHook
var instrumentAfterUpdateHooks []InstrumentHook
var instrumentAfterDeleteHooks []InstrumentHook
var instrumentAfterUpsertHooks []InstrumentHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Instrument) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range instrumentBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Instrument) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range instrumentBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Instrument) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range instrumentBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Instrument) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range instrumentBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Instrument) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range instrumentAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Instrument) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range instrumentAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Instrument) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range instrumentAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Instrument) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range instrumentAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Instrument) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range instrumentAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddInstrumentHook registers your hook function for all future operations.
func AddInstrumentHook(hookPoint boil.HookPoint, instrumentHook InstrumentHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		instrumentBeforeInsertHooks = append(instrumentBeforeInsertHooks, instrumentHook)
	case boil.BeforeUpdateHook:
		instrumentBeforeUpdateHooks = append(instrumentBeforeUpdateHooks, instrumentHook)
	case boil.BeforeDeleteHook:
		instrumentBeforeDeleteHooks = append(instrumentBeforeDeleteHooks, instrumentHook)
	case boil.BeforeUpsertHook:
		instrumentBeforeUpsertHooks = append(instrumentBeforeUpsertHooks, instrumentHook)
	case boil.AfterInsertHook:
		instrumentAfterInsertHooks = append(instrumentAfterInsertHooks, instrumentHook)
	case boil.AfterSelectHook:
		instrumentAfterSelectHooks = append(instrumentAfterSelectHooks, instrumentHook)
	case boil.AfterUpdateHook:
		instrumentAfterUpdateHooks = append(instrumentAfterUpdateHooks, instrumentHook)
	case boil.AfterDeleteHook:
		instrumentAfterDeleteHooks = append(instrumentAfterDeleteHooks, instrumentHook)
	case boil.AfterUpsertHook:
		instrumentAfterUpsertHooks = append(instrumentAfterUpsertHooks, instrumentHook)
	}
}

// One returns a single instrument record from the query.
func (q instrumentQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Instrument, error) {
	o := &Instrument{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "postgres: failed to execute a one query for instrument")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Instrument records from the query.
func (q instrumentQuery) All(ctx context.Context, exec boil.ContextExecutor) (InstrumentSlice, error) {
	var o []*Instrument

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "postgres: failed to assign all query results to Instrument slice")
	}

	if len(instrumentAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Instrument records in the query.
func (q instrumentQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: failed to count instrument rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q instrumentQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "postgres: failed to check if instrument exists")
	}

	return count > 0, nil
}

// Instruments retrieves all the records using an executor.
func Instruments(mods ...qm.QueryMod) instrumentQuery {
	mods = append(mods, qm.From("\"instrument\""))
	return instrumentQuery{NewQuery(mods...)}
}

// FindInstrument retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindInstrument(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*Instrument, error) {
	instrumentObj := &Instrument{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"instrument\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, instrumentObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "postgres: unable to select from instrument")
	}

	if err = instrumentObj.doAfterSelectHooks(ctx, exec); err != nil {
		return instrumentObj, err
	}

	return instrumentObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Instrument) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("postgres: no instrument provided for insertion")
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

	nzDefaults := queries.NonZeroDefaultSet(instrumentColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	instrumentInsertCacheMut.RLock()
	cache, cached := instrumentInsertCache[key]
	instrumentInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			instrumentAllColumns,
			instrumentColumnsWithDefault,
			instrumentColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(instrumentType, instrumentMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(instrumentType, instrumentMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"instrument\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"instrument\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "postgres: unable to insert into instrument")
	}

	if !cached {
		instrumentInsertCacheMut.Lock()
		instrumentInsertCache[key] = cache
		instrumentInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Instrument.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Instrument) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	instrumentUpdateCacheMut.RLock()
	cache, cached := instrumentUpdateCache[key]
	instrumentUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			instrumentAllColumns,
			instrumentPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("postgres: unable to update instrument, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"instrument\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, instrumentPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(instrumentType, instrumentMapping, append(wl, instrumentPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "postgres: unable to update instrument row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: failed to get rows affected by update for instrument")
	}

	if !cached {
		instrumentUpdateCacheMut.Lock()
		instrumentUpdateCache[key] = cache
		instrumentUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q instrumentQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to update all for instrument")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to retrieve rows affected for instrument")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o InstrumentSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), instrumentPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"instrument\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, instrumentPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to update all in instrument slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to retrieve rows affected all in update all instrument")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Instrument) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("postgres: no instrument provided for upsert")
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

	nzDefaults := queries.NonZeroDefaultSet(instrumentColumnsWithDefault, o)

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

	instrumentUpsertCacheMut.RLock()
	cache, cached := instrumentUpsertCache[key]
	instrumentUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			instrumentAllColumns,
			instrumentColumnsWithDefault,
			instrumentColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			instrumentAllColumns,
			instrumentPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("postgres: unable to upsert instrument, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(instrumentPrimaryKeyColumns))
			copy(conflict, instrumentPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"instrument\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(instrumentType, instrumentMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(instrumentType, instrumentMapping, ret)
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
		return errors.Wrap(err, "postgres: unable to upsert instrument")
	}

	if !cached {
		instrumentUpsertCacheMut.Lock()
		instrumentUpsertCache[key] = cache
		instrumentUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Instrument record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Instrument) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("postgres: no Instrument provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), instrumentPrimaryKeyMapping)
	sql := "DELETE FROM \"instrument\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to delete from instrument")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: failed to get rows affected by delete for instrument")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q instrumentQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("postgres: no instrumentQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to delete all from instrument")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: failed to get rows affected by deleteall for instrument")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o InstrumentSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(instrumentBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), instrumentPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"instrument\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, instrumentPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "postgres: unable to delete all from instrument slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "postgres: failed to get rows affected by deleteall for instrument")
	}

	if len(instrumentAfterDeleteHooks) != 0 {
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
func (o *Instrument) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindInstrument(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *InstrumentSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := InstrumentSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), instrumentPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"instrument\".* FROM \"instrument\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, instrumentPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "postgres: unable to reload all in InstrumentSlice")
	}

	*o = slice

	return nil
}

// InstrumentExists checks if the Instrument row exists.
func InstrumentExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"instrument\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "postgres: unable to check if instrument exists")
	}

	return exists, nil
}
