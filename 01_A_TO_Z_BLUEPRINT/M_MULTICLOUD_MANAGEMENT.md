# M — Multi-Cloud Management — RETIRED (BYOC/Provisioning Removed by 2026-11-26 Pivot)

**Document ID:** `CS-AZ-M` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (marked RETIRED, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Status:** **RETIRED** — BYOC account connection, server provisioning, day-2 server operations (power/resize/backup/firewall/SSH/IP/delete), cost-approval gates, and `08_CLOUD_MANAGEMENT/` + `/api/v1/cloud/*` + console server pages are retired. No new connections, no new work. Management adapter (`internal/adapters/management`) is frozen. See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0551 ~ WO-CS-0700`

---

## 1. 기능 목적 (Functional Purpose)

HISTORICAL (CloudScanr era): Design the complete Bring-Your-Own-Cloud (BYOC) account connection and server lifecycle management engine across supported cloud providers, enforcing strict tenant isolation, explicit cost approval, idempotent provisioning, timeout reconciliation, and destructive-action confirmation. This purpose is retired; the sections below preserve the original spec for traceability.

---

## 2. User Roles & Actors

HISTORICAL: Cloud Management Agent, Security Reviewer, Authenticated User / Org Operator. No live actors.

---

## 3. User Scenarios

HISTORICAL (retired, preserved in intent):

1. User connects their DigitalOcean account via OAuth 2.0 (or Hetzner/Vultr/UpCloud/Cloudways v2 via scoped API token), selects `Read-Only` or `Read+Write` mode, and verifies account ownership.
2. User selects a Droplet plan, chooses OS image and SSH key, reviews the itemized cost modal, signs the explicit approval, and provisions the server idempotently.
3. On `POST /servers` timeout, the worker reconciles by provider list/tags (`cloudscanr-op-id`, historical tag key) before any retry, linking `external_resource_id` instead of creating a duplicate billable server.

---

## 4. Technology Stack

HISTORICAL: Go `internal/adapters/management` (per-provider server adapters), KMS Envelope Encryption, Asynq `cloud_ops` state machine. Frozen — no new providers, no new calls.

---

## 3. Database Schema & Tables

RETIRED tables (historical; see `P_POSTGRESQL.md`):  
`cloud_connections`, `provider_credentials_metadata`, `cloud_resources`, `resource_operations`, `operation_approvals`, `cloud_usage_snapshots`, `budget_policies`, `resource_reconciliations`

- No new rows for live flows. Retention/deletion of any stored user-cloud credentials is a security-handled migration concern (see `S_SECURITY.md`).

---

## 6. API 구조 (API Contracts & Endpoints)

RETIRED endpoints (do not implement; remove from live OpenAPI surface):

- `POST /api/v1/cloud/connections`
- `POST /api/v1/cloud/connections/{id}/verify`
- `DELETE /api/v1/cloud/connections/{id}`
- `GET /api/v1/cloud/servers`
- `POST /api/v1/cloud/approvals`
- `POST /api/v1/cloud/servers`
- `POST /api/v1/cloud/servers/{id}/actions`
- `DELETE /api/v1/cloud/servers/{id}`

---

## 7. User Interface & Experience

RETIRED pages: Connected Accounts Console (`/dashboard/connections`), Server Provisioning Wizard (`/dashboard/servers/new`), Server Inventory (`/dashboard/servers`), Server Detail Console (`/dashboard/servers/[id]`). Do not build.

---

## 8. Core Processing Logic

HISTORICAL (retired): provisioning state machine (`DRAFT` -> `AWAITING_COST_APPROVAL` -> `APPROVED_QUEUED` -> `PREFLIGHT_VERIFYING` -> `DISPATCHING_TO_PROVIDER` -> `POLLING_ACTIVE`/`TIMEOUT_RECONCILING` -> `SUCCEEDED`/`FAILED_SAFE`); mandatory `cloudscanr-managed` resource tagging (historical key); destructive-operation double-gate (typed hostname + step-up MFA + approval record).

---

## 9. Security & Access Control

HISTORICAL: cross-tenant RLS isolation, read-only connection guards. Live credential-vault rules for affiliate-network keys continue under `S_SECURITY.md`; BYOC credential flows are retired.

---

## 10. Failure Modes & Resilience

HISTORICAL: 22 cloud-safety scenarios (expired token, insufficient funds, quota exceeded, mid-operation disconnect, drift reconciliation). Retired with this module (see `Q_QA.md`).

---

## 11. Inter-Module Dependencies

Formerly depended on Identity (`I`), Security (`S`), Jobs (`J`), and `08_CLOUD_MANAGEMENT/PROVIDER_CAPABILITIES.md`. No live dependents; nothing may newly depend on this module.

---

## 12. Testing & Verification Strategy

HISTORICAL: mock-provider fault-injection suite + sandbox live tests. Retired; preserved as frozen work-order evidence.

---

## 13. Acceptance Criteria

- [x] Module formally marked RETIRED with live/retired scope recorded in `TIRPSCANNER_PIVOT.md`.
- [x] No live OpenAPI group, page, job, or role references BYOC/provisioning/server management.
- [x] Historical spec preserved above for audit traceability.

---

## 14. Required Verification Evidence

Pivot decision record (`TIRPSCANNER_PIVOT.md`); frozen work-order history (`11_WORK_ORDERS_0001_TO_1000/`, read-only).

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0551 ~ WO-CS-0700`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
