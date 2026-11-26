# E — Deal & Estimation Engine — Plan-Price Normalization, Promo Stacking, Renewal Math, Taxes & FX

**Document ID:** `CS-AZ-E` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (TirpScanner pivot reframe, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Pivot Note (2026-11-26):** 36-field cloud TCO engine (730h billing, egress, IPv4) is RETIRED. Live engine normalizes SaaS plan pricing and scores deals. See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0251 ~ WO-CS-0300`

---

## 1. 기능 목적 (Functional Purpose)

Implement the SaaS Plan Normalization & Deal-Scan Engine. Separate raw vendor sticker prices from scenario-based effective monthly/annual cost, accounting for monthly vs annual per-seat billing, intro vs renewal pricing, trial/freemium value, coupon/promo stacking, credit-pack overages, VAT/GST, and real-time FX conversion — and compute a reproducible deal score per plan.

---

## 2. User Roles & Actors

Pricing Engine Agent, Financial Analyst, End User Comparing SaaS Plans.

---

## 3. User Scenarios

1. User compares an AI writer's `Pro` plan ($20/seat/mo monthly vs $12/seat/mo annual, 3-seat minimum, renewing at full rate) against a competitor's flat `$49/mo up to 10 seats` plan with an active 20%-off coupon.
2. User toggles between `Monthly Billing`, `Annual Billing (per-seat)`, `Annual Upfront Total`, and `Year-2 Renewal Cost` to expose intro-discount traps.
3. User in South Korea selects `KRW (₩)` and `Include 10% VAT`; the engine converts USD/EUR raw rates using timestamped FX rates and clearly displays the FX snapshot timestamp and tax breakdown.

---

## 4. Technology Stack

Go `internal/domain/estimation` using `shopspring/decimal` (12 decimal places internal precision, Half-Even rounding for display), PostgreSQL `price_components`, `billing_terms`, `fx_rates`, `tax_policies`.

---

## 5. Database Schema & Tables

**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`price_snapshots`, `price_components`, `billing_terms`, `promotions`, `fx_rates`, `tax_policies`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)

**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**

- `POST /api/v1/pricing/estimate`
- `GET /api/v1/pricing/fx-rates`
- `GET /api/v1/pricing/scenarios`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. User Interface & Experience

Interactive Cost Breakdown Drawer & Plan Calculator (`/calculator`) displaying two distinct columns: **Raw Vendor Price** (exact currency & unit billed by vendor) vs. **TirpScanner Effective Cost** (itemized stacked bar: Base Seats + Add-ons/Credits + Coupon Discount + Renewal Delta + Tax).

---

## 8. Core Processing Logic

Formula for Scenario Effective Cost:

`BaseSeats = SeatCount * PerSeatRate(BillingInterval)` (honoring `min_seats` and flat-team caps).  
`PromoDiscount = StackableCoupons(BaseSeats, promotions)` (percent-then-fixed order, expiry-checked).  
`CreditOverage = MAX(0, EstimatedCredits - IncludedCredits) * CreditOverageRate`.  
`SubtotalNative = BaseSeats - PromoDiscount + CreditOverage + AddonCost`.  
`RenewalCost = SubtotalNative recomputed with renewal (non-intro) rates`.  
`TotalTargetCurrency = SubtotalNative * FXRate(Native -> Target) * (1 + TaxRate)`.  
`DealScore = f(discount depth vs 90-day median, coupon stackability, renewal fairness, trial value)` (inputs defined in `R_RANKING.md`; affiliate commission is never an input).

---

## 9. Security & Access Control

Deterministic calculation inputs validated against injection and overflow (`SeatCount` bounded `1..10`, `EstimatedCredits` bounded `0..10`).

---

## 10. Failure Modes & Resilience

If live FX rate feed is older than 24 hours, display `FX Rate Last Updated: YYYY-MM-DD HH:MM UTC (Stale Warning)` alongside the native currency price.

---

## 11. Inter-Module Dependencies

Depends on Catalog (`C_CATALOG.md`) and Data Acquisition (`D_DATA_ACQUISITION.md`).

---

## 12. Testing & Verification Strategy

Table-driven Go unit tests covering monthly vs annual normalization, seat-minimum and team-cap edge cases, coupon stacking/expiry, intro-vs-renewal deltas, credit-pack overages, and FX/tax composition.

---

## 13. Acceptance Criteria

- [x] All plan-normalization attributes (billing interval, per-seat rate, minimums/caps, trial/freemium, quotas, renewal terms, promos, FX, tax) are modeled in the database and estimation payload.
- [x] Raw vendor billing price and TirpScanner estimated effective cost are never conflated in API responses or UI.
- [x] Intro vs renewal pricing is explicitly labeled and comparable by the user.

---

## 14. Required Verification Evidence

Go decimal unit test results (`go test ./internal/domain/estimation/...`), JSON output of `/api/v1/pricing/estimate`, and calculator UI verification.

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0251 ~ WO-CS-0300`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
