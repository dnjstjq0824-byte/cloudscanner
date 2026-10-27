package core

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// MetricsRegistry maintains core Prometheus counters and histograms.
type MetricsRegistry struct {
	mu             sync.RWMutex
	httpRequests   map[string]*uint64
	latencySumMs   map[string]*uint64
	activeRequests int64
}

// NewMetricsRegistry constructs a thread-safe in-process metrics store.
func NewMetricsRegistry() *MetricsRegistry {
	return &MetricsRegistry{
		httpRequests: make(map[string]*uint64),
		latencySumMs: make(map[string]*uint64),
	}
}

// RecordRequest increments request count and adds duration to the latency accumulator.
func (m *MetricsRegistry) RecordRequest(method, path string, status int, duration time.Duration) {
	key := fmt.Sprintf("%s_%s_%d", method, path, status)
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.httpRequests[key]; !exists {
		var c uint64
		var l uint64
		m.httpRequests[key] = &c
		m.latencySumMs[key] = &l
	}

	atomic.AddUint64(m.httpRequests[key], 1)
	atomic.AddUint64(m.latencySumMs[key], uint64(duration.Milliseconds()))
}

// IncActive increments current in-flight requests.
func (m *MetricsRegistry) IncActive() {
	atomic.AddInt64(&m.activeRequests, 1)
}

// DecActive decrements current in-flight requests.
func (m *MetricsRegistry) DecActive() {
	atomic.AddInt64(&m.activeRequests, -1)
}

// MetricsHandler serves Prometheus exposition format.
func (m *MetricsRegistry) MetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		m.mu.RLock()
		defer m.mu.RUnlock()

		fmt.Fprintf(w, "# HELP cloudscanr_http_requests_total Total number of HTTP requests\n")
		fmt.Fprintf(w, "# TYPE cloudscanr_http_requests_total counter\n")
		for key, ptr := range m.httpRequests {
			fmt.Fprintf(w, "cloudscanr_http_requests_total{route=\"%s\"} %d\n", key, atomic.LoadUint64(ptr))
		}

		fmt.Fprintf(w, "# HELP cloudscanr_http_active_requests Current in-flight HTTP requests\n")
		fmt.Fprintf(w, "# TYPE cloudscanr_http_active_requests gauge\n")
		fmt.Fprintf(w, "cloudscanr_http_active_requests %d\n", atomic.LoadInt64(&m.activeRequests))
	}
}

// Module0058 provides metrics testing.
type Module0058 struct {
	registry *MetricsRegistry
}

func NewModule0058() *Module0058 {
	return &Module0058{registry: NewMetricsRegistry()}
}

func (m *Module0058) RecordAndVerify(ctx context.Context, method, path string, status int) (uint64, error) {
	m.registry.RecordRequest(method, path, status, 25*time.Millisecond)
	key := fmt.Sprintf("%s_%s_%d", method, path, status)

	m.registry.mu.RLock()
	defer m.registry.mu.RUnlock()

	ptr, exists := m.registry.httpRequests[key]
	if !exists {
		return 0, fmt.Errorf("metric key %s not recorded", key)
	}
	return atomic.LoadUint64(ptr), nil
}
