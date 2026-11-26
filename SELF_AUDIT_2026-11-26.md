# 자체 감사 액션 목록 (2026-11-26)

4레인(문서·백엔드·프론트엔드·배포/보안) 감사 합성. 총 76건 (P0 6 / P1 25 / P2 45).
형식: `등급 | 파일:행 | 제목 | 조치 [owner]`. 검증 baseline 전부 통과
(go build/vet/test 18pkg, pnpm typecheck, vitest 5/5, 동결파일·CSV·YAML·문서ID 정상, 시크릿 없음).

## P0 — 즉시 (6건)

- P0 | cmd/api/main.go:67-71 | 짧은 AUTH_TOKEN_SECRET 시 하드코딩 폴백 시크릿으로 기동 지속 → 관리자 토큰 위조 가능 | 시크릿 검증 실패 시 os.Exit(1)로 기동 중단(최소 prod), 폴백 시크릿 제거 [owner: backend]
- P0 | internal/platform/bootstrap.go:64-66 | prod APP_DOMAIN 검증이 구 cloudscanr만 허용 → tirpscanr.com prod 부팅 불가 | TirpScanner 도메인 허용(또는 env 기반 allowlist)으로 수정 후 prod 부팅 재테스트, DEPLOY.md에 APP_DOMAIN/APP_ENV 문서화 [owner: backend]
- P0 | web/pages/dashboard/connections.vue:107 등 + web/layouts/dashboard.vue:10-14 | 퇴역 BYOC/서버 콘솔 4종이 인증 내비 뒤에서 접근 가능 | /dashboard/connections|servers|provision|budgets 제거·리다이렉트, 대시보드 BYOC 내비 삭제 [owner: frontend]
- P0 | web/pages/calculator.vue 등 + web/layouts/public.vue:19-22 | 퇴역 클라우드 카탈로그(VPS/TCO/검색/랭킹)가 공개 내비에 노출 | TirpScanner SaaS/딜 카탈로그 페이지로 교체, TCO/estimate UI 퇴역 [owner: frontend]
- P0 | docker-compose.yml:67-74 + DEPLOY.md:42-48 | 문서상 prod 경로가 하드코딩 dev env로 기동, .env 무시(godotenv/${} 없음) | compose를 env_file/${}로 매개변수화하고 실제 시크릿을 요구하는 prod 오버라이드 추가 [owner: devops-security]
- P0 | cmd/api/main.go:132-138 | 퇴역 BYOC /cloud/* 라우트(servers/connections/approvals/budgets/resources)가 여전히 동작 | 라우트 제거 또는 피처게이트, 또는 문서 정정(삭제 vs 게이트는 제품 결정) [owner: backend]

## P1 — 조속 (25건)

- P1 | cmd/api/main.go:104 + internal/domain/affiliate/service.go:66-68 | AFFILIATE_SIGNING_SECRET 문서화됐으나 코드 미배선, nil이면 하드코딩 기본값 | 플랫폼 설정에서 읽어 NewService에 전달, 길이 강제, prod 누락 시 부팅 거부 [owner: backend]
- P1 | 05_UXUI/SCREEN_INVENTORY.csv:5 | live 화면이 퇴역 GPU 엔드포인트 바인딩 | SCR-PUB-04에서 GET /api/v1/catalog/gpu 제거 또는 행 재표기, V:94 무참조 주장과 정합 [owner: docs]
- P1 | 05_UXUI/SCREEN_INVENTORY.csv:16 | live 화면이 퇴역 K8s 엔드포인트 바인딩 | SCR-PUB-15에서 POST /api/v1/pricing/estimate-kubernetes 제거 또는 행 재표기, K:94와 정합 [owner: docs]
- P1 | 05_UXUI/SCREEN_INVENTORY.csv:15 | live 화면이 퇴역 regions 엔드포인트 바인딩 | SCR-PUB-14에서 GET /api/v1/regions/comparison 제거 또는 행 재표기, N:93과 정합 [owner: docs]
- P1 | 04_API_SPECIFICATION/openapi.yaml:54 | 03 태그 접두사 중복, 05 결번 | 03_Cloud_Catalog을 05_Cloud_Catalog으로 변경(:54,370,377,384) 및 주석 :367 수정 [owner: docs]
- P1 | 04_API_SPECIFICATION/openapi.yaml:145 | live 스키마가 클라우드 하드웨어 모델 유지 | SKUItem(:145-159)·TCOEstimateRequest(:171-187)·search params(:465-473)를 PRICING_CONTRACT.md:11-39 plan/seat/usage 차원으로 재작성 [owner: docs]
- P1 | 04_API_SPECIFICATION/PROVIDER_ADAPTER_CONTRACT.md:29 | CatalogAdapter 개명이 pricing 유지 피벗법과 충돌 | 계약서·COMPONENTS·DATA_FLOW·어댑터 문서를 피벗법(PricingAdapter 유지)에 맞추거나 피벗 개정 [owner: docs]
- P1 | 02_SYSTEM_ARCHITECTURE/THREAT_MODEL.md:13 | KMS 폐기 주장이 live 볼트 문서와 충돌 | THREAT_MODEL:13·SYSTEM_DIAGRAMS:5·COMPONENTS:25,56을 BYOC-토큰 용도만 퇴역으로 재작성(envelope 암호화는 live 유지) [owner: docs]
- P1 | 03_DATABASE/TABLE_DICTIONARY.md:43 | provider_credentials_metadata 퇴역-vs-live 분기 | S_SECURITY.md:41 live 전용과 사전·ERD·S 중 하나의 진실로 통일 [owner: docs]
- P1 | 04_API_SPECIFICATION/openapi.yaml:33 | apex를 ASSUMED으로 오기재 | apex REGISTERED(배선 대기)·api/www ASSUMED으로 재작성(피벗:31-37, openapi:5) [owner: docs]
- P1 | 01_A_TO_Z_BLUEPRINT/I_IDENTITY.md:33 | PASETO/JWT 토큰 구 주장 | v1 HMAC-SHA256 토큰(COMMERCIALIZATION:46-50, openapi:95-97)으로 교체 [owner: docs]
- P1 | 01_A_TO_Z_BLUEPRINT/P_POSTGRESQL.md:25 | nonexistent migration filename | replace 000012_seed_reference_data.up.sql with 000011_seed_providers_and_roles.up.sql and note 000012 relocation to down_migrations_reference/ [owner: docs]
- P1 | README.md:35 | duplicate Service 3 number | renumber Admin & Affiliate Operations row to Service 5 [owner: docs]
- P1 | 01_A_TO_Z_BLUEPRINT/Q_QA.md:106 | reference to nonexistent file | replace DEAL_AFFILIATE_E2E.md (absent repo-wide) with SEARCH_E2E.md/AFFILIATE_E2E.md/SECURITY_E2E.md [owner: docs]
- P1 | web/composables/useApi.ts:55 | stale production host api.cloudscanr.com in API-client doc comment; zero api.tirpscanr.com refs in web/ | update comment/example to https://api.tirpscanr.com per DEPLOY.md:25 [owner: frontend]
- P1 | web/nuxt.config.ts:27 | stale canonical https://cloudscanr.com | change to https://tirpscanr.com [owner: frontend]
- P1 | web/package.json:2,14 + layouts + error.vue:6 + privacy/terms/affiliate-disclosure pages | CloudScanr brand, cloudscanr-web package/Pages project, @cloudscanr.com contact emails stale | rebrand to TirpScanner/tirpscanr.com across UI, meta, emails, package name, deploy:pages project [owner: frontend]
- P1 | web/nuxt.config.ts:61,63 | dead routeRules /compare/** and /articles/** with no matching pages | delete the stale rules or ship the pages [owner: frontend]
- P1 | web/pages/admin/index.vue:12-14,25,55,63 | admin ops health shows retired BYOC/server/approval/management-API signals | rebuild ops cards around pivot LIVE scope (deal ingestion, affiliate ledger, payout reconciliation, catalog moderation) [owner: frontend]
- P1 | docker-compose.yml:absent + DEPLOY.md:87-88 | no worker service though guide documents 15-min harvest cadence | add worker service building ./cmd/worker or document single-container ops explicitly [owner: devops-security]
- P1 | .github/workflows/ci.yml:1-76 vs 14_DEVOPS/CI_CD.md:5-9 | CI runs only go test/build + pnpm typecheck/test/build; none of the 5 documented stages exist and no lint configs are tracked | implement the gates or rewrite CI_CD.md to describe actual ci.yml [owner: devops-security]
- P1 | 14_DEVOPS/BACKUP_RECOVERY.md:4-5 | pgBackRest/S3 PITR + hourly snapshots claimed with zero artifacts | add backup jobs + tested restore procedure or relabel doc as target architecture [owner: devops-security]
- P1 | 14_DEVOPS/MONITORING.md:5-9 vs module_0058.go:62-70 | docs claim tirpscanner_* signals + PagerDuty; code emits only stale-namespace cloudscanr_http_* | emit documented metrics under tirpscanner_* or correct the doc [owner: devops-security]
- P1 | 09_SECURITY/KMS_SECRETS.md:7 + 14_DEVOPS/PRODUCTION_RUNBOOK.md:7 + 14_DEVOPS/COST_MANAGEMENT.md:5 | reference redactor.go, cs-migrator, cmd/janitor — none exist | implement them or mark each reference as planned with a target date [owner: devops-security]
- P1 | docker-compose.yml:45,12-13,27-28,47-48 | DISABLE_SECURITY_PLUGIN + published DB/redis/OS ports on the file DEPLOY uses for prod | split dev/prod compose files; bind dev ports to 127.0.0.1 [owner: devops-security]

## P2 — 여유 될 때 (45건)

- P2 | 04_API_SPECIFICATION/openapi.yaml | 구현된 GET /api/v1/auth/me가 스펙에 없음 | 실제 경로를 스펙에 추가(또는 /users/me 구현) [owner: backend]
- P2 | 04_API_SPECIFICATION/openapi.yaml | 구현된 GET /api/v1/cloud/resources가 스펙에 없음 | 스펙에 문서화 또는 라우트 제거 [owner: backend]
- P2 | 04_API_SPECIFICATION/openapi.yaml:553-567 | connections/{id} 스펙에 구현된 GET 누락 | GET 오퍼레이션 추가 [owner: backend]
- P2 | .env.example:1-50 | 구 cloudscanr 기본값·퇴역 BYOC OAUTH_*·미배선 SMTP/NETWORK_SECRETS 잔존 | 배선된 9개 변수로 정리+도메인 갱신, 계획 변수는 별도 섹션 분리 [owner: devops-security]
- P2 | cmd/api/main.go:126 | 멱등성 스토어가 프로세스-로컬 메모리 | 멀티인스턴스 제한 고지 또는 Redis 백엔드 전환 [owner: backend]
- P2 | cmd/worker 등 4개 패키지 | 테스트 파일 없음(worker·DB 풀·마이그레이션 미검증) | 커버리지 추가(퇴역 management/reseller는 삭제 대안) [owner: backend]
- P2 | web/pages/calculator.vue:292 vs provision.vue:181 | TCO 추정 엔드포인트 2종 혼용 | 하나로 통일(둘 다 퇴역 범위, UI 존속 시에만) [owner: frontend]
- P2 | web/layouts/public.vue:63 | 푸터 고지 오타 FT/KFTC | FTC/KFTC로 수정 [owner: frontend]
- P2 | web/pages/calculator.vue:252 등 4곳 | /go 딜 링크 수제 조립 | dealUrl(slug, src)로 일원화 [owner: frontend]
- P2 | web/package.json | lint 게이트 없음 | eslint 설정+스크립트 추가 또는 typecheck+vitest만 명시 [owner: frontend]
- P2 | web/tests/useNum.test.ts | useNum만 테스트 | resolveUrl/dealUrl/auth 쿠키 단위 테스트 추가 [owner: frontend]
- P2 | deploy/docker/Dockerfile.web:1-9 | npm 사용(pnpm 고정과 충돌), compose 미참조, Pages 전략과 충돌 | 삭제 또는 pnpm으로 재작성+용도 문서화 [owner: devops-security]
- P2 | DEPLOY.md:75,78 | ## 3. 중복 + lockfile 문장 단절 | 번호 수정, lockfileVersion 9.0 문장 완성 [owner: devops-security]
- P2 | .github/workflows/ci.yml:56 vs DEPLOY.md:78-79 | CI가 pnpm 9 부동, 문서의 9.15.0 고정과 충돌 | 9.15.0으로 핀 [owner: devops-security]
- P2 | 15_FINAL_DELIVERY/KNOWN_LIMITATIONS.md:1-32 | 현행 배포 차단 이슈 누락 | 도메인 게이트·미배선 시크릿·백업 미구현·CI 공백 일자별 추록 [owner: devops-security]
- P2 | docker-compose.yml:1 | 구식 version 키 | 제거 [owner: devops-security]
- P2 | 14_DEVOPS/PRODUCTION_RUNBOOK.md:5 vs cmd/worker/main.go:41-54 | SOP-01 격리/재처리에 스키마만 있고 worker는 로그만 기록 | 상태 전이 구현 또는 aspirational 표기 [owner: devops-security]
- P2 | COMMERCIALIZATION.md:114 | 섹션 번호 중복 | 두 번째 ## 3을 ## 5로(DEPLOY:75·OFFICIAL_REFERENCE_LINKS:40·X:19·A-Z 26종 동일) [owner: docs]
- P2 | 02_SYSTEM_ARCHITECTURE/COMPONENTS.md:12 | 숫자 단절 오타 | worker2.·internal2.·application2.·counts1 복원 [owner: docs]
- P2 | 02_SYSTEM_ARCHITECTURE/DEPLOYMENT.md:13 | localhost:300 오타 | localhost:3000으로 수정 [owner: docs]
- P2 | 01_A_TO_Z_BLUEPRINT/L_LEADS_AFFILIATE.md:70 | 상태코드 단절 | 20/301/302 → 200/301/302 [owner: docs]
- P2 | 01_A_TO_Z_BLUEPRINT/J_JOBS.md:26 | 퍼센트 누락 | 35 price drop → 35% [owner: docs]
- P2 | 01_A_TO_Z_BLUEPRINT/O_OPENSEARCH.md:67 | 날짜 스탬프 단절 | tirpscanner_skus_2026112를 8자리로 수정 [owner: docs]
- P2 | 01_A_TO_Z_BLUEPRINT/C_CATALOG.md:103 | ISO 번호 오기 | ISO-3126 → ISO-3166 [owner: docs]
- P2 | 01_A_TO_Z_BLUEPRINT/Q_QA.md:65 | 결함 항목 단절 | 429, 30, truncated feed payload 확장(아마 30s 타임아웃) [owner: docs]
- P2 | 01_A_TO_Z_BLUEPRINT/E_ESTIMATION_ENGINE.md:81 | 검증 범위 의심 | SeatCount 1..10/EstimatedCredits 0..10 확인, 20-off → 20%-off [owner: docs]
- P2 | 01_A_TO_Z_BLUEPRINT/R_RANKING.md:25 | 가중치 합 96% | 40/23/13/10/10을 100%로 재조정 [owner: docs]
- P2 | 01_A_TO_Z_BLUEPRINT/D_DATA_ACQUISITION.md:27 | 이상 임계·백오프 드리프트 | >23% → >25%, 백오프 12s → 16s [owner: docs]
- P2 | 01_A_TO_Z_BLUEPRINT/V_VPS_GPU.md:13 | GPU명 단절 의심 | NVIDIA H10(→H100?)·AMD MI30X(→MI300X?) 확인 [owner: docs]
- P2 | 01_A_TO_Z_BLUEPRINT/N_NETWORK_REGION.md:62 | 개수 불일치 | 3 bandwidth 구조(목록 5개) 수정 [owner: docs]
- P2 | 01_A_TO_Z_BLUEPRINT/A_ARCHITECTURE.md:88 | 오류 클래스 오기 | 3xx errors → 5xx [owner: docs]
- P2 | 15_FINAL_DELIVERY/FINAL_MASTER_REPORT.md:18 | 감사 수치 불일치 | pricing 13→14·CSV 9→10·ops 33→65·item2 10/5→11/4 일자별 추록 [owner: docs]
- P2 | 04_API_SPECIFICATION/openapi.yaml:110 | Idempotency-Key minLength 드리프트 | 12 → 16(16-64 계약) [owner: docs]
- P2 | 04_API_SPECIFICATION/API_TEST_CASES.md:34 | 없는 엔드포인트 대상 테스트 | ATC-22 quarantine resolve 경로 추가 또는 케이스 재작성 [owner: docs]
- P2 | 04_API_SPECIFICATION/API_TEST_CASES.md:26 | removed-vs-deprecated 분기 | ATC-14..19의 404 기대와 deprecated:true 유지 중 택일, CLOUD_MANAGEMENT_CONTRACT:26 정합 [owner: docs]
- P2 | 02_SYSTEM_ARCHITECTURE/COMPONENTS.md:6 | 도메인 상태 과장 | 3종 confirmed → apex REGISTERED+api/www ASSUMED(DELIVERY_CHECKLIST:33·FINAL_MASTER_REPORT:56 동일) [owner: docs]
- P2 | 02_SYSTEM_ARCHITECTURE/COMPONENTS.md:17 | 구 툴체인 버전 | Node 22/Go 1.23 → Go 1.27·Node 24·pnpm 9.15.0, 3셀 행 :19,:25 복구 [owner: docs]
- P2 | 01_A_TO_Z_BLUEPRINT/B_BUSINESS.md:52 | openapi에 없는 엔드포인트 바인딩 다수 | 경로 추가 또는 참조 정정 [owner: docs]
- P2 | 01_A_TO_Z_BLUEPRINT/B_BUSINESS.md:80 | 역할명 드리프트 | FINANCE_ADMIN·DEAL_CURATOR·에이전트명을 IAM 행렬·AGENT_ROLE_MATRIX와 통일 [owner: docs]
- P2 | 01_A_TO_Z_BLUEPRINT/S_SECURITY.md:7 | live 범위가 퇴역 대역 주장 | S·T·Y 범위와 TST-SEC-01 매핑을 퇴역 BYOC/프로비저닝 대역에서 제외 [owner: docs]
- P2 | 01_A_TO_Z_BLUEPRINT/S_SECURITY.md:68 | 승인 TTL 분기 | 3분 승인과 300s payout 토큰(COST_RISK_CONTROLS:5) 통일 [owner: docs]
- P2 | 07_AFFILIATE_ENGINE/COMMISSION_LEDGER.md:9 | 리워드 enum 분기 | SERVICE_CREDIT을 HOSTING_CREDIT과 통일 [owner: docs]
- P2 | 01_A_TO_Z_BLUEPRINT/Z_ZERO_TRUST.md:106 | 섹션 수+동결 목록 오류 | 14 → 15, FINAL_MASTER_REPORT 동결 목록에서 제외, live 추적 추록 추가 [owner: docs]
- P2 | 12_AGENT_EXECUTION/HANDOFF_PROTOCOL.md:5 | 검증 스크립트 부재 | verify_blueprint.py 추가 또는 실제 게이트 명령으로 재작성 [owner: docs]
- P2 | 12_AGENT_EXECUTION/AGENT_ROLE_MATRIX.md:18 | 구 live 화면 수 | 57-screen → 44 live(16+7+21, 퇴역 13) [owner: docs]
- P2 | 02_SYSTEM_ARCHITECTURE/THREAT_MODEL.md:14 | live 위협이 퇴역 테스트로 검증 | TM-02·TM-04 검증을 live SEARCH/AFFILIATE/SECURITY E2E 게이트로 전환 [owner: docs]

## 조치 현황 (2026-11-26 코드 재대조)

코드 소유 항목은 전부 해소 확인. 문서 소유 항목(오타·드리프트·스펙 추록)은 OPEN.
형식: `상태 | 항목 | 근거`.

- DONE | P0 시크릿 폴백 | identity/token.go:54-57 길이 강제 + cmd/api/main.go:83-88 os.Exit(1)
- DONE | P0 APP_DOMAIN | internal/platform/bootstrap.go:70-72 prod tirpscanr 강제 + compose prod 오버라이드 필수화
- DONE | P0 퇴역 콘솔 4종 | dashboard 하위 connections/servers/provision/budgets 페이지 삭제 + server/plugins/retired-redirect.ts 301 (curl 301 확인)
- DONE | P0 퇴역 카탈로그 내비 | layouts/public.vue SaaS 카탈로그 내비 + provision.vue 삭제
- DONE | P0 compose prod 분리 | docker-compose.yml ${} 매개변수화 + docker-compose.prod.yml 필수 시크릿
- DONE | P0 /cloud/* 라우트 | main.go:157-167 피처게이트(기본 off, ENABLE_RETIRED_CLOUD_ROUTES)
- DONE | P1 AFFILIATE_SIGNING_SECRET 배선 | main.go:38-46 resolveAffiliateKey + NewService 전달, prod 32B 강제
- DONE | P1 useApi 주석 호스트 | useApi.ts:52-53 api.tirpscanr.com
- DONE | P1 canonical | nuxt.config.ts:15 tirpscanr.com
- DONE | P1 리브랜드 | package tirpscanner-web, web/ 내 cloudscanr 잔존 0건
- DONE | P1 dead routeRules | /compare/**·/articles/** 규칙 제거, compare 페이지 실존
- DONE | P1 admin ops 카드 | admin/index.vue 공급사/피드/제휴/컴플라이언스 재구성
- DONE | P1 worker 서비스 | docker-compose.yml:85 worker + prod 오버라이드
- DONE | P1 metrics 네임스페이스 | module_0058.go tirpscanner_http_* + 테스트 단언
- DONE | P1 compose 포트 | 전체 127.0.0.1 바인드, DISABLE_SECURITY_PLUGIN은 localhost OpenSearch dev 한정(주석 고지)
- DONE | P2 calculator TCO | 구독 계산기로 재작성, provision.vue 없음
- DONE | P2 푸터 FTC | en.json footer.disclosureNote "FTC/KFTC" (5로케일)
- DONE | P2 /go 일원화 | utils/api-urls.ts buildDealUrl + api-urls.test.ts, 수제 조립 0건
- DONE | P2 lint 게이트 | .github/workflows/ci.yml이 typecheck+test+build를 명시 게이트로 고정 (lint 없음이 명시 상태)
- DONE | P2 단위 테스트 | api-urls.test.ts resolveUrl/buildDealUrl/세션헤더 커버
- DONE | P2 Dockerfile.web | 삭제됨 (deploy/docker는 api+worker만)
- DONE | P2 pnpm 핀 | ci.yml 9.15.0 + packageManager 9.15.0
- DONE | P2 version 키 | docker-compose.yml에 없음
- DONE | P1 #38 CI 문서 | CI_CD.md가 실제 ci.yml 2게이트 + Planned 5단계로 기술됨
- DONE | P1 #39 백업 | BACKUP_RECOVERY.md TARGET 선언 + 미구현 명시
- DONE | P1 #41 미구현 참조 | cs-migrator PLANNED·cmd/janitor unimplemented 표기, redaction은 module_0055 실구현 참조
- DONE | P1 문서 13건 | CSV 3건 무잔존·openapi 태그/apex/스키마·어댑터/KMS/사전/토큰/마이그레이션/README/QA 전건 대조 확인
- DONE | P2 문서 오타·드리프트 전건 | 번호·오타·수치·스펙추록 30건 대조 확인 (FINAL/KNOWN_LIMITATIONS 추록 포함)
- OPEN(후속) | OpenSearch prod TLS | docker-compose.yml:46-47 주석으로 고지됨
- OPEN(외부) | 감사 범위 밖 5항 | DNS/TLS·제휴승인·실DB 등 운영자 실행 필요

## 감사 범위 밖 (미확인)

- 동결 WO 본문 1000건 내용·TRACEABILITY 본문 (개수·파싱만 확인)
- 외부 URL 유효성, 제휴 프로그램 팩트 재검증
- 실 DB 연동 경로, Pages 대시보드 env, api.tirpscanr.com DNS/TLS 상태
- 퇴역 라우트 삭제 vs 게이트 (제품 결정 필요)
- web 소비 엔드포인트의 백엔드 실체 일부, COMMERCIALIZATION/START_HERE 배포 주장 일부
