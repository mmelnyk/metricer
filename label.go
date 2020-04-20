package metricer

import "sync/atomic"

type label struct {
	name  string
	help  string
	value atomic.Value
}

func (metric *label) Name() string {
	return metric.name
}

func (metric *label) Help() string {
	return metric.help
}

func (metric *label) Update(v string) {
	metric.value.Store(v)
}

func (metric *label) Value() string {
	v := metric.value.Load()
	if str, ok := v.(string); ok {
		return str
	}
	return "(n/a)"
}
