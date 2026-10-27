# A — Architecture — Independent Platform Topology, Service Boundaries & Fault Isolation

**Document ID:** `CS-AZ-A`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0051 ~ WO-CS-0100`

---

## 1. 기능 목적 (Functional Purpose)
Define the complete, standalone system architecture of CloudScanr (`cloudscanr.com`) with zero operational, database, or identity coupling to TreeScan. Establish strict service boundaries across the 8 core CloudScanr services and 4 provider adapter layers (Affiliate, Pricing, Management, Reseller), ensuring high availability, horizontal scalability, and blast-radius containment.

---

## 2. 사용자 역할 (User Roles & Actors)
Anonymous Visitor, Authenticated Member, Organization Owner/Admin/Viewer, Affiliate Operations Manager, Cloud Infrastructure Operator, Platform SuperAdmin, External Search Crawler (Googlebot, Bingbot, TreeScanbot).

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. Anonymous user searches 18+ cloud providers for a 4 vCPU / 16 GB NVMe server in Tokyo/Frankfurt with sub-50ms facet response via Nuxt 3 SSR and OpenSearch.
2. User clicks an affiliate deal link (`/go/do-droplet-s4vcpu`), triggering non-blocking click telemetry and instant 302 redirect while an outbox worker persists attribution state.
3. Authenticated Organization Admin connects a DigitalOcean or UpCloud account via BYOC and provisions a server through an isolated asynchronous job queue without impacting public search latency.
4. External Provider API (e.g., Vultr or Hetzner) experiences a 503 outage; the Pricing and Management circuit breakers trip independently while serving cached/stale-flagged catalog data.

---

## 4. 기술 스택 (Technology Stack)
Go 1.23+ (Chi/Echo HTTP router, pgx/v5 PostgreSQL driver, Asynq distributed task queue), Nuxt 3.13+ (Vue 3.5, TypeScript 5.6, Pinia, TailwindCSS), PostgreSQL 16 (Primary + Read Replicas, Row-Level Security), OpenSearch 2.17 (BM25 + Dynamic Facet Aggregations), Redis 7.4 (Rate limiting, cache, Asynq broker), AWS KMS / HashiCorp Vault (Envelope Encryption), Cloudflare Enterprise WAF/CDN.

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`users`, `organizations`, `cloud_providers`, `product_skus`, `price_snapshots`, `affiliate_links`, `cloud_connections`, `cloud_resources`, `outbox_events`, `audit_logs`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)
**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**
- `GET /api/v1/health`
- `GET /api/v1/readiness`
- `GET /api/v1/search/skus`
- `GET /go/{slug}`
- `POST /api/v1/cloud/connections`
- `POST /api/v1/cloud/servers`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. UI/UX 설계 (User Interface & Experience)
Three physically and logically separated frontend layouts: (1) Public Web (`layouts/public.vue`, SSR-first, zero auth bundle overhead), (2) User Cloud & Personalization Console (`layouts/dashboard.vue`, CSR/SSR hybrid with MFA/approval modals), and (3) Admin & Affiliate Operations Console (`layouts/admin.vue`, strict RBAC + IP/WebAuthn enforcement).

---

## 8. 상세 처리 로직 (Core Processing Logic)
1. Edge Layer: Cloudflare WAF terminates TLS 1.3, blocks bot floods, caches static assets, and routes `/api/*` and `/go/*` to the Go API Gateway.
2. Modular Go Backend: Structured into `cmd/api` (stateless HTTP servers), `cmd/worker` (Asynq background processors), and `cmd/scheduler` (cron price/FX/reconciliation dispatchers).
3. Transactional Outbox Pattern: Every mutating state change (server creation request, affiliate conversion, price change) writes atomically to its domain table and `outbox_events` within a single PostgreSQL transaction, relayed to Redis/OpenSearch within <200ms.
4. Four-Adapter Isolation: `AffiliateAdapter`, `PricingAdapter`, `ManagementAdapter`, and `ResellerAdapter` live in separate Go packages with distinct credentials and circuit breakers (`sony/gobreaker`).

---

## 9. 보안 통제 (Security & Access Control)
Complete network and credential isolation from TreeScan. Mutual TLS (mTLS) between internal pods. Strict CORS policy (`https://cloudscanr.com`). Zero raw secrets in environment variables or database tables (all BYOC tokens encrypted via KMS Envelope Encryption). Per-tenant PostgreSQL Row-Level Security (`app.current_org_id`).

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
Circuit breakers trip after 5 consecutive provider 5xx errors (half-open after 60s). If OpenSearch is unreachable, `SearchService` automatically falls back to PostgreSQL indexed queries (`pg_trgm` + composite B-tree on `product_skus`). If Redis is degraded, `/go/{slug}` reads from an in-memory LRU cache synced every 30s from PostgreSQL so affiliate monetization never drops a click.

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Foundation for all modules B through Z. Must be initialized in Phase 1 before Catalog (C), Affiliate (L), or BYOC Management (M).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
Architecture fitness tests (Go package import graph verification ensuring `pricing` never imports `management` credentials), chaos testing (killing OpenSearch/Redis containers during k6 load test at 1,500 RPS), and health/readiness probe verification.

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] Zero network, DB, or code dependencies on TreeScan verified via `go.mod`, `package.json`, and Terraform state inspection.
- [x] Public search API maintains p95 < 80ms at 1,000 concurrent requests even when OpenSearch is simulated offline (PostgreSQL fallback active).
- [x] Four provider adapter interfaces (`AffiliateAdapter`, `PricingAdapter`, `ManagementAdapter`, `ResellerAdapter`) compile as isolated Go packages.

---

## 14. 검증 증빙 (Required Verification Evidence)
Go module dependency graph output (`go list -m all`), Docker Compose / Kubernetes deployment logs, k6 load test summary JSON, and fallback integration test logs.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0051 ~ WO-CS-0100`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
