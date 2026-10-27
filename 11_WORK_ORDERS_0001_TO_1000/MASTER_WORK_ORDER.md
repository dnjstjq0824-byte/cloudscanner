# MASTER_WORK_ORDER.md — CloudScanr 1,000 Executable Work Orders Master Governance

**Document ID:** `CS-WO-MASTER`  
**Total Work Orders:** `1,000` (`WO-CS-0001` through `WO-CS-1000` — Zero Gaps, Zero Duplicates)

---

## 1. Work Order Domain Band Allocation (Article 20 Compliance)

| Work Order Range | Count | Domain Area (작업 분야) | Primary Phase | Lead Responsible Agent |
| :--- | :--- | :--- | :--- | :--- |
| `WO-CS-0001` – `WO-CS-0050` | 50 | 글로벌 시장·경쟁·제휴 조사 (Market, Competitor & Affiliate Due Diligence) | Phase 0 | Affiliate Research Agent / Provider API Agent |
| `WO-CS-0051` – `WO-CS-0100` | 50 | 요구사항·아키텍처·인프라 설계 (Architecture, DB DDL & Infra Foundation) | Phase 1 | Lead Architect / Database Agent / DevOps Agent |
| `WO-CS-0101` – `WO-CS-0150` | 50 | 회원·인증·권한·조직 (Identity, OIDC, Passkey MFA, RBAC & Organizations) | Phase 4 | Backend Agent / Security Agent |
| `WO-CS-0151` – `WO-CS-0200` | 50 | 공급업체 등록·상품 카탈로그 (18+ Providers & 16 Product Categories Catalog) | Phase 2 | Backend Agent |
| `WO-CS-0201` – `WO-CS-0250` | 50 | 공식 가격 API·피드 커넥터 (Official Pricing API Connectors & Anomaly Guard) | Phase 2 | Provider API Agent |
| `WO-CS-0251` – `WO-CS-0300` | 50 | 가격 정규화·환율·세금·총비용 (36-Field Normalization, FX, VAT/GST & TCO Engine) | Phase 2 | Pricing Engine Agent |
| `WO-CS-0301` – `WO-CS-0350` | 50 | OpenSearch·필터·검색·랭킹 (Sub-50ms Facets, NL Spec Parser & Unbiased Ranking) | Phase 2 | Search Agent |
| `WO-CS-0351` – `WO-CS-0400` | 50 | Affiliate 가입·계약·링크 관리 (Affiliate Programs, Contracts, Campaigns & Links) | Phase 3 | Backend Agent / Affiliate Research Agent |
| `WO-CS-0401` – `WO-CS-0450` | 50 | 클릭·전환·커미션·정산 대사 (`/go/{slug}` Router, Postbacks, CSV Reconcile & Ledger) | Phase 3 | Backend Agent |
| `WO-CS-0451` – `WO-CS-0500` | 50 | 사용자 웹·홈·탐색·상품 상세 (Public Web Nuxt 3 SSR, VPS/GPU Views & Alerts) | Phase 2 / 4 | Frontend Agent / UX/UI Agent |
| `WO-CS-0501` – `WO-CS-0550` | 50 | 비교 UI·가격 계산기·SEO (6-SKU Compare Matrix, TCO Calculator & Technical SEO) | Phase 2 / 8 | Frontend Agent |
| `WO-CS-0551` – `WO-CS-0600` | 50 | 클라우드 계정 연결·토큰 보안 (BYOC OAuth/Token Connection & KMS Envelope Security) | Phase 5 | Security Agent / Cloud Management Agent |
| `WO-CS-0601` – `WO-CS-0650` | 50 | 서버 생성·조회·전원·관리 (Cost-Approved Server Provisioning, Power & Safe Delete) | Phase 6 | Cloud Management Agent |
| `WO-CS-0651` – `WO-CS-0700` | 50 | 리사이징·네트워크·백업·예산 (Resize, Backups, Firewalls, SSH, Budgets & Drift) | Phase 7 | Cloud Management Agent |
| `WO-CS-0701` – `WO-CS-0750` | 50 | 관리자 콘솔·공급업체 운영 (23-Screen Admin Console, Quarantine & Ledger Ops) | Phase 3 / 7 | Backend Agent / Frontend Agent |
| `WO-CS-0751` – `WO-CS-0800` | 50 | 리뷰·콘텐츠·다국어·접근성 (Benchmarks, Verified Reviews, 5-Locale i18n & WCAG AA) | Phase 8 | Frontend Agent |
| `WO-CS-0801` – `WO-CS-0850` | 50 | DB·작업 큐·보안·운영 (Partitioning, Asynq DLQ, Audit Hash-Chain & DR Drills) | Phase 7 / 9 | Database Agent / Security Agent / DevOps Agent |
| `WO-CS-0851` – `WO-CS-0900` | 50 | 단위·통합·E2E·API 테스트 (Unit, Testcontainers, OpenAPI Contract & 22 Safety Tests) | Phase 9 | QA Agent |
| `WO-CS-0901` – `WO-CS-0950` | 50 | 비용·부하·보안·장애 테스트 (k6 1500 RPS Load, Chaos, IDOR/SSRF Penetration & Cost Cap) | Phase 9 | QA Agent / Security Agent |
| `WO-CS-0951` – `WO-CS-1000` | 50 | 독립 감사·출시·인수인계 (Independent Zero-Trust Audit, Traceability & Launch) | Phase 10 | Independent Audit Agent / Lead Architect |
