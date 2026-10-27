package core_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModule0060_RateLimiterAndBackoff(t *testing.T) {
	m := core.NewModule0060()
	ctx := context.Background()

	// 1. Token bucket exhaustion test
	err := m.TestRateLimit(ctx, "hetzner", 3.0)
	require.NoError(t, err)

	// 2. Retry-After header parsing test
	dSecs := core.ParseRetryAfter("120")
	assert.Equal(t, 120*time.Second, dSecs)

	futureTime := time.Now().Add(45 * time.Second).UTC().Format(http.TimeFormat)
	dDate := core.ParseRetryAfter(futureTime)
	assert.True(t, dDate > 40*time.Second && dDate <= 46*time.Second)

	// 3. Jittered exponential backoff testing
	base := 100 * time.Millisecond
	max := 5 * time.Second
	for attempt := 1; attempt <= 5; attempt++ {
		backoff := core.CalculateBackoff(attempt, base, max)
		assert.True(t, backoff >= 0)
		assert.True(t, backoff <= max)
	}
}
