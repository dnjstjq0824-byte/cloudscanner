# SECURITY_E2E.md — Tenant Isolation, KMS Cryptography, SSRF & RBAC Penetration Suite
**Document ID:** `CS-TST-005`
> TirpScanner pivot (2026-11-26): LIVE scope — reframed from BYOC/server endpoints to catalog/affiliate/admin endpoints. See `TIRPSCANNER_PIVOT.md`.
## 1. Security E2E Verification Cases
1. **Cross-Tenant IDOR Matrix**: Automated harness iterates through every `/api/v1/deals/*`, `/api/v1/catalog/*`, `/api/v1/affiliate/*`, and `/api/v1/orgs/*` endpoint using Tenant A's JWT against Tenant B's UUIDs -> Asserts 100% `403`/`404` responses and zero data leakage.
2. **Step-Up MFA Enforcement**: Asserts that catalog/affiliate admin mutations (`POST /api/v1/admin/catalog/*`, `POST /api/v1/admin/affiliate/*`, payout-batch approval, listing deletion) reject requests lacking a fresh (<10m) `X-Step-Up-MFA-Token`.
