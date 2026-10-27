# C — Catalog — Multi-Category Cloud Product Taxonomy, Versioning & Regional Availability

**Document ID:** `CS-AZ-C`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0151 ~ WO-CS-0200`

---

## 1. 기능 목적 (Functional Purpose)
Architect a unified, extensible product catalog capable of modeling all 16 cloud infrastructure categories (VPS, Shared CPU VM, Dedicated CPU VM, GPU Server, Bare Metal, Object Storage, Block Storage, Managed DB, Kubernetes, Load Balancer, CDN, Backup, Managed Hosting, WordPress Hosting, Developer Hosting, AI Infrastructure) across 18+ providers and 250+ global datacenter regions.

---

## 2. 사용자 역할 (User Roles & Actors)
Catalog Architect, Provider API Agent, Admin Catalog Curator, Search User.

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. DigitalOcean `s-2vcpu-4gb` Droplet is ingested with availability across `nyc1`, `nyc3`, `sfo3`, `ams3`, `sgp1`, `fra1`, `syd1` and linked to family `Basic Regular Droplet`.
2. Vultr `vcg-a100-1c-6g-20vram` (Fractional/Full GPU instance) is cataloged with GPU model `NVIDIA A100`, `20 GB VRAM`, and regional stock status.
3. Hetzner updates `CX23` specs or introduces a new ARM64 `CAX11` generation; the catalog creates a new immutable SKU version while preserving historical price snapshots.

---

## 4. 기술 스택 (Technology Stack)
PostgreSQL 16 (Relational core + JSONB category-specific extension attributes validated by JSON Schema), Go catalog domain service, Redis catalog entity cache.

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`cloud_providers`, `provider_regions`, `provider_products`, `product_skus`, `product_specs`

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

## 7. UI/UX 설계 (User Interface & Experience)
Provider Directory page (`/providers`), Provider Detail page (`/providers/[slug]`), Product Detail page (`/products/[sku]`), and Category Hub pages (`/vps`, `/gpu`, `/bare-metal`, `/managed-hosting`).

---

## 8. 상세 처리 로직 (Core Processing Logic)
1. Three-Tier Product Hierarchy: `cloud_providers` (e.g., DigitalOcean) -> `provider_products` (e.g., Premium AMD Droplets, Managed PostgreSQL) -> `product_skus` (e.g., `s-4vcpu-8gb-amd` in region `fra1`).
2. Strongly Typed Hardware Specs (`product_specs`): Stores `vcpu_count`, `cpu_architecture` (`x86_64`, `arm64`), `cpu_allocation` (`SHARED`, `BURSTABLE`, `DEDICATED`), `cpu_model`, `ram_mb`, `storage_type` (`NVME`, `SSD`, `HDD`, `NETWORK_BLOCK`), `storage_gb`, `raid_config`, `bandwidth_tb`, `network_gbps`, `gpu_model`, `gpu_count`, `gpu_vram_gb`, `ipv4_included`, `ipv6_included`, plus category-specific `extra_specs_jsonb` (e.g., K8s control plane fee, Object Storage $/GB/mo).

---

## 9. 보안 통제 (Security & Access Control)
Public read access to active catalog items; write mutations restricted to signed internal `PricingAdapter` sync workers or `CATALOG_ADMIN` role with mandatory audit logging.

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
When a provider deprecates a SKU, `product_skus.status` transitions to `DEPRECATED` (never hard-deleted) so existing user favorites, comparison lists, and historical price charts remain intact with a 'Legacy Plan — No Longer Available for New Deployment' badge.

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Depends on PostgreSQL (`P_POSTGRESQL.md`); feeds Data Acquisition (`D_DATA_ACQUISITION.md`), Estimation Engine (`E_ESTIMATION_ENGINE.md`), and OpenSearch (`O_OPENSEARCH.md`).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
Schema validation tests across all 16 product categories, region-mapping normalization tests (mapping `fra1`, `eu-central-1`, `fsn1`, `de-fra` to canonical country `DE` and continent `EU`), and SKU lifecycle deprecation tests.

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] All 16 product categories defined in Article 2 are supported without requiring schema-breaking DDL migrations.
- [x] Every SKU accurately records CPU architecture (`x86_64` vs `arm64`), allocation type (`SHARED` vs `DEDICATED`), storage medium, bandwidth allowance, and GPU VRAM.
- [x] Canonical region mapping normalizes 250+ provider-specific datacenter codes into ISO-3166 country codes and geographic coordinates.

---

## 14. 검증 증빙 (Required Verification Evidence)
Seed catalog database dump, SQL query output across all 16 categories, and OpenAPI schema validation passes.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0151 ~ WO-CS-0200`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
