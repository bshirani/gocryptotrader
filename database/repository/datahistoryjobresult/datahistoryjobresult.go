package datahistoryjobresult

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gocryptotrader/database"
	"gocryptotrader/database/models/postgres"
	"gocryptotrader/log"

	"github.com/gofrs/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// Setup returns a DBService
func Setup(db database.IDatabase) (*DBService, error) {
	if db == nil {
		return nil, nil
	}
	if !db.IsConnected() {
		return nil, nil
	}
	cfg := db.GetConfig()
	dbCon, err := db.GetSQL()
	if err != nil {
		return nil, err
	}
	return &DBService{
		sql:    dbCon,
		driver: cfg.Driver,
	}, nil
}

// Upsert inserts or updates jobs into the database
func (db *DBService) Upsert(jobs ...*DataHistoryJobResult) error {
	if len(jobs) == 0 {
		return nil
	}
	ctx := context.Background()

	tx, err := db.sql.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginTx %w", err)
	}
	defer func() {
		if err != nil {
			errRB := tx.Rollback()
			if errRB != nil {
				log.Errorf(log.DatabaseMgr, "Insert tx.Rollback %v", errRB)
			}
		}
	}()

	err = upsertPostgres(ctx, tx, jobs...)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetByJobID returns a job by its related JobID
func (db *DBService) GetByJobID(jobID string) ([]DataHistoryJobResult, error) {
	var err error
	var job []DataHistoryJobResult
	job, err = db.getByJobIDPostgres(jobID)
	if err != nil {
		return nil, err
	}
	return job, nil
}

// GetJobResultsBetween will return all jobs between two dates
func (db *DBService) GetJobResultsBetween(jobID string, startDate, endDate time.Time) ([]DataHistoryJobResult, error) {
	var err error
	var jobs []DataHistoryJobResult
	jobs, err = db.getJobResultsBetweenPostgres(jobID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func upsertPostgres(ctx context.Context, tx *sql.Tx, results ...*DataHistoryJobResult) error {
	var err error
	for i := range results {
		if results[i].ID == "" {
			var freshUUID uuid.UUID
			freshUUID, err = uuid.NewV4()
			if err != nil {
				return err
			}
			results[i].ID = freshUUID.String()
		}
		var tempEvent = postgres.Datahistoryjobresult{
			ID:                results[i].ID,
			JobID:             results[i].JobID,
			Result:            null.NewString(results[i].Result, results[i].Result != ""),
			Status:            float64(results[i].Status),
			IntervalStartTime: results[i].IntervalStartDate.UTC(),
			IntervalEndTime:   results[i].IntervalEndDate.UTC(),
			RunTime:           results[i].Date.UTC(),
		}
		err = tempEvent.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer())
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DBService) getByJobIDPostgres(jobID string) ([]DataHistoryJobResult, error) {
	query := postgres.Datahistoryjobresults(qm.Where("job_id = ?", jobID))
	results, err := query.All(context.Background(), db.sql)
	if err != nil {
		return nil, err
	}
	var resp []DataHistoryJobResult
	for i := range results {
		resp = append(resp, DataHistoryJobResult{
			ID:                results[i].ID,
			JobID:             results[i].JobID,
			IntervalStartDate: results[i].IntervalStartTime,
			IntervalEndDate:   results[i].IntervalEndTime,
			Status:            int64(results[i].Status),
			Result:            results[i].Result.String,
			Date:              results[i].RunTime,
		})
	}

	return resp, nil
}

func (db *DBService) getJobResultsBetweenPostgres(jobID string, startDate, endDate time.Time) ([]DataHistoryJobResult, error) {
	var jobs []DataHistoryJobResult
	query := postgres.Datahistoryjobresults(qm.Where("job_id = ? AND run_time BETWEEN ? AND  ? ", jobID, startDate, endDate))
	results, err := query.All(context.Background(), db.sql)
	if err != nil {
		return jobs, err
	}

	for i := range results {
		jobs = append(jobs, DataHistoryJobResult{
			ID:                results[i].ID,
			JobID:             results[i].JobID,
			IntervalStartDate: results[i].IntervalStartTime,
			IntervalEndDate:   results[i].IntervalEndTime,
			Status:            int64(results[i].Status),
			Result:            results[i].Result.String,
			Date:              results[i].RunTime,
		})
	}

	return jobs, nil
}
