# FX_TAXES.md — Multi-Currency Foreign Exchange & Global Indirect Tax Engine
**Document ID:** `CS-PE-002`
## 1. Foreign Exchange (FX) Pipeline
- Polls ECB and OpenExchangeRates every 6 hours into `fx_rates` using `NUMERIC(18,8)` precision.
- Every API response and UI tooltip shows:
  1. `native_currency` and `raw_price` (e.g., `€5.49 EUR`).
  2. `target_currency` converted amount (e.g., `$5.94 USD` or `₩8,210 KRW`).
  3. `fx_rate_applied` (`1 EUR = 1.08166577 USD`) and `fx_rate_timestamp` (`2026-10-28T00:00:00Z`).
## 2. Indirect Tax (VAT / GST / Consumption Tax) Rules
- Most cloud providers publish B2B tax-exclusive prices (`tax_included = false`).
- When a user enables `Include Local Tax` for country `c`, CloudScanr looks up `tax_policies` (`KR = 10% VAT`, `DE = 19% VAT`, `FI = 25.5% VAT`, `JP = 10% CT`, `SG = 9% GST`) and appends the tax line item separately from base infrastructure cost.
