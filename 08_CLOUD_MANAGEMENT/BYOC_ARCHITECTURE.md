> **RETIRED SCOPE — TirpScanner pivot (2026-11-26):** BYOC account connection, server provisioning, and cloud server management were removed from live scope. This file is CloudScanr-era history, preserved as record; do not build from it. Live scope (SaaS/AI-tool catalog, deal scan / price-drop detection, affiliate monetization) is defined in `TIRPSCANNER_PIVOT.md` at the repository root.

# BYOC_ARCHITECTURE.md — Bring-Your-Own-Cloud Multi-Tenant Architecture (Article 12)
**Document ID:** `CS-CM-001`
## 1. Core BYOC Principles
- CloudScanr acts strictly as an **authorized control plane delegate** on the user's own cloud account.
- The user retains 100% legal ownership and direct billing relationship with the cloud provider.
- Authentication Preference Order:
  1. **Official OAuth 2.0** (DigitalOcean OAuth, Linode OAuth, GCP/Azure OAuth/OIDC) with granular `read` or `read+write` scopes.
  2. **IAM Cross-Account Role Assumption** (AWS STS `AssumeRole` with unique per-tenant `ExternalId`).
  3. **Scoped API Token / Sub-Account Key** (Hetzner Project API Token, UpCloud API Sub-account with IP restriction, Vultr API Key with IP ACL, Cloudways API v2 OAuth key).
