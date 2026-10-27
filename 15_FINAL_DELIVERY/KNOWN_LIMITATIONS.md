# KNOWN_LIMITATIONS.md — Transparent Empirical Verification Boundaries & External Dependencies (Articles 4, 27, 28, 29)

**Document ID:** `CS-FD-002`  
**Verification Date:** `2026-10-28`

---

## 1. Strict Honesty & Non-Fabrication Disclosures (Article 29 Compliance)

In strict compliance with Articles 4, 11, 24, 27, and 29, CloudScanr explicitly separates **verified specification/architecture facts** from **items requiring external merchant approval, corporate legal contracts, or live funded cloud accounts**:

1. **External Affiliate Network Publisher Approvals (`REQUIRES_APPLICATION`)**:
   - While official affiliate programs for **DigitalOcean (Impact)**, **UpCloud (PartnerStack)**, **Cloudways**, **Vultr**, **Contabo (CJ)**, **Hostinger**, **Kamatera**, **IONOS (Awin/CJ)**, and **Google Cloud (Impact)** are empirically `VERIFIED` to exist and accept publisher applications, **CloudScanr does NOT falsely claim that live publisher accounts have already been approved**.
   - Actual affiliate commission generation requires the platform operator to deploy the Phase 2 MVP on `https://cloudscanr.com`, submit publisher applications + IRS Form W-8BEN, and pass each merchant's human content review (2–10 business days). Until approved, `affiliate_accounts.account_status` remains `PENDING_APPROVAL` and `/go/{slug}` safely routes to `fallback_clean_url`.

2. **Providers Without Cash Affiliate Programs (`HOSTING_CREDIT` / `CONTRACT_REQUIRED` / `NOT_SUPPORTED`)**:
   - **Hetzner Cloud**: Offers a Customer Referral Program (`€10 Hetzner Cloud Credit` per qualified referral after the referrer has paid >=3 monthly invoices), **not** a cash-payout affiliate network. Modeled strictly as `reward_type = HOSTING_CREDIT`.
   - **AWS (APN)**, **Microsoft Azure (CSP)**, **Scaleway Partner**, **Cloudflare Partner**, and **Oracle Cloud (OPN)**: Do **not** offer open individual cash affiliate programs; they require B2B corporate entity registration and formal Partner/Reseller contracts (`CONTRACT_REQUIRED` / `NOT_SUPPORTED`).

3. **Providers Without Public Unauthenticated Pricing APIs (`DOCUMENTATION_UNAVAILABLE` / `REQUIRES_APPLICATION`)**:
   - **Contabo**, **IONOS**, and **Cloudflare**: Provide authenticated compute/edge management APIs, but do not expose a single unauthenticated public SKU price-book endpoint equivalent to Vultr (`GET /v2/plans`) or Linode (`GET /v4/linode/types`). For these providers, CloudScanr uses the documented **Curated Rate-Card + Checksum Verification Pipeline (`WO-CS-0245`)** rather than inventing non-existent pricing endpoints.

4. **Live External Billable Server Provisioning Tests (`BLOCKED_EXTERNAL_APPROVAL` Until Live Sandbox Authorized)**:
   - All 22 mandatory server provisioning safety scenarios (Article 24, `SAF-01` to `SAF-22`) are fully specified and designed to run deterministically against the local/CI Provider Fault Simulator (`internal/testutil/providersim`).
   - Executing live billable VM creation/deletion against production DigitalOcean, UpCloud, Vultr, or Hetzner accounts (`Tier C Live Sandbox Testing`) is gated behind explicit budget authorization (`ALLOW_LIVE_BILLABLE_TESTS=true`) to prevent unauthorized cloud charges.
