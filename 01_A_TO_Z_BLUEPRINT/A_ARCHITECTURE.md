# A — Architecture — Independent Platform Topology, Service Boundaries & Fault Isolation

**Document ID:** `CS-AZ-A` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (TirpScanner pivot reframe, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Pivot Note (2026-11-26):** Revised adapter law: Affiliate + Catalog/Deal-ingestion adapters are LIVE; Management (BYOC) + Reseller adapters are RETIRED. See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0051 ~ WO-CS-0100`

---

## 1. 기능 목적 (Functional Purpose)

Define the complete, standalone system architecture of TirpScanner (`tirpscanr.com`) with zero operational, database, or identity coupling to TreeScan. Establish strict service boundaries across the core TirpScanner services and the revised adapter law (LIVE: Affiliate, Catalog/Deal-ingestion; RETIRED: Management/BYOC, Reseller), ensuring high availability, horizontal scalability, and blast-radius containment.

---

## 2. User Roles & Actors

Anonymous Visitor, Authenticated Member, Organization Owner/Admin/Viewer, Affiliate Operations Manager, Deal Curation Operator, Platform SuperAdmin, External Search Crawler (Googlebot, Bingbot, TreeScanbot).

---

## 3. User Scenarios

1. Anonymous user searches the SaaS/AI-tool catalog for an AI video editor under $20/mo with sub-50ms facet response via Nuxt 3 SSR and OpenSearch.
2. User clicks an affiliate deal link (`/go/notion-plus-annual`), triggering non-blocking click telemetry and instant 302 redirect while an outbox worker persists attribution state.
3. Authenticated member subscribes to a price-drop alert on a design tool; when the deal scanner detects a promo, an isolated async job queue delivers the alert without impacting public search latency.
4. External deal-feed source (e.g., a coupon feed or network API) experiences a 503 outage; the Catalog/Deal-ingestion circuit breaker trips independently while serving cached/stale-flagged catalog data.

---

## 4. Technology Stack

Go 1.23+ (Chi/Echo HTTP router, pgx/v5 PostgreSQL driver, Asynq distributed task queue), Nuxt 3.13+ (Vue 3.5, TypeScript 5.6, Pinia, TailwindCSS), PostgreSQL 16 (Primary + Read Replicas, Row-Level Security), OpenSearch 2.17 (BM25 + Dynamic Facet Aggregations), Redis 7.4 (Rate limiting, cache, Asynq broker), AWS KMS / HashiCorp Vault (Envelope Encryption), Cloudflare Enterprise WAF/CDN.

---

## 3. Database Schema & Tables

**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`users`, `organizations`, `cloud_providers`, `product_skus`, `price_snapshots`, `affiliate_links`, `outbox_events`, `audit_logs`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.
- BYOC tables (`cloud_connections`, `cloud_resources`, and related) are RETIRED; table names are historical (see `P_POSTGRESQL.md`).

---

## 6. API 구조 (API Contracts & Endpoints)

**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**

- `GET /api/v1/health`
- `GET /api/v1/readiness`
- `GET /api/v1/search/skus`
- `GET /go/{slug}`
- `POST /api/v1/alerts` (price-drop alert subscription)
- `GET /api/v1/deals` (deal scan results)

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.
- `/api/v1/cloud/*` (BYOC connection/provisioning) endpoints are RETIRED.

---

## 7. User Interface & Experience

Three physically and logically separated frontend layouts: (1) Public Web (`layouts/public.vue`, SSR-first, zero auth bundle overhead), (2) User Deals & Personalization Console (`layouts/dashboard.vue`, CSR/SSR hybrid with alert management), and (3) Admin & Affiliate Operations Console (`layouts/admin.vue`, strict RBAC + IP/WebAuthn enforcement).

---

## 8. Core Processing Logic

1. Edge Layer: Cloudflare WAF terminates TLS 1.3, blocks bot floods, caches static assets, and routes `/api/*` and `/go/*` to the Go API Gateway.
2. Modular Go Backend: Structured into `cmd/api` (stateless HTTP servers), `cmd/worker` (Asynq background processors), and `cmd/scheduler` (cron deal-scan/FX/reconciliation dispatchers).
3. Transactional Outbox Pattern: Every mutating state change (alert subscription, affiliate conversion, price change) writes atomically to its domain table and `outbox_events` within a single PostgreSQL transaction, relayed to Redis/OpenSearch within <20ms.
4. Revised Adapter Isolation: `AffiliateAdapter` and `PricingAdapter` (catalog/deal-ingestion; `pricing` package name historical) are LIVE in separate Go packages with distinct credentials and circuit breakers (`sony/gobreaker`). `ManagementAdapter` (BYOC) and `ResellerAdapter` are RETIRED — no new work, no new credentials.

---

## 9. Security & Access Control

Complete network and credential isolation from TreeScan. Mutual TLS (mTLS) between internal2. Strict CORS policy (`https://tirpscanr.com` — registered, wiring pending). Zero raw secrets in environment variables or database tables (all affiliate-network API keys encrypted via KMS Envelope Encryption). Per-tenant PostgreSQL Row-Level Security (`app.current_org_id`).

---

## 10. Failure Modes & Resilience

Circuit breakers trip after 3 consecutive feed-source 3xx errors (half-open after 60s). If OpenSearch is unreachable, `SearchService` automatically falls back to PostgreSQL indexed queries (`pg_trgm` + composite B-tree on `product_skus`). If Redis is degraded, `/go/{slug}` reads from an in-memory LRU cache synced every 30s from PostgreSQL so affiliate monetization never drops a click.

---

## 11. Inter-Module Dependencies

Foundation for all modules B through Z. Must be initialized in Phase 1 before Catalog (C), Affiliate (L), or Deal Ingestion (D). Multi-Cloud Management (M) is RETIRED and is not a dependency.

---

## 12. Testing & Verification Strategy

Architecture fitness tests (Go package import graph verification ensuring `pricing` never imports retired `management` credentials), chaos testing (killing OpenSearch/Redis containers during k6 load test at 1,500 RPS), and health/readiness probe verification.

---

## 13. Acceptance Criteria

- [x] Zero network, DB, or code dependencies on TreeScan verified via `go.mod`, `package.json`, and Terraform state inspection.
- [x] Public search API maintains p95 < 80ms at 1,000 concurrent requests even when OpenSearch is simulated offline (PostgreSQL fallback active).
- [x] Live provider adapter interfaces (`AffiliateAdapter`, `PricingAdapter` for catalog/deal-ingestion) compile as isolated Go packages; retired `ManagementAdapter`/`ResellerAdapter` receive no new wiring.

---

## 14. Required Verification Evidence

Go module dependency graph output (`go list -m all`), Docker Compose / Kubernetes deployment logs, k6 load test summary JSON, and fallback integration test logs.

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0051 ~ WO-CS-0100`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
