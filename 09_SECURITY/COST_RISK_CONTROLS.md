# COST_RISK_CONTROLS.md — Cloud Financial Blast-Radius & Runaway Spend Prevention
**Document ID:** `CS-SEC-004`
## 1. Mandatory Financial Safeguards
1. **5-Minute Single-Use Approval Token**: `operation_approvals` expires in 300s and can only be consumed once (`consumed_at IS NULL`).
2. **Velocity Circuit Breaker**: Max 5 server creation operations per organization per 10 minutes by default (`MAX_PROVISION_VELOCITY_10M = 5`), preventing Compromised-Session or Script Loop bill shock.
3. **Timeout Reconciliation Before Retry**: Eliminates duplicate server creation when provider API responses are delayed.
