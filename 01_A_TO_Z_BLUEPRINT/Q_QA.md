# Q — Quality Assurance — Unit, Integration, Contract, E2E, Performance & 22-Scenario Cloud Safety Testing

**Document ID:** `CS-AZ-Q`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0851 ~ WO-CS-0950`

---

## 1. 기능 목적 (Functional Purpose)
Establish a comprehensive, deterministic verification pyramid ensuring every Work Order (`WO-CS-0001` to `WO-CS-1000`) is backed by automated Unit, Integration, OpenAPI Contract, Playwright E2E, k6 Load, Security, and Cloud Provisioning Safety tests.

---

## 2. 사용자 역할 (User Roles & Actors)
QA Agent, Security Agent, Independent Audit Agent.

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. CI pipeline runs `go test -race ./...` and `vitest` on every commit, failing if domain coverage drops below 85%.
2. Cloud Provisioning Safety Suite executes all 22 mandatory edge cases from Article 24 (invalid key, expired token, read-only token attempting server creation, unsupported region/plan, insufficient balance, quota exceeded, provider 503, request timeout + reconciliation, lost response, duplicate request, user approval cancellation, server deletion gate, cross-user IDOR attempt, mid-flight token revocation, mid-flight account disconnect, resize cost change, backup failure, budget exceeded, orphan resource reconciliation, and audit log verification).

---

## 4. 기술 스택 (Technology Stack)
Go `testing` + `testify` + `testcontainers-go` (real PostgreSQL 16, Redis 7.4, OpenSearch 2.17 in CI), WireMock / Go `httptest` provider simulators, Playwright (TypeScript), k6 load testing, `oasdiff` / `vacuum` OpenAPI linter.

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`audit_logs`, `resource_operations`, `resource_reconciliations`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)
**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**
- `All 22 OpenAPI groups verified via contract test suite`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. UI/UX 설계 (User Interface & Experience)
Every screen in `05_UXUI/SCREEN_INVENTORY.csv` has dedicated Playwright desktop and mobile E2E assertions.

---

## 8. 상세 처리 로직 (Core Processing Logic)
1. Three-Tier Provider API Testing:
   - **Tier A (Deterministic Simulator)**: Local Go `httptest` server simulating exact DigitalOcean, UpCloud, Vultr, Hetzner, Linode, Cloudways v2 responses and fault injections (latency, 429, 500, dropped TCP connection after 201 Created).
   - **Tier B (Live Read-Only Verification)**: Safe catalog/region polling against live provider endpoints.
   - **Tier C (Authorized Live Sandbox Provisioning)**: Strictly gated by `ALLOW_LIVE_BILLABLE_TESTS=true` and `MAX_TEST_BUDGET_USD=5.00`, tearing down any created test VM within 120 seconds.

---

## 9. 보안 통제 (Security & Access Control)
Static analysis (`gosec`, `govulncheck`, `semgrep`, `trivy`), dynamic OWASP ZAP baseline scan, and automated IDOR matrix fuzzing.

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
Any test requiring external merchant approval or ungranted live billing authority is explicitly marked `BLOCKED_EXTERNAL_APPROVAL` in test reports rather than faked or skipped silently.

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Governs `13_TESTING/` and `WO-CS-0851 ~ WO-CS-0950`.

---

## 12. 테스트 전략 (Testing & Verification Strategy)
Self-validating test harness verifying test coverage, fixture integrity, and OpenAPI schema conformance.

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] All 22 mandatory server safety test cases in Article 24 have dedicated test specifications and automated harness cases.
- [x] Zero fake test logs or fabricated Git SHAs allowed.
- [x] 100% of OpenAPI 3.1 endpoints mapped to contract tests.

---

## 14. 검증 증빙 (Required Verification Evidence)
`13_TESTING/MASTER_TEST_MATRIX.csv`, `CLOUD_PROVISIONING_E2E.md`, and CI test runner logs.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0851 ~ WO-CS-0950`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
