# N — Network & Region — Global Datacenter Taxonomy, Egress Pricing Rules, IPv4/IPv6 & Latency

**Document ID:** `CS-AZ-N`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0170 ~ WO-CS-0185, WO-CS-0260 ~ WO-CS-0270`

---

## 1. 기능 목적 (Functional Purpose)
Standardize global cloud datacenter regions, inter-region/internet bandwidth egress pricing tiers, pooled vs. per-instance transfer allowances, IPv4 exhaustion surcharges, and IPv6 readiness across all providers.

---

## 2. 사용자 역할 (User Roles & Actors)
Network & Pricing Architect, Search User.

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. User filters for servers in `Asia-Pacific -> Northeast Asia (Seoul, Tokyo, Osaka)` with at least `4 TB Included Transfer` and `IPv4 Included`.
2. User compares bandwidth overage for a 15 TB/month streaming server: Hetzner EU (20 TB included -> $0 overage) vs. DigitalOcean ($0.01/GB overage) vs. AWS EC2 ($0.09/GB overage -> $1,341/mo egress bill).

---

## 4. 기술 스택 (Technology Stack)
Go network cost normalizer, PostgreSQL `provider_regions` with PostGIS/lat-lon coordinates and ISO-3166 country/city subdivisions.

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`provider_regions`, `product_specs`, `price_components`, `billing_terms`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)
**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**
- `GET /api/v1/regions`
- `GET /api/v1/regions/comparison`
- `POST /api/v1/pricing/estimate-network`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. UI/UX 설계 (User Interface & Experience)
Interactive Global Region Map & Region Price Comparison Matrix (`/region-comparison`) showing cheapest VPS/GPU by continent, country, city, green energy status, and egress cost per TB.

---

## 8. 상세 처리 로직 (Core Processing Logic)
Models 5 distinct bandwidth billing structures: (1) Outbound-only metered with free inbound (DO, Vultr, AWS, GCP, Azure), (2) Bidirectional metered, (3) Unmetered fixed port speed (OVH/Contabo fair-use tiers), (4) Account-wide pooled transfer (DigitalOcean/Linode/UpCloud pool across active VMs), and (5) Free egress up to high threshold (UpCloud zero-cost egress policy / OCI 10 TB free).

---

## 9. 보안 통제 (Security & Access Control)
Validates region availability per user account during BYOC preflight (some cloud accounts have regional locks until support raises quota).

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
Distinguishes between `Region Cataloged` and `Region Currently Out of Stock (Capacity Full)` when provider API returns stock availability flags.

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Feeds Catalog (`C`), Estimation Engine (`E`), and OpenSearch (`O`).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
Unit tests for pooled bandwidth calculation, AWS/GCP tiered egress brackets, and IPv4 vs IPv6-only pricing differentials.

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] Every provider region maps to canonical Continent, Country ISO-2, City, and Lat/Lon.
- [x] IPv4 hourly/monthly add-on costs ($0.005/hr on AWS/GCP/Azure, €0.50–€0.60/mo on Hetzner) are accurately factored into TCO.

---

## 14. 검증 증빙 (Required Verification Evidence)
Region taxonomy seed table verification and network egress unit test results.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0170 ~ WO-CS-0185, WO-CS-0260 ~ WO-CS-0270`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
