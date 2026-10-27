package core_test

import (
	"context"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/platform"
	"github.com/cloudscanr/cloudscanr/internal/platform/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModule0051_VerifyBoundary(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *platform.PlatformConfig
		expectError bool
		errMsg      string
	}{
		{
			name: "Valid standalone CloudScanr configuration",
			cfg: &platform.PlatformConfig{
				Environment:   "production",
				Domain:        "cloudscanr.com",
				DatabaseURL:   "postgres://cloudscanr:pass@db.internal:5432/cloudscanr_prod",
				RedisURL:      "redis://redis.internal:6379/0",
				OpenSearchURL: "http://opensearch.internal:9200",
			},
			expectError: false,
		},
		{
			name: "Violation: Database URL points to TreeScan DB",
			cfg: &platform.PlatformConfig{
				Environment:   "production",
				Domain:        "cloudscanr.com",
				DatabaseURL:   "postgres://user:pass@db.internal:5432/treescan_prod",
				RedisURL:      "redis://redis.internal:6379/0",
				OpenSearchURL: "http://opensearch.internal:9200",
			},
			expectError: true,
			errMsg:      "article 1 isolation violation",
		},
		{
			name: "Violation: Domain contains TreeScan",
			cfg: &platform.PlatformConfig{
				Environment:   "production",
				Domain:        "cloudscanr.treescan.com",
				DatabaseURL:   "postgres://user:pass@db.internal:5432/cloudscanr_prod",
				RedisURL:      "redis://redis.internal:6379/0",
				OpenSearchURL: "http://opensearch.internal:9200",
			},
			expectError: true,
			errMsg:      "article 1 isolation violation",
		},
		{
			name: "Violation: Invalid production domain",
			cfg: &platform.PlatformConfig{
				Environment:   "production",
				Domain:        "otherbrand.com",
				DatabaseURL:   "postgres://user:pass@db.internal:5432/cloudscanr_prod",
				RedisURL:      "redis://redis.internal:6379/0",
				OpenSearchURL: "http://opensearch.internal:9200",
			},
			expectError: true,
			errMsg:      "invalid production domain",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := core.NewModule0051(tc.cfg)
			err := m.VerifyBoundary(context.Background())
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
