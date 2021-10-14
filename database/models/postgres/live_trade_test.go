// Code generated by SQLBoiler 4.7.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package postgres

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/randomize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testLiveTrades(t *testing.T) {
	t.Parallel()

	query := LiveTrades()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testLiveTradesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LiveTrade{}
	if err = randomize.Struct(seed, o, liveTradeDBTypes, true, liveTradeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := LiveTrades().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testLiveTradesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LiveTrade{}
	if err = randomize.Struct(seed, o, liveTradeDBTypes, true, liveTradeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := LiveTrades().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := LiveTrades().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testLiveTradesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LiveTrade{}
	if err = randomize.Struct(seed, o, liveTradeDBTypes, true, liveTradeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := LiveTradeSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := LiveTrades().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testLiveTradesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LiveTrade{}
	if err = randomize.Struct(seed, o, liveTradeDBTypes, true, liveTradeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := LiveTradeExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if LiveTrade exists: %s", err)
	}
	if !e {
		t.Errorf("Expected LiveTradeExists to return true, but got false.")
	}
}

func testLiveTradesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LiveTrade{}
	if err = randomize.Struct(seed, o, liveTradeDBTypes, true, liveTradeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	liveTradeFound, err := FindLiveTrade(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if liveTradeFound == nil {
		t.Error("want a record, got nil")
	}
}

func testLiveTradesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LiveTrade{}
	if err = randomize.Struct(seed, o, liveTradeDBTypes, true, liveTradeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = LiveTrades().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testLiveTradesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LiveTrade{}
	if err = randomize.Struct(seed, o, liveTradeDBTypes, true, liveTradeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := LiveTrades().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testLiveTradesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	liveTradeOne := &LiveTrade{}
	liveTradeTwo := &LiveTrade{}
	if err = randomize.Struct(seed, liveTradeOne, liveTradeDBTypes, false, liveTradeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}
	if err = randomize.Struct(seed, liveTradeTwo, liveTradeDBTypes, false, liveTradeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = liveTradeOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = liveTradeTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := LiveTrades().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testLiveTradesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	liveTradeOne := &LiveTrade{}
	liveTradeTwo := &LiveTrade{}
	if err = randomize.Struct(seed, liveTradeOne, liveTradeDBTypes, false, liveTradeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}
	if err = randomize.Struct(seed, liveTradeTwo, liveTradeDBTypes, false, liveTradeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = liveTradeOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = liveTradeTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := LiveTrades().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func liveTradeBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *LiveTrade) error {
	*o = LiveTrade{}
	return nil
}

func liveTradeAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *LiveTrade) error {
	*o = LiveTrade{}
	return nil
}

func liveTradeAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *LiveTrade) error {
	*o = LiveTrade{}
	return nil
}

func liveTradeBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *LiveTrade) error {
	*o = LiveTrade{}
	return nil
}

func liveTradeAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *LiveTrade) error {
	*o = LiveTrade{}
	return nil
}

func liveTradeBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *LiveTrade) error {
	*o = LiveTrade{}
	return nil
}

func liveTradeAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *LiveTrade) error {
	*o = LiveTrade{}
	return nil
}

func liveTradeBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *LiveTrade) error {
	*o = LiveTrade{}
	return nil
}

func liveTradeAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *LiveTrade) error {
	*o = LiveTrade{}
	return nil
}

func testLiveTradesHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &LiveTrade{}
	o := &LiveTrade{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, liveTradeDBTypes, false); err != nil {
		t.Errorf("Unable to randomize LiveTrade object: %s", err)
	}

	AddLiveTradeHook(boil.BeforeInsertHook, liveTradeBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	liveTradeBeforeInsertHooks = []LiveTradeHook{}

	AddLiveTradeHook(boil.AfterInsertHook, liveTradeAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	liveTradeAfterInsertHooks = []LiveTradeHook{}

	AddLiveTradeHook(boil.AfterSelectHook, liveTradeAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	liveTradeAfterSelectHooks = []LiveTradeHook{}

	AddLiveTradeHook(boil.BeforeUpdateHook, liveTradeBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	liveTradeBeforeUpdateHooks = []LiveTradeHook{}

	AddLiveTradeHook(boil.AfterUpdateHook, liveTradeAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	liveTradeAfterUpdateHooks = []LiveTradeHook{}

	AddLiveTradeHook(boil.BeforeDeleteHook, liveTradeBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	liveTradeBeforeDeleteHooks = []LiveTradeHook{}

	AddLiveTradeHook(boil.AfterDeleteHook, liveTradeAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	liveTradeAfterDeleteHooks = []LiveTradeHook{}

	AddLiveTradeHook(boil.BeforeUpsertHook, liveTradeBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	liveTradeBeforeUpsertHooks = []LiveTradeHook{}

	AddLiveTradeHook(boil.AfterUpsertHook, liveTradeAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	liveTradeAfterUpsertHooks = []LiveTradeHook{}
}

func testLiveTradesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LiveTrade{}
	if err = randomize.Struct(seed, o, liveTradeDBTypes, true, liveTradeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := LiveTrades().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testLiveTradesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LiveTrade{}
	if err = randomize.Struct(seed, o, liveTradeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(liveTradeColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := LiveTrades().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testLiveTradesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LiveTrade{}
	if err = randomize.Struct(seed, o, liveTradeDBTypes, true, liveTradeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testLiveTradesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LiveTrade{}
	if err = randomize.Struct(seed, o, liveTradeDBTypes, true, liveTradeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := LiveTradeSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testLiveTradesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LiveTrade{}
	if err = randomize.Struct(seed, o, liveTradeDBTypes, true, liveTradeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := LiveTrades().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	liveTradeDBTypes = map[string]string{`ID`: `integer`, `Side`: `text`, `EntryOrderID`: `text`, `EntryPrice`: `real`, `EntryTime`: `timestamp with time zone`, `StopLossPrice`: `real`, `StrategyID`: `text`, `Status`: `text`, `Pair`: `text`, `ExitTime`: `timestamp with time zone`, `TakeProfitPrice`: `real`, `ProfitLossPoints`: `real`, `ExitPrice`: `real`, `CreatedAt`: `timestamp without time zone`, `UpdatedAt`: `timestamp without time zone`}
	_                = bytes.MinRead
)

func testLiveTradesUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(liveTradePrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(liveTradeAllColumns) == len(liveTradePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &LiveTrade{}
	if err = randomize.Struct(seed, o, liveTradeDBTypes, true, liveTradeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := LiveTrades().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, liveTradeDBTypes, true, liveTradePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testLiveTradesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(liveTradeAllColumns) == len(liveTradePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &LiveTrade{}
	if err = randomize.Struct(seed, o, liveTradeDBTypes, true, liveTradeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := LiveTrades().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, liveTradeDBTypes, true, liveTradePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(liveTradeAllColumns, liveTradePrimaryKeyColumns) {
		fields = liveTradeAllColumns
	} else {
		fields = strmangle.SetComplement(
			liveTradeAllColumns,
			liveTradePrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := LiveTradeSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testLiveTradesUpsert(t *testing.T) {
	t.Parallel()

	if len(liveTradeAllColumns) == len(liveTradePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := LiveTrade{}
	if err = randomize.Struct(seed, &o, liveTradeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert LiveTrade: %s", err)
	}

	count, err := LiveTrades().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, liveTradeDBTypes, false, liveTradePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize LiveTrade struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert LiveTrade: %s", err)
	}

	count, err = LiveTrades().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
