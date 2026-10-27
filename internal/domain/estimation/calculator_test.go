package estimation_test

import (
	"testing"

	"github.com/cloudscanr/cloudscanr/internal/domain/estimation"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEstimation_CalculateTCO(t *testing.T) {
	// Scenario 1: DigitalOcean Droplet with 672h monthly cap
	// Hourly: $0.0357 * 730h = $26.061. But capped at $24.00!
	doSKU := &estimation.NormalizedSKU{
		ID:                    "sku-do-s2vcpu-4gb",
		ProviderSlug:          "digitalocean",
		PlanCode:              "s-2vcpu-4gb",
		HasMonthlyCap:         true,
		CapHours:              672,
		RawHourlyPriceUSD:     decimal.RequireFromString("0.0357"),
		RawMonthlyPriceUSD:    decimal.RequireFromString("24.00"),
		IncludedBandwidthGB:   4000,
		IPv4IncludedCount:     1,
		IPv4MonthlyAddonPrice: decimal.RequireFromString("4.00"),
		BackupIsPercentage:    true,
		BackupPriceOrPct:      decimal.RequireFromString("20.0"), // 20% of base
		SnapshotPerGBMonth:    decimal.RequireFromString("0.06"),
	}

	t.Run("DO standard 730h capped at $24.00", func(t *testing.T) {
		in := estimation.DefaultScenarioInput()
		in.UsageHours = 730
		in.RequiredIPv4Count = 1

		res, err := estimation.CalculateTCO(doSKU, in)
		require.NoError(t, err)
		assert.True(t, res.BaseComputeUSD.Equal(decimal.RequireFromString("24.00")), "Expected $24.00 capped compute, got %s", res.BaseComputeUSD)
		assert.True(t, res.IPv4CostUSD.IsZero())
		assert.True(t, res.TotalTarget.Equal(decimal.RequireFromString("24.00")))
	})

	t.Run("DO with 20% backup and 1 extra IPv4", func(t *testing.T) {
		in := estimation.DefaultScenarioInput()
		in.UsageHours = 730
		in.RequiredIPv4Count = 2 // 1 extra IPv4 = +$4.00
		in.EnableBackups = true  // 20% of $24 = +$4.80

		res, err := estimation.CalculateTCO(doSKU, in)
		require.NoError(t, err)
		assert.True(t, res.BaseComputeUSD.Equal(decimal.RequireFromString("24.00")))
		assert.True(t, res.IPv4CostUSD.Equal(decimal.RequireFromString("4.00")))
		assert.True(t, res.BackupCostUSD.Equal(decimal.RequireFromString("4.80")))
		// Total: 24 + 4 + 4.8 = 32.80
		assert.True(t, res.SubtotalUSD.Equal(decimal.RequireFromString("32.80")))
	})

	t.Run("AWS EC2 uncapped hourly with 10% VAT in KRW", func(t *testing.T) {
		// AWS t3.medium: $0.0416/hr uncapped. 730h = $30.368
		awsSKU := &estimation.NormalizedSKU{
			ID:                    "sku-aws-t3-medium",
			ProviderSlug:          "aws",
			PlanCode:              "t3.medium",
			HasMonthlyCap:         false,
			RawHourlyPriceUSD:     decimal.RequireFromString("0.0416"),
			IncludedBandwidthGB:   100,
			EgressOveragePerGB:    decimal.RequireFromString("0.09"),
			IPv4IncludedCount:     0, // AWS charges $0.005/hr for IPv4 ($3.65/mo)
			IPv4MonthlyAddonPrice: decimal.RequireFromString("3.65"),
		}

		in := estimation.ScenarioInput{
			UsageHours:        730,
			RequiredIPv4Count: 1,
			EstimatedEgressGB: decimal.NewFromInt(200), // 100GB overage = $9.00
			TargetCurrency:    "KRW",
			FXRateToTarget:    decimal.NewFromInt(1400),
			TaxPercentage:     decimal.NewFromInt(10), // 10% VAT
		}

		res, err := estimation.CalculateTCO(awsSKU, in)
		require.NoError(t, err)

		// Base Compute: 0.0416 * 730 = 30.368
		assert.True(t, res.BaseComputeUSD.Equal(decimal.RequireFromString("30.3680")))
		// IPv4: 3.65
		assert.True(t, res.IPv4CostUSD.Equal(decimal.RequireFromString("3.6500")))
		// Egress: (200-100)*0.09 = 9.00
		assert.True(t, res.EgressCostUSD.Equal(decimal.RequireFromString("9.0000")))

		// Subtotal USD: 30.368 + 3.65 + 9.0 = 43.018 USD
		// Subtotal KRW: 43.018 * 1400 = 60,225.20 KRW
		// Tax (10%): 6,022.52 KRW
		// Total KRW: 66,247.72 KRW
		assert.True(t, res.Currency == "KRW")
		assert.True(t, res.SubtotalTarget.Equal(decimal.RequireFromString("60225.20")))
		assert.True(t, res.TotalTarget.Equal(decimal.RequireFromString("66247.72")))
	})
}
