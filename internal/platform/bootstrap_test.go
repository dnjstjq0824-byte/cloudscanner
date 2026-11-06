package platform

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBootstrap_DatabaseFallback verifies that development boot continues with
// in-memory stores when PostgreSQL is unreachable (Docker-free CI/dev).
func TestBootstrap_DatabaseFallback(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Environment = "development"
	cfg.DatabaseURL = "postgres://invalid:invalid@127.0.0.1:1/cloudscanr_prod?sslmode=disable&connect_timeout=1"

	p, err := Bootstrap(context.Background(), cfg)
	require.NoError(t, err)
	require.NotNil(t, p)
	assert.True(t, p.IsReady())
	assert.Nil(t, p.DB)
}

// TestBootstrap_ProductionRequiresDatabase verifies production refuses to
// serve traffic without a reachable database.
func TestBootstrap_ProductionRequiresDatabase(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Environment = "production"
	cfg.Domain = "cloudscanr.com"
	cfg.DatabaseURL = "postgres://invalid:invalid@127.0.0.1:1/cloudscanr_prod?sslmode=disable&connect_timeout=1"

	p, err := Bootstrap(context.Background(), cfg)
	require.NoError(t, err)
	require.NotNil(t, p)
	assert.False(t, p.IsReady())
	assert.Nil(t, p.DB)
}
