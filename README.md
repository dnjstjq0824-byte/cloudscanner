# CloudScanr — Global Cloud Price Comparison, Affiliate Monetization & BYOC Multi-Cloud Management Platform

**Project Code:** `BN-CLOUDSCANR-GEMINI4-ARGON-MASTER-BLUEPRINT-1000-001`  
**Version:** `1.0.0`  
**Architecture Paradigm:** Modular Go Monolith + Async Worker Fleet + Nuxt 3 SSR + PostgreSQL 16 + OpenSearch 2.17  
**Platform Independence:** 100% Independent of TreeScan (Dedicated Domain, Auth, DB, Search, Backend, Admin, CI/CD)

---

## 1. Platform Overview & Core Services (Article 2)

CloudScanr unifies global cloud infrastructure discovery, transparent total-cost-of-ownership (TCO) comparison, compliant affiliate monetization, and Bring-Your-Own-Cloud (BYOC) server lifecycle management in a single platform:

| Service ID | Service Name | Core Capabilities |
| :--- | :--- | :--- |
| **Service 1** | **Global Cloud Search** | Real-time multi-provider search across VPS, Shared/Dedicated CPU Cloud VMs, GPU Servers, Bare Metal, Object/Block Storage, Managed DB, Kubernetes, Load Balancers, CDN, Backups, Managed/WordPress/Developer Hosting, and AI Infrastructure. |
| **Service 2** | **Cloud Price Comparison** | 36-dimension price normalization engine separating raw provider rates from usage-scenario estimates (hourly, 730h vs calendar month, bandwidth egress overage, IPv4 fees, backups, VAT/GST, live FX rates, freshness indicators). |
| **Service 3** | **Affiliate Marketplace** | Compliant affiliate link routing (`/go/{link_slug}`), `SubID`/`ClickID` injection, bot filtering, conversion postback webhook & CSV reconciliation, double-entry Commission Ledger, and FTC/Fair-Trade disclosure badges. |
| **Service 4** | **Cloud Account Connection (BYOC)** | Tenant-isolated connection of user-owned cloud accounts via official OAuth 2.0 or scoped API tokens, protected with KMS Envelope Encryption, ownership verification, and read/write scope separation. |
| **Service 5** | **Cloud Provisioning** | Explicit cost-approved server creation across supported providers with idempotency keys, asynchronous job state tracking, timeout reconciliation (preventing double provisioning), and immutable audit trails. |
| **Service 6** | **Cloud Management** | Unified multi-cloud console for instance status, specs, power actions (boot, graceful shutdown, reboot, power cycle), resizing, backups/snapshots, firewalls, SSH keys, IP management, and safe deletion gates. |
| **Service 7** | **Creator & Content Hub** | Reproducible YABS/Sysbench benchmark reports, deployment tutorials, provider reviews, and regional latency guides with Schema.org structured data. |
| **Service 8** | **Admin & Affiliate Operations** | Central operations console for provider catalog sync health, stale price alerts, affiliate campaign/payout reconciliation, BYOC operation audits, security incident monitoring, and SEO quality control. |

---

## 2. Strict Separation of the 4 Integration Adapters (Article 6)

CloudScanr enforces strict interface and package separation across four distinct provider integration adapters in `internal/adapters/`:

```
+-----------------------------------------------------------------------------------+
|                        CloudScanr Core Application Layer                          |
+-------------------+-------------------+---------------------+---------------------+
| Affiliate Adapter |  Pricing Adapter  | Management Adapter  |  Reseller Adapter   |
| (Monetization)    | (Catalog & Rates) | (BYOC User Account) | (Future Contracted) |
+-------------------+-------------------+---------------------+---------------------+
| - Link Builder    | - Public API Poll | - OAuth2 / API Key  | - Disabled by       |
| - SubID/ClickID   | - Auth API Poll   | - Create/Boot/Stop  |   default in V1     |
| - Postback Ingest | - Price Snapshot  | - Resize/Backup     | - Requires formal   |
| - CSV Reconcile   | - Spec Normalizer | - Delete Approval   |   Reseller Agreement|
+-------------------+-------------------+---------------------+---------------------+
```

1. **Affiliate Adapter (`internal/adapters/affiliate`)**: Handles outbound tracking URLs, campaign parameters, click telemetry, conversion postbacks, and network reporting APIs (Impact, CJ, PartnerStack, Custom). Does **not** touch server infrastructure.
2. **Pricing Adapter (`internal/adapters/pricing`)**: Harvests SKU specifications, regional availability, and raw pricing from official APIs or documented public price books. Runs on CloudScanr's system credentials or unauthenticated public endpoints.
3. **Management Adapter (`internal/adapters/management`)**: Executes infrastructure operations strictly on behalf of an authenticated user who has connected their own cloud account (BYOC) and granted explicit approval.
4. **Reseller Adapter (`internal/adapters/reseller`)**: Reserved interface for future direct resale under formal legal agreements. Never activated without explicit contract verification.

---

## 3. Repository & Blueprint Directory Map (Article 25)

- `START_HERE.md` — Mandatory 19-step onboarding and execution protocol for engineering agents.
- `README.md` — Executive architecture overview and package guide.
- `MANIFEST.json` — Machine-readable inventory of all files, sizes, categories, and SHA-256 digests.
- `FILE_CHECKSUMS.sha256` — Standard `sha256sum -c` compatible integrity file.
- `00_MARKET_RESEARCH/` — Deep due diligence on 18+ cloud providers, 5+ competitors, individual Korean developer affiliate eligibility, pricing APIs, management APIs, and partner/reseller programs.
- `01_A_TO_Z_BLUEPRINT/` — 26 exhaustive technical specifications (`A_ARCHITECTURE.md` through `Z_ZERO_TRUST.md`).
- `02_SYSTEM_ARCHITECTURE/` — Component topology, Mermaid/ASCII system diagrams, data flows, Kubernetes/Cloudflare deployment architecture, and STRIDE threat model.
- `03_DATABASE/` — Complete Entity-Relationship Diagram (`ERD.md`), `TABLE_DICTIONARY.md`, `INDEX_DESIGN.md`, and 12 executable PostgreSQL 16 DDL migration files (`migrations/`).
- `04_API_SPECIFICATION/` — Full OpenAPI 3.1 specification (`openapi.yaml`) covering all 22 API groups, adapter contracts, and API test suites.
- `05_UXUI/` — Desktop & Mobile wireframes, interaction specs, component states, design system tokens, and `SCREEN_INVENTORY.csv` for all 57 screens.
- `06_PRICING_ENGINE/` — 36-field normalization rules, multi-currency FX & tax policies, TCO scenario calculator, and freshness state machine.
- `07_AFFILIATE_ENGINE/` — Click attribution, `SubID`/`ClickID` cryptographic signing, postback/CSV reconciliation, double-entry commission ledger, and anti-fraud compliance.
- `08_CLOUD_MANAGEMENT/` — BYOC architecture, KMS envelope encryption, idempotent provisioning state machine, day-2 resource operations, budget controls, and provider capability matrix.
- `09_SECURITY/` — RBAC/ABAC IAM model, KMS secret lifecycle, PostgreSQL RLS tenant isolation, financial blast-radius controls, and security test suites.
- `10_ADMIN_OPERATIONS/` — Admin console workflows, catalog discrepancy resolution, affiliate payout reconciliation, and cloud operation auditing.
- `11_WORK_ORDERS_0001_TO_1000/` — Master index, dependency DAG, and 1,000 individual executable Work Orders (`individual/WO-CS-0001.md` to `WO-CS-1000.md`).
- `12_AGENT_EXECUTION/` — Multi-agent orchestration directive, 14-agent role matrix, file ownership locks, handoff protocol, and audit reporting templates.
- `13_TESTING/` — Master test matrix, acceptance criteria, Search/Affiliate/Provisioning/Security E2E specs (including all 22 mandatory server safety test scenarios), and independent audit procedures.
- `14_DEVOPS/` — GitHub Actions CI/CD pipelines, Prometheus/Grafana/OpenTelemetry monitoring, PITR backup/recovery, infrastructure cost management, and production runbooks.
- `15_FINAL_DELIVERY/` — Phase 0–10 roadmap, end-to-end requirements traceability matrix, known limitations & verification disclosures, delivery checklist, and final master audit report.
