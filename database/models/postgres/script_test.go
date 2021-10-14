// Code generated by SQLBoiler 3.5.0-gct (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package postgres

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/randomize"
	"github.com/volatiletech/sqlboiler/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testScripts(t *testing.T) {
	t.Parallel()

	query := Scripts()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testScriptsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Script{}
	if err = randomize.Struct(seed, o, scriptDBTypes, true, scriptColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
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

	count, err := Scripts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testScriptsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Script{}
	if err = randomize.Struct(seed, o, scriptDBTypes, true, scriptColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Scripts().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Scripts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testScriptsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Script{}
	if err = randomize.Struct(seed, o, scriptDBTypes, true, scriptColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ScriptSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Scripts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testScriptsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Script{}
	if err = randomize.Struct(seed, o, scriptDBTypes, true, scriptColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := ScriptExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if Script exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ScriptExists to return true, but got false.")
	}
}

func testScriptsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Script{}
	if err = randomize.Struct(seed, o, scriptDBTypes, true, scriptColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	scriptFound, err := FindScript(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if scriptFound == nil {
		t.Error("want a record, got nil")
	}
}

func testScriptsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Script{}
	if err = randomize.Struct(seed, o, scriptDBTypes, true, scriptColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Scripts().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testScriptsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Script{}
	if err = randomize.Struct(seed, o, scriptDBTypes, true, scriptColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Scripts().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testScriptsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	scriptOne := &Script{}
	scriptTwo := &Script{}
	if err = randomize.Struct(seed, scriptOne, scriptDBTypes, false, scriptColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
	}
	if err = randomize.Struct(seed, scriptTwo, scriptDBTypes, false, scriptColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = scriptOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = scriptTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Scripts().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testScriptsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	scriptOne := &Script{}
	scriptTwo := &Script{}
	if err = randomize.Struct(seed, scriptOne, scriptDBTypes, false, scriptColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
	}
	if err = randomize.Struct(seed, scriptTwo, scriptDBTypes, false, scriptColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = scriptOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = scriptTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Scripts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func scriptBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *Script) error {
	*o = Script{}
	return nil
}

func scriptAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *Script) error {
	*o = Script{}
	return nil
}

func scriptAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *Script) error {
	*o = Script{}
	return nil
}

func scriptBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Script) error {
	*o = Script{}
	return nil
}

func scriptAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Script) error {
	*o = Script{}
	return nil
}

func scriptBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Script) error {
	*o = Script{}
	return nil
}

func scriptAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Script) error {
	*o = Script{}
	return nil
}

func scriptBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Script) error {
	*o = Script{}
	return nil
}

func scriptAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Script) error {
	*o = Script{}
	return nil
}

func testScriptsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &Script{}
	o := &Script{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, scriptDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Script object: %s", err)
	}

	AddScriptHook(boil.BeforeInsertHook, scriptBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	scriptBeforeInsertHooks = []ScriptHook{}

	AddScriptHook(boil.AfterInsertHook, scriptAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	scriptAfterInsertHooks = []ScriptHook{}

	AddScriptHook(boil.AfterSelectHook, scriptAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	scriptAfterSelectHooks = []ScriptHook{}

	AddScriptHook(boil.BeforeUpdateHook, scriptBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	scriptBeforeUpdateHooks = []ScriptHook{}

	AddScriptHook(boil.AfterUpdateHook, scriptAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	scriptAfterUpdateHooks = []ScriptHook{}

	AddScriptHook(boil.BeforeDeleteHook, scriptBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	scriptBeforeDeleteHooks = []ScriptHook{}

	AddScriptHook(boil.AfterDeleteHook, scriptAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	scriptAfterDeleteHooks = []ScriptHook{}

	AddScriptHook(boil.BeforeUpsertHook, scriptBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	scriptBeforeUpsertHooks = []ScriptHook{}

	AddScriptHook(boil.AfterUpsertHook, scriptAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	scriptAfterUpsertHooks = []ScriptHook{}
}

func testScriptsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Script{}
	if err = randomize.Struct(seed, o, scriptDBTypes, true, scriptColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Scripts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testScriptsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Script{}
	if err = randomize.Struct(seed, o, scriptDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(scriptColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Scripts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testScriptToManyScriptExecutions(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Script
	var b, c ScriptExecution

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, scriptDBTypes, true, scriptColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, scriptExecutionDBTypes, false, scriptExecutionColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, scriptExecutionDBTypes, false, scriptExecutionColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	queries.Assign(&b.ScriptID, a.ID)
	queries.Assign(&c.ScriptID, a.ID)
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.ScriptExecutions().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if queries.Equal(v.ScriptID, b.ScriptID) {
			bFound = true
		}
		if queries.Equal(v.ScriptID, c.ScriptID) {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := ScriptSlice{&a}
	if err = a.L.LoadScriptExecutions(ctx, tx, false, (*[]*Script)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ScriptExecutions); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.ScriptExecutions = nil
	if err = a.L.LoadScriptExecutions(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ScriptExecutions); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testScriptToManyAddOpScriptExecutions(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Script
	var b, c, d, e ScriptExecution

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, scriptDBTypes, false, strmangle.SetComplement(scriptPrimaryKeyColumns, scriptColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ScriptExecution{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, scriptExecutionDBTypes, false, strmangle.SetComplement(scriptExecutionPrimaryKeyColumns, scriptExecutionColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*ScriptExecution{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddScriptExecutions(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if !queries.Equal(a.ID, first.ScriptID) {
			t.Error("foreign key was wrong value", a.ID, first.ScriptID)
		}
		if !queries.Equal(a.ID, second.ScriptID) {
			t.Error("foreign key was wrong value", a.ID, second.ScriptID)
		}

		if first.R.Script != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Script != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.ScriptExecutions[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.ScriptExecutions[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.ScriptExecutions().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testScriptToManySetOpScriptExecutions(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Script
	var b, c, d, e ScriptExecution

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, scriptDBTypes, false, strmangle.SetComplement(scriptPrimaryKeyColumns, scriptColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ScriptExecution{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, scriptExecutionDBTypes, false, strmangle.SetComplement(scriptExecutionPrimaryKeyColumns, scriptExecutionColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err = a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.SetScriptExecutions(ctx, tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ScriptExecutions().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetScriptExecutions(ctx, tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.ScriptExecutions().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if !queries.IsValuerNil(b.ScriptID) {
		t.Error("want b's foreign key value to be nil")
	}
	if !queries.IsValuerNil(c.ScriptID) {
		t.Error("want c's foreign key value to be nil")
	}
	if !queries.Equal(a.ID, d.ScriptID) {
		t.Error("foreign key was wrong value", a.ID, d.ScriptID)
	}
	if !queries.Equal(a.ID, e.ScriptID) {
		t.Error("foreign key was wrong value", a.ID, e.ScriptID)
	}

	if b.R.Script != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Script != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Script != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Script != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.ScriptExecutions[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.ScriptExecutions[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testScriptToManyRemoveOpScriptExecutions(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Script
	var b, c, d, e ScriptExecution

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, scriptDBTypes, false, strmangle.SetComplement(scriptPrimaryKeyColumns, scriptColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ScriptExecution{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, scriptExecutionDBTypes, false, strmangle.SetComplement(scriptExecutionPrimaryKeyColumns, scriptExecutionColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.AddScriptExecutions(ctx, tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ScriptExecutions().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveScriptExecutions(ctx, tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.ScriptExecutions().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if !queries.IsValuerNil(b.ScriptID) {
		t.Error("want b's foreign key value to be nil")
	}
	if !queries.IsValuerNil(c.ScriptID) {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.Script != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Script != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Script != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.Script != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.ScriptExecutions) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.ScriptExecutions[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.ScriptExecutions[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testScriptsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Script{}
	if err = randomize.Struct(seed, o, scriptDBTypes, true, scriptColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
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

func testScriptsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Script{}
	if err = randomize.Struct(seed, o, scriptDBTypes, true, scriptColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ScriptSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testScriptsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Script{}
	if err = randomize.Struct(seed, o, scriptDBTypes, true, scriptColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Scripts().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	scriptDBTypes = map[string]string{`ID`: `uuid`, `ScriptID`: `text`, `ScriptName`: `character varying`, `ScriptPath`: `character varying`, `ScriptData`: `bytea`, `LastExecutedAt`: `timestamp without time zone`, `CreatedAt`: `timestamp without time zone`}
	_             = bytes.MinRead
)

func testScriptsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(scriptPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(scriptAllColumns) == len(scriptPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Script{}
	if err = randomize.Struct(seed, o, scriptDBTypes, true, scriptColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Scripts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, scriptDBTypes, true, scriptPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testScriptsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(scriptAllColumns) == len(scriptPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Script{}
	if err = randomize.Struct(seed, o, scriptDBTypes, true, scriptColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Scripts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, scriptDBTypes, true, scriptPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(scriptAllColumns, scriptPrimaryKeyColumns) {
		fields = scriptAllColumns
	} else {
		fields = strmangle.SetComplement(
			scriptAllColumns,
			scriptPrimaryKeyColumns,
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

	slice := ScriptSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testScriptsUpsert(t *testing.T) {
	t.Parallel()

	if len(scriptAllColumns) == len(scriptPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Script{}
	if err = randomize.Struct(seed, &o, scriptDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Script: %s", err)
	}

	count, err := Scripts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, scriptDBTypes, false, scriptPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Script struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Script: %s", err)
	}

	count, err = Scripts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
