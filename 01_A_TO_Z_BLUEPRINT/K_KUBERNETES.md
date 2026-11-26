# K — Kubernetes Cost Comparison — RETIRED (Cloud-Infra Scope Removed by 2026-11-26 Pivot)

**Document ID:** `CS-AZ-K` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (marked RETIRED, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Status:** **RETIRED** — Managed-Kubernetes cost comparison (DOKS/UKS/VKE/LKE/MKS/EKS/GKE/AKS control-plane fees, node-pool TCO, `/kubernetes-comparison` page, `estimate-kubernetes` API) belonged to the retired cloud-infrastructure catalog. No new work. TirpScanner's *own* hosting on Kubernetes remains live and is specified in `H_HOSTING.md`. See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0281 ~ WO-CS-0295, WO-CS-0520 ~ WO-CS-0530`

---

## 1. 기능 목적 (Functional Purpose)

HISTORICAL (CloudScanr era): Model and compare Managed Kubernetes services including hidden control-plane fees, HA control-plane surcharges, worker node pool pricing, load balancer attachments, and persistent volume CSI costs. This purpose is retired; the section below preserves the original spec for traceability.

Original live purpose (retired): compare DOKS, UKS, VKE, LKE, MKS, EKS, GKE, AKS cluster TCO via `internal/domain/estimation/kubernetes.go`, `product_specs.extra_specs_jsonb` K8s schema, and a Nuxt Cluster Builder UI.

---

## 2. User Roles & Actors

HISTORICAL: Cloud Architect, DevOps Buyer, Pricing Engine Agent. No live actors.

---

## 3. User Scenarios

HISTORICAL (retired, preserved verbatim in intent):

1. User configures a 3-node Kubernetes cluster (each node 4 vCPU / 16 GB RAM) + 1 HA Control Plane + 1 Cloud Load Balancer + 200 GB Block Storage.
2. The engine calculates that AWS EKS charges $0.10/hr ($73/mo) for the control plane whereas DOKS, VKE, LKE, and MKS provide a free standard control plane, shifting total 3-node cluster TCO by 35–60%.

---

## 4. Technology Stack

HISTORICAL: Go K8s TCO calculator module, PostgreSQL `extra_specs_jsonb` K8s schema, Nuxt 3 K8s Cluster Builder UI. No live stack.

---

## 3. Database Schema & Tables

HISTORICAL tables (`product_skus`, `product_specs`, `price_components` K8s rows) are retired for this purpose; the tables themselves remain live for the SaaS catalog (see `C_CATALOG.md`, `P_POSTGRESQL.md`).

---

## 6. API 구조 (API Contracts & Endpoints)

RETIRED endpoints (do not implement; remove from live OpenAPI surface):

- `POST /api/v1/pricing/estimate-kubernetes`
- `GET /api/v1/catalog/kubernetes`

---

## 7. User Interface & Experience

RETIRED page: Kubernetes Cluster Cost Comparison (`/kubernetes-comparison`). Do not build.

---

## 8. Core Processing Logic

HISTORICAL formula (retired): `ClusterMonthlyTCO = ControlPlaneFee(Tier) + SUM(NodePool_i.Count * NodeSKU_i.MonthlyPrice) + (LoadBalancerCount * LBMonthlyPrice) + (PersistentStorageGiB * BlockStorageRatePerGiB) + EgressOverageCost`.

---

## 9. Security & Access Control

RETIRED: BYOC K8s cluster discovery/deletion rules retired with `M_MULTICLOUD_MANAGEMENT.md`.

---

## 10. Failure Modes & Resilience

HISTORICAL: region-specific K8s version availability and node-pool minimum size flags. Retired.

---

## 11. Inter-Module Dependencies

Formerly extended Catalog (`C_CATALOG.md`) and Estimation Engine (`E_ESTIMATION_ENGINE.md`). No live dependents; nothing may newly depend on this module.

---

## 12. Testing & Verification Strategy

HISTORICAL: 3-node and 10-node cluster TCO tests across DOKS, LKE, VKE, UKS, EKS, GKE, AKS. Retired; keep only as frozen work-order evidence.

---

## 13. Acceptance Criteria

- [x] Module formally marked RETIRED with live/retired scope recorded in `TIRPSCANNER_PIVOT.md`.
- [x] No live OpenAPI group, page, or job references K8s cost comparison.
- [x] Historical spec preserved above for audit traceability.

---

## 14. Required Verification Evidence

Pivot decision record (`TIRPSCANNER_PIVOT.md`); frozen work-order history (`11_WORK_ORDERS_0001_TO_1000/`, read-only).

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0281 ~ WO-CS-0295, WO-CS-0520 ~ WO-CS-0530`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
