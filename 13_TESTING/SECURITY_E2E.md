# SECURITY_E2E.md — Tenant Isolation, KMS Cryptography, SSRF & RBAC Penetration Suite
**Document ID:** `CS-TST-005`
## 1. Security E2E Verification Cases
1. **Cross-Tenant IDOR Matrix**: Automated harness iterates through every `/api/v1/cloud/*` and `/api/v1/orgs/*` endpoint using Tenant A's JWT against Tenant B's UUIDs -> Asserts 100% `403`/`404` responses and zero data leakage.
2. **Step-Up MFA Enforcement**: Asserts that `POST /api/v1/cloud/connections`, `POST /api/v1/cloud/approvals`, `POST /api/v1/cloud/servers`, and `DELETE /api/v1/cloud/servers/{id}` reject requests lacking a fresh (<10m) `X-Step-Up-MFA-Token`.
