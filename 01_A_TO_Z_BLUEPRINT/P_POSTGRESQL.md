# P — PostgreSQL — Complete Relational Schema, Domain ERD, DDL Migrations, RLS & Ledger Integrity

**Document ID:** `CS-AZ-P` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (TirpScanner pivot reframe, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Pivot Note (2026-11-26):** BYOC/cloud-management tables retired (names historical). Catalog, plan/deal, affiliate, identity, content, and ops tables are live. Stack unchanged (PostgreSQL 16). See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0060 ~ WO-CS-0074, WO-CS-0801 ~ WO-CS-0815`

---

## 1. 기능 목적 (Functional Purpose)

Define the authoritative PostgreSQL 16 database architecture covering required tables across Identity, Vendors, Plans/Deals, Affiliate, Search & Content, and Operations, with strict foreign key integrity, `NUMERIC(18,8)` financial precision, Row-Level Security (RLS), and zero plaintext secret storage. Cloud Management tables are retired (see §3).

---

## 2. User Roles & Actors

Database Agent, Backend Agent, Security Auditor.

---

## 3. User Scenarios

1. Running `golang-migrate` applies `000001_extensions_and_types.up.sql` through `000011_seed_providers_and_roles.up.sql` cleanly on a fresh PostgreSQL 16 instance (`000012` rollback reference lives in `03_DATABASE/down_migrations_reference/`, outside initdb reach).
2. An application transaction sets `SET LOCAL app.current_org_id = '...';` and queries `price_alerts`; PostgreSQL RLS guarantees zero rows from other organizations can ever be returned.
3. An attempt to `UPDATE` or `DELETE` a finalized row in `commissions` or `audit_logs` is blocked at the database engine level by an immutability trigger.

---

## 4. Technology Stack

PostgreSQL 16, `pgcrypto`, `pg_trgm`, `btree_gist`, `golang-migrate/migrate/v4`, `pgx/v5` connection pool with PgBouncer transaction pooling.

---

## 5. Database Schema & Tables

**Live PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`users`, `organizations`, `organization_members`, `roles`, `sessions`, `auth_identities`, `cloud_providers`, `provider_regions`, `provider_products`, `provider_capabilities`, `provider_data_sources`, `product_skus`, `product_specs`, `price_snapshots`, `price_components`, `billing_terms`, `promotions`, `fx_rates`, `tax_policies`, `price_sync_jobs`, `affiliate_programs`, `affiliate_accounts`, `affiliate_campaigns`, `affiliate_links`, `affiliate_clicks`, `affiliate_conversions`, `commissions`, `commission_adjustments`, `affiliate_payouts`, `search_events`, `comparison_lists`, `favorites`, `price_alerts`, `articles`, `provider_reviews`, `audit_logs`, `background_jobs`, `outbox_events`, `admin_actions`, `feature_flags`, `system_settings`

**Retired Tables (BYOC era; names historical; no new live flows):**  
`cloud_connections`, `provider_credentials_metadata`, `cloud_resources`, `resource_operations`, `operation_approvals`, `cloud_usage_snapshots`, `budget_policies`, `resource_reconciliations`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)

**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**

- `Internal DB layer backing all live API groups (cloud groups retired)`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. User Interface & Experience

Reflected in Admin Database & Migration Health view (`/admin/system`).

---

## 8. Core Processing Logic

1. Strict Financial & Timestamp Types: All monetary columns use `NUMERIC(18,8)` (never `FLOAT` or `REAL`); all timestamps use `TIMESTAMPTZ` stored in UTC.
2. Secret Prohibition Rule: Credential metadata stores ONLY `kms_key_id`, `encrypted_dek`, `ciphertext_blob`, `nonce`, `token_fingerprint_sha256`, and `scopes`. Plaintext secret columns do not exist.
3. Partitioning: High-volume append tables (`price_snapshots`, `affiliate_clicks`, `search_events`, `audit_logs`) support monthly range partitioning by `created_at`.

---

## 9. Security & Access Control

Separate PostgreSQL roles: `tirpscanner_migrator` (DDL owner), `tirpscanner_app` (DML runtime user subject to RLS), `tirpscanner_readonly` (analytics replica reader).

---

## 10. Failure Modes & Resilience

Strict statement timeouts (`statement_timeout = 5000ms` for web requests, `60000ms` for batch workers) and deadlock retry wrappers.

---

## 11. Inter-Module Dependencies

Core persistence layer for all modules (`03_DATABASE/`).

---

## 12. Testing & Verification Strategy

Automated up/down migration idempotency tests, RLS cross-tenant isolation tests, and immutability trigger tests.

---

## 13. Acceptance Criteria

- [x] All live domain tables are created with primary keys, foreign keys, check constraints, and indexes in `03_DATABASE/migrations/`; retired tables are not extended.
- [x] Zero plaintext API secret columns exist anywhere in the schema.
- [x] All SQL migrations pass syntax and referential integrity checks.

---

## 14. Required Verification Evidence

Complete SQL DDL files in `03_DATABASE/migrations/`, `ERD.md`, `TABLE_DICTIONARY.md`, and `INDEX_DESIGN.md`.

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0060 ~ WO-CS-0074, WO-CS-0801 ~ WO-CS-0815`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
