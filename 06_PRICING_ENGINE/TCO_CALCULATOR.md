# TCO_CALCULATOR.md — Total Cost of Ownership Mathematical Model
**Document ID:** `CS-PE-003`
## 1. Deterministic Decimal Equations (`shopspring/decimal`)
```go
func CalculateMonthlyScenarioTCO(s PriceSnapshot, spec ProductSpec, req ScenarioInput, fx Decimal, taxPct Decimal) TCOResult {
    // 1. Base Compute with Cap Awareness
    uncappedCompute := s.RawHourlyPrice.Mul(DecimalFromInt(req.UsageHours))
    baseCompute := uncappedCompute
    if s.HasMonthlyCap && s.RawMonthlyPrice.GreaterThan(Zero) && uncappedCompute.GreaterThan(s.RawMonthlyPrice) {
        baseCompute = s.RawMonthlyPrice
    }
    // 2. Public IPv4 Surcharge
    billableIPv4 := max(0, req.RequiredIPv4Count - spec.IPv4IncludedCount)
    ipv4Cost := s.IPv4MonthlyAddonPrice.Mul(DecimalFromInt(billableIPv4))
    // 3. Egress Overage
    overageGB := DecimalMax(Zero, req.EstimatedEgressGB.Sub(DecimalFromInt(spec.IncludedBandwidthGB)))
    egressCost := overageGB.Mul(s.EgressOveragePerGBPrice)
    // 4. Backups & Snapshots
    backupCost := Zero
    if req.EnableBackups {
        if s.BackupIsPercentage {
            backupCost = baseCompute.Mul(s.BackupMonthlyPriceOrPct).Div(DecimalFromInt(100))
        } else {
            backupCost = s.BackupMonthlyPriceOrPct
        }
    }
    snapshotCost := req.SnapshotStorageGB.Mul(s.SnapshotPerGBMonthPrice)
    // 5. Native Subtotal -> Target Currency -> Tax
    subtotalNative := baseCompute.Add(ipv4Cost).Add(egressCost).Add(backupCost).Add(snapshotCost)
    subtotalTarget := subtotalNative.Mul(fx)
    taxAmount := subtotalTarget.Mul(taxPct).Div(DecimalFromInt(100))
    return TCOResult{
        SubtotalNative: subtotalNative.RoundBank(4),
        SubtotalTarget: subtotalTarget.RoundBank(2),
        TaxTarget:      taxAmount.RoundBank(2),
        TotalTarget:    subtotalTarget.Add(taxAmount).RoundBank(2),
    }
}
```
