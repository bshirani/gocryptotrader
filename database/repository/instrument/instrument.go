package instrument

import (
	"context"
	"database/sql"

	"gocryptotrader/database"
	modelPSQL "gocryptotrader/database/models/postgres"
	"gocryptotrader/log"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

// func OneByPair(pair currency.Pair) (Details, error) {
// 	out := Details{}
// 	if database.DB.SQL == nil {
// 		return out, database.ErrDatabaseSupportDisabled
// 	}
//
// 	exchange.One("gateio")
// 	whereQM := qm.Where("base = ?", pair.Base)
// 	ret, err := modelPSQL.Instruments(whereQM).One(context.Background(), database.DB.SQL)
// 	// out.Base = currency.NewCode(ret.Base)
// 	// out.Quote = currency.NewCode(ret.Quote)
//
// 	return out, err
// }

// // Upsert inserts or updates jobs into the database
// func (db *DBService) Upsert(jobs ...*DataHistoryJob) error {
// 	ctx := context.Background()
//
// 	tx, err := db.sql.BeginTx(ctx, nil)
// 	if err != nil {
// 		return fmt.Errorf("beginTx %w", err)
// 	}
// 	defer func() {
// 		if err != nil {
// 			errRB := tx.Rollback()
// 			if errRB != nil {
// 				log.Errorf(log.DatabaseMgr, "Insert tx.Rollback %v", errRB)
// 			}
// 		}
// 	}()
//
// 	err = upsertPostgres(ctx, tx, jobs...)
// 	if err != nil {
// 		return err
// 	}
//
// 	return tx.Commit()
// }

// func InsertCMC(in CryptoCurrencyMap) error {
// 	if database.DB.SQL == nil {
// 		return database.ErrDatabaseSupportDisabled
// 	}
//
// 	ctx := context.Background()
// 	tx, err := database.DB.SQL.BeginTx(ctx, nil)
// 	if err != nil {
// 		return err
// 	}
// 	err = insertPostgresql(ctx, tx, []Details{in})
//
// 	if err != nil {
// 		errRB := tx.Rollback()
// 		if errRB != nil {
// 			log.Errorln(log.DatabaseMgr, errRB)
// 		}
// 		return err
// 	}
//
// 	err = tx.Commit()
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

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
	err = insertPostgresql(ctx, tx, []Details{in})

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

// InsertMany writes multiple entries into database
func InsertMany(in []Details) error {
	if database.DB.SQL == nil {
		return database.ErrDatabaseSupportDisabled
	}

	ctx := context.Background()
	tx, err := database.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = insertPostgresql(ctx, tx, in)

	if err != nil {
		errRB := tx.Rollback()
		if errRB != nil {
			log.Errorln(log.DatabaseMgr, errRB)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		errRB := tx.Rollback()
		if errRB != nil {
			log.Errorln(log.DatabaseMgr, errRB)
		}
		return err
	}
	return nil
}

func insertPostgresql(ctx context.Context, tx *sql.Tx, in []Details) (err error) {
	for x := range in {
		// fmt.Println("inserting base", in[x].Base)
		var tempInsert = modelPSQL.Instrument{
			ID:                  in[x].ID,
			Symbol:              in[x].Symbol,
			CMCID:               in[x].CMCID,
			FirstHistoricalData: in[x].FirstHistoricalData,
			LastHistoricalData:  in[x].LastHistoricalData,
			Name:                in[x].Name,
			Active:              in[x].Active,
			Status:              in[x].Status,
			// Quote:               in[x].Quote.String(),
			// Base:                in[x].Base.String(),
			UpdatedAt: in[x].UpdatedAt,
			CreatedAt: in[x].CreatedAt,
		}

		err = tempInsert.Upsert(ctx, tx, true, []string{"id"}, boil.Infer(), boil.Infer())
		if err != nil {
			errRB := tx.Rollback()
			if errRB != nil {
				log.Errorln(log.DatabaseMgr, errRB)
			}
			return
		}
	}
	return nil
}
