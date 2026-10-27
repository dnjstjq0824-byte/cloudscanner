# B — Business & Affiliate — Monetization Models, Commission Accounting, Individual Eligibility & Reseller Boundaries

**Document ID:** `CS-AZ-B`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0001 ~ WO-CS-0050, WO-CS-0351 ~ WO-CS-0450`

---

## 1. 기능 목적 (Functional Purpose)
Define CloudScanr's multi-stream revenue engine (Affiliate CPA/RevShare, Sponsored Placement with strict editorial separation, Pro SaaS BYOC alerting/team features, and future Contracted Reseller expansion), double-entry commission accounting, and practical operational rules for individual Korean developers vs. corporate entities.

---

## 2. 사용자 역할 (User Roles & Actors)
Affiliate Business Analyst, Finance & Reconciliation Admin, Partner Manager, Platform Operator.

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. Operator registers individual publisher accounts for DigitalOcean (Impact), Cloudways (Hybrid $30 + 7% recurring or Slab up to $125), UpCloud, Contabo (CJ), Hostinger, Vultr, and Kamatera using `cloudscanr.com` and W-8BEN.
2. Hetzner referrals generate €10 cloud hosting credits; the system records them under `reward_type = HOSTING_CREDIT` rather than cash revenue.
3. A user clicks a Cloudways affiliate link; the click is recorded, later matched via CSV/API reconciliation as `PENDING`, approved after the 45-day lock period as `APPROVED`, and settled in a monthly payout batch as `PAID`.
4. A cloud provider requests a sponsored banner on the GPU comparison page; the ad is rendered with an explicit `Sponsored` disclosure badge and excluded from altering the algorithmic `value_score`.

---

## 4. 기술 스택 (Technology Stack)
Go `shopspring/decimal` for exact financial arithmetic, PostgreSQL double-entry ledger tables (`commissions`, `commission_adjustments`, `affiliate_payouts`), CSV/JSON reconciliation parsers for Impact, CJ, PartnerStack, and HasOffers.

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
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

## 7. UI/UX 설계 (User Interface & Experience)
Public FTC/KFTC Affiliate Disclosure banner on all comparison/product pages and `/affiliate-disclosure` policy page. Admin Affiliate Revenue Cockpit displaying EPC (Earnings Per 100 Clicks), Conversion Rate, Pending vs. Approved vs. Paid Commission Pipeline, and Cash vs. Hosting-Credit separation.

---

## 8. 상세 처리 로직 (Core Processing Logic)
1. Three-Model Strict Separation (Article 14):
   - **Affiliate**: Outbound referral link -> user signs up & pays provider directly -> CloudScanr earns CPA/RevShare.
   - **BYOC**: User connects their own existing cloud account -> CloudScanr executes authorized API calls -> provider bills user directly.
   - **Reseller**: Disabled in V1 until formal Reseller Agreement, tax/billing compliance, and Tier-1 support desk are established.
2. Commission State Machine: `CLICKED` -> `CONVERSION_PENDING` -> `APPROVED` (or `REVERSED`/`CLAWBACK`) -> `SETTLED_PAID`.
3. Currency & Tax Handling: Preserves original network payout currency (`USD`, `EUR`, `GBP`) alongside normalized `USD` and `KRW` reporting equivalents at settlement FX rate.

---

## 9. 보안 통제 (Security & Access Control)
Immutable append-only ledger for `commissions` and `commission_adjustments` (PostgreSQL trigger prevents `UPDATE`/`DELETE` on finalized ledger entries; reversals must insert a compensating adjustment row). Role-based access restricted to `FINANCE_ADMIN` and `SUPER_ADMIN`.

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
Idempotent conversion ingestion keyed on `(program_id, external_conversion_id)`. Duplicate postbacks or overlapping CSV uploads are safely deduplicated and logged without double-counting revenue.

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Depends on Market Research (`00_MARKET_RESEARCH`) and PostgreSQL (`P_POSTGRESQL.md`); feeds into Leads & Attribution (`L_LEADS_AFFILIATE.md`) and Analytics (`Y_ANALYTICS.md`).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
Unit tests for Slab vs. Hybrid commission calculation, clawback adjustment ledger balancing, cash vs. credit segregation, and idempotent CSV reconciliation imports.

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] Every provider's affiliate status (`VERIFIED`, `REQUIRES_APPLICATION`, `CONTRACT_REQUIRED`, `NOT_SUPPORTED`) matches `00_MARKET_RESEARCH/AFFILIATE_PROGRAMS.csv`.
- [x] Hetzner credit referrals are strictly segregated from cash revenue totals.
- [x] Ledger sum of `commissions` + `commission_adjustments` equals exact reconciled payout batch totals with zero floating-point rounding drift.

---

## 14. 검증 증빙 (Required Verification Evidence)
PostgreSQL ledger constraint test logs, sample reconciled Impact/CJ CSV import audit report, and rendered FTC/KFTC disclosure screenshots.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0001 ~ WO-CS-0050, WO-CS-0351 ~ WO-CS-0450`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
