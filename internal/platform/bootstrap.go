package platform

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
)

// PlatformConfig holds environment configuration for CloudScanr.
type PlatformConfig struct {
	Environment      string        `json:"environment"`
	Domain           string        `json:"domain"`
	HTTPPort         int           `json:"http_port"`
	DatabaseURL      string        `json:"database_url"`
	RedisURL         string        `json:"redis_url"`
	OpenSearchURL    string        `json:"opensearch_url"`
	KMSKeyID         string        `json:"kms_key_id"`
	ShutdownTimeout  time.Duration `json:"shutdown_timeout"`
	TreeScanEnforced bool          `json:"treescan_enforced"` // Must always be true for Article 1 compliance
}

// DefaultConfig provides secure defaults for CloudScanr.
func DefaultConfig() *PlatformConfig {
	return &PlatformConfig{
		Environment:      getEnv("APP_ENV", "development"),
		Domain:           getEnv("APP_DOMAIN", "cloudscanr.com"),
		HTTPPort:         8080,
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://cloudscanr:cloudscanr@localhost:5432/cloudscanr_prod?sslmode=disable"),
		RedisURL:         getEnv("REDIS_URL", "redis://localhost:6379/0"),
		OpenSearchURL:    getEnv("OPENSEARCH_URL", "http://localhost:9200"),
		KMSKeyID:         getEnv("KMS_KEY_ID", "local-dev-mock-key-001"),
		ShutdownTimeout:  15 * time.Second,
		TreeScanEnforced: true,
	}
}

// Validate ensures zero TreeScan coupling according to Article 1.
func (c *PlatformConfig) Validate() error {
	if c.Domain == "" {
		return errors.New("platform domain must not be empty")
	}

	// Article 1 strict isolation: Never allow TreeScan domain, database, or connection strings
	forbiddenKeywords := []string{"treescan", "tree-scan", "treescan_prod", "treescan.com"}
	combined := strings.ToLower(fmt.Sprintf("%s %s %s", c.Domain, c.DatabaseURL, c.RedisURL))

	for _, kw := range forbiddenKeywords {
		if strings.Contains(combined, kw) {
			return fmt.Errorf("article 1 violation: coupling detected with forbidden keyword '%s'", kw)
		}
	}

	if !strings.Contains(c.Domain, "cloudscanr") && c.Environment == "production" {
		return fmt.Errorf("production domain must be cloudscanr.com, got: %s", c.Domain)
	}

	return nil
}

// Platform manages system state, service lifecycles, and health signals.
type Platform struct {
	Config *PlatformConfig
	Logger *slog.Logger
	mu     sync.RWMutex
	ready  bool
}

// Bootstrap initializes the CloudScanr core platform.
func Bootstrap(ctx context.Context, cfg *PlatformConfig) (*Platform, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("bootstrap validation failure: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})).With("service", "cloudscanr-api", "domain", cfg.Domain)

	p := &Platform{
		Config: cfg,
		Logger: logger,
		ready:  true,
	}

	logger.InfoContext(ctx, "cloudscanr platform bootstrapped successfully",
		"domain", cfg.Domain,
		"env", cfg.Environment,
		"treescan_isolation", "enforced",
	)

	return p, nil
}

// IsReady returns whether the platform is ready to accept traffic.
func (p *Platform) IsReady() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.ready
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
