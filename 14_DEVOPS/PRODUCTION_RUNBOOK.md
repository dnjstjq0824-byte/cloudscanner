# PRODUCTION_RUNBOOK.md — Incident Response & Operational Runbooks
**Document ID:** `CS-DO-005`
> TirpScanner pivot (2026-11-26): SOP-02 (BYOC freeze) retired; SOP-01/SOP-03 reframed for deal feeds and platform credentials. See `TIRPSCANNER_PIVOT.md`.
## 1. Standard Operating Procedures (SOPs)
- **SOP-01: Deal-Feed API Breaking Schema Change**: Sync job enters `QUARANTINED_ANOMALY` -> Last good deal snapshot remains live -> Update vendor ingestion adapter -> Replay sync job.
- **~~SOP-02: Emergency BYOC Freeze~~ (RETIRED)**: BYOC scope removed; incident freeze now means disabling affected deal-feed sync and affiliate payout batches via feature flag.
- **SOP-03: KMS Key Rotation**: Trigger `cs-migrator kms-rotate-deks` to re-encrypt platform credential `encrypted_dek` values online with zero downtime.
