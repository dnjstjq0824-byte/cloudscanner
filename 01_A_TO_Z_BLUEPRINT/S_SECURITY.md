# S — Security — KMS Envelope Encryption, Credential Lifecycle, Tenant Isolation, Approval Gates & Audit Trail

**Document ID:** `CS-AZ-S` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (TirpScanner pivot reframe, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Pivot Note (2026-11-26):** BYOC user-cloud token flows retired. Live vault scope: affiliate-network API keys and system feed credentials. RLS/approval/audit mechanics unchanged. See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0551 ~ WO-CS-0600, WO-CS-0826 ~ WO-CS-0840, WO-CS-0921 ~ WO-CS-0940`

---

## 1. 기능 목적 (Functional Purpose)

Enforce defense-in-depth security across TirpScanner, with special emphasis on protecting affiliate-network and feed API credentials, preventing cross-tenant access (IDOR), stopping unauthorized payout or link changes, and maintaining tamper-evident audit logs.

---

## 2. User Roles & Actors

Security Agent, Affiliate Operations Agent, Independent Audit Agent.

---

## 3. User Scenarios

1. Admin submits a PartnerStack/Impact/CJ API key; the backend generates a unique 256-bit Data Encryption Key (DEK) via KMS `GenerateDataKey`, encrypts the key with `AES-256-GCM` using `org_id || credential_id` as Additional Authenticated Data (AAD), stores only the KMS-encrypted DEK and ciphertext, and zeroes the plaintext memory buffer.
2. Attacker compromises User A's session and attempts `GET /api/v1/admin/affiliate/payouts/{org_b_payout_id}`; both the Go middleware and PostgreSQL RLS reject the query, returning `404 Not Found` and logging a `SECURITY_TENANT_VIOLATION_ATTEMPT` event.

---

## 4. Technology Stack

AWS KMS / HashiCorp Vault Transit Engine, Go `crypto/aes` + `crypto/cipher` (AES-256-GCM with AAD), PostgreSQL 16 Row-Level Security (RLS), Structured JSON Logger (`slog`) with automated secret redaction filter.

---

## 3. Database Schema & Tables

**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`provider_credentials_metadata`, `affiliate_accounts`, `affiliate_payouts`, `audit_logs`, `sessions`

- `provider_credentials_metadata` now vaults affiliate-network/feed credentials (never user cloud accounts — BYOC retired). Any residual BYOC credential rows are quarantined for secure deletion, never used.
- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)

**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**

- `POST /api/v1/admin/affiliate/credentials/{id}/rotate`
- `DELETE /api/v1/admin/affiliate/credentials/{id}`
- `GET /api/v1/audit-logs`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. User Interface & Experience

Security & Audit Log Console (`/dashboard/security-logs` and `/admin/security-events`) showing credential usage timestamps, IP addresses, payout approvals, and one-click emergency `Revoke & Disconnect All Network Keys` button.

---

## 8. Core Processing Logic

1. Envelope Encryption with AAD Binding: Passing `org_id:credential_id:network_id` as GCM Additional Authenticated Data (AAD) makes it cryptographically impossible for an attacker with DB access to copy a ciphertext blob from one tenant row to another tenant row and decrypt it.
2. Cryptographic Operation Approval Token: When a user approves a payout or credential rotation, the backend records a signed approval over `(user_id, org_id, action, target_id, amount_or_scope, expires_at)` valid for 3 minutes.

---

## 9. Security & Access Control

Zero-log policy for headers `Authorization`, `X-Api-Key`, `Cookie`, and body fields matching `*token*`, `*secret*`, `*password*`, `*key*`.

---

## 10. Failure Modes & Resilience

If KMS is unreachable, credential-dependent operations fail closed (`503 Service Unavailable — Security Subsystem Unreachable`) without ever caching decrypted keys on disk.

---

## 11. Inter-Module Dependencies

Governs `09_SECURITY/` and affiliate-credential work orders. BYOC work orders (`WO-CS-0551 ~ WO-CS-0700`) are retired history.

---

## 12. Testing & Verification Strategy

AAD swap tampering test, cross-tenant RLS penetration test, log-redaction fuzz test (injecting simulated API tokens into error messages and verifying they are masked as `[REDACTED_SECRET]`), and expired approval token rejection test.

---

## 13. Acceptance Criteria

- [x] 100% of network/feed credentials encrypted via AES-256-GCM Envelope Encryption with tenant AAD binding.
- [x] Zero raw secrets ever written to PostgreSQL, Redis, OpenSearch, or application logs.
- [x] Tamper-evident hash chain (`prev_hash`, `entry_hash`) on `audit_logs`.

---

## 14. Required Verification Evidence

Security test suite logs (`09_SECURITY/SECURITY_TESTS.md`), KMS encryption unit tests, and log sanitizer test output.

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0551 ~ WO-CS-0600, WO-CS-0826 ~ WO-CS-0840, WO-CS-0921 ~ WO-CS-0940`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
