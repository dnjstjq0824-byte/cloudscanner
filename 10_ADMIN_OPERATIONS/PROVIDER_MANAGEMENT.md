# PROVIDER_MANAGEMENT.md — Provider Onboarding, Sync Quarantine & Rate-Card Maintenance
**Document ID:** `CS-ADM-002`
## 1. Quarantine Resolution Workflow (`/admin/pricing-failures`)
1. Automated sync worker detects an anomaly (e.g., SKU price changed >25% or HTTP 403/500 on Pricing API) and sets `price_sync_jobs.status = 'QUARANTINED_ANOMALY'`.
2. Catalog Admin inspects the JSON diff and official provider URL in `/admin/pricing-failures`.
3. Admin either **Approves** (promoting the new snapshot to `FRESH` and re-indexing OpenSearch) or **Rejects** (keeping prior snapshot and filing an adapter parser fix).
