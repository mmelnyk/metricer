package metricer

import (
	"errors"
	"testing"
)

// TestHealthName validates returning name for health checker
func TestHealthName(t *testing.T) {
	metric := &health{name: "health-name"}
	if val := metric.Name(); val != "health-name" {
		t.Errorf("Expected: health-name but got %s", val)
	}
}

// TestHealthHelp validates returning help for health checker
func TestHealthHelp(t *testing.T) {
	metric := &health{name: "health-name", help: "help"}
	if val := metric.Help(); val != "help" {
		t.Errorf("Expected: help but got %s", val)
	}
}

// TestHealthCheckBasic validates returning help for health checker
func TestHealthCheckBasic(t *testing.T) {
	called := false
	metric := &health{name: "health-name", checker: func() error {
		called = true
		return nil
	}}
	if err := metric.Check(); err != nil {
		t.Errorf("Expected: no errors but got %s", err.Error())
	}
	if !called {
		t.Error("Expected: coroutine expected to be called, but it wasn't happened")
	}
}

// TestHealthCheckWithError validates returning help for health checker
func TestHealthCheckWithError(t *testing.T) {
	called := false
	metric := &health{name: "health-name", checker: func() error {
		called = true
		return errors.New("failed")
	}}
	err := metric.Check()
	if err == nil {
		t.Fatal("Expected: error expected but got no error")
	}
	if err.Error() != "failed" {
		t.Fatalf("Expected: error failed expected but got error %s", err.Error())
	}
	if !called {
		t.Error("Expected: coroutine expected to be called, but it wasn't happened")
	}
}

// TestHealthCheckWithPanic validates returning help for health checker
func TestHealthCheckWithPanic(t *testing.T) {
	var counternil Counter
	metric := &health{name: "health-name", checker: func() error {
		counternil.Name()
		return nil
	}}
	err := metric.Check()
	if err == nil {
		t.Fatal("Expected: error expected but got no error")
	}
	if err != errHealthCheckerPanic {
		t.Fatalf("Expected: error errHealthCheckerPanic expected but got error %s", err.Error())
	}
}
