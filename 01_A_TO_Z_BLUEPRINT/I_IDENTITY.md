# I — Identity — Authentication, OIDC/Social Login, MFA/Passkeys, Sessions, Organizations & RBAC

**Document ID:** `CS-AZ-I`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0101 ~ WO-CS-0150`

---

## 1. 기능 목적 (Functional Purpose)
Implement a standalone, hardened Identity and Access Management (IAM) system for CloudScanr (`cloudscanr.com`) supporting email/password registration with email verification, Argon2id password hashing, TOTP/WebAuthn Passkey MFA, OAuth2/OIDC social login (GitHub, Google), device session management, multi-tenant Organizations, and fine-grained RBAC for cloud provisioning and budget approvals.

---

## 2. 사용자 역할 (User Roles & Actors)
Identity Engineer, Security Agent, Individual User, Organization Owner, Organization Admin, Billing/Budget Approver, Organization Viewer.

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. Developer registers with email, verifies ownership via a single-use SHA-256 hashed token, enables TOTP + WebAuthn Passkey, and creates an Organization (`Acme Infra Team`).
2. Organization Owner invites a Junior Engineer with the `ORG_VIEWER` role (can view connected servers and costs, cannot create or delete servers) and a Senior SRE with `CLOUD_OPERATOR` role (can provision servers within approved budget policies).
3. User requests GDPR/PIPA account export or deletion; the system exports personal data as JSON, revokes all active sessions, disconnects BYOC credentials, and anonymizes audit pointers.

---

## 4. 기술 스택 (Technology Stack)
Go `internal/domain/identity`, Argon2id (`m=64MB, t=3, p=4`), PASETO v4 / JWT short-lived access tokens (15m) + rotating refresh tokens in `HttpOnly` cookies, WebAuthn (`go-webauthn`), TOTP (`pquerna/otp`), OAuth2 (`golang.org/x/oauth2`).

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`users`, `organizations`, `organization_members`, `roles`, `sessions`, `auth_identities`, `audit_logs`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)
**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/verify-email`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/mfa/verify`
- `POST /api/v1/auth/passkeys/register`
- `POST /api/v1/orgs`
- `POST /api/v1/orgs/{org_id}/members`
- `DELETE /api/v1/users/me`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. UI/UX 설계 (User Interface & Experience)
Clean standalone Auth pages (`/register`, `/login`, `/forgot-password`, `/mfa-challenge`), User Security & Devices page (`/dashboard/security`), and Organization RBAC Matrix manager (`/dashboard/org/members`).

---

## 8. 상세 처리 로직 (Core Processing Logic)
1. Strict Role Matrix:
   - `ORG_OWNER`: Full org control, billing/budget approval, server deletion approval, member management.
   - `ORG_ADMIN`: Connect cloud accounts, approve server creation/deletion within org limits.
   - `CLOUD_OPERATOR`: Request/create servers, reboot/power actions; requires `operation_approvals` sign-off for destructive actions or over-budget plans.
   - `ORG_VIEWER`: Read-only access to connected cloud inventory, cost estimates, and comparison lists.
2. Step-Up Authentication: Connecting a new write-enabled BYOC token, deleting a server, or viewing security logs requires a fresh MFA/Passkey verification within the last 10 minutes (`step_up_verified_at`).

---

## 9. 보안 통제 (Security & Access Control)
Argon2id hashing; refresh token reuse detection (automatic family revocation if a rotated token is replayed); rate limiting (5 failed logins/15m per IP+account); strict separation between CloudScanr user identity and external cloud account ownership.

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
Graceful lockout recovery via email + backup recovery codes; atomic session revocation across Redis and PostgreSQL.

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Foundation for User Dashboard (`U_UXUI.md`), BYOC Connection (`M_MULTICLOUD_MANAGEMENT.md`), and Security (`S_SECURITY.md`).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
Unit and integration tests for Argon2id, token rotation replay detection, RBAC permission boundaries (`ORG_VIEWER` attempting `POST /api/v1/cloud/servers` must return `403 Forbidden`), and step-up MFA enforcement.

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] All 20 identity and organization capabilities required in Article 15 are fully specified and backed by DDL and OpenAPI endpoints.
- [x] Step-up MFA is enforced before any server deletion or write-token attachment.
- [x] Zero dependency on TreeScan user accounts or SSO.

---

## 14. 검증 증빙 (Required Verification Evidence)
Auth integration test logs, RBAC matrix test results, and OpenAPI auth security scheme validation.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0101 ~ WO-CS-0150`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
