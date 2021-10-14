package trade

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"gocryptotrader/database"
	"gocryptotrader/database/models/postgres"
	"gocryptotrader/database/repository/exchange"
	"gocryptotrader/exchange/kline"
	"gocryptotrader/log"

	"github.com/gofrs/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// Insert saves trade data to the database
func Insert(trades ...Data) error {
	for i := range trades {
		if trades[i].ExchangeNameID == "" && trades[i].Exchange != "" {
			exchangeUUID, err := exchange.UUIDByName(trades[i].Exchange)
			if err != nil {
				return err
			}
			trades[i].ExchangeNameID = exchangeUUID.String()
		} else if trades[i].ExchangeNameID == "" && trades[i].Exchange == "" {
			return errors.New("exchange name/uuid not set, cannot insert")
		}
	}

	ctx := context.Background()
	ctx = boil.SkipTimestamps(ctx)

	tx, err := database.DB.SQL.BeginTx(ctx, nil)
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

	err = insertPostgres(ctx, tx, trades...)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// VerifyTradeInIntervals will query for ONE trade within each kline interval and verify if data exists
// if it does, it will set the range holder property "HasData" to true
func VerifyTradeInIntervals(exchangeName, assetType, base, quote string, irh *kline.IntervalRangeHolder) error {
	ctx := context.Background()
	ctx = boil.SkipTimestamps(ctx)

	tx, err := database.DB.SQL.BeginTx(ctx, nil)
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

	err = verifyTradeInIntervalsPostgres(ctx, tx, exchangeName, assetType, base, quote, irh)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func verifyTradeInIntervalsPostgres(ctx context.Context, tx *sql.Tx, exchangeName, assetType, base, quote string, irh *kline.IntervalRangeHolder) error {
	exch, err := postgres.Exchanges(qm.Where("name = ?", exchangeName)).One(ctx, tx)
	if err != nil {
		return err
	}
	for i := range irh.Ranges {
		for j := range irh.Ranges[i].Intervals {
			result, err := postgres.Trades(qm.Where("exchange_name_id = ? AND asset = ? AND base = ? AND quote = ? AND timestamp between ? AND ?",
				exch.ID,
				assetType,
				base,
				quote,
				irh.Ranges[i].Intervals[j].Start.Time.UTC(),
				irh.Ranges[i].Intervals[j].End.Time.UTC())).One(ctx, tx)
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			if result != nil {
				irh.Ranges[i].Intervals[j].HasData = true
			}
		}
	}

	return nil
}

func insertPostgres(ctx context.Context, tx *sql.Tx, trades ...Data) error {
	var err error
	for i := range trades {
		if trades[i].ID == "" {
			var freshUUID uuid.UUID
			freshUUID, err = uuid.NewV4()
			if err != nil {
				return err
			}
			trades[i].ID = freshUUID.String()
		}
		var tempEvent = postgres.Trade{
			ExchangeNameID: trades[i].ExchangeNameID,
			Base:           strings.ToUpper(trades[i].Base),
			Quote:          strings.ToUpper(trades[i].Quote),
			Asset:          strings.ToLower(trades[i].AssetType),
			Price:          trades[i].Price,
			Amount:         trades[i].Amount,
			Timestamp:      trades[i].Timestamp.UTC(),
			ID:             trades[i].ID,
		}
		if trades[i].Side != "" {
			tempEvent.Side.SetValid(strings.ToUpper(trades[i].Side))
		}
		if trades[i].TID != "" {
			tempEvent.Tid.SetValid(trades[i].TID)
		}

		err = tempEvent.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer())
		if err != nil {
			return err
		}
	}

	return nil
}

// GetByUUID returns a trade by its unique ID
func GetByUUID(uuid string) (td Data, err error) {
	td, err = getByUUIDPostgres(uuid)
	if err != nil {
		return td, fmt.Errorf("trade.Get getByUUIDPostgres %w", err)
	}

	return td, nil
}

func getByUUIDPostgres(uuid string) (td Data, err error) {
	query := postgres.Trades(qm.Where("id = ?", uuid))
	var result *postgres.Trade
	result, err = query.One(context.Background(), database.DB.SQL)
	if err != nil {
		return td, err
	}

	td = Data{
		ID:        result.ID,
		Timestamp: result.Timestamp.UTC(),
		Exchange:  result.ExchangeNameID,
		Base:      strings.ToUpper(result.Base),
		Quote:     strings.ToUpper(result.Quote),
		AssetType: strings.ToLower(result.Asset),
		Price:     result.Price,
		Amount:    result.Amount,
	}
	if result.Side.Valid {
		td.Side = result.Side.String
	}
	return td, nil
}

// GetInRange returns all trades by an exchange in a date range
func GetInRange(exchangeName, assetType, base, quote string, startDate, endDate time.Time) (td []Data, err error) {
	td, err = getInRangePostgres(exchangeName, assetType, base, quote, startDate, endDate)
	if err != nil {
		return td, fmt.Errorf("trade.GetByExchangeInRange getInRangePostgres %w", err)
	}

	return td, nil
}

func getInRangePostgres(exchangeName, assetType, base, quote string, startDate, endDate time.Time) (td []Data, err error) {
	var exchangeUUID uuid.UUID
	exchangeUUID, err = exchange.UUIDByName(exchangeName)
	if err != nil {
		return nil, err
	}
	wheres := map[string]interface{}{
		"exchange_name_id": exchangeUUID,
		"asset":            strings.ToLower(assetType),
		"base":             strings.ToUpper(base),
		"quote":            strings.ToUpper(quote),
	}

	q := generateQuery(wheres, startDate, endDate)
	query := postgres.Trades(q...)
	var result []*postgres.Trade
	result, err = query.All(context.Background(), database.DB.SQL)
	if err != nil {
		return td, err
	}
	for i := range result {
		t := Data{
			ID:        result[i].ID,
			Timestamp: result[i].Timestamp,
			Exchange:  strings.ToLower(exchangeName),
			Base:      strings.ToUpper(result[i].Base),
			Quote:     strings.ToUpper(result[i].Quote),
			AssetType: strings.ToLower(result[i].Asset),
			Price:     result[i].Price,
			Amount:    result[i].Amount,
		}
		if result[i].Side.Valid {
			t.Side = result[i].Side.String
		}
		td = append(td, t)
	}
	return td, nil
}

// DeleteTrades will remove trades from the database using trade.Data
func DeleteTrades(trades ...Data) error {
	ctx := context.Background()
	ctx = boil.SkipTimestamps(ctx)

	tx, err := database.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginTx %w", err)
	}
	defer func() {
		if err != nil {
			errRB := tx.Rollback()
			if errRB != nil {
				log.Errorf(log.DatabaseMgr, "DeleteTrades tx.Rollback %v", errRB)
			}
		}
	}()
	err = deleteTradesPostgres(context.Background(), tx, trades...)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func deleteTradesPostgres(ctx context.Context, tx *sql.Tx, trades ...Data) error {
	var tradeIDs []interface{}
	for i := range trades {
		tradeIDs = append(tradeIDs, trades[i].ID)
	}
	query := postgres.Trades(qm.WhereIn(`id in ?`, tradeIDs...))
	_, err := query.DeleteAll(ctx, tx)
	return err
}

func generateQuery(clauses map[string]interface{}, start, end time.Time) []qm.QueryMod {
	query := []qm.QueryMod{
		qm.OrderBy("timestamp"),
	}
	query = append(query, qm.Where("timestamp BETWEEN ? AND ?", start.UTC(), end.UTC()))
	for k, v := range clauses {
		query = append(query, qm.Where(k+` = ?`, v))
	}

	return query
}
