package metricer

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"

	"go.melnyk.org/mlog"
	"go.melnyk.org/mlog/nolog"
)

type host struct {
	started time.Time // time of structure initialization

	config Config

	logbook mlog.Logbook
	log     mlog.Logger

	server *http.Server
	wg     sync.WaitGroup

	mu           sync.RWMutex
	metrics      []interface{}
	healthchecks []Health

	rtgoroutines       Gauge
	rtmemalloc         Gauge
	failedhealthchecks Counter
}

// NewHost creates new instance of metricer host with basic initialization
func NewHost(config *Config, lb mlog.Logbook) Host {
	h := &host{}
	h.metrics = make([]interface{}, 0, 16)
	h.healthchecks = make([]Health, 0, 8)
	h.started = time.Now()

	// initialize logger and keep logbook for exposing api
	if lb == nil {
		lb = nolog.NewLogbook()
	}
	h.logbook = lb
	h.log = lb.Joiner().Join(logname)

	if err := config.Validate(); err != nil {
		h.log.Event(mlog.Warning, func(e mlog.Event) {
			e.String("msg", "Config validation problem")
			e.String("error", err.Error())
		})
		config = &Config{
			Port: defaultPort,
		}
	}
	h.config = *config

	// create runtime metrics
	rtos := h.NewLabel(rtmetricos, rtmetricoshelp)
	rtos.Update(runtime.GOOS)
	rtnumcpu := h.NewGauge(rtmetricnumcpu, rtmetricnumcpuhelp)
	rtnumcpu.Update(int64(runtime.NumCPU()))

	h.rtgoroutines = h.NewGauge(rtmetricgoroutines, rtmetricgoroutineshelp)
	h.rtmemalloc = h.NewGauge(rtmetricmemalloc, rtmetricmemallochelp)
	h.failedhealthchecks = h.NewCounter(healthcheckfailed, healthcheckfailedhelp)

	return h
}

// NewLabel creates new named label metric inside metrics collection
func (h *host) NewLabel(name string, help string) Label {
	metric := &label{name: name, help: help}
	h.mu.Lock()
	h.metrics = append(h.metrics, metric)
	h.mu.Unlock()
	return metric
}

// NewCounter creates new named counter metric inside metrics collection
func (h *host) NewCounter(name string, help string) Counter {
	metric := &counter{name: name, help: help}
	h.mu.Lock()
	h.metrics = append(h.metrics, metric)
	h.mu.Unlock()
	return metric
}

// NewGauge creates new named gauge metric inside metrics collection
func (h *host) NewGauge(name string, help string) Gauge {
	metric := &gauge{name: name, help: help}
	h.mu.Lock()
	h.metrics = append(h.metrics, metric)
	h.mu.Unlock()
	return metric
}

// NewHealthCheck creates new named health checker
func (h *host) NewHealthCheck(name string, help string, checker HealthcheckFunc) {
	metric := &health{name: name, help: help, checker: checker}
	h.mu.Lock()
	h.healthchecks = append(h.healthchecks, metric)
	h.mu.Unlock()
}

func (h *host) Start() error {
	if h.server != nil {
		h.log.Warning("Metricer start function is called more than once")
		return errRunMoreOnce
	}

	h.log.Info("Starting Metricer...")

	h.wg.Add(1)
	wait := make(chan struct{}, 1)

	go func(lh *host) {
		defer lh.wg.Done()

		addrtmpl := "127.0.0.1:%d"
		if lh.config.AllowExternal {
			addrtmpl = ":%d"
		}

		lh.server = &http.Server{Handler: lh.buildMuxer()}

	loop:
		for i := uint(0); i < 24; i++ {
			port := fmt.Sprintf(addrtmpl, lh.config.Port+i)

			ln, err := net.Listen("tcp", port)
			if err != nil {
				lh.log.Event(mlog.Verbose, func(e mlog.Event) {
					e.String("msg", "Moving to next try during net.Listen error")
					e.String("err", err.Error())
				})
				continue
			}

			lh.server.Addr = port

			// continue
			wait <- struct{}{}

			lh.log.Info("Metricer interface should be available at " + port)
			if lh.config.EnableDebug {
				lh.log.Verbose("HTTP interface for debugging is active")
			}

			switch err := lh.server.Serve(ln); {
			case err == http.ErrServerClosed:
				break loop
			case err != nil:
				lh.log.Warning(err.Error())
			}
		}

		// check if server was initialized correctly
		if lh.server.Addr == "" {
			// ups...
			wait <- struct{}{}
		}

		lh.log.Verbose("Metricer interface is not available")
	}(h)

	<-wait

	return nil
}

func (h *host) Stop() error {
	h.log.Info("Stopping Merticer...")

	if h.server != nil {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		h.server.Shutdown(ctx)
	}

	h.wg.Wait()

	h.log.Info("Metricer has been stopped")

	return nil
}
