# PRODUCTION_RUNBOOK.md — Incident Response & Operational Runbooks
**Document ID:** `CS-DO-005`
## 1. Standard Operating Procedures (SOPs)
- **SOP-01: Provider Pricing API Breaking Schema Change**: Sync job enters `QUARANTINED_ANOMALY` -> Last good price remains live -> Update provider `PricingAdapter` -> Replay sync job.
- **SOP-02: Emergency BYOC Freeze**: Execute `POST /api/v1/admin/security/connections/{id}/freeze` or global feature flag `disable_all_byoc_mutations = true`.
- **SOP-03: KMS Key Rotation**: Trigger `cs-migrator kms-rotate-deks` to re-encrypt `provider_credentials_metadata.encrypted_dek` online with zero downtime.
