package currencypair

import (
	"context"
	"fmt"
	"time"

	"gocryptotrader/database"
	"gocryptotrader/database/models/postgres"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Details struct {
	ID           int
	KrakenSymbol string
	CreatedAt    time.Time
	UpdatedAt    time.Time
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

func One(id int) (st Details, err error) {
	// query := postgres.Datahistoryjobs(qm.Where("nickname = ?", strings.ToLower(nickname)))
	query := postgres.Currencies(qm.Where("id=?", id))
	var base *postgres.Currency
	var quote *postgres.Currency
	base, err = query.One(context.Background(), database.DB.SQL)
	quote, err = query.One(context.Background(), database.DB.SQL)

	fmt.Println("base", base, "quote", quote)

	// var r *postgres.CurrencyPair
	// r, err = query.One(context.Background(), database.DB.SQL)
	// if err != nil {
	// 	return st, err
	// }
	//
	// st = Details{
	// 	ID:        r.ID,
	// 	Base
	// 	CreatedAt: r.CreatedAt,
	// 	UpdatedAt: r.UpdatedAt,
	// }
	//
	return st, nil
}
