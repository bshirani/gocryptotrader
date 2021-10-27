package livetrade

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
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

func ActiveForStrategyName(sname string) (out []Details, err error) {
	// boil.DebugMode = true
	// defer func() { boil.DebugMode = false }()
	if database.DB.SQL == nil {
		return out, database.ErrDatabaseSupportDisabled
	}

	whereQM := qm.Where("status IN ('OPEN') AND strategy_name = ?", sname)
	ret, errS := postgres.LiveTrades(whereQM).All(context.Background(), database.DB.SQL)
	for _, x := range ret {
		out = append(out, Details{
			ID: x.ID,
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

		pair, _ := currency.NewPairFromString(x.Pair)

		out = append(out, Details{
			EntryPrice:    decimal.NewFromFloat(x.EntryPrice),
			StopLossPrice: decimal.NewFromFloat(x.StopLossPrice),
			// TakeProfitPrice: decimal.NewFromFloat(x.TakeProfitPrice),
			// ExitTime:     x.ExitTime,
			Pair:         pair,
			EntryTime:    x.EntryTime,
			Amount:       decimal.NewFromFloat(x.Amount),
			ID:           x.ID,
			StrategyName: x.StrategyName,
			Status:       order.Status(x.Status),
			Side:         order.Side(x.Side),
			UpdatedAt:    x.UpdatedAt,
			CreatedAt:    x.CreatedAt,
			EntryOrderID: x.EntryOrderID,
			RiskedQuote:  x.RiskedQuote,
			RiskedPoints: x.RiskedPoints,
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
	id, err = upsertPostgresql(ctx, tx, in)

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
func Update(in *Details) (int, error) {
	if database.DB.SQL == nil {
		return 0, database.ErrDatabaseSupportDisabled
	}

	ctx := context.Background()
	tx, err := database.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	id, err := upsertPostgresql(ctx, tx, *in)

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

func upsertPostgresql(ctx context.Context, tx *sql.Tx, in Details) (id int, err error) {
	// boil.DebugMode = true
	// defer func() { boil.DebugMode = false }()
	entryPrice, _ := in.EntryPrice.Float64()
	exitPrice, _ := in.ExitPrice.Float64()
	stopLossPrice, _ := in.StopLossPrice.Float64()
	takeProfitPrice, _ := in.TakeProfitPrice.Float64()
	profitLossPoints, _ := in.ProfitLossPoints.Float64()
	profitLossQuote, _ := in.ProfitLossQuote.Float64()
	amount, _ := in.Amount.Float64()
	// riskedQuote, _ := in.RiskedQuote.Float64()
	// riskedPoints, _ := in.RiskedPoints.Float64()

	if stopLossPrice < 0 {
		return 0, fmt.Errorf("stop loss price cannot be below zero")
	}

	if in.EntryOrderID == 0 {
		panic("entry order id cannot be 0")
	}

	var tempInsert = postgres.LiveTrade{
		ID:            in.ID,
		EntryPrice:    entryPrice,
		EntryTime:     in.EntryTime,
		StopLossPrice: stopLossPrice,
		Status:        fmt.Sprintf("%s", in.Status),
		StrategyName:  in.StrategyName,
		Pair:          in.Pair.String(),
		EntryOrderID:  in.EntryOrderID,
		Side:          in.Side.String(),
		Amount:        amount,
		RiskedQuote:   in.RiskedQuote,
		RiskedPoints:  in.RiskedPoints,

		ExitTime:         null.NewTime(in.ExitTime, !in.ExitTime.IsZero()),
		ExitPrice:        null.NewFloat64(exitPrice, exitPrice != 0),
		TakeProfitPrice:  null.NewFloat64(takeProfitPrice, takeProfitPrice != 0),
		ProfitLossPoints: null.NewFloat64(profitLossPoints, profitLossPoints != 0 || exitPrice != 0),
		ProfitLossQuote:  null.NewFloat64(profitLossQuote, profitLossQuote != 0 || exitPrice != 0),
	}

	err = tempInsert.Upsert(
		ctx,
		tx,
		true,
		[]string{"entry_order_id"},
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

func WriteCSV(trades []*Details, newpath string) {
	// var nickName string
	// if d.Config.Nickname != "" {
	// 	nickName = d.Config.Nickname + "-"
	// }
	fmt.Println("writing to", newpath)

	file, err := os.OpenFile(newpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println("error", err)
	}

	header := "strategy,pair,direction,entry_time,exit_time,entry_price,exit_price,stop_loss,amount,take_profit\n"
	file.WriteString(header)
	for _, t := range trades {
		s := fmt.Sprintf(
			"%d,%s,%s,%v,%v,%v,%v,%v,%v,%v\n",
			t.StrategyID,
			t.Pair,
			t.Side,
			t.EntryTime.Format(common.SimpleTimeFormat),
			t.ExitTime.Format(common.SimpleTimeFormat),
			t.EntryPrice,
			t.ExitPrice,
			t.StopLossPrice,
			t.Amount,
			t.TakeProfitPrice,
		)
		file.WriteString(s)
	}
	file.Close()
}

// LoadCSV loads & parses a CSV list of exchanges
func LoadCSV(file string) (out []*Details, err error) {
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
			continue
		}
		count += 1

		sid, err := strconv.ParseInt(row[0], 10, 32)
		pair, err := currency.NewPairFromString(row[1])
		entryTime, err := time.Parse(common.SimpleTimeFormat, row[3])
		exitTime, err := time.Parse(common.SimpleTimeFormat, row[4])
		entryPrice, err := decimal.NewFromString(row[5])
		exitPrice, err := decimal.NewFromString(row[6])
		stop, err := decimal.NewFromString(row[7])
		amount, err := decimal.NewFromString(row[8])
		takeProfit, err := decimal.NewFromString(row[9])

		out = append(out, &Details{
			StrategyID:      int(sid),
			Pair:            pair,
			Side:            order.Side(row[2]),
			EntryTime:       entryTime,
			ExitTime:        exitTime,
			EntryPrice:      entryPrice,
			ExitPrice:       exitPrice,
			StopLossPrice:   stop,
			TakeProfitPrice: takeProfit,
			Amount:          amount,
		})
		if err != nil {
			fmt.Println("error", err)
		}
	}
	return out, err
}
