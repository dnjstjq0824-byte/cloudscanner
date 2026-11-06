package vultr_test

import (
	"context"
	"strings"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/adapters/pricing/vultr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVultr_PricingAdapter(t *testing.T) {
	adapter := vultr.NewAdapter(nil)
	ctx := context.Background()

	assert.Equal(t, "vultr", adapter.ProviderSlug())

	regions, err := adapter.FetchRegions(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, regions)
	assert.Equal(t, "icn", regions[3].RegionCode)

	sample := `{
		"plans": [
			{
				"id": "vc2-2c-4gb",
				"vcpu_count": 2,
				"ram": 4096,
				"disk": 80,
				"bandwidth": 3000,
				"monthly_cost": 20.0,
				"type": "vc2",
				"locations": ["ewr", "fra", "icn"]
			}
		]
	}`

	payload, err := adapter.ParsePlansPayload(strings.NewReader(sample))
	require.NoError(t, err)
	require.Len(t, payload.RawItems, 1)

	item := payload.RawItems[0]
	assert.Equal(t, "vc2-2c-4gb", item.PlanCode)
	assert.Equal(t, 2, item.VCPU)
	assert.Equal(t, int64(4096), item.RAMMB)
	assert.Equal(t, "20", item.MonthlyPriceUSD.String())
	assert.Equal(t, "0.0298", item.HourlyPriceUSD.String())

	defaultPayload, err := adapter.FetchPlansAndPrices(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(defaultPayload.RawItems), 2)
}

