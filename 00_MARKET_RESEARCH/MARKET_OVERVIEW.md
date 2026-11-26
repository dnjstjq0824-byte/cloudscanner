# MARKET_OVERVIEW.md — Global SaaS, AI-Tool & Software Affiliate Market Due Diligence

**Document ID:** `CS-MR-001` (stable ID, unchanged since CloudScanr era)  
**Original Verification Date:** `2026-10-28` (CloudScanr cloud-IaaS cut — superseded, see §4)  
**Pivot Reframe Date:** `2026-11-26`  
**Scope:** Global SaaS & AI-tool discovery, software deal/coupon tracking, affiliate monetization ecosystem, and content-driven buyer acquisition

> **Pivot Note (2026-11-26):** TirpScanner (`tirpscanr.com` — registered, wiring pending) is a global
> SaaS/AI-tool & software affiliate deal scanner. The 2026-10-28 CloudScanr
> cloud/VPS/GPU market cut below is superseded for live scope and preserved in
> §4 as history. See `TIRPSCANNER_PIVOT.md`.

---

## 1. Executive Market Synthesis

The global SaaS and AI-tool market exhibits extreme pricing fragmentation and opaque promotional structures across three distinct monetization tiers:

1. **Productivity & Business SaaS Vendors**: Project management, CRM, marketing automation, design, video, developer tools, and security suites. Most sell per-seat monthly/annual plans with intro-vs-renewal gaps, add-on upsells, and partner referral programs paying 20–30% recurring commissions (commonly 12 months or lifetime) via PartnerStack, FirstPromoter, Rewardful, Impact, or in-house portals.
2. **AI-Tool Vendors (Highest Growth)**: AI writers, image/video generators, voice, coding assistants, and AI search. Pricing mixes subscriptions with credit/token packs; promos, lifetime deals (AppSumo, PitchGround), and launch-day coupons create constant price movement that static directories never track.
3. **Software Deal & Coupon Channels**: AppSumo, StackSocial, affiliate-network coupon feeds, vendor promo pages, and seasonal sales (Black Friday, back-to-school). Discounts of 20–30% appear and expire in days; buyers lack a single scanner that verifies which codes still work and what the renewal price will be.

### Key Market Pain Points Addressed by TirpScanner

- **Intro-vs-Renewal Traps**: Many tools advertise `$12/mo` that is actually `$12/mo billed annually for year one, renewing at $29/mo`. Buyers compare sticker prices without seeing renewal terms, per-seat minimums, or overage/credit pricing.
- **Subscription Sprawl & Stale Comparisons**: Teams hold 10–30+ paid seats across vendors. Static blog "best of" lists go stale within weeks as plans, limits, and promos change (verified pattern through 2026: vendors routinely re-tier AI plans and credit bundles).
- **Dead Coupons & Unverified Deals**: Coupon aggregators list expired codes with no verification timestamp. TirpScanner re-scans deal sources on a schedule and labels every deal `FRESH`/`STALE`/`UNVERIFIED` with a UTC timestamp.
- **Disconnect Between Discovery and Deal Timing**: Existing directories (`FutureTools`, `SaaSHub`, `AlternativeTo`) stop at listing. Buyers want to know *when* to buy. TirpScanner bridges discovery with **deal scanning**: price-drop detection, promo tracking, and alerts.

---

## 2. Vendor & Network Integration Tier Classification

TirpScanner never assumes uniform capability across vendors. Each vendor/network is classified across three independent axes (**Deal-Feed Ingestion**, **Affiliate Program**, and **Content Eligibility**):

| Vendor / Network ID | Vendor / Network Name | Deal-Feed Ingestion Method | Affiliate Monetization Status | Content / Listing Eligibility | V1 Integration Level |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `partnerstack` | **PartnerStack (network)** | `VERIFIED` (REST API: programs, coupons, conversions) | `VERIFIED` (20–30% recurring typical; program-dependent) | `VERIFIED` (marketplace listings + partner content) | **Full Tier-1 (Feed + Affiliate + Content)** |
| `impact` | **Impact Radius (network)** | `VERIFIED` (REST API: actions, promo codes) | `VERIFIED` (CPA / RevShare, advertiser-dependent) | `VERIFIED` (publisher content allowed) | **Full Tier-1 (Feed + Affiliate + Content)** |
| `cj_affiliate` | **CJ Affiliate (network)** | `VERIFIED` (REST API: links, coupons, transactions) | `VERIFIED` (CPA / RevShare, advertiser-dependent) | `VERIFIED` (publisher content allowed) | **Full Tier-1 (Feed + Affiliate + Content)** |
| `awin` | **Awin (network)** | `VERIFIED` (REST API: programmes, promotions, transactions) | `VERIFIED` (CPA / RevShare, advertiser-dependent) | `VERIFIED` (publisher content allowed) | **Full Tier-1 (Feed + Affiliate + Content)** |
| `firstpromoter` | **FirstPromoter (SaaS referrals)** | `REQUIRES_APPLICATION` (per-vendor API keys) | `VERIFIED` (SaaS recurring, vendor-set rates) | `VERIFIED` | **Tier-1 Feed + Affiliate** |
| `rewardful` | **Rewardful (Stripe affiliates)** | `REQUIRES_APPLICATION` (per-vendor API keys) | `VERIFIED` (Stripe-linked recurring) | `VERIFIED` | **Tier-1 Feed + Affiliate** |
| `shareasale` | **ShareASale (network)** | `VERIFIED` (API + datafeeds) | `VERIFIED` (merchant-set CPA/RevShare) | `VERIFIED` | **Tier-2 Feed + Affiliate** |
| `appsumo` | **AppSumo (deal marketplace)** | `DOCUMENTATION_UNAVAILABLE` (curated deal curation + partner feed on approval) | `REQUIRES_APPLICATION` (partner/affiliate application) | `VERIFIED` (deal reviews + roundups) | **Tier-2 Curated Deals + Content** |
| `stacksocial` | **StackSocial (deal marketplace)** | `DOCUMENTATION_UNAVAILABLE` (curated feed) | `REQUIRES_APPLICATION` (affiliate application) | `VERIFIED` | **Tier-2 Curated Deals + Content** |
| `g2` | **G2 (reviews/discovery)** | `NOT_SUPPORTED` (no public pricing feed; data licensed) | `NOT_SUPPORTED` (PPC vendor model, not publisher affiliate) | `VERIFIED` (outbound citation only; no review copying) | **Discovery Reference Only** |
| `capterra` | **Capterra (reviews/discovery)** | `NOT_SUPPORTED` (no public feed) | `NOT_SUPPORTED` (PPC vendor model) | `VERIFIED` (outbound citation only) | **Discovery Reference Only** |
| `producthunt` | **Product Hunt (launches)** | `VERIFIED` (GraphQL API, rate-limited) | `NOT_SUPPORTED` | `VERIFIED` (launch tracking + outbound links) | **Tier-2 Launch Feed** |

---

## 3. Individual Developer & Korean Resident Affiliate Eligibility

A critical operational requirement for TirpScanner is verifying what an **individual developer residing in South Korea** can legally and practically join without fabricating business credentials:

1. **Directly Eligible for Korean Individual Developers (`VERIFIED` / `REQUIRES_APPLICATION`)**:
   - **PartnerStack**: Open to international individual publishers/affiliates with a live content website (`tirpscanr.com` — registered, wiring pending); per-program approval by each SaaS vendor. Payouts via PayPal or Stripe (thresholds vary by program). Typical SaaS offers: 20–30% recurring for 12 months or lifetime.
   - **Impact Radius**: Individual publishers in South Korea can register on Impact, submit a live website, file an electronic **IRS Form W-8BEN** (claiming US–Korea tax treaty status for non-US marketing income), and receive payouts via PayPal or bank transfer (minimums vary by advertiser).
   - **CJ Affiliate**: Korean individuals can open a CJ Publisher account free of charge, complete the non-US W-8BEN tax interview, link a Payoneer USD/EUR receiving account, and apply to advertiser programs (CPA/RevShare, cookie windows typically 7–45 days).
   - **Awin**: Individual publishers accepted globally subject to a small refundable signup deposit and website review; payouts to bank/PayPal.
   - **FirstPromoter / Rewardful**: Per-vendor referral programs; individuals join each SaaS vendor's program directly (no separate network entity needed). Payouts via PayPal/Stripe per vendor terms.
   - **ShareASale**: Individual publishers accepted with website review and W-8BEN; payouts via direct deposit/Payoneer.

2. **Credit-Only or Non-Cash Rewards (Must Not Be Misclassified as Cash Revenue)**:
   - Some SaaS vendors pay referral rewards as **account credits** (free months, workspace credit) rather than cash. TirpScanner models such rewards as `reward_type = HOSTING_CREDIT` (enum name historical; covers SaaS account credit) so financial reporting never inflates cash revenue with non-cash credit.

3. **Entity / Contract-Gated Programs (`CONTRACT_REQUIRED` / `NOT_SUPPORTED` for Individuals)**:
   - Enterprise reseller tiers, agency-white-label programs, and managed-partner levels at major SaaS vendors require a registered corporate entity, commercial agreements, and sometimes revenue commitments.
   - **Operational Guidance for Korean Operator**: TirpScanner can launch with individual affiliate accounts across the Tier-1/Tier-2 networks above, but registering a Korean sole proprietorship is strongly recommended before scaling commercial payouts, issuing foreign-currency receipts, and applying for contract-gated partner tiers.

---

## 4. Superseded 2026-10-28 CloudScanr Cut (History — Not Live Scope)

On 2026-10-28 this document covered the cloud IaaS/VPS/GPU market: 18+ cloud
providers (DigitalOcean, UpCloud, Vultr, Hetzner, Linode, hyperscalers, GPU clouds),
four-axis classification including BYOC Management API and reseller eligibility,
and cloud-provider affiliate programs (Cloudways slab/hybrid, DigitalOcean Impact
RevShare, Contabo CJ CPA, Hetzner credit-only referrals). That market is **RETIRED**
under the 2026-11-26 pivot. The sections above replace it for all live planning;
retired detail survives in frozen work-order history (`11_WORK_ORDERS_0001_TO_1000/`,
read-only).
