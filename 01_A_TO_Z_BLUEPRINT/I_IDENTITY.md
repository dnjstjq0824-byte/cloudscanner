# I — Identity — Authentication, OIDC/Social Login, MFA/Passkeys, Sessions, Organizations & RBAC

**Document ID:** `CS-AZ-I` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (TirpScanner pivot reframe, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Pivot Note (2026-11-26):** Cloud-provisioning roles and BYOC step-up rules retired; RBAC reframed around deal curation, alerts, and affiliate ops. See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0101 ~ WO-CS-0150`

---

## 1. 기능 목적 (Functional Purpose)

Implement a standalone, hardened Identity and Access Management (IAM) system for TirpScanner (`tirpscanr.com`) supporting email/password registration with email verification, Argon2id password hashing, TOTP/WebAuthn Passkey MFA, OAuth2/OIDC social login (GitHub, Google), device session management, multi-tenant Organizations, and fine-grained RBAC for deal curation, alert management, and payout approvals.

---

## 2. User Roles & Actors

Identity Engineer, Security Agent, Individual User, Organization Owner, Organization Admin, Payout Approver, Organization Viewer.

---

## 3. User Scenarios

1. Developer registers with email, verifies ownership via a single-use SHA-256 hashed token, enables TOTP + WebAuthn Passkey, and creates an Organization (`Acme Deals Team`).
2. Organization Owner invites a Junior Analyst with the `ORG_VIEWER` role (can view alerts, comparisons, and deal reports; cannot curate catalog or approve payouts) and a Senior Editor with the `DEAL_CURATOR` role (can curate deals and publish buying guides within policy).
3. User requests GDPR/PIPA account export or deletion; the system exports personal data as JSON, revokes all active sessions, deletes alert subscriptions, and anonymizes audit pointers.

---

## 4. Technology Stack

Go `internal/domain/identity`, Argon2id (`m=64MB, t=3, p=4`), PASETO v4 / JWT short-lived access tokens (15m) + rotating refresh tokens in `HttpOnly` cookies, WebAuthn (`go-webauthn`), TOTP (`pquerna/otp`), OAuth2 (`golang.org/x/oauth2`).

---

## 3. Database Schema & Tables

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

## 7. User Interface & Experience

Clean standalone Auth pages (`/register`, `/login`, `/forgot-password`, `/mfa-challenge`), User Security & Devices page (`/dashboard/security`), and Organization RBAC Matrix manager (`/dashboard/org/members`).

---

## 8. Core Processing Logic

1. Strict Role Matrix:
   - `ORG_OWNER`: Full org control, payout approval, member management.
   - `ORG_ADMIN`: Manage alerts/comparisons, curate deals, approve affiliate links.
   - `DEAL_CURATOR`: Curate deals and guides; payout approval requires a second approver.
   - `ORG_VIEWER`: Read-only access to alerts, cost reports, and comparison lists.
2. Step-Up Authentication: Approving an affiliate payout, rotating a network API credential, or viewing security logs requires a fresh MFA/Passkey verification within the last 10 minutes (`step_up_verified_at`).

---

## 9. Security & Access Control

Argon2id hashing; refresh token reuse detection (automatic family revocation if a rotated token is replayed); rate limiting (3 failed logins/15m per IP+account); strict separation between TirpScanner user identity and external affiliate-network account ownership.

---

## 10. Failure Modes & Resilience

Graceful lockout recovery via email + backup recovery codes; atomic session revocation across Redis and PostgreSQL.

---

## 11. Inter-Module Dependencies

Foundation for User Dashboard (`U_UXUI.md`), Affiliate Operations (`W_WORKFLOW_ADMIN.md`), and Security (`S_SECURITY.md`). BYOC Connection (`M_MULTICLOUD_MANAGEMENT.md`) is RETIRED.

---

## 12. Testing & Verification Strategy

Unit and integration tests for Argon2id, token rotation replay detection, RBAC permission boundaries (`ORG_VIEWER` attempting `POST /api/v1/admin/affiliate/payouts` must return `403 Forbidden`), and step-up MFA enforcement.

---

## 13. Acceptance Criteria

- [x] All identity and organization capabilities are fully specified and backed by DDL and OpenAPI endpoints.
- [x] Step-up MFA is enforced before any payout approval or network-credential rotation.
- [x] Zero dependency on TreeScan user accounts or SSO.

---

## 14. Required Verification Evidence

Auth integration test logs, RBAC matrix test results, and OpenAPI auth security scheme validation.

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0101 ~ WO-CS-0150`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
