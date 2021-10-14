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

func testAuditEvents(t *testing.T) {
	t.Parallel()

	query := AuditEvents()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testAuditEventsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuditEvent{}
	if err = randomize.Struct(seed, o, auditEventDBTypes, true, auditEventColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
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

	count, err := AuditEvents().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAuditEventsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuditEvent{}
	if err = randomize.Struct(seed, o, auditEventDBTypes, true, auditEventColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := AuditEvents().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := AuditEvents().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAuditEventsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuditEvent{}
	if err = randomize.Struct(seed, o, auditEventDBTypes, true, auditEventColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := AuditEventSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := AuditEvents().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAuditEventsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuditEvent{}
	if err = randomize.Struct(seed, o, auditEventDBTypes, true, auditEventColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := AuditEventExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if AuditEvent exists: %s", err)
	}
	if !e {
		t.Errorf("Expected AuditEventExists to return true, but got false.")
	}
}

func testAuditEventsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuditEvent{}
	if err = randomize.Struct(seed, o, auditEventDBTypes, true, auditEventColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	auditEventFound, err := FindAuditEvent(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if auditEventFound == nil {
		t.Error("want a record, got nil")
	}
}

func testAuditEventsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuditEvent{}
	if err = randomize.Struct(seed, o, auditEventDBTypes, true, auditEventColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = AuditEvents().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testAuditEventsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuditEvent{}
	if err = randomize.Struct(seed, o, auditEventDBTypes, true, auditEventColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := AuditEvents().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testAuditEventsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	auditEventOne := &AuditEvent{}
	auditEventTwo := &AuditEvent{}
	if err = randomize.Struct(seed, auditEventOne, auditEventDBTypes, false, auditEventColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
	}
	if err = randomize.Struct(seed, auditEventTwo, auditEventDBTypes, false, auditEventColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = auditEventOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = auditEventTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := AuditEvents().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testAuditEventsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	auditEventOne := &AuditEvent{}
	auditEventTwo := &AuditEvent{}
	if err = randomize.Struct(seed, auditEventOne, auditEventDBTypes, false, auditEventColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
	}
	if err = randomize.Struct(seed, auditEventTwo, auditEventDBTypes, false, auditEventColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = auditEventOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = auditEventTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := AuditEvents().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func auditEventBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *AuditEvent) error {
	*o = AuditEvent{}
	return nil
}

func auditEventAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *AuditEvent) error {
	*o = AuditEvent{}
	return nil
}

func auditEventAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *AuditEvent) error {
	*o = AuditEvent{}
	return nil
}

func auditEventBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *AuditEvent) error {
	*o = AuditEvent{}
	return nil
}

func auditEventAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *AuditEvent) error {
	*o = AuditEvent{}
	return nil
}

func auditEventBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *AuditEvent) error {
	*o = AuditEvent{}
	return nil
}

func auditEventAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *AuditEvent) error {
	*o = AuditEvent{}
	return nil
}

func auditEventBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *AuditEvent) error {
	*o = AuditEvent{}
	return nil
}

func auditEventAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *AuditEvent) error {
	*o = AuditEvent{}
	return nil
}

func testAuditEventsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &AuditEvent{}
	o := &AuditEvent{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, auditEventDBTypes, false); err != nil {
		t.Errorf("Unable to randomize AuditEvent object: %s", err)
	}

	AddAuditEventHook(boil.BeforeInsertHook, auditEventBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	auditEventBeforeInsertHooks = []AuditEventHook{}

	AddAuditEventHook(boil.AfterInsertHook, auditEventAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	auditEventAfterInsertHooks = []AuditEventHook{}

	AddAuditEventHook(boil.AfterSelectHook, auditEventAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	auditEventAfterSelectHooks = []AuditEventHook{}

	AddAuditEventHook(boil.BeforeUpdateHook, auditEventBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	auditEventBeforeUpdateHooks = []AuditEventHook{}

	AddAuditEventHook(boil.AfterUpdateHook, auditEventAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	auditEventAfterUpdateHooks = []AuditEventHook{}

	AddAuditEventHook(boil.BeforeDeleteHook, auditEventBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	auditEventBeforeDeleteHooks = []AuditEventHook{}

	AddAuditEventHook(boil.AfterDeleteHook, auditEventAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	auditEventAfterDeleteHooks = []AuditEventHook{}

	AddAuditEventHook(boil.BeforeUpsertHook, auditEventBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	auditEventBeforeUpsertHooks = []AuditEventHook{}

	AddAuditEventHook(boil.AfterUpsertHook, auditEventAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	auditEventAfterUpsertHooks = []AuditEventHook{}
}

func testAuditEventsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuditEvent{}
	if err = randomize.Struct(seed, o, auditEventDBTypes, true, auditEventColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := AuditEvents().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testAuditEventsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuditEvent{}
	if err = randomize.Struct(seed, o, auditEventDBTypes, true); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(auditEventColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := AuditEvents().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testAuditEventsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuditEvent{}
	if err = randomize.Struct(seed, o, auditEventDBTypes, true, auditEventColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
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

func testAuditEventsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuditEvent{}
	if err = randomize.Struct(seed, o, auditEventDBTypes, true, auditEventColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := AuditEventSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testAuditEventsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &AuditEvent{}
	if err = randomize.Struct(seed, o, auditEventDBTypes, true, auditEventColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := AuditEvents().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	auditEventDBTypes = map[string]string{`ID`: `bigint`, `Type`: `character varying`, `Identifier`: `character varying`, `Message`: `text`, `CreatedAt`: `timestamp without time zone`}
	_                 = bytes.MinRead
)

func testAuditEventsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(auditEventPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(auditEventAllColumns) == len(auditEventPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &AuditEvent{}
	if err = randomize.Struct(seed, o, auditEventDBTypes, true, auditEventColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := AuditEvents().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, auditEventDBTypes, true, auditEventPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testAuditEventsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(auditEventAllColumns) == len(auditEventPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &AuditEvent{}
	if err = randomize.Struct(seed, o, auditEventDBTypes, true, auditEventColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := AuditEvents().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, auditEventDBTypes, true, auditEventPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(auditEventAllColumns, auditEventPrimaryKeyColumns) {
		fields = auditEventAllColumns
	} else {
		fields = strmangle.SetComplement(
			auditEventAllColumns,
			auditEventPrimaryKeyColumns,
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

	slice := AuditEventSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testAuditEventsUpsert(t *testing.T) {
	t.Parallel()

	if len(auditEventAllColumns) == len(auditEventPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := AuditEvent{}
	if err = randomize.Struct(seed, &o, auditEventDBTypes, true); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert AuditEvent: %s", err)
	}

	count, err := AuditEvents().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, auditEventDBTypes, false, auditEventPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize AuditEvent struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert AuditEvent: %s", err)
	}

	count, err = AuditEvents().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
