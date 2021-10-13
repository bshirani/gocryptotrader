package base

import (
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"gocryptotrader/common"
	"gocryptotrader/data"
	datakline "gocryptotrader/data/kline"
	"gocryptotrader/eventtypes/event"
	"gocryptotrader/eventtypes/kline"
	"gocryptotrader/currency"
	"gocryptotrader/exchange/asset"
	gctkline "gocryptotrader/exchange/kline"
)

func TestGetBase(t *testing.T) {
	s := Strategy{}
	_, err := s.GetBaseData(nil)
	if !errors.Is(err, eventtypes.ErrNilArguments) {
		t.Errorf("received: %v, expected: %v", err, eventtypes.ErrNilArguments)
	}

	_, err = s.GetBaseData(&datakline.DataFromKline{})
	if !errors.Is(err, eventtypes.ErrNilEvent) {
		t.Errorf("received: %v, expected: %v", err, eventtypes.ErrNilEvent)
	}
	tt := time.Now()
	exch := "binance"
	a := asset.Spot
	p := currency.NewPair(currency.BTC, currency.USDT)
	d := data.Base{}
	d.SetStream([]eventtypes.DataEventHandler{&kline.Kline{
		Base: event.Base{
			Exchange:     exch,
			Time:         tt,
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
	_, err = s.GetBaseData(&datakline.DataFromKline{
		Item:        gctkline.Item{},
		Base:        d,
		RangeHolder: &gctkline.IntervalRangeHolder{},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestSetSimultaneousProcessing(t *testing.T) {
	s := Strategy{}
	is := s.UsingSimultaneousProcessing()
	if is {
		t.Error("expected false")
	}
	s.SetSimultaneousProcessing(true)
	is = s.UsingSimultaneousProcessing()
	if !is {
		t.Error("expected true")
	}
}
