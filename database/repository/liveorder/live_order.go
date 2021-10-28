package liveorder

import (
	"context"
	"database/sql"
	"fmt"

	"gocryptotrader/database"
	"gocryptotrader/database/models/postgres"
	"gocryptotrader/exchange/order"
	"gocryptotrader/log"

	"github.com/shopspring/decimal"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func Count() int64 {
	i, _ := postgres.LiveOrders().Count(context.Background(), database.DB.SQL)
	return i
}

func OneByStrategyID(in int) (Details, error) {
	return one(in, "strategy_id")
}

func OneByID(in int) (Details, error) {
	return one(in, "id")
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
				log.Errorf(log.DatabaseMgr, "DeleteOrders tx.Rollback %v", errRB)
			}
		}
	}()

	query := postgres.LiveOrders()
	_, err = query.DeleteAll(ctx, tx)
	if err != nil {
		panic(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
		return err
	}

	if Count() > 0 {
		panic("did not delete")
	}

	return nil
}

func ActiveForStrategyName(sname string) (out []Details, err error) {
	// boil.DebugMode = true
	// defer func() { boil.DebugMode = false }()
	if database.DB.SQL == nil {
		return out, database.ErrDatabaseSupportDisabled
	}

	whereQM := qm.Where("status IN ('NEW', 'ACTIVE') AND strategy_name = ?", sname)
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

func Active() (out []Details, err error) {
	// boil.DebugMode = true
	if database.DB.SQL == nil {
		return out, database.ErrDatabaseSupportDisabled
	}

	whereQM := qm.Where("status IN ('NEW', 'ACTIVE')")
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

func Create(in *order.Submit) (int, error) {
	return Upsert(in)
}

// Insert writes a single entry into database
func Upsert(in *order.Detail) (int, error) {
	// boil.DebugMode = true
	// defer func() { boil.DebugMode = false }()
	if database.DB.SQL == nil {
		return 0, database.ErrDatabaseSupportDisabled
	}
	if in.Amount == 0 {
		panic(fmt.Errorf("order amount cannot be zero"))
	}
	if in.Price == 0 {
		panic(fmt.Errorf("order price cannot be zero"))
	}

	ctx := context.Background()
	tx, err := database.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	id, err := upsertPostgresql(ctx, tx, in)

	if err != nil {
		panic(err)
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

// Insert writes a single entry into database
func Update(in *Details) (int64, error) {
	//
	// boil.DebugMode = true
	if database.DB.SQL == nil {
		return 0, database.ErrDatabaseSupportDisabled
	}

	if in.Amount.IsZero() {
		panic(fmt.Errorf("order amount cannot be zero"))
	}
	if in.Price.IsZero() {
		panic(fmt.Errorf("order price cannot be zero"))
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
		stopLossPrice, _ := in[x].StopLossPrice.Float64()
		takeProfitPrice, _ := in[x].TakeProfitPrice.Float64()
		price, _ := in[x].Price.Float64()
		amount, _ := in[x].Amount.Float64()

		// if in[x].EntryTime.IsZero() {
		// 	fmt.Println("entrytimezero")
		// 	log.Errorln(log.DatabaseMgr, "entry time zero")
		// 	os.Exit(2)
		// }
		var tempUpdate = postgres.LiveOrder{
			ID:                in[x].ID,
			Status:            in[x].Status.String(),
			OrderType:         in[x].OrderType.String(),
			FilledAt:          null.NewTime(in[x].FilledAt, !in[x].FilledAt.IsZero()),
			InternalOrderType: in[x].InternalOrderType.String(),
			StopLossPrice:     stopLossPrice,
			TakeProfitPrice:   takeProfitPrice,
			Amount:            amount,
			Price:             price,
			Exchange:          in[x].Exchange,
			InternalID:        in[x].InternalID,
			StrategyName:      in[x].StrategyName,
			UpdatedAt:         in[x].UpdatedAt,
			CreatedAt:         in[x].CreatedAt,
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

// func insertPostgresql(ctx context.Context, tx *sql.Tx, in *order.Submit) (id int, err error) {
// 	// boil.DebugMode = true
// 	// defer func() { boil.DebugMode = false }()
// 	// fmt.Println("insert order", in.Type, in.StopLossPrice, in.Status, "active_at", in.Date)
// 	var tempInsert = postgres.LiveOrder{
// 		Status:            in.Status.String(),
// 		OrderType:         in.Type.String(),
// 		Amount:            in.Amount,
// 		InternalOrderType: in.InternalOrderType.String(),
// 		Exchange:          in.Exchange,
// 		Side:              in.Side.String(),
// 		Price:             in.Price,
// 		StopLossPrice:     in.StopLossPrice,
// 		TakeProfitPrice:   in.TakeProfitPrice,
// 		ClientOrderID:     in.ID,
// 		StrategyName:      in.StrategyName,
// 	}
//
// 	if in.Status == order.Active {
// 		tempInsert.ActiveAt = null.NewTime(in.Date, true)
// 	}
//
// 	err = tempInsert.Insert(ctx, tx, boil.Infer())
// 	// err = tempInsert.Upsert(ctx, tx, true, []string{"id"}, boil.Infer(), boil.Infer())
// 	if err != nil {
// 		log.Errorln(log.DatabaseMgr, err)
// 		errRB := tx.Rollback()
// 		panic(err)
// 		if errRB != nil {
// 			log.Errorln(log.DatabaseMgr, errRB)
// 		}
// 		return 0, err
// 	}
//
// 	return tempInsert.ID, nil
// }

func upsertPostgresql(ctx context.Context, tx *sql.Tx, in *order.Detail) (id int, err error) {
	// fmt.Println("upsert order!!!!!", "type", in.Type, "st", in.Status, "fa", in.FilledAt, "fanull")
	// boil.DebugMode = true
	// defer func() { boil.DebugMode = false }()

	var tempInsert = postgres.LiveOrder{
		ID:              in.InternalOrderID,
		Status:          in.Status.String(),
		OrderType:       in.Type.String(),
		Exchange:        in.Exchange,
		Side:            in.Side.String(),
		Price:           in.Price,
		StopLossPrice:   in.StopLossPrice,
		TakeProfitPrice: in.TakeProfitPrice,
		ClientOrderID:   in.ID,
		StrategyName:    in.StrategyName,
	}

	if in.Status == order.Filled {
		tempInsert.FilledAt = null.NewTime(in.Date, true)
	} else if in.Status == order.Cancelled {
		tempInsert.CancelledAt = null.NewTime(in.Date, true)
	} else if in.Status == order.Active {
		tempInsert.ActiveAt = null.NewTime(in.Date, true)
	}

	err = tempInsert.Upsert(ctx, tx, true, []string{"id"}, boil.Infer(), boil.Infer())
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

func one(in int, clause string) (out Details, err error) {
	if database.DB.SQL == nil {
		return out, database.ErrDatabaseSupportDisabled
	}
	// boil.DebugMode = true

	whereQM := qm.Where("id=?", in)
	ret, errS := postgres.LiveOrders(whereQM).One(context.Background(), database.DB.SQL)

	out = Details{
		ID:                ret.ID,
		Status:            order.Status(ret.Status),
		OrderType:         order.Type(ret.OrderType),
		Amount:            decimal.NewFromFloat(ret.Amount),
		Price:             decimal.NewFromFloat(ret.Price),
		Exchange:          ret.Exchange,
		InternalID:        ret.InternalID,
		InternalOrderType: order.InternalOrderType(ret.InternalOrderType),
		StrategyName:      ret.StrategyName,
		UpdatedAt:         ret.UpdatedAt,
		CreatedAt:         ret.CreatedAt,
	}

	if errS != nil {
		return out, errS
	}

	return out, err
}
