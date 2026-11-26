# B — Business & Affiliate — Monetization Models, Commission Accounting, Individual Eligibility & Deal Boundaries

**Document ID:** `CS-AZ-B` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (TirpScanner pivot reframe, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Pivot Note (2026-11-26):** Live revenue: SaaS/AI-tool affiliate CPA/RevShare + sponsored placement + deal-alert subscriptions. BYOC/cloud resale models retired. See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0001 ~ WO-CS-0050, WO-CS-0351 ~ WO-CS-0450`

---

## 1. 기능 목적 (Functional Purpose)

Define TirpScanner's multi-stream revenue engine (SaaS/AI-tool Affiliate CPA/RevShare, Sponsored Placement with strict editorial separation, Pro deal-alert/subscription features, and content-hub yield), double-entry commission accounting, and practical operational rules for individual Korean developers vs. corporate entities.

---

## 2. User Roles & Actors

Affiliate Business Analyst, Finance & Reconciliation Admin, Partner Manager, Platform Operator.

---

## 3. User Scenarios

1. Operator registers individual publisher accounts on PartnerStack, Impact, CJ Affiliate, Awin, FirstPromoter, Rewardful, and ShareASale using `tirpscanr.com` (registered; wiring pending) and W-8BEN, then applies to SaaS/AI vendor programs (typically 20–30% recurring).
2. A SaaS vendor pays referral rewards as account credit (free months); the system records them under `reward_type = HOSTING_CREDIT` (enum name historical) rather than cash revenue.
3. A user clicks a SaaS affiliate link; the click is recorded, later matched via postback/API reconciliation as `PENDING`, approved after the network lock period as `APPROVED`, and settled in a monthly payout batch as `PAID`.
4. A SaaS vendor requests a sponsored banner on the AI video-editor comparison page; the ad is rendered with an explicit `Sponsored` disclosure badge and excluded from altering the algorithmic deal score.

---

## 4. Technology Stack

Go `shopspring/decimal` for exact financial arithmetic, PostgreSQL double-entry ledger tables (`commissions`, `commission_adjustments`, `affiliate_payouts`), CSV/JSON reconciliation parsers for PartnerStack, Impact, CJ, Awin, and ShareASale.

---

## 5. Database Schema & Tables

**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`affiliate_programs`, `affiliate_accounts`, `affiliate_campaigns`, `affiliate_links`, `affiliate_clicks`, `affiliate_conversions`, `commissions`, `commission_adjustments`, `affiliate_payouts`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)

**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**

- `GET /api/v1/admin/affiliate/programs`
- `POST /api/v1/admin/affiliate/reconcile`
- `GET /api/v1/admin/affiliate/ledger`
- `POST /api/v1/admin/affiliate/payouts`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. User Interface & Experience

Public FTC/KFTC Affiliate Disclosure banner on all comparison/product pages and `/affiliate-disclosure` policy page. Admin Affiliate Revenue Cockpit displaying EPC (Earnings Per 10 Clicks), Conversion Rate, Pending vs. Approved vs. Paid Commission Pipeline, and Cash vs. Account-Credit separation.

---

## 8. Core Processing Logic

1. Strict Model Separation:
   - **Affiliate (LIVE)**: Outbound referral link -> user signs up & pays vendor directly -> TirpScanner earns CPA/RevShare.
   - **Deal Scan (LIVE)**: Scheduled ingestion of vendor pricing pages, network feeds, and coupon sources -> normalized deals, price-drop alerts.
   - **BYOC (RETIRED)**: Former user-account server ops; no new connections. **Reseller (RETIRED)**: Never activated; requires an executed contract if ever revived.
2. Commission State Machine: `CLICKED` -> `CONVERSION_PENDING` -> `APPROVED` (or `REVERSED`/`CLAWBACK`) -> `SETTLED_PAID`.
3. Currency & Tax Handling: Preserves original network payout currency (`USD`, `EUR`, `GBP`) alongside normalized `USD` and `KRW` reporting equivalents at settlement FX rate.

---

## 9. Security & Access Control

Immutable append-only ledger for `commissions` and `commission_adjustments` (PostgreSQL trigger prevents `UPDATE`/`DELETE` on finalized ledger entries; reversals must insert a compensating adjustment row). Role-based access restricted to `PAYOUT_APPROVER` and `SUPER_ADMIN`.

---

## 10. Failure Modes & Resilience

Idempotent conversion ingestion keyed on `(program_id, external_conversion_id)`. Duplicate postbacks or overlapping CSV uploads are safely deduplicated and logged without double-counting revenue.

---

## 11. Inter-Module Dependencies

Depends on Market Research (`00_MARKET_RESEARCH`) and PostgreSQL (`P_POSTGRESQL.md`); feeds into Leads & Attribution (`L_LEADS_AFFILIATE.md`) and Analytics (`Y_ANALYTICS.md`).

---

## 12. Testing & Verification Strategy

Unit tests for recurring-vs-CPA commission calculation, clawback adjustment ledger balancing, cash vs. credit segregation, and idempotent CSV reconciliation imports.

---

## 13. Acceptance Criteria

- [x] Every network/vendor affiliate status (`VERIFIED`, `REQUIRES_APPLICATION`, `CONTRACT_REQUIRED`, `NOT_SUPPORTED`) matches `00_MARKET_RESEARCH` matrices.
- [x] Non-cash (account-credit) referrals are strictly segregated from cash revenue totals.
- [x] Ledger sum of `commissions` + `commission_adjustments` equals exact reconciled payout batch totals with zero floating-point rounding drift.

---

## 14. Required Verification Evidence

PostgreSQL ledger constraint test logs, sample reconciled PartnerStack/Impact/CJ import audit report, and rendered FTC/KFTC disclosure screenshots.

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0001 ~ WO-CS-0050, WO-CS-0351 ~ WO-CS-0450`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
