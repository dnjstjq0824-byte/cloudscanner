# MARKET_OVERVIEW.md — Global Cloud, VPS & GPU Infrastructure Market Due Diligence

**Document ID:** `CS-MR-001`  
**Verification Date:** `2026-10-28`  
**Scope:** Global Cloud IaaS, VPS, Dedicated CPU, GPU Cloud, Bare Metal, Managed Hosting, and Affiliate Monetization Ecosystem

---

## 1. Executive Market Synthesis

The global cloud compute, Virtual Private Server (VPS), and GPU infrastructure market exhibits extreme pricing fragmentation and opaque billing structures across three distinct provider tiers:
1. **Developer-Focused Cloud & VPS Providers (Tier 1 Core)**: DigitalOcean, UpCloud, Vultr, Hetzner, Akamai Cloud (Linode), Contabo, Scaleway, OVHcloud, Kamatera, Hostinger, and IONOS. These providers offer predictable monthly caps or flat hourly rates, generous bundled bandwidth (1 TB to 32 TB), and accessible REST APIs.
2. **Managed Cloud Hosting Platforms**: Cloudways (by DigitalOcean), Kinsta, WP Engine. They layer managed platform orchestration on top of raw IaaS (DigitalOcean, AWS, GCP, Vultr, Linode) and offer high-paying affiliate commissions ($50–$125 CPA or $30 + 7% lifetime recurring).
3. **Hyperscalers & AI GPU Clouds**: AWS, Google Cloud Platform (GCP), Microsoft Azure, Oracle Cloud Infrastructure (OCI), Cloudflare, and specialized GPU providers (RunPod). Hyperscalers bill modularly across dozens of micro-dimensions (vCPU hours, GiB-month storage, provisioned IOPS, public IPv4 hourly fees of $0.005/hr, and steep egress rates of $0.08–$0.12/GB, except OCI which includes 10 TB/month free egress).

### Key Market Pain Points Addressed by CloudScanr
- **Hidden Egress & IPv4 Traps**: Many users compare base vCPU/RAM prices without realizing that public IPv4 addresses cost $3.65/month extra on hyperscalers (or €0.50–€0.60/month difference on Hetzner IPv4 vs IPv6), or that 2 TB of outbound traffic costs $0 on DigitalOcean/UpCloud/Hetzner vs. $180+ on AWS/GCP/Azure.
- **CPU Steal & Performance Variance**: A "2 vCPU / 4 GB RAM" instance can mean throttled burstable shared cores (`t3.medium`), regular shared cores, or high-frequency dedicated AMD EPYC Genoa / Intel Emerald Rapids cores.
- **Recent Market Price Shifts (Verified 2026)**: Providers regularly adjust regional and instance pricing (e.g., Hetzner Cloud adjusted CX/CAX/CPX/CCX pricing structures in June 2026 and applies regional differentials across EU vs. US vs. Singapore; DigitalOcean introduced Agent Droplets and expanded GPU Droplets). Static blog comparisons become obsolete within weeks.
- **Disconnect Between Discovery and Day-2 Operations**: Existing price comparison tools (`instances.vantage.sh`, `vpsbenchmarks.com`) stop at discovery. Once a developer chooses a server, they must manage 3–5 separate cloud consoles. CloudScanr bridges discovery with **BYOC (Bring Your Own Cloud)** multi-cloud provisioning and management.

---

## 2. Provider Integration & Monetization Tier Classification (Article 3 & 5)

CloudScanr never assumes uniform capability across all 18+ providers. Each provider is classified across four independent axes (**Pricing API**, **Affiliate Program**, **Management API (BYOC)**, and **Reseller Eligibility**):

| Provider ID | Provider Name | Pricing Ingestion Method | Affiliate Monetization Status | BYOC Management API | V1 Integration Level |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `digitalocean` | **DigitalOcean** | `VERIFIED` (`GET /v2/sizes` via OAuth/PAT) | `VERIFIED` (10% RevShare 12mo via Impact) | `VERIFIED` (Droplets API v2 + OAuth2/PAT) | **Full Tier-1 (Catalog + Affiliate + BYOC Provisioning)** |
| `upcloud` | **UpCloud** | `VERIFIED` (`GET /1.3/price` & `/1.3/plan`) | `VERIFIED` (Up to €300 CPA via PartnerStack) | `VERIFIED` (REST API 1.3 `/server`) | **Full Tier-1 (Catalog + Affiliate + BYOC Provisioning)** |
| `cloudways` | **Cloudways** | `VERIFIED` (API v2 `GET /api/v2/server_sizes`) | `VERIFIED` (Slab up to $125 or Hybrid $30+7%) | `VERIFIED` (API v2 OAuth `/api/v2/server`) | **Full Tier-1 (Catalog + Affiliate + BYOC Provisioning)** |
| `vultr` | **Vultr** | `VERIFIED` (Public `GET /v2/plans` & `/v2/plans-metal`) | `VERIFIED` (Referral $35–$100 CPA) | `VERIFIED` (REST API v2 `/v2/instances`) | **Full Tier-1 (Catalog + Affiliate + BYOC Provisioning)** |
| `hetzner` | **Hetzner Cloud** | `VERIFIED` (`GET /v1/pricing` & `/v1/server_types`) | `VERIFIED` (Customer Referral €10 Cloud Credit only — **No Cash Affiliate**) | `VERIFIED` (Cloud API v1 `/v1/servers`) | **Full Tier-1 (Catalog + Credit Referral Disclosure + BYOC)** |
| `akamai_linode` | **Akamai (Linode)** | `VERIFIED` (Public `GET /v4/linode/types`) | `REQUIRES_APPLICATION` (Referral Credit / Impact Partner) | `VERIFIED` (Linode API v4 `/v4/linode/instances`) | **Full Tier-1 (Catalog + Referral/Affiliate + BYOC)** |
| `gcp` | **Google Cloud** | `VERIFIED` (Cloud Billing Catalog API v1) | `VERIFIED` (Google Cloud Affiliate via Impact) | `VERIFIED` (Compute Engine REST API v1) | **Tier-1 Catalog & Affiliate + Read/Scoped BYOC** |
| `aws` | **Amazon Web Services** | `VERIFIED` (AWS Bulk Price List API JSON) | `CONTRACT_REQUIRED` (APN B2B Partner Only) | `VERIFIED` (EC2 API via IAM STS `AssumeRole`) | **Tier-1 Catalog + BYOC Read/Manage (Affiliate = APN Only)** |
| `azure` | **Microsoft Azure** | `VERIFIED` (Public `GET https://prices.azure.com/api/retail/prices`) | `CONTRACT_REQUIRED` (Microsoft AI Cloud Partner CSP) | `VERIFIED` (Azure ARM Compute REST API) | **Tier-2 Catalog + BYOC Read/Manage** |
| `contabo` | **Contabo** | `DOCUMENTATION_UNAVAILABLE` (Curated Feed + Admin Sync) | `VERIFIED` (€20–€250 CPA via CJ Affiliate) | `VERIFIED` (Contabo REST API v1 `/v1/compute/instances`) | **Tier-2 Catalog + CJ Affiliate + BYOC** |
| `hostinger` | **Hostinger** | `REQUIRES_APPLICATION` (Curated Feed + VPS API) | `VERIFIED` (60%+ CPA via In-House/CJ) | `VERIFIED` (Hostinger VPS API v1) | **Tier-2 Catalog + Affiliate + BYOC** |
| `ovhcloud` | **OVHcloud** | `VERIFIED` (Public `GET /1.0/order/catalog/public/cloud`) | `REQUIRES_APPLICATION` (Affiliate / Partner Referral Path) | `VERIFIED` (OVH API v1 `/cloud/project/{id}/instance`) | **Tier-2 Catalog + Affiliate + BYOC** |
| `scaleway` | **Scaleway** | `VERIFIED` (`GET /instance/v1/zones/{z}/products/servers`) | `CONTRACT_REQUIRED` (B2B Partner Program) | `VERIFIED` (Scaleway Instance API v1) | **Tier-2 Catalog + BYOC** |
| `kamatera` | **Kamatera** | `VERIFIED` (Cloud API `/service/server` options) | `VERIFIED` ($50–$150 CPA Direct Affiliate) | `VERIFIED` (Kamatera Cloud API) | **Tier-2 Catalog + Affiliate + BYOC** |
| `ionos` | **IONOS Cloud** | `DOCUMENTATION_UNAVAILABLE` (Curated Rate Card) | `VERIFIED` (CPA via Awin / CJ) | `VERIFIED` (IONOS Cloud API v6) | **Tier-2 Catalog + Affiliate + BYOC** |
| `oracle_cloud` | **Oracle Cloud (OCI)** | `VERIFIED` (Public Cost Estimator API `/cetools/api/v1/products/`) | `NOT_SUPPORTED` (OPN Enterprise Partner Only) | `VERIFIED` (OCI Core Services API) | **Tier-2 Catalog + BYOC Read** |
| `cloudflare` | **Cloudflare** | `DOCUMENTATION_UNAVAILABLE` (Formula Rate Cards: R2/Workers) | `CONTRACT_REQUIRED` (Partner Network Only) | `VERIFIED` (Cloudflare API v4) | **Tier-2 Edge/Storage Catalog** |
| `runpod` | **RunPod** | `VERIFIED` (GraphQL `gpuTypes` query) | `VERIFIED` (5% Referral Program) | `VERIFIED` (GraphQL Pod Deployment API) | **Tier-2 GPU Catalog + Referral + BYOC** |

---

## 3. Individual Developer & Korean Resident Affiliate Eligibility (Article 4 & 11)

A critical operational requirement for CloudScanr is verifying what an **individual developer residing in South Korea (대한민국 거주 개인 개발자)** can legally and practically join without fabricating business credentials:

1. **Directly Eligible for Korean Individual Developers (`VERIFIED` / `REQUIRES_APPLICATION`)**:
   - **Cloudways**: Open to international individuals; requires a live content website (`cloudscanr.com`), email verification, and account manager review. Payouts via PayPal or Bank Wire ($250 threshold). Offers choice of **Slab** ($50–$125/sale) or **Hybrid** ($30 + 7% lifetime recurring).
   - **DigitalOcean**: Hosted on **Impact Radius**. Individual publishers in South Korea can register on Impact, submit a live website (`cloudscanr.com`), file an electronic **IRS Form W-8BEN** (claiming US–Korea tax treaty status for non-US service/marketing income), and receive payouts in KRW/USD via PayPal or local bank wire ($100 minimum). Earns **10% of referred customer spend for 12 months**.
   - **UpCloud**: Hosted on **PartnerStack / Direct**. Accepts individual content creators and comparison platforms subject to application review. Payouts via PayPal or Stripe.
   - **Contabo**: Managed on **CJ Affiliate (Commission Junction)**. Korean individuals can open a CJ Publisher account free of charge, complete the non-US W-8BEN tax interview, link a Payoneer USD/EUR receiving account, and apply to the Contabo program (€20–€250 CPA, 45-day cookie).
   - **Hostinger**: Accepts individual publishers globally with an active website or tech channel. Minimum 60% CPA on initial hosting/VPS packages, $100 minimum payout via PayPal.
   - **Vultr**: Any verified customer can access the Customer Referral Program ($35–$100 CPA, paid via PayPal on the 1st and 15th of the month, $100 minimum). Higher-volume media affiliate tiers require application.
   - **Google Cloud**: Available via Impact Radius (`https://cloud.google.com/affiliate-program`) offering tiered CPA and a $350 trial credit link for new eligible billing accounts.

2. **Credit-Only Referral Programs (Must Not Be Misclassified as Cash Revenue)**:
   - **Hetzner Cloud**: Hetzner does **not** run a cash-payout affiliate network on CJ or Impact. It offers a **Customer Referral Program** (`Give €20, Get €10 Hetzner Cloud Credit`) strictly available to existing customers who have paid at least 3 monthly invoices. CloudScanr models Hetzner's reward type as `HOSTING_CREDIT` in the `commission_type` enum so financial reporting never inflates cash revenue with hosting credits.

3. **B2B Entity / Partner Contract Required (`CONTRACT_REQUIRED` / `NOT_SUPPORTED` for Individuals)**:
   - **AWS (APN)**, **Microsoft Azure (AI Cloud Partner / CSP)**, **Oracle Cloud (OPN)**, **Scaleway Partner Program**, and **Cloudflare Partner Network** require a registered corporate entity (사업자등록/법인), commercial partner agreements, and often annual membership or certification commitments.
   - **Operational Guidance for Korean Operator**: While CloudScanr can launch Phase 1–3 with individual affiliate accounts for DigitalOcean, Cloudways, UpCloud, Contabo, Hostinger, Vultr, Kamatera, and GCP, registering a Korean Sole Proprietorship (**개인사업자 — 전자상거래/정보통신업/광고대행업**) is strongly recommended before scaling commercial payouts, issuing foreign currency tax invoices (영세율 외화입금증명), and applying for B2B Partner programs.
