# V — VPS, GPU & Infra Expansion — RETIRED (Hardware Scope Removed by 2026-11-26 Pivot)

**Document ID:** `CS-AZ-V` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (marked RETIRED, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Status:** **RETIRED** — VPS/GPU/bare-metal/storage/managed-DB comparison (GPU unit economics, `/gpu-comparison`, hardware spec columns, 12 infra categories) belonged to the retired cloud-infrastructure catalog. SaaS/AI-tool category expansion continues under Catalog (`C`). See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0186 ~ WO-CS-0200, WO-CS-0480 ~ WO-CS-0510`

---

## 1. 기능 목적 (Functional Purpose)

HISTORICAL (CloudScanr era): Extend the comparison and management engines beyond basic VPS to support high-growth AI GPU servers (NVIDIA H100, H200, A100, L40S, RTX 4090, AMD MI300X), Bare Metal servers, Managed Databases, Object/Block Storage, and Managed WordPress/App Hosting. This purpose is retired; the sections below preserve the original spec for traceability.

---

## 2. User Roles & Actors

HISTORICAL: Catalog Architect, AI Infrastructure Buyer, Pricing Engine Agent. No live actors.

---

## 3. User Scenarios

HISTORICAL (retired, preserved in intent):

1. AI engineer searches for `80GB VRAM GPU` across Vultr, DigitalOcean GPU Droplets, Scaleway, OVHcloud, GCP, AWS, and RunPod, sorting by `$/VRAM-hour`.
2. Enterprise architect compares 10 TB Object Storage with 30 TB/month public egress across Cloudflare R2, Hetzner Object Storage, DigitalOcean Spaces, Vultr Objects, and AWS S3.

---

## 4. Technology Stack

HISTORICAL: PostgreSQL `product_specs` specialized hardware columns (`gpu_model`, `gpu_count`, `gpu_vram_gb`, `gpu_interconnect`, `fp12_tflops`) + category-specific JSON Schema validators. Hardware columns retired; `product_specs` now rows SaaS plan specs (see `C_CATALOG.md`).

---

## 5. Database Schema & Tables

HISTORICAL hardware use of `provider_products`, `product_skus`, `product_specs`, `price_components` is retired. Tables remain live for the SaaS catalog (see `P_POSTGRESQL.md`).

---

## 6. API 구조 (API Contracts & Endpoints)

RETIRED endpoints (do not implement; remove from live OpenAPI surface):

- `GET /api/v1/catalog/gpu`
- `GET /api/v1/catalog/bare-metal`
- `GET /api/v1/catalog/storage`
- `GET /api/v1/catalog/managed-db`

---

## 7. User Interface & Experience

RETIRED page: GPU Comparison Matrix (`/gpu-comparison`). Do not build.

---

## 8. Core Processing Logic

HISTORICAL (retired): GPU/storage/managed-DB unit economics (`PricePerGPUHour`, `PricePerTBMonthStored`, `HAStandbyMultiplier`, etc.).

---

## 9. Security & Access Control

RETIRED: GPU BYOC budget-cap rules retired with `M_MULTICLOUD_MANAGEMENT.md`.

---

## 10. Failure Modes & Resilience

HISTORICAL: GPU stock/capacity indicator (`AVAILABLE`, `LIMITED_CAPACITY`, `WAITLIST_ONLY`). Retired.

---

## 11. Inter-Module Dependencies

Formerly extended Catalog (`C_CATALOG.md`) and Estimation Engine (`E_ESTIMATION_ENGINE.md`). No live dependents; nothing may newly depend on this module.

---

## 12. Testing & Verification Strategy

HISTORICAL: hardware category normalization and unit-economics tests. Retired; preserved as frozen work-order evidence.

---

## 13. Acceptance Criteria

- [x] Module formally marked RETIRED with live/retired scope recorded in `TIRPSCANNER_PIVOT.md`.
- [x] No live OpenAPI group, page, or job references hardware/infra comparison.
- [x] Historical spec preserved above for audit traceability.

---

## 14. Required Verification Evidence

Pivot decision record (`TIRPSCANNER_PIVOT.md`); frozen work-order history (`11_WORK_ORDERS_0001_TO_1000/`, read-only).

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0186 ~ WO-CS-0200, WO-CS-0480 ~ WO-CS-0510`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
