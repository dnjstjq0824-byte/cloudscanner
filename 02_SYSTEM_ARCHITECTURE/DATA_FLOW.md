# DATA_FLOW.md — Core Data Flows (Pricing Sync, Affiliate Attribution & BYOC Provisioning)

**Document ID:** `CS-ARCH-003`

---

## 1. Flow A: Official Pricing Ingestion & 36-Field Normalization Flow
1. **Trigger**: `cs-worker-pricing` cron triggers hourly (Tier-1 REST APIs) or every 6 hours (Hyperscaler bulk catalogs).
2. **Fetch**: `PricingAdapter` calls official endpoint (e.g., `GET https://api.digitalocean.com/v2/sizes`, `GET https://api.vultr.com/v2/plans`, `GET https://api.hetzner.cloud/v1/pricing`, `GET https://api.linode.com/v4/linode/types`).
3. **Archive & Hash**: Raw response is hashed (`SHA-256`) and archived. If unchanged from last snapshot, `price_snapshots.verified_at` is bumped.
4. **Anomaly Guard**: If price delta > 25% or SKU count drops > 20%, payload is routed to `QUARANTINED` for Admin review (`WO-CS-0712`).
5. **Normalize & Index**: Validated SKUs update `product_skus`, `product_specs`, and insert immutable `price_snapshots`. Outbox event triggers bulk update on OpenSearch `cloudscanr_skus_read` and evaluates user `price_alerts`.

---

## 2. Flow B: Affiliate Click, Attribution & Commission Ledger Flow
1. **Outbound Click**: User clicks an affiliate CTA (`/go/upcloud-dev-2cpu-4gb?pos=compare_row_1`).
2. **Bot & Fraud Check**: API checks Cloudflare Bot Score, User-Agent blocklist, and IP subnet rate bucket in Redis.
3. **ClickID Generation**: Generates ULID `click_id`, signs `subId1 = click_id` + `subId2 = HMAC(click_id, secret)`, and immediately returns `HTTP 302` to the official affiliate network URL (`< 10ms`).
4. **Async Click Persistence**: Worker writes `affiliate_clicks` (storing hashed IP prefix, country code, device type, referrer page, SKU ID).
5. **Conversion Reconciliation**: Network webhook or Admin CSV import submits `(external_conversion_id, subId1, order_amount, commission_amount, status)`. Worker verifies `subId1` exists in `affiliate_clicks`, writes `affiliate_conversions`, and posts an immutable entry to `commissions` (`PENDING` -> `APPROVED` -> `PAID` or `CLAWBACK` via `commission_adjustments`).

---

## 3. Flow C: BYOC Account Connection & Server Lifecycle Flow
1. **Account Connection**: User authenticates via Provider OAuth 2.0 or submits a scoped API token over TLS 1.3 after step-up MFA.
2. **Preflight Scope Check**: Backend queries provider read endpoint (`/v2/account` or `/v1/servers?per_page=1`) to verify token validity, account ID, and read vs. write scope.
3. **Envelope Encryption**: KMS generates a 256-bit DEK; token is encrypted via `AES-256-GCM` with `AAD = org_id:connection_id:provider_id`. Only `encrypted_dek` and `ciphertext_blob` are saved to `provider_credentials_metadata`.
4. **Approved Provisioning**: User completes the 4-step wizard, reviews the itemized TCO estimate, and signs `POST /api/v1/cloud/approvals`. `POST /api/v1/cloud/servers` verifies the approval token + `Idempotency-Key`, creates a `resource_operations` record, and dispatches to `cs-worker-cloud`.
