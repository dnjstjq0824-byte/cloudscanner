# IAM.md — Role-Based & Attribute-Based Access Control Specification
**Document ID:** `CS-SEC-001`
## 1. Organization Permission Matrix
| Action | `ORG_OWNER` | `ORG_ADMIN` | `CLOUD_OPERATOR` | `BUDGET_APPROVER` | `ORG_VIEWER` |
| :--- | :--- | :--- | :--- | :--- | :--- |
| View Connected Servers & Costs | ALLOW | ALLOW | ALLOW | ALLOW | ALLOW |
| Connect / Disconnect BYOC Account | ALLOW | ALLOW (MFA) | DENY | DENY | DENY |
| Create / Reboot / Resize Server (Within Budget) | ALLOW (MFA) | ALLOW (MFA) | ALLOW (MFA) | DENY | DENY |
| Approve Provisioning Over Standard Cap | ALLOW (MFA) | DENY | DENY | ALLOW (MFA) | DENY |
| Delete / Rebuild Server | ALLOW (MFA) | ALLOW (MFA) | REQUIRES_APPROVAL | DENY | DENY |
