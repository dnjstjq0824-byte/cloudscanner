package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBPool wraps pgxpool.Pool with tenant-isolated transaction execution.
type DBPool struct {
	pool *pgxpool.Pool
}

// Config specifies pool tuning parameters.
type Config struct {
	DatabaseURL     string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

// DefaultPoolConfig provides production defaults for PostgreSQL 16.
func DefaultPoolConfig(dbURL string) *Config {
	return &Config{
		DatabaseURL:     dbURL,
		MaxConns:        25,
		MinConns:        5,
		MaxConnLifetime: 1 * time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	}
}

// NewPool initializes a new PostgreSQL connection pool.
func NewPool(ctx context.Context, cfg *Config) (*DBPool, error) {
	if cfg == nil || cfg.DatabaseURL == "" {
		return nil, errors.New("database URL is required")
	}

	poolCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = cfg.MinConns
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return &DBPool{pool: pool}, nil
}

// Close closes all pool connections.
func (db *DBPool) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

// Ping verifies active database connectivity.
func (db *DBPool) Ping(ctx context.Context) error {
	if db.pool == nil {
		return errors.New("database pool is nil")
	}
	return db.pool.Ping(ctx)
}

// RawPool exposes underlying pgxpool.Pool.
func (db *DBPool) RawPool() *pgxpool.Pool {
	return db.pool
}

// WithTenantTx executes a database transaction with PostgreSQL Row-Level Security (RLS)
// locked to the specified organization ID via `SET LOCAL app.current_org_id = $1`.
func (db *DBPool) WithTenantTx(ctx context.Context, orgID string, fn func(ctx context.Context, tx pgx.Tx) error) error {
	if db.pool == nil {
		return errors.New("database pool is nil")
	}
	if orgID == "" {
		return errors.New("organization ID is required for tenant-scoped transaction")
	}

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Inject RLS tenant context strictly for this transaction
	// SET LOCAL prevents leakage to connection pool upon transaction end
	_, err = tx.Exec(ctx, "SET LOCAL app.current_org_id = $1", orgID)
	if err != nil {
		return fmt.Errorf("failed to set tenant RLS context: %w", err)
	}

	if err := fn(ctx, tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit tenant transaction: %w", err)
	}
	return nil
}

// Migration represents a single SQL migration file.
type Migration struct {
	ID   string
	Path string
}

// RunMigrations executes pending .up.sql migrations from migrationsDir in
// lexicographical order. Applied files are tracked in schema_migrations so
// every boot is idempotent; each migration runs inside a transaction.
func (db *DBPool) RunMigrations(ctx context.Context, migrationsDir string) ([]string, error) {
	if db.pool == nil {
		return nil, errors.New("database pool is nil")
	}

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var upFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".up.sql") {
			upFiles = append(upFiles, f.Name())
		}
	}
	sort.Strings(upFiles)

	_, err = db.pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		filename TEXT PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`)
	if err != nil {
		return nil, fmt.Errorf("failed to init schema_migrations: %w", err)
	}

	rows, err := db.pool.Query(ctx, `SELECT filename FROM schema_migrations`)
	if err != nil {
		return nil, fmt.Errorf("failed to list applied migrations: %w", err)
	}
	done := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			rows.Close()
			return nil, err
		}
		done[name] = true
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Baseline: docker-compose seeds the schema via /docker-entrypoint-initdb.d
	// on first volume creation. If core tables already exist but nothing is
	// tracked, mark all present files applied instead of re-running DDL.
	if len(done) == 0 && len(upFiles) > 0 {
		var usersRegclass string
		err := db.pool.QueryRow(ctx, `SELECT to_regclass('public.users')::text`).Scan(&usersRegclass)
		if err == nil && usersRegclass != "" {
			for _, fileName := range upFiles {
				if _, err := db.pool.Exec(ctx,
					`INSERT INTO schema_migrations (filename) VALUES ($1) ON CONFLICT DO NOTHING`, fileName); err != nil {
					return nil, fmt.Errorf("failed to baseline migration %s: %w", fileName, err)
				}
				done[fileName] = true
			}
		}
	}

	var applied []string
	for _, fileName := range upFiles {
		if done[fileName] {
			continue
		}
		fullPath := filepath.Join(migrationsDir, fileName)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			return applied, fmt.Errorf("failed to read migration file %s: %w", fileName, err)
		}

		tx, err := db.pool.Begin(ctx)
		if err != nil {
			return applied, fmt.Errorf("failed to begin migration %s: %w", fileName, err)
		}
		if _, err := tx.Exec(ctx, string(content)); err != nil {
			_ = tx.Rollback(ctx)
			return applied, fmt.Errorf("migration failure on %s: %w", fileName, err)
		}
		if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (filename) VALUES ($1)`, fileName); err != nil {
			_ = tx.Rollback(ctx)
			return applied, fmt.Errorf("failed to record migration %s: %w", fileName, err)
		}
		if err := tx.Commit(ctx); err != nil {
			return applied, fmt.Errorf("failed to commit migration %s: %w", fileName, err)
		}
		applied = append(applied, fileName)
	}

	return applied, nil
}

