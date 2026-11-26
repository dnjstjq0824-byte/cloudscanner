# Y — Yield & Analytics — Affiliate Revenue Attribution, Funnel Conversion, EPC & Deal Performance

**Document ID:** `CS-AZ-Y` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (TirpScanner pivot reframe, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Pivot Note (2026-11-26):** BYOC spend telemetry and `/api/v1/cloud/costs/*` retired. Live analytics: affiliate yield, funnels, search trends, deal performance. See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0435 ~ WO-CS-0450, WO-CS-0685 ~ WO-CS-0700`

---

## 1. 기능 목적 (Functional Purpose)

Provide real-time and historical business intelligence on affiliate revenue yield (EPC, RPM, conversion rate by vendor/category/page/country), search demand trends, zero-result queries, and deal performance (drop depth, coupon take-rate, alert conversion).

---

## 2. User Roles & Actors

Affiliate Business Analyst, Product Manager, Admin Operator, End User (Personal Deal Analytics).

---

## 3. User Scenarios

1. Affiliate Ops Manager analyzes which comparison page generates the highest EPC: discovers `/ai-video-editors` drives $420 EPC on annual-plan links while `/deals` drives high clicks with broader low-ticket volume.
2. Search Product Manager reviews the top 30 zero-result natural language queries in `search_events` (e.g., users searching for `ai podcast editor with lifetime deal`) to prioritize new catalog categories.

---

## 4. Technology Stack

PostgreSQL analytical materialized views (`mv_daily_affiliate_yield`, `mv_provider_funnel_stats`), Go Analytics Service (`internal/domain/analytics`), privacy-first first-party event collector (zero third-party ad-tracker cookies required).

---

## 3. Database Schema & Tables

**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`affiliate_clicks`, `affiliate_conversions`, `commissions`, `search_events`, `price_snapshots`

- (`cloud_usage_snapshots` retired with BYOC; replaced by `price_snapshots` for deal-trend analytics.)
- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)

**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**

- `GET /api/v1/analytics/affiliate-yield`
- `GET /api/v1/analytics/search-trends`
- `GET /api/v1/analytics/deal-performance`
- `GET /api/v1/analytics/deal-forecast`

- (`GET /api/v1/cloud/costs/summary` and `/api/v1/cloud/costs/forecast` are RETIRED.)
- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. User Interface & Experience

Admin Yield & Funnel Analytics Dashboard (`/admin/analytics`) and User Deal Watchlist & Savings Dashboard (`/dashboard/savings`).

---

## 8. Core Processing Logic

1. Key Yield Metrics:
   - `CTR = OutboundAffiliateClicks / ProductImpressions`
   - `CVR = VerifiedConversions / ValidHumanClicks`
   - `EPC = (ApprovedCommissionsUSD / ValidHumanClicks) * 100`
   - `NetRealizedRevenue = PaidCommissionsUSD - ClawbacksUSD` (strictly excluding account-credit referrals from cash KPIs).
2. User Deal Savings Tracker: Computes realized savings across alerted purchases (`SUM(median_90d_price - paid_price)`) and projects upcoming renewal exposure from tracked annual plans.

---

## 9. Security & Access Control

First-party analytics complies with GDPR/PIPA/CCPA without invasive cross-site tracking; user deal analytics are strictly isolated per organization.

---

## 10. Failure Modes & Resilience

Materialized views refresh concurrently (`REFRESH MATERIALIZED VIEW CONCURRENTLY`) every 15 minutes so heavy analytics queries never lock transactional tables.

---

## 11. Inter-Module Dependencies

Consumes data from Leads & Attribution (`L`), Search (`O`), and Deal/Estimation Engine (`E`). Multi-Cloud Management (`M`) is RETIRED.

---

## 12. Testing & Verification Strategy

SQL aggregation accuracy tests for EPC/CVR/Clawback calculations and deal-savings/renewal-exposure forecasting tests.

---

## 13. Acceptance Criteria

- [x] Accurate separation of Pending, Approved, Paid, and Reversed commissions in all yield reports.
- [x] Accurate separation of Cash Commissions vs. Account Credits.
- [x] Sub-100ms dashboard loading via concurrent materialized views.

---

## 14. Required Verification Evidence

Analytics SQL view definitions, unit test outputs, and OpenAPI analytics schema.

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0435 ~ WO-CS-0450, WO-CS-0685 ~ WO-CS-0700`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
