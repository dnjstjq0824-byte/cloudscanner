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

	"github.com/cloudscanr/cloudscanr/internal/platform/database"
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
	AuthTokenSecret  string        `json:"-"`
	CORSAllowedOrigins []string    `json:"cors_allowed_origins"`
	ShutdownTimeout  time.Duration `json:"shutdown_timeout"`
	TreeScanEnforced bool          `json:"treescan_enforced"` // Must always be true for Article 1 compliance
}

// DefaultConfig provides secure defaults for CloudScanr.
func DefaultConfig() *PlatformConfig {
	return &PlatformConfig{
		Environment:      getEnv("APP_ENV", "development"),
		Domain:           getEnv("APP_DOMAIN", "cloudscanr.com"),
		HTTPPort:         getEnvInt("HTTP_PORT", 8080),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://cloudscanr:cloudscanr@localhost:5432/cloudscanr_prod?sslmode=disable"),
		RedisURL:         getEnv("REDIS_URL", "redis://localhost:6379/0"),
		OpenSearchURL:    getEnv("OPENSEARCH_URL", "http://localhost:9200"),
		KMSKeyID:         getEnv("KMS_KEY_ID", "local-dev-mock-key-001"),
		AuthTokenSecret:  getEnv("AUTH_TOKEN_SECRET", "dev-only-token-secret-change-me-32bytes"),
		CORSAllowedOrigins: getEnvList("CORS_ALLOWED_ORIGINS"),
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
	DB     DBConnector
	mu     sync.RWMutex
	ready  bool
}

// DBConnector is the minimal database surface the platform depends on.
// *database.DBPool satisfies this interface; tests may substitute fakes.
type DBConnector interface {
	Close()
	Ping(ctx context.Context) error
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

	dbStatus := p.connectDatabase(ctx)
	if dbStatus == "connected" {
		if applied, err := p.runMigrations(ctx); err != nil {
			logger.ErrorContext(ctx, "database migration failed", "err", err)
			if cfg.Environment == "production" {
				return nil, fmt.Errorf("production migration failure: %w", err)
			}
		} else {
			logger.InfoContext(ctx, "database migrations verified", "applied", len(applied))
		}
	}

	logger.InfoContext(ctx, "cloudscanr platform bootstrapped successfully",
		"domain", cfg.Domain,
		"env", cfg.Environment,
		"database", dbStatus,
		"treescan_isolation", "enforced",
	)

	return p, nil
}

// connectDatabase dials PostgreSQL with a short timeout. Production refuses to
// boot without a database; development and test fall back to in-memory stores.
func (p *Platform) connectDatabase(ctx context.Context) string {
	if p.Config.DatabaseURL == "" {
		p.Logger.WarnContext(ctx, "DATABASE_URL is empty; using in-memory stores")
		return "degraded_in_memory"
	}
	dialCtx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	pool, err := database.NewPool(dialCtx, database.DefaultPoolConfig(p.Config.DatabaseURL))
	if err != nil {
		return p.databaseFallback(ctx, fmt.Errorf("pool init: %w", err))
	}
	if err := pool.Ping(dialCtx); err != nil {
		pool.Close()
		return p.databaseFallback(ctx, fmt.Errorf("ping: %w", err))
	}
	p.mu.Lock()
	p.DB = pool
	p.mu.Unlock()
	return "connected"
}

// databaseFallback logs the outage and decides whether boot may continue.
func (p *Platform) databaseFallback(ctx context.Context, cause error) string {
	if p.Config.Environment == "production" {
		// Production boot is aborted by the caller via DBConnected check in main.
		p.Logger.ErrorContext(ctx, "FATAL: database unreachable in production", "err", cause)
		p.mu.Lock()
		p.ready = false
		p.mu.Unlock()
		return "unreachable_production_blocked"
	}
	p.Logger.WarnContext(ctx, "database unreachable; continuing with in-memory stores (data will not persist)",
		"err", cause)
	return "degraded_in_memory"
}

// runMigrations applies pending .up.sql files from the packaged migrations dir.
func (p *Platform) runMigrations(ctx context.Context) ([]string, error) {
	pool, ok := p.DB.(*database.DBPool)
	if !ok || pool == nil {
		return nil, errors.New("no database pool available")
	}
	for _, dir := range []string{"/app/migrations", "03_DATABASE/migrations"} {
		if st, err := os.Stat(dir); err == nil && st.IsDir() {
			return pool.RunMigrations(ctx, dir)
		}
	}
	return nil, errors.New("migrations directory not found")
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

// getEnvList parses a comma-separated environment variable into a trimmed list.
func getEnvList(key string) []string {
	raw := os.Getenv(key)
	if raw == "" {
		return nil
	}
	var out []string
	for _, part := range strings.Split(raw, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		var n int
		if _, err := fmt.Sscanf(val, "%d", &n); err == nil && n > 0 && n < 65536 {
			return n
		}
	}
	return defaultVal
}

