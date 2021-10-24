package currencypairstrategy

import (
	"context"
	"time"

	"gocryptotrader/currency"
	"gocryptotrader/database"
	"gocryptotrader/database/models/postgres"
	"gocryptotrader/database/repository/currencypair"
	"gocryptotrader/exchange/order"
)

type Details struct {
	ID           int
	CurrencyPair currency.Pair
	StrategyID   int
	Side         order.Side
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func All() (st []Details, err error) {
	query := postgres.CurrencyPairStrategies()
	var result []*postgres.CurrencyPairStrategy
	result, err = query.All(context.Background(), database.DB.SQL)
	if err != nil {
		return st, err
	}

	// // r.CurrencyPairID
	// // query2 := postgres.CurrencyPairs()
	// query2 := postgres.CurrencyPairs(qm.Where("id=?", id))
	// var cp *postgres.CurrencyPair
	// cp, err = query2.One(context.Background(), database.DB.SQL)
	// if err != nil {
	// 	return st, err
	// }
	// pair, _ := currency.NewPairFromString(cp.KrakenSymbol)
	// fmt.Println("loaded pair", cp.KrakenSymbol, pair, pair.Base, pair.Quote)

	for _, r := range result {
		pair, _ := currencypair.One(r.CurrencyPairID)
		// pair.Base, pair.Quote

		st = append(st, Details{
			ID:           r.ID,
			CurrencyPair: pair,
			Side:         order.Side(r.Side),
			StrategyID:   r.StrategyID,
			CreatedAt:    r.CreatedAt,
			UpdatedAt:    r.UpdatedAt,
		})
	}

	return st, nil
}
