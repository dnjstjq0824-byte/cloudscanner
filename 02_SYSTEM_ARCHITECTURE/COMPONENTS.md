# COMPONENTS.md — CloudScanr Service Component & Package Topology

**Document ID:** `CS-ARCH-001`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent of TreeScan)

---

## 1. High-Level Runtime Binaries & Services

CloudScanr is architected as a **Modular Go Monolith with Isolated Worker Pools** paired with a **Nuxt 3 SSR Edge Frontend**, giving the development velocity and transactional consistency of a single codebase while isolating heavy background workloads (catalog ingestion, BYOC server provisioning) in dedicated worker pods.

| Component ID | Binary / Service | Language / Runtime | Primary Responsibilities | Scaling Strategy |
| :--- | :--- | :--- | :--- | :--- |
| `cs-edge-waf` | **Cloudflare Enterprise Edge** | Cloudflare Workers / WAF | DNS, Anycast CDN, TLS 1.3 termination, Bot Management, DDoS mitigation, `/go/*` edge rate limiting. | Global Anycast |
| `cs-web-ssr` | **CloudScanr Web Frontend** | Nuxt 3.13 + Node 22 LTS | Server-Side Rendering (SSR) for 16 Public SEO pages, Hydration for 18 User Dashboard pages, 23 Admin Console pages. | HPA 3–25 pods (CPU 65%) |
| `cs-api-gateway` | **CloudScanr Core API** (`cmd/api`) | Go 1.23 (`pgx/v5`, `chi`) | Serves all 22 REST API groups (`/api/v1/*`), `/go/{slug}` affiliate redirect engine, RBAC/RLS enforcement, KMS envelope encryption/decryption. | HPA 4–30 pods (RPS / p95 latency) |
| `cs-worker-cloud` | **BYOC Cloud Ops Worker** (`cmd/worker-cloud`) | Go 1.23 (`asynq`) | Executes user-approved BYOC server creation, power actions, resizing, backups, timeout reconciliation, and drift checks. | HPA 2–15 pods (Queue depth) |
| `cs-worker-pricing` | **Pricing & FX Sync Worker** (`cmd/worker-pricing`) | Go 1.23 (`asynq`) | Polls 18+ provider Pricing APIs (`PricingAdapter`), computes SHA-256 diffs, normalizes 36 TCO fields, updates OpenSearch aliases. | 2–6 pods (Static Egress NAT) |
| `cs-worker-affiliate` | **Affiliate & Outbox Worker** (`cmd/worker-affiliate`) | Go 1.23 (`asynq`) | Relays `outbox_events`, processes `/go/{slug}` click telemetry, ingests Impact/CJ/PartnerStack postbacks & CSVs, updates Commission Ledger. | HPA 2–10 pods |
| `cs-db-primary` | **PostgreSQL 16 HA Cluster** | PostgreSQL 16.4 | Primary ACID relational store (49 tables), Row-Level Security (RLS), Double-Entry Commission Ledger, Outbox table. | Primary + 2 Sync/Async Read Replicas |
| `cs-search-cluster` | **OpenSearch 2.17 Cluster** | OpenSearch 2.17 | Sub-50ms faceted SKU search, Natural Language spec query execution, Best-Value ranking index. | 3 Data/Master nodes |
| `cs-redis-cluster` | **Redis 7.4 Sentinel** | Redis 7.4 | Asynq job queues, rate-limit token buckets, `/go/{slug}` hot redirect cache, FX rate cache. | 3-node Sentinel HA |
| `cs-kms-vault` | **KMS / Secret Engine** | AWS KMS / HashiCorp Vault | Hardware-backed Key Encryption Keys (KEK) for `provider_credentials_metadata` Envelope Encryption (`AES-256-GCM`). | Managed Regional HA |

---

## 2. Go Monorepo Internal Package Layout (`internal/`)

```text
cloudscanr/
├── cmd/
│   ├── api/                  # Stateless HTTP API & /go/{slug} Redirect Server
│   ├── worker-cloud/         # BYOC Provisioning & Reconciliation Worker
│   ├── worker-pricing/       # Official Pricing API Harvester & Normalizer
│   ├── worker-affiliate/     # Click Attribution, Postback & Outbox Relay Worker
│   └── migrator/             # Database DDL Migration CLI
├── internal/
│   ├── adapters/             # STRICT 4-WAY ADAPTER SEPARATION (Article 6)
│   │   ├── affiliate/        # Impact, CJ, PartnerStack, HasOffers, Direct Referral
│   │   ├── pricing/          # DO, UpCloud, Cloudways v2, Vultr, Hetzner, Linode, GCP, AWS, Azure, OVH, Scaleway, OCI, etc.
│   │   ├── management/       # BYOC Server Create/Control/Delete per Provider Capability Matrix
│   │   └── reseller/         # Disabled stub in V1; gated by formal Partner/Reseller Contract
│   ├── domain/
│   │   ├── identity/         # Users, Orgs, RBAC, MFA, Passkeys, Sessions
│   │   ├── catalog/          # Providers, Regions, Products, SKUs, 16-Category Specs
│   │   ├── estimation/       # 36-Dimension Normalization, 730h/Calendar TCO, FX, VAT/GST
│   │   ├── search/           # OpenSearch Query Builder, NL Spec Parser, PostgreSQL Fallback
│   │   ├── ranking/          # Unbiased Value Score Engine (Zero Affiliate Weight)
│   │   ├── affiliate/        # SubID HMAC Signer, Fraud Filter, Commission Ledger
│   │   ├── cloudmgmt/        # BYOC Connections, Cost Approval Gate, Provisioning FSM
│   │   ├── content/          # Benchmarks, Reviews, Tutorials, SEO Metadata
│   │   └── admin/            # Quarantine Resolution, Four-Eyes Approval, Feature Flags
│   ├── platform/
│   │   ├── kms/              # Envelope Encryption (GenerateDataKey, SealWithAAD, OpenWithAAD)
│   │   ├── database/         # pgxpool with automatic RLS tenant context injection
│   │   ├── observability/    # OpenTelemetry traces, Prometheus metrics, Secret-Redacting slog
│   │   └── httputil/         # RFC 9457 Problem Details, Idempotency-Key middleware, Rate Limiter
├── web/                      # Nuxt 3.13 SSR Frontend Application
└── migrations/               # PostgreSQL 16 SQL DDL Migrations (000001 ~ 000012)
```
