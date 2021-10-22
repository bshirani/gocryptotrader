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

	whereQM := qm.Where("status IN ('NEW')")
	ret, errS := postgres.LiveOrders(whereQM).All(context.Background(), database.DB.SQL)

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
		return id, err
	}

	err = tx.Commit()
	if err != nil {
		return id, err
	}

	return id, nil
}

func updatePostgresql(ctx context.Context, tx *sql.Tx, in []Details) (id int64, err error) {
	// boil.DebugMode = true
	for x := range in {
		// entryPrice, _ := in[x].EntryPrice.Float64()
		// exitPrice, _ := in[x].ExitPrice.Float64()
		// stopLossPrice, _ := in[x].StopLossPrice.Float64()

		// if in[x].EntryTime.IsZero() {
		// 	fmt.Println("entrytimezero")
		// 	log.Errorln(log.DatabaseMgr, "entry time zero")
		// 	os.Exit(2)
		// }
		var tempUpdate = postgres.LiveOrder{
			ID:         in[x].ID,
			Status:     in[x].Status.String(),
			OrderType:  in[x].OrderType.String(),
			Exchange:   in[x].Exchange,
			InternalID: in[x].InternalID,
			StrategyID: in[x].StrategyID,
			UpdatedAt:  in[x].UpdatedAt,
			CreatedAt:  in[x].CreatedAt,
		}

		id, err = tempUpdate.Update(ctx, tx, boil.Infer())
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

// Insert writes a single entry into database
func Insert(in Details) (string, error) {
	if database.DB.SQL == nil {
		return "", database.ErrDatabaseSupportDisabled
	}

	ctx := context.Background()
	tx, err := database.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return "", err
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

func insertPostgresql(ctx context.Context, tx *sql.Tx, in Details) (id string, err error) {
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
		return "", err
	}

	return tempInsert.ID, nil
}
