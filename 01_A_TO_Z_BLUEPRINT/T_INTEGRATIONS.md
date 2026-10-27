# T — Third-Party Integrations — Unified Adapter Contracts, Affiliate Networks, FX Feeds & Notification Channels

**Document ID:** `CS-AZ-T`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0201 ~ WO-CS-0250, WO-CS-0570 ~ WO-CS-0620`

---

## 1. 기능 목적 (Functional Purpose)
Standardize all external third-party integrations behind clean Go interfaces with uniform timeouts, retries, circuit breakers, rate limiters, and observability metrics.

---

## 2. 사용자 역할 (User Roles & Actors)
Provider API Agent, Backend Agent, DevOps Agent.

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. The `PricingAdapter` registry executes `SyncCatalog(ctx)` across 18 providers using provider-specific rate limiters (e.g., DigitalOcean 5,000 req/hr, Hetzner 3,600 req/hr, Vultr 30 req/sec).
2. The `AffiliateNetworkAdapter` ingests conversion webhooks and daily settlement deltas from Impact Radius, CJ Affiliate, PartnerStack, and Awin.
3. The Notification Engine dispatches price-drop and budget-threshold alerts via Email (SES/Postmark) and Webhook (Slack/Discord).

---

## 4. 기술 스택 (Technology Stack)
Go `golang.org/x/time/rate`, `sony/gobreaker`, OpenTelemetry HTTP client instrumentation, Webhook HMAC-SHA256 signer.

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`provider_data_sources`, `affiliate_programs`, `fx_rates`, `price_alerts`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)
**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**
- `POST /api/v1/webhooks/affiliate/{network}`
- `POST /api/v1/notifications/test-webhook`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. UI/UX 설계 (User Interface & Experience)
Admin Third-Party Integration Status Matrix (`/admin/integrations`) displaying real-time circuit breaker state (`CLOSED`, `HALF-OPEN`, `OPEN`), remaining API rate-limit quota, and error rates per provider.

---

## 8. 상세 처리 로직 (Core Processing Logic)
Strict Go Interface Segregation (`04_API_SPECIFICATION/PROVIDER_ADAPTER_CONTRACT.md`):
- `PricingAdapter`: `FetchPlans(ctx) ([]RawPlan, error)`, `FetchRegions(ctx) ([]RawRegion, error)`.
- `ManagementAdapter`: `VerifyCredentials(ctx, cred)`, `ListServers(ctx, cred)`, `CreateServer(ctx, cred, req)`, `ExecutePowerAction(ctx, cred, id, action)`, `DeleteServer(ctx, cred, id)`.
- `AffiliateAdapter`: `BuildTrackingURL(link, clickID, subID)`, `VerifyPostback(r *http.Request)`, `FetchConversions(ctx, window)`.

---

## 9. 보안 통제 (Security & Access Control)
Outbound SSRF protection on user-supplied notification webhooks (blocks RFC1918 private IPs `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`, link-local `169.254.169.254` metadata endpoints, and loopback).

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
Per-provider circuit breakers isolate failures so a slow OVHcloud API never blocks DigitalOcean or UpCloud requests.

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Connects external providers to Core Domains (`D`, `E`, `L`, `M`).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
SSRF blocklist unit tests for user webhooks, rate-limiter throttling tests, and circuit-breaker state transition tests.

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] Every external HTTP call enforces a strict context timeout (`<= 15s` for read, `<= 30s` for provisioning) and passes through a rate limiter and circuit breaker.
- [x] User webhook dispatcher blocks all internal/metadata IP ranges (SSRF protection).

---

## 14. 검증 증빙 (Required Verification Evidence)
Adapter interface contract tests and SSRF test suite output.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0201 ~ WO-CS-0250, WO-CS-0570 ~ WO-CS-0620`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
