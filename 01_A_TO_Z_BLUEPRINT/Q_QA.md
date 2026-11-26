# Q — Quality Assurance — Unit, Integration, Contract, E2E, Performance & Deal/Affiliate Safety Testing

**Document ID:** `CS-AZ-Q` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (TirpScanner pivot reframe, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Pivot Note (2026-11-26):** 22-scenario cloud-provisioning safety suite retired. Live QA targets the deal scanner: ingestion accuracy, plan math, affiliate attribution, and content/SEO. See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0851 ~ WO-CS-0950`

---

## 1. 기능 목적 (Functional Purpose)

Establish a comprehensive, deterministic verification pyramid ensuring every live TirpScanner behavior is backed by automated Unit, Integration, OpenAPI Contract, Playwright E2E, k6 Load, Security, and Deal/Affiliate Safety tests.

---

## 2. User Roles & Actors

QA Agent, Security Agent, Independent Audit Agent.

---

## 3. User Scenarios

1. CI pipeline runs `go test -race ./...` and `vitest` on every commit, failing if domain coverage drops below 85%.
2. Deal/Affiliate Safety Suite executes mandatory edge cases: feed timeout + STALE flip, 0-price/anomaly quarantine, expired-coupon suppression, `/go` open-redirect blocking, postback HMAC tamper rejection, duplicate conversion dedupe, clawback ledger balance, ranking independence from commission rates, FX-stale labeling, and audit log verification.

---

## 4. Technology Stack

Go `testing` + `testify` + `testcontainers-go` (real PostgreSQL 16, Redis 7.4, OpenSearch 2.17 in CI), WireMock / Go `httptest` feed simulators, Playwright (TypeScript), k6 load testing, `oasdiff` / `vacuum` OpenAPI linter.

---

## 3. Database Schema & Tables

**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`audit_logs`, `price_snapshots`, `affiliate_conversions`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)

**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**

- `All live OpenAPI groups verified via contract test suite (cloud groups retired)`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. User Interface & Experience

Every screen in `05_UXUI/SCREEN_INVENTORY.csv` has dedicated Playwright desktop and mobile E2E assertions.

---

## 8. Core Processing Logic

1. Three-Tier Feed API Testing:
   - **Tier A (Deterministic Simulator)**: Local Go `httptest` server simulating exact PartnerStack, Impact, CJ, Awin, and Product Hunt responses and fault injections (latency, 429, 30, truncated feed payload).
   - **Tier B (Live Read-Only Verification)**: Safe coupon/offer polling against live network endpoints.
   - **Tier C (Authorized Live Conversion Checks)**: Strictly gated test conversions reconciled against sandbox advertiser accounts, reversed/voided after verification.
2. Cloud-provisioning safety tiers (billable VM tests) are RETIRED.

---

## 9. Security & Access Control

Static analysis (`gosec`, `govulncheck`, `semgrep`, `trivy`), dynamic OWASP ZAP baseline scan, and automated IDOR matrix fuzzing.

---

## 10. Failure Modes & Resilience

Any test requiring external merchant approval or ungranted live billing authority is explicitly marked `BLOCKED_EXTERNAL_APPROVAL` in test reports rather than faked or skipped silently.

---

## 11. Inter-Module Dependencies

Governs `13_TESTING/` and `WO-CS-0851 ~ WO-CS-0950`.

---

## 12. Testing & Verification Strategy

Self-validating test harness verifying test coverage, fixture integrity, and OpenAPI schema conformance.

---

## 13. Acceptance Criteria

- [x] All mandatory deal/affiliate safety test cases have dedicated test specifications and automated harness cases.
- [x] Zero fake test logs or fabricated Git SHAs allowed.
- [x] 100% of live OpenAPI 3.1 endpoints mapped to contract tests.

---

## 14. Required Verification Evidence

`13_TESTING/MASTER_TEST_MATRIX.csv`, `DEAL_AFFILIATE_E2E.md`, and CI test runner logs.

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0851 ~ WO-CS-0950`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
