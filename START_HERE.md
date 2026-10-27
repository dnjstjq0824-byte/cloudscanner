# START_HERE.md — CloudScanr Autonomous Engineering & Agent Handoff Protocol

**Document ID:** `CS-DOC-ROOT-001`  
**Version:** `1.0.0`  
**Target Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Price Comparison, Affiliate Monetization & BYOC Multi-Cloud Management Platform)  
**Master Package:** `BN_CLOUDSCANR_MASTER_BLUEPRINT_1000_V1_0.zip`  
**Verification Date:** `2026-10-28`

---

## 0. Executive Mandate & Platform Independence

Welcome, Engineering Agent. You have received the authoritative **CloudScanr Master Blueprint (`BN_CLOUDSCANR_MASTER_BLUEPRINT_1000_V1_0.zip`)**.

### Strict Independence from TreeScan (Article 1 Compliance)
CloudScanr is **NOT** an internal sub-service of TreeScan. CloudScanr is a standalone, independent global SaaS and marketplace brand:
1. **Separate Domain & Brand**: Operates on its own dedicated domain (`cloudscanr.com`), TLS certificates, and CDN/WAF edge.
2. **Separate Identity & Database**: Uses its own dedicated PostgreSQL 16 cluster (`cloudscanr_prod`) and OpenSearch 2.x cluster. Do **NOT** import, link, or query any TreeScan database, user table, or SSO session.
3. **Zero Code Modification to TreeScan**: Never touch any TreeScan repository, database, or deployment pipeline.
4. **Search Engine Discoverability Only**: Long-term, TreeScan's web crawler may index CloudScanr's public pages just like Googlebot or Bingbot via standard `sitemap.xml`, `robots.txt`, Canonical URLs, and Schema.org `Product`/`Offer` JSON-LD. Private routes (`/dashboard/*`, `/api/v1/cloud/*`, `/admin/*`, tokens, billing, server inventory) are strictly disallowed from public indexing.

---

## 1. Mandatory 19-Step Execution Procedure for Downstream Agents (Article 26)

Every autonomous development agent or human engineering team ingesting this blueprint **MUST** execute the following 19 steps sequentially before marking any milestone complete:

### Step 1: ZIP & Archive Integrity Verification
Run SHA-256 verification against all packaged files:
```bash
sha256sum -c FILE_CHECKSUMS.sha256
```
Verify that zero files are corrupted, missing, or truncated. Compare file counts against `MANIFEST.json`.

### Step 2: Package File Inventory Audit
Parse `MANIFEST.json` and confirm all 16 top-level directories (`00_MARKET_RESEARCH` through `15_FINAL_DELIVERY`), 26 A–Z specification files, 12 PostgreSQL SQL migration files, 1 OpenAPI 3.1 YAML specification, and exactly 1,000 Work Order files (`WO-CS-0001.md` to `WO-CS-1000.md`) are present.

### Step 3: README & Core Principles Review
Read `README.md` and internalize the **Four-Adapter Separation Law (Article 6)**:
- **Affiliate Adapter**: Outbound tracking links, `SubID`/`ClickID`, conversion reconciliation, and commission ledger.
- **Pricing Adapter**: Public/authenticated catalog ingestion, hourly/monthly price snapshots, bandwidth/IPv4/storage normalization.
- **Management Adapter (BYOC)**: User-authorized OAuth2 or scoped API token actions on the user's own cloud account.
- **Reseller Adapter**: Separated future contract module; never resell or re-bill third-party cloud infrastructure without an executed Partner/Reseller contract.

### Step 4: A–Z Technical Blueprint Review
Read `01_A_TO_Z_BLUEPRINT/A_ARCHITECTURE.md` through `Z_ZERO_TRUST.md`. Verify the domain boundaries, technology stack (Go 1.23+, Nuxt 3.13+ SSR, PostgreSQL 16, OpenSearch 2.17, Redis 7.4 / Asynq, AWS KMS / HashiCorp Vault), and failure isolation rules.

### Step 5: Database Schema & OpenAPI Contract Lock
Inspect `03_DATABASE/ERD.md`, `03_DATABASE/migrations/*.sql`, and `04_API_SPECIFICATION/openapi.yaml`. These files are the **Single Source of Truth (SSOT)** for table names, column types, constraints, and REST payloads. Do not invent undocumented database columns or ad-hoc API routes without following the Contract Change Protocol in `12_AGENT_EXECUTION/HANDOFF_PROTOCOL.md`.

### Step 6: Work Order Index Loading
Load `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv` and `MASTER_WORK_ORDER.md`. Each Work Order (`WO-CS-0001` through `WO-CS-1000`) defines exact target files, new files, DB migrations, API endpoints, prerequisites, security rules, and verification commands.

### Step 7: Execution Dependency Graph Analysis
Consult `11_WORK_ORDERS_0001_TO_1000/DEPENDENCY_GRAPH.md` and `12_AGENT_EXECUTION/EXECUTION_SEQUENCE.md`. Identify critical-path blockers (e.g., KMS envelope encryption `WO-CS-0551` must be completed and tested before any cloud credential storage `WO-CS-0555` or server provisioning `WO-CS-0601`).

### Step 8: Target Repository & Environment Due Diligence
Inspect the current workspace/repository state. Check existing Go modules (`go.mod`), Nuxt configuration (`nuxt.config.ts`), Docker Compose / Kubernetes manifests, and database migration state.

### Step 9: Implementation Gap Analysis
Compare the target repository against the 1,000 Work Orders. Classify every Work Order into `NOT_STARTED`, `IN_PROGRESS`, `BLOCKED`, or `VERIFIED_COMPLETE`. Record findings in your execution state tracker.

### Step 10: Priority Work Order Selection (Phase-Driven)
Select the next unblocked Work Order batch according to `15_FINAL_DELIVERY/PHASE_ROADMAP.md` (Phase 0 Market Verification -> Phase 1 Foundation -> Phase 2 Price Comparison MVP -> Phase 3 Affiliate Engine -> Phase 4 Identity & Personalization -> Phase 5 BYOC Connection -> Phase 6 Provisioning & Management -> Phase 7 Advanced Ops -> Phase 8 Global & SEO -> Phase 9 Commercial Audit -> Phase 10 Launch).

### Step 11: Production-Grade Code Implementation
Implement the selected Work Order strictly adhering to its `Target Files`, `Backend Requirements`, `Frontend Requirements`, and `DB Changes`. Zero placeholder `TODO` stubs or hardcoded mock credentials in production code paths.

### Step 12: Unit & Integration Testing
Execute unit tests (`go test -v -race ./...`, `vitest run`) and integration tests against local PostgreSQL/OpenSearch test containers. Ensure all edge cases specified in the Work Order pass.

### Step 13: External Provider API Scope & Rate-Limit Verification
Before invoking any live external Provider API (DigitalOcean, UpCloud, Vultr, Hetzner, Linode, Cloudways, AWS, GCP, Azure, etc.), verify rate limits, authentication scopes, and read-only vs. write permissions against `00_MARKET_RESEARCH/PRICING_API_MATRIX.csv` and `MANAGEMENT_API_MATRIX.csv`.

### Step 14: Server Provisioning Cost Control & Safety Gate Verification
Verify that any mutating Cloud Management operation (`CreateServer`, `ResizeServer`, `EnableBackups`, `DeleteServer`, `RebuildServer`) enforces:
1. Pre-execution cost estimation display (`hourly_cost`, `monthly_cap`, `backup_fee`, `ipv4_fee`).
2. Explicit cryptographic or database-backed user approval token (`operation_approvals` table).
3. Client-supplied `Idempotency-Key` header checked before calling external APIs.
4. Timeout reconciliation (`ReconcileResourceState`) prior to any retry to prevent duplicate server creation and double billing.

### Step 15: Full-Stack UX/UI End-to-End (E2E) Verification
Run Playwright E2E test suites (`13_TESTING/SEARCH_E2E.md`, `AFFILIATE_E2E.md`, `CLOUD_PROVISIONING_E2E.md`, `SECURITY_E2E.md`) across Desktop (`1440x900`) and Mobile (`390x844`) viewports.

### Step 16: Independent Zero-Trust Audit
Run the Independent Audit Agent checklist (`01_A_TO_Z_BLUEPRINT/Z_ZERO_TRUST.md` and `13_TESTING/INDEPENDENT_AUDIT.md`). Verify tenant isolation (Cross-Tenant IDOR test must return `403/404`), secret redaction in logs, and SQL ledger integrity.

### Step 17: Approved Git Commit & Deployment Gate
Commit code with the standardized Work Order prefix (`feat(WO-CS-XXXX): ...`) only after all automated checks pass.

### Step 18: Per-Work-Order Evidence Archival
Record test output logs, migration verification outputs, HTTP response traces, and screenshot artifacts as required by the `Evidence` field of each Work Order.

### Step 19: Residual Work & Blocker Reporting
Generate an updated status report using `12_AGENT_EXECUTION/REPORTING_TEMPLATE.md`, explicitly flagging any items that are `BLOCKED` due to external merchant approval, legal entity registration, or live budget authorization.

---

## 2. Hard Safety & Honesty Rules (Article 29)

1. **Never fabricate affiliate approval**: Do not hardcode active affiliate status for programs requiring manual review (`REQUIRES_APPLICATION` or `CONTRACT_REQUIRED`).
2. **Never invent non-existent APIs**: Where a provider lacks a public pricing API or conversion postback API, use the documented fallback (manual/CSV import or verified catalog feed) and mark API capability accurately.
3. **Never expose user API secrets**: Customer BYOC tokens must be encrypted via Envelope Encryption (DEK encrypted by KMS KEK) and stored in `provider_credentials_metadata` with ciphertext/key-ID only; never log raw tokens.
4. **Never execute unapproved billable or destructive cloud actions**: Live server creation/deletion tests require an explicitly authorized sandbox account with strict budget caps.
