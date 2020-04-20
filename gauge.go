package metricer

import "sync/atomic"

type gauge struct {
	name  string
	help  string
	value int64
}

func (metric *gauge) Name() string {
	return metric.name
}

func (metric *gauge) Help() string {
	return metric.help
}

func (metric *gauge) Update(v int64) {
	atomic.StoreInt64(&metric.value, v)
}

func (metric *gauge) Value() int64 {
	return atomic.LoadInt64(&metric.value)
}
