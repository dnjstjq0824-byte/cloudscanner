# J — Jobs — Distributed Task Queue, Deal-Sync Schedulers, Retries, DLQ & Idempotency

**Document ID:** `CS-AZ-J` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (TirpScanner pivot reframe, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Pivot Note (2026-11-26):** `cloud:provision` and BYOC lifecycle polling retired. Live queues: deal sync, FX, reconciliation, alerts. See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0801 ~ WO-CS-0825`

---

## 1. 기능 목적 (Functional Purpose)

Architect a high-reliability background job processing engine that orchestrates scheduled deal-feed synchronization, FX updates, affiliate conversion reconciliation, price-drop alert delivery, and content publishing with strict idempotency and dead-letter queue (DLQ) recovery.

---

## 2. User Roles & Actors

Backend Agent, Deal-Feed Agent, SRE.

---

## 3. User Scenarios

1. User subscribes to a price-drop alert; the HTTP handler validates permissions, records an idempotent subscription, and returns `202 Accepted` with an alert ID.
2. The deal-scan worker detects a 35 price drop on a tracked plan, enqueues a high-priority `alerts:price_drop` task, and the notification worker delivers email/webhook alerts with dedupe by `(alert_id, snapshot_id)`.
3. A deal-sync job encounters a transient network reset; Asynq retries with exponential backoff up to 3 times, and if still failing, moves the task to the DLQ and creates an Admin alert.

---

## 4. Technology Stack

Go `hibiken/asynq` (Redis-backed distributed queue) + PostgreSQL durable job state table (`background_jobs`) + Transactional Outbox relay (`outbox_events`).

---

## 3. Database Schema & Tables

**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`background_jobs`, `outbox_events`, `price_sync_jobs`, `resource_operations`

- `resource_operations` is RETIRED (BYOC-era table, name historical); live deal/alert pipelines use `background_jobs` + `price_sync_jobs`.
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

## 7. User Interface & Experience

User Alert Delivery Status (alert history with sent/opened state) and Admin Job Queue Inspector (`/admin/jobs`) showing queue depths (`critical`, `deal_sync`, `reconciliation`, `notifications`, `dlq`).

---

## 8. Core Processing Logic

1. Priority Queues: `critical` (weight 6: credential revocations, security alerts), `deal_sync` (weight 4: catalog/deal/FX updates), `reconciliation` (weight 2: affiliate conversion matching), `default` (weight 1: email/webhook alerts).
2. Strict Idempotency Guard: Every job carries a deterministic `idempotency_key`. Before calling any external API, the worker acquires a PostgreSQL advisory lock or checks durable job state to guarantee a task is never executed twice.

---

## 9. Security & Access Control

Job payloads in Redis never contain plaintext API tokens; workers fetch and decrypt KMS-wrapped credentials from PostgreSQL only in memory during execution and zero the byte slice (`crypto/subtle` / explicit wipe) immediately after.

---

## 10. Failure Modes & Resilience

Special rule for alert fan-out: NEVER send duplicate alerts on worker retry. Always check `(alert_id, snapshot_id)` dedupe state before dispatching, so a retried job delivers exactly once.

---

## 11. Inter-Module Dependencies

Depends on PostgreSQL (`P_POSTGRESQL.md`) and Redis; powers Data Acquisition (`D`) and Deal/Estimation Engine (`E`). Multi-cloud Management (`M`) is RETIRED.

---

## 12. Testing & Verification Strategy

Idempotency race-condition tests (30 concurrent workers attempting the same `idempotency_key`), alert-dedupe simulation tests, and DLQ replay tests.

---

## 13. Acceptance Criteria

- [x] Zero duplicate alert deliveries when network timeouts or worker retries occur.
- [x] Zero plaintext credentials stored in Redis job payloads.
- [x] Full visibility and manual retry/cancel controls in Admin Console.

---

## 14. Required Verification Evidence

Asynq worker test logs, concurrent idempotency test output, and PostgreSQL `background_jobs` state transition traces.

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0801 ~ WO-CS-0825`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
