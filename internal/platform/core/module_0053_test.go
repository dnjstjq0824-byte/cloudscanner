package core_test

import (
	"context"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModule0053_VerifyFrontend(t *testing.T) {
	// From internal/platform/core, root is "../../.."
	m := core.NewModule0053("../../..")
	err := m.VerifyFrontend(context.Background())
	require.NoError(t, err, "Nuxt 3 frontend workspace should verify cleanly")

	badVerifier := core.NewModule0053("/invalid/path")
	errBad := badVerifier.VerifyFrontend(context.Background())
	assert.Error(t, errBad)
}

