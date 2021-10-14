package livetrade

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"gocryptotrader/database"
	modelSQLite "gocryptotrader/database/models/sqlite3"
	"gocryptotrader/exchange/order"
	"gocryptotrader/log"

	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/sqlboiler/boil"
	"github.com/thrasher-corp/sqlboiler/queries/qm"
	"github.com/volatiletech/null"
)

func Count() int64 {
	i, _ := modelSQLite.LiveTrades().Count(context.Background(), database.DB.SQL)
	return i
}

func OneByStrategyID(in string) (Details, error) {
	return one(in, "strategy_id")
}

func OneByID(in string) (Details, error) {
	return one(in, "id")
}

func one(in, clause string) (out Details, err error) {
	if database.DB.SQL == nil {
		return out, database.ErrDatabaseSupportDisabled
	}
	// boil.DebugMode = true

	whereQM := qm.Where(clause+"= ?", in)
	ret, errS := modelSQLite.LiveTrades(whereQM).One(context.Background(), database.DB.SQL)
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
	ret, errS := modelSQLite.LiveTrades(whereQM).All(context.Background(), database.DB.SQL)
	ret.ReloadAll(context.Background(), database.DB.SQL)
	layout2 := time.RFC3339

	for _, x := range ret {
		fmt.Println("parsing entry time", x.EntryTime, x.CreatedAt, x.UpdatedAt)
		// IntervalStartTime: results[i].IntervalStartDate.UTC().Format(time.RFC3339),
		// fmt.Println("parsing entry time", x.EntryTime, entryTime)
		entryTime, _ := time.Parse(layout2, x.EntryTime)
		if entryTime.IsZero() {
			fmt.Println("ERROR entryTime is zero")
			os.Exit(2)
		}
		// exitTime, _ := time.Parse(layout2, x.ExitTime)
		updatedAt, _ := time.Parse(layout2, x.UpdatedAt)
		createdAt, _ := time.Parse(layout2, x.CreatedAt)
		out = append(out, Details{
			EntryPrice: decimal.NewFromFloat(x.EntryPrice),
			EntryTime:  entryTime,
			// ExitTime:   exitTime,
			ID:         x.ID,
			StrategyID: x.StrategyID,
			Status:     order.Status(x.Status),
			Side:       order.Side(x.Side),
			UpdatedAt:  updatedAt,
			CreatedAt:  createdAt,
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
func Insert(in Details) (id int64, err error) {
	if database.DB.SQL == nil {
		return 0, database.ErrDatabaseSupportDisabled
	}

	ctx := context.Background()
	tx, err := database.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	id, err = insertSQLite(ctx, tx, in)

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

	return id, nil
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
	fmt.Println("updating trade status", in.Status, in.ID)
	id, err := updateSQLite(ctx, tx, []Details{*in})

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

	return id, nil
}

func insertSQLite(ctx context.Context, tx *sql.Tx, in Details) (id int64, err error) {
	// boil.DebugMode = true
	entryPrice, _ := in.EntryPrice.Float64()
	exitPrice, _ := in.ExitPrice.Float64()
	stopLossPrice, _ := in.StopLossPrice.Float64()
	fmt.Println("inserting sqlite trade", in.EntryTime)

	var tempInsert = modelSQLite.LiveTrade{
		EntryPrice:    entryPrice,
		CreatedAt:     time.Now().UTC().Format(time.RFC3339),
		EntryTime:     in.EntryTime.String(),
		ExitTime:      null.String{String: in.ExitTime.String()},
		ExitPrice:     null.Float64{Float64: exitPrice},
		StopLossPrice: stopLossPrice,
		Status:        fmt.Sprintf("%s", in.Status),
		StrategyID:    in.StrategyID,
		Pair:          in.Pair.String(),
		EntryOrderID:  in.EntryOrderID,
		Side:          in.Side.String(),
	}

	err = tempInsert.Insert(ctx, tx, boil.Infer())
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

func updateSQLite(ctx context.Context, tx *sql.Tx, in []Details) (id int64, err error) {
	// boil.DebugMode = true
	for x := range in {
		entryPrice, _ := in[x].EntryPrice.Float64()
		exitPrice, _ := in[x].ExitPrice.Float64()
		stopLossPrice, _ := in[x].StopLossPrice.Float64()

		if in[x].EntryTime.IsZero() {
			fmt.Println("entrytimezero")
			os.Exit(2)
		} else {
			fmt.Println("saving trade", in[x].EntryTime)
		}
		var tempInsert = modelSQLite.LiveTrade{
			ID:            in[x].ID,
			UpdatedAt:     time.Now().String(),
			EntryPrice:    entryPrice,
			EntryTime:     in[x].EntryTime.UTC().String(),
			ExitTime:      null.String{String: in[x].ExitTime.String()},
			ExitPrice:     null.Float64{Float64: exitPrice},
			StopLossPrice: stopLossPrice,
			Status:        fmt.Sprintf("%s", in[x].Status),
			StrategyID:    in[x].StrategyID,
			Pair:          in[x].Pair.String(),
			EntryOrderID:  in[x].EntryOrderID,
			Side:          in[x].Side.String(),
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
