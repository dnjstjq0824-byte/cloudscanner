# SUBID_CLICKID.md — Cryptographic ClickID & SubID Attribution Specification
**Document ID:** `CS-AE-002`

> **Pivot addendum (2026-11-26):** Core TirpScanner monetization surface — kept; network mappings unchanged, merchant examples reframed to SaaS vendors.

## 1. Parameter Mapping Across Affiliate Networks
- **Impact Radius (SaaS vendors)**: `?subId1={click_ulid}&subId2={hmac_sig16}&subId3={placement_slot}`
- **CJ Affiliate (SaaS vendors)**: `?sid={click_ulid}_{hmac_sig16}`
- **PartnerStack (B2B SaaS)**: `?utm_content={click_ulid}&ps_xid={click_ulid}`
- **PostAffiliatePro-based programs**: `?a_aid={pub_id}&data1={click_ulid}&data2={hmac_sig16}`
- **HasOffers / Tune (SaaS vendors)**: `?aff_sub={click_ulid}&aff_sub2={hmac_sig16}`
