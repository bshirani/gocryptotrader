package strategy

import (
	"context"
	"time"

	"gocryptotrader/database"
	"gocryptotrader/database/models/postgres"
	"gocryptotrader/exchange/order"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Details struct {
	ID            int
	Side          order.Side
	Capture       string
	TimeframeDays int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func All() (st []Details, err error) {
	query := postgres.Strategies()
	var result []*postgres.Strategy
	result, err = query.All(context.Background(), database.DB.SQL)
	if err != nil {
		return st, err
	}

	for _, r := range result {
		st = append(st, Details{
			ID:            r.ID,
			Side:          order.Side(r.Side),
			Capture:       r.Capture,
			TimeframeDays: r.TimeframeDays,
			CreatedAt:     r.CreatedAt,
			UpdatedAt:     r.UpdatedAt,
		})
	}

	return st, nil
}

func One(id int) (st Details, err error) {
	// query := postgres.Datahistoryjobs(qm.Where("nickname = ?", strings.ToLower(nickname)))
	query := postgres.Strategies(qm.Where("id=?", id))
	var r *postgres.Strategy
	r, err = query.One(context.Background(), database.DB.SQL)
	if err != nil {
		return st, err
	}

	st = Details{
		ID:            r.ID,
		Side:          order.Side(r.Side),
		Capture:       r.Capture,
		TimeframeDays: r.TimeframeDays,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}

	return st, nil
}
