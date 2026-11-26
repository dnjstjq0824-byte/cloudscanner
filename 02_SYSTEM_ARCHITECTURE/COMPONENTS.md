# COMPONENTS.md — TirpScanner Service Component & Package Topology

**Document ID:** `CS-ARCH-001` (retained for traceability)
**Platform:** TirpScanner (`tirpscanr.com` — global SaaS/AI-tool & software affiliate deal scanner)

> **Pivot addendum (2026-11-26):** Refreshed for the pivot from CloudScanr (cloud price comparison + BYOC multi-cloud management) to TirpScanner. Live scope is the SaaS/AI-tool catalog, deal scan / price-drop detection, and affiliate monetization. Cloud-infrastructure catalog, BYOC provisioning, and server management are REMOVED scope. Domain `tirpscanr.com` is confirmed (apex registered on Cloudflare, DNS wiring pending); `api.`/`www.` subdomains are ASSUMED.

---

## 1. High-Level Runtime Binaries & Services

TirpScanner is architected as a **Modular Go Monolith with Isolated Worker Pools** paired with a **Nuxt 3 SSR Edge Frontend**, giving the development velocity and transactional consistency of a single codebase while isolating heavy background workloads (deal ingestion, price-drop evaluation) in dedicated worker pods.

| Component ID | Binary / Service | Language / Runtime | Primary Responsibilities | Scaling Strategy |
| :--- | :--- | :--- | :--- | :--- |
| `cs-edge-waf` | **Cloudflare Enterprise Edge** | Cloudflare Workers / WAF | DNS, Anycast CDN, TLS 1.3 termination, Bot Management, DDoS mitigation, `/go/*` edge rate limiting. | Global Anycast |
| `cs-web-ssr` | **TirpScanner Web Frontend** | Nuxt 3.13 + Node 24 LTS (pnpm 9.15.0) | Server-Side Rendering (SSR) for Public SEO deal pages, Hydration for User Dashboard pages, Admin Console pages. | HPA 3–25 pods (CPU 65%) |
| `cs-api-gateway` | **TirpScanner Core API** (`cmd/api`) | Go 1.27 (`pgx/v5`, `chi`) | Serves the live REST API surface (affiliate/catalog/deal/search/admin), `/go/{slug}` affiliate redirect engine, RBAC/RLS enforcement. | HPA 4–30 pods (RPS / p95 latency) |
| `cs-worker-cloud` | **BYOC Cloud Ops Worker** (`cmd/worker-cloud`) | — (retired) | **RETIRED (2026-11-26).** BYOC server provisioning removed with the pivot; do not deploy or extend. | — |
| `cs-worker-deals` | **Deal & Plan-Price Sync Worker** (`cmd/worker-deals`, succeeds `worker-pricing`) | Go 1.27 (`asynq`) | Polls SaaS vendor pricing/deal feeds and affiliate product feeds (`PricingAdapter`; name historical), computes SHA-256 diffs, normalizes subscription plan prices, detects price drops, updates OpenSearch aliases. | 2–6 pods (Static Egress NAT) |
| `cs-worker-affiliate` | **Affiliate & Outbox Worker** (`cmd/worker-affiliate`) | Go 1.27 (`asynq`) | Relays `outbox_events`, processes `/go/{slug}` click telemetry, ingests Impact/CJ/PartnerStack postbacks & CSVs, updates Commission Ledger. | HPA 2–10 pods |
| `cs-db-primary` | **PostgreSQL 16 HA Cluster** | PostgreSQL 16.4 | Primary ACID relational store (SaaS catalog, deals, price alerts, affiliate ledger), Row-Level Security (RLS), Double-Entry Commission Ledger, Outbox table. | Primary + 2 Sync/Async Read Replicas |
| `cs-search-cluster` | **OpenSearch 2.17 Cluster** | OpenSearch 2.17 | Sub-50ms faceted SaaS/AI-tool & deal search, Best-Value ranking index. | 3 Data/Master nodes |
| `cs-redis-cluster` | **Redis 7.4 Sentinel** | Redis 7.4 | Asynq job queues, rate-limit token buckets, `/go/{slug}` hot redirect cache, FX rate cache. | 3-node Sentinel HA |
| `cs-kms-vault` | **KMS / Secret Engine** | AWS KMS / HashiCorp Vault | Envelope encryption for affiliate/integration credentials (BYOC-token use retired 2026-11-26). | Managed KEK rotation |

---

## 2. Go Monorepo Internal Package Layout (`internal/`)

```text
tirpscanner/
├── cmd/
│   ├── api/                  # Stateless HTTP API & /go/{slug} Redirect Server (LIVE)
│   ├── worker-deals/         # SaaS Deal & Plan-Price Harvester, Normalizer & Price-Drop Detector (LIVE; succeeds worker-pricing)
│   ├── worker-affiliate/     # Click Attribution, Postback & Outbox Relay Worker (LIVE)
│   ├── worker-cloud/         # RETIRED — BYOC Provisioning worker removed with the pivot
│   └── migrator/             # Database DDL Migration CLI
├── internal/
│   ├── adapters/             # TWO LIVE ADAPTERS (+ two retired, see 04_API_SPECIFICATION/PROVIDER_ADAPTER_CONTRACT.md)
│   │   ├── affiliate/        # LIVE — Impact, CJ, PartnerStack, HasOffers, Direct Referral
│   │   ├── pricing/          # LIVE — SaaS/AI-tool catalog & deal-ingestion adapters (name historical)
│   │   ├── management/       # RETIRED — BYOC adapters removed; do not extend
│   │   └── reseller/         # RETIRED — never activated; removed with the pivot
│   ├── domain/
│   │   ├── identity/         # Users, Orgs, RBAC, MFA, Passkeys, Sessions (LIVE)
│   │   ├── catalog/          # SaaS vendors, AI tools, subscription plans, deals (LIVE)
│   │   ├── estimation/       # Plan-price normalization, price-drop detection, FX, VAT/GST (LIVE)
│   │   ├── search/           # OpenSearch Query Builder, deal search, PostgreSQL Fallback (LIVE)
│   │   ├── ranking/          # Unbiased Value Score Engine (Zero Affiliate Weight) (LIVE)
│   │   ├── affiliate/        # SubID HMAC Signer, Fraud Filter, Commission Ledger (LIVE)
│   │   ├── cloudmgmt/        # RETIRED — BYOC domain removed with the pivot
│   │   ├── content/          # Reviews, tutorials, deal guides, SEO Metadata (LIVE)
│   │   └── admin/            # Quarantine Resolution, Approvals, Feature Flags (LIVE)
│   ├── platform/
│   │   ├── kms/              # LIVE — envelope encryption for affiliate/integration credentials (BYOC use retired)
│   │   ├── database/         # pgxpool with automatic RLS tenant context injection
│   │   ├── observability/    # OpenTelemetry traces, Prometheus metrics, Secret-Redacting slog
│   │   └── httputil/         # RFC 9457 Problem Details, Idempotency-Key middleware, Rate Limiter
├── web/                      # Nuxt 3.13 SSR Frontend Application
└── migrations/               # PostgreSQL 16 SQL DDL Migrations (historical; BYOC tables retained as legacy, no live reads/writes)
```
