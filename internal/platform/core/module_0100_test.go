package core_test

import (
	"context"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModule0100_Phase1IntegrationSmoke(t *testing.T) {
	m := core.NewModule0100("../../..")
	ctx := context.Background()

	report, err := m.RunPhase1SmokeTest(ctx)
	require.NoError(t, err, "Phase 1 master smoke test must pass completely")
	assert.True(t, report.AllGatesGreen)
	assert.True(t, report.TreeScanIndependence)
	assert.True(t, report.AdaptersSegregated)
	assert.True(t, report.DatabaseMigrations12)
	assert.True(t, report.RLSTenantIsolation)
	assert.True(t, report.KMSEnvelopeAADPassed)
	assert.True(t, report.SSRFFilterActive)
	assert.True(t, report.IdempotencyEnforced)
}
