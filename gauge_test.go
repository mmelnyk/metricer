package metricer

import (
	"testing"
)

// TestGaugeUpdate validates update call for gauge metric
func TestGaugeUpdate(t *testing.T) {
	metric := &gauge{value: 50}
	metric.Update((int64)(88))
	if metric.value != 88 {
		t.Errorf(" Expected: 88 but got %d", metric.value)
	}
}

// TestGaugeValue validages value call for gauge metric
func TestGaugeValue(t *testing.T) {
	metric := &gauge{value: 50}
	if val := metric.Value(); val != 50 {
		t.Errorf("Expected: 50 but got %d", val)
	}
}

// TestGaugeName validates returning name for conter
func TestGaugeName(t *testing.T) {
	metric := &gauge{name: "counter-name"}
	if val := metric.Name(); val != "counter-name" {
		t.Errorf("Expected: counter-name but got %s", val)
	}
}

// TestGaugeHelp validates returning help for counter
func TestGaugeHelp(t *testing.T) {
	metric := &gauge{name: "counter-name", help: "help"}
	if val := metric.Help(); val != "help" {
		t.Errorf("Expected: help but got %s", val)
	}
}
