package metricer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/pprof"
	"runtime"
	"strings"
	"time"

	"go.melnyk.org/mlog"
)

const (
	pathHealthCheck   = "/health/check"
	pathMetricsValues = "/metrics/values"

	pathDebug             = "/debug/"
	pathDebugPprof        = pathDebug + "pprof/"
	pathDebugPprofCmdline = pathDebugPprof + "cmdline"
	pathDebugPprofProfile = pathDebugPprof + "profile"
	pathDebugPprofSymbol  = pathDebugPprof + "symbol"
	pathDebugPprofTrace   = pathDebugPprof + "trace"
	pathDebugLoggerLevels = pathDebug + "logger/levels"
)

func (h *host) handlerError(w http.ResponseWriter, r *http.Request, status int, message string) {
	h.log.Event(mlog.Warning, func(e mlog.Event) {
		e.String("msg", message)
		e.Int("status", status)
		e.String("remote", r.RemoteAddr)
		e.String("method", r.Method)
		e.String("path", r.RequestURI)
	})

	e := struct {
		Error interface{} `json:"error"`
	}{
		Error: struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{
			Code:    status,
			Message: message,
		},
	}

	w.Header().Set("Content-Type", contenttypeJSON)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(e)
}

func (h *host) healthCheck(w http.ResponseWriter, r *http.Request) {
	h.wg.Add(1)
	defer h.wg.Done()

	// only GET method is allowed
	if r.Method != http.MethodGet {
		h.handlerError(w, r, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	h.log.Event(mlog.Verbose, func(e mlog.Event) {
		e.String("msg", "Health check request")
		e.String("remote", r.RemoteAddr)
	})

	data := struct {
		Status string `json:"status"`
		Metric string `json:"metric,omitempty"`
		Msg    string `json:"message,omitempty"`
	}{
		Status: "ok",
	}

	h.mu.RLock()
	healthchecks := make([]Health, len(h.healthchecks))
	copy(healthchecks, h.healthchecks)
	h.mu.RUnlock()

iterations:
	for _, v := range healthchecks {
		if err := v.Check(); err != nil {
			if err == errHealthCheckerPanic {
				h.log.Event(mlog.Error, func(e mlog.Event) {
					e.String("msg", "Panic in healthchecker callback")
					e.String("metric", v.Name())
				})
			}
			// update internal metrics
			h.failedhealthchecks.Inc(1)

			h.log.Event(mlog.Warning, func(e mlog.Event) {
				e.String("msg", "Health check failed")
				e.String("metric", v.Name())
				e.String("reason", err.Error())
			})

			// build message
			data.Status = "failed"
			data.Metric = v.Name()
			data.Msg = err.Error()
			w.WriteHeader(http.StatusServiceUnavailable)
			break iterations
		}
	}

	w.Header().Set("Content-Type", contenttypeJSON)
	json.NewEncoder(w).Encode(data)
}

func (h *host) metricsInJSON(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	metrics := make([]interface{}, len(h.metrics))
	copy(metrics, h.metrics)
	h.mu.RUnlock()

	// Build response
	data := make(map[string]interface{})
	m := make(map[string]interface{})
	for _, v := range metrics {
		switch v := v.(type) {
		case Counter:
			m[v.Name()] = v.Count()
		case Gauge:
			m[v.Name()] = v.Value()
		case Label:
			data[v.Name()] = v.Value()
		}
	}

	data["uptime"] = time.Since(h.started)
	data["metrics"] = m

	json.NewEncoder(w).Encode(data)
}

func (h *host) metricsInOpenMetrics(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	metrics := make([]interface{}, len(h.metrics))
	copy(metrics, h.metrics)
	h.mu.RUnlock()

	// Extract labels first
	labels := make([]string, 0)
	for _, v := range metrics {
		switch v := v.(type) {
		case Label:
			var b strings.Builder
			fmt.Fprintf(&b, `%s="%s"`, v.Name(), v.Value())
			labels = append(labels, b.String())
		}
	}
	label := strings.Join(labels, ",")

	for _, v := range metrics {
		switch v := v.(type) {
		case Counter:
			fmt.Fprintf(w, "# HELP %s %s\n", v.Name(), v.Help())
			fmt.Fprintf(w, "# TYPE %s counter\n", v.Name())
			fmt.Fprintf(w, "%s{%s} %d\n", v.Name(), label, v.Count())
		case Gauge:
			fmt.Fprintf(w, "# HELP %s %s\n", v.Name(), v.Help())
			fmt.Fprintf(w, "# TYPE %s gauge\n", v.Name())
			fmt.Fprintf(w, "%s{%s} %d\n", v.Name(), label, v.Value())
		}
	}

	fmt.Fprintf(w, "# HELP %s %s\n", rtuptime, rtuptimehelp)
	fmt.Fprintf(w, "# TYPE %s gauge\n", rtuptime)
	fmt.Fprintf(w, "%s{%s} %d\n", rtuptime, label, time.Since(h.started))
}

func (h *host) metricsValues(w http.ResponseWriter, r *http.Request) {
	h.wg.Add(1)
	defer h.wg.Done()

	// only GET method is allowed
	if r.Method != http.MethodGet {
		h.handlerError(w, r, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	accept := ""
	if accepts, ok := r.Header["Accept"]; ok {
		accept = strings.Join(accepts, ";")
	}

	h.log.Event(mlog.Verbose, func(e mlog.Event) {
		e.String("msg", "Get metrics values")
		e.String("remote", r.RemoteAddr)
		e.String("accept", accept)
	})

	// update runtime metrics
	var memstats runtime.MemStats
	runtime.ReadMemStats(&memstats)
	h.rtmemalloc.Update(int64(memstats.Alloc))
	h.rtgoroutines.Update(int64(runtime.NumGoroutine()))

	switch {
	case strings.Contains(accept, acceptJSON):
		w.Header().Set("Content-Type", contenttypeJSON)
		h.metricsInJSON(w, r)
	default:
		w.Header().Set("Content-Type", contenttypeText)
		h.metricsInOpenMetrics(w, r)
	}
}

func (h *host) loggerGetLevels(w http.ResponseWriter, r *http.Request) {
	h.log.Event(mlog.Verbose, func(e mlog.Event) {
		e.String("msg", "Get logger levels")
		e.String("remote", r.RemoteAddr)
	})

	// get levels from mlog
	levels := h.logbook.Levels()
	tostring := make(map[string]string)
	for k, v := range levels {
		tostring[k] = strings.ToLower(v.String())
	}

	w.Header().Set("Content-Type", contenttypeJSON)
	json.NewEncoder(w).Encode(tostring)
}

func (h *host) loggerPatchLevels(w http.ResponseWriter, r *http.Request) {
	h.log.Event(mlog.Verbose, func(e mlog.Event) {
		e.String("msg", "Patch logger levels")
		e.String("remote", r.RemoteAddr)
	})

	newlevels := make(map[string]string)

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&newlevels); err != nil {
		h.log.Event(mlog.Warning, func(e mlog.Event) {
			e.String("msg", "Parsing request body error")
			e.String("err", err.Error())
		})
		h.handlerError(w, r, http.StatusBadRequest, "Incorrect request format")
		return
	}

	currentlevels := h.logbook.Levels()

	lvm := map[string]mlog.Level{
		"verbose": mlog.Verbose,
		"info":    mlog.Info,
		"warning": mlog.Warning,
		"error":   mlog.Error,
		"fatal":   mlog.Fatal,
	}

	// First pass - check request structure
	for logger, level := range newlevels {
		if _, ok := currentlevels[logger]; !ok {
			h.log.Event(mlog.Error, func(e mlog.Event) {
				e.String("msg", "Requested logger does not exist")
				e.String("logger", logger)
			})
			h.handlerError(w, r, http.StatusBadRequest, "Incorrect request format")
			return
		}
		if _, ok := lvm[level]; !ok {
			h.log.Event(mlog.Error, func(e mlog.Event) {
				e.String("msg", "Requested incorect logger level")
				e.String("logger", logger)
				e.String("level", level)
			})
			h.handlerError(w, r, http.StatusBadRequest, "Incorrect request format")
			return
		}
	}

	// 2nd pass - apply changes
	for logger, level := range newlevels {
		if err := h.logbook.SetLevel(logger, lvm[level]); err != nil {
			h.log.Event(mlog.Warning, func(e mlog.Event) {
				e.String("msg", "Failed call to logbook's SetLevel")
				e.String("err", err.Error())
			})
		}
	}

	w.Header().Set("Content-Type", contenttypeJSON)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}")) // Empty message
}

func (h *host) loggerLevels(w http.ResponseWriter, r *http.Request) {
	h.wg.Add(1)
	defer h.wg.Done()

	switch r.Method {
	case http.MethodGet:
		h.loggerGetLevels(w, r)
	case http.MethodPatch:
		h.loggerPatchLevels(w, r)
	default:
		h.handlerError(w, r, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *host) buildMuxer() *http.ServeMux {
	mux := http.NewServeMux()

	// main handlers
	mux.HandleFunc(pathHealthCheck, h.healthCheck)
	mux.HandleFunc(pathMetricsValues, h.metricsValues)

	// enable debug interface
	if h.config.EnableDebug {
		mux.HandleFunc(pathDebugPprof, pprof.Index)
		mux.HandleFunc(pathDebugPprofCmdline, pprof.Cmdline)
		mux.HandleFunc(pathDebugPprofProfile, pprof.Profile)
		mux.HandleFunc(pathDebugPprofSymbol, pprof.Symbol)
		mux.HandleFunc(pathDebugPprofTrace, pprof.Trace)

		mux.HandleFunc(pathDebugLoggerLevels, h.loggerLevels)

		// enable collecting data for block and mutex
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)
	}

	// catch all other requests
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		h.handlerError(w, r, http.StatusNotFound, "Not Found")
	})

	return mux
}
