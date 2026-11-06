package upcloud

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpCloud_PricingAdapter(t *testing.T) {
	adapter := NewAdapter(nil, "mock_user", "mock_pass")
	assert.Equal(t, "upcloud", adapter.ProviderSlug())

	regions, err := adapter.FetchRegions(nil)
	require.NoError(t, err)
	assert.NotEmpty(t, regions)
	assert.Equal(t, "de-fra1", regions[0].RegionCode)

	sampleJSON := `{
		"plans": {
			"plan": [
				{
					"name": "1xCPU-1GB",
					"core_number": 1,
					"memory_amount": 1024,
					"storage_size": 25,
					"storage_tier": "maxiops",
					"public_traffic_out": 1000,
					"price_hourly": 0.0075,
					"price_monthly": 5.00
				},
				{
					"name": "2xCPU-4GB",
					"core_number": 2,
					"memory_amount": 4096,
					"storage_size": 80,
					"storage_tier": "maxiops",
					"public_traffic_out": 4000,
					"price_hourly": 0.0300,
					"price_monthly": 20.00
				}
			]
		}
	}`

	payload, err := adapter.ParsePlansPayload(strings.NewReader(sampleJSON))
	require.NoError(t, err)
	assert.Equal(t, "upcloud", payload.ProviderSlug)
	assert.Len(t, payload.RawItems, 2)

	p1 := payload.RawItems[0]
	assert.Equal(t, "1xCPU-1GB", p1.PlanCode)
	assert.Equal(t, 1, p1.VCPU)
	assert.Equal(t, int64(1024), p1.RAMMB)
	assert.Equal(t, int64(25), p1.DiskGB)
	assert.Equal(t, "NVME", p1.DiskType)
	assert.Equal(t, "1", p1.BandwidthTB.String())
	assert.Equal(t, "5", p1.MonthlyPriceUSD.String())
}

