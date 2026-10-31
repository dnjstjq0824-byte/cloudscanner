package pricing

import (
	"context"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/adapters/pricing/aws"
	"github.com/cloudscanr/cloudscanr/internal/adapters/pricing/digitalocean"
	"github.com/cloudscanr/cloudscanr/internal/adapters/pricing/hetzner"
	"github.com/cloudscanr/cloudscanr/internal/adapters/pricing/upcloud"
	"github.com/cloudscanr/cloudscanr/internal/adapters/pricing/vultr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPricing_Registry(t *testing.T) {
	reg := NewRegistry()

	require.NoError(t, reg.Register(aws.NewAdapter(nil, "us-east-1")))
	require.NoError(t, reg.Register(digitalocean.NewAdapter(nil, "mock-token")))
	require.NoError(t, reg.Register(hetzner.NewAdapter(nil, "mock-token")))
	require.NoError(t, reg.Register(upcloud.NewAdapter(nil, "mock-user", "mock-pass")))
	require.NoError(t, reg.Register(vultr.NewAdapter(nil)))

	// Retrieve by slug
	doAdapter, err := reg.Get("digitalocean")
	require.NoError(t, err)
	assert.Equal(t, "digitalocean", doAdapter.ProviderSlug())

	// Unknown slug error
	_, err = reg.Get("nonexistent")
	assert.Error(t, err)

	// Harvest AWS (which provides baseline mock catalog without network)
	payloads, err := reg.HarvestAll(context.Background())
	require.NoError(t, err)
	assert.Contains(t, payloads, "aws")
	assert.Equal(t, "aws", payloads["aws"].ProviderSlug)
}
