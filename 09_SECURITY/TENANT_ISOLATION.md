# TENANT_ISOLATION.md — Multi-Layer Tenant Isolation & PostgreSQL RLS
**Document ID:** `CS-SEC-003`
> TirpScanner pivot (2026-11-26): isolation layers stay; BYOC table specifics dropped and reframed for tenant-scoped catalog/affiliate data. See `TIRPSCANNER_PIVOT.md`.
## 1. Defense-in-Depth Isolation Layers
1. **Layer 1 (API Middleware)**: Validates `X-Organization-Id` against `organization_members` for the authenticated `user_id`.
2. **Layer 2 (Database Transaction Wrapper)**: Executes `SET LOCAL app.current_org_id = $1;` on every connection checked out from `pgxpool` before running any query.
3. **Layer 3 (PostgreSQL RLS Policy)**: Enforces `organization_id = current_setting('app.current_org_id')::uuid` on all tenant-scoped tables (catalog favorites, alerts, affiliate accounts, commissions, credential metadata).
4. **Layer 4 (Cryptographic AAD Binding)**: Even if a SQL bug returned another tenant's platform credential row, `AES-256-GCM` decryption fails with `cipher: message authentication failed` because `org_id` in the AAD does not match.
