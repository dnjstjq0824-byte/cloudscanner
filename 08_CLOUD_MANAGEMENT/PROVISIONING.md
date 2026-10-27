# PROVISIONING.md — Idempotent Server Provisioning & Timeout Reconciliation Engine (Article 13 & 21)
**Document ID:** `CS-CM-003`
## 1. 20-Step Server Provisioning Workflow (Matching Article 21 `WO-CS-0601`)
1. Verify Provider API capability & connection `READ_WRITE` access level.
2. Verify `organization_id` ownership via JWT + PostgreSQL RLS.
3. Fetch live supported regions for the connected account.
4. Fetch live supported plans/SKUs available in the selected region.
5. Fetch OS images / distributions (`ubuntu-24-04-x64`, `debian-12-x64`, etc.).
6. Validate user hardware specification selection.
7. Compute itemized hourly & monthly TCO (including IPv4, backups, extra storage).
8. Validate RFC-1123 server hostname.
9. Select or register SSH public key IDs.
10. Configure optional block storage volumes & automated backups.
11. Require signed `operation_approvals` record + Step-Up MFA token.
12. Lock client `Idempotency-Key` in `resource_operations`.
13. Enqueue `cloud:provision` job in `background_jobs`.
14. Inject deterministic tag `cloudscanr-op-id:{operation_id}` and invoke Provider Create API; persist `external_resource_id` immediately upon response.
15. Poll provider resource state (`new` -> `active` / `running`) with exponential backoff.
16. **Timeout Reconciliation Rule**: If the initial `POST` times out or drops connection, query provider list API filtered by tag `cloudscanr-op-id:{operation_id}` BEFORE ever attempting a retry.
17. Prevent duplicate server creation and double billing under all network fault conditions.
18. Record structured failure reason if provider rejects request (e.g., quota exceeded).
19. Emit real-time completion status + assigned IPv4/IPv6 to User Dashboard.
20. Append hash-chained record to `audit_logs`.
