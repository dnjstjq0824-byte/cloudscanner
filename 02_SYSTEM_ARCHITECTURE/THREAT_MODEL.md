# THREAT_MODEL.md — STRIDE Threat Model & Blast-Radius Mitigation

**Document ID:** `CS-ARCH-005`

---

## 1. STRIDE Threat Matrix Across Critical Assets

| Threat ID | STRIDE Category | Target Asset | Attack Vector | Mandatory Architectural Mitigation | Verification Test |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `TM-01` | **Information Disclosure** | User BYOC Cloud API Tokens | SQL Injection, DB snapshot leak, or log exposure | AES-256-GCM KMS Envelope Encryption with `org_id:connection_id` AAD; zero plaintext columns; automated `slog` secret redaction filter. | `WO-CS-0554`, `SECURITY_E2E.md` |
| `TM-02` | **Elevation of Privilege (IDOR)** | Connected Cloud Servers of Tenant B | Authenticated User A manipulates `{server_id}` or `{connection_id}` in API request | Dual enforcement: Go middleware org membership check + PostgreSQL Row-Level Security (`RLS`) policy on `organization_id`. | `WO-CS-0560`, Article 24 Test #15 |
| `TM-03` | **Tampering / Financial Loss** | Unauthorized Cloud Spend / Double Server Creation | Double-click, network timeout retry, or CSRF triggering billable VM creation | Mandatory cryptographic `operation_approvals` token + client `Idempotency-Key` + `TIMEOUT_RECONCILING` tag check before retry. | `WO-CS-0602`, Article 24 Tests #10, #11, #12 |
| `TM-04` | **Repudiation** | Destructive Server Deletion or Admin Payout Approval | Malicious insider or compromised user deletes production VM and denies action | Step-up MFA required + append-only hash-chained `audit_logs` (`prev_hash`, `entry_hash`) protected by DB immutability trigger. | `WO-CS-0830`, Article 24 Test #22 |
| `TM-05` | **Spoofing** | Affiliate Conversion Postbacks | Attacker sends fake webhook to `/api/v1/webhooks/affiliate/impact` to inflate commissions | HMAC-SHA256 webhook signature verification + network IP allowlist + mandatory join to existing signed `affiliate_clicks.click_id`. | `WO-CS-0415`, `AFFILIATE_E2E.md` |
| `TM-06` | **Denial of Service** | Public Search & `/go/{slug}` Redirect | Botnet floods search facets or affiliate redirect links | Cloudflare WAF rate limiting + Redis token bucket + PostgreSQL fallback + in-memory LRU cache for `/go/{slug}`. | `WO-CS-0915` |
| `TM-07` | **Server-Side Request Forgery (SSRF)** | Internal VPC / Cloud Metadata (`169.254.169.254`) | User registers custom notification webhook pointing to internal IP | Custom Go `net.Dialer` control hook blocking RFC1918, loopback, link-local (`169.254.0.0/16`), and IPv6 unique-local addresses. | `WO-CS-0925` |
