package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultPoolConfig(t *testing.T) {
	cfg := DefaultPoolConfig("postgres://localhost:5432/x")
	assert.Equal(t, "postgres://localhost:5432/x", cfg.DatabaseURL)
	assert.Equal(t, int32(25), cfg.MaxConns)
	assert.Equal(t, int32(5), cfg.MinConns)
	assert.Equal(t, time.Hour, cfg.MaxConnLifetime)
	assert.Equal(t, 30*time.Minute, cfg.MaxConnIdleTime)
}

func TestNewPool_Validation(t *testing.T) {
	ctx := context.Background()

	_, err := NewPool(ctx, nil)
	require.Error(t, err)

	_, err = NewPool(ctx, &Config{DatabaseURL: ""})
	require.Error(t, err)

	_, err = NewPool(ctx, &Config{DatabaseURL: "://bad-url"})
	require.Error(t, err)
}
