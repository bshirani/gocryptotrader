package livetrade

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gocryptotrader/common"
	"gocryptotrader/currency"
	"gocryptotrader/database"
	"gocryptotrader/exchange/order"

	"gocryptotrader/database/models/postgres"
	"gocryptotrader/log"

	"github.com/shopspring/decimal"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func Count() int64 {
	i, _ := postgres.LiveTrades().Count(context.Background(), database.DB.SQL)
	return i
}

func OneByStrategyID(in int) (Details, error) {
	return one(in, "strategy_id")
}

func OneByID(in int) (Details, error) {
	return one(in, "id")
}

func one(in int, clause string) (out Details, err error) {
	if database.DB.SQL == nil {
		return out, database.ErrDatabaseSupportDisabled
	}
	// boil.DebugMode = true

	whereQM := qm.Where(clause+"= ?", in)
	ret, errS := postgres.LiveTrades(whereQM).One(context.Background(), database.DB.SQL)
	out.ID = ret.ID
	if errS != nil {
		return out, errS
	}

	return out, err
}

func Active() (out []Details, err error) {
	return ByStatus(order.Open)
}

func Closed() (out []Details, err error) {
	return ByStatus(order.Closed)
}

func ByStatus(status order.Status) (out []Details, err error) {
	// boil.DebugMode = true
	if database.DB.SQL == nil {
		return out, database.ErrDatabaseSupportDisabled
	}

	whereQM := qm.Where(fmt.Sprintf("status IN ('%s')", status))
	ret, errS := postgres.LiveTrades(whereQM).All(context.Background(), database.DB.SQL)

	for _, x := range ret {
		// fmt.Printf("ByStatus EntryTime  %d ep: %v en:%s up:%s cr:%s\n",
		// 	x.ID,
		// 	x.EntryPrice,
		// 	x.EntryTime,
		// 	x.UpdatedAt,
		// 	x.CreatedAt)

		if x.EntryTime.IsZero() {
			fmt.Println("ERROR entryTime is zero")
			os.Exit(2)
		}

		out = append(out, Details{
			EntryPrice: decimal.NewFromFloat(x.EntryPrice),
			// ExitPrice:     decimal.NewFromFloat(x.ExitPrice),
			// ExitTime:      x.ExitTime,
			// ExitPrice:     null.Float64{Float64: x.ExitPrice},
			StopLossPrice: decimal.NewFromFloat(x.StopLossPrice),
			// Pair:         x.Pair,
			EntryTime:    x.EntryTime,
			Amount:       decimal.NewFromFloat(x.Amount),
			ID:           x.ID,
			StrategyID:   x.StrategyID,
			Status:       order.Status(x.Status),
			Side:         order.Side(x.Side),
			UpdatedAt:    x.UpdatedAt,
			CreatedAt:    x.CreatedAt,
			EntryOrderID: x.EntryOrderID,
		})
	}
	if errS != nil {
		return out, errS
	}

	if errS != nil {
		return out, errS
	}

	return out, err
}

// Insert writes a single entry into database
func Insert(in Details) (id int, err error) {
	if database.DB.SQL == nil {
		return 0, database.ErrDatabaseSupportDisabled
	}

	ctx := context.Background()
	tx, err := database.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	id, err = insertPostgresql(ctx, tx, in)

	if err != nil {
		errRB := tx.Rollback()
		if errRB != nil {
			log.Errorln(log.DatabaseMgr, errRB)
		}
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("error committing insert", err)
		return 0, err
	}

	return id, nil
}

func DeleteAll() error {
	ctx := context.Background()

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

	query := postgres.LiveTrades(qm.Where(`1=1`))
	_, err = query.DeleteAll(ctx, tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Insert writes a single entry into database
func Update(in *Details) (int64, error) {
	if database.DB.SQL == nil {
		return 0, database.ErrDatabaseSupportDisabled
	}

	ctx := context.Background()
	tx, err := database.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	id, err := updatePostgresql(ctx, tx, []Details{*in})

	if err != nil {
		errRB := tx.Rollback()
		if errRB != nil {
			log.Errorln(log.DatabaseMgr, errRB)
		}
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("error committing update", err)
		return 0, err
	}
	// record, err := postgres.LiveTrades().One(context.Background(), database.DB.SQL)
	// if err != nil {
	// 	fmt.Println("error retrieving update", err)
	// 	return 0, err
	// }
	// fmt.Println("after update record", record.EntryTime, record.EntryPrice, record.UpdatedAt)

	return id, nil
}

func insertPostgresql(ctx context.Context, tx *sql.Tx, in Details) (id int, err error) {
	// boil.DebugMode = true
	entryPrice, _ := in.EntryPrice.Float64()
	exitPrice, _ := in.ExitPrice.Float64()
	stopLossPrice, _ := in.StopLossPrice.Float64()
	amount, _ := in.Amount.Float64()

	if stopLossPrice < 0 {
		return 0, fmt.Errorf("stop loss price cannot be below zero")
	}

	var tempInsert = postgres.LiveTrade{
		EntryPrice: entryPrice,
		// CreatedAt:     time.Now(),
		EntryTime:     in.EntryTime,
		ExitTime:      null.Time{Time: in.ExitTime},
		ExitPrice:     null.Float64{Float64: exitPrice},
		StopLossPrice: stopLossPrice,
		Status:        fmt.Sprintf("%s", in.Status),
		StrategyID:    in.StrategyID,
		Pair:          in.Pair.String(),
		EntryOrderID:  in.EntryOrderID,
		Side:          in.Side.String(),
		Amount:        amount,
	}

	err = tempInsert.Upsert(
		ctx,
		tx,
		true,
		[]string{},
		boil.Infer(),
		boil.Infer())
	if err != nil {
		log.Errorln(log.DatabaseMgr, err)
		errRB := tx.Rollback()
		if errRB != nil {
			log.Errorln(log.DatabaseMgr, errRB)
		}
		return 0, err
	}

	return tempInsert.ID, nil
}

func updatePostgresql(ctx context.Context, tx *sql.Tx, in []Details) (id int64, err error) {
	// boil.DebugMode = true
	for x := range in {
		entryPrice, _ := in[x].EntryPrice.Float64()
		exitPrice, _ := in[x].ExitPrice.Float64()
		stopLossPrice, _ := in[x].StopLossPrice.Float64()
		amount, _ := in[x].StopLossPrice.Float64()

		if in[x].EntryTime.IsZero() {
			fmt.Println("entrytimezero")
			log.Errorln(log.DatabaseMgr, "entry time zero")
			os.Exit(2)
		}
		var tempInsert = postgres.LiveTrade{
			ID:            in[x].ID,
			UpdatedAt:     time.Now(),
			EntryPrice:    entryPrice,
			EntryTime:     in[x].EntryTime,
			ExitTime:      null.Time{Time: in[x].ExitTime},
			ExitPrice:     null.Float64{Float64: exitPrice},
			StopLossPrice: stopLossPrice,
			Status:        fmt.Sprintf("%s", in[x].Status),
			StrategyID:    in[x].StrategyID,
			Pair:          in[x].Pair.String(),
			EntryOrderID:  in[x].EntryOrderID,
			Side:          in[x].Side.String(),
			Amount:        amount,
		}

		id, err = tempInsert.Update(ctx, tx, boil.Infer())
		if err != nil {
			log.Errorln(log.DatabaseMgr, err)
			errRB := tx.Rollback()
			if errRB != nil {
				log.Errorln(log.DatabaseMgr, errRB)
			}
			return 0, err
		}
	}

	return id, nil
}

func WriteCSV(trades []*Details) {
	// var nickName string
	// if d.Config.Nickname != "" {
	// 	nickName = d.Config.Nickname + "-"
	// }
	fileName := fmt.Sprintf(
		"results/trades-%v.csv",
		time.Now().Format("2006-01-02-15-04-05"))
	newpath := filepath.Join(".", fileName)
	fmt.Println("writing to", newpath)

	file, err := os.OpenFile(newpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println("error", err)
	}

	header := "strategy,pair,direction,entry_time,exit_time,entry_price,exit_price,stop_loss,amount\n"
	file.WriteString(header)
	for _, t := range trades {
		s := fmt.Sprintf(
			"%d,%s,%s,%v,%v,%v,%v,%v,%v\n",
			t.StrategyID,
			t.Pair,
			t.Side,
			t.EntryTime.Format(common.SimpleTimeFormat),
			t.ExitTime.Format(common.SimpleTimeFormat),
			t.EntryPrice,
			t.ExitPrice,
			t.StopLossPrice,
			t.Amount,
		)
		file.WriteString(s)
	}
	fmt.Println("wrote trades CSV", newpath)
	file.Close()
}

// LoadCSV loads & parses a CSV list of exchanges
func LoadCSV(file string) (out []Details, err error) {
	csvFile, err := os.Open(file)
	if err != nil {
		return out, err
	}

	defer func() {
		err = csvFile.Close()
		if err != nil {
			log.Errorln(log.Global, err)
		}
	}()

	fmt.Println("rading", file)
	csvData := csv.NewReader(csvFile)
	count := 0
	for {
		row, errCSV := csvData.Read()
		if errCSV != nil {
			if errCSV == io.EOF {
				return out, err
			}
			return out, errCSV
		}

		if count == 0 {
			count += 1
			fmt.Println("HERE", row)
			continue
		}
		fmt.Println("row", row)
		count += 1
		id, err := strconv.ParseInt(row[0], 10, 64)
		pair, err := currency.NewPairFromString(row[1])
		entryTime, err := time.Parse(common.SimpleTimeFormat, row[3])
		exitTime, err := time.Parse(common.SimpleTimeFormat, row[4])
		entryPrice, err := decimal.NewFromString(row[5])
		exitPrice, err := decimal.NewFromString(row[6])
		stop, err := decimal.NewFromString(row[7])
		amount, err := decimal.NewFromString(row[8])
		fmt.Println("amount", amount)
		out = append(out, Details{
			StrategyID:    int(id),
			Pair:          pair,
			Side:          order.Side(row[2]),
			EntryTime:     entryTime,
			ExitTime:      exitTime,
			EntryPrice:    entryPrice,
			ExitPrice:     exitPrice,
			StopLossPrice: stop,
			Amount:        amount,
		})
		if err != nil {
			fmt.Println("error", err)
		}
	}
	return out, err
}

func LastResult() string {
	// return os.MkdirAll(dir, 0770)
	wd, err := os.Getwd()
	dir := filepath.Join(wd, "../backtest/results")
	lf := lastFileInDir(dir)

	if err != nil {
		fmt.Println(err)
	}
	return filepath.Join(wd, "../backtest/results", lf)
}

func lastFileInDir(dir string) string {
	files, _ := ioutil.ReadDir(dir)
	var modTime time.Time
	var names []string
	for _, fi := range files {
		if fi.Mode().IsRegular() {
			if !fi.ModTime().Before(modTime) {
				if fi.ModTime().After(modTime) {
					modTime = fi.ModTime()
					names = names[:0]
				}
				names = append(names, fi.Name())
			}
		}
	}
	if len(names) > 0 {
		fmt.Println(modTime, names)
	}
	return names[len(names)-1]
}
