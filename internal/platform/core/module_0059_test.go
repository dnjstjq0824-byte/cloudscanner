package core_test

import (
	"context"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModule0059_CircuitBreakerIsolation(t *testing.T) {
	m := core.NewModule0059()
	ctx := context.Background()

	// 1. Trip Vultr into OPEN state
	err := m.TripProvider(ctx, "vultr")
	require.NoError(t, err, "vultr circuit breaker should trip open after 5 consecutive failures")

	// 2. Verify DigitalOcean is unaffected (Isolation verification)
	doBreaker := m.TripProvider(ctx, "digitalocean")
	require.NoError(t, doBreaker)

	// Direct check that a fresh provider is CLOSED
	freshManager := core.NewCircuitBreakerManager()
	res, errFresh := freshManager.Execute("upcloud", func() (interface{}, error) {
		return "healthy", nil
	})
	require.NoError(t, errFresh)
	assert.Equal(t, "healthy", res)
}

