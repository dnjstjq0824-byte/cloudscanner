# O — OpenSearch — Sub-50ms Faceted Tool/Deal Search, Natural Language Query Parser & Index Lifecycle

**Document ID:** `CS-AZ-O` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (TirpScanner pivot reframe, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Pivot Note (2026-11-26):** Search reframed from server SKUs to SaaS/AI tools, plans, and deals. Stack unchanged (OpenSearch/BM25/facets). See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0301 ~ WO-CS-0350`

---

## 1. 기능 목적 (Functional Purpose)

Power TirpScanner's instant global SaaS/AI-tool search, natural language plan/deal query parsing, multi-attribute range/term faceting, autocomplete, and zero-downtime index alias swapping.

---

## 2. User Roles & Actors

Search Agent, Backend Engineer, End User.

---

## 3. User Scenarios

1. User types `ai video editor under $20 annual with free trial` into the global search bar; the Natural Language Query Parser extracts structured filters (`category = AI_VIDEO`, `effective_monthly_usd <= 20`, `billing = ANNUAL`, `trial = true`) and executes a hybrid OpenSearch query in 18ms.
2. Hourly deal sync updates 80 plans; the worker bulk-indexes changes into `tirpscanner_skus_v2` via `_bulk` API in <2 seconds without dropping active queries.

---

## 4. Technology Stack

OpenSearch 2.17, Go `opensearch-go/v4` client, deterministic NL Query Parser (`internal/domain/search/nl_parser.go`), PostgreSQL `search_events` analytics sink.

---

## 3. Database Schema & Tables

**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`product_skus`, `price_snapshots`, `search_events`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)

**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**

- `GET /api/v1/search/skus`
- `GET /api/v1/search/autocomplete`
- `POST /api/v1/search/parse-query`
- `POST /api/v1/admin/search/reindex`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. User Interface & Experience

Global Command-K Search Bar (`Cmd+K` / `/search`) with instant parsed-filter chips (e.g., `[Category: AI Video] [Seats: 3] [Billing: Annual] [Max: $20/mo]`) and live facet counts on every sidebar filter.

---

## 8. Core Processing Logic

1. Index Alias Pattern: `tirpscanner_skus_read` points to the active timestamped index (`tirpscanner_skus_2026112`). Full reindexing builds a new index in the background and atomically flips the alias.
2. Deterministic Natural Language Query Parser: Uses fast regex/lexer token rules (zero LLM latency in the hot search path) to parse categories, seat counts, billing intervals (`monthly`/`annual`/`lifetime`), price ceilings, trial/freemium flags, coupon presence, platforms, and renewal-cap hints.

---

## 9. Security & Access Control

Read-only public search query builder sanitizes all filter inputs; OpenSearch cluster is isolated in private data subnet with no public endpoint.

---

## 10. Failure Modes & Resilience

Automatic circuit-breaker fallback to PostgreSQL composite indexes if OpenSearch cluster latency exceeds 500ms or returns 5xx.

---

## 11. Inter-Module Dependencies

Consumes Catalog (`C`) and Deal/Estimation (`E`) data via Outbox relay (`J`).

---

## 12. Testing & Verification Strategy

100+ natural language query parser unit test cases, zero-downtime alias swap integration test, and k6 search load test (p95 < 50ms at 500 RPS).

---

## 13. Acceptance Criteria

- [x] Sub-50ms p95 OpenSearch query execution across 50,000+ tool/plan/deal documents.
- [x] Natural language deal queries accurately decompose into structured filters with visible UI filter chips.
- [x] Zero search downtime during full catalog reindex.

---

## 14. Required Verification Evidence

OpenSearch index mapping JSON, NL parser test suite results, and k6 latency histogram.

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0301 ~ WO-CS-0350`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
