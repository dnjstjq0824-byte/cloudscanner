# W — Workflow & Administration — Operations Console, Deal Quarantine, Affiliate Reconciliation & Security Ops

**Document ID:** `CS-AZ-W` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (TirpScanner pivot reframe, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Pivot Note (2026-11-26):** BYOC audit/freeze workflows retired. Live admin: vendors, deal quarantine, affiliate ops, content moderation, flags. See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0701 ~ WO-CS-0750`

---

## 1. 기능 목적 (Functional Purpose)

Provide platform operators with a comprehensive, auditable Admin Console (`/admin/*`) to govern vendors, resolve deal-sync anomalies, manage affiliate contracts/links/commissions, moderate content/reviews, and control system feature flags.

---

## 2. User Roles & Actors

SuperAdmin, Catalog Admin, Affiliate Ops Manager, Security Auditor, Content Editor.

---

## 3. User Scenarios

1. An automated deal sync flags a SaaS plan price discrepancy; the Catalog Admin reviews the side-by-side diff in `/admin/pricing-failures`, verifies the official vendor page, approves the update, and records a signed `admin_actions` audit entry.
2. Security Admin detects abnormal affiliate-link modification attempts from an organization and triggers an emergency freeze on that organization's publishing rights via `/admin/security-events`.

---

## 4. Technology Stack

Go Admin API (`internal/api/admin`), Nuxt 3 Admin Console (`pages/admin/**`), PostgreSQL `admin_actions`, `feature_flags`, `system_settings`.

---

## 5. Database Schema & Tables

**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`admin_actions`, `audit_logs`, `feature_flags`, `system_settings`, `price_sync_jobs`, `affiliate_conversions`, `commissions`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)

**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**

- `GET /api/v1/admin/overview`
- `POST /api/v1/admin/providers`
- `POST /api/v1/admin/pricing/quarantine/{id}/resolve`
- `POST /api/v1/admin/affiliate/payouts/{id}/approve`
- `POST /api/v1/admin/security/orgs/{id}/freeze`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. User Interface & Experience

Dedicated Admin Console screens (`05_UXUI/ADMIN_CONSOLE.md`) with dense operational tables, diff viewers, CSV export/import wizards, and mandatory reason-code prompts for every mutating admin action.

---

## 8. Core Processing Logic

1. Four-Eyes Principle for High-Risk Admin Actions: Modifying affiliate payout bank metadata, overriding an alert-rate ceiling, or changing a vendor's base feed domain requires approval from a second `SUPER_ADMIN` or step-up WebAuthn signature.
2. Every admin mutation writes to `admin_actions` with `admin_user_id`, `action_type`, `target_entity`, `before_state_jsonb`, `after_state_jsonb`, `reason`, and `ip_address`.

---

## 9. Security & Access Control

Admin Console routes require `ADMIN_*` role + WebAuthn hardware passkey MFA + optional Cloudflare Access / IP allowlist enforcement.

---

## 10. Failure Modes & Resilience

Admin UI never marks a backend feature as complete if the underlying worker or feed adapter is unimplemented.

---

## 11. Inter-Module Dependencies

Depends on Identity (`I`), Catalog (`C`), and Affiliate (`L`). Cloud Management (`M`) is RETIRED.

---

## 12. Testing & Verification Strategy

Admin RBAC authorization matrix tests, audit log before/after diff capture tests, and quarantine resolution workflow tests.

---

## 13. Acceptance Criteria

- [x] All live Admin Console screens are fully specified and wired to real backend APIs and database tables.
- [x] Every mutating admin endpoint requires a non-empty `reason` field and logs `before_state`/`after_state` in `admin_actions`.

---

## 14. Required Verification Evidence

Admin API integration test logs and `10_ADMIN_OPERATIONS/` documentation.

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0701 ~ WO-CS-0750`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
