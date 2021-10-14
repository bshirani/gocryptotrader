package liveorder

import (
	"context"
	"database/sql"

	"gocryptotrader/database"
	"gocryptotrader/database/models/postgres"
	"gocryptotrader/log"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func Count() int64 {
	i, _ := postgres.LiveOrders().Count(context.Background(), database.DB.SQL)
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
	ret, errS := postgres.LiveOrders(whereQM).One(context.Background(), database.DB.SQL)
	out.ID = int64(ret.ID)
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
	ret, errS := postgres.LiveOrders(whereQM).All(context.Background(), database.DB.SQL)

	for _, x := range ret {
		out = append(out, Details{
			ID: int64(x.ID),
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
func Insert(in Details) (int64, error) {
	if database.DB.SQL == nil {
		return 0, database.ErrDatabaseSupportDisabled
	}

	ctx := context.Background()
	tx, err := database.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	id, err := insertPostgresql(ctx, tx, in)

	if err != nil {
		errRB := tx.Rollback()
		if errRB != nil {
			log.Errorln(log.DatabaseMgr, errRB)
		}
		return id, err
	}

	err = tx.Commit()
	if err != nil {
		return id, err
	}

	return id, nil
}

func insertPostgresql(ctx context.Context, tx *sql.Tx, in Details) (id int64, err error) {
	var tempInsert = postgres.LiveOrder{
		Status:     in.Status.String(),
		OrderType:  in.OrderType.String(),
		Exchange:   in.Exchange,
		InternalID: in.InternalID,
		StrategyID: in.StrategyID,
		UpdatedAt:  in.UpdatedAt,
		CreatedAt:  in.CreatedAt,
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

	return int64(tempInsert.ID), nil
}
