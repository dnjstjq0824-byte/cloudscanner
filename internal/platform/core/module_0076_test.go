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

func TestMockKMSServiceWithSecret_RestartSafe(t *testing.T) {
	secret := []byte("test-only-kms-secret-that-is-32bytes!")
	boot1, err := core.NewMockKMSServiceWithSecret("key-001", secret)
	require.NoError(t, err)

	payload, err := boot1.SealWithAAD([]byte("dop_v1_token"), "org:1:conn:2:prov:do")
	require.NoError(t, err)

	// Simulate a process restart: a fresh instance from the same secret reopens.
	boot2, err := core.NewMockKMSServiceWithSecret("key-001", secret)
	require.NoError(t, err)
	plain, err := boot2.OpenWithAAD(payload, "org:1:conn:2:prov:do")
	require.NoError(t, err)
	require.Equal(t, "dop_v1_token", string(plain))

	// A different secret must not reopen the payload.
	other, err := core.NewMockKMSServiceWithSecret("key-001", []byte("a-different-32byte-secret-00000000"))
	require.NoError(t, err)
	_, err = other.OpenWithAAD(payload, "org:1:conn:2:prov:do")
	require.Error(t, err)

	_, err = core.NewMockKMSServiceWithSecret("key-001", []byte("short"))
	require.Error(t, err)
}

