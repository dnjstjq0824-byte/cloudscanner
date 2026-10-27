# AFFILIATE_E2E.md — Affiliate Click, Attribution, Reconciliation & Anti-Fraud E2E Specification
**Document ID:** `CS-TST-003`
## 1. End-to-End Monetization Verification
- **Scenario 1**: Click `/go/{slug}` -> Verify `HTTP 302` redirect with signed `subId1={click_ulid}` and `X-Robots-Tag: noindex, nofollow`.
- **Scenario 2**: Bot User-Agent (`Googlebot`, `curl`, headless scraper) hits `/go/{slug}` -> Redirect succeeds for UX, but `is_bot_filtered = true` so click statistics are never polluted.
- **Scenario 3**: Signed postback and CSV reconciliation match `click_ulid`, transition `commissions` through `PENDING` -> `APPROVED` -> `SETTLED_PAID`, and handle `CLAWBACK_REFUND` via `commission_adjustments`.
