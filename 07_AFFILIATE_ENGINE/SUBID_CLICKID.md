# SUBID_CLICKID.md — Cryptographic ClickID & SubID Attribution Specification
**Document ID:** `CS-AE-002`
## 1. Parameter Mapping Across Affiliate Networks
- **Impact Radius (DigitalOcean, GCP)**: `?subId1={click_ulid}&subId2={hmac_sig16}&subId3={placement_slot}`
- **CJ Affiliate (Contabo, IONOS)**: `?sid={click_ulid}_{hmac_sig16}`
- **PartnerStack (UpCloud)**: `?utm_content={click_ulid}&ps_xid={click_ulid}`
- **Cloudways (PostAffiliatePro)**: `?a_aid={pub_id}&data1={click_ulid}&data2={hmac_sig16}`
- **HasOffers / Tune (Hostinger)**: `?aff_sub={click_ulid}&aff_sub2={hmac_sig16}`
