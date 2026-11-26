# C — Catalog — SaaS/AI-Tool Product Taxonomy, Plan Specs, Deal Shelf & Availability

**Document ID:** `CS-AZ-C` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (TirpScanner pivot reframe, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Pivot Note (2026-11-26):** Cloud-infrastructure catalog (VPS/GPU/bare-metal/TCO, hardware specs, datacenter regions) is RETIRED. Live catalog covers SaaS/AI tools, plans, deals. See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0151 ~ WO-CS-0200`

---

## 1. 기능 목적 (Functional Purpose)

Architect a unified, extensible product catalog capable of modeling SaaS/AI-tool and software products across vendors: plan tiers (Free/Trial/Starter/Pro/Business/Enterprise), per-seat and usage-based pricing, feature specs, integrations, platform availability, and linked deals/coupons — with versioned plan history and availability tracking.

---

## 2. User Roles & Actors

Catalog Architect, Vendor Feed Agent, Admin Catalog Curator, Search User.

---

## 3. User Scenarios

1. A project-management SaaS product is ingested with plans `Free`, `Starter ($10/seat/mo annual)`, `Pro ($20/seat/mo annual)`, trial terms, and platform availability (Web, iOS, Android).
2. An AI video generator is cataloged with its subscription tiers plus credit-pack add-ons (`10 credits/mo included`, `$13 per 10 extra credits`) and linked active coupons.
3. A vendor re-tiers plans or renames a tier; the catalog creates a new immutable plan version while preserving historical price snapshots and deal history.

---

## 4. Technology Stack

PostgreSQL 16 (Relational core + JSONB category-specific extension attributes validated by JSON Schema), Go catalog domain service, Redis catalog entity cache.

---

## 5. Database Schema & Tables

**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`cloud_providers`, `provider_regions`, `provider_products`, `product_skus`, `product_specs`

- Table names are historical (CloudScanr era): `cloud_providers` now rows SaaS vendors/networks, `product_skus` rows vendor plans, `product_specs` rows plan specs, `provider_regions` rows availability locales.
- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)

**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**

- `GET /api/v1/providers`
- `GET /api/v1/providers/{slug}`
- `GET /api/v1/catalog/categories`
- `GET /api/v1/catalog/skus`
- `GET /api/v1/catalog/skus/{sku_id}`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. User Interface & Experience

Vendor Directory page (`/providers`), Vendor Detail page (`/providers/[slug]`), Tool Detail page (`/products/[sku]`), and Category Hub pages (`/ai-writing`, `/video-editing`, `/project-management`, `/deals`).

---

## 8. Core Processing Logic

1. Three-Tier Product Hierarchy: `cloud_providers` (e.g., a SaaS vendor — table name historical) -> `provider_products` (e.g., AI Video Suite, Team Chat) -> `product_skus` (e.g., `pro-annual-3seat` plan variant).
2. Strongly Typed Plan Specs (`product_specs`): Stores `billing_interval` (`MONTHLY`, `ANNUAL`, `LIFETIME`, `CREDIT_PACK`), `price_per_seat`, `min_seats`, `trial_days`, `freemium_limits`, `usage_quotas` (credits/tokens/storage), `renewal_terms`, `platforms`, `integrations`, plus category-specific `extra_specs_jsonb` (e.g., AI model access tier, export resolution caps, SSO availability).

---

## 9. Security & Access Control

Public read access to active catalog items; write mutations restricted to signed internal catalog/deal-ingestion sync workers or `CATALOG_ADMIN` role with mandatory audit logging.

---

## 10. Failure Modes & Resilience

When a vendor retires a plan, `product_skus.status` transitions to `DEPRECATED` (never hard-deleted) so existing user favorites, comparison lists, and historical price charts remain intact with a 'Legacy Plan — No Longer Available' badge.

---

## 11. Inter-Module Dependencies

Depends on PostgreSQL (`P_POSTGRESQL.md`); feeds Data Acquisition (`D_DATA_ACQUISITION.md`), Deal/Estimation Engine (`E_ESTIMATION_ENGINE.md`), and OpenSearch (`O_OPENSEARCH.md`).

---

## 12. Testing & Verification Strategy

Schema validation tests across all SaaS/AI-tool categories, availability-locale normalization tests (mapping vendor-stated availability to canonical country/currency sets), and plan lifecycle deprecation tests.

---

## 13. Acceptance Criteria

- [x] All live SaaS/AI-tool categories are supported without requiring schema-breaking DDL migrations.
- [x] Every plan accurately records billing interval, per-seat price, trial/freemium terms, usage quotas, and renewal terms.
- [x] Canonical availability mapping normalizes vendor-stated locales into ISO-3166 country codes and currencies.

---

## 14. Required Verification Evidence

Seed catalog database dump, SQL query output across all live categories, and OpenAPI schema validation passes.

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0151 ~ WO-CS-0200`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
