package factors

import (
	"encoding/csv"
	"fmt"
	"gocryptotrader/database/repository/livetrade"
	"io"

	"github.com/shopspring/decimal"
)

func WriteCSV(w io.Writer, calcs []*Calculation, trades []*livetrade.Details) {
	cw := csv.NewWriter(w)

	headers := []string{"time", "id"}
	for _, h := range calcs[0].CSVHeader() {
		headers = append(headers, h)
	}
	headers = append(headers, "profit_loss_quote")
	headers = append(headers, "risked_quote")
	cw.Write(headers)

	for _, t := range trades {
		calc, _ := getCalcForTrade(calcs, t)
		strings := []string{fmt.Sprintf("%d", t.EntryTime.Unix()), fmt.Sprintf("%d", t.ID)}
		for _, s := range calc.ToStrings() {
			strings = append(strings, s)
		}
		strings = append(strings, t.ProfitLossQuote.String())
		strings = append(strings, decimal.NewFromFloat(t.RiskedQuote).String())
		cw.Write(strings)
	}

	cw.Flush()
}

func getCalcForTrade(calcs []*Calculation, t *livetrade.Details) (*Calculation, error) {
	for _, c := range calcs {
		if c.Time == t.EntryTime {
			return c, nil
		}
	}
	return nil, fmt.Errorf("calcuation not found for trade %v", t)
}
