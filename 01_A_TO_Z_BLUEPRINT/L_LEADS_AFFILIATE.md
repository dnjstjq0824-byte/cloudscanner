# L — Leads & Attribution — Affiliate Link Router, SubID/ClickID Signing, Conversion Reconciliation & Anti-Fraud

**Document ID:** `CS-AZ-L`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0351 ~ WO-CS-0450`

---

## 1. 기능 목적 (Functional Purpose)
Engineer the complete Affiliate Attribution and Revenue Tracking Engine mandated by Article 10, covering affiliate link generation, cryptographically signed `ClickID`/`SubID` injection, sub-10ms `/go/{slug}` redirect routing, bot/click-fraud protection, network webhook/CSV conversion reconciliation, and policy compliance.

---

## 2. 사용자 역할 (User Roles & Actors)
Affiliate Engine Agent, Marketing Operations Admin, Security Reviewer.

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. User on `/vps-comparison` clicks 'Deploy on DigitalOcean (Affiliate)'; the browser hits `GET /go/digitalocean-s-2vcpu-4gb?src=compare_table`.
2. The Go Redirect Engine verifies the request is not a known bot (`User-Agent` + Cloudflare Bot Score), generates a ULID `click_id` (`01J...`), signs an HMAC-SHA256 `SubID` token, logs the click asynchronously, and issues a `302 Found` redirect to Impact's tracking URL with `?subId1={click_id}&subId2=compare_table`.
3. When Impact fires a conversion postback or the Admin uploads a Cloudways/Contabo CSV report, the Reconciliation Engine matches `subId1` -> `affiliate_clicks.click_id`, validates the signature, updates `affiliate_conversions`, and posts a `PENDING` entry to `commissions`.

---

## 4. 기술 스택 (Technology Stack)
Go `internal/domain/affiliate`, ULID generator (`oklog/ulid/v2`), HMAC-SHA256 SubID signer, MaxMind GeoIP2 (country/ASN lookup without storing raw PII IP), Redis token-bucket click fraud limiter.

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`affiliate_programs`, `affiliate_accounts`, `affiliate_campaigns`, `affiliate_links`, `affiliate_clicks`, `affiliate_conversions`, `commissions`, `commission_adjustments`, `affiliate_payouts`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)
**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**
- `GET /go/{slug}`
- `POST /api/v1/webhooks/affiliate/{network}`
- `POST /api/v1/admin/affiliate/links`
- `POST /api/v1/admin/affiliate/reconcile-csv`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. UI/UX 설계 (User Interface & Experience)
Clear UI distinction between `Official Direct Link` and `Partner Affiliate Link (Supports CloudScanr at no extra cost)` with promo badge (`$200 Free Trial Credit`). Admin Campaign & Conversion Reconciliation UI.

---

## 8. 상세 처리 로직 (Core Processing Logic)
1. Zero Fake Conversions Rule (Article 10 & 29): Clicking `/go/{slug}` ONLY inserts an `affiliate_clicks` record. A conversion in `affiliate_conversions` is NEVER created unless confirmed by an authenticated network postback or verified CSV/API report from the affiliate network.
2. Anti-Fraud Filter: Blocks automated headless browsers, known datacenter crawler ASNs, >10 clicks/min from the same IP subnet (stores privacy-preserving keyed hash `HMAC(ip, daily_salt)`), and self-referrals.
3. Link Health Monitor: Hourly HEAD/GET check verifies destination URLs return 200/301/302; if a campaign expires or 404s, automatically falls back to the clean provider URL and alerts Admin.

---

## 9. 보안 통제 (Security & Access Control)
Webhook postback endpoints verify network HMAC signatures + IP allowlists; strict open-redirect prevention (`/go/{slug}` only redirects to pre-registered, admin-approved `affiliate_links` database records, never arbitrary `?url=` query strings).

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
If database write queue is full during traffic spike, `/go/{slug}` still redirects immediately using cached link target and buffers click events in local disk/Redis ring buffer.

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Depends on Business (`B_BUSINESS.md`) and PostgreSQL (`P_POSTGRESQL.md`).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
Open-redirect penetration test, SubID HMAC tamper rejection test, bot User-Agent filtering test, and end-to-end click -> postback -> commission ledger integration test.

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] All 26 mandatory affiliate engine features in Article 10 are implemented.
- [x] `/go/{slug}` p99 latency < 15ms and strictly blocks arbitrary open redirects.
- [x] Zero unverified clicks are ever recorded as conversions or commissions.

---

## 14. 검증 증빙 (Required Verification Evidence)
Click redirect benchmark logs, webhook signature verification test results, and reconciliation audit trail.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0351 ~ WO-CS-0450`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
