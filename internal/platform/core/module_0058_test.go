package core_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModule0058_PrometheusMetrics(t *testing.T) {
	m := core.NewModule0058()
	ctx := context.Background()

	// 1. Record requests
	count, err := m.RecordAndVerify(ctx, http.MethodGet, "/api/v1/search/skus", http.StatusOK)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), count)

	// 2. Metrics endpoint output verification
	reg := core.NewMetricsRegistry()
	reg.RecordRequest(http.MethodGet, "/api/v1/health", http.StatusOK, 5*time.Millisecond)
	reg.IncActive()

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	handler := reg.MetricsHandler()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, "cloudscanr_http_requests_total")
	assert.Contains(t, body, "cloudscanr_http_active_requests 1")
}

