package strategy

import (
	"context"
	"time"

	"gocryptotrader/database"
	"gocryptotrader/database/models/postgres"
	"gocryptotrader/exchange/order"
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
