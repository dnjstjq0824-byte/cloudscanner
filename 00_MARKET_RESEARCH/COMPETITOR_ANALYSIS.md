# COMPETITOR_ANALYSIS.md — Global SaaS/AI-Tool Affiliate & Deal-Scanner Competitor Teardown

**Document ID:** `CS-MR-002` (stable ID, unchanged since CloudScanr era)  
**Original Verification Date:** `2026-10-28` (CloudScanr cloud-comparison cut — superseded, see §4)  
**Pivot Reframe Date:** `2026-11-26`  
**Analyzed Platforms:** `g2.com`, `capterra.com`, `trustradius.com`, `producthunt.com`, `appsumo.com`, `futuretools.io`, `theresanaiforthat.com`, `saashub.com`, `alternativeto.net`, `stacksocial.com`

> **Pivot Note (2026-11-26):** TirpScanner (`tirpscanr.com` — registered, wiring pending) is a global
> SaaS/AI-tool & software affiliate deal scanner. The 2026-10-28 CloudScanr
> cloud/VPS/GPU comparison cut below is superseded for live scope and preserved
> in §4 as history. See `TIRPSCANNER_PIVOT.md`.

---

## 1. Feature-by-Feature Competitor Matrix

| Analysis Dimension | `G2` / `Capterra` / `GetApp` | `TrustRadius` / `Peer Insights` | `Product Hunt` | `AppSumo` / `StackSocial` / `PitchGround` | `FutureTools` / `There's An AI For That` / `SaaSHub` / `AlternativeTo` | **TirpScanner (Target Architecture)** |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **Tool Coverage** | 150k+ SaaS products, review-led | Enterprise software, verified reviews | Daily launches (~50/day) | Curated lifetime deals (hundreds active) | 5k–15k AI tools / alternatives each | **Global SaaS/AI-tool & software catalog with plan-level specs** |
| **Product Categories** | 2k+ software categories | Enterprise segments | All tech launches | Deals only (no full catalog) | AI use-case tags / alternatives graph | **SaaS categories + AI use-case taxonomy + deal shelf** |
| **Deal / Price Tracking** | Static pricing pages, often stale | Static pricing summaries | None (launch-day promos only) | Deal price + discount % + countdown | None or static price mentions | **Deal scan: coupon/promo tracking, plan normalization, price-drop detection, alerts** |
| **Price Freshness** | Vendor-claimed, months stale | Vendor-claimed | N/A | Deal-window accurate | Rarely updated | **Scheduled deal-feed sync + explicit `FRESH`/`STALE`/`UNVERIFIED` badge + UTC timestamp** |
| **Plan Normalization** | Free text pricing blurbs | Feature lists | None | Single deal tier | None | **Normalized plans: monthly/annual per-seat, trial/freemium, renewal terms, FX + tax** |
| **Filter Capabilities** | Category + feature + rating filters | Segment + rating filters | Tag search only | Category + price filters | Tag/keyword search | **Sub-50ms OpenSearch facets + NL query parser (`ai video editor under $20/mo annual`)** |
| **Search UX/UI** | Review grids + comparison pages | Review reports | Feed + leaderboard | Deal cards + countdown | Directory cards | **Hybrid: Engineer Data-Grid + Guided Card Mode + Sticky Compare Tray (up to 6 plans)** |
| **Reviews / UGC** | Millions of verified reviews | Verified enterprise reviews | Community votes/comments | Buyer reviews + ratings | Community votes/lists | **Curated comparisons + rankings + buying guides; community ratings post-launch** |
| **Affiliate Model** | PPC pay-per-click + sponsored placements | Vendor-paid research | Promoted launches | Direct deal margin / rev-share | Outbound affiliate links | **Full Affiliate Engine (`/go/` router, signed `ClickID`/`SubID`, postback/CSV reconcile, commission ledger)** |
| **Ad Display Policy** | Sponsored slots labeled, mixed with organic | Vendor-sponsored reports | "Promoted" launches | Deals are the product | Minimal / none | **Strict organic vs sponsored separation (`is_sponsored` badge; ads never alter Deal Score)** |
| **Public API** | No (G2 data sold, not open) | No | GraphQL API (rate-limited) | No | No | **Documented OpenAPI 3.1 REST API (deal/catalog groups; cloud groups retired)** |
| **SEO Architecture** | Programmatic `/category/{x}` + `/compare/{a}-vs-{b}` | Review landing pages | `/leaderboard`, `/topics` | Deal landing pages | `/ai/{tool}`, `/alternative/{x}` pages | **Nuxt 3 SSR + scenario/vendor/plan pages + Schema.org `Product`/`Offer` + anti-thin-content guard** |
| **Global Localization** | English-first, partial locales | English | English | English + USD | English | **Multi-language (`en`, `ko`, `ja`, `de`, `es`) + multi-currency live FX + local VAT/GST toggle** |

---

## 2. Key Takeaways & Architectural Adoption for TirpScanner

1. **Adopt G2/Capterra Programmatic Category & Comparison SEO**:
   - Buyers search "best {category} software" and "{a} vs {b}". TirpScanner implements
     programmatic category hubs (`/ai-video-editors`), vendor pages (`/vendors/{slug}`),
     tool pages (`/tools/{slug}`), and head-to-head comparisons (`/compare/{a}-vs-{b}`)
     with live plan pricing instead of stale vendor blurbs.
2. **Adopt AppSumo Deal Mechanics Without the Marketplace Risk**:
   - Countdown urgency, discount depth (`-62%`), and stackable-coupon math convert.
     TirpScanner tracks third-party deals and coupons (never selling licenses itself),
     so there is no inventory, refund, or vendor-fulfillment risk — pure affiliate yield.
3. **Establish Unmatched Differentiation via Live Deal Scanning & Honest Plan Math**:
   - **Zero Competitors Offer Continuous Deal Scanning**: Directories list tools once and
     let pricing rot. TirpScanner re-scans vendor pricing pages, affiliate-network feeds,
     and coupon sources on a schedule, flagging price drops, new promos, and renewal hikes.
   - **Zero Competitors Separate Sticker Price from Effective Cost Transparently**:
     TirpScanner normalizes monthly vs annual per-seat math, trial/freemium limits, renewal
     vs intro pricing, and estimated tax/FX, with source URL and collection timestamp.

---

## 3. Clean-Room Intellectual Property Compliance

- TirpScanner **never** scrapes or copies proprietary review text from `G2`, `Capterra`,
  `TrustRadius`, or any competitor, and never copies frontend code from competitors.
- All plan/deal data is harvested from **official vendor pricing pages**, **official
  affiliate-network feeds/APIs**, and **authorized coupon sources** documented in
  `OFFICIAL_REFERENCE_LINKS.md`, with full provenance stored in `price_snapshots`.

---

## 4. Superseded 2026-10-28 CloudScanr Cut (History — Not Live Scope)

On 2026-10-28 this document compared cloud/VPS/GPU price-comparison platforms
(`instances.vantage.sh`, `vpsbenchmarks.com`, `rackrate.dev`, `cloudpriceindex.org`,
`costgoat.com`, `deploybase.ai`) across 18+ cloud providers, 12 infrastructure
categories, 36-field TCO normalization, and BYOC provisioning. That market
(cloud-infrastructure catalog, BYOC, provisioning, server management, resale) is
**RETIRED** under the 2026-11-26 pivot; the matrix above replaces it for all live
planning. The retired detail survives in frozen work-order history
(`11_WORK_ORDERS_0001_TO_1000/`, read-only).
