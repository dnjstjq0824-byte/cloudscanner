# PROVIDER_PRICING_ADAPTERS.md — SaaS/AI-Tool Deal-Ingestion Adapter Specification
**Document ID:** `CS-PE-005`

> **Pivot addendum (2026-11-26):** Rewritten for TirpScanner. The 12 per-cloud-provider endpoint mappings (DigitalOcean, Vultr, Hetzner, AWS, etc.) are RETIRED with the cloud catalog; preserved in git history. Live ingestion follows the tier model below.

## 1. Ingestion Tiers & Adapter Rules
1. **Tier-1 — Vendor pricing feeds & APIs**: Where a SaaS vendor exposes a pricing API, structured pricing page, or changelog feed, the `CatalogAdapter` polls it directly (hourly). Raw payloads are SHA-256 hashed and archived for provenance.
2. **Tier-2 — Affiliate network product catalogs**: Impact, CJ, and PartnerStack product feeds provide plan names, list prices, and promotion windows; synced every 6 hours and cross-checked against vendor sources before going live.
3. **Tier-3 — Curated deal intake**: Time-boxed coupons and launch deals without machine-readable feeds enter via the admin curation flow with mandatory `official_source_url`, `expires_at`, and eligibility notes; every curated deal requires source-link verification before publish.
4. **Mapping rules (all tiers)**: Every ingested record must map to the canonical plan model (`vendor.slug`, `plan.tier_code`, `billing_cycle`, `seat_model`, `monthly_equivalent`, `native_currency`, `tax_included`, `official_source_url`); records missing required fields route to `QUARANTINED`, never live.
