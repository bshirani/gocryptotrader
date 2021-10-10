package holdings

import (
	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/backtester/common"
	"github.com/thrasher-corp/gocryptotrader/backtester/eventtypes/fill"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
)

// Create makes a Holding struct to track total values of strategy holdings over the course of a backtesting run
func Create(ev common.EventHandler, initialFunds, riskFreeRate decimal.Decimal) (Holding, error) {
	if ev == nil {
		return Holding{}, common.ErrNilEvent
	}

	if initialFunds.LessThan(decimal.NewFromFloat(0)) {
		return Holding{}, ErrInitialFundsZero
	}

	return Holding{
		Offset:            ev.GetOffset(),
		Pair:              ev.Pair(),
		Asset:             ev.GetAssetType(),
		Exchange:          ev.GetExchange(),
		Timestamp:         ev.GetTime(),
		QuoteInitialFunds: initialFunds,
		QuoteSize:         initialFunds,
		BaseInitialFunds:  initialFunds,
		BaseSize:          initialFunds,
		RiskFreeRate:      riskFreeRate,
		TotalInitialValue: initialFunds,
	}, nil
}

// Update calculates holding statistics for the events time
func (h *Holding) Update(e fill.Event) {
	h.Timestamp = e.GetTime()
	h.Offset = e.GetOffset()
	h.update(e)
}

// UpdateValue calculates the holding's value for a data event's time and price
func (h *Holding) UpdateValue(d common.DataEventHandler) {
	h.Timestamp = d.GetTime()
	latest := d.ClosePrice()
	h.Offset = d.GetOffset()
	h.updateValue(latest)
}

// HasInvestments determines whether there are any holdings in the base funds
func (h *Holding) HasInvestments() bool {
	return h.BaseSize.GreaterThan(decimal.Zero)
}

// HasFunds determines whether there are any holdings in the quote funds
func (h *Holding) HasFunds() bool {
	return h.QuoteSize.GreaterThan(decimal.Zero)
}

func (h *Holding) update(e fill.Event) {
	direction := e.GetDirection()
	o := e.GetOrder()
	if o != nil {
		amount := decimal.NewFromFloat(o.Amount)
		fee := decimal.NewFromFloat(o.Fee)
		price := decimal.NewFromFloat(o.Price)
		h.BaseSize = decimal.NewFromFloat(1000.0)  //f.BaseAvailable()
		h.QuoteSize = decimal.NewFromFloat(1000.0) //f.QuoteAvailable()
		h.BaseValue = h.BaseSize.Mul(price)
		h.TotalFees = h.TotalFees.Add(fee)
		switch direction {
		case order.Buy:
			h.BoughtAmount = h.BoughtAmount.Add(amount)
			h.BoughtValue = h.BoughtAmount.Mul(price)
		case order.Sell:
			h.SoldAmount = h.SoldAmount.Add(amount)
			h.SoldValue = h.SoldAmount.Mul(price)
		case common.DoNothing, common.CouldNotSell, common.CouldNotBuy, common.MissingData, common.TransferredFunds, "":
		}
	}
	h.TotalValueLostToVolumeSizing = h.TotalValueLostToVolumeSizing.Add(e.GetClosePrice().Sub(e.GetVolumeAdjustedPrice()).Mul(e.GetAmount()))
	h.TotalValueLostToSlippage = h.TotalValueLostToSlippage.Add(e.GetVolumeAdjustedPrice().Sub(e.GetPurchasePrice()).Mul(e.GetAmount()))
	h.updateValue(e.GetClosePrice())
}

func (h *Holding) updateValue(latestPrice decimal.Decimal) {
	origPosValue := h.BaseValue
	origBoughtValue := h.BoughtValue
	origSoldValue := h.SoldValue
	origTotalValue := h.TotalValue
	h.BaseValue = h.BaseSize.Mul(latestPrice)
	h.BoughtValue = h.BoughtAmount.Mul(latestPrice)
	h.SoldValue = h.SoldAmount.Mul(latestPrice)
	h.TotalValue = h.BaseValue.Add(h.QuoteSize)

	h.TotalValueDifference = h.TotalValue.Sub(origTotalValue)
	h.BoughtValueDifference = h.BoughtValue.Sub(origBoughtValue)
	h.PositionsValueDifference = h.BaseValue.Sub(origPosValue)
	h.SoldValueDifference = h.SoldValue.Sub(origSoldValue)

	if !origTotalValue.IsZero() {
		h.ChangeInTotalValuePercent = h.TotalValue.Sub(origTotalValue).Div(origTotalValue)
	}
}
