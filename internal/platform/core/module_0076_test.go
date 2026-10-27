package core_test

import (
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
	"github.com/stretchr/testify/require"
)

func TestModule0076_KMSEnvelopeEncryption(t *testing.T) {
	m := core.NewModule0076(nil)

	err := m.TestKMSTamperResistance(
		"dop_v1_very_secret_cloud_token_1234567890",
		"org-uuid-001",
		"conn-uuid-002",
		"digitalocean",
	)
	require.NoError(t, err, "KMS Envelope encryption with AAD tamper-resistance must pass")
}

