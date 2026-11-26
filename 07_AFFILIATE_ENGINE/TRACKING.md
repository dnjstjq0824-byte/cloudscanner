# TRACKING.md — Outbound Affiliate Link Router & Click Telemetry Engine
**Document ID:** `CS-AE-001`

> **Pivot addendum (2026-11-26):** Core TirpScanner monetization surface — kept and rebranded.

## 1. Architecture of `/go/{slug}`
- Designed for `< 10ms` p99 response time and zero lost clicks.
- Never trusts client-supplied target URLs; resolves `slug` strictly against `affiliate_links`.
- Injects `Rel: nofollow sponsored` HTTP headers (`X-Robots-Tag: noindex, nofollow`) so search engines never penalize TirpScanner for outbound affiliate redirects.
