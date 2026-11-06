package core

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestIDMiddleware_AssignsAndPropagates(t *testing.T) {
	var gotID string
	next := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotID = RequestIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()
	next.ServeHTTP(rec, req)

	require.NotEmpty(t, gotID)
	assert.Equal(t, gotID, rec.Header().Get("X-Request-Id"))

	// Client-supplied IDs are honored for trace continuity.
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	req2.Header.Set("X-Request-Id", "trace-abc-123")
	rec2 := httptest.NewRecorder()
	next.ServeHTTP(rec2, req2)
	assert.Equal(t, "trace-abc-123", rec2.Header().Get("X-Request-Id"))
}

func TestAccessLogMiddleware_EmitsStatus(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	next := AccessLogMiddleware(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/providers", nil)
	rec := httptest.NewRecorder()
	next.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusTeapot, rec.Code)
}

func TestRateLimiter_AllowsUnderLimitThen429(t *testing.T) {
	cfg := RateLimitConfig{RequestsPerMinute: 3, AuthRequestsPerMinute: 1, Window: time.Minute}
	limiter := NewRateLimiter(cfg)
	next := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/providers", nil)
		req.RemoteAddr = "10.0.0.7:1234"
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code, "request %d must pass", i+1)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/providers", nil)
	req.RemoteAddr = "10.0.0.7:1234"
	rec := httptest.NewRecorder()
	next.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
	assert.Equal(t, "application/problem+json", rec.Header().Get("Content-Type"))
	assert.NotEmpty(t, rec.Header().Get("X-RateLimit-Reset"))
}

func TestRateLimiter_StricterBudgetForAuth(t *testing.T) {
	cfg := RateLimitConfig{RequestsPerMinute: 100, AuthRequestsPerMinute: 1, Window: time.Minute}
	limiter := NewRateLimiter(cfg)
	next := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	newLogin := func() (*httptest.ResponseRecorder, *http.Request) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
		req.RemoteAddr = "10.0.0.9:443"
		return httptest.NewRecorder(), req
	}

	rec, req := newLogin()
	next.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	rec, req = newLogin()
	next.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
}

func TestRateLimiter_HealthNeverThrottled(t *testing.T) {
	cfg := RateLimitConfig{RequestsPerMinute: 1, AuthRequestsPerMinute: 1, Window: time.Minute}
	limiter := NewRateLimiter(cfg)
	next := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
		req.RemoteAddr = "10.0.0.11:80"
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestRateLimiter_WindowReset(t *testing.T) {
	cfg := RateLimitConfig{RequestsPerMinute: 1, AuthRequestsPerMinute: 1, Window: 50 * time.Millisecond}
	limiter := NewRateLimiter(cfg)

	allowed, _, _ := limiter.allow("api:1.2.3.4", 1)
	require.True(t, allowed)
	allowed, _, _ = limiter.allow("api:1.2.3.4", 1)
	require.False(t, allowed)

	time.Sleep(60 * time.Millisecond)
	allowed, remaining, _ := limiter.allow("api:1.2.3.4", 1)
	assert.True(t, allowed)
	assert.Equal(t, 0, remaining)
}
