# T — Third-Party Integrations — Revised Adapter Contracts, Affiliate Networks, FX Feeds & Notification Channels

**Document ID:** `CS-AZ-T` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (TirpScanner pivot reframe, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Pivot Note (2026-11-26):** Revised adapter law: Affiliate + Catalog/Deal-ingestion LIVE; Management (BYOC) + Reseller RETIRED. See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0201 ~ WO-CS-0250, WO-CS-0570 ~ WO-CS-0620`

---

## 1. 기능 목적 (Functional Purpose)

Standardize all external third-party integrations behind clean Go interfaces with uniform timeouts, retries, circuit breakers, rate limiters, and observability metrics — covering affiliate networks, deal/coupon feeds, FX feeds, and notification channels.

---

## 2. User Roles & Actors

Vendor Feed Agent, Backend Agent, DevOps Agent.

---

## 3. User Scenarios

1. The catalog/deal-ingestion adapter registry executes `SyncCatalog(ctx)` across networks and feeds using source-specific rate limiters (e.g., Impact, CJ, Awin, PartnerStack quotas; Product Hunt GraphQL limits).
2. The `AffiliateNetworkAdapter` ingests conversion webhooks and daily settlement deltas from PartnerStack, Impact, CJ Affiliate, Awin, and ShareASale.
3. The Notification Engine dispatches price-drop and deal alerts via Email (SES/Postmark) and Webhook (Slack/Discord).

---

## 4. Technology Stack

Go `golang.org/x/time/rate`, `sony/gobreaker`, OpenTelemetry HTTP client instrumentation, Webhook HMAC-SHA256 signer.

---

## 5. Database Schema & Tables

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

## 7. User Interface & Experience

Admin Third-Party Integration Status Matrix (`/admin/integrations`) displaying real-time circuit breaker state (`CLOSED`, `HALF-OPEN`, `OPEN`), remaining API rate-limit quota, and error rates per network/feed.

---

## 8. Core Processing Logic

Strict Go Interface Segregation (`04_API_SPECIFICATION/PROVIDER_ADAPTER_CONTRACT.md`):

- `PricingAdapter` (LIVE — catalog/deal-ingestion; name historical): `FetchPlans(ctx) ([]RawPlan, error)`, `FetchRegions(ctx) ([]RawRegion, error)`, `FetchPromotions(ctx) ([]RawPromo, error)`.
- `AffiliateAdapter` (LIVE): `BuildTrackingURL(link, clickID, subID)`, `VerifyPostback(r *http.Request)`, `FetchConversions(ctx, window)`.
- `ManagementAdapter` (RETIRED — BYOC): `VerifyCredentials`, `ListServers`, `CreateServer`, `ExecutePowerAction`, `DeleteServer` — frozen, no new calls.
- `ResellerAdapter` (RETIRED): never activated; no contract surface.

---

## 9. Security & Access Control

Outbound SSRF protection on user-supplied notification webhooks (blocks RFC1918 private IPs `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`, link-local `169.254.169.254` metadata endpoints, and loopback).

---

## 10. Failure Modes & Resilience

Per-source circuit breakers isolate failures so a slow coupon feed never blocks PartnerStack or Impact syncs.

---

## 11. Inter-Module Dependencies

Connects external networks/feeds to Core Domains (`D`, `E`, `L`). Cloud Management (`M`) is RETIRED.

---

## 12. Testing & Verification Strategy

SSRF blocklist unit tests for user webhooks, rate-limiter throttling tests, and circuit-breaker state transition tests.

---

## 13. Acceptance Criteria

- [x] Every external HTTP call enforces a strict context timeout (`<= 13s` for read, `<= 30s` for conversion fetch) and passes through a rate limiter and circuit breaker.
- [x] User webhook dispatcher blocks all internal/metadata IP ranges (SSRF protection).
- [x] Retired adapters expose no new live wiring or credentials.

---

## 14. Required Verification Evidence

Adapter interface contract tests and SSRF test suite output.

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0201 ~ WO-CS-0250, WO-CS-0570 ~ WO-CS-0620`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
