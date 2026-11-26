# IAM.md — Role-Based & Attribute-Based Access Control Specification
**Document ID:** `CS-SEC-001`
> TirpScanner pivot (2026-11-26): BYOC/server actions retired; matrix reframed for catalog, deal, affiliate, and payout operations. See `TIRPSCANNER_PIVOT.md`.
## 1. Organization Permission Matrix
| Action | `ORG_OWNER` | `ORG_ADMIN` | `CATALOG_OPERATOR` | `PAYOUT_APPROVER` | `ORG_VIEWER` |
| :--- | :--- | :--- | :--- | :--- | :--- |
| View Catalog, Deals & Commission Reports | ALLOW | ALLOW | ALLOW | ALLOW | ALLOW |
| Create / Edit / Retire Catalog Listings & Deals | ALLOW (MFA) | ALLOW (MFA) | ALLOW (MFA) | DENY | DENY |
| Manage Affiliate Campaigns & Reconciliation Uploads | ALLOW (MFA) | ALLOW (MFA) | ALLOW (MFA) | DENY | DENY |
| Approve Payout Batch Over Standard Cap | ALLOW (MFA) | DENY | DENY | ALLOW (MFA) | DENY |
| Moderate Content / Delete Listing | ALLOW (MFA) | ALLOW (MFA) | REQUIRES_APPROVAL | DENY | DENY |
