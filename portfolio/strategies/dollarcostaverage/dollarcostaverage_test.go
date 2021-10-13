package dollarcostaverage

import (
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"gtc/common"
	"gtc/data"
	"gtc/data/kline"
	"gtc/strategies/base"
	"gtc/eventtypes/event"
	eventkline "gtc/eventtypes/kline"
	"gtc/eventtypes/signal"
	"gtc/currency"
	"gtc/exchanges/asset"
	gctkline "gtc/exchanges/kline"
	gctorder "gtc/exchanges/order"
)

func TestName(t *testing.T) {
	d := Strategy{}
	n := d.Name()
	if n != Name {
		t.Errorf("expected %v", Name)
	}
}

func TestSupportsSimultaneousProcessing(t *testing.T) {
	s := Strategy{}
	if !s.SupportsSimultaneousProcessing() {
		t.Error("expected true")
	}
}

func TestSetCustomSettings(t *testing.T) {
	s := Strategy{}
	err := s.SetCustomSettings(nil)
	if !errors.Is(err, base.ErrCustomSettingsUnsupported) {
		t.Errorf("received: %v, expected: %v", err, base.ErrCustomSettingsUnsupported)
	}
}

func TestOnSignal(t *testing.T) {
	s := Strategy{}
	_, err := s.OnSignal(nil)
	if !errors.Is(err, eventtypes.ErrNilEvent) {
		t.Errorf("received: %v, expected: %v", err, eventtypes.ErrNilEvent)
	}

	dStart := time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC)
	dInsert := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	dEnd := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	exch := "binance"
	a := asset.Spot
	p := currency.NewPair(currency.BTC, currency.USDT)
	d := data.Base{}
	d.SetStream([]eventtypes.DataEventHandler{&eventkline.Kline{
		Base: event.Base{
			Exchange:     exch,
			Time:         dInsert,
			Interval:     gctkline.OneDay,
			CurrencyPair: p,
			AssetType:    a,
		},
		Open:   decimal.NewFromInt(1337),
		Close:  decimal.NewFromInt(1337),
		Low:    decimal.NewFromInt(1337),
		High:   decimal.NewFromInt(1337),
		Volume: decimal.NewFromInt(1337),
	}})
	d.Next()
	da := &kline.DataFromKline{
		Item:        gctkline.Item{},
		Base:        d,
		RangeHolder: &gctkline.IntervalRangeHolder{},
	}
	var resp signal.Event
	resp, err = s.OnSignal(da)
	if err != nil {
		t.Error(err)
	}
	if resp.GetDirection() != common.MissingData {
		t.Error("expected missing data")
	}

	da.Item = gctkline.Item{
		Exchange: exch,
		Pair:     p,
		Asset:    a,
		Interval: gctkline.OneDay,
		Candles: []gctkline.Candle{
			{
				Time:   dInsert,
				Open:   1337,
				High:   1337,
				Low:    1337,
				Close:  1337,
				Volume: 1337,
			},
		},
	}
	err = da.Load()
	if err != nil {
		t.Error(err)
	}

	ranger, err := gctkline.CalculateCandleDateRanges(dStart, dEnd, gctkline.OneDay, 100000)
	if err != nil {
		t.Error(err)
	}
	da.RangeHolder = ranger
	da.RangeHolder.SetHasDataFromCandles(da.Item.Candles)
	resp, err = s.OnSignal(da, nil)
	if err != nil {
		t.Error(err)
	}
	if resp.GetDirection() != gctorder.Buy {
		t.Errorf("expected buy, received %v", resp.GetDirection())
	}
}

func TestOnSignals(t *testing.T) {
	s := Strategy{}
	_, err := s.OnSignal(nil)
	if !errors.Is(err, eventtypes.ErrNilEvent) {
		t.Errorf("received: %v, expected: %v", err, eventtypes.ErrNilEvent)
	}
	dStart := time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC)
	dInsert := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	dEnd := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	exch := "binance"
	a := asset.Spot
	p := currency.NewPair(currency.BTC, currency.USDT)
	d := data.Base{}
	d.SetStream([]eventtypes.DataEventHandler{&eventkline.Kline{
		Base: event.Base{
			Offset:       1,
			Exchange:     exch,
			Time:         dInsert,
			Interval:     gctkline.OneDay,
			CurrencyPair: p,
			AssetType:    a,
		},
		Open:   decimal.NewFromInt(1337),
		Close:  decimal.NewFromInt(1337),
		Low:    decimal.NewFromInt(1337),
		High:   decimal.NewFromInt(1337),
		Volume: decimal.NewFromInt(1337),
	}})
	d.Next()
	da := &kline.DataFromKline{
		Item:        gctkline.Item{},
		Base:        d,
		RangeHolder: &gctkline.IntervalRangeHolder{},
	}
	var resp []signal.Event
	resp, err = s.OnSimultaneousSignals([]data.Handler{da}, nil)
	if err != nil {
		t.Error(err)
	}
	if len(resp) != 1 {
		t.Fatal("expected 1 response")
	}
	if resp[0].GetDirection() != common.MissingData {
		t.Error("expected missing data")
	}

	da.Item = gctkline.Item{
		Exchange: exch,
		Pair:     p,
		Asset:    a,
		Interval: gctkline.OneDay,
		Candles: []gctkline.Candle{
			{
				Time:   dInsert,
				Open:   1337,
				High:   1337,
				Low:    1337,
				Close:  1337,
				Volume: 1337,
			},
		},
	}
	err = da.Load()
	if err != nil {
		t.Error(err)
	}

	ranger, err := gctkline.CalculateCandleDateRanges(dStart, dEnd, gctkline.OneDay, 100000)
	if err != nil {
		t.Error(err)
	}
	da.RangeHolder = ranger
	da.RangeHolder.SetHasDataFromCandles(da.Item.Candles)
	resp, err = s.OnSimultaneousSignals([]data.Handler{da}, nil)
	if err != nil {
		t.Error(err)
	}
	if len(resp) != 1 {
		t.Fatal("expected 1 response")
	}
	if resp[0].GetDirection() != gctorder.Buy {
		t.Error("expected buy")
	}
}

func TestSetDefaults(t *testing.T) {
	s := Strategy{}
	s.SetDefaults()
	if s != (Strategy{}) {
		t.Error("expected no changes")
	}
}
