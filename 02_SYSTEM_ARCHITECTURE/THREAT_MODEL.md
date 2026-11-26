# THREAT_MODEL.md — STRIDE Threat Model & Blast-Radius Mitigation

**Document ID:** `CS-ARCH-005`

> **Pivot addendum (2026-11-26):** BYOC token (`TM-01`), cross-tenant server (`TM-02` original), and unauthorized cloud-spend (`TM-03`) threats are RETIRED with BYOC removal — no user cloud credentials are stored and no billable provisioning exists. `TM-02` is retargeted to tenant deal/alert data and `TM-04` to admin financial approvals. `TM-05`–`TM-07` kept.

---

## 1. STRIDE Threat Matrix Across Critical Assets

| Threat ID | STRIDE Category | Target Asset | Attack Vector | Mandatory Architectural Mitigation | Verification Test |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `TM-01` | **RETIRED** | (Formerly: user BYOC cloud API tokens) | — | **Retired 2026-11-26.** No BYOC credentials are collected or stored; KMS envelope encryption now protects affiliate/integration credentials (BYOC-token use retired). | — |
| `TM-02` | **Elevation of Privilege (IDOR)** | Tenant deal alerts, saved lists & favorites | Authenticated User A manipulates `{alert_id}` or `{org_id}` in API request | Dual enforcement: Go middleware org membership check + PostgreSQL Row-Level Security (`RLS`) policy on `organization_id`/`user_id`. | Live `SECURITY_E2E` IDOR matrix + `SEARCH_E2E` gates |
| `TM-03` | **RETIRED** | (Formerly: unauthorized cloud spend / double server creation) | — | **Retired 2026-11-26.** No billable provisioning endpoints exist; `Idempotency-Key` remains mandatory on affiliate reconciliation and admin mutations. | — |
| `TM-04` | **Repudiation** | Admin payout approval & destructive admin actions | Malicious insider or compromised admin approves a payout and denies action | Step-up MFA required + append-only hash-chained `audit_logs` (`prev_hash`, `entry_hash`) protected by DB immutability trigger. | Live `AFFILIATE_E2E` + `SECURITY_E2E` gates |
| `TM-05` | **Spoofing** | Affiliate Conversion Postbacks | Attacker sends fake webhook to `/api/v1/webhooks/affiliate/impact` to inflate commissions | HMAC-SHA256 webhook signature verification + network IP allowlist + mandatory join to existing signed `affiliate_clicks.click_id`. | `AFFILIATE_E2E.md` |
| `TM-06` | **Denial of Service** | Public Search & `/go/{slug}` Redirect | Botnet floods search facets or affiliate redirect links | Cloudflare WAF rate limiting + Redis token bucket + PostgreSQL fallback + in-memory LRU cache for `/go/{slug}`. | Load-test gate |
| `TM-07` | **Server-Side Request Forgery (SSRF)** | Internal VPC / Cloud Metadata (`169.254.169.254`) | User registers custom notification webhook pointing to internal IP | Custom Go `net.Dialer` control hook blocking RFC1918, loopback, link-local (`169.254.0.0/16`), and IPv6 unique-local addresses. | SSRF test family |
