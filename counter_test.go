package metricer

import (
	"testing"
)

// TestCounterReset validates reset counter metric
func TestCounterReset(t *testing.T) {
	metric := &counter{value: 50}
	metric.Reset()
	if metric.value != 0 {
		t.Errorf(" Expected: 0 but got %d", metric.value)
	}
}

// TestCounterCount validages count call for counter metric
func TestCounterCount(t *testing.T) {
	metric := &counter{value: 50}
	if val := metric.Count(); val != 50 {
		t.Errorf("Expected: 50 but got %d", val)
	}
}

// TestCounterInc validates incrementing conter call
func TestCounterInc(t *testing.T) {
	metric := &counter{value: 50}
	metric.Inc(20)
	if metric.value != 70 {
		t.Errorf("Expected: 70 but got %d", metric.value)
	}
}

// TestCounterDec validates decrementing conter call
func TestCounterDec(t *testing.T) {
	metric := &counter{value: 50}
	metric.Dec(20)
	if metric.value != 30 {
		t.Errorf("Expected: 30 but got %d", metric.value)
	}
}

// TestCounterName validates returning name for conter
func TestCounterName(t *testing.T) {
	metric := &counter{name: "counter-name"}
	if val := metric.Name(); val != "counter-name" {
		t.Errorf("Expected: counter-name but got %s", val)
	}
}

// TestCounterHelp validates returning help for counter
func TestCounterHelp(t *testing.T) {
	metric := &counter{name: "counter-name", help: "help"}
	if val := metric.Help(); val != "help" {
		t.Errorf("Expected: help but got %s", val)
	}
}
