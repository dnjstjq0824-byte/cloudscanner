# J — Jobs — Distributed Task Queue, Price Sync Schedulers, Retries, DLQ & Idempotency

**Document ID:** `CS-AZ-J`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0801 ~ WO-CS-0825`

---

## 1. 기능 목적 (Functional Purpose)
Architect a high-reliability background job processing engine that orchestrates scheduled price synchronization, FX updates, affiliate conversion reconciliation, price-drop alert delivery, and asynchronous BYOC server provisioning/lifecycle polling with strict idempotency and dead-letter queue (DLQ) recovery.

---

## 2. 사용자 역할 (User Roles & Actors)
Backend Agent, Cloud Management Agent, SRE.

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. User submits a DigitalOcean Droplet creation request; the HTTP handler validates permissions, records an idempotent job in `background_jobs` and `resource_operations`, enqueues a high-priority `cloud:provision` task, and returns `202 Accepted` with an operation ID.
2. The worker polls DigitalOcean `GET /v2/droplets/{id}` every 5s until status transitions from `new` to `active`, then emits a real-time SSE/WebSocket event to the user dashboard.
3. A price sync job encounters a transient network reset; Asynq retries with exponential backoff up to 5 times, and if still failing, moves the task to the DLQ and creates an Admin alert.

---

## 4. 기술 스택 (Technology Stack)
Go `hibiken/asynq` (Redis-backed distributed queue) + PostgreSQL durable job state table (`background_jobs`) + Transactional Outbox relay (`outbox_events`).

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`background_jobs`, `outbox_events`, `price_sync_jobs`, `resource_operations`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)
**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**
- `GET /api/v1/admin/jobs`
- `POST /api/v1/admin/jobs/{job_id}/retry`
- `DELETE /api/v1/admin/jobs/{job_id}/cancel`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. UI/UX 설계 (User Interface & Experience)
User Operation Progress Drawer (real-time step progress for server creation/resizing) and Admin Job Queue Inspector (`/admin/jobs`) showing queue depths (`critical`, `cloud_ops`, `pricing_sync`, `notifications`, `dlq`).

---

## 8. 상세 처리 로직 (Core Processing Logic)
1. Priority Queues: `critical` (weight 6: token revocations, security alerts), `cloud_ops` (weight 4: BYOC provisioning, power actions, reconciliation), `pricing_sync` (weight 2: catalog/FX updates), `default` (weight 1: email/webhook alerts).
2. Strict Idempotency Guard: Every job carries a deterministic `idempotency_key`. Before calling any external API, the worker acquires a PostgreSQL advisory lock or checks `resource_operations.status` to guarantee a task is never executed twice.

---

## 9. 보안 통제 (Security & Access Control)
Job payloads in Redis never contain plaintext API tokens; workers fetch and decrypt KMS-wrapped credentials from PostgreSQL only in memory during execution and zero the byte slice (`crypto/subtle` / explicit wipe) immediately after.

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
Special rule for `cloud:provision` timeout: NEVER auto-retry a timed-out `POST /servers` call blindly. Always execute `ReconcileBeforeRetry` (query provider list/tags by `cloudscanr_op_id`) first to check whether the server was actually created before the connection dropped.

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Depends on PostgreSQL (`P_POSTGRESQL.md`) and Redis; powers Data Acquisition (`D`) and Multi-cloud Management (`M`).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
Idempotency race-condition tests (50 concurrent workers attempting the same `idempotency_key`), timeout-reconciliation simulation tests, and DLQ replay tests.

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] Zero duplicate server creation when network timeouts occur after provider receives `POST` request.
- [x] Zero plaintext credentials stored in Redis job payloads.
- [x] Full visibility and manual retry/cancel controls in Admin Console.

---

## 14. 검증 증빙 (Required Verification Evidence)
Asynq worker test logs, concurrent idempotency test output, and PostgreSQL `background_jobs` state transition traces.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0801 ~ WO-CS-0825`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
