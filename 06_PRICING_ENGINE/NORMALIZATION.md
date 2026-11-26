# NORMALIZATION.md — SaaS Subscription Plan-Price & Deal Normalization Specification
**Document ID:** `CS-PE-001`

> **Pivot addendum (2026-11-26):** Reframed for TirpScanner. The cloud hardware/billing-cap model (RAM/CPU classes, 672h caps) is superseded by subscription normalization below; the retired model is preserved in git history.

## 1. Canonical Plan & Billing Normalization
- **Billing Cycle**: Vendor cycles (`monthly`, `annual`, `2-year`, `lifetime`, `usage-based`) are normalized to a canonical `billing_cycle` plus a comparable `monthly_equivalent` price (annual total / 12; lifetime deals flagged non-comparable for monthly math and shown as one-time cost).
- **Seat Model Classification**:
  - `FLAT`: One price covers the whole team (e.g., flat team plans).
  - `PER_SEAT`: Price multiplies by seat count (e.g., per-user SaaS).
  - `TIERED_SEATS`: Unit price changes at seat thresholds; TirpScanner stores the tier table.
  - `USAGE`: Metered billing (AI credits, API calls); estimates use declared usage inputs.
- **Discount vs Deal Separation**:
  - Structural discounts (annual billing, multi-year commitment) are normalized into `annual_discount_pct` and always comparable.
  - Time-boxed deals (coupons, seasonal sales, launch pricing) are stored as `promotions` with `expires_at` and never blended into the normalized baseline — deal savings are shown as a separate line item.
- **Trial Normalization**: Trial length (`trial_days`) and card requirement (`trial_requires_card`) are captured uniformly so trials compare fairly across vendors.
