# PRICING_CONTRACT.md — 36-Dimension Price Normalization & TCO API Contract

**Document ID:** `CS-API-003`

---

## 1. Mandatory 36 Normalization Attributes in API Payloads (Article 8)

Every SKU returned by `/api/v1/catalog/skus/{sku_id}` and `/api/v1/pricing/estimate` exposes all 36 required dimensions grouped into logical sub-objects:

| # | Attribute Name | JSON Field Path | Type | Description |
| :--- | :--- | :--- | :--- | :--- |
| 1 | Provider | `provider.slug` | `string` | Canonical provider slug (`digitalocean`, `hetzner`, etc.) |
| 2 | Product Family | `product.family_code` | `string` | e.g., `Premium AMD Droplets`, `CCX Dedicated` |
| 3 | SKU | `external_sku_code` | `string` | Exact provider SKU (`s-2vcpu-4gb`, `cx23`, `t4g.medium`) |
| 4 | vCPU | `specs.vcpu_count` | `integer` | Virtual CPU count |
| 5 | CPU Architecture | `specs.cpu_architecture` | `string` | `x86_64` or `arm64` |
| 6 | Shared / Dedicated CPU | `specs.cpu_allocation_type` | `string` | `SHARED`, `BURSTABLE`, or `DEDICATED` |
| 7 | RAM | `specs.ram_mb` | `integer` | Memory in MiB |
| 8 | Storage Type | `specs.storage_type` | `string` | `NVME`, `SSD`, `HDD`, `NETWORK_BLOCK`, `NONE` |
| 9 | Storage Capacity | `specs.storage_gb` | `integer` | Bundled local storage in GB |
| 10 | Disk Performance | `specs.disk_iops_read` | `integer` | Provisioned/Max Read/Write IOPS & MB/s |
| 11 | Included Bandwidth | `specs.included_bandwidth_gb` | `integer` | Monthly included data transfer in GB |
| 12 | Data Transfer Egress | `specs.bandwidth_billing_direction` | `string` | `OUTBOUND_ONLY`, `BIDIRECTIONAL`, `UNMETERED_PORT` |
| 13 | Overage Charge | `pricing.egress_overage_per_gb_price` | `decimal_str` | Per-GB overage fee in native currency |
| 14 | Region | `region.region_code` | `string` | Provider region + ISO Country + City |
| 15 | Availability | `availability_status` | `string` | `AVAILABLE`, `LOW_STOCK`, `OUT_OF_STOCK`, `DEPRECATED` |
| 16 | IPv4 Cost | `pricing.ipv4_monthly_addon_price` | `decimal_str` | `0.00` if included, else monthly surcharge |
| 17 | IPv6 Support | `specs.ipv6_supported` | `boolean` | True if native `/64` or `/128` IPv6 is free |
| 18 | Backup Price | `pricing.backup_monthly_price_or_pct` | `decimal_str` | Automated backup rate (e.g., `20%` of base compute) |
| 19 | Snapshot Price | `pricing.snapshot_per_gb_month_price` | `decimal_str` | Per-GB/month snapshot storage rate |
| 20 | OS License | `pricing.os_license_hourly_price` | `decimal_str` | Windows/RHEL license surcharge (`0.00` for Linux) |
| 21 | GPU Model | `specs.gpu_model` | `string?` | e.g., `NVIDIA H100 SXM5`, `RTX 4090`, `null` |
| 22 | GPU VRAM | `specs.gpu_vram_gb` | `integer` | Total GPU VRAM in GB |
| 23 | Hourly Billing | `pricing.raw_hourly_price` | `decimal_str` | Exact native hourly billing rate |
| 24 | Monthly Billing | `pricing.raw_monthly_price` | `decimal_str` | Exact native monthly cap or fixed price |
| 25 | Annual Commitment | `pricing.raw_annual_price` | `decimal_str` | 1-Year commitment rate if offered |
| 26 | Reserved/Spot Pricing | `billing_term.term_code` | `string` | `ON_DEMAND`, `RESERVED_1YR`, `RESERVED_3YR`, `SPOT` |
| 27 | Setup Fee | `pricing.setup_fee` | `decimal_str` | One-time setup fee (common on Bare Metal) |
| 28 | Tax Inclusion | `pricing.tax_included` | `boolean` | Whether raw price includes VAT |
| 29 | Currency | `pricing.native_currency` | `string` | ISO-4217 currency code (`USD`, `EUR`) |
| 30 | FX Rate | `pricing.fx_rate_to_target` | `decimal_str` | Applied exchange rate + `fx_rate_timestamp` |
| 31 | Promotion | `promotion.credit_amount` | `decimal_str` | Eligible trial credit (e.g., `$200 Free Credit`) |
| 32 | Promotion Expiry | `promotion.validity_days` | `integer` | Trial validity window in days |
| 33 | Price Effective Date | `pricing.price_effective_date` | `datetime` | UTC date provider rate took effect |
| 34 | Collection Timestamp | `pricing.collected_at` | `datetime` | Exact UTC timestamp synced by CloudScanr |
| 35 | Pricing Source | `pricing.official_source_url` | `uri` | Direct link to official API doc or price page |
| 36 | Confidence / Verification | `pricing.verification_status` | `string` | `VERIFIED`, `FRESH`, `STALE`, `UNVERIFIED` |
