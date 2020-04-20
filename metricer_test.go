package metricer

import (
	"testing"

	"go.melnyk.org/mlog/testlog"
)

func TestMetricerBasicNoConfig(t *testing.T) {
	host := NewHost(nil, nil)
	if host == nil {
		t.Fatal("Expected: Metricer host, but got nil")
	}
}

func TestMetricerBasicConfig(t *testing.T) {
	cfg := &Config{}
	host := NewHost(cfg, nil)
	if host == nil {
		t.Fatal("Expected: Metricer host, but got nil")
	}
}

func TestMetricerBasic(t *testing.T) {
	mhost := NewHost(nil, testlog.NewLogbook())
	if mhost == nil {
		t.Fatal("Expected: Metricer host, but got nil")
	}
	counter := mhost.NewCounter("counter", "counter help")
	if counter == nil {
		t.Fatal("Expected: Counter, but got nil")
	}
	gauge := mhost.NewGauge("gauge", "gauge help")
	if gauge == nil {
		t.Fatal("Expected: Gauge, but got nil")
	}
	label := mhost.NewLabel("label", "label help")
	if label == nil {
		t.Fatal("Expected: Label, but got nil")
	}
	mhost.NewHealthCheck("health", "health help", func() error {
		return nil
	})
	if len(mhost.(*host).healthchecks) != 1 {
		t.Fatalf("Expected: size of healthcheck list equal 1, but size %d", len(mhost.(*host).healthchecks))
	}
}

func TestMetricerStop(t *testing.T) {
	mhost := NewHost(nil, testlog.NewLogbook())
	// TBD: should we have an error on stop?
	mhost.Stop()
}

func TestMetricerStartFailed(t *testing.T) {
	cfg := &Config{Port: 1000000, AllowExternal: true}
	mhost := NewHost(cfg, testlog.NewLogbook())
	// TBD: should we have an error on stop?
	mhost.Start()
	// should fail?
	mhost.Start()
	// should fail?
	mhost.Stop()
}

func TestMetricerStartStop(t *testing.T) {
	cfg := &Config{EnableDebug: true}
	mhost := NewHost(cfg, testlog.NewLogbook())
	// TBD: should we have an error on stop?
	mhost.Start()
	mhost.Stop()
}
