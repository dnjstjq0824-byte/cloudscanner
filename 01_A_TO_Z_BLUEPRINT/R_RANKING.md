# R — Ranking & Recommendation — Transparent Deal Score, Use-Case Presets & Editorial Independence

**Document ID:** `CS-AZ-R` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (TirpScanner pivot reframe, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Pivot Note (2026-11-26):** Server value-score reframed as SaaS deal score; anti-bias law unchanged and still central. See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0331 ~ WO-CS-0350`

---

## 1. 기능 목적 (Functional Purpose)

Define a transparent, mathematically reproducible Ranking & Recommendation Engine that computes objective Deal Scores (`0–100`) per use-case profile while strictly separating paid sponsorships or affiliate commission rates from organic rankings.

---

## 2. User Roles & Actors

Pricing Engine Agent, Data Scientist, Search User.

---

## 3. User Scenarios

1. User views `/best-value-rankings` and selects the `AI Video Editing for Small Teams` preset; plans are ranked by weighted Effective Price-per-Seat (`40%`), Feature Coverage (`23%`), Trial/Freemium Generosity (`13%`), Renewal Fairness (`10%`), and Vendor Reliability (`10%`).
2. A plan with $0 affiliate commission outranks a high-commission plan whenever its objective features-per-dollar score is higher.

---

## 4. Technology Stack

Go `internal/domain/ranking`, PostgreSQL precomputed score snapshots synced to OpenSearch sort fields.

---

## 3. Database Schema & Tables

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

## 7. User Interface & Experience

Best Value Rankings Leaderboard (`/best-value-rankings`) with interactive weight sliders so users can customize how much they care about price vs. features vs. trial vs. renewal fairness, plus a 'How Deal Score is Calculated' transparency modal.

---

## 8. Core Processing Logic

1. Anti-Bias Law: `affiliate_commission_rate` is NEVER an input variable to `ComputeDealScore()`. Organic ranking and `DealScore` depend solely on plan specs, verified deals, and normalized effective cost.
2. Composite Score Formula:
   `DealScore = w_price * Norm(FeatureScore / EffectiveCost) + w_trial * TrialScore + w_renew * RenewalFairness + w_deal * PromoDepthVsMedian + w_rel * Vendor_Reliability_Score`.

---

## 9. Security & Access Control

Prevents manipulation of rankings via review bombing (only verified users with email verification + anti-sybil checks count toward weighted community scores).

---

## 10. Failure Modes & Resilience

When verified deal data is unavailable for a newly launched plan, computes a `Spec-Derived Estimated Score` clearly badged as `Est. from Plan Specs` rather than `Verified Deal`.

---

## 11. Inter-Module Dependencies

Depends on Catalog (`C`) and Deal/Estimation Engine (`E`).

---

## 12. Testing & Verification Strategy

Unit tests proving changing a vendor's affiliate commission from 0% to 80% produces a `0.00` change in `DealScore` and organic rank order.

---

## 13. Acceptance Criteria

- [x] Mathematical independence between affiliate commission rates and organic `DealScore` proven by unit test.
- [x] Support for 6 distinct use-case recommendation presets (AI Writing, AI Video & Image, Project Management, Marketing & CRM, Developer Tools, Team Communication).

---

## 14. Required Verification Evidence

Ranking independence unit test logs and `/api/v1/rankings/best-value` JSON payloads.

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0331 ~ WO-CS-0350`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
