package digitalocean_test

import (
	"context"
	"strings"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/adapters/pricing/digitalocean"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDigitalOcean_PricingAdapter(t *testing.T) {
	adapter := digitalocean.NewAdapter(nil, "")
	ctx := context.Background()

	// 1. Verify Provider Slug
	assert.Equal(t, "digitalocean", adapter.ProviderSlug())

	// 2. Test Regions
	regions, err := adapter.FetchRegions(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, regions)
	assert.Equal(t, "nyc1", regions[0].RegionCode)

	// 3. Test Parsing JSON fixture
	sampleJSON := `{
		"sizes": [
			{
				"slug": "s-2vcpu-4gb",
				"memory": 4096,
				"vcpus": 2,
				"disk": 80,
				"transfer": 4.0,
				"price_monthly": 24.0,
				"price_hourly": 0.0357,
				"regions": ["nyc1", "fra1"],
				"available": true
			},
			{
				"slug": "old-unavailable-plan",
				"memory": 512,
				"vcpus": 1,
				"disk": 10,
				"transfer": 0.5,
				"price_monthly": 3.0,
				"price_hourly": 0.0045,
				"regions": ["nyc1"],
				"available": false
			}
		]
	}`

	payload, err := adapter.ParseSizesPayload(strings.NewReader(sampleJSON))
	require.NoError(t, err)
	require.Len(t, payload.RawItems, 1, "Unavailable sizes should be filtered out")

	item := payload.RawItems[0]
	assert.Equal(t, "s-2vcpu-4gb", item.PlanCode)
	assert.Equal(t, 2, item.VCPU)
	assert.Equal(t, int64(4096), item.RAMMB)
	assert.Equal(t, int64(80), item.DiskGB)
	assert.True(t, item.MonthlyPriceUSD.Equal(item.MonthlyPriceUSD))
	assert.Equal(t, "24", item.MonthlyPriceUSD.String())
	assert.Equal(t, "0.0357", item.HourlyPriceUSD.String())
	assert.True(t, item.IPv4Included)

	// 4. Test FetchPlansAndPrices fallback
	fallbackPayload, err := adapter.FetchPlansAndPrices(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(fallbackPayload.RawItems), 3)
}
