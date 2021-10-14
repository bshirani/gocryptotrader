package fill

import (
	"testing"

	gctorder "gocryptotrader/exchange/order"

	"github.com/shopspring/decimal"
)

func TestSetDirection(t *testing.T) {
	t.Parallel()
	f := Fill{
		Direction: gctorder.Sell,
	}
	f.SetDirection(gctorder.Buy)
	if f.Direction() != gctorder.Buy {
		t.Error("expected buy")
	}
}

func TestSetAmount(t *testing.T) {
	t.Parallel()
	f := Fill{
		Amount: decimal.NewFromInt(1),
	}
	f.SetAmount(decimal.NewFromInt(1337))
	if !f.GetAmount().Equal(decimal.NewFromInt(1337)) {
		t.Error("expected decimal.NewFromInt(1337)")
	}
}

func TestGetClosePrice(t *testing.T) {
	t.Parallel()
	f := Fill{
		ClosePrice: decimal.NewFromInt(1337),
	}
	if !f.GetClosePrice().Equal(decimal.NewFromInt(1337)) {
		t.Error("expected decimal.NewFromInt(1337)")
	}
}

func TestGetVolumeAdjustedPrice(t *testing.T) {
	t.Parallel()
	f := Fill{
		VolumeAdjustedPrice: decimal.NewFromInt(1337),
	}
	if !f.GetVolumeAdjustedPrice().Equal(decimal.NewFromInt(1337)) {
		t.Error("expected decimal.NewFromInt(1337)")
	}
}

func TestGetPurchasePrice(t *testing.T) {
	t.Parallel()
	f := Fill{
		PurchasePrice: decimal.NewFromInt(1337),
	}
	if !f.GetPurchasePrice().Equal(decimal.NewFromInt(1337)) {
		t.Error("expected decimal.NewFromInt(1337)")
	}
}

func TestSetExchangeFee(t *testing.T) {
	t.Parallel()
	f := Fill{
		ExchangeFee: decimal.NewFromInt(1),
	}
	f.SetExchangeFee(decimal.NewFromInt(1337))
	if !f.GetExchangeFee().Equal(decimal.NewFromInt(1337)) {
		t.Error("expected decimal.NewFromInt(1337)")
	}
}

func TestGetOrder(t *testing.T) {
	t.Parallel()
	f := Fill{
		Order: &gctorder.Detail{},
	}
	if f.GetOrder() == nil {
		t.Error("expected not nil")
	}
}

func TestGetSlippageRate(t *testing.T) {
	t.Parallel()
	f := Fill{
		Slippage: decimal.NewFromInt(1),
	}
	if !f.GetSlippageRate().Equal(decimal.NewFromInt(1)) {
		t.Error("expected 1")
	}
}
