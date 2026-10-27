# O — OpenSearch — Sub-50ms Faceted Product Search, Natural Language Spec Parser & Index Lifecycle

**Document ID:** `CS-AZ-O`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0301 ~ WO-CS-0350`

---

## 1. 기능 목적 (Functional Purpose)
Power CloudScanr's instant global cloud search, natural language specification query parsing, multi-attribute range/term faceting, autocomplete, and zero-downtime index alias swapping.

---

## 2. 사용자 역할 (User Roles & Actors)
Search Agent, Backend Engineer, End User.

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. User types `8 vcpu 32gb ram nvme frankfurt under $60` into the global search bar; the Natural Language Query Parser extracts structured filters (`vcpu >= 8`, `ram_gb >= 32`, `storage_type = NVME`, `city = Frankfurt`, `monthly_price_usd <= 60`) and executes a hybrid OpenSearch query in 18ms.
2. Hourly price sync updates 800 SKUs; the worker bulk-indexes changes into `cloudscanr_skus_v2` via `_bulk` API in <2 seconds without dropping active queries.

---

## 4. 기술 스택 (Technology Stack)
OpenSearch 2.17, Go `opensearch-go/v4` client, deterministic NL Spec Parser (`internal/domain/search/nl_parser.go`), PostgreSQL `search_events` analytics sink.

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
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

## 7. UI/UX 설계 (User Interface & Experience)
Global Command-K Search Bar (`Cmd+K` / `/search`) with instant parsed-filter chips (e.g., `[vCPU: >=8] [RAM: >=32GB] [Region: Frankfurt] [Max: $60/mo]`) and live facet counts on every sidebar filter.

---

## 8. 상세 처리 로직 (Core Processing Logic)
1. Index Alias Pattern: `cloudscanr_skus_read` points to the active timestamped index (`cloudscanr_skus_20261028`). Full reindexing builds a new index in the background and atomically flips the alias.
2. Deterministic Natural Language Spec Parser: Uses fast regex/lexer token rules (zero LLM latency in the hot search path) to parse CPU counts, RAM (`GB`/`MB`), disk sizes (`GB`/`TB`), disk types (`NVMe`, `SSD`), GPU models (`H100`, `A100`, `L40S`, `RTX 4090`), regions/countries, billing intervals, and price ceilings.

---

## 9. 보안 통제 (Security & Access Control)
Read-only public search query builder sanitizes all filter inputs; OpenSearch cluster is isolated in private data subnet with no public endpoint.

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
Automatic circuit-breaker fallback to PostgreSQL composite indexes if OpenSearch cluster latency exceeds 500ms or returns 5xx.

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Consumes Catalog (`C`) and Estimation (`E`) data via Outbox relay (`J`).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
100+ natural language query parser unit test cases, zero-downtime alias swap integration test, and k6 search load test (p95 < 50ms at 500 RPS).

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] Sub-50ms p95 OpenSearch query execution across 50,000+ regional SKU documents.
- [x] Natural language spec queries accurately decompose into structured filters with visible UI filter chips.
- [x] Zero search downtime during full catalog reindex.

---

## 14. 검증 증빙 (Required Verification Evidence)
OpenSearch index mapping JSON, NL parser test suite results, and k6 latency histogram.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0301 ~ WO-CS-0350`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
