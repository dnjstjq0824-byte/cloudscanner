# AGENT_ROLE_MATRIX.md — 14 Specialized Agent Roles, File Ownership & Test Responsibilities (Article 22)
**Document ID:** `CS-AG-002`
> TirpScanner pivot (2026-11-26): Cloud Management Agent row retired; ingestion/pricing/QA rows reframed for the deal scanner. Work-order ranges are frozen history (`11_WORK_ORDERS_0001_TO_1000/`). See `TIRPSCANNER_PIVOT.md`.

| # | Agent Role | Primary Work Order Ranges | Owned Repository Paths | Prerequisite Lock | Verification & Test Responsibility |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | **Lead Architect** | `0051–0060`, `0951–0960` | `START_HERE.md`, `README.md`, `01_A_TO_Z_BLUEPRINT/`, `02_SYSTEM_ARCHITECTURE/` | Phase 0 Market Due Diligence | Architecture fitness tests, package boundary import checks |
| 2 | **Affiliate Research Agent** | `0001–0025`, `0351–0375` | `00_MARKET_RESEARCH/AFFILIATE_*.csv`, `07_AFFILIATE_ENGINE/POLICY_COMPLIANCE.md` | None (Phase 0 Root) | Official URL reachability & status classification audit |
| 3 | **Catalog/Deal-Feed Ingestion Agent** | `0026–0050`, `0201–0250` | `internal/adapters/pricing/` (name historical), `00_MARKET_RESEARCH/*_API_MATRIX.csv` | `WO-CS-0001..0025` | Recorded JSON fixture contract tests for SaaS/AI-tool vendors & deal feeds |
| 4 | **Deal/Plan Normalization Agent** | `0251–0300` | `internal/domain/estimation/`, `internal/domain/ranking/`, `06_PRICING_ENGINE/` | `WO-CS-0151..0250` | Plan normalization, coupon/deal math, FX, tax & unbiased ranking tests |
| 5 | **Frontend Agent** | `0451–0550`, `0751–0800` | `web/pages/`, `web/components/`, `web/stores/`, `nuxt.config.ts` | OpenAPI 3.1 + Design System Lock | Vitest component tests, Lighthouse SEO = 100, WCAG 2.1 AA axe-core |
| 6 | **Backend Agent** | `0101–0200`, `0376–0450`, `0701–0750` | `cmd/api/`, `internal/domain/identity/`, `internal/domain/catalog/`, `internal/domain/affiliate/` | DB Migrations `000001..000008` | Go race-detector unit/integration tests + OpenAPI contract validation |
| 7 | **Database Agent** | `0061–0074`, `0801–0815` | `03_DATABASE/`, `migrations/*.sql`, `internal/platform/database/` | ERD & Domain Spec Lock | Up/Down migration idempotency, RLS tenant isolation, immutability triggers |
| 8 | **Search Agent** | `0301–0350` | `internal/domain/search/`, `opensearch/` | Catalog & Pricing Schema Lock | Sub-50ms p95 facet benchmark, 100+ NL spec query parser tests, alias swap |
| 9 | **~~Cloud Management Agent~~ (RETIRED)** | `0551–0700` | `internal/adapters/management/` (retired), `08_CLOUD_MANAGEMENT/` (retired history) | Identity (`0101–0150`) + KMS (`0551`) | ~~All 22 mandatory Server Safety Scenarios~~ — retired by pivot |
| 10 | **Security Agent** | `0551–0565`, `0826–0840`, `0921–0940` | `internal/platform/kms/`, `09_SECURITY/` | Core Architecture Lock | KMS AAD tamper test, IDOR matrix fuzzing, SSRF blocklist, log redaction |
| 11 | **DevOps Agent** | `0075–0100`, `0841–0850` | `14_DEVOPS/`, `.github/workflows/`, `deploy/k8s/`, `deploy/terraform/` | Architecture Topology Lock | CI pipeline execution, Terraform validate, PITR restore drill |
| 12 | **UX/UI Agent** | `0451–0460`, `0701–0710` | `05_UXUI/`, `web/assets/css/`, `tailwind.config.ts` | Product Requirements Lock | 44-screen live inventory coverage (16 public + 7 user + 21 admin; 13 retired), Desktop/Mobile viewport & dark-mode parity |
| 13 | **QA Agent** | `0851–0950` | `13_TESTING/`, `tests/e2e/`, `tests/load/` | Feature Implementation WOs | Playwright E2E (search/affiliate/security), k6 load testing (provisioning fault injection retired) |
| 14 | **Independent Audit Agent** | `0951–1000` | `15_FINAL_DELIVERY/`, `MANIFEST.json`, `FILE_CHECKSUMS.sha256` | All prior phases | End-to-end traceability verification, SHA-256 integrity, Article 29 audit |
