package estimation

import (
	"errors"

	"github.com/shopspring/decimal"
)

var (
	zero    = decimal.Zero
	hundred = decimal.NewFromInt(100)
)

// CalculateTCO computes accurate monthly TCO according to the Article 8 specifications.
func CalculateTCO(sku *NormalizedSKU, in ScenarioInput) (*TCOResult, error) {
	if sku == nil {
		return nil, errors.New("sku cannot be nil")
	}

	if in.UsageHours <= 0 {
		in.UsageHours = 730
	}
	if in.FXRateToTarget.IsZero() {
		in.FXRateToTarget = decimal.NewFromInt(1)
	}
	if in.TargetCurrency == "" {
		in.TargetCurrency = "USD"
	}

	// 1. Base Compute with Monthly Cap awareness
	hoursDec := decimal.NewFromInt(int64(in.UsageHours))
	uncappedCompute := sku.RawHourlyPriceUSD.Mul(hoursDec)
	baseCompute := uncappedCompute

	if sku.HasMonthlyCap && sku.RawMonthlyPriceUSD.GreaterThan(zero) && uncappedCompute.GreaterThan(sku.RawMonthlyPriceUSD) {
		baseCompute = sku.RawMonthlyPriceUSD
	}

	// 2. Public IPv4 Surcharge
	billableIPv4 := in.RequiredIPv4Count - sku.IPv4IncludedCount
	if billableIPv4 < 0 {
		billableIPv4 = 0
	}
	ipv4Cost := sku.IPv4MonthlyAddonPrice.Mul(decimal.NewFromInt(int64(billableIPv4)))

	// 3. Egress Bandwidth Overage
	includedGB := decimal.NewFromInt(sku.IncludedBandwidthGB)
	overageGB := in.EstimatedEgressGB.Sub(includedGB)
	if overageGB.LessThan(zero) {
		overageGB = zero
	}
	egressCost := overageGB.Mul(sku.EgressOveragePerGB)

	// 4. Backups and Snapshots
	backupCost := zero
	if in.EnableBackups && sku.BackupPriceOrPct.GreaterThan(zero) {
		if sku.BackupIsPercentage {
			backupCost = baseCompute.Mul(sku.BackupPriceOrPct).Div(hundred)
		} else {
			backupCost = sku.BackupPriceOrPct
		}
	}

	snapshotCost := zero
	if in.SnapshotStorageGB.GreaterThan(zero) && sku.SnapshotPerGBMonth.GreaterThan(zero) {
		snapshotCost = in.SnapshotStorageGB.Mul(sku.SnapshotPerGBMonth)
	}

	// 5. Total Native USD
	subtotalUSD := baseCompute.Add(ipv4Cost).Add(egressCost).Add(backupCost).Add(snapshotCost)

	// 6. Currency Conversion and Tax Application
	subtotalTarget := subtotalUSD.Mul(in.FXRateToTarget)
	taxTarget := zero
	if in.TaxPercentage.GreaterThan(zero) {
		taxTarget = subtotalTarget.Mul(in.TaxPercentage).Div(hundred)
	}
	totalTarget := subtotalTarget.Add(taxTarget)

	// Effective hourly rate in target currency
	effectiveHourly := totalTarget.Div(hoursDec)

	return &TCOResult{
		BaseComputeUSD:  baseCompute.RoundBank(4),
		IPv4CostUSD:     ipv4Cost.RoundBank(4),
		EgressCostUSD:   egressCost.RoundBank(4),
		BackupCostUSD:   backupCost.RoundBank(4),
		SnapshotCostUSD: snapshotCost.RoundBank(4),
		SubtotalUSD:     subtotalUSD.RoundBank(4),
		SubtotalTarget:  subtotalTarget.RoundBank(2),
		TaxTarget:       taxTarget.RoundBank(2),
		TotalTarget:     totalTarget.RoundBank(2),
		EffectiveHourly: effectiveHourly.RoundBank(4),
		Currency:        in.TargetCurrency,
	}, nil
}
