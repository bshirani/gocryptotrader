package livesignal

import (
	"context"
	"fmt"
	"gocryptotrader/database"
	"gocryptotrader/database/models/postgres"
	"gocryptotrader/log"

	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

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
	var tempInsert = postgres.LiveSignal{
		ID:           in.ID,
		SignalTime:   in.SignalTime,
		StrategyName: in.StrategyName,
		Prediction:   in.Prediction,
		ValidUntil:   in.ValidUntil,
	}

	err = tempInsert.Upsert(
		ctx,
		tx,
		true,
		[]string{"id"},
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

	return tempInsert.ID, nil
}

func Active() (out []Details, err error) {
	boil.DebugMode = true
	defer func() { boil.DebugMode = false }()
	if database.DB.SQL == nil {
		return out, database.ErrDatabaseSupportDisabled
	}
	// if currentTime == nil {
	// 	return []Details{}, nil
	// }

	ret, errS := postgres.LiveSignals(
		Select("strategy_name", "valid_until", "signal_time"),
		// Where("valid_until < ?", currentTime),
		OrderBy("signal_time desc"),
		GroupBy("strategy_name"),
		GroupBy("valid_until"),
		GroupBy("signal_time"),
		Limit(1),
	).All(context.Background(), database.DB.SQL)

	for _, x := range ret {
		out = append(out, Details{
			StrategyName: x.StrategyName,
			SignalTime:   x.SignalTime,
			ValidUntil:   x.ValidUntil,
		})
	}

	if errS != nil {
		return out, errS
	}
	fmt.Println("signals returning", len(out))

	return out, err
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

	query := postgres.LiveSignals(Where(`1=1`))
	_, err = query.DeleteAll(ctx, tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
