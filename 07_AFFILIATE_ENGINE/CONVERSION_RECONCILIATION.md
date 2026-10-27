# CONVERSION_RECONCILIATION.md — Webhook Postback & CSV Reconciliation Engine
**Document ID:** `CS-AE-003`
## 1. Three-Source Reconciliation Hierarchy
1. **Real-Time S2S Postback Webhook**: Supported by Impact, CJ, and HasOffers. Verifies network HMAC signature, matches `click_ulid`, and inserts `affiliate_conversions` (`status = PENDING`).
2. **Scheduled Network Reporting API**: Polls Impact `/Mediapartners/{AccountId}/Actions` and CJ GraphQL Commission Detail API every 6 hours to capture delayed conversions, lock-period approvals, and clawbacks.
3. **Deterministic CSV Reconciliation Import**: For providers without public conversion APIs (Cloudways, Kamatera, Vultr portal exports), Admin uploads the official CSV in `/admin/affiliate/reconciliation`. The parser normalizes columns, matches `data1`/`subId` to `affiliate_clicks.click_ulid`, and generates an auditable reconciliation summary.
