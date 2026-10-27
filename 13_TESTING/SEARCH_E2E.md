# SEARCH_E2E.md — Global Search, Comparison & TCO Calculator E2E Specification
**Document ID:** `CS-TST-002`
## 1. Playwright & API E2E Scenarios
- **Scenario 1**: Natural language spec search (`4 vcpu 16gb nvme tokyo under $40`) decomposes into exact filter chips and matches OpenSearch facets.
- **Scenario 2**: Multi-currency & VAT switch (`USD` -> `KRW` + `10% VAT`) updates all comparison rows and TCO stacked bars reactively while keeping original billing currency visible.
- **Scenario 3**: Pinning 6 SKUs across DigitalOcean, UpCloud, Vultr, Hetzner, Linode, and AWS generates a permanent `/compare/{share_slug}` URL.
