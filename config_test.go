package metricer

import "testing"

func TestConfigBasic(t *testing.T) {
	cfg := &Config{}
	if err := cfg.Validate(); err != nil {
		t.Errorf("Expected: no errors, but got %s", err.Error())
	}

	if cfg.Port != defaultPort {
		t.Errorf("Expected: default port value, but got %d", cfg.Port)
	}
}

func TestConfigNil(t *testing.T) {
	var cfg *Config
	if err := cfg.Validate(); err == nil {
		t.Error("Expected: error, but got nil")
	}
}
