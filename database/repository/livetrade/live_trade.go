package livetrade

import (
	"context"
	"database/sql"
	"fmt"

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
	// boil.DebugMode = true
	if database.DB.SQL == nil {
		return out, database.ErrDatabaseSupportDisabled
	}

	whereQM := qm.Where("status IN ('OPEN')")
	ret, errS := modelSQLite.LiveTrades(whereQM).All(context.Background(), database.DB.SQL)

	for _, x := range ret {
		out = append(out, Details{
			EntryPrice: decimal.NewFromFloat(x.EntryPrice),
			ID:         x.ID,
			StrategyID: x.StrategyID,
			Status:     order.Status(x.Status),
			Side:       order.Side(x.Side),
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
func Insert(in Details) error {
	if database.DB.SQL == nil {
		return database.ErrDatabaseSupportDisabled
	}

	ctx := context.Background()
	tx, err := database.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	err = insertSQLite(ctx, tx, []Details{in})

	if err != nil {
		errRB := tx.Rollback()
		if errRB != nil {
			log.Errorln(log.DatabaseMgr, errRB)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func insertSQLite(ctx context.Context, tx *sql.Tx, in []Details) (err error) {
	// boil.DebugMode = true
	for x := range in {
		entryPrice, _ := in[x].EntryPrice.Float64()
		exitPrice, _ := in[x].ExitPrice.Float64()
		stopLossPrice, _ := in[x].StopLossPrice.Float64()

		var tempInsert = modelSQLite.LiveTrade{
			EntryPrice:    entryPrice,
			ExitPrice:     null.Float64{Float64: exitPrice},
			StopLossPrice: stopLossPrice,
			Status:        fmt.Sprintf("%s", in[x].Status),
			StrategyID:    in[x].StrategyID,
			Pair:          in[x].Pair.String(),
			EntryOrderID:  in[x].EntryOrderID,
			Side:          in[x].Side.String(),
		}

		err = tempInsert.Insert(ctx, tx, boil.Infer())
		if err != nil {
			log.Errorln(log.DatabaseMgr, err)
			errRB := tx.Rollback()
			if errRB != nil {
				log.Errorln(log.DatabaseMgr, errRB)
			}
			return err
		}
	}

	return nil
}
