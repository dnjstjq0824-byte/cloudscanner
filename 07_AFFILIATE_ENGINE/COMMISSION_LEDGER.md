# COMMISSION_LEDGER.md — Double-Entry Commission Accounting & Payout Settlement
**Document ID:** `CS-AE-004`
## 1. Immutable Financial Ledger Rules
- `commissions` records the initial recognized commission upon conversion verification (`PENDING` -> `APPROVED` -> `SETTLED_PAID`).
- Any modification after initial creation (e.g., customer refund within 30–45 day lock period, tier slab upgrade bonus on Cloudways, or FX settlement variance) MUST be recorded as an append-only row in `commission_adjustments`.
- `reward_type` strictly separates `CASH_COMMISSION` (DigitalOcean, Cloudways, UpCloud, Contabo, Hostinger, Vultr, Kamatera, GCP) from `HOSTING_CREDIT` (Hetzner €10 referral credit).
