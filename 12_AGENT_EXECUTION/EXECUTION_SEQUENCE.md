# EXECUTION_SEQUENCE.md — Critical Path & Parallel Wave Execution Plan
**Document ID:** `CS-AG-003`
> TirpScanner pivot (2026-11-26): Wave 5 (BYOC/provisioning) retired; Wave 7 safety-suite scope reduced to search/affiliate/security/load E2E. Work-order ranges are frozen history. See `TIRPSCANNER_PIVOT.md`.
## 1. Wave-Based Parallel Execution Plan
- **Wave 0 (Market & API Verification)**: `WO-CS-0001` ~ `WO-CS-0050` (Affiliate Research Agent + Provider API Agent in parallel).
- **Wave 1 (Foundation, DB DDL & OpenAPI Lock)**: `WO-CS-0051` ~ `WO-CS-0100` (Lead Architect + Database Agent + DevOps Agent).
- **Wave 2 (Identity, Catalog & Pricing Ingestion)**: `WO-CS-0101` ~ `WO-CS-0250` (Backend Agent + Provider API Agent in parallel).
- **Wave 3 (Deal/Plan Normalization, OpenSearch & Affiliate Core)**: `WO-CS-0251` ~ `WO-CS-0450` (Deal/Plan Normalization Agent + Search Agent + Backend Agent in parallel).
- **Wave 4 (Public Web, Comparison UX & SEO)**: `WO-CS-0451` ~ `WO-CS-0550` (UX/UI Agent + Frontend Agent).
- **~~Wave 5 (KMS Security, BYOC Connection & Server Provisioning/Management)~~ (RETIRED by pivot)**: `WO-CS-0551` ~ `WO-CS-0700` preserved as frozen history only; KMS/platform-credential hardening continues under live security work.
- **Wave 6 (Admin Console, Content, i18n & Operations)**: `WO-CS-0701` ~ `WO-CS-0850` (Backend Agent + Frontend Agent + Database/DevOps Agents).
- **Wave 7 (Full QA, Search/Affiliate/Security E2E, Load & Security Audits)**: `WO-CS-0851` ~ `WO-CS-0950` (QA Agent + Security Agent; 22-scenario provisioning safety tests retired).
- **Wave 8 (Independent Zero-Trust Audit & Release Packaging)**: `WO-CS-0951` ~ `WO-CS-1000` (Independent Audit Agent + Lead Architect).
