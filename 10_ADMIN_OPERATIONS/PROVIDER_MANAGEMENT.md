# PROVIDER_MANAGEMENT.md — Vendor Onboarding, Sync Quarantine & Deal-Feed Maintenance
**Document ID:** `CS-ADM-002`
> TirpScanner pivot (2026-11-26): LIVE scope — reframed from cloud providers to SaaS/AI-tool vendors and deal feeds. See `TIRPSCANNER_PIVOT.md`.
## 1. Quarantine Resolution Workflow (`/admin/pricing-failures`)
1. Automated sync worker detects an anomaly (e.g., plan price changed >25%, coupon terms mismatch, or HTTP 403/500 on a deal-feed API) and sets `price_sync_jobs.status = 'QUARANTINED_ANOMALY'`.
2. Catalog Admin inspects the JSON diff and official vendor URL in `/admin/pricing-failures`.
3. Admin either **Approves** (promoting the new snapshot to `FRESH` and re-indexing OpenSearch) or **Rejects** (keeping prior snapshot and filing an adapter parser fix).
