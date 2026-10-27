package core_test

import (
	"context"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModule0052_VerifyWorkspace(t *testing.T) {
	// Root of the project from internal/platform/core is "../../.."
	m := core.NewModule0052("../../..")
	err := m.VerifyWorkspace(context.Background())
	require.NoError(t, err, "workspace layout should meet monorepo architectural requirements")

	// Testing non-existent directory fails
	badVerifier := core.NewModule0052("/non/existent/path")
	errBad := badVerifier.VerifyWorkspace(context.Background())
	assert.Error(t, errBad)
}
