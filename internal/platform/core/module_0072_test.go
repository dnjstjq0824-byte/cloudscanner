package core_test

import (
	"context"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModule0072_DatabaseAndRLS(t *testing.T) {
	// From internal/platform/core, root is "../../.."
	m := core.NewModule0072("../../../03_DATABASE/migrations")
	ctx := context.Background()

	// 1. Verify all 12 migration files integrity
	err := m.VerifyMigrationsIntegrity(ctx)
	require.NoError(t, err, "PostgreSQL migrations should pass integrity verification")

	// 2. Verify RLS tenant statement generator
	sqlStmt, errStmt := m.VerifyRLSInjectionSyntax("01HY0TESTORG1234567890123")
	require.NoError(t, errStmt)
	assert.Equal(t, "SET LOCAL app.current_org_id = '01HY0TESTORG1234567890123'", sqlStmt)

	// 3. Reject empty org ID
	_, errEmpty := m.VerifyRLSInjectionSyntax("")
	assert.Error(t, errEmpty)
}

