package metricer

// Metric provides interface to general metrics
type Metric interface {
	Name() string
	Help() string
}

// HealthcheckFunc provides callback function to check current health status
type HealthcheckFunc func() error

// Health provides interface to update health status
type Health interface {
	Metric
	Check() error
}

// Label provides interface to label metrics
type Label interface {
	Metric
	Update(string)
	Value() string
}

// Counter provides interface to metrics with increment and decrement methods
type Counter interface {
	Metric
	Reset()
	Inc(int64)
	Dec(int64)
	Count() int64
}

// Gauge provides interface to metrics with update methods
type Gauge interface {
	Metric
	Update(int64)
	Value() int64
}

// Host represents metric host interace
type Host interface {
	Start() error
	Stop() error

	NewHealthCheck(string, string, HealthcheckFunc)

	NewLabel(string, string) Label
	NewGauge(string, string) Gauge
	NewCounter(string, string) Counter
}
