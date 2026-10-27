# S — Security — KMS Envelope Encryption, BYOC Token Lifecycle, Tenant Isolation, Approval Gates & Audit Trail

**Document ID:** `CS-AZ-S`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0551 ~ WO-CS-0600, WO-CS-0826 ~ WO-CS-0840, WO-CS-0921 ~ WO-CS-0940`

---

## 1. 기능 목적 (Functional Purpose)
Enforce defense-in-depth security across CloudScanr, with special emphasis on protecting user BYOC cloud API credentials, preventing cross-tenant access (IDOR), stopping unauthorized cloud spend, and maintaining tamper-evident audit logs.

---

## 2. 사용자 역할 (User Roles & Actors)
Security Agent, Cloud Management Agent, Independent Audit Agent.

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. User submits a DigitalOcean or Hetzner API token; the backend generates a unique 256-bit Data Encryption Key (DEK) via KMS `GenerateDataKey`, encrypts the token with `AES-256-GCM` using `org_id || connection_id` as Additional Authenticated Data (AAD), stores only the KMS-encrypted DEK and ciphertext in `provider_credentials_metadata`, and zeroes the plaintext memory buffer.
2. Attacker compromises User A's session and attempts `GET /api/v1/cloud/servers/{user_b_server_id}` or `POST /api/v1/cloud/servers/{user_b_server_id}/actions`; both the Go middleware and PostgreSQL RLS reject the query, returning `404 Not Found` and logging a `SECURITY_TENANT_VIOLATION_ATTEMPT` event.

---

## 4. 기술 스택 (Technology Stack)
AWS KMS / HashiCorp Vault Transit Engine, Go `crypto/aes` + `crypto/cipher` (AES-256-GCM with AAD), PostgreSQL 16 Row-Level Security (RLS), Structured JSON Logger (`slog`) with automated secret redaction filter.

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`provider_credentials_metadata`, `cloud_connections`, `operation_approvals`, `audit_logs`, `sessions`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)
**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**
- `POST /api/v1/cloud/connections/{id}/rotate`
- `DELETE /api/v1/cloud/connections/{id}`
- `GET /api/v1/audit-logs`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. UI/UX 설계 (User Interface & Experience)
Security & Audit Log Console (`/dashboard/security-logs` and `/admin/security-events`) showing token usage timestamps, IP addresses, operation approvals, and one-click emergency `Revoke & Disconnect All Cloud Tokens` button.

---

## 8. 상세 처리 로직 (Core Processing Logic)
1. Envelope Encryption with AAD Binding: Passing `org_id:connection_id:provider_id` as GCM Additional Authenticated Data (AAD) makes it cryptographically impossible for an attacker with DB access to copy a ciphertext blob from one tenant row to another tenant row and decrypt it.
2. Cryptographic Operation Approval Token: When a user approves a billable operation (`CreateServer`, `ResizeServer`), the backend records a signed approval over `(user_id, org_id, connection_id, action, sku_id, region, estimated_monthly_cost, expires_at)` valid for 5 minutes.

---

## 9. 보안 통제 (Security & Access Control)
Zero-log policy for headers `Authorization`, `X-Api-Key`, `Cookie`, and body fields matching `*token*`, `*secret*`, `*password*`, `*key*`.

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
If KMS is unreachable, BYOC operations fail closed (`503 Service Unavailable — Security Subsystem Unreachable`) without ever caching decrypted keys on disk.

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Governs `09_SECURITY/` and all BYOC work orders (`WO-CS-0551 ~ WO-CS-0700`).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
AAD swap tampering test, cross-tenant RLS penetration test, log-redaction fuzz test (injecting simulated API tokens into error messages and verifying they are masked as `[REDACTED_SECRET]`), and expired approval token rejection test.

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] 100% of BYOC credentials encrypted via AES-256-GCM Envelope Encryption with tenant AAD binding.
- [x] Zero raw secrets ever written to PostgreSQL, Redis, OpenSearch, or application logs.
- [x] Tamper-evident hash chain (`prev_hash`, `entry_hash`) on `audit_logs`.

---

## 14. 검증 증빙 (Required Verification Evidence)
Security test suite logs (`09_SECURITY/SECURITY_TESTS.md`), KMS encryption unit tests, and log sanitizer test output.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0551 ~ WO-CS-0600, WO-CS-0826 ~ WO-CS-0840, WO-CS-0921 ~ WO-CS-0940`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
