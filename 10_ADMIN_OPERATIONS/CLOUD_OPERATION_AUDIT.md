# CLOUD_OPERATION_AUDIT.md — BYOC Operation Auditing, Drift Reconciliation & Incident Freeze
**Document ID:** `CS-ADM-004`
## 1. Global BYOC Audit & Emergency Freeze
- Platform admins can monitor aggregate provisioning error rates, timeout reconciliations, and security anomalies in `/admin/cloud/operations-audit` **without** ever having access to decrypt user BYOC tokens.
- If a tenant account exhibits compromised behavior, `SUPER_ADMIN` can click **Freeze Connection**, setting `cloud_connections.connection_status = 'FROZEN_SECURITY'` and blocking all worker tasks for that connection immediately.
