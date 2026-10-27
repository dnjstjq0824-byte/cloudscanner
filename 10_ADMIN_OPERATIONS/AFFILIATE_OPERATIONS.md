# AFFILIATE_OPERATIONS.md — Campaign Lifecycle, Link Health & Monthly Payout Reconciliation
**Document ID:** `CS-ADM-003`
## 1. Monthly Affiliate Settlement Workflow
1. On the 1st and 16th of each month, Affiliate Admin downloads settlement CSVs from Cloudways, Kamatera, and Vultr, and verifies automated API syncs from Impact (DigitalOcean, GCP) and CJ (Contabo, IONOS).
2. Admin uploads CSVs to `/admin/affiliate/reconciliation`.
3. System reconciles `PENDING` conversions to `APPROVED` or records `CLAWBACK_REFUND` adjustments in `commission_adjustments`.
4. Upon bank/PayPal receipt, Admin records the `affiliate_payouts` batch reference and transitions matched `commissions` to `SETTLED_PAID`.
