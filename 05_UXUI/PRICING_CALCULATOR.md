# PRICING_CALCULATOR.md — Interactive 36-Dimension TCO Calculator UX Specification

**Document ID:** `CS-UX-004`

---

## 1. Dual-Column Price Honesty Layout (Article 8)
Every calculation result card displays two side-by-side panels:
- **Panel A: Provider Raw Price (원본 공급업체 요금)**
  - Exact native currency (`USD`, `EUR`) and billing unit (`$0.0357/hr` up to `$24.00/mo cap` or `€5.49/mo excl. VAT`).
  - Billing rule note (e.g., `672h cap` vs `730h mean month` vs `Pure hourly uncapped`).
- **Panel B: CloudScanr Scenario Estimate (시나리오 예상 총비용)**
  - Stacked bar and line-item breakdown:
    1. Base Compute (`UsageHours` applied)
    2. Public IPv4 Address Surcharge (`$0.00` if included, or `$3.65/mo` on AWS/GCP/Azure)
    3. Extra Block Storage & Provisioned IOPS
    4. Outbound Bandwidth Egress Overage (`MAX(0, ScenarioEgressGB - IncludedGB) * OverageRate`)
    5. Automated Backups & Snapshots
    6. OS License Surcharge
    7. Subtotal in Target Currency (with exact FX Rate & UTC timestamp)
    8. Estimated Local VAT/GST (`0%` / `10% KR VAT` / `19% DE MwSt`)
