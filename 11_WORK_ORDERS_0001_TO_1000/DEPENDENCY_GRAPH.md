# DEPENDENCY_GRAPH.md — 1,000 Work Order Directed Acyclic Graph (DAG) & Critical Path

**Document ID:** `CS-WO-DAG`

---

## 1. Macro Dependency DAG Across the 20 Work Order Bands

```mermaid
graph TD
    B01["0001-0050: Market & Affiliate Research"] --> B02["0051-0100: Architecture, DB DDL & Infra"]
    B02 --> B03["0101-0150: Identity, Auth, MFA & Orgs"]
    B02 --> B04["0151-0200: Providers & 16-Category Catalog"]
    B04 --> B05["0201-0250: Official Pricing API Connectors"]
    B05 --> B06["0251-0300: 36-Field Normalization, FX, Tax & TCO"]
    B06 --> B07["0301-0350: OpenSearch, NL Parser & Ranking"]
    B02 --> B08["0351-0400: Affiliate Signup, Contracts & Links"]
    B08 --> B09["0401-0450: Clicks, Postbacks & Commission Ledger"]
    B07 --> B10["0451-0500: Public Web, Home, Search & SKU Detail"]
    B09 --> B10
    B10 --> B11["0501-0550: Comparison Matrix, Calculator & SEO"]
    B03 --> B12["0551-0600: BYOC Account Connection & KMS Security"]
    B12 --> B13["0601-0650: Server Provisioning, Power & Safe Delete"]
    B13 --> B14["0651-0700: Resize, Network, Backups & Budgets"]
    B09 --> B15["0701-0750: Admin Console & Provider Operations"]
    B14 --> B15
    B11 --> B16["0751-0800: Reviews, Benchmarks, i18n & WCAG AA"]
    B14 --> B17["0801-0850: DB Partitioning, Job Queue & SRE Ops"]
    B15 --> B18["0851-0900: Unit, Integration, E2E & 22 Safety Tests"]
    B16 --> B18
    B17 --> B19["0901-0950: Cost, Load, Chaos & Security Pen-Tests"]
    B18 --> B19
    B19 --> B20["0951-1000: Independent Zero-Trust Audit & Launch"]
```

## 2. Critical Path Lock Gates
1. **Lock Gate 1 (`WO-CS-0071`)**: All 12 PostgreSQL DDL migrations (`49 tables`) must be applied before domain repositories (`0101+`) execute.
2. **Lock Gate 2 (`WO-CS-0551`)**: KMS Envelope Encryption with `org_id:connection_id:provider_id` AAD must pass `SEC-02` before any BYOC token (`0565+`) or server provisioning (`0601+`) is enabled.
3. **Lock Gate 3 (`WO-CS-0601` & `WO-CS-0615`)**: Explicit Cost Approval (`operation_approvals`) and Timeout Reconciliation (`TIMEOUT_RECONCILING`) must pass all 22 Article 24 safety tests (`0881–0902`) before live BYOC provisioning is unlocked.
