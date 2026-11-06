package aws

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAWS_PricingAdapter(t *testing.T) {
	adapter := NewAdapter(nil, "us-east-1")
	assert.Equal(t, "aws", adapter.ProviderSlug())

	regions, err := adapter.FetchRegions(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, regions)

	payload, err := adapter.FetchPlansAndPrices(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "aws", payload.ProviderSlug)
	assert.NotEmpty(t, payload.RawItems)

	// Check t4g.nano
	found := false
	for _, item := range payload.RawItems {
		if item.PlanCode == "t4g.nano" {
			found = true
			assert.Equal(t, 2, item.VCPU)
			assert.Equal(t, int64(512), item.RAMMB)
			assert.False(t, item.IPv4Included, "AWS must not mark IPv4 as included (charges $0.005/hr surcharge)")
			assert.True(t, item.HourlyPriceUSD.GreaterThan(item.MonthlyPriceUSD.Div(item.MonthlyPriceUSD).Mul(item.HourlyPriceUSD).Sub(item.HourlyPriceUSD)))
			break
		}
	}
	assert.True(t, found, "t4g.nano must be present in parsed offers")
}

