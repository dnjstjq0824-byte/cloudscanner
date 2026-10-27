# PROVIDER_PRICING_ADAPTERS.md — Per-Provider Pricing Ingestion Specifications
**Document ID:** `CS-PE-005`
## 1. Tier-1 & Tier-2 Adapter Endpoint & Mapping Rules
1. **DigitalOcean (`digitalocean`)**: `GET https://api.digitalocean.com/v2/sizes` -> maps `slug`, `memory`, `vcpus`, `disk`, `transfer`, `price_monthly`, `price_hourly`, `regions[]`.
2. **UpCloud (`upcloud`)**: `GET https://api.upcloud.com/1.3/plan` + `GET /1.3/price` -> maps zone-specific hourly pricing and zero-cost egress transfer policy.
3. **Cloudways (`cloudways`)**: Strictly **API v2** `GET https://api.cloudways.com/api/v2/server_sizes` + `/api/v2/providers` (DO, Vultr, Linode, AWS, GCP managed tiers).
4. **Vultr (`vultr`)**: Unauthenticated `GET https://api.vultr.com/v2/plans?type=all` and `GET /v2/plans-metal` -> maps `vc2`, `vhf`, `vhp`, `voc`, `vcg` (GPU).
5. **Hetzner (`hetzner`)**: `GET https://api.hetzner.cloud/v1/pricing` & `GET /v1/server_types` -> maps per-location (`fsn1`, `nbg1`, `hel1`, `ash`, `hil`, `sin`) net/gross hourly & monthly caps + IPv4 add-on price.
6. **Akamai / Linode (`akamai_linode`)**: Unauthenticated `GET https://api.linode.com/v4/linode/types` -> maps `price.hourly`, `price.monthly`, and regional price overrides (`region_prices[]`).
7. **Google Cloud (`gcp`)**: Cloud Billing Catalog API `GET https://cloudbilling.googleapis.com/v1/services/6F81-5844-456A/skus` -> combines vCPU + RAM SKU components per machine family (`e2`, `n2`, `c3`, `a2`/`a3` GPU).
8. **AWS (`aws`)**: Bulk Price List API `https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEC2/current/{region}/index.json` streamed via Go `json.Decoder` token stream (zero OOM).
9. **Microsoft Azure (`azure`)**: Unauthenticated `GET https://prices.azure.com/api/retail/prices?$filter=serviceName eq 'Virtual Machines'`.
10. **OVHcloud (`ovhcloud`)**: `GET https://api.ovh.com/1.0/order/catalog/public/cloud?ovhSubsidiary=US`.
11. **Scaleway (`scaleway`)**: `GET https://api.scaleway.com/instance/v1/zones/{zone}/products/servers`.
12. **Oracle Cloud (`oracle_cloud`)**: `GET https://apexapps.oracle.com/pls/apex/cetools/api/v1/products/` (normalizing `1 OCPU = 2 vCPUs` on x86).
