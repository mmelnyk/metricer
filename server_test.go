package metricer

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.melnyk.org/mlog/testlog"
)

func TestServerMuxerDefault(t *testing.T) {
	mhost := NewHost(nil, testlog.NewLogbook()).(*host)
	muxer := mhost.buildMuxer()

	tests := []struct {
		method string
		url    string
		code   int
	}{
		{"GET", "http://test/health/check", http.StatusOK},
		{"POST", "http://test/health/check", http.StatusMethodNotAllowed},
		{"GET", "http://test/something/else", http.StatusNotFound},
		{"GET", "http://test/debug/pprof/cmdline", http.StatusNotFound},
	}

	for i, v := range tests {
		req := httptest.NewRequest(v.method, v.url, nil)
		h, _ := muxer.Handler(req)
		if h == nil {
			t.Fatalf("Expected (%d): HealthCheck handler, but got nil", i)
		}

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		resp := w.Result()

		if resp.StatusCode != v.code {
			t.Errorf("Expected (%d): %d, but got %d", i, v.code, resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != contenttypeJSON {
			t.Errorf("Expected (%d): Content-type JSON, but got %s", i, resp.Header.Get("Content-Type"))
		}
	}
}

func TestServerMuxerWithDebug(t *testing.T) {
	cfg := &Config{EnableDebug: true}
	mhost := NewHost(cfg, testlog.NewLogbook()).(*host)
	mhost.NewHealthCheck("health", "health help", func() error {
		return nil
	})
	muxer := mhost.buildMuxer()

	tests := []struct {
		method string
		url    string
		code   int
	}{
		{"GET", "http://test/health/check", http.StatusOK},
		{"POST", "http://test/health/check", http.StatusMethodNotAllowed},
		{"GET", "http://test/something/else", http.StatusNotFound},
		{"GET", "http://test/debug/pprof/cmdline", http.StatusOK},
	}

	for i, v := range tests {
		req := httptest.NewRequest(v.method, v.url, nil)
		h, _ := muxer.Handler(req)
		if h == nil {
			t.Fatalf("Expected (%d): HealthCheck handler, but got nil", i)
		}

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		resp := w.Result()

		if resp.StatusCode != v.code {
			t.Errorf("Expected (%d): %d, but got %d", i, v.code, resp.StatusCode)
		}
	}
}

func TestServerHealthCheckOk(t *testing.T) {
	mhost := NewHost(nil, testlog.NewLogbook()).(*host)
	mhost.NewHealthCheck("health", "health help", func() error {
		return nil
	})

	tests := []struct {
		method string
		url    string
		code   int
	}{
		{"GET", "http://test/health/check", http.StatusOK},
		{"POST", "http://test/health/check", http.StatusMethodNotAllowed},
	}

	for i, v := range tests {
		req := httptest.NewRequest(v.method, v.url, nil)
		w := httptest.NewRecorder()
		mhost.healthCheck(w, req)

		resp := w.Result()

		if resp.StatusCode != v.code {
			t.Errorf("Expected (%d): %d, but got %d", i, v.code, resp.StatusCode)
		}
	}
}

func TestServerHealthCheckFailed(t *testing.T) {
	mhost := NewHost(nil, testlog.NewLogbook()).(*host)
	mhost.NewHealthCheck("health", "health help", func() error {
		return errors.New("failed")
	})

	tests := []struct {
		method string
		url    string
		code   int
	}{
		{"GET", "http://test/health/check", http.StatusServiceUnavailable},
		{"POST", "http://test/health/check", http.StatusMethodNotAllowed},
	}

	for i, v := range tests {
		req := httptest.NewRequest(v.method, v.url, nil)
		w := httptest.NewRecorder()
		mhost.healthCheck(w, req)

		resp := w.Result()

		if resp.StatusCode != v.code {
			t.Errorf("Expected (%d): %d, but got %d", i, v.code, resp.StatusCode)
		}
	}
}

func TestServerHealthCheckFailedWithPanic(t *testing.T) {
	mhost := NewHost(nil, testlog.NewLogbook()).(*host)
	mhost.NewHealthCheck("health", "health help", func() error {
		var counter Counter
		counter.Inc(1)
		return nil
	})

	tests := []struct {
		method string
		url    string
		code   int
	}{
		{"GET", "http://test/health/check", http.StatusServiceUnavailable},
		{"POST", "http://test/health/check", http.StatusMethodNotAllowed},
	}

	for i, v := range tests {
		req := httptest.NewRequest(v.method, v.url, nil)
		w := httptest.NewRecorder()
		mhost.healthCheck(w, req)

		resp := w.Result()

		if resp.StatusCode != v.code {
			t.Errorf("Expected (%d): %d, but got %d", i, v.code, resp.StatusCode)
		}
	}
}

func TestServerMetricsValuesPlain(t *testing.T) {
	mhost := NewHost(nil, testlog.NewLogbook()).(*host)

	tests := []struct {
		method string
		url    string
		code   int
	}{
		{"GET", "http://test/metrics/values", http.StatusOK},
		{"POST", "http://test/metrics/values", http.StatusMethodNotAllowed},
	}

	for i, v := range tests {
		req := httptest.NewRequest(v.method, v.url, nil)
		req.Header.Set("Accept", acceptText)
		w := httptest.NewRecorder()
		mhost.metricsValues(w, req)

		resp := w.Result()

		if resp.StatusCode != v.code {
			t.Errorf("Expected (%d): %d, but got %d", i, v.code, resp.StatusCode)
		}
	}
}

func TestServerMetricsValuesJson(t *testing.T) {
	mhost := NewHost(nil, testlog.NewLogbook()).(*host)

	tests := []struct {
		method string
		url    string
		code   int
	}{
		{"GET", "http://test/metrics/values", http.StatusOK},
		{"POST", "http://test/metrics/values", http.StatusMethodNotAllowed},
	}

	for i, v := range tests {
		req := httptest.NewRequest(v.method, v.url, nil)
		req.Header.Set("Accept", acceptJSON)
		w := httptest.NewRecorder()
		mhost.metricsValues(w, req)

		resp := w.Result()

		if resp.StatusCode != v.code {
			t.Errorf("Expected (%d): %d, but got %d", i, v.code, resp.StatusCode)
		}
	}
}

func TestServerLoggerLevels(t *testing.T) {
	cfg := &Config{EnableDebug: true}
	mhost := NewHost(cfg, testlog.NewLogbook()).(*host)

	tests := []struct {
		method  string
		url     string
		code    int
		payload string
	}{
		{http.MethodGet, "http://test/debug/logger/levels", http.StatusOK, ""},
		{http.MethodPost, "http://test/debug/logger/levels", http.StatusMethodNotAllowed, ""},
		{http.MethodPatch, "http://test/debug/logger/levels", http.StatusBadRequest, ""},
		{http.MethodPatch, "http://test/debug/logger/levels", http.StatusOK, "{}"},
		{http.MethodPatch, "http://test/debug/logger/levels", http.StatusBadRequest, `{"test":"code"}`},
		{http.MethodPatch, "http://test/debug/logger/levels", http.StatusBadRequest, `{"test":"fatal"}`},
		{http.MethodPatch, "http://test/debug/logger/levels", http.StatusOK, `{"DEFAULT":"fatal"}`},
		{http.MethodPatch, "http://test/debug/logger/levels", http.StatusBadRequest, `{"DEFAULT":"foo"}`},
	}

	for i, v := range tests {
		req := httptest.NewRequest(v.method, v.url, bytes.NewReader([]byte(v.payload)))
		w := httptest.NewRecorder()
		mhost.loggerLevels(w, req)

		resp := w.Result()

		if resp.StatusCode != v.code {
			t.Errorf("Expected (%d): %d, but got %d", i, v.code, resp.StatusCode)
		}
	}
}
