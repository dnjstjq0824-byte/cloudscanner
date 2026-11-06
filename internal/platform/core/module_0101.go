package core

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ctxKeyRequestID is the context key carrying the per-request trace identifier.
type ctxKeyRequestID struct{}

// RequestIDFromContext returns the request ID stored by RequestIDMiddleware.
func RequestIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyRequestID{}).(string); ok {
		return v
	}
	return ""
}

// RequestIDMiddleware assigns a UUIDv4 request ID, propagates it via
// X-Request-Id, and stores it in the request context for log correlation.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-Id")
		if reqID == "" {
			reqID = uuid.NewString()
		}
		w.Header().Set("X-Request-Id", reqID)
		ctx := context.WithValue(r.Context(), ctxKeyRequestID{}, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// statusRecorder captures the response status code for access logging.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// AccessLogMiddleware emits one structured JSON log line per request.
func AccessLogMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r)
			logger.InfoContext(r.Context(), "http request",
				"request_id", RequestIDFromContext(r.Context()),
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration_ms", time.Since(start).Milliseconds(),
				"remote", clientIP(r),
			)
		})
	}
}

// clientIP extracts the originating client IP, preferring X-Forwarded-For.
func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		if host, _, err := net.SplitHostPort(fwd); err == nil {
			return host
		}
		// Single IP without port.
		if ip := net.ParseIP(fwd); ip != nil {
			return fwd
		}
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}

// RateLimitConfig tunes the fixed-window per-IP rate limiter.
type RateLimitConfig struct {
	// RequestsPerMinute is the default budget for general API routes.
	RequestsPerMinute int
	// AuthRequestsPerMinute is the stricter budget for /api/v1/auth/* routes.
	AuthRequestsPerMinute int
	// Window is the fixed window duration.
	Window time.Duration
}

// DefaultRateLimitConfig provides production-safe defaults.
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerMinute:     240,
		AuthRequestsPerMinute: 20,
		Window:                time.Minute,
	}
}

// rateBucket tracks one window of request counts for a single key.
type rateBucket struct {
	count     int
	expiresAt time.Time
}

// RateLimiter is a dependency-free fixed-window per-IP rate limiter.
type RateLimiter struct {
	mu      sync.Mutex
	cfg     RateLimitConfig
	buckets map[string]*rateBucket
	now     func() time.Time
}

// NewRateLimiter constructs a rate limiter with the given config.
func NewRateLimiter(cfg RateLimitConfig) *RateLimiter {
	if cfg.Window <= 0 {
		cfg.Window = time.Minute
	}
	return &RateLimiter{
		cfg:     cfg,
		buckets: make(map[string]*rateBucket),
		now:     time.Now,
	}
}

// allow checks and consumes one request token for key with the given limit.
func (l *RateLimiter) allow(key string, limit int) (allowed bool, remaining int, resetAt time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	b, ok := l.buckets[key]
	if !ok || !now.Before(b.expiresAt) {
		b = &rateBucket{count: 0, expiresAt: now.Add(l.cfg.Window)}
		l.buckets[key] = b
	}
	// Opportunistic cleanup to bound memory.
	if len(l.buckets) > 10000 {
		for k, v := range l.buckets {
			if !now.Before(v.expiresAt) {
				delete(l.buckets, k)
			}
		}
	}

	b.count++
	remaining = limit - b.count
	if remaining < 0 {
		remaining = 0
	}
	return b.count <= limit, remaining, b.expiresAt
}

// Middleware enforces per-IP rate limits and returns RFC 9457 429 on excess.
func (l *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Health probes and metrics scraping are never throttled.
		if r.URL.Path == "/api/v1/health" || r.URL.Path == "/api/v1/readiness" || r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		limit := l.cfg.RequestsPerMinute
		scope := "api"
		if len(r.URL.Path) >= 13 && r.URL.Path[:13] == "/api/v1/auth/" {
			limit = l.cfg.AuthRequestsPerMinute
			scope = "auth"
		}
		if r.URL.Path == "/api/v1/auth/register" || r.URL.Path == "/api/v1/auth/login" {
			limit = l.cfg.AuthRequestsPerMinute
			scope = "auth"
		}

		key := scope + ":" + clientIP(r)
		allowed, remaining, resetAt := l.allow(key, limit)
		w.Header().Set("X-RateLimit-Limit", itoa(limit))
		w.Header().Set("X-RateLimit-Remaining", itoa(remaining))
		w.Header().Set("X-RateLimit-Reset", resetAt.UTC().Format(time.RFC3339))
		if !allowed {
			p := NewProblem(http.StatusTooManyRequests, "Too Many Requests",
				"Rate limit exceeded. Retry after the window resets.", r.URL.Path, RequestIDFromContext(r.Context()))
			WriteProblemDetails(w, p)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// itoa renders small non-negative integers without importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}

// Module0101 wraps production HTTP hardening verification (request ID,
// access logging, and per-IP rate limiting).
type Module0101 struct{}

func NewModule0101() *Module0101 { return &Module0101{} }
