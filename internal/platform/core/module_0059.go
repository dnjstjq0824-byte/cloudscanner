package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/sony/gobreaker"
)

// CircuitBreakerManager manages isolated circuit breakers per cloud provider.
type CircuitBreakerManager struct {
	mu       sync.RWMutex
	breakers map[string]*gobreaker.CircuitBreaker
}

// NewCircuitBreakerManager constructs a manager with Article 6/10 failure isolation.
func NewCircuitBreakerManager() *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*gobreaker.CircuitBreaker),
	}
}

// GetBreaker retrieves or initializes an independent circuit breaker for a provider.
func (m *CircuitBreakerManager) GetBreaker(providerSlug string) *gobreaker.CircuitBreaker {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cb, exists := m.breakers[providerSlug]; exists {
		return cb
	}

	st := gobreaker.Settings{
		Name:        fmt.Sprintf("cloudscanr-cb-%s", providerSlug),
		MaxRequests: 3,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Trip after 5 consecutive failures
			return counts.ConsecutiveFailures >= 5
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			// State transition callback
		},
	}

	cb := gobreaker.NewCircuitBreaker(st)
	m.breakers[providerSlug] = cb
	return cb
}

// Execute wraps an external provider call in its dedicated circuit breaker.
func (m *CircuitBreakerManager) Execute(providerSlug string, fn func() (interface{}, error)) (interface{}, error) {
	cb := m.GetBreaker(providerSlug)
	return cb.Execute(fn)
}

// Module0059 provides test harness for circuit breaker verification.
type Module0059 struct {
	manager *CircuitBreakerManager
}

func NewModule0059() *Module0059 {
	return &Module0059{manager: NewCircuitBreakerManager()}
}

// TripProvider forces 5 failures to trip the circuit breaker into OPEN state.
func (m *Module0059) TripProvider(ctx context.Context, providerSlug string) error {
	cb := m.manager.GetBreaker(providerSlug)
	for i := 0; i < 5; i++ {
		_, _ = cb.Execute(func() (interface{}, error) {
			return nil, errors.New("upstream provider 503 gateway unavailable")
		})
	}

	// 6th call should immediately return ErrOpenState
	_, err := cb.Execute(func() (interface{}, error) {
		return "should not run", nil
	})

	if !errors.Is(err, gobreaker.ErrOpenState) {
		return fmt.Errorf("expected ErrOpenState, got: %v", err)
	}
	return nil
}
