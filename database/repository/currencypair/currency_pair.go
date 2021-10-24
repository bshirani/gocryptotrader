package currencypair

import (
	"context"
	"time"

	"gocryptotrader/currency"
	"gocryptotrader/database"
	"gocryptotrader/database/models/postgres"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Details struct {
	ID        int
	Pair      currency.Pair
	CreatedAt time.Time
	UpdatedAt time.Time
}

func All() (st []Details, err error) {
	query := postgres.CurrencyPairs()
	var result []*postgres.CurrencyPair
	result, err = query.All(context.Background(), database.DB.SQL)
	if err != nil {
		return st, err
	}

	for _, r := range result {
		st = append(st, Details{
			ID:        r.ID,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}

	return st, nil
}

func One(id int, liveMode bool) (pair currency.Pair, err error) {
	query := postgres.CurrencyPairs(qm.Where("id=?", id))
	var r *postgres.CurrencyPair
	r, err = query.One(context.Background(), database.DB.SQL)
	if err != nil {
		return pair, err
	}

	if liveMode {
		pair, err = currency.NewPairFromString(r.KrakenSymbol)
	} else {
		pair, err = currency.NewPairFromString(r.GateioSymbol)
	}

	return pair, nil
}
