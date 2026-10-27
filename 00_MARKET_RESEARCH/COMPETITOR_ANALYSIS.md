# COMPETITOR_ANALYSIS.md — Global Cloud Price Comparison Competitor Teardown

**Document ID:** `CS-MR-002`  
**Verification Date:** `2026-10-28`  
**Analyzed Platforms:** `instances.vantage.sh`, `vpsbenchmarks.com`, `rackrate.dev`, `cloudpriceindex.org`, `costgoat.com`, `deploybase.ai`

---

## 1. Feature-by-Feature Competitor Matrix (Article 3)

| Analysis Dimension | `instances.vantage.sh` | `vpsbenchmarks.com` | `rackrate.dev` | `cloudpriceindex.org` | `costgoat.com` / `deploybase.ai` | **CloudScanr (Target Architecture)** |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **Supported Providers** | AWS, Azure, GCP, RDS, ElastiCache | 25+ VPS & Cloud hosts (DO, Vultr, Hetzner, Linode, OVH, AWS, Oracle) | Developer VPS & Bare Metal providers | Hyperscalers + Major Regional Cloud IaaS | Hetzner, Scaleway, OVH, GPU Clouds | **18+ Global Providers (Tier-1 VPS, Managed Hosting, Hyperscalers, GPU)** |
| **Product Categories** | Cloud VMs, Managed DBs, Cache | VPS, Cloud VMs | VPS, Dedicated, Bare Metal | Compute VMs, GPU, Storage | Cloud VPS, GPU Instances, LLM Tokens | **16 Categories (VPS, Shared/Ded VM, GPU, Bare Metal, DB, K8s, Storage, Hosting)** |
| **Filter Capabilities** | Fast client-side DataTables grid (vCPU, Memory, Network, Arch, Region) | Price brackets ($8–$120), CPU/IO/Web score filters, Feature toggles | Spec sliders (CPU, RAM, Disk, Bandwidth), Location filters | Region, Currency, Spec tier index filters | Calculator sliders per provider | **Sub-50ms OpenSearch Facets + Natural Language Spec Parser (`4 vcpu 16gb nvme tokyo under $40`)** |
| **Search UX/UI** | Dense single-screen engineering table; zero fluff; instant column sort | Multi-tab scorecards, scatter plots (Price vs Score), side-by-side compare | Clean dark-mode developer cards + tabular toggle | Index charts & regional heatmaps | Interactive cost breakdown calculators | **Hybrid Mode: Engineer Data-Grid Mode + Guided Card Mode + Sticky Compare Tray (up to 6 SKUs)** |
| **Price Freshness** | Automated daily/weekly scrape of AWS/Azure/GCP JSON APIs | Updated during monthly trial runs & periodic crawls | Periodic catalog polling | Periodic index snapshots | Weekly/Monthly verified snapshots | **Hourly/6-Hour Official API Sync + Explicit `FRESH`/`STALE`/`UNVERIFIED` Badge + UTC Timestamp** |
| **Price Normalization** | Hourly, Monthly (730h), Reserved (1yr/3yr), Spot | Monthly base price + trial grade normalization | Base monthly/hourly rate | Standardized basket index | Includes bandwidth overage & IPv4 toggles | **36-Field TCO Engine (Raw Rate vs 730h/Calendar Month + Egress + IPv4 + Backup + Tax + FX)** |
| **Benchmarks Source** | Official provider ECU/vCPU specs only | Proprietary automated 25-metric multi-day live trials (Sysbench, FIO, wrk) | Spec-derived + community YABS | None (Price focus) | Provider specs + GPU TFLOPS/VRAM bandwidth | **Provenance-Tagged Benchmarks (Official Spec vs Verified YABS/Sysbench Run with Raw Log Hash)** |
| **Affiliate Model** | Lead-gen funnel to Vantage SaaS FinOps platform | Outbound affiliate links (`?ref=...`) on "Visit Host" CTAs + Trial promos | Outbound affiliate referral links | Research/Ad supported | Outbound referral links & calculators | **Full Affiliate Engine (`/go/` router, Signed `ClickID`/`SubID`, Postback/CSV Reconcile, Commission Ledger)** |
| **Ad Display Policy** | First-party Vantage SaaS banner only | Minimal native sponsorship; strict separation from trial scores | Clean / Minimal | Informational | Minimal | **Strict Organic vs Sponsored Separation (`is_sponsored` badge; ads never alter Value Score)** |
| **Live Server Management** | **None** | **None** | **None** | **None** | **None** | **Full BYOC Multi-Cloud Connection, Cost-Approved Server Provisioning & Day-2 Lifecycle Console** |
| **Public API** | Open-source GitHub JSON build artifacts | Private (AI chat assistant only) | Limited JSON | No | No | **Documented OpenAPI 3.1 REST API (22 API Groups)** |
| **SEO Architecture** | Programmatic `/aws/ec2/{instance}` static pages | High-authority `/compare/{a}_vs_{b}` & `/best_vps/{year}` pages | Clean semantic URLs | Informational | Programmatic `/pricing/{provider}` pages | **Nuxt 3 SSR + High-Value Scenario/Provider/SKU Pages + Schema.org `Product`/`Offer` + Anti-Thin-Content Guard** |
| **Global Localization** | English + Multi-currency toggle | English only (USD) | English (USD/EUR) | English | English | **Multi-language (`en`, `ko`, `ja`, `de`, `es`) + Multi-currency Live FX + Local VAT/GST Toggle** |

---

## 2. Key Takeaways & Architectural Adoption for CloudScanr

1. **Adopt `instances.vantage.sh` Speed & Engineer Ergonomics**:
   - Engineers love `instances.vantage.sh` because it loads instantly, allows sorting/filtering by `$/vCPU`, `$/GB RAM`, architecture (`x86_64` vs `arm64`), and shows exact SKU codes (`s-2vcpu-4gb`, `cpx21`, `t4g.medium`). CloudScanr implements an **Engineer Data-Grid View** alongside a **Visual Comparison View**, backed by OpenSearch sub-50ms aggregations and client-side virtualized tables.
2. **Adopt `vpsbenchmarks.com` Price-Bracket & Scenario Methodology**:
   - Buyers rarely ask "what is the cheapest server overall?"; they ask "what is the best server under $15/mo in Tokyo or Frankfurt with NVMe and 2 TB transfer?". CloudScanr implements **Scenario Presets** (e.g., * Hobby Web Server < $10*, *Production K8s Worker 4vCPU/16GB*, *LLM Inference 24GB+ VRAM GPU*, *High-Egress Media Node 10TB/mo*) that dynamically recalculate total monthly cost across all providers.
3. **Establish Unmatched Differentiation via BYOC Provisioning & Honest TCO**:
   - **Zero Competitors Offer Direct BYOC Provisioning**: Neither `instances.vantage.sh` nor `vpsbenchmarks.com` allows a user to connect their DigitalOcean, UpCloud, Vultr, or Hetzner account and spin up the exact compared server in 30 seconds with pre-approved cost guardrails.
   - **Zero Competitors Separate Raw Price from Full TCO Transparently**: CloudScanr explicitly breaks down Base Compute + Mandatory IPv4 Fee + Egress Overage + Backup Cost + Estimated Tax, displaying the exact official Pricing Source URL and collection timestamp.

---

## 3. Clean-Room Intellectual Property Compliance
- CloudScanr **never** scrapes or copies proprietary benchmark text from `vpsbenchmarks.com` or steals frontend code from competitors.
- All pricing data is harvested directly from **Official Provider APIs** or **Official Public Price Books** documented in `OFFICIAL_REFERENCE_LINKS.md`, with full provenance stored in `price_snapshots`.
