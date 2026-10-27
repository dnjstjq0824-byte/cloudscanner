package core_test

import (
	"context"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModule0055_MaskingLogger(t *testing.T) {
	m := core.NewModule0055()
	ctx := context.Background()

	tests := []struct {
		name     string
		key      string
		val      string
		expected string
	}{
		{
			name:     "Redact authorization bearer token",
			key:      "authorization",
			val:      "Bearer secret-paseto-token-12345",
			expected: core.RedactedPlaceholder,
		},
		{
			name:     "Redact cloud API token",
			key:      "token",
			val:      "dop_v1_abcdef1234567890",
			expected: core.RedactedPlaceholder,
		},
		{
			name:     "Redact user password",
			key:      "password",
			val:      "SuperSecretPassword!99",
			expected: core.RedactedPlaceholder,
		},
		{
			name:     "Preserve safe metadata",
			key:      "provider",
			val:      "digitalocean",
			expected: "digitalocean",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := m.SanitizeLogEntry(ctx, "test audit event", tc.key, tc.val)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, output)
		})
	}
}
