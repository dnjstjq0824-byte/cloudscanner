# COST_MANAGEMENT.md — CloudScanr Own Infrastructure FinOps & Test Account Budget Controls
**Document ID:** `CS-DO-004`
## 1. Sandbox Live Testing Budget Cap
- Any live provider sandbox test account used in Phase 9 QA enforces a hard **$5.00/day** and **$25.00/month** ceiling, with an automated janitor cron (`cmd/janitor`) that terminates any test instance tagged `cloudscanr-ci-ephemeral=true` older than 10 minutes.
