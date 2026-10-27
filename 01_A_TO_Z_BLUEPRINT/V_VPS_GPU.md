# V — VPS, GPU & Product Expansion — Specialized Hardware Specs, AI GPU Clusters, Bare Metal & Storage

**Document ID:** `CS-AZ-V`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0186 ~ WO-CS-0200, WO-CS-0480 ~ WO-CS-0510`

---

## 1. 기능 목적 (Functional Purpose)
Extend the comparison and management engines beyond basic VPS to support high-growth AI GPU servers (NVIDIA H100, H200, A100, L40S, RTX 4090, AMD MI300X), Bare Metal servers, Managed Databases, Object/Block Storage, and Managed WordPress/App Hosting.

---

## 2. 사용자 역할 (User Roles & Actors)
Catalog Architect, AI Infrastructure Buyer, Pricing Engine Agent.

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. AI engineer searches for `80GB VRAM GPU` (`NVIDIA A100 80GB` or `H100 80GB`) with `NVLink` or `PCIe` interconnect across Vultr, DigitalOcean GPU Droplets, Scaleway, OVHcloud, GCP, AWS, and RunPod, sorting by `$/VRAM-hour`.
2. Enterprise architect compares 100 TB Object Storage with 50 TB/month public egress across Cloudflare R2 ($0 egress), Hetzner Object Storage, DigitalOcean Spaces, Vultr Objects, and AWS S3.

---

## 4. 기술 스택 (Technology Stack)
PostgreSQL `product_specs` specialized columns (`gpu_model`, `gpu_count`, `gpu_vram_gb`, `gpu_interconnect`, `fp16_tflops`) + category-specific JSON Schema validators in Go.

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`provider_products`, `product_skus`, `product_specs`, `price_components`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)
**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**
- `GET /api/v1/catalog/gpu`
- `GET /api/v1/catalog/bare-metal`
- `GET /api/v1/catalog/storage`
- `GET /api/v1/catalog/managed-db`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. UI/UX 설계 (User Interface & Experience)
Dedicated GPU Comparison Matrix (`/gpu-comparison`) featuring `$/GPU-Hour`, `$/GB VRAM-Hour`, Interconnect (`SXM5`/`NVLink` vs `PCIe`), CUDA/Driver pre-install badges, and Spot/Preemptible vs On-Demand toggles.

---

## 8. 상세 처리 로직 (Core Processing Logic)
Computes specialized unit economics per category:
- **GPU**: `PricePerGPUHour`, `PricePerGBVRAMMonth`.
- **Storage**: `PricePerTBMonthStored` + `PricePer10kReadWriteOps` + `PricePerTBEgress`.
- **Managed DB**: `BaseNodePrice` + `HAStandbyMultiplier` + `BackupStorageOverage`.

---

## 9. 보안 통제 (Security & Access Control)
GPU provisioning in BYOC enforces stricter default budget caps (e.g., warning threshold at $2.00/hr or $500/mo) due to rapid hourly burn rates of multi-GPU instances.

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
Real-time GPU stock/capacity indicator (`AVAILABLE`, `LIMITED_CAPACITY`, `WAITLIST_ONLY`) since GPU SKUs frequently experience regional stockouts.

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Extends Catalog (`C_CATALOG.md`) and Estimation Engine (`E_ESTIMATION_ENGINE.md`).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
Category-specific normalization and unit-economics calculation tests for GPU, Bare Metal, Storage, and Managed DB.

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] Dedicated `/gpu-comparison` specification supports filtering by GPU model, total VRAM, interconnect type, and Spot vs On-Demand billing.
- [x] All 16 product categories in Article 2 have validated JSON schemas and comparison metrics.

---

## 14. 검증 증빙 (Required Verification Evidence)
GPU & multi-category schema test logs and sample GPU comparison payloads.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0186 ~ WO-CS-0200, WO-CS-0480 ~ WO-CS-0510`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
