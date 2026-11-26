# PRICING_CONTRACT.md — SaaS Subscription Plan-Price Normalization & Deal API Contract

**Document ID:** `CS-API-003`

> **Pivot addendum (2026-11-26):** Reframed for TirpScanner. The 36-dimension cloud-hardware normalization model is superseded by subscription plan-price normalization for SaaS/AI tools: billing-cycle, seat, discount, trial, FX, and tax normalization plus price-drop provenance. The retired 36-field table is preserved in git history.

---

## 1. Mandatory Plan-Price Normalization Attributes in API Payloads

Every plan returned by `/api/v1/catalog/skus/{sku_id}` and `/api/v1/pricing/estimate` exposes the normalized dimensions below, grouped into logical sub-objects:

| # | Attribute Name | JSON Field Path | Type | Description |
| :--- | :--- | :--- | :--- | :--- |
| 1 | Vendor | `vendor.slug` | `string` | Canonical vendor slug (`notion`, `figma`, etc.) |
| 2 | Tool / Product | `product.slug` | `string` | SaaS/AI tool name and product line |
| 3 | Plan Tier | `plan.tier_code` | `string` | e.g., `FREE`, `PLUS`, `PRO`, `TEAM`, `ENTERPRISE` |
| 4 | Plan Code | `external_plan_code` | `string` | Exact vendor plan identifier |
| 5 | Billing Cycle | `plan.billing_cycle` | `string` | `MONTHLY`, `ANNUAL`, `BIENNIAL`, `LIFETIME`, `USAGE_BASED` |
| 6 | Monthly-Equivalent Price | `pricing.monthly_equivalent` | `decimal_str` | Cycle-normalized price per month in native currency |
| 7 | List Monthly Price | `pricing.raw_monthly_price` | `decimal_str` | Vendor month-to-month list price |
| 8 | List Annual Price | `pricing.raw_annual_price` | `decimal_str` | Vendor annual list price (total, not per-month) |
| 9 | Annual Discount | `pricing.annual_discount_pct` | `decimal_str` | Effective % saved on annual vs monthly |
| 10 | Seat Model | `plan.seat_model` | `string` | `FLAT`, `PER_SEAT`, `TIERED_SEATS`, `USAGE` |
| 11 | Included Seats | `plan.included_seats` | `integer` | Seats included at list price (`1` for individual plans) |
| 12 | Extra Seat Price | `pricing.per_seat_price` | `decimal_str` | Price per additional seat per month |
| 13 | Free Trial | `promotion.trial_days` | `integer` | Trial length in days (`0` if none) |
| 14 | Trial Card Requirement | `promotion.trial_requires_card` | `boolean` | True if a card is required to start the trial |
| 15 | Current Deal | `promotion.deal_label` | `string?` | Active discount label (e.g., `20% off annual`) or `null` |
| 16 | Deal Expiry | `promotion.expires_at` | `datetime?` | UTC expiry of the current deal or `null` |
| 17 | Price Drop | `pricing.drop_pct_30d` | `decimal_str?` | % change vs 30 days ago (negative = drop) or `null` |
| 18 | Availability | `availability_status` | `string` | `AVAILABLE`, `GRANDFATHERED`, `DEPRECATED` |
| 19 | Tax Inclusion | `pricing.tax_included` | `boolean` | Whether raw price includes VAT/GST |
| 20 | Currency | `pricing.native_currency` | `string` | ISO-4217 currency code (`USD`, `EUR`) |
| 21 | FX Rate | `pricing.fx_rate_to_target` | `decimal_str` | Applied exchange rate + `fx_rate_timestamp` |
| 22 | Price Effective Date | `pricing.price_effective_date` | `datetime` | UTC date vendor rate took effect |
| 23 | Collection Timestamp | `pricing.collected_at` | `datetime` | Exact UTC timestamp synced by TirpScanner |
| 24 | Pricing Source | `pricing.official_source_url` | `uri` | Direct link to official pricing page or feed |
| 25 | Confidence / Verification | `pricing.verification_status` | `string` | `VERIFIED`, `FRESH`, `STALE`, `UNVERIFIED` |
