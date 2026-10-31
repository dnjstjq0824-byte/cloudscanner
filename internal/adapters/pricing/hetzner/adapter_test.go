package hetzner_test

import (
	"context"
	"strings"
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/adapters/pricing/hetzner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHetzner_PricingAdapter(t *testing.T) {
	adapter := hetzner.NewAdapter(nil, "")
	ctx := context.Background()

	assert.Equal(t, "hetzner", adapter.ProviderSlug())

	regions, err := adapter.FetchRegions(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, regions)
	assert.Equal(t, "fsn1", regions[0].RegionCode)

	sample := `{
		"server_types": [
			{
				"id": 1,
				"name": "cx22",
				"cores": 2,
				"memory": 4.0,
				"disk": 40,
				"storage_type": "local",
				"cpu_type": "shared",
				"prices": [
					{"location": "fsn1", "price_hourly": {"net": "0.0055"}, "price_monthly": {"net": "3.79"}}
				]
			}
		]
	}`

	payload, err := adapter.ParseServerTypesPayload(strings.NewReader(sample))
	require.NoError(t, err)
	require.Len(t, payload.RawItems, 1)

	item := payload.RawItems[0]
	assert.Equal(t, "cx22", item.PlanCode)
	assert.Equal(t, 2, item.VCPU)
	assert.Equal(t, int64(4096), item.RAMMB)
	assert.Equal(t, int64(40), item.DiskGB)
	assert.True(t, item.MonthlyPriceUSD.GreaterThan(item.HourlyPriceUSD))

	defaultPayload, err := adapter.FetchPlansAndPrices(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(defaultPayload.RawItems), 2)
}
