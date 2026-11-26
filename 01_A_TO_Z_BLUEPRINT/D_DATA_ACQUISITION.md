# D — Data Acquisition — Deal-Feed Connectors, Price-Drop Detection & Licensing Compliance

**Document ID:** `CS-AZ-D` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (TirpScanner pivot reframe, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Pivot Note (2026-11-26):** Cloud pricing-API harvesting is RETIRED. Live ingestion: vendor pricing pages, affiliate-network feeds/APIs, coupon sources, launch feeds. See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0201 ~ WO-CS-0250`

---

## 1. 기능 목적 (Functional Purpose)

Design resilient, legal, and rate-limit-compliant data acquisition pipelines that harvest plan pricing, promos, and coupon data exclusively from official vendor pricing pages, official affiliate-network feeds/APIs, and authorized coupon/launch sources, with cryptographic payload hashing, diff detection, price-drop detection, and staleness tracking.

---

## 2. User Roles & Actors

Vendor Feed Agent, Data Engineering Lead, Operations Admin.

---

## 3. User Scenarios

1. Hourly scheduled worker polls affiliate-network APIs (PartnerStack, Impact, CJ, Awin) for coupon/promo deltas and conversion-capable offer changes.
2. 6-hour scheduled worker re-scans vendor pricing pages and curated deal sources (AppSumo, StackSocial, Product Hunt launches) for plan re-tiering, new promos, and renewal-term changes.
3. A vendor changes a plan price by >25% or a feed returns an unexpected empty array; the Anomaly Circuit Breaker halts automatic overwrite, keeps the last verified snapshot, and alerts the Admin Console (`Deal Discrepancy Queue`).

---

## 4. Technology Stack

Go `internal/adapters/pricing` worker pool (catalog/deal-ingestion; `pricing` package name historical), Asynq cron scheduler, HTTP client with exponential backoff + jitter and ETag/If-Modified-Since caching, SHA-256 raw payload hashing, S3/R2 raw payload archive.

---

## 5. Database Schema & Tables

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

## 7. User Interface & Experience

Admin Deal-Sync Health Dashboard (`/admin/pricing-sync`) showing per-source sync status, last success timestamp, plan delta counts (+added / ~modified / -deprecated), and anomaly quarantine queue.

---

## 8. Core Processing Logic

1. Pull Pipeline: Fetch raw feed/page -> compute SHA-256 hash -> if hash matches last run's hash, update `last_verified_at` timestamp without inserting duplicate rows -> if changed, parse via source-specific ingestion adapter, validate bounds, insert new `price_snapshots` row, and emit `catalog.price_updated` event.
2. Triggering Price Alerts: When a plan's normalized effective price drops below a user's threshold in `price_alerts` (or a new coupon matches), enqueue a notification job.

---

## 9. Security & Access Control

Read-only system credentials for authenticated network APIs (PartnerStack, Impact, CJ, Awin, FirstPromoter) stored in KMS; strict outbound allowlist so ingestion workers can only connect to verified feed domains.

---

## 10. Failure Modes & Resilience

Exponential backoff (2s, 4s, 8s, 16s, 32s) on HTTP 429/5xx. If a feed source fails for >24 hours, all plans from that source automatically transition from `FRESH` to `STALE` in the UI with a visible warning banner.

---

## 11. Inter-Module Dependencies

Depends on Catalog (`C_CATALOG.md`) and Jobs (`J_JOBS.md`); feeds Deal/Estimation Engine (`E_ESTIMATION_ENGINE.md`).

---

## 12. Testing & Verification Strategy

Contract tests using recorded fixtures from real PartnerStack, Impact, CJ, Awin, and Product Hunt API responses; anomaly detector unit tests (simulating 0-price corruption or 90% drop).

---

## 13. Acceptance Criteria

- [x] Zero unauthorized scraping of competitor review content (`G2`, `Capterra`, `TrustRadius`).
- [x] Every `price_snapshots` row records `source_url`, `collected_at` (UTC), `raw_payload_sha256`, and `verification_status`.
- [x] Every deal source honors `robots.txt`, rate limits, and feed terms of use.

---

## 14. Required Verification Evidence

Recorded JSON fixture test logs, `price_sync_jobs` execution records, and SHA-256 verification hashes.

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0201 ~ WO-CS-0250`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
