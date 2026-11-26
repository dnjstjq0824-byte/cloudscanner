# CLOUD_MANAGEMENT_CONTRACT.md — BYOC Provisioning, Approval & Safety API Contract

**Document ID:** `CS-API-005`

> **RETIRED 2026-11-26 — TirpScanner pivot.** The entire `/api/v1/cloud/*` surface (connections, approvals, servers, budgets, costs) is REMOVED scope. This document is preserved below as dated history and must not be implemented or extended.

---

## 1. Two-Phase Approval & Execution Protocol (Articles 13 & 21) — RETIRED HISTORY

No mutating billable or destructive endpoint (`POST /api/v1/cloud/servers`, `DELETE /api/v1/cloud/servers/{id}`) can be invoked directly without a pre-signed `approval_id`:

### Phase A: Request Cost / Destructive Approval (`POST /api/v1/cloud/approvals`)
- Requires `BearerAuth` + `X-Step-Up-MFA-Token` + `X-Organization-Id`.
- Backend calculates exact hourly/monthly cost and checks `budget_policies`.
- Returns `approval_id` valid for **300 seconds (5 minutes)** bound to exact `(connection_id, operation_type, sku_id, region_code, acknowledged_monthly_cost_usd)`.

### Phase B: Execute Idempotent Operation (`POST /api/v1/cloud/servers`)
- Requires `Idempotency-Key` header + `approval_id` in body.
- Atomically marks `operation_approvals.consumed_at = NOW()` and inserts `resource_operations` in state `APPROVED_QUEUED`.
- Returns `202 Accepted` with `operation_id` and polling URI `/api/v1/cloud/operations/{operation_id}`.

---

## 2. Retirement Addendum (2026-11-26)
TirpScanner stores no user cloud credentials and offers no provisioning or server management. All `/api/v1/cloud/*` paths are marked `deprecated: true` / REMOVED in `openapi.yaml`, their test cases (`ATC-14`–`ATC-19`) expect `404`, and the backing tables (`cloud_connections`, `provider_credentials_metadata`, `cloud_resources`, `resource_operations`, `operation_approvals`, `cloud_usage_snapshots`, `budget_policies`, `resource_reconciliations`) are retired legacy with no live reads/writes.
