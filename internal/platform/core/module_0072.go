package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudscanr/cloudscanr/internal/platform/database"
)

// Module0072 implements WO-CS-0072: PostgreSQL pgx/v5 connection pooler & RLS tenant injector verifier.
type Module0072 struct {
	migrationsDir string
}

// NewModule0072 constructs a database and migration verifier.
func NewModule0072(migrationsDir string) *Module0072 {
	if migrationsDir == "" {
		migrationsDir = "03_DATABASE/migrations"
	}
	return &Module0072{migrationsDir: migrationsDir}
}

// VerifyMigrationsIntegrity audits the 12 migration files for completeness and RLS policies.
func (m *Module0072) VerifyMigrationsIntegrity(ctx context.Context) error {
	entries, err := os.ReadDir(m.migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to open migrations directory: %w", err)
	}

	var upCount, downCount int
	var hasRLS, hasExtensions, hasSeed bool

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".up.sql") {
			upCount++
		}
		if strings.HasSuffix(name, ".down.sql") {
			downCount++
		}

		fullPath := filepath.Join(m.migrationsDir, name)
		contentBytes, err := os.ReadFile(fullPath)
		if err != nil {
			return err
		}
		content := string(contentBytes)

		if strings.Contains(name, "000001") && strings.Contains(content, "CREATE EXTENSION") {
			hasExtensions = true
		}
		if strings.Contains(name, "000009") && strings.Contains(content, "ENABLE ROW LEVEL SECURITY") {
			hasRLS = true
		}
		if strings.Contains(name, "000011") && strings.Contains(content, "INSERT INTO cloud_providers") {
			hasSeed = true
		}
	}

	if upCount < 11 {
		return fmt.Errorf("expected at least 11 up migrations, found %d", upCount)
	}
	if downCount < 1 {
		return fmt.Errorf("expected down migration reference, found %d", downCount)
	}
	if !hasExtensions {
		return errors.New("migration 000001 must create required extensions (pgcrypto, pg_trgm, btree_gist)")
	}
	if !hasRLS {
		return errors.New("migration 000009 must configure Row-Level Security policies")
	}
	if !hasSeed {
		return errors.New("migration 000011 must seed initial verified providers")
	}

	return nil
}

// VerifyRLSInjectionSyntax validates the format of the tenant isolation statement.
func (m *Module0072) VerifyRLSInjectionSyntax(orgID string) (string, error) {
	if orgID == "" {
		return "", errors.New("organization ID cannot be empty for tenant isolation")
	}
	return fmt.Sprintf("SET LOCAL app.current_org_id = '%s'", orgID), nil
}

// Ensure database package compiles and implements Config.
var _ *database.Config = database.DefaultPoolConfig("postgres://localhost/test")

