# N — Network & Region — RETIRED (Datacenter/Egress Scope Removed by 2026-11-26 Pivot)

**Document ID:** `CS-AZ-N` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (marked RETIRED, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Status:** **RETIRED** — Datacenter taxonomy, egress pricing tiers, pooled transfer, IPv4/IPv6 surcharges, `/region-comparison` page, and `estimate-network` API belonged to the retired cloud-infrastructure catalog. SaaS availability locales and currencies continue under Catalog (`C`) and Globalization (`G`). See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0170 ~ WO-CS-0185, WO-CS-0260 ~ WO-CS-0270`

---

## 1. 기능 목적 (Functional Purpose)

HISTORICAL (CloudScanr era): Standardize global cloud datacenter regions, inter-region/internet bandwidth egress pricing tiers, pooled vs. per-instance transfer allowances, IPv4 exhaustion surcharges, and IPv6 readiness across all providers. This purpose is retired; the sections below preserve the original spec for traceability.

---

## 2. User Roles & Actors

HISTORICAL: Network & Pricing Architect, Search User. No live actors.

---

## 3. User Scenarios

HISTORICAL (retired, preserved in intent):

1. User filters for servers in `Asia-Pacific -> Northeast Asia (Seoul, Tokyo, Osaka)` with at least `4 TB Included Transfer` and `IPv4 Included`.
2. User compares bandwidth overage for a 15 TB/month streaming server across Hetzner EU, DigitalOcean, and AWS EC2 egress rates.

---

## 4. Technology Stack

HISTORICAL: Go network cost normalizer, PostgreSQL `provider_regions` with PostGIS/lat-lon coordinates. The `provider_regions` table itself remains live for SaaS availability locales (see `C_CATALOG.md`); datacenter/egress modeling is retired.

---

## 5. Database Schema & Tables

HISTORICAL use of `provider_regions`, `product_specs`, `price_components`, `billing_terms` for network modeling is retired. Tables remain live for SaaS catalog purposes (see `P_POSTGRESQL.md`).

---

## 6. API 구조 (API Contracts & Endpoints)

RETIRED endpoints (do not implement; remove from live OpenAPI surface):

- `GET /api/v1/regions`
- `GET /api/v1/regions/comparison`
- `POST /api/v1/pricing/estimate-network`

---

## 7. User Interface & Experience

RETIRED page: Global Region Map & Region Price Comparison Matrix (`/region-comparison`). Do not build.

---

## 8. Core Processing Logic

HISTORICAL (retired): 5 bandwidth billing structures (outbound-only metered, bidirectional metered, unmetered fair-use, account-pooled transfer, high free-egress thresholds).

---

## 9. Security & Access Control

RETIRED: BYOC preflight region-availability validation retired with `M_MULTICLOUD_MANAGEMENT.md`.

---

## 10. Failure Modes & Resilience

HISTORICAL: `Region Cataloged` vs `Region Out of Stock` distinction. Retired.

---

## 11. Inter-Module Dependencies

Formerly fed Catalog (`C`), Estimation Engine (`E`), and OpenSearch (`O`). No live dependents; availability-locale needs are served by `C_CATALOG.md` and `G_GLOBALIZATION.md`.

---

## 12. Testing & Verification Strategy

HISTORICAL: pooled bandwidth, tiered egress, IPv4/IPv6 differential tests. Retired; preserved as frozen work-order evidence.

---

## 13. Acceptance Criteria

- [x] Module formally marked RETIRED with live/retired scope recorded in `TIRPSCANNER_PIVOT.md`.
- [x] No live OpenAPI group, page, or job references datacenter/egress modeling.
- [x] Historical spec preserved above for audit traceability.

---

## 14. Required Verification Evidence

Pivot decision record (`TIRPSCANNER_PIVOT.md`); frozen work-order history (`11_WORK_ORDERS_0001_TO_1000/`, read-only).

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0170 ~ WO-CS-0185, WO-CS-0260 ~ WO-CS-0270`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
