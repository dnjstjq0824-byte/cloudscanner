# BUDGET_CONTROLS.md — Pre-Provisioning Spend Guardrails & Burn-Rate Alerts
**Document ID:** `CS-CM-005`
## 1. Pre-Execution Budget Enforcement
Before issuing an `operation_approvals` token for `CREATE_SERVER` or `RESIZE_SERVER`, the Budget Engine evaluates:
1. `NewServerMonthlyCost <= budget_policy.max_single_server_monthly_usd`
2. `CurrentOrgMonthlyRunRate + NewServerMonthlyCost <= budget_policy.monthly_limit_usd`
If exceeded and `hard_block_provisioning_on_exceed = true`, the API returns `HTTP 402 Payment Required` (`BUDGET_POLICY_EXCEEDED`) unless an `ORG_OWNER` or `BUDGET_APPROVER` signs an explicit budget override.
