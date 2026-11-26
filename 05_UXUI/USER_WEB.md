# USER_WEB.md — Public Web & Personalization UX/UI Specification

**Document ID:** `CS-UX-002`

> **Pivot addendum (2026-11-26):** Reframed for TirpScanner. The BYOC deploy CTA is removed; the tertiary action is now a price-drop alert subscription.

---

## 1. Public Discovery & Affiliate Transparency UX (Screens `SCR-PUB-01` to `SCR-PUB-16`)
- **Dual CTA Pattern on Every Plan Row & Tool Detail Page**:
  1. **Primary Affiliate CTA (`Visit Vendor & Claim Deal ↗`)**: Routes via `/go/{slug}` with `rel="sponsored noopener"`, displaying an active deal badge if verified.
  2. **Secondary Direct Source Link (`Official Pricing ↗`)**: Direct link to the vendor's official pricing page (`rel="noopener"`) + UTC collection timestamp (`Verified 22m ago`).
  3. **Tertiary Alert Action (`Alert Me on Price Drop 🔔`)**: Opens the price-alert subscription for the plan (`/dashboard/price-alerts?plan={plan_id}`).
- **Freshness & Confidence Indicator**: Every price cell displays a micro-dot indicator:
  - Green Dot (`FRESH`): Synced within SLA (`< 6h` for vendor feeds, `< 24h` for bulk/network feeds).
  - Amber Dot (`STALE`): Vendor feed temporarily unreachable; displays last verified price and exact UTC timestamp.
