# DATA_FLOW.md — Core Data Flows (Deal Ingestion, Affiliate Attribution & Price-Drop Alerts)

**Document ID:** `CS-ARCH-003`

> **Pivot addendum (2026-11-26):** Refreshed for TirpScanner. Flow A reframed from cloud pricing ingestion to SaaS/AI-tool deal ingestion and plan-price normalization. Flow B (affiliate) kept and rebranded. Former Flow C (BYOC provisioning) removed — see Section 3 for the retired reference and its replacement deal-scan flow.

---

## 1. Flow A: SaaS Deal Ingestion & Plan-Price Normalization Flow
1. **Trigger**: `cs-worker-deals` cron triggers hourly (Tier-1 vendor pricing/deal feeds) or every 6 hours (affiliate network product feeds and bulk catalogs).
2. **Fetch**: `PricingAdapter` (catalog/deal-ingestion role; name historical) calls official endpoints (e.g., vendor pricing APIs, Impact/CJ product catalog feeds, curated deal sources).
3. **Archive & Hash**: Raw response is hashed (`SHA-256`) and archived. If unchanged from last snapshot, `price_snapshots.verified_at` is bumped.
4. **Anomaly Guard**: If price delta > 25% or plan count drops > 20%, payload is routed to `QUARANTINED` for Admin review.
5. **Normalize & Index**: Validated plans update the SaaS catalog and insert immutable `price_snapshots`. Outbox event triggers bulk update on OpenSearch `tirpscanner_deals_read`, runs price-drop detection against previous snapshots, and evaluates user `price_alerts`.

---

## 2. Flow B: Affiliate Click, Attribution & Commission Ledger Flow
1. **Outbound Click**: User clicks an affiliate CTA (`/go/notion-plus-annual?pos=compare_row_1`).
2. **Bot & Fraud Check**: API checks Cloudflare Bot Score, User-Agent blocklist, and IP subnet rate bucket in Redis.
3. **ClickID Generation**: Generates ULID `click_id`, signs `subId1 = click_id` + `subId2 = HMAC(click_id, secret)`, and immediately returns `HTTP 302` to the official affiliate network URL (`< 10ms`).
4. **Async Click Persistence**: Worker writes `affiliate_clicks` (storing hashed IP prefix, country code, device type, referrer page, plan/deal ID).
5. **Conversion Reconciliation**: Network webhook or Admin CSV import submits `(external_conversion_id, subId1, order_amount, commission_amount, status)`. Worker verifies `subId1` exists in `affiliate_clicks`, writes `affiliate_conversions`, and posts an immutable entry to `commissions` (`PENDING` -> `APPROVED` -> `PAID` or `CLAWBACK` via `commission_adjustments`).

---

## 3. Flow C: Deal Scan & Price-Drop Alert Flow (replaces retired BYOC flow)
1. **Deal Scan**: After each ingestion run, the worker compares new `price_snapshots` against the previous current snapshot per plan to detect price drops, new discounts, and expiring promotions.
2. **Alert Matching**: Price drops are matched against active `price_alerts` (plan-targeted or filter-targeted subscriptions).
3. **Notification Dispatch**: Matched users receive in-app/email notifications; SSRF-safe webhook dispatch applies to custom webhook channels.
4. **Historical Note (retired 2026-11-26):** Former Flow C covered BYOC account connection (OAuth/token + KMS envelope encryption) and a two-phase cost-approval server provisioning protocol. That flow is REMOVED scope under TirpScanner; the BYOC state machine is preserved as retired history in `SYSTEM_DIAGRAMS.md`.
