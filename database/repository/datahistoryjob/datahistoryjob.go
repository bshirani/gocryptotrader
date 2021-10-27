package datahistoryjob

import (
	"context"
	"database/sql"
	"fmt"
	"gocryptotrader/database"
	"gocryptotrader/database/models/postgres"
	"gocryptotrader/database/repository/datahistoryjobresult"
	"gocryptotrader/log"
	"strings"
	"time"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// Setup returns a DBService
func Setup(db database.IDatabase) (*DBService, error) {
	if db == nil {
		return nil, database.ErrNilInstance
	}
	if !db.IsConnected() {
		return nil, database.ErrDatabaseNotConnected
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

func (db *DBService) ClearJobs() error {
	_, err := queries.Raw("update datahistoryjob set status = 1 where status = 0").Exec(db.sql)
	return err
}

// Upsert inserts or updates jobs into the database
func (db *DBService) Upsert(jobs ...*DataHistoryJob) error {
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

// GetByNickName returns a job by its nickname
func (db *DBService) GetByNickName(nickname string) (*DataHistoryJob, error) {
	// boil.DebugMode = true
	return db.getByNicknamePostgres(nickname)
}

// GetByNickName returns a job by its nickname
func (db *DBService) CountActive() (int64, error) {
	// boil.DebugMode = true
	query := postgres.Datahistoryjobs(qm.Where("status=0"))
	result, err := query.Count(context.Background(), db.sql)
	if err != nil {
		return -1, err
	}
	return result, err
}

// GetByID returns a job by its id
func (db *DBService) GetByID(id string) (*DataHistoryJob, error) {
	return db.getByIDPostgres(id)
}

// GetJobsBetween will return all jobs between two dates
func (db *DBService) GetJobsBetween(startDate, endDate time.Time) ([]DataHistoryJob, error) {
	return db.getJobsBetweenPostgres(startDate, endDate)
}

// GetAllIncompleteJobsAndResults returns all jobs that have the status "active"
func (db *DBService) GetAllIncompleteJobsAndResults() ([]DataHistoryJob, error) {
	// boil.DebugMode = true
	// defer func() { boil.DebugMode = false }()
	query := postgres.Datahistoryjobs(
		qm.Load(postgres.DatahistoryjobRels.JobDatahistoryjobresults),
		qm.Where("status = ?", 0))
	results, err := query.All(context.Background(), db.sql)
	if err != nil {
		return nil, err
	}

	var jobs []DataHistoryJob
	for i := range results {
		job, err := db.createPostgresDataHistoryJobResponse(results[i])
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, *job)
	}

	return jobs, nil
}

// GetJobAndAllResults returns a job and joins all job results
func (db *DBService) GetJobAndAllResults(nickname string) (*DataHistoryJob, error) {
	return db.getJobAndAllResultsPostgres(nickname)
}

// GetRelatedUpcomingJobs will return related jobs
func (db *DBService) GetRelatedUpcomingJobs(nickname string) ([]*DataHistoryJob, error) {
	return db.getRelatedUpcomingJobsPostgres(nickname)
}

// GetPrerequisiteJob will return the job that must complete before the
// referenced job
func (db *DBService) GetPrerequisiteJob(nickname string) (*DataHistoryJob, error) {
	return db.getPrerequisiteJobPostgres(nickname)
}

// SetRelationshipByID removes a relationship in the event of a changed
// relationship during upsertion
func (db *DBService) SetRelationshipByID(prerequisiteJobID, followingJobID string, status int64) error {
	ctx := context.Background()
	if strings.EqualFold(prerequisiteJobID, followingJobID) {
		return errCannotSetSamePrerequisite
	}
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

	err = setRelationshipByIDPostgres(ctx, tx, prerequisiteJobID, followingJobID, status)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// SetRelationshipByNickname removes a relationship in the event of a changed
// relationship during upsertion
func (db *DBService) SetRelationshipByNickname(prerequisiteNickname, followingNickname string, status int64) error {
	ctx := context.Background()
	if strings.EqualFold(prerequisiteNickname, followingNickname) {
		return errCannotSetSamePrerequisite
	}
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

	err = setRelationshipByNicknamePostgres(ctx, tx, prerequisiteNickname, followingNickname, status)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func upsertPostgres(ctx context.Context, tx *sql.Tx, jobs ...*DataHistoryJob) error {
	for i := range jobs {
		exch, err := postgres.Exchanges(
			qm.Where("name = ?", strings.ToLower(jobs[i].ExchangeName))).One(ctx, tx)
		if err != nil {
			return fmt.Errorf("could not retrieve exchange '%v', %w", jobs[i].ExchangeName, err)
		}
		var secondaryExch *postgres.Exchange
		if jobs[i].SecondarySourceExchangeName != "" {
			secondaryExch, err = postgres.Exchanges(
				qm.Where("name = ?", strings.ToLower(jobs[i].SecondarySourceExchangeName))).One(ctx, tx)
			if err != nil {
				return fmt.Errorf("could not retrieve secondary exchange '%v', %w", jobs[i].SecondarySourceExchangeName, err)
			}
		}

		var tempEvent = postgres.Datahistoryjob{
			ID:                       jobs[i].ID,
			Nickname:                 strings.ToLower(jobs[i].Nickname),
			ExchangeNameID:           exch.ID,
			Asset:                    strings.ToLower(jobs[i].Asset),
			Base:                     strings.ToUpper(jobs[i].Base),
			Quote:                    strings.ToUpper(jobs[i].Quote),
			StartTime:                jobs[i].StartDate.UTC(),
			EndTime:                  jobs[i].EndDate.UTC(),
			Interval:                 float64(jobs[i].Interval),
			DataType:                 float64(jobs[i].DataType),
			BatchCount:               float64(jobs[i].BatchSize),
			RequestSize:              float64(jobs[i].RequestSizeLimit),
			MaxRetries:               float64(jobs[i].MaxRetryAttempts),
			Status:                   float64(jobs[i].Status),
			Created:                  time.Now().UTC(),
			ConversionInterval:       null.Float64{Float64: float64(jobs[i].ConversionInterval), Valid: jobs[i].ConversionInterval > 0},
			OverwriteData:            null.Bool{Bool: jobs[i].OverwriteData, Valid: jobs[i].OverwriteData},
			DecimalPlaceComparison:   null.Int{Int: int(jobs[i].DecimalPlaceComparison), Valid: jobs[i].DecimalPlaceComparison > 0},
			ReplaceOnIssue:           null.Bool{Bool: jobs[i].ReplaceOnIssue, Valid: jobs[i].ReplaceOnIssue},
			IssueTolerancePercentage: null.Float64{Float64: jobs[i].IssueTolerancePercentage, Valid: jobs[i].IssueTolerancePercentage > 0},
		}
		if secondaryExch != nil {
			tempEvent.SecondaryExchangeID = null.String{String: secondaryExch.ID, Valid: true}
		}
		err = tempEvent.Upsert(ctx, tx, true, []string{"nickname"}, boil.Infer(), boil.Infer())
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DBService) getByNicknamePostgres(nickname string) (*DataHistoryJob, error) {
	query := postgres.Datahistoryjobs(qm.Where("nickname = ?", strings.ToLower(nickname)))
	result, err := query.One(context.Background(), db.sql)
	if err != nil {
		return nil, err
	}
	return db.createPostgresDataHistoryJobResponse(result)
}

func (db *DBService) getByIDPostgres(id string) (*DataHistoryJob, error) {
	query := postgres.Datahistoryjobs(qm.Where("id = ?", id))
	result, err := query.One(context.Background(), db.sql)
	if err != nil {
		return nil, err
	}

	return db.createPostgresDataHistoryJobResponse(result)
}

func (db *DBService) getJobsBetweenPostgres(startDate, endDate time.Time) ([]DataHistoryJob, error) {
	var jobs []DataHistoryJob
	query := postgres.Datahistoryjobs(qm.Where("created BETWEEN ? AND  ? ", startDate, endDate))
	results, err := query.All(context.Background(), db.sql)
	if err != nil {
		return jobs, err
	}

	for i := range results {
		job, err := db.createPostgresDataHistoryJobResponse(results[i])
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, *job)
	}

	return jobs, nil
}

func (db *DBService) getJobAndAllResultsPostgres(nickname string) (*DataHistoryJob, error) {
	query := postgres.Datahistoryjobs(
		qm.Load(postgres.DatahistoryjobRels.JobDatahistoryjobresults),
		qm.Where("nickname = ?", strings.ToLower(nickname)))
	result, err := query.One(context.Background(), db.sql)
	if err != nil {
		return nil, err
	}

	return db.createPostgresDataHistoryJobResponse(result)
}

func (db *DBService) getRelatedUpcomingJobsPostgres(nickname string) ([]*DataHistoryJob, error) {
	q := postgres.Datahistoryjobs(qm.Load(postgres.DatahistoryjobRels.JobDatahistoryjobs), qm.Where("nickname = ?", nickname))
	jobWithRelations, err := q.One(context.Background(), db.sql)
	if err != nil {
		return nil, err
	}
	var response []*DataHistoryJob
	for i := range jobWithRelations.R.JobDatahistoryjobs {
		job, err := db.getByIDPostgres(jobWithRelations.R.JobDatahistoryjobs[i].ID)
		if err != nil {
			return nil, err
		}
		response = append(response, job)
	}
	return response, nil
}

func setRelationshipByIDPostgres(ctx context.Context, tx *sql.Tx, prerequisiteJobID, followingJobID string, status int64) error {
	job, err := postgres.Datahistoryjobs(qm.Where("id = ?", followingJobID)).One(ctx, tx)
	if err != nil {
		return err
	}
	job.Status = float64(status)
	_, err = job.Update(ctx, tx, boil.Infer())
	if err != nil {
		return err
	}

	if prerequisiteJobID == "" {
		return job.SetPrerequisiteJobDatahistoryjobs(ctx, tx, false)
	}
	result, err := postgres.Datahistoryjobs(qm.Where("id = ?", prerequisiteJobID)).One(ctx, tx)
	if err != nil {
		return err
	}

	return job.SetPrerequisiteJobDatahistoryjobs(ctx, tx, false, result)
}

func (db *DBService) getPrerequisiteJobPostgres(nickname string) (*DataHistoryJob, error) {
	job, err := postgres.Datahistoryjobs(qm.Where("nickname = ?", nickname)).One(context.Background(), db.sql)
	if err != nil {
		return nil, err
	}
	result, err := job.PrerequisiteJobDatahistoryjobs().One(context.Background(), db.sql)
	if err != nil {
		return nil, err
	}

	return db.createPostgresDataHistoryJobResponse(result)
}

func setRelationshipByNicknamePostgres(ctx context.Context, tx *sql.Tx, prerequisiteJobNickname, followingJobNickname string, status int64) error {
	job, err := postgres.Datahistoryjobs(qm.Where("nickname = ?", followingJobNickname)).One(ctx, tx)
	if err != nil {
		return err
	}
	job.Status = float64(status)
	_, err = job.Update(ctx, tx, boil.Infer())
	if err != nil {
		return err
	}

	if prerequisiteJobNickname == "" {
		return job.SetPrerequisiteJobDatahistoryjobs(ctx, tx, false)
	}
	result, err := postgres.Datahistoryjobs(qm.Where("nickname = ?", prerequisiteJobNickname)).One(ctx, tx)
	if err != nil {
		return err
	}
	return job.SetPrerequisiteJobDatahistoryjobs(ctx, tx, false, result)
}

// helpers
func (db *DBService) createPostgresDataHistoryJobResponse(result *postgres.Datahistoryjob) (*DataHistoryJob, error) {
	var exchange *postgres.Exchange
	var err error
	if result.R != nil && result.R.ExchangeName != nil {
		exchange = result.R.ExchangeName
	} else {
		exchange, err = result.ExchangeName().One(context.Background(), db.sql)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve exchange '%v' %w", result.ExchangeNameID, err)
		}
	}

	var secondaryExchangeName string
	if result.SecondaryExchangeID.String != "" {
		var secondaryExchangeResult *postgres.Exchange
		secondaryExchangeResult, err = result.SecondaryExchange().One(context.Background(), db.sql)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve secondary exchange '%v' %w", result.SecondaryExchangeID, err)
		}
		if secondaryExchangeResult != nil {
			secondaryExchangeName = secondaryExchangeResult.Name
		}
	}

	prereqJob, err := result.PrerequisiteJobDatahistoryjobs().One(context.Background(), db.sql)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	var prereqNickname, prereqID string
	if prereqJob != nil {
		prereqID = prereqJob.ID
		prereqNickname = prereqJob.Nickname
	}

	var jobResults []*datahistoryjobresult.DataHistoryJobResult
	if result.R != nil {
		for i := range result.R.JobDatahistoryjobresults {
			jobResults = append(jobResults, &datahistoryjobresult.DataHistoryJobResult{
				ID:                result.R.JobDatahistoryjobresults[i].ID,
				JobID:             result.R.JobDatahistoryjobresults[i].JobID,
				IntervalStartDate: result.R.JobDatahistoryjobresults[i].IntervalStartTime,
				IntervalEndDate:   result.R.JobDatahistoryjobresults[i].IntervalEndTime,
				Status:            int64(result.R.JobDatahistoryjobresults[i].Status),
				Result:            result.R.JobDatahistoryjobresults[i].Result.String,
				Date:              result.R.JobDatahistoryjobresults[i].RunTime,
			})
		}
	}

	return &DataHistoryJob{
		ID:                          result.ID,
		Nickname:                    result.Nickname,
		ExchangeID:                  exchange.ID,
		ExchangeName:                exchange.Name,
		Asset:                       result.Asset,
		Base:                        result.Base,
		Quote:                       result.Quote,
		StartDate:                   result.StartTime,
		EndDate:                     result.EndTime,
		Interval:                    int64(result.Interval),
		RequestSizeLimit:            int64(result.RequestSize),
		DataType:                    int64(result.DataType),
		MaxRetryAttempts:            int64(result.MaxRetries),
		BatchSize:                   int64(result.BatchCount),
		Status:                      int64(result.Status),
		CreatedDate:                 result.Created,
		Results:                     jobResults,
		PrerequisiteJobID:           prereqID,
		PrerequisiteJobNickname:     prereqNickname,
		ConversionInterval:          int64(result.ConversionInterval.Float64),
		OverwriteData:               result.OverwriteData.Bool,
		DecimalPlaceComparison:      int64(result.DecimalPlaceComparison.Int),
		SecondarySourceExchangeName: secondaryExchangeName,
		IssueTolerancePercentage:    result.IssueTolerancePercentage.Float64,
		ReplaceOnIssue:              result.ReplaceOnIssue.Bool,
	}, nil
}
