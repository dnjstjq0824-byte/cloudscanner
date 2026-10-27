# E — Estimation Engine — 36-Dimension Price Normalization, Usage Scenarios, Bandwidth, Taxes & FX

**Document ID:** `CS-AZ-E`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0251 ~ WO-CS-0300`

---

## 1. 기능 목적 (Functional Purpose)
Implement the 36-attribute Cloud Price Normalization & TCO Estimation Engine mandated by Article 8. Separate raw provider rates from scenario-based estimated monthly/annual TCO, accounting for 730-hour vs. calendar-month billing, bandwidth egress overages, public IPv4 surcharges, backup/snapshot storage, OS licenses, promotional credits, VAT/GST, and real-time FX conversion.

---

## 2. 사용자 역할 (User Roles & Actors)
Pricing Engine Agent, Financial Analyst, End User Comparing Cloud TCO.

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. User compares an AWS `t4g.medium` ($0.0336/hr compute + $0.005/hr IPv4 + $0.08/GB EBS gp3 + $0.09/GB egress beyond 100GB) against a Hetzner `CAX21` (€6.49/mo cap including 20 TB traffic) and DigitalOcean `s-2vcpu-4gb` ($24/mo cap including 4 TB traffic and free IPv4).
2. User toggles between `730 Hours (Standard Mean Month)`, `672 Hours (28-Day Feb)`, `744 Hours (31-Day Month)`, and `Custom Hours (e.g., 160h Dev/Test Weekday Server)`.
3. User in South Korea selects `KRW (₩)` and `Include 10% VAT`; the engine converts USD/EUR raw rates using timestamped FX rates and clearly displays the FX snapshot timestamp and tax breakdown.

---

## 4. 기술 스택 (Technology Stack)
Go `internal/domain/estimation` using `shopspring/decimal` (12 decimal places internal precision, Half-Even rounding for display), PostgreSQL `price_components`, `billing_terms`, `fx_rates`, `tax_policies`.

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
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

## 7. UI/UX 설계 (User Interface & Experience)
Interactive TCO Breakdown Drawer & Pricing Calculator (`/calculator`) displaying two distinct columns: **Raw Provider Price** (exact currency & unit billed by provider) vs. **CloudScanr Estimated TCO** (itemized stacked bar: Compute + Storage + IPv4 + Egress Overage + Backups + OS License + Tax).

---

## 8. 상세 처리 로직 (Core Processing Logic)
Formula for Scenario Monthly Estimate:
`BaseCompute = MIN(HourlyRate * UsageHours, MonthlyCap)` (if provider enforces a monthly cap like DO/Vultr/Linode/UpCloud/Hetzner; if pure hourly like AWS/GCP/Azure, `HourlyRate * UsageHours`).
`IPv4Cost = RequiredIPv4Count * (IPv4Included ? 0 : IPv4MonthlyOrHourlyRate)`.
`EgressCost = MAX(0, EstimatedEgressGB - IncludedBandwidthGB) * EgressOverageRatePerGB`.
`BackupCost = BackupEnabled ? (BaseCompute * BackupPct + SnapshotGB * SnapshotRatePerGB) : 0`.
`SubtotalNative = BaseCompute + ExtraStorageCost + IPv4Cost + EgressCost + BackupCost + OSLicenseCost + SetupFeeAmortized`.
`TotalTargetCurrency = SubtotalNative * FXRate(Native -> Target) * (1 + TaxRate)`.

---

## 9. 보안 통제 (Security & Access Control)
Deterministic calculation inputs validated against injection and overflow (`UsageHours` bounded `1..8784`, `EstimatedEgressGB` bounded `0..10000000`).

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
If live FX rate feed is older than 24 hours, display `FX Rate Last Updated: YYYY-MM-DD HH:MM UTC (Stale Warning)` alongside the native currency price.

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Depends on Catalog (`C_CATALOG.md`) and Data Acquisition (`D_DATA_ACQUISITION.md`).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
Table-driven Go unit tests covering all 36 normalization fields, 730h vs. 672h/744h edge cases, monthly price cap enforcement, IPv6-only discount on Hetzner, and AWS/GCP egress overage tiers.

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] All 36 attributes listed in Article 8 are modeled in the database and estimation payload.
- [x] Raw provider billing price and CloudScanr estimated scenario price are never conflated in API responses or UI.
- [x] 730-hour assumption is explicitly labeled and configurable by the user.

---

## 14. 검증 증빙 (Required Verification Evidence)
Go decimal unit test results (`go test ./internal/domain/estimation/...`), JSON output of `/api/v1/pricing/estimate`, and calculator UI verification.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0251 ~ WO-CS-0300`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
