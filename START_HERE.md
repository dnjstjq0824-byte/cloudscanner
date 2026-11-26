# START_HERE.md — TirpScanner Autonomous Engineering & Agent Handoff Protocol

**Document ID:** `TS-DOC-ROOT-001`  
**Version:** `2.0.0` (TirpScanner pivot, 2026-11-26)  
**Target Platform:** TirpScanner (`tirpscanr.com` — Independent Global SaaS/AI-Tool & Software Affiliate Deal Scanner; domain registered, wiring pending; see `TIRPSCANNER_PIVOT.md`)  
**Master Package:** `BN_TIRPSCANNER_MASTER_BLUEPRINT_10_V2_0.zip` (successor to `BN_CLOUDSCANR_MASTER_BLUEPRINT_10_V1_0.zip`)  
**Verification Date:** `2026-10-28` (original blueprint cut; pivot addendum `2026-11-26`)

> **Pivot note (2026-11-26):** this protocol was rewritten for the TirpScanner
> pivot. Live scope is the SaaS/AI-tool catalog, deal scan / price-drop
> detection, and affiliate monetization. Cloud-infrastructure catalog, BYOC
> account connection, server provisioning, and server management are
> **RETIRED** — steps referencing them are generalized below to deal-scanner
> scope or marked retired. The 2026-11-18 CloudScanr cut is preserved as
> history in [`COMMERCIALIZATION.md`](COMMERCIALIZATION.md).

---

## 0. Executive Mandate & Platform Independence

Welcome, Engineering Agent. You have received the authoritative **TirpScanner Master Blueprint**.

### Strict Independence as a Standalone Brand
TirpScanner is a standalone, independent global SaaS and marketplace brand:
1. **Separate Domain & Brand**: Operates on its own dedicated domain (`tirpscanr.com` — registered; wiring pending), TLS certificates, and CDN/WAF edge, with API on `api.tirpscanr.com` (ASSUMED) and `www.tirpscanr.com` (ASSUMED) as the www alias.
2. **Separate Identity & Database**: Uses its own dedicated PostgreSQL 16 cluster (`tirpscanner_prod`) and OpenSearch 2.x cluster. Do **NOT** import, link, or query any other product's database, user table, or SSO session.
3. **Zero Modification to Other Products**: Never touch any other repository, database, or deployment pipeline.
4. **Search Engine Discoverability Only**: Web crawlers may index TirpScanner's public pages just like Googlebot or Bingbot via standard `sitemap.xml`, `robots.txt`, Canonical URLs, and Schema.org `Product`/`Offer`/`SoftwareApplication` JSON-LD. Private routes (`/dashboard/*`, `/admin/*`, tokens, billing, affiliate internals) are strictly disallowed from public indexing.

---

## 1. Mandatory 19-Step Execution Procedure for Downstream Agents

Every autonomous development agent or human engineering team ingesting this blueprint **MUST** execute the following 19 steps sequentially before marking any milestone complete:

### Step 1: ZIP & Archive Integrity Verification
Run SHA-256 verification against all packaged files:
```bash
sha256sum -c FILE_CHECKSUMS.sha256
```
Verify that zero files are corrupted, missing, or truncated. Compare file counts against `MANIFEST.json`. Note: `MANIFEST.json` and `FILE_CHECKSUMS.sha256` are **FROZEN** blueprint-delivery records — verify against them, never edit them.

### Step 2: Package File Inventory Audit
Parse `MANIFEST.json` and confirm all 12 top-level directories (`00_MARKET_RESEARCH` through `15_FINAL_DELIVERY`), 26 A–Z specification files, PostgreSQL SQL migration files, the OpenAPI 3.1 YAML specification, and exactly 1,000 Work Order files (`WO-CS-0001.md` to `WO-CS-1000.md`) are present. Note: `11_WORK_ORDERS_0001_TO_1000/` is **FROZEN** — audit it, never edit it.

### Step 3: README & Core Principles Review
Read `README.md` and `TIRPSCANNER_PIVOT.md`, and internalize the **Revised Adapter Law**:
- **Affiliate Adapter (LIVE)**: Outbound tracking links, `SubID`/`ClickID`, conversion reconciliation, and commission ledger.
- **Catalog/Deal-Ingestion Adapter (LIVE)**: SaaS/AI-tool catalog ingestion, deal/coupon feed snapshots, plan-price normalization, price-drop detection.
- **Management Adapter (BYOC) — RETIRED**: No new account connections, provisioning, or server operations.
- **Reseller Adapter — RETIRED**: Never resell or re-bill third-party infrastructure or software without an executed Partner/Reseller contract.

### Step 4: A–Z Technical Blueprint Review
Read `01_A_TO_Z_BLUEPRINT/A_ARCHITECTURE.md` through `Z_ZERO_TRUST.md`. Verify the domain boundaries, technology stack (Go 1.23+, Nuxt 3.13+ SSR, PostgreSQL 16, OpenSearch 2.17, Redis 7.4 / Asynq, AWS KMS / HashiCorp Vault), and failure isolation rules. Treat BYOC/provisioning/server-management sections as retired history.

### Step 5: Database Schema & OpenAPI Contract Lock
Inspect `03_DATABASE/ERD.md`, `03_DATABASE/migrations/*.sql`, and `04_API_SPECIFICATION/openapi.yaml`. These files are the **Single Source of Truth (SSOT)** for table names, column types, constraints, and REST payloads. Do not invent undocumented database columns or ad-hoc API routes without following the Contract Change Protocol in `12_AGENT_EXECUTION/HANDOFF_PROTOCOL.md`. `/cloud/*` server groups are retired; catalog/deal/affiliate/admin groups are live.

### Step 6: Work Order Index Loading
Load `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv` and `MASTER_WORK_ORDER.md`. Each Work Order (`WO-CS-0001` through `WO-CS-1000`) defines exact target files, new files, DB migrations, API endpoints, prerequisites, security rules, and verification commands. This directory is **FROZEN** — load it as a historical index; map live deal-scanner work forward without editing it.

### Step 7: Execution Dependency Graph Analysis
Consult `11_WORK_ORDERS_0001_TO_1000/DEPENDENCY_GRAPH.md` and `12_AGENT_EXECUTION/EXECUTION_SEQUENCE.md`. Identify critical-path blockers for live scope (e.g., catalog/deal-ingestion pipeline and affiliate attribution must be completed and tested before deal-alert and payout-reconciliation work). BYOC/provisioning dependency chains are retired.

### Step 8: Target Repository & Environment Due Diligence
Inspect the current workspace/repository state. Check existing Go modules (`go.mod`), Nuxt configuration (`nuxt.config.ts`), Docker Compose / Kubernetes manifests, and database migration state.

### Step 9: Implementation Gap Analysis
Compare the target repository against live-scope work (SaaS/AI-tool catalog, deal scan / price-drop detection, affiliate engine, content hub, admin console). Classify every item into `NOT_STARTED`, `IN_PROGRESS`, `BLOCKED`, or `VERIFIED_COMPLETE`. Record findings in your execution state tracker. Do not plan new BYOC/provisioning/server-management work.

### Step 10: Priority Work Selection (Phase-Driven)
Select the next unblocked work batch for deal-scanner scope (Market Verification -> Foundation -> Catalog & Deal-Scan MVP -> Affiliate Engine -> Identity & Personalization -> Deal Alerts & Content -> Advanced Ops -> Global & SEO -> Commercial Audit -> Launch). Retired BYOC/provisioning phases are skipped.

### Step 11: Production-Grade Code Implementation
Implement the selected work strictly adhering to its `Target Files`, `Backend Requirements`, `Frontend Requirements`, and `DB Changes`. Zero placeholder `TODO` stubs or hardcoded mock credentials in production code paths.

### Step 12: Unit & Integration Testing
Execute unit tests (`go test -v -race ./...`, `vitest run`) and integration tests against local PostgreSQL/OpenSearch test containers. Ensure all edge cases specified in the work item pass.

### Step 13: External Deal-Feed & Affiliate-Network Scope & Rate-Limit Verification
Before invoking any live external deal feed or affiliate network API (affiliate networks, SaaS marketplaces, vendor feeds), verify rate limits, authentication scopes, and feed terms against the market-research matrices and each network's documented terms. Cloud-provider server management APIs are retired scope — do not integrate them.

### Step 14: Deal-Mutation Accuracy & Affiliate-Safety Gate Verification
Verify that any deal/affiliate mutation (deal publish/expire, coupon attach, payout reconcile, notice publish) enforces:
1. Pre-execution accuracy display (plan/price shown, discount math, expiry, disclosure badge).
2. Explicit approval or moderation step where the workflow requires it.
3. Client-supplied `Idempotency-Key` header checked before mutating.
4. Staleness reconciliation (re-verify feed state) prior to any retry to prevent duplicate or expired deal publication.
Retired: server provisioning cost-control and safety gates (no billable server actions exist in live scope).

### Step 15: Full-Stack UX/UI End-to-End (E2E) Verification
Run Playwright E2E test suites for live scope (search/deals, affiliate, security) across Desktop (`1440x900`) and Mobile (`390x844`) viewports. Cloud-provisioning E2E suites are retired.

### Step 16: Independent Zero-Trust Audit
Run the Independent Audit Agent checklist (`01_A_TO_Z_BLUEPRINT/Z_ZERO_TRUST.md` and `13_TESTING/INDEPENDENT_AUDIT.md`). Verify tenant isolation (Cross-Tenant IDOR test must return `403/404`), secret redaction in logs, and SQL ledger integrity.

### Step 17: Approved Git Commit & Deployment Gate
Commit code with the standardized work-item prefix only after all automated checks pass.

### Step 18: Per-Work-Item Evidence Archival
Record test output logs, migration verification outputs, HTTP response traces, and screenshot artifacts as required by the `Evidence` field of each work item.

### Step 19: Residual Work & Blocker Reporting
Generate an updated status report using `12_AGENT_EXECUTION/REPORTING_TEMPLATE.md`, explicitly flagging any items that are `BLOCKED` due to external merchant approval, legal entity registration, or affiliate-network admission.

---

## 2. Hard Safety & Honesty Rules

1. **Never fabricate affiliate approval**: Do not hardcode active affiliate status for programs requiring manual review (`REQUIRES_APPLICATION` or `CONTRACT_REQUIRED`).
2. **Never invent non-existent APIs or feeds**: Where a vendor lacks a public catalog/deal API or conversion postback API, use the documented fallback (manual/CSV import or verified catalog feed) and mark API capability accurately.
3. **Never expose secrets**: Tokens and credentials must be encrypted via Envelope Encryption (DEK encrypted by KMS KEK) with ciphertext/key-ID storage only; never log raw secrets. (Retired specificity: BYOC cloud-account tokens no longer exist in live scope.)
4. **Never publish unverified deals or payouts**: Live deal publication and payout reconciliation require verified feed state and an authorized account; never fake availability, pricing, or commission status.
