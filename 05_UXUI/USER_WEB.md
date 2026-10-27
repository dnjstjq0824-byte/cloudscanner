# USER_WEB.md — Public Web & Personalization UX/UI Specification

**Document ID:** `CS-UX-002`

---

## 1. Public Discovery & Affiliate Transparency UX (Screens `SCR-PUB-01` to `SCR-PUB-16`)
- **Dual CTA Pattern on Every SKU Row & Product Detail Page**:
  1. **Primary Affiliate CTA (`Visit Provider & Claim $200 Credit ↗`)**: Routes via `/go/{slug}` with `rel="sponsored noopener"`, displaying an active promo badge if verified.
  2. **Secondary Direct Source Link (`Official Pricing Doc ↗`)**: Direct link to the provider's official documentation or price book (`rel="noopener"`) + UTC collection timestamp (`Verified 22m ago`).
  3. **Tertiary BYOC Action (`Deploy on My Account ⚡`)**: Opens `/dashboard/servers/new?sku={sku_id}` for users with a connected BYOC account.
- **Freshness & Confidence Indicator**: Every price cell displays a micro-dot indicator:
  - Green Dot (`FRESH`): Synced within SLA (`< 6h` for REST API, `< 24h` for bulk catalog).
  - Amber Dot (`STALE`): Provider API temporarily unreachable; displays last verified price and exact UTC timestamp.
