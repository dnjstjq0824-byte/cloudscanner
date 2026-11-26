# X — eXecution & Agent Orchestration — Specialized Agents, File Ownership Locks & Work Order DAG

**Document ID:** `CS-AZ-X` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (TirpScanner pivot reframe, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Pivot Note (2026-11-26):** Cloud Management Agent scope retired; Vendor Feed Agent covers deal ingestion. `WO-CS-*` IDs and frozen work-order history unchanged. See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0001 ~ WO-CS-1000 (Specifically orchestrated in WO-CS-0951 ~ WO-CS-1000)`

---

## 1. 기능 목적 (Functional Purpose)

Define the deterministic multi-agent orchestration framework governing how specialized engineering agents execute the 1,000 Work Orders (`WO-CS-0001` through `WO-CS-1000`) in parallel without merge conflicts, schema drift, or unverified completion claims — with retired-scope orders (`M`, `K`, `N`, `V`, `Q`-cloud-safety) marked terminal history.

---

## 2. User Roles & Actors

1. Lead Architect, 2. Affiliate Research Agent, 3. Catalog/Deal-Feed Ingestion Agent, 4. Deal/Plan Normalization Agent, 5. Frontend Agent, 6. Backend Agent, 7. Database Agent, 8. Search Agent, 9. ~~Cloud Management Agent~~ RETIRED, 10. Security Agent, 11. DevOps Agent, 12. UX/UI Agent, 13. QA Agent, 14. Independent Audit Agent.

---

## 3. User Scenarios

1. Database Agent completes `WO-CS-0061` (Identity & Providers DDL) and locks `03_DATABASE/migrations/000002_*.sql`; Backend Agent and Frontend Agent then execute `WO-CS-0101` (Auth API) and `WO-CS-0451` (Home & Auth UI) concurrently against the locked OpenAPI 3.1 contract.
2. Vendor Feed Agent implements a PartnerStack coupon-delta ingestor; QA Agent and Independent Audit Agent run the deal/affiliate safety tests before the Work Order is allowed to transition to `VERIFIED_COMPLETE`.

---

## 4. Technology Stack

Git branch-per-WO workflow, automated CI contract gates (`oasdiff`, `sqlfluff`, `golangci-lint`), `WORK_ORDER_INDEX.csv` DAG scheduler.

---

## 5. Database Schema & Tables

**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`All schema changes governed by Database Agent lock protocol`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)

**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**

- `All OpenAPI changes governed by Lead Architect + Backend Agent contract lock`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. User Interface & Experience

Standardized Work Order Markdown specification and Handoff Reporting Template (`12_AGENT_EXECUTION/`).

---

## 8. Core Processing Logic

1. Strict Directory & File Ownership: Each agent owns designated repository paths (`12_AGENT_EXECUTION/AGENT_ROLE_MATRIX.md`). Cross-boundary edits require a formal Contract Change Request (CCR). No agent may extend retired modules (`M`, `K`, `N`, `V`) or retired adapters.
2. Definition of Done (DoD) Gate: A Work Order is only complete when Code + Unit Test + Integration/E2E Test + Security Check + Independent Audit Evidence are all green.

---

## 9. Security & Access Control

Independent Audit Agent has read-only access to all code and test logs, and holds veto power over any Work Order that lacks real test evidence or fabricates completion.

---

## 10. Failure Modes & Resilience

If a Work Order is blocked by external third-party credentials (e.g., live Partner contract), the agent implements the full adapter against the specification + simulator harness, and marks live production credential verification as `BLOCKED_EXTERNAL_APPROVAL`.

---

## 11. Inter-Module Dependencies

Orchestrates `11_WORK_ORDERS_0001_TO_1000/` (FROZEN, read-only) and `12_AGENT_EXECUTION/`.

---

## 12. Testing & Verification Strategy

DAG cycle-detection script verifying the 1,000 Work Orders form a valid Directed Acyclic Graph with zero circular dependencies or orphan prerequisites.

---

## 13. Acceptance Criteria

- [x] All agent roles have explicit scopes, file ownership rules, and test responsibilities (Cloud Management Agent retired).
- [x] All 1,000 Work Orders (`WO-CS-0001` to `WO-CS-1000`) are mapped to a primary responsible agent and verified free of circular dependencies.

---

## 14. Required Verification Evidence

`12_AGENT_EXECUTION/` documents, `WORK_ORDER_INDEX.csv`, and DAG verification script output.

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0001 ~ WO-CS-1000 (Specifically orchestrated in WO-CS-0951 ~ WO-CS-1000)`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
