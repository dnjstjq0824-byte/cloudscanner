# COMMISSION_LEDGER.md — Double-Entry Commission Accounting & Payout Settlement
**Document ID:** `CS-AE-004`

> **Pivot addendum (2026-11-26):** Core TirpScanner monetization surface — kept; merchant examples reframed to SaaS vendors.

## 1. Immutable Financial Ledger Rules
- `commissions` records the initial recognized commission upon conversion verification (`PENDING` -> `APPROVED` -> `SETTLED_PAID`).
- Any modification after initial creation (e.g., customer refund within the network lock period, tier slab upgrade bonus, or FX settlement variance) MUST be recorded as an append-only row in `commission_adjustments`.
- `reward_type` strictly separates `CASH_COMMISSION` (standard SaaS affiliate payouts via Impact/CJ/PartnerStack) from `SERVICE_CREDIT` (vendors paying in account credit instead of cash).
