# AFFILIATE_OPERATIONS.md — Campaign Lifecycle, Link Health & Monthly Payout Reconciliation
**Document ID:** `CS-ADM-003`
> TirpScanner pivot (2026-11-26): LIVE scope — settlement workflow kept; merchant set reframed from cloud hosts to SaaS/AI-tool vendors and affiliate networks. See `TIRPSCANNER_PIVOT.md`.
## 1. Monthly Affiliate Settlement Workflow
1. On the 1st and 16th of each month, Affiliate Admin downloads settlement CSVs from direct SaaS vendor programs and verifies automated API syncs from affiliate networks (Impact, CJ, PartnerStack, Awin).
2. Admin uploads CSVs to `/admin/affiliate/reconciliation`.
3. System reconciles `PENDING` conversions to `APPROVED` or records `CLAWBACK_REFUND` adjustments in `commission_adjustments`.
4. Upon bank/PayPal receipt, Admin records the `affiliate_payouts` batch reference and transitions matched `commissions` to `SETTLED_PAID`.
