# ADMIN_CONSOLE.md — Operational Workflows & Governance Runbook
**Document ID:** `CS-ADM-001`
> TirpScanner pivot (2026-11-26): LIVE scope — these governance rules now cover catalog, deal, affiliate, payout, and content operations. BYOC/server admin pages are retired. See `TIRPSCANNER_PIVOT.md`.
## 1. Admin Operational Principles
- Every state-changing action in `/admin/*` requires:
  1. Authenticated Admin Role (`SUPER_ADMIN`, `CATALOG_ADMIN`, `AFFILIATE_ADMIN`).
  2. Hardware WebAuthn Passkey step-up token (`X-Step-Up-MFA-Token`).
  3. Non-empty `justification_reason` string saved to `admin_actions` along with `before_state_jsonb` and `after_state_jsonb`.
