package core

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// TokenBucket implements a standard token bucket rate limiter.
type TokenBucket struct {
	mu         sync.Mutex
	capacity   float64
	tokens     float64
	refillRate float64 // tokens per second
	lastRefill time.Time
}

// NewTokenBucket creates a bucket with given capacity and refill rate per second.
func NewTokenBucket(capacity float64, refillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow attempts to consume 1 token. Returns true if permitted.
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens = math.Min(tb.capacity, tb.tokens+elapsed*tb.refillRate)
	tb.lastRefill = now

	if tb.tokens >= 1.0 {
		tb.tokens -= 1.0
		return true
	}
	return false
}

// RateLimiterRegistry manages rate limiters for external cloud providers.
type RateLimiterRegistry struct {
	mu       sync.RWMutex
	limiters map[string]*TokenBucket
}

// NewRateLimiterRegistry constructs provider rate limiter store.
func NewRateLimiterRegistry() *RateLimiterRegistry {
	return &RateLimiterRegistry{
		limiters: make(map[string]*TokenBucket),
	}
}

// GetLimiter gets or creates a token bucket for a provider.
func (r *RateLimiterRegistry) GetLimiter(providerSlug string, cap float64, rate float64) *TokenBucket {
	r.mu.Lock()
	defer r.mu.Unlock()

	if tb, exists := r.limiters[providerSlug]; exists {
		return tb
	}

	tb := NewTokenBucket(cap, rate)
	r.limiters[providerSlug] = tb
	return tb
}

// ParseRetryAfter parses standard HTTP Retry-After headers (seconds or date).
func ParseRetryAfter(headerVal string) time.Duration {
	if headerVal == "" {
		return 0
	}

	// Try seconds integer first
	if secs, err := strconv.Atoi(headerVal); err == nil && secs > 0 {
		return time.Duration(secs) * time.Second
	}

	// Try HTTP-date format
	if t, err := http.ParseTime(headerVal); err == nil {
		d := time.Until(t)
		if d > 0 {
			return d
		}
	}

	return 0
}

// CalculateBackoff computes exponential backoff with full jitter.
func CalculateBackoff(attempt int, baseBackoff time.Duration, maxBackoff time.Duration) time.Duration {
	if attempt <= 0 {
		attempt = 1
	}
	// Exp = base * 2^(attempt-1)
	exp := float64(baseBackoff) * math.Pow(2, float64(attempt-1))
	if exp > float64(maxBackoff) {
		exp = float64(maxBackoff)
	}

	// Full jitter: random between [0, exp)
	jittered := rand.Float64() * exp
	return time.Duration(jittered)
}

// Module0060 wraps rate limiter and retry tests.
type Module0060 struct {
	registry *RateLimiterRegistry
}

func NewModule0060() *Module0060 {
	return &Module0060{registry: NewRateLimiterRegistry()}
}

func (m *Module0060) TestRateLimit(ctx context.Context, slug string, capacity float64) error {
	limiter := m.registry.GetLimiter(slug, capacity, 1.0)
	// Consume all capacity
	for i := 0; i < int(capacity); i++ {
		if !limiter.Allow() {
			return errors.New("expected token consumption to succeed within capacity")
		}
	}
	// Next call should be rate-limited
	if limiter.Allow() {
		return errors.New("expected token consumption to be rate-limited when exhausted")
	}
	return nil
}
