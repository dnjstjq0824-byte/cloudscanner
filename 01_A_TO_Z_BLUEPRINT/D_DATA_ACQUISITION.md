# D — Data Acquisition — Official Pricing API Connectors, Change Detection & Licensing Compliance

**Document ID:** `CS-AZ-D`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0201 ~ WO-CS-0250`

---

## 1. 기능 목적 (Functional Purpose)
Design resilient, legal, and rate-limit-compliant data acquisition pipelines that harvest pricing and specification data exclusively from official Provider APIs, official public price books, and authorized feeds, with cryptographic payload hashing, diff detection, and staleness tracking.

---

## 2. 사용자 역할 (User Roles & Actors)
Provider API Agent, Data Engineering Lead, Operations Admin.

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. Hourly scheduled worker invokes DigitalOcean `GET /v2/sizes`, Vultr `GET /v2/plans`, Linode `GET /v4/linode/types`, Hetzner `GET /v1/pricing`, and UpCloud `GET /1.3/price`.
2. 6-hour scheduled worker streams AWS Bulk EC2 JSON index, GCP Cloud Billing Catalog SKUs, Azure Retail Prices REST API (`prices.azure.com`), OVHcloud Catalog API, and Oracle OCI Cost Estimator API.
3. A provider changes the monthly price of a SKU by >25% or returns an unexpected empty array; the Anomaly Circuit Breaker halts automatic overwrite, keeps the last verified snapshot, and alerts the Admin Console (`Price Discrepancy Queue`).

---

## 4. 기술 스택 (Technology Stack)
Go `internal/adapters/pricing` worker pool, Asynq cron scheduler, HTTP client with exponential backoff + jitter and ETag/If-Modified-Since caching, SHA-256 raw payload hashing, S3/R2 raw payload archive.

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`provider_data_sources`, `price_sync_jobs`, `price_snapshots`, `price_alerts`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)
**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**
- `POST /api/v1/admin/pricing/sync/{provider_id}`
- `GET /api/v1/admin/pricing/jobs`
- `GET /api/v1/admin/pricing/failures`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. UI/UX 설계 (User Interface & Experience)
Admin Price Sync Health Dashboard (`/admin/pricing-sync`) showing per-provider sync status, last success timestamp, SKU delta counts (+added / ~modified / -deprecated), and anomaly quarantine queue.

---

## 8. 상세 처리 로직 (Core Processing Logic)
1. Pull Pipeline: Fetch raw JSON -> compute SHA-256 hash -> if hash matches last run's hash, update `last_verified_at` timestamp without inserting duplicate rows -> if changed, parse via provider-specific `PricingAdapter`, validate bounds, insert new `price_snapshots` row, and emit `catalog.price_updated` event.
2. Triggering Price Alerts: When a SKU's normalized price drops below a user's threshold in `price_alerts`, enqueue a notification job.

---

## 9. 보안 통제 (Security & Access Control)
Read-only system credentials for authenticated pricing APIs (DigitalOcean, Hetzner, UpCloud, Cloudways v2, GCP) stored in KMS; strict outbound allowlist so pricing workers can only connect to verified provider API domains.

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
Exponential backoff (2s, 4s, 8s, 16s, 32s) on HTTP 429/5xx. If a provider API fails for >24 hours, all SKUs from that source automatically transition from `FRESH` to `STALE` in the UI with a visible warning banner.

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Depends on Catalog (`C_CATALOG.md`) and Jobs (`J_JOBS.md`); feeds Estimation Engine (`E_ESTIMATION_ENGINE.md`).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
Contract tests using recorded fixtures from real DigitalOcean, Vultr, Linode, Hetzner, UpCloud, Azure, and OVHcloud JSON responses; anomaly detector unit tests (simulating 0-price corruption or 90% drop).

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] Zero unauthorized scraping of competitor websites (`instances.vantage.sh`, `vpsbenchmarks.com`).
- [x] Every `price_snapshots` row records `source_url`, `collected_at` (UTC), `raw_payload_sha256`, and `verification_status`.
- [x] Cloudways adapter strictly uses API v2 (`/api/v2/server_sizes`) and rejects v1 endpoints.

---

## 14. 검증 증빙 (Required Verification Evidence)
Recorded JSON fixture test logs, `price_sync_jobs` execution records, and SHA-256 verification hashes.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0201 ~ WO-CS-0250`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
