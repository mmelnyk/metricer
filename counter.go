package metricer

import "sync/atomic"

type counter struct {
	name  string
	help  string
	value int64
}

func (metric *counter) Name() string {
	return metric.name
}

func (metric *counter) Help() string {
	return metric.help
}

func (metric *counter) Reset() {
	atomic.StoreInt64(&metric.value, 0)
}

func (metric *counter) Count() int64 {
	return atomic.LoadInt64(&metric.value)
}

func (metric *counter) Inc(v int64) {
	atomic.AddInt64(&metric.value, v)
}

func (metric *counter) Dec(v int64) {
	atomic.AddInt64(&metric.value, -v)
}
