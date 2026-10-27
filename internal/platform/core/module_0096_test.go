package core_test

import (
	"context"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
	"github.com/stretchr/testify/require"
)

func TestModule0096_ArchitectureFitness(t *testing.T) {
	// Root of repo is "../../.."
	m := core.NewModule0096("../../..")
	ctx := context.Background()

	// 1. Verify 100% TreeScan independence
	errTS := m.VerifyTreeScanIndependence(ctx)
	require.NoError(t, errTS, "Codebase must maintain 100% independence from TreeScan")

	// 2. Verify Adapter isolation
	errAdapters := m.VerifyAdapterPackageSeparation(ctx)
	require.NoError(t, errAdapters, "Four adapters must maintain strict package separation")
}
