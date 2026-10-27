package core

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cloudscanr/cloudscanr/internal/platform"
)

// Module0051 implements WO-CS-0051: CloudScanr independent domain & TreeScan complete isolation boundary.
type Module0051 struct {
	config *platform.PlatformConfig
}

// NewModule0051 constructs a new isolation boundary verifier.
func NewModule0051(cfg *platform.PlatformConfig) *Module0051 {
	return &Module0051{config: cfg}
}

// VerifyBoundary checks all domain and database boundaries against TreeScan cross-contamination.
func (m *Module0051) VerifyBoundary(ctx context.Context) error {
	if m.config == nil {
		return errors.New("platform config is nil")
	}

	// 1. Verify Article 1 zero-coupling across all configuration values first
	blockedTokens := []string{"treescan", "tree_scan", "ts_prod", "treescan_db"}
	inspectTargets := []string{
		m.config.Domain,
		m.config.DatabaseURL,
		m.config.RedisURL,
		m.config.OpenSearchURL,
	}

	for _, target := range inspectTargets {
		lower := strings.ToLower(target)
		for _, token := range blockedTokens {
			if strings.Contains(lower, token) {
				return fmt.Errorf("article 1 isolation violation: target '%s' contains forbidden token '%s'", target, token)
			}
		}
	}

	// 2. Verify Domain for production
	if !strings.EqualFold(m.config.Domain, "cloudscanr.com") && m.config.Environment == "production" {
		return fmt.Errorf("invalid production domain: %s (expected cloudscanr.com)", m.config.Domain)
	}

	return nil
}
