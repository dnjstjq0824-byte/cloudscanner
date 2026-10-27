# SECURITY_TESTS.md — Automated Security & Penetration Verification Suite
**Document ID:** `CS-SEC-005`
## 1. Mandatory Security Test Cases
- `SEC-01`: Cross-Tenant IDOR Fuzzing across all `/api/v1/cloud/*` routes -> Must return `404 Not Found`.
- `SEC-02`: KMS AAD Swap Attack -> Swapping `ciphertext_blob` between Org A and Org B must fail GCM authentication.
- `SEC-03`: Read-Only Token Mutation Attempt -> Using a `READ_ONLY` connection to call `POST /api/v1/cloud/servers` must abort prior to KMS decryption with `403 Forbidden`.
- `SEC-04`: Log Secret Injection Test -> Passing `dop_v1_abcdef123456...` in an error path must output `[REDACTED_SECRET]` in logs.
- `SEC-05`: Open Redirect & Webhook Tamper Test on `/go/{slug}` and `/api/v1/webhooks/affiliate/*`.
