package accountlog

import (
	"context"
	"gocryptotrader/database"
	"gocryptotrader/database/models/postgres"
	"gocryptotrader/log"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

func Insert(in Details) error {
	if database.DB.SQL == nil {
		return database.ErrDatabaseSupportDisabled
	}

	ctx := context.Background()
	tx, err := database.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var tempInsert = postgres.AccountLog{
		ID:         in.ID,
		UsdBalance: in.USDBalance,
		XRPBalance: in.XRPBalance,
		BTCBalance: in.BTCBalance,
		OpenTrades: in.OpenTrades,
		Timestamp:  in.Timestamp,
		CreatedAt:  in.CreatedAt,
		UpdatedAt:  in.UpdatedAt,
	}

	err = tempInsert.Upsert(ctx, tx, true, []string{"id"}, boil.Infer(), boil.Infer())
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
