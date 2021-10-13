package strategies

import (
	"errors"
	"testing"

	"gocryptotradere/portfolio/strategies/base"
	"gocryptotradere/portfolio/strategies/dollarcostaverage"
	"gocryptotradere/portfolio/strategies/rsi"
)

func TestGetStrategies(t *testing.T) {
	t.Parallel()
	resp := GetStrategies()
	if len(resp) < 2 {
		t.Error("expected at least 2 strategies to be loaded")
	}
}

func TestLoadStrategyByName(t *testing.T) {
	t.Parallel()
	var resp Handler
	_, err := LoadStrategyByName("test", false)
	if !errors.Is(err, base.ErrStrategyNotFound) {
		t.Errorf("received: %v, expected: %v", err, base.ErrStrategyNotFound)
	}
	_, err = LoadStrategyByName("test", true)
	if !errors.Is(err, base.ErrStrategyNotFound) {
		t.Errorf("received: %v, expected: %v", err, base.ErrStrategyNotFound)
	}

	resp, err = LoadStrategyByName(dollarcostaverage.Name, false)
	if err != nil {
		t.Error(err)
	}
	if resp.Name() != dollarcostaverage.Name {
		t.Error("expected dca")
	}
	resp, err = LoadStrategyByName(dollarcostaverage.Name, true)
	if err != nil {
		t.Error(err)
	}
	if !resp.UsingSimultaneousProcessing() {
		t.Error("expected true")
	}

	resp, err = LoadStrategyByName(rsi.Name, false)
	if err != nil {
		t.Error(err)
	}
	if resp.Name() != rsi.Name {
		t.Error("expected rsi")
	}
	_, err = LoadStrategyByName(rsi.Name, true)
	if !errors.Is(err, nil) {
		t.Errorf("received: %v, expected: %v", err, nil)
	}
}
