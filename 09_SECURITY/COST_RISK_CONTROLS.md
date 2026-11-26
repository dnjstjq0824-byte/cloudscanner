# COST_RISK_CONTROLS.md — Platform Financial Blast-Radius & Runaway Spend Prevention
**Document ID:** `CS-SEC-004`
> TirpScanner pivot (2026-11-26): server-provisioning spend guardrails are retired with BYOC scope. The controls below are reframed for affiliate payout and deal-ingestion integrity; see `TIRPSCANNER_PIVOT.md`.
## 1. Mandatory Financial Safeguards
1. **5-Minute Single-Use Approval Token**: payout-batch and commission-adjustment approvals expire in 300s and can only be consumed once (`consumed_at IS NULL`).
2. **Velocity Circuit Breaker**: Max 5 payout-batch or bulk commission-adjustment operations per organization per 10 minutes by default (`MAX_PAYOUT_VELOCITY_10M = 5`), preventing Compromised-Session or Script Loop payout fraud.
3. **Reconciliation Before Retry**: Eliminates duplicate payouts or duplicate commission postings when affiliate network or banking responses are delayed; every retry re-queries settlement state first.
