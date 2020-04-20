package metricer

import (
	"testing"
)

// TestLabelUpdate validates update call for gauge metric
func TestLabelUpdate(t *testing.T) {
	metric := &label{}
	metric.Update("new value")
	if val := metric.Value(); val != "new value" {
		t.Errorf(" Expected: new value but got %s", val)
	}
}

// TestLabelValue validages value call for gauge metric
func TestLabelValue(t *testing.T) {
	metric := &label{}
	if val := metric.Value(); val != "(n/a)" {
		t.Errorf("Expected: n/a but got %s", val)
	}
}

// TestLabelName validates returning name for conter
func TestLabelName(t *testing.T) {
	metric := &label{name: "label-name"}
	if val := metric.Name(); val != "label-name" {
		t.Errorf("Expected: label-name but got %s", val)
	}
}

// TestLabelHelp validates returning help for counter
func TestLabelHelp(t *testing.T) {
	metric := &label{name: "label-name", help: "help"}
	if val := metric.Help(); val != "help" {
		t.Errorf("Expected: help but got %s", val)
	}
}
