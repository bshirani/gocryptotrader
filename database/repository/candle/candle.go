package candle

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"gocryptotrader/database"
	modelPSQL "gocryptotrader/database/models/postgres"
	"gocryptotrader/database/repository/exchange"
	"gocryptotrader/log"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func CountExchange(exchangeName string, interval int64, asset string) (count int64, err error) {
	// boil.DebugMode = true
	// defer func() { boil.DebugMode = false }()
	if exchangeName == "" || asset == "" || interval <= 0 {
		return count, errInvalidInput
	}

	queries := []qm.QueryMod{
		qm.Where("interval = ?", interval),
		qm.Where("asset = ?", strings.ToLower(asset)),
	}

	exchangeUUID, errS := exchange.UUIDByName(exchangeName)
	if errS != nil {
		return count, errS
	}
	queries = append(queries, qm.Where("exchange_name_id = ?", exchangeUUID.String()))
	return modelPSQL.Candles(queries...).Count(context.Background(), database.DB.SQL)
}

func CountTo(exchangeName, base, quote string, interval int64, asset string, to time.Time) (count int64, err error) {
	// boil.DebugMode = true
	// defer func() { boil.DebugMode = false }()
	if exchangeName == "" || base == "" || quote == "" || asset == "" || interval <= 0 {
		return count, errInvalidInput
	}

	queries := []qm.QueryMod{
		qm.Where("base = ?", strings.ToUpper(base)),
		qm.Where("quote = ?", strings.ToUpper(quote)),
		qm.Where("interval = ?", interval),
		qm.Where("asset = ?", strings.ToLower(asset)),
		qm.Where("timestamp <= ?", to),
	}

	exchangeUUID, errS := exchange.UUIDByName(exchangeName)
	if errS != nil {
		return count, errS
	}
	queries = append(queries, qm.Where("exchange_name_id = ?", exchangeUUID.String()))
	count, err = modelPSQL.Candles(queries...).Count(context.Background(), database.DB.SQL)
	// fmt.Println("candlecount returning", count)
	return count, err
}

func CountFrom(exchangeName, base, quote string, interval int64, asset string, from time.Time) (count int64, err error) {
	// boil.DebugMode = true
	// defer func() { boil.DebugMode = false }()
	if exchangeName == "" || base == "" || quote == "" || asset == "" || interval <= 0 {
		return count, errInvalidInput
	}

	queries := []qm.QueryMod{
		qm.Where("base = ?", strings.ToUpper(base)),
		qm.Where("quote = ?", strings.ToUpper(quote)),
		qm.Where("interval = ?", interval),
		qm.Where("asset = ?", strings.ToLower(asset)),
		qm.Where("timestamp >= ?", from),
	}

	exchangeUUID, errS := exchange.UUIDByName(exchangeName)
	if errS != nil {
		return count, errS
	}
	queries = append(queries, qm.Where("exchange_name_id = ?", exchangeUUID.String()))
	count, err = modelPSQL.Candles(queries...).Count(context.Background(), database.DB.SQL)
	// fmt.Println("candlecount returning", count)
	return count, err
}

func Count(exchangeName, base, quote string, interval int64, asset string) (count int64, err error) {
	// boil.DebugMode = true
	// defer func() { boil.DebugMode = false }()
	if exchangeName == "" || base == "" || quote == "" || asset == "" || interval <= 0 {
		return count, errInvalidInput
	}

	queries := []qm.QueryMod{
		qm.Where("base = ?", strings.ToUpper(base)),
		qm.Where("quote = ?", strings.ToUpper(quote)),
		qm.Where("interval = ?", interval),
		qm.Where("asset = ?", strings.ToLower(asset)),
	}

	exchangeUUID, errS := exchange.UUIDByName(exchangeName)
	if errS != nil {
		return count, errS
	}
	queries = append(queries, qm.Where("exchange_name_id = ?", exchangeUUID.String()))
	count, err = modelPSQL.Candles(queries...).Count(context.Background(), database.DB.SQL)
	// fmt.Println("candlecount returning", count)
	return count, err
}

func Counts() (pairCounts []PairCandleCountResponse, err error) {
	queries.Raw(`
		select t.date, base, quote, exchange_name_id as exchange_id, asset as asset_type, count(c.*) as "count"
		FROM (
			select generate_series(
			(date '2013-10-01')::timestamp,
			(date '2013-10-02')::timestamp,
			interval '1 day') AS "date"
			) t
		LEFT OUTER JOIN candle c ON t.date = date_trunc('day', c.timestamp)
		where
		(timestamp >= '2013-10-01' and timestamp <= '2013-10-02') or timestamp is null
		group by t.date, base, quote
		order by t.date, base, quote;
	`).Bind(context.TODO(), database.DB.SQL, &pairCounts)
	if err != nil {
		fmt.Println("candle.counts error", err)
	}
	return pairCounts, err
}

func Last(exchangeName, base, quote string, interval int64, asset string) (out Candle, err error) {
	// boil.DebugMode = true
	// defer func() { boil.DebugMode = false }()
	if exchangeName == "" || base == "" || quote == "" || asset == "" || interval <= 0 {
		return out, errInvalidInput
	}

	queries := []qm.QueryMod{
		qm.Where("base = ?", strings.ToUpper(base)),
		qm.Where("quote = ?", strings.ToUpper(quote)),
		qm.Where("interval = ?", interval),
		qm.Where("asset = ?", strings.ToLower(asset)),
		qm.OrderBy("timestamp desc"),
		qm.Limit(1),
	}

	exchangeUUID, errS := exchange.UUIDByName(exchangeName)
	if errS != nil {
		return out, errS
	}
	queries = append(queries, qm.Where("exchange_name_id = ?", exchangeUUID.String()))
	retCandle, errC := modelPSQL.Candles(queries...).All(context.Background(), database.DB.SQL)
	if errC != nil {
		return out, errC
	}

	if len(retCandle) > 0 {
		out = Candle{
			Timestamp:        retCandle[0].Timestamp,
			Open:             retCandle[0].Open,
			High:             retCandle[0].High,
			Low:              retCandle[0].Low,
			Close:            retCandle[0].Close,
			Volume:           retCandle[0].Volume,
			SourceJobID:      retCandle[0].SourceJobID.String,
			ValidationJobID:  retCandle[0].ValidationJobID.String,
			ValidationIssues: retCandle[0].ValidationIssues.String,
		}

		if out.Timestamp.IsZero() {
			return out, fmt.Errorf("%w: %s %s %s %v %s", ErrNoCandleDataFound, exchangeName, base, quote, interval, asset)
		}
	} else {
		// fmt.Println("candle.Last: no candle data for", base, quote)
		// os.Exit(123)
		return out, errNoCandleData
	}

	return out, err
}

// Series returns candle data
func Series(exchangeName, base, quote string, interval int64, asset string, start, end time.Time) (out Item, err error) {
	// defer func() { boil.DebugMode = false }()
	// boil.DebugMode = true
	if exchangeName == "" || base == "" || quote == "" || asset == "" || interval <= 0 {
		return out, errInvalidInput
	}

	queries := []qm.QueryMod{
		qm.Where("base = ?", strings.ToUpper(base)),
		qm.Where("quote = ?", strings.ToUpper(quote)),
		qm.Where("interval = ?", interval),
		qm.Where("asset = ?", strings.ToLower(asset)),
		qm.OrderBy("timestamp"),
	}

	exchangeUUID, errS := exchange.UUIDByName(exchangeName)
	if errS != nil {
		return out, errS
	}
	queries = append(queries, qm.Where("exchange_name_id = ?", exchangeUUID.String()))
	queries = append(queries, qm.Where("timestamp between ? and ?", start.UTC(), end.UTC()))
	retCandle, errC := modelPSQL.Candles(queries...).All(context.Background(), database.DB.SQL)
	if errC != nil {
		return out, errC
	}

	for x := range retCandle {
		out.Candles = append(out.Candles, Candle{
			Timestamp:        retCandle[x].Timestamp,
			Open:             retCandle[x].Open,
			High:             retCandle[x].High,
			Low:              retCandle[x].Low,
			Close:            retCandle[x].Close,
			Volume:           retCandle[x].Volume,
			SourceJobID:      retCandle[x].SourceJobID.String,
			ValidationJobID:  retCandle[x].ValidationJobID.String,
			ValidationIssues: retCandle[x].ValidationIssues.String,
		})
	}
	if len(out.Candles) < 1 {
		return out, fmt.Errorf("%w: %s %s %s %v %s", ErrNoCandleDataFound, exchangeName, base, quote, interval, asset)
	}

	out.ExchangeID = exchangeName
	out.Interval = interval
	out.Base = base
	out.Quote = quote
	out.Asset = asset
	return out, err
}

// DeleteCandles will delete all existing matching candles
func DeleteCandles(in *Item) (int64, error) {
	if database.DB.SQL == nil {
		return 0, database.ErrDatabaseSupportDisabled
	}
	if len(in.Candles) < 1 {
		return 0, errNoCandleData
	}

	ctx := context.Background()
	queries := []qm.QueryMod{
		qm.Where("base = ?", strings.ToUpper(in.Base)),
		qm.Where("quote = ?", strings.ToUpper(in.Quote)),
		qm.Where("interval = ?", in.Interval),
		qm.Where("asset = ?", strings.ToLower(in.Asset)),
		qm.Where("exchange_name_id = ?", in.ExchangeID),
	}

	queries = append(queries, qm.Where("timestamp between ? and ?", in.Candles[0].Timestamp.UTC(), in.Candles[len(in.Candles)-1].Timestamp.UTC()))
	return deletePostgres(ctx, queries)
}

func deletePostgres(ctx context.Context, queries []qm.QueryMod) (int64, error) {
	retCandle, err := modelPSQL.Candles(queries...).All(context.Background(), database.DB.SQL)
	if err != nil {
		return 0, err
	}
	var tx *sql.Tx
	tx, err = database.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	var totalDeleted int64
	totalDeleted, err = retCandle.DeleteAll(ctx, tx)
	if err != nil {
		errRB := tx.Rollback()
		if errRB != nil {
			log.Errorln(log.DatabaseMgr, errRB)
		}
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return totalDeleted, nil
}

// Insert series of candles
func Insert(in *Item) (uint64, error) {
	if database.DB.SQL == nil {
		return 0, database.ErrDatabaseSupportDisabled
	}

	if len(in.Candles) < 1 {
		return 0, errNoCandleData
	}

	ctx := context.Background()
	tx, err := database.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	var totalInserted uint64
	totalInserted, err = insertPostgresSQL(ctx, tx, in)
	if err != nil {
		errRB := tx.Rollback()
		if errRB != nil {
			log.Errorln(log.DatabaseMgr, errRB)
		}
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return totalInserted, nil
}

func insertPostgresSQL(ctx context.Context, tx *sql.Tx, in *Item) (uint64, error) {
	var totalInserted uint64
	for x := range in.Candles {
		var tempCandle = modelPSQL.Candle{
			ExchangeNameID: in.ExchangeID,
			Base:           strings.ToUpper(in.Base),
			Quote:          strings.ToUpper(in.Quote),
			Interval:       in.Interval,
			Asset:          strings.ToLower(in.Asset),
			Timestamp:      in.Candles[x].Timestamp,
			Open:           in.Candles[x].Open,
			High:           in.Candles[x].High,
			Low:            in.Candles[x].Low,
			Close:          in.Candles[x].Close,
			Volume:         in.Candles[x].Volume,
		}
		tempCandle.ValidationJobID = null.String{String: in.Candles[x].ValidationJobID, Valid: in.Candles[x].ValidationJobID != ""}
		tempCandle.ValidationIssues = null.String{String: in.Candles[x].ValidationIssues, Valid: in.Candles[x].ValidationIssues != ""}
		tempCandle.SourceJobID = null.String{String: in.Candles[x].SourceJobID, Valid: in.Candles[x].SourceJobID != ""}
		err := tempCandle.Upsert(ctx, tx, true, []string{"timestamp", "exchange_name_id", "base", "quote", "interval", "asset"}, boil.Infer(), boil.Infer())
		if err != nil {
			return 0, err
		}
		if totalInserted < math.MaxUint64 {
			totalInserted++
		}
	}
	return totalInserted, nil
}

// InsertFromCSV load a CSV list of candle data and insert into database
func InsertFromCSV(exchangeName, base, quote string, interval int64, asset, file string) (uint64, error) {
	csvFile, err := os.Open(file)
	if err != nil {
		return 0, err
	}

	defer func() {
		err = csvFile.Close()
		if err != nil {
			log.Errorln(log.Global, err)
		}
	}()

	csvData := csv.NewReader(csvFile)

	exchangeUUID, err := exchange.UUIDByName(exchangeName)
	if err != nil {
		return 0, err
	}

	tempCandle := &Item{
		ExchangeID: exchangeUUID.String(),
		Base:       base,
		Quote:      quote,
		Interval:   interval,
		Asset:      asset,
	}

	for {
		row, errCSV := csvData.Read()
		if errCSV != nil {
			if errCSV == io.EOF {
				break
			}
			return 0, errCSV
		}

		tempTick := Candle{}
		v, errParse := strconv.ParseInt(row[0], 10, 32)
		if errParse != nil {
			return 0, errParse
		}
		tempTick.Timestamp = time.Unix(v, 0).UTC()
		if tempTick.Timestamp.IsZero() {
			err = fmt.Errorf("invalid timestamp received on row %v", row)
			break
		}

		tempTick.Volume, err = strconv.ParseFloat(row[1], 64)
		if err != nil {
			break
		}

		tempTick.Open, err = strconv.ParseFloat(row[2], 64)
		if err != nil {
			break
		}

		tempTick.High, err = strconv.ParseFloat(row[3], 64)
		if err != nil {
			break
		}

		tempTick.Low, err = strconv.ParseFloat(row[4], 64)
		if err != nil {
			break
		}

		tempTick.Close, err = strconv.ParseFloat(row[5], 64)
		if err != nil {
			break
		}
		tempCandle.Candles = append(tempCandle.Candles, tempTick)
	}
	if err != nil {
		return 0, err
	}

	return Insert(tempCandle)
}
