# W — Workflow & Administration — Operations Console, Catalog Quarantine, Affiliate Reconciliation & Security Ops

**Document ID:** `CS-AZ-W`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0701 ~ WO-CS-0750`

---

## 1. 기능 목적 (Functional Purpose)
Provide platform operators with a comprehensive, auditable Admin Console (`/admin/*`) to govern providers, resolve pricing sync anomalies, manage affiliate contracts/links/commissions, audit BYOC server operations, moderate content/reviews, and control system feature flags.

---

## 2. 사용자 역할 (User Roles & Actors)
SuperAdmin, Catalog Admin, Affiliate Ops Manager, Security Auditor, Content Editor.

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. An automated price sync flags a Contabo SKU price discrepancy; the Catalog Admin reviews the side-by-side diff in `/admin/pricing-failures`, verifies the official rate card, approves the update, and records a signed `admin_actions` audit entry.
2. Security Admin detects abnormal API token verification failures from an organization and triggers an emergency freeze on that organization's BYOC connections via `/admin/security-events`.

---

## 4. 기술 스택 (Technology Stack)
Go Admin API (`internal/api/admin`), Nuxt 3 Admin Console (`pages/admin/**`), PostgreSQL `admin_actions`, `feature_flags`, `system_settings`.

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`admin_actions`, `audit_logs`, `feature_flags`, `system_settings`, `price_sync_jobs`, `affiliate_conversions`, `commissions`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)
**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**
- `GET /api/v1/admin/dashboard`
- `POST /api/v1/admin/providers`
- `POST /api/v1/admin/pricing/quarantine/{id}/resolve`
- `POST /api/v1/admin/affiliate/payouts/{id}/approve`
- `POST /api/v1/admin/security/connections/{id}/freeze`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. UI/UX 설계 (User Interface & Experience)
23 dedicated Admin Console screens (`05_UXUI/ADMIN_CONSOLE.md`) with dense operational tables, diff viewers, CSV export/import wizards, and mandatory reason-code prompts for every mutating admin action.

---

## 8. 상세 처리 로직 (Core Processing Logic)
1. Four-Eyes Principle for High-Risk Admin Actions: Modifying affiliate payout bank metadata, overriding a global budget ceiling, or changing a provider's base API domain requires approval from a second `SUPER_ADMIN` or step-up WebAuthn signature.
2. Every admin mutation writes to `admin_actions` with `admin_user_id`, `action_type`, `target_entity`, `before_state_jsonb`, `after_state_jsonb`, `reason`, and `ip_address`.

---

## 9. 보안 통제 (Security & Access Control)
Admin Console routes require `ADMIN_*` role + WebAuthn hardware passkey MFA + optional Cloudflare Access / IP allowlist enforcement.

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
Admin UI never marks a backend feature as complete if the underlying worker or provider adapter is unimplemented (strict compliance with Article 29).

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Depends on Identity (`I`), Catalog (`C`), Affiliate (`L`), and Cloud Management (`M`).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
Admin RBAC authorization matrix tests, audit log before/after diff capture tests, and quarantine resolution workflow tests.

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] All 23 Admin Console screens listed in Article 16 are fully specified and wired to real backend APIs and database tables.
- [x] Every mutating admin endpoint requires a non-empty `reason` field and logs `before_state`/`after_state` in `admin_actions`.

---

## 14. 검증 증빙 (Required Verification Evidence)
Admin API integration test logs and `10_ADMIN_OPERATIONS/` documentation.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0701 ~ WO-CS-0750`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
