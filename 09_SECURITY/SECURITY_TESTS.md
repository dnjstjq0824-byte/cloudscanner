# SECURITY_TESTS.md — Automated Security & Penetration Verification Suite
**Document ID:** `CS-SEC-005`
> TirpScanner pivot (2026-11-26): BYOC/server cases retired; tenant, crypto, and affiliate cases reframed below stay mandatory. See `TIRPSCANNER_PIVOT.md`.
## 1. Mandatory Security Test Cases
- `SEC-01`: Cross-Tenant IDOR Fuzzing across all `/api/v1/deals/*`, `/api/v1/catalog/*`, and `/api/v1/orgs/*` routes -> Must return `404 Not Found`.
- `SEC-02`: KMS AAD Swap Attack -> Swapping `ciphertext_blob` between Org A and Org B must fail GCM authentication.
- `SEC-03`: Read-Only Role Mutation Attempt -> Using an `ORG_VIEWER` session to call catalog/affiliate admin mutation endpoints must abort prior to KMS decryption with `403 Forbidden`.
- `SEC-04`: Log Secret Injection Test -> Passing an affiliate network API key value in an error path must output `[REDACTED_SECRET]` in logs.
- `SEC-05`: Open Redirect & Webhook Tamper Test on `/go/{slug}` and `/api/v1/webhooks/affiliate/*`.
