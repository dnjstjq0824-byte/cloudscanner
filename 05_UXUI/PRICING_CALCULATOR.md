# PRICING_CALCULATOR.md — Interactive Subscription Plan & Deal Comparison UX Specification

**Document ID:** `CS-UX-004`

> **Pivot addendum (2026-11-26):** Reframed for TirpScanner — from 36-dimension cloud TCO to SaaS subscription plan/deal comparison (billing cycle, seats, discounts, trials).

---

## 1. Dual-Column Price Honesty Layout
Every calculation result card displays two side-by-side panels:
- **Panel A: Vendor List Price**
  - Exact native currency (`USD`, `EUR`) and billing cycle (`$12/seat/mo billed annually` or `$15/seat/mo billed monthly`).
  - Billing rule note (e.g., `annual commitment` vs `month-to-month` vs `lifetime deal`).
- **Panel B: TirpScanner Scenario Estimate**
  - Stacked bar and line-item breakdown:
    1. Base Plan (`BillingCycle` + `SeatCount` applied)
    2. Extra Seats Beyond Included Count
    3. Add-Ons & Usage Overage
    4. Annual/Loyalty Discount Applied
    5. Active Deal / Coupon Savings
    6. Free Trial Credit (prorated first term)
    7. Subtotal in Target Currency (with exact FX Rate & UTC timestamp)
    8. Estimated Local VAT/GST (`0%` / `10% KR VAT` / `19% DE MwSt`)
