# M — Multi-Cloud Management — BYOC Account Connection, Server Provisioning, Day-2 Operations & Safety Controls

**Document ID:** `CS-AZ-M`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0551 ~ WO-CS-0700`

---

## 1. 기능 목적 (Functional Purpose)
Design the complete Bring-Your-Own-Cloud (BYOC) account connection and server lifecycle management engine (Articles 12, 13, 21, 24) across supported cloud providers, enforcing strict tenant isolation, explicit cost approval, idempotent provisioning, timeout reconciliation, and destructive-action confirmation.

---

## 2. 사용자 역할 (User Roles & Actors)
Cloud Management Agent, Security Reviewer, Authenticated User / Org Operator.

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. User connects their DigitalOcean account via OAuth 2.0 (or Hetzner/Vultr/UpCloud/Cloudways v2 via scoped API token), selects `Read-Only` or `Read+Write` mode, and verifies account ownership via `GET /v2/account`.
2. User selects a DigitalOcean `s-2vcpu-4gb` Droplet in `fra1` ($24/mo), chooses Ubuntu 24.04 LTS and an SSH key, reviews the itemized cost modal, signs the explicit approval, and provisions the Droplet idempotently.
3. A network timeout occurs 10 seconds after sending `POST https://api.digitalocean.com/v2/droplets`. Instead of blindly re-POSTing and creating a duplicate billable Droplet, the worker queries `GET /v2/droplets?tag_name=cs-op-{op_id}`, discovers the server was created, links the `external_resource_id`, and transitions the operation to `SUCCEEDED`.

---

## 4. 기술 스택 (Technology Stack)
Go `internal/adapters/management` (DigitalOcean, UpCloud, Vultr, Hetzner, Linode, Cloudways v2, AWS STS, GCP, Azure, Scaleway, OVH, Contabo, Hostinger, Kamatera, IONOS, OCI, RunPod), KMS Envelope Encryption, Asynq `cloud_ops` state machine.

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`cloud_connections`, `provider_credentials_metadata`, `cloud_resources`, `resource_operations`, `operation_approvals`, `cloud_usage_snapshots`, `budget_policies`, `resource_reconciliations`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)
**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**
- `POST /api/v1/cloud/connections`
- `POST /api/v1/cloud/connections/{id}/verify`
- `DELETE /api/v1/cloud/connections/{id}`
- `GET /api/v1/cloud/servers`
- `POST /api/v1/cloud/approvals`
- `POST /api/v1/cloud/servers`
- `POST /api/v1/cloud/servers/{id}/actions`
- `DELETE /api/v1/cloud/servers/{id}`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. UI/UX 설계 (User Interface & Experience)
Connected Accounts Console (`/dashboard/connections`), Guided 4-Step Server Provisioning Wizard (`/dashboard/servers/new`) with mandatory Cost Approval Gate, Unified Server Inventory (`/dashboard/servers`), and Server Detail & Power/Backup/Firewall Console (`/dashboard/servers/[id]`).

---

## 8. 상세 처리 로직 (Core Processing Logic)
1. Provisioning State Machine: `DRAFT` -> `AWAITING_COST_APPROVAL` -> `APPROVED_QUEUED` -> `PREFLIGHT_VERIFYING` -> `DISPATCHING_TO_PROVIDER` -> `POLLING_ACTIVE` (or `TIMEOUT_RECONCILING`) -> `SUCCEEDED` / `FAILED_SAFE`.
2. Mandatory Resource Tagging: Every server created via CloudScanr is tagged with `cloudscanr-managed=true` and `cloudscanr-op-id={operation_uuid}` so timeout reconciliation can deterministically identify resources even if the HTTP response was lost.
3. Destructive Operation Double-Gate: `DeleteServer`, `RebuildServer`, and `DetachVolume` require (a) typing the exact server hostname, (b) step-up MFA verification, and (c) a valid `operation_approvals` record.

---

## 9. 보안 통제 (Security & Access Control)
Cross-tenant isolation enforced at API middleware AND PostgreSQL RLS (`organization_id = current_setting('app.current_org_id')::uuid`). Read-only connections physically reject any mutating adapter call before decrypting credentials.

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
Covers all 22 mandatory safety test scenarios defined in Article 24 (expired token, insufficient funds, quota exceeded, mid-operation account disconnect, drift reconciliation).

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Depends on Identity (`I`), Security (`S`), Jobs (`J`), and Provider Capabilities Matrix (`08_CLOUD_MANAGEMENT/PROVIDER_CAPABILITIES.md`).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
Mock provider HTTP server fault-injection suite testing all 22 Article 24 failure modes + sandbox account live integration tests.

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] All 21 minimum server management features in Article 13 are specified and mapped to provider capabilities.
- [x] All 22 safety scenarios in Article 24 pass deterministic automated tests.
- [x] Zero possibility of double server creation on HTTP timeout.

---

## 14. 검증 증빙 (Required Verification Evidence)
22-scenario safety test execution report, state machine transition logs, and KMS encryption audit logs.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0551 ~ WO-CS-0700`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
