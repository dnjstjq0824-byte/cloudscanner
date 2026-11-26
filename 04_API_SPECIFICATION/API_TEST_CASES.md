# API_TEST_CASES.md — Contract & Security Test Suite Across Live API Groups (BYOC Groups Retired)

**Document ID:** `CS-API-006`

> **Pivot addendum (2026-11-26):** Refreshed for TirpScanner. `ATC-14`–`ATC-19` covered BYOC provisioning/budget/cost endpoints and are REMOVED. Remaining cases reframed to the SaaS/AI-tool deal surface.

---

## 1. Mandatory API Test Matrix

| Test ID | API Group | Endpoint | Scenario Tested | Expected HTTP Status & Assertion |
| :--- | :--- | :--- | :--- | :--- |
| `ATC-01` | 01_Auth | `POST /api/v1/auth/register` | Valid registration vs Duplicate email | `201 Created` then `409 Conflict` |
| `ATC-02` | 01_Auth | `POST /api/v1/auth/login` | Brute-force >5 failed attempts in 15m | `429 Too Many Requests` + `Retry-After` header |
| `ATC-03` | 02_Users | `DELETE /api/v1/users/me` | Delete account without `X-Step-Up-MFA-Token` | `403 Forbidden` (`STEP_UP_MFA_REQUIRED`) |
| `ATC-04` | 03_Orgs | `POST /api/v1/orgs/{id}/members` | `ORG_VIEWER` attempts to invite an admin | `403 Forbidden` (`INSUFFICIENT_RBAC_ROLE`) |
| `ATC-05` | 04_Vendors | `GET /api/v1/providers/notion` | Verify vendor affiliate status & reward type | `200 OK` (`reward_type: CASH_COMMISSION`) |
| `ATC-06` | 05_Catalog | `GET /api/v1/catalog/skus/{id}` | Verify all plan-price normalization attributes present | `200 OK` (Passes JSON Schema for normalization fields) |
| `ATC-07` | 06_Pricing | `POST /api/v1/pricing/estimate` | Compare monthly vs annual billing + seat scaling | `200 OK` (Exact decimal match to 8 places) |
| `ATC-08` | 07_Comparison | `POST /api/v1/comparisons` | Create comparison with 4 plans across currencies | `201 Created` with `share_slug` |
| `ATC-09` | 08_Search | `GET /api/v1/search/skus?q=ai+video+editor+under+%2420` | Natural language deal query + OpenSearch facet | `200 OK` (`p95 < 50ms`, filter chips extracted) |
| `ATC-10` | 09_Promotions | `GET /api/v1/promotions` | Verify expired promotions are excluded | `200 OK` (All items `expires_at > NOW()`) |
| `ATC-11` | 10_Affiliate | `GET /api/v1/affiliate/disclosure` | Verify FTC/KFTC disclosure text returned | `200 OK` |
| `ATC-12` | 11_Tracking | `GET /go/unknown-malicious-slug` | Attempt open redirect via unknown slug | `404 Not Found` (Zero external redirect) |
| `ATC-13` | 11_Tracking | `POST /api/v1/webhooks/affiliate/impact` | Replay identical conversion postback twice | First `200 OK (created)`, Second `200 OK (idempotent_noop)` |
| `ATC-14` | 12_Connection | `POST /api/v1/cloud/connections` | **REMOVED 2026-11-26** — BYOC retired | Endpoint removed; expect `404` |
| `ATC-15` | 13_Provision | `POST /api/v1/cloud/servers` | **REMOVED 2026-11-26** — BYOC retired | Endpoint removed; expect `404` |
| `ATC-16` | 13_Provision | `POST /api/v1/cloud/servers` | **REMOVED 2026-11-26** — BYOC retired | Endpoint removed; expect `404` |
| `ATC-17` | 14_Management | `GET /api/v1/cloud/servers/{id}` | **REMOVED 2026-11-26** — BYOC retired | Endpoint removed; expect `404` |
| `ATC-18` | 15_Budget | `POST /api/v1/cloud/servers` | **REMOVED 2026-11-26** — BYOC retired | Endpoint removed; expect `404` |
| `ATC-19` | 16_Cost | `GET /api/v1/cloud/costs/summary` | **REMOVED 2026-11-26** — BYOC retired | Endpoint removed; expect `404` |
| `ATC-20` | 17_Alert | `POST /api/v1/price-alerts` | Register webhook pointing to `http://169.254.169.254` | `422 Unprocessable Entity` (`SSRF_PRIVATE_IP_BLOCKED`) |
| `ATC-21` | 18_Analytics | `GET /api/v1/analytics/affiliate-yield` | Verify cash vs credit separation | `200 OK` |
| `ATC-22` | 19_Admin | `POST /api/v1/admin/pricing/quarantine/{id}/resolve` | Resolve anomaly without `justification_reason` | `400 Bad Request` (`REASON_REQUIRED`) |
| `ATC-23` | 20_Content | `POST /api/v1/content/reviews` | Submit review as verified customer | `201 Created` (`verified_customer: true`) |
| `ATC-24` | 21_Audit | `GET /api/v1/audit-logs` | Verify cryptographic hash chain (`prev_hash` -> `entry_hash`) | `200 OK` (`chain_valid: true`) |
| `ATC-25` | 22_Notify | `GET /api/v1/notifications` | Fetch unread user notifications | `200 OK` |
