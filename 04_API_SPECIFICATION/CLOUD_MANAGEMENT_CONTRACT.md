# CLOUD_MANAGEMENT_CONTRACT.md — BYOC Provisioning, Approval & Safety API Contract

**Document ID:** `CS-API-005`

---

## 1. Two-Phase Approval & Execution Protocol (Articles 13 & 21)

No mutating billable or destructive endpoint (`POST /api/v1/cloud/servers`, `DELETE /api/v1/cloud/servers/{id}`) can be invoked directly without a pre-signed `approval_id`:

### Phase A: Request Cost / Destructive Approval (`POST /api/v1/cloud/approvals`)
- Requires `BearerAuth` + `X-Step-Up-MFA-Token` + `X-Organization-Id`.
- Backend calculates exact hourly/monthly cost and checks `budget_policies`.
- Returns `approval_id` valid for **300 seconds (5 minutes)** bound to exact `(connection_id, operation_type, sku_id, region_code, acknowledged_monthly_cost_usd)`.

### Phase B: Execute Idempotent Operation (`POST /api/v1/cloud/servers`)
- Requires `Idempotency-Key` header + `approval_id` in body.
- Atomically marks `operation_approvals.consumed_at = NOW()` and inserts `resource_operations` in state `APPROVED_QUEUED`.
- Returns `202 Accepted` with `operation_id` and polling URI `/api/v1/cloud/operations/{operation_id}`.
