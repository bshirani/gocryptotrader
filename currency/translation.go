package currency

import (
	"strings"
)

// GetTranslation returns similar strings for a particular currency if not found
// returns the code back
func GetTranslation(currency Code) Code {
	val, ok := translations[currency]
	if !ok {
		return currency
	}
	return val
}

func GetPairTranslation(exch string, pair Pair) Pair {
	// fmt.Println("pair", pair.String(), pair.Upper().String())
	if strings.EqualFold(exch, "kraken") && strings.EqualFold(pair.Upper().String(), "BTC_USDT") {
		p, _ := NewPairFromString("XBT_USDT")
		return p
	}
	return pair
}

func ArePairsEqual(p1, p2 Pair) bool {
	return strings.EqualFold(p1.Quote.String(), p2.Quote.String()) && strings.EqualFold(p1.Base.String(), p2.Base.String())
}

var translations = map[Code]Code{
	BTC:  XBT,
	ETH:  XETH,
	DOGE: XDG,
	USD:  USDT,
	XBT:  BTC,
	XETH: ETH,
	XDG:  DOGE,
	USDT: USD,
}
