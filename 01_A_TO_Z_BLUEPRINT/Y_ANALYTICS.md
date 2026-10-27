# Y — Yield & Analytics — Affiliate Revenue Attribution, Funnel Conversion, EPC & Cloud Spend Telemetry

**Document ID:** `CS-AZ-Y`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0435 ~ WO-CS-0450, WO-CS-0685 ~ WO-CS-0700`

---

## 1. 기능 목적 (Functional Purpose)
Provide real-time and historical business intelligence on affiliate revenue yield (EPC, RPM, conversion rate by provider/category/page/country), search demand trends, zero-result queries, and aggregated BYOC infrastructure spend.

---

## 2. 사용자 역할 (User Roles & Actors)
Affiliate Business Analyst, Product Manager, Admin Operator, End User (Personal Cost Analytics).

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. Affiliate Ops Manager analyzes which comparison page generates the highest EPC: discovers `/gpu-comparison` drives $420 EPC on Vultr/DigitalOcean GPU links while `/vps-comparison` under $5 drives high clicks but higherContabo/Hostinger volume.
2. Search Product Manager reviews the top 50 zero-result natural language queries in `search_events` (e.g., users searching for `Mac Mini M4 cloud server`) to prioritize new catalog categories.

---

## 4. 기술 스택 (Technology Stack)
PostgreSQL analytical materialized views (`mv_daily_affiliate_yield`, `mv_provider_funnel_stats`), Go Analytics Service (`internal/domain/analytics`), privacy-first first-party event collector (zero third-party ad-tracker cookies required).

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`affiliate_clicks`, `affiliate_conversions`, `commissions`, `search_events`, `cloud_usage_snapshots`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)
**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**
- `GET /api/v1/analytics/affiliate-yield`
- `GET /api/v1/analytics/search-trends`
- `GET /api/v1/cloud/costs/summary`
- `GET /api/v1/cloud/costs/forecast`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. UI/UX 설계 (User Interface & Experience)
Admin Yield & Funnel Analytics Dashboard (`/admin/analytics`) and User Cloud Cost & Forecast Dashboard (`/dashboard/costs`).

---

## 8. 상세 처리 로직 (Core Processing Logic)
1. Key Yield Metrics:
   - `CTR = OutboundAffiliateClicks / ProductImpressions`
   - `CVR = VerifiedConversions / ValidHumanClicks`
   - `EPC = (ApprovedCommissionsUSD / ValidHumanClicks) * 100`
   - `NetRealizedRevenue = PaidCommissionsUSD - ClawbacksUSD` (strictly excluding `HOSTING_CREDIT` referrals from cash KPIs).
2. User BYOC Cost Forecast: Computes current daily run-rate across all connected accounts (`SUM(active_server.hourly_rate * 24)`) and projects end-of-month spend against `budget_policies`.

---

## 9. 보안 통제 (Security & Access Control)
First-party analytics complies with GDPR/PIPA/CCPA without invasive cross-site tracking; user BYOC cost analytics are strictly isolated per organization.

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
Materialized views refresh concurrently (`REFRESH MATERIALIZED VIEW CONCURRENTLY`) every 15 minutes so heavy analytics queries never lock transactional tables.

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Consumes data from Leads & Attribution (`L`), Search (`O`), and Multi-Cloud Management (`M`).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
SQL aggregation accuracy tests for EPC/CVR/Clawback calculations and BYOC monthly burn-rate forecasting tests.

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] Accurate separation of Pending, Approved, Paid, and Reversed commissions in all yield reports.
- [x] Accurate separation of Cash Commissions vs. Hosting Credits.
- [x] Sub-100ms dashboard loading via concurrent materialized views.

---

## 14. 검증 증빙 (Required Verification Evidence)
Analytics SQL view definitions, unit test outputs, and OpenAPI analytics schema.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0435 ~ WO-CS-0450, WO-CS-0685 ~ WO-CS-0700`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
