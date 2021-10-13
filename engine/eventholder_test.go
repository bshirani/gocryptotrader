package engine

import (
	"testing"

	"gtc/eventtypes"
	"gtc/eventtypes/order"
)

func TestReset(t *testing.T) {
	t.Parallel()
	e := Holder{Queue: []eventtypes.EventHandler{}}
	e.Reset()
	if e.Queue != nil {
		t.Error("expected nil")
	}
}

func TestAppendEvent(t *testing.T) {
	t.Parallel()
	e := Holder{Queue: []eventtypes.EventHandler{}}
	e.AppendEvent(&order.Order{})
	if len(e.Queue) != 1 {
		t.Error("expected 1")
	}
}

func TestNextEvent(t *testing.T) {
	t.Parallel()
	e := Holder{Queue: []eventtypes.EventHandler{}}
	ev := e.NextEvent()
	if ev != nil {
		t.Error("expected not ok")
	}

	e = Holder{Queue: []eventtypes.EventHandler{
		&order.Order{},
		&order.Order{},
		&order.Order{},
	}}
	if len(e.Queue) != 3 {
		t.Error("expected 3")
	}
	e.NextEvent()
	if len(e.Queue) != 2 {
		t.Error("expected 2")
	}
}
