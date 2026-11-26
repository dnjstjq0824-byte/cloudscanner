# SEARCH_E2E.md — Global Search, Comparison & Deal-Scan E2E Specification
**Document ID:** `CS-TST-002`
> TirpScanner pivot (2026-11-26): LIVE scope — reframed from cloud SKUs/TCO to SaaS/AI-tool search, comparison, and price-drop detection. See `TIRPSCANNER_PIVOT.md`.
## 1. Playwright & API E2E Scenarios
- **Scenario 1**: Natural language tool search (`ai video editor team plan under $20`) decomposes into exact filter chips and matches OpenSearch facets.
- **Scenario 2**: Multi-currency & VAT switch (`USD` -> `KRW` + `10% VAT`) updates all comparison rows and deal price-drop badges reactively while keeping original billing currency visible.
- **Scenario 3**: Pinning 6 SaaS/AI-tool plans across vendors generates a permanent `/compare/{share_slug}` URL.
