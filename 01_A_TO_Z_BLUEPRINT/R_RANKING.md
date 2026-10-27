# R — Ranking & Recommendation — Transparent Value Score, Workload Presets & Editorial Independence

**Document ID:** `CS-AZ-R`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0331 ~ WO-CS-0350`

---

## 1. 기능 목적 (Functional Purpose)
Define a transparent, mathematically reproducible Ranking & Recommendation Engine that computes objective Value Scores (`0–100`) per workload profile while strictly separating paid sponsorships or affiliate commission rates from organic rankings.

---

## 2. 사용자 역할 (User Roles & Actors)
Pricing Engine Agent, Data Scientist, Search User.

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. User views `/best-value-rankings` and selects the `Low-Latency Web & API Server` preset; SKUs are ranked by weighted Price-Performance (`40%`), Network/Egress Generosity (`25%`), Disk IOPS/NVMe (`20%`), and Provider Reliability/SLA (`15%`).
2. A provider with $0 affiliate commission (e.g., Hetzner or Oracle Cloud) outranks a high-commission provider (e.g., Cloudways or Hostinger) whenever its objective specs-per-dollar score is higher.

---

## 4. 기술 스택 (Technology Stack)
Go `internal/domain/ranking`, PostgreSQL precomputed score snapshots synced to OpenSearch sort fields.

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`product_skus`, `product_specs`, `price_snapshots`, `provider_reviews`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)
**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**
- `GET /api/v1/rankings/best-value`
- `POST /api/v1/rankings/recommend`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. UI/UX 설계 (User Interface & Experience)
Best Value Rankings Leaderboard (`/best-value-rankings`) with interactive weight sliders so users can customize how much they care about CPU vs. RAM vs. NVMe vs. Bandwidth, plus a 'How Value Score is Calculated' transparency modal.

---

## 8. 상세 처리 로직 (Core Processing Logic)
1. Anti-Bias Law (Article 7-R & 17): `affiliate_commission_rate` is NEVER an input variable to `ComputeValueScore()`. Organic ranking and `ValueScore` depend solely on technical specifications, verified benchmarks, and normalized TCO.
2. Composite Score Formula:
   `ValueScore = w_cpu * Norm(CPU_Score / TCO) + w_ram * Norm(RAM_GB / TCO) + w_disk * Norm(Storage_Score / TCO) + w_net * Norm(Bandwidth_Score / TCO) + w_rel * Provider_SLA_Score`.

---

## 9. 보안 통제 (Security & Access Control)
Prevents manipulation of rankings via review bombing (only verified users with email verification + anti-sybil check or verified BYOC connection badge count toward weighted community scores).

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
When benchmark trial data is unavailable for a newly launched SKU, computes a `Spec-Derived Estimated Score` clearly badged as `Est. from Hardware Specs` rather than `Verified Benchmark`.

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Depends on Catalog (`C`) and Estimation Engine (`E`).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
Unit tests proving changing a provider's affiliate commission from 0% to 80% produces a `0.0000` change in `ValueScore` and organic rank order.

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] Mathematical independence between affiliate commission rates and organic `ValueScore` proven by unit test.
- [x] Support for 6 distinct workload recommendation presets (General Web, High-Memory DB, CPU-Intensive Build, High-Egress Media, K8s Worker, AI/GPU Inference).

---

## 14. 검증 증빙 (Required Verification Evidence)
Ranking independence unit test logs and `/api/v1/rankings/best-value` JSON payloads.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0331 ~ WO-CS-0350`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
