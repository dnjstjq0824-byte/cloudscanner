# P — PostgreSQL — Complete Relational Schema, 6+ Domain ERD, DDL Migrations, RLS & Ledger Integrity

**Document ID:** `CS-AZ-P`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0060 ~ WO-CS-0074, WO-CS-0801 ~ WO-CS-0815`

---

## 1. 기능 목적 (Functional Purpose)
Define the authoritative PostgreSQL 16 database architecture covering all 45+ required tables across Identity, Providers, Pricing, Affiliate, Cloud Management, Search & Content, and Operations (Article 18), with strict foreign key integrity, `NUMERIC(18,8)` financial precision, Row-Level Security (RLS), and zero plaintext secret storage.

---

## 2. 사용자 역할 (User Roles & Actors)
Database Agent, Backend Agent, Security Auditor.

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. Running `golang-migrate` applies `000001_extensions_and_types.up.sql` through `000012_seed_reference_data.up.sql` cleanly on a fresh PostgreSQL 16 instance.
2. An application transaction sets `SET LOCAL app.current_org_id = '...';` and queries `cloud_resources`; PostgreSQL RLS guarantees zero rows from other organizations can ever be returned.
3. An attempt to `UPDATE` or `DELETE` a finalized row in `commissions` or `audit_logs` is blocked at the database engine level by an immutability trigger.

---

## 4. 기술 스택 (Technology Stack)
PostgreSQL 16, `pgcrypto`, `pg_trgm`, `btree_gist`, `golang-migrate/migrate/v4`, `pgx/v5` connection pool with PgBouncer transaction pooling.

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`users`, `organizations`, `organization_members`, `roles`, `sessions`, `auth_identities`, `cloud_providers`, `provider_regions`, `provider_products`, `provider_capabilities`, `provider_data_sources`, `product_skus`, `product_specs`, `price_snapshots`, `price_components`, `billing_terms`, `promotions`, `fx_rates`, `tax_policies`, `price_sync_jobs`, `affiliate_programs`, `affiliate_accounts`, `affiliate_campaigns`, `affiliate_links`, `affiliate_clicks`, `affiliate_conversions`, `commissions`, `commission_adjustments`, `affiliate_payouts`, `cloud_connections`, `provider_credentials_metadata`, `cloud_resources`, `resource_operations`, `operation_approvals`, `cloud_usage_snapshots`, `budget_policies`, `resource_reconciliations`, `search_events`, `comparison_lists`, `favorites`, `price_alerts`, `articles`, `provider_reviews`, `audit_logs`, `background_jobs`, `outbox_events`, `admin_actions`, `feature_flags`, `system_settings`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)
**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**
- `Internal DB layer backing all 22 API groups`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. UI/UX 설계 (User Interface & Experience)
Reflected in Admin Database & Migration Health view (`/admin/system`).

---

## 8. 상세 처리 로직 (Core Processing Logic)
1. Strict Financial & Timestamp Types: All monetary columns use `NUMERIC(18,8)` (never `FLOAT` or `REAL`); all timestamps use `TIMESTAMPTZ` stored in UTC.
2. Secret Prohibition Rule: `provider_credentials_metadata` stores ONLY `kms_key_id`, `encrypted_dek`, `ciphertext_blob`, `nonce`, `token_fingerprint_sha256`, and `scopes`. Plaintext secret columns do not exist.
3. Partitioning: High-volume append tables (`price_snapshots`, `affiliate_clicks`, `search_events`, `audit_logs`) support monthly range partitioning by `created_at`.

---

## 9. 보안 통제 (Security & Access Control)
Separate PostgreSQL roles: `cloudscanr_migrator` (DDL owner), `cloudscanr_app` (DML runtime user subject to RLS), `cloudscanr_readonly` (analytics replica reader).

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
Strict statement timeouts (`statement_timeout = 5000ms` for web requests, `60000ms` for batch workers) and deadlock retry wrappers.

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Core persistence layer for all modules (`03_DATABASE/`).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
Automated up/down migration idempotency tests, RLS cross-tenant isolation tests, and immutability trigger tests.

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] All 49 required domain tables in Article 18 are created with primary keys, foreign keys, check constraints, and indexes in `03_DATABASE/migrations/`.
- [x] Zero plaintext API secret columns exist anywhere in the schema.
- [x] All SQL migrations pass syntax and referential integrity checks.

---

## 14. 검증 증빙 (Required Verification Evidence)
Complete SQL DDL files in `03_DATABASE/migrations/`, `ERD.md`, `TABLE_DICTIONARY.md`, and `INDEX_DESIGN.md`.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0060 ~ WO-CS-0074, WO-CS-0801 ~ WO-CS-0815`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
