package script

import (
	"context"
	"time"

	"gocryptotrader/database"
	modelPSQL "gocryptotrader/database/models/postgres"
	"gocryptotrader/log"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// Event inserts a new script event into database with execution details (script name time status hash of script)
func Event(id, name, path string, data null.Bytes, executionType, status string, time time.Time) {
	if database.DB.SQL == nil {
		return
	}

	ctx := context.Background()
	ctx = boil.SkipTimestamps(ctx)
	tx, err := database.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		log.Errorf(log.DatabaseMgr, "Event transaction begin failed: %v", err)
		return
	}

	var tempEvent = modelPSQL.Script{
		ScriptID:   id,
		ScriptName: name,
		ScriptPath: path,
		ScriptData: data,
	}
	err = tempEvent.Upsert(ctx, tx, true, []string{"script_id"}, boil.Whitelist("last_executed_at"), boil.Infer())
	if err != nil {
		log.Errorf(log.DatabaseMgr, "Event insert failed: %v", err)
		err = tx.Rollback()
		if err != nil {
			log.Errorf(log.DatabaseMgr, "Event Transaction rollback failed: %v", err)
		}
		return
	}

	tempScriptExecution := &modelPSQL.ScriptExecution{
		ExecutionTime:   time.UTC(),
		ExecutionStatus: status,
		ExecutionType:   executionType,
	}

	err = tempEvent.AddScriptExecutions(ctx, tx, true, tempScriptExecution)
	if err != nil {
		log.Errorf(log.DatabaseMgr, "Event insert failed: %v", err)
		err = tx.Rollback()
		if err != nil {
			log.Errorf(log.DatabaseMgr, "Event Transaction rollback failed: %v", err)
		}
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Errorf(log.DatabaseMgr, "Event Transaction commit failed: %v", err)
	}
}
