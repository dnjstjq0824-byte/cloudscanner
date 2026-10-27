# AFFILIATE_CONTRACT.md — Affiliate Redirect, Attribution & Ledger API Contract

**Document ID:** `CS-API-004`

---

## 1. Redirect Endpoint Contract (`GET /go/{slug}`)
- **Request**: `GET /go/digitalocean-s-2vcpu-4gb?src=compare_matrix&slot=row_1`
- **Processing**:
  1. Look up `slug` in Redis hot-cache (fallback to PostgreSQL `affiliate_links`). If unknown, return `404 Not Found` (strictly no arbitrary URL redirect).
  2. Evaluate Bot/Fraud Score.
  3. Generate ULID `click_id` (e.g. `01JA8X9K2M...`) and HMAC signature `sig = HMAC_SHA256(click_id, SUBID_SECRET)[:16]`.
  4. Populate destination template: `https://digitalocean.pxf.io/c/...?subId1=01JA8X9K2M...&subId2=a9f8e1...`.
  5. Enqueue non-blocking `affiliate_clicks` insert.
- **Response**: `HTTP/1.1 302 Found`, `Cache-Control: no-store, no-cache, must-revalidate`, `Referrer-Policy: strict-origin-when-cross-origin`.

---

## 2. Postback & CSV Reconciliation Contract
- **Rule**: Clicking a link never creates a conversion. Conversions are only created via `POST /api/v1/webhooks/affiliate/{network}` (with verified HMAC header) or `POST /api/v1/admin/affiliate/reconcile-csv`.
- **Idempotency**: Keyed on `(program_id, external_conversion_id)`.
