package currencypairstrategy

import (
	"context"
	"fmt"
	"time"

	"gocryptotrader/currency"
	"gocryptotrader/database"
	"gocryptotrader/database/models/postgres"
	"gocryptotrader/database/repository/currencypair"
	"gocryptotrader/exchange/order"

	"github.com/shopspring/decimal"
)

type Details struct {
	ID           int
	CurrencyPair currency.Pair
	StrategyID   int
	Active       bool
	Side         order.Side
	Weight       decimal.Decimal
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func ActivePairs(liveMode bool) (currency.Pairs, error) {
	if liveMode {
		pair, err := currency.NewPairFromString("XBT_USDT")
		return currency.Pairs{pair}, err
	}

	pair, err := currency.NewPairFromString("BTC_USDT")
	return currency.Pairs{pair}, err

	// return pairs.Format(pairFormat.Delimiter,
	// 		pairFormat.Index,
	// 		pairFormat.Uppercase),
	// 	nil
	// return
}

func All(liveMode bool) (st []Details, err error) {
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
		pair, err := currencypair.One(r.CurrencyPairID, liveMode)
		if err != nil {
			fmt.Println("error getting cp", err)
		}
		fmt.Println("got pair", pair, "for id", r.CurrencyPairID)
		// pair.Base, pair.Quote

		st = append(st, Details{
			ID:           r.ID,
			CurrencyPair: pair,
			Active:       r.Active,
			Side:         order.Side(r.Side),
			Weight:       decimal.NewFromFloat(r.Weight),
			StrategyID:   r.StrategyID,
			CreatedAt:    r.CreatedAt,
			UpdatedAt:    r.UpdatedAt,
		})
	}

	return st, nil
}
