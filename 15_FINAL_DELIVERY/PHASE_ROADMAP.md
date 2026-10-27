# PHASE_ROADMAP.md — Phase 0 to Phase 10 Master Execution Roadmap (Article 23)
**Document ID:** `CS-FD-001`

| Phase | Phase Name | Work Order Scope | Entry Prerequisites | Completion Gate & Required Evidence | Deployment Scope |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Phase 0** | **시장 실사 (Market & API Due Diligence)** | `WO-CS-0001` ~ `0050` | Project Charter | Verified `00_MARKET_RESEARCH/*.csv` with official URLs & status enums (`VERIFIED`, `REQUIRES_APPLICATION`, `CONTRACT_REQUIRED`). | Internal Blueprint |
| **Phase 1** | **프로젝트 기반 (Platform Foundation)** | `WO-CS-0051` ~ `0100` | Phase 0 Sign-off | Go 1.23 + Nuxt 3 + PG 16 (49 tables) + OpenSearch 2.17 + Redis 7.4 + CI/CD green; zero TreeScan dependency. | Staging Cluster |
| **Phase 2** | **가격 비교 MVP (Price Comparison MVP)** | `WO-CS-0151` ~ `0350`, `0451` ~ `0550` | Phase 1 Sign-off | 18+ providers cataloged, 36-field TCO calculator, sub-50ms OpenSearch facets, 6-SKU compare matrix live. | Public Beta (`cloudscanr.com`) |
| **Phase 3** | **Affiliate 수익화 (Affiliate Monetization)** | `WO-CS-0351` ~ `0450`, `0715` ~ `0727` | Phase 2 Live Website | `/go/{slug}` router (<15ms), signed `ClickID`/`SubID`, postback/CSV reconciliation, double-entry Commission Ledger, FTC/KFTC disclosures. | Production Monetization |
| **Phase 4** | **회원 및 개인화 (Identity & Personalization)** | `WO-CS-0101` ~ `0150`, `0496` ~ `0499` | Phase 2 Sign-off | Argon2id Auth, TOTP + WebAuthn Passkeys, Orgs & RBAC, Favorites, Price Drop Alerts. | Production Auth |
| **Phase 5** | **BYOC 연결 (BYOC Account Connection)** | `WO-CS-0551` ~ `0600` | Phase 4 + KMS Audit | OAuth2 & Scoped Token connections, AES-256-GCM KMS Envelope Encryption with AAD, Read-Only inventory view. | Production BYOC Read-Only |
| **Phase 6** | **서버 생성·관리 (Server Provisioning & Control)** | `WO-CS-0601` ~ `0650` | Phase 5 + 22 Safety Tests | Cost-approved server provisioning, timeout reconciliation, power actions, double-gated safe deletion. | Production BYOC Read-Write |
| **Phase 7** | **고급 클라우드 운영 (Advanced Cloud Ops)** | `WO-CS-0651` ~ `0750`, `0801` ~ `0850` | Phase 6 Sign-off | Resizing, automated backups/snapshots, firewalls, SSH keys, budget hard-blocks, drift reconciliation, full Admin Console. | Production Full Suite |
| **Phase 8** | **글로벌 확장 (Global Expansion & SEO)** | `WO-CS-0751` ~ `0800` | Phase 3 & 7 Sign-off | 5 locales (`en`, `ko`, `ja`, `de`, `es`), live multi-currency FX + local VAT/GST, verified customer reviews, YABS benchmarks. | Global Edge CDN |
| **Phase 9** | **상용화 감사 (Commercial Security & QA Audit)** | `WO-CS-0851` ~ `0950` | Phase 1–8 Complete | 100% pass on 22 safety scenarios, k6 1,500 RPS load test, IDOR/SSRF penetration test, PITR DR restore drill. | Pre-Launch Gate |
| **Phase 10** | **정식 운영·확장 (GA Launch & Future Reseller)** | `WO-CS-0951` ~ `1000` | Phase 9 Independent Audit | Independent Zero-Trust Audit sign-off, GA production launch, and post-contract B2B Reseller evaluation. | Global General Availability |
