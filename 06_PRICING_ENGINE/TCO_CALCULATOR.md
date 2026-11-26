# TCO_CALCULATOR.md — Subscription Scenario Total Mathematical Model
**Document ID:** `CS-PE-003`

> **Pivot addendum (2026-11-26):** Reframed for TirpScanner — from hourly-cap cloud TCO (IPv4, egress, backups) to subscription scenario totals (billing cycle, seats, discounts, deals, trial credit). The retired model is preserved in git history.

## 1. Deterministic Decimal Equations (`shopspring/decimal`)
```go
func CalculateSubscriptionScenarioTotal(p PlanSnapshot, req ScenarioInput, fx Decimal, taxPct Decimal) TCOResult {
    // 1. Base Plan with Billing-Cycle & Seat Awareness
    cycleMonths := DecimalFromInt(req.BillingCycleMonths) // 1 = monthly, 12 = annual
    billableSeats := max(p.IncludedSeats, req.SeatCount)
    basePlan := p.MonthlyEquivalent.Mul(cycleMonths).Mul(DecimalFromInt(billableSeats))
    // 2. Add-Ons & Usage
    addonCost := req.AddonMonthly.Mul(cycleMonths)
    usageCost := req.EstimatedUsageUnits.Mul(p.UsageUnitPrice)
    // 3. Structural Discount (annual/multi-year) — always comparable
    discountAmount := basePlan.Mul(p.StructuralDiscountPct).Div(DecimalFromInt(100))
    // 4. Time-Boxed Deal Savings — separate line item, never blended into baseline
    dealSavings := req.DealAmountCapped(basePlan.Sub(discountAmount))
    // 5. Trial Credit (prorated first term only)
    trialCredit := Zero
    if req.IsFirstTerm && p.TrialDays > 0 {
        trialCredit = basePlan.Mul(DecimalFromInt(p.TrialDays)).Div(DecimalFromInt(30 * int(cycleMonths.IntPart())))
    }
    // 6. Native Subtotal -> Target Currency -> Tax
    subtotalNative := basePlan.Add(addonCost).Add(usageCost).Sub(discountAmount).Sub(dealSavings).Sub(trialCredit)
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
