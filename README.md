# TirpScanner — Global SaaS/AI-Tool & Software Affiliate Deal Scanner

**Project Code:** `BN-TIRPSCANNER-MASTER-BLUEPRINT-10-02`  
**Version:** `2.0.0` (TirpScanner pivot, 2026-11-26)  
**Architecture Paradigm:** Modular Go Monolith + Async Worker Fleet + Nuxt 3 SSR + PostgreSQL 16 + OpenSearch 2.17  
**Platform Independence:** 100% Independent standalone brand (Dedicated Domain, Auth, DB, Search, Backend, Admin, CI/CD)

> **Pivot notice (2026-11-26):** this repository pivoted from CloudScanr
> (cloud price comparison + BYOC multi-cloud management) to **TirpScanner**,
> a global SaaS/AI-tool & software affiliate deal scanner. The pivot decision
> record is [`TIRPSCANNER_PIVOT.md`](TIRPSCANNER_PIVOT.md); live domains
> are `tirpscanr.com` (registered; wiring pending), `api.tirpscanr.com`
> and `www.tirpscanr.com` (ASSUMED).
> The [`COMMERCIALIZATION.md`](COMMERCIALIZATION.md) 2026-11-18 CloudScanr cut
> is preserved as history; its 2026-11-26 addendum is the forward record.
> Blueprint directories below remain the design archive; where they describe
> cloud-infrastructure catalog, BYOC, provisioning, or server management as
> live scope, that scope is **RETIRED** — the live scope is the SaaS/AI-tool
> catalog, deal scan / price-drop detection, and affiliate monetization.

---

## 1. Platform Overview & Core Services

TirpScanner (**Tirp** = Tool Indexing & Ranking Platform) unifies global SaaS/AI-tool and software discovery, deal scanning
with price-drop detection, transparent plan/price comparison, compliant
affiliate monetization, and content-driven rankings in a single platform:

| Service ID | Service Name | Core Capabilities |
| :--- | :--- | :--- |
| **Service 1** | **Global SaaS/AI-Tool Search** | Real-time search across SaaS categories, AI tools (LLM apps, agents, image/video/audio, productivity, dev tools), and downloadable software, with category, platform, pricing-model, and region facets. |
| **Service 2** | **Deal Scan & Price-Drop Detection** | Continuous deal ingestion, coupon/promo tracking, plan-price normalization (monthly vs annual vs lifetime, seats, usage tiers, VAT/GST, live FX rates), price-drop alerts, and freshness indicators. |
| **Service 3** | **Affiliate Marketplace** | Compliant affiliate link routing (`/go/{link_slug}`), `SubID`/`ClickID` injection, bot filtering, conversion postback webhook & CSV reconciliation, double-entry Commission Ledger, and FTC/Fair-Trade disclosure badges. |
| **Service 4** | **Creator & Content Hub** | Tool comparisons, deal roundups, review methodology, buying guides, and alternatives pages with Schema.org structured data. |
| **Service 5** | **Admin & Affiliate Operations** | Central operations console for catalog/deal-ingestion sync health, stale-deal alerts, affiliate campaign/payout reconciliation, security incident monitoring, and SEO quality control. |

### Retired scope (CloudScanr era, no longer live)

Cloud Account Connection (BYOC), Cloud Provisioning (server creation), and
Cloud Management (instance lifecycle, power actions, resizing, backups,
firewalls, SSH keys, IP management, deletion gates) are **RETIRED**. Docs,
APIs, and UI flows describing them are historical. No new work targets them.

---

## 2. Revised Adapter Law (affiliate + catalog/deal-ingestion live; management/reseller retired)

TirpScanner enforces strict interface and package separation across provider
integration adapters in `internal/adapters/`. Two adapters are live; two are
retired:

```
+-----------------------------------------------------------------------------------+
|                        TirpScanner Core Application Layer                         |
+-------------------+-------------------+---------------------+---------------------+
| Affiliate Adapter | Catalog/Deal-     | Management Adapter  |  Reseller Adapter   |
| (Monetization)    | Ingestion Adapter | (BYOC) — RETIRED    | (Future) — RETIRED  |
| LIVE              | (Catalog & Deals) |                     |                     |
|                   | LIVE              |                     |                     |
+-------------------+-------------------+---------------------+---------------------+
| - Link Builder    | - Public API Poll | - RETIRED: no new   | - RETIRED: never    |
| - SubID/ClickID   | - Feed/CSV Import |   BYOC connections, |   activated; no     |
| - Postback Ingest | - Deal Snapshot   |   provisioning, or  |   resale or re-bill |
| - CSV Reconcile   | - Plan Normalizer |   server operations |   without executed  |
+-------------------+-------------------+---------------------+---------------------+
```

1. **Affiliate Adapter (`internal/adapters/affiliate`) — LIVE**: Handles outbound tracking URLs, campaign parameters, click telemetry, conversion postbacks, and network reporting APIs (Impact, CJ, PartnerStack, Custom). Does **not** touch server infrastructure.
2. **Catalog/Deal-Ingestion Adapter (`internal/adapters/pricing`) — LIVE**: Harvests SaaS/AI-tool listings, plan specifications, deal/coupon feeds, regional availability, and raw pricing from official APIs or documented public feeds. Runs on TirpScanner's system credentials or unauthenticated public endpoints. (Package name `pricing` is a historical label; its live role is catalog and deal ingestion.)
3. **Management Adapter (`internal/adapters/management`) — RETIRED**: Formerly executed infrastructure operations on user-connected cloud accounts (BYOC). Retired with the pivot; no new connections, provisioning, or server operations.
4. **Reseller Adapter (`internal/adapters/reseller`) — RETIRED**: Reserved interface that was never activated. Retired; never resell or re-bill third-party infrastructure or software without an executed Partner/Reseller contract.

---

## 3. Repository & Blueprint Directory Map

- `START_HERE.md` — Mandatory 19-step onboarding and execution protocol for engineering agents (deal-scanner scope).
- `README.md` — Executive architecture overview and package guide (this file).
- `TIRPSCANNER_PIVOT.md` — Pivot decision record, assumed domains, live-vs-retired scope map, frozen-doc policy, doc-change log.
- `COMMERCIALIZATION.md` — Implementation ledger; 2026-11-18 CloudScanr cut preserved as history, 2026-11-26 pivot addendum is forward record.
- `DEPLOY.md` — Deployment guide for `tirpscanr.com` / `api.tirpscanr.com` (toolchain unchanged).
- `MANIFEST.json` — Machine-readable inventory of all files, sizes, categories, and SHA-256 digests. **FROZEN — never edit.**
- `FILE_CHECKSUMS.sha256` — Standard `sha256sum -c` compatible integrity file. **FROZEN — never edit.**
- `00_MARKET_RESEARCH/` — Market due-diligence archive. Cloud-provider pricing/management API matrices are retired-scope history; affiliate-network research remains relevant.
- `01_A_TO_Z_BLUEPRINT/` — 26 technical specifications (`A_ARCHITECTURE.md` through `Z_ZERO_TRUST.md`). Live where they cover catalog, deals, affiliate, auth, admin, and platform; retired where they cover BYOC/provisioning/server management.
- `02_SYSTEM_ARCHITECTURE/` — Component topology, diagrams, data flows, deployment architecture, and threat model. Read with the pivot scope map; BYOC/provisioning paths are retired.
- `03_DATABASE/` — Entity-Relationship Diagram (`ERD.md`), `TABLE_DICTIONARY.md`, `INDEX_DESIGN.md`, executable PostgreSQL 16 DDL migrations, and rollback reference. Catalog/deal/affiliate/admin tables are live; BYOC/server tables are retired history.
- `04_API_SPECIFICATION/` — OpenAPI 3.1 specification (`openapi.yaml`) with per-path implemented/planned status. Catalog/deal/affiliate/admin groups are live; `/cloud/*` server groups are retired.
- `05_UXUI/` — Wireframes, interaction specs, component states, design system tokens, and screen inventory. Deal-scanner storefront, search, rankings, disclosure, and admin screens are live; cloud console/provisioning screens are retired.
- `06_PRICING_ENGINE/` — Normalization rules, multi-currency FX & tax policies, scenario calculator, and freshness state machine. Live as the deal/plan normalization engine (naming is historical).
- `07_AFFILIATE_ENGINE/` — Click attribution, `SubID`/`ClickID` cryptographic signing, postback/CSV reconciliation, double-entry commission ledger, and anti-fraud compliance. **LIVE.**
- `08_CLOUD_MANAGEMENT/` — BYOC architecture, provisioning state machine, and day-2 operations. **RETIRED scope — historical archive only.**
- `09_SECURITY/` — IAM model, secret lifecycle, tenant isolation, financial controls, and security test suites. Live as platform security; BYOC-credential sections are retired history.
- `10_ADMIN_OPERATIONS/` — Admin console workflows, catalog discrepancy resolution, affiliate payout reconciliation. Live; cloud-operation auditing sections are retired history.
- `11_WORK_ORDERS_0001_TO_1000/` — Master index, dependency DAG, and 1,000 individual Work Orders. **FROZEN — never edit.**
- `12_AGENT_EXECUTION/` — Multi-agent orchestration directive, role matrix, file ownership locks, handoff protocol, and audit reporting templates.
- `13_TESTING/` — Master test matrix, acceptance criteria, E2E specs, and independent audit procedures. Search/affiliate/security suites are live; provisioning suites are retired history.
- `14_DEVOPS/` — CI/CD pipelines, monitoring, backup/recovery, cost management, and production runbooks. Live as platform operations.
- `15_FINAL_DELIVERY/` — Phase roadmap, requirements traceability matrix, limitations & disclosures, delivery checklist, and final audit report. `REQUIREMENTS_TRACEABILITY.csv` is **FROZEN — never edit**; dated addenda carry forward status.
