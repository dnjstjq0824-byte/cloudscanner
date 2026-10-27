package core_test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModule0056_RFC9457AndPanicRecovery(t *testing.T) {
	m := core.NewModule0056()
	ctx := context.Background()

	// 1. Validate RFC 9457 model
	prob, err := m.ValidateRFC9457(ctx, http.StatusNotFound, "Resource Not Found", "The requested SKU was not found")
	require.NoError(t, err)
	assert.Equal(t, 404, prob.Status)
	assert.Equal(t, "https://cloudscanr.com/errors/http-404", prob.Type)

	// 2. Validate Panic Recovery Middleware
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	recoveryMw := core.PanicRecoveryMiddleware(logger)

	panickingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("database connection pool exhausted")
	})

	wrapped := recoveryMw(panickingHandler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cloud/servers", nil)
	rec := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		wrapped.ServeHTTP(rec, req)
	})

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, "application/problem+json", res.Header.Get("Content-Type"))
}

