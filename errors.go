package metricer

import "errors"

var (
	errNilConfig          = errors.New("Config cannot be nil")
	errRunMoreOnce        = errors.New("Metricer start function is called more than once")
	errHealthCheckerPanic = errors.New("Panic in health check callback")
)
