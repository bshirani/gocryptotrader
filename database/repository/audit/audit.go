package audit

import (
	"context"
	"time"

	"gocryptotrader/database"
	modelPSQL "gocryptotrader/database/models/postgres"
	"gocryptotrader/log"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// Event inserts a new audit event to database
func Event(id, msgtype, message string) {
	if database.DB.SQL == nil {
		return
	}

	ctx := context.Background()
	ctx = boil.SkipTimestamps(ctx)

	tx, err := database.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		log.Errorf(log.Global, "Event transaction begin failed: %v", err)
		return
	}

	var tempEvent = modelPSQL.AuditEvent{
		Type:       msgtype,
		Identifier: id,
		Message:    message,
	}
	err = tempEvent.Insert(ctx, tx, boil.Blacklist("created_at"))

	if err != nil {
		log.Errorf(log.Global, "Event insert failed: %v", err)
		err = tx.Rollback()
		if err != nil {
			log.Errorf(log.Global, "Event Transaction rollback failed: %v", err)
		}
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Errorf(log.Global, "Event Transaction commit failed: %v", err)
		return
	}
}

// GetEvent () returns list of order events matching query
func GetEvent(startTime, endTime time.Time, order string, limit int) (interface{}, error) {
	if database.DB.SQL == nil {
		return nil, database.ErrDatabaseSupportDisabled
	}

	query := qm.Where("created_at BETWEEN ? AND ?", startTime, endTime)

	orderByQueryString := "id"
	if order == "desc" {
		orderByQueryString += " desc"
	}

	orderByQuery := qm.OrderBy(orderByQueryString)
	limitQuery := qm.Limit(limit)

	ctx := context.Background()

	return modelPSQL.AuditEvents(query, orderByQuery, limitQuery).All(ctx, database.DB.SQL)
}
