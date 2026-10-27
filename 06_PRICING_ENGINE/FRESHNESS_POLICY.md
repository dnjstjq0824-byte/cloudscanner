# FRESHNESS_POLICY.md — Price Snapshot Freshness SLA & Stale Quarantine Policy
**Document ID:** `CS-PE-004`
## 1. Freshness State Transitions
- `FRESH`: Snapshot verified within `<= 6 hours` (Tier-1 REST APIs) or `<= 24 hours` (Bulk feeds).
- `STALE`: Snapshot age exceeds `24 hours` due to upstream API outage; UI displays amber warning badge `Price verified YYYY-MM-DD HH:MM UTC — Check official site before purchase`.
- `QUARANTINED`: Automated sync detected a `> 25%` price swing or missing required fields; held for Catalog Admin inspection in `/admin/pricing-failures` while keeping the last known good snapshot live with a `PENDING_REVERIFICATION` notice.
