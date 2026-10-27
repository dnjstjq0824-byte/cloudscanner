# NORMALIZATION.md — 36-Dimension Cloud Price & Hardware Normalization Specification
**Document ID:** `CS-PE-001`
## 1. Canonical Hardware & Billing Unit Normalization
- **Memory**: Provider responses in `MB` (1000-based) vs `MiB` (1024-based) vs `GB` are normalized to canonical `ram_mb` (MiB, where `1 GiB = 1024 MiB`) and displayed with clear `GB` labels.
- **CPU Allocation Classification**:
  - `BURSTABLE`: AWS `t3`/`t4g`, GCP `e2-micro/small/medium`, Azure `B-series`.
  - `SHARED`: DigitalOcean Basic Regular/Premium, Vultr Regular/High-Performance, Hetzner `CX`/`CAX`/`CPX`, Linode Shared, UpCloud Developer/General Purpose.
  - `DEDICATED`: DigitalOcean CPU-Optimized/General Purpose, Hetzner `CCX`, Vultr Optimized Cloud, Linode Dedicated, UpCloud High CPU, AWS `m7i`/`c7g`.
- **Billing Cap vs Uncapped Hourly**:
  - Providers with monthly caps (DigitalOcean 672h cap, Vultr 672h cap, Linode, Hetzner, UpCloud) cap monthly compute charges at `raw_monthly_price` even in a 744-hour (31-day) month.
  - Pure hourly providers (AWS, GCP, Azure) charge `raw_hourly_price * actual_month_hours` (`672h` in Feb, `720h` in 30d month, `744h` in 31d month, `730h` annual mean). CloudScanr models both and lets users toggle the exact billing assumption.
