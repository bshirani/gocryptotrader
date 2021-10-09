package livetrade

import (
	"context"
	"database/sql"

	"github.com/thrasher-corp/gocryptotrader/database"
	modelSQLite "github.com/thrasher-corp/gocryptotrader/database/models/sqlite3"
	"github.com/thrasher-corp/gocryptotrader/log"
	"github.com/thrasher-corp/sqlboiler/boil"
	"github.com/thrasher-corp/sqlboiler/queries/qm"
	"github.com/volatiletech/null"
)

func Count() int64 {
	// whereQM := qm.Where("id = ?", 5)
	i, _ := modelSQLite.LiveTrades().Count(context.Background(), database.DB.SQL)
	return i
}

func OneByID(in string) (Details, error) {
	return one(in, "id")
}

// one returns one exchange by clause
func one(in, clause string) (out Details, err error) {
	if database.DB.SQL == nil {
		return out, database.ErrDatabaseSupportDisabled
	}

	whereQM := qm.Where(clause+"= ?", in)
	ret, errS := modelSQLite.Exchanges(whereQM).One(context.Background(), database.DB.SQL)
	out.ID = ret.ID
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
	for x := range in {
		var tempInsert = modelSQLite.LiveTrade{
			EntryPrice: in[x].EntryPrice,
			ExitPrice:  null.Float64{Float64: in[x].ExitPrice},
			StopPrice:  in[x].StopPrice,
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
