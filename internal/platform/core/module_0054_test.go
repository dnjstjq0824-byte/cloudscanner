package core_test

import (
	"context"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/platform/core"
	"github.com/stretchr/testify/require"
)

func TestModule0054_VerifyAdapters(t *testing.T) {
	m := core.NewModule0054()
	providers := []string{"digitalocean", "upcloud", "vultr", "hetzner"}

	for _, p := range providers {
		t.Run("Adapter separation for "+p, func(t *testing.T) {
			err := m.VerifyAdapters(context.Background(), p)
			require.NoError(t, err)
		})
	}
}
