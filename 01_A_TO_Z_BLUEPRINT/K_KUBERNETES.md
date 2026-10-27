# K — Kubernetes & Containers — Managed K8s Cluster Cost Comparison & Node Pool TCO

**Document ID:** `CS-AZ-K`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0281 ~ WO-CS-0295, WO-CS-0520 ~ WO-CS-0530`

---

## 1. 기능 목적 (Functional Purpose)
Model and compare Managed Kubernetes services (DigitalOcean DOKS, UpCloud UKS, Vultr VKE, Linode LKE, OVHcloud MKS, Scaleway Kapsule, AWS EKS, GCP GKE, Azure AKS) including hidden control-plane fees, HA control-plane surcharges, worker node pool pricing, load balancer attachments, and persistent volume CSI costs.

---

## 2. 사용자 역할 (User Roles & Actors)
Cloud Architect, DevOps Buyer, Pricing Engine Agent.

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. User configures a 3-node Kubernetes cluster (each node 4 vCPU / 16 GB RAM) + 1 HA Control Plane + 1 Cloud Load Balancer + 200 GB Block Storage.
2. CloudScanr calculates that AWS EKS charges $0.10/hr ($73/mo) for the control plane whereas DigitalOcean DOKS, Vultr VKE, Linode LKE, and OVHcloud MKS provide a free standard control plane, shifting total 3-node cluster TCO by 35–60%.

---

## 4. 기술 스택 (Technology Stack)
Go K8s TCO calculator module (`internal/domain/estimation/kubernetes.go`), PostgreSQL `product_specs.extra_specs_jsonb` K8s schema, Nuxt 3 K8s Cluster Builder UI.

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`provider_products`, `product_skus`, `product_specs`, `price_components`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)
**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**
- `POST /api/v1/pricing/estimate-kubernetes`
- `GET /api/v1/catalog/kubernetes`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. UI/UX 설계 (User Interface & Experience)
Kubernetes Cluster Cost Comparison page (`/kubernetes-comparison`) with sliders for Control Plane Tier (Standard vs HA SLA), Worker Node Count & SKU spec, Load Balancer count, Persistent Volume GiB, and Container Registry storage.

---

## 8. 상세 처리 로직 (Core Processing Logic)
`ClusterMonthlyTCO = ControlPlaneFee(Tier) + SUM(NodePool_i.Count * NodeSKU_i.MonthlyPrice) + (LoadBalancerCount * LBMonthlyPrice) + (PersistentStorageGiB * BlockStorageRatePerGiB) + EgressOverageCost`.

---

## 9. 보안 통제 (Security & Access Control)
Read-only cluster discovery supported in BYOC for DOKS/LKE/VKE/EKS/GKE when user grants read-only K8s scope; cluster deletion disabled by default unless multi-party org approval is enabled.

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
Clearly flags region-specific K8s version availability and node-pool minimum size constraints (e.g., minimum 1 vs 3 nodes).

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Extends Catalog (`C_CATALOG.md`) and Estimation Engine (`E_ESTIMATION_ENGINE.md`).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
Unit tests comparing 3-node and 10-node cluster TCO across DOKS, LKE, VKE, UKS, EKS, GKE, and AKS.

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] Accurately models both Free Control Plane providers (DOKS, VKE, LKE, MKS) and Paid/Tiered Control Plane providers (EKS $73/mo, GKE/AKS free tier vs SLA tier, UpCloud UKS).
- [x] Includes mandatory Load Balancer and CSI Block Storage costs in cluster TCO breakdown.

---

## 14. 검증 증빙 (Required Verification Evidence)
Kubernetes TCO unit test outputs and `/api/v1/pricing/estimate-kubernetes` response payloads.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0281 ~ WO-CS-0295, WO-CS-0520 ~ WO-CS-0530`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
