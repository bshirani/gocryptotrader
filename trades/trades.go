package trades

import (
	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/eventtypes"
	"github.com/thrasher-corp/gocryptotrader/eventtypes/signal"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
)

// Create makes a Trade struct to track total values of strategy holdings over the course of a backtesting run
func Create(ev signal.Event) (Trade, error) {
	if ev == nil {
		return Trade{}, eventtypes.ErrNilEvent
	}

	return Trade{
		Offset:     ev.GetOffset(),
		Pair:       ev.Pair(),
		Asset:      ev.GetAssetType(),
		Exchange:   ev.GetExchange(),
		Timestamp:  ev.GetTime(),
		Direction:  ev.GetDirection(),
		EntryPrice: ev.GetPrice(),
	}, nil
}

// Update calculates holding statistics for the events time
func (t *Trade) Update(e signal.Event) {
	t.Timestamp = e.GetTime()
	t.Offset = e.GetOffset()
	t.update(e)
}

// UpdateValue calculates the trades's value for a data event's time and price
func (t *Trade) UpdateValue(d eventtypes.DataEventHandler) {
	t.Timestamp = d.GetTime()
	t.CurrentPrice = d.ClosePrice()

	if t.Direction == order.Buy {
		t.NetProfit = t.CurrentPrice.Sub(t.EntryPrice)
	} else if t.Direction == order.Sell {
		t.NetProfit = t.EntryPrice.Sub(t.CurrentPrice)
	}
	// latest := d.ClosePrice()
	t.Offset = d.GetOffset()
	// t.updateValue(latest)
}

func (t *Trade) update(e signal.Event) {
	// fmt.Println("update trade")
	// direction := e.GetDirection()
	// o := e.GetOrder()
	// if o != nil {
	// 	amount := decimal.NewFromFloat(o.Amount)
	// 	fee := decimal.NewFromFloat(o.Fee)
	// 	price := decimal.NewFromFloat(o.Price)
	// 	t.BaseSize = decimal.NewFromFloat(1000.0)  //f.BaseAvailable()
	// 	t.QuoteSize = decimal.NewFromFloat(1000.0) //f.QuoteAvailable()
	// 	t.BaseValue = t.BaseSize.Mul(price)
	// 	t.TotalFees = t.TotalFees.Add(fee)
	// 	switch direction {
	// 	case order.Buy:
	// 		t.BoughtAmount = t.BoughtAmount.Add(amount)
	// 		t.BoughtValue = t.BoughtAmount.Mul(price)
	// 	case order.Sell:
	// 		t.SoldAmount = t.SoldAmount.Add(amount)
	// 		t.SoldValue = t.SoldAmount.Mul(price)
	// 	case eventtypes.DoNothing, eventtypes.CouldNotSell, eventtypes.CouldNotBuy, eventtypes.MissingData, eventtypes.TransferredFunds, "":
	// 	}
	// }
	// t.TotalValueLostToVolumeSizing = t.TotalValueLostToVolumeSizing.Add(e.GetClosePrice().Sub(e.GetVolumeAdjustedPrice()).Mul(e.GetAmount()))
	// t.TotalValueLostToSlippage = t.TotalValueLostToSlippage.Add(e.GetVolumeAdjustedPrice().Sub(e.GetPurchasePrice()).Mul(e.GetAmount()))
	// t.updateValue(e.GetClosePrice())
	return
}

func (t *Trade) updateValue(latestPrice decimal.Decimal) {
	// fmt.Println("updatevalue")
	// 	origPosValue := t.BaseValue
	// 	origBoughtValue := t.BoughtValue
	// 	origSoldValue := t.SoldValue
	// 	origTotalValue := t.TotalValue
	t.CurrentPrice = latestPrice
	// 	t.BoughtValue = t.BoughtAmount.Mul(latestPrice)
	// 	t.SoldValue = t.SoldAmount.Mul(latestPrice)
	// 	t.TotalValue = t.BaseValue.Add(t.QuoteSize)
	//
	// 	t.TotalValueDifference = t.TotalValue.Sub(origTotalValue)
	// 	t.BoughtValueDifference = t.BoughtValue.Sub(origBoughtValue)
	// 	t.PositionsValueDifference = t.BaseValue.Sub(origPosValue)
	// 	t.SoldValueDifference = t.SoldValue.Sub(origSoldValue)
	//
	// 	if !origTotalValue.IsZero() {
	// 		t.ChangeInTotalValuePercent = t.TotalValue.Sub(origTotalValue).Div(origTotalValue)
	// 	}
}
