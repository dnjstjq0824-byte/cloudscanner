# F — Frontend — Nuxt 3 SSR Web Application, Technical SEO, Interactive Comparison & Accessibility

**Document ID:** `CS-AZ-F`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0451 ~ WO-CS-0550`

---

## 1. 기능 목적 (Functional Purpose)
Build a blazing-fast, SEO-optimized, accessible (WCAG 2.1 AA) Nuxt 3 / Vue 3 frontend supporting server-side rendering (SSR) for public discovery pages and rich interactive client-side state for multi-SKU comparison, TCO calculation, and BYOC server management.

---

## 2. 사용자 역할 (User Roles & Actors)
Frontend Agent, UX/UI Designer, SEO Specialist, End User.

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. Search engine crawler (`Googlebot`, `TreeScanbot`) requests `/providers/digitalocean` or `/vps-comparison`; Nuxt 3 returns complete SSR HTML with canonical tags, hreflang, and Schema.org `Product`/`AggregateOffer` JSON-LD in <150ms.
2. User filters 5,000+ SKUs in Engineer Data-Grid Mode with virtual scrolling at 60 FPS, pins 4 SKUs to the sticky comparison tray, and opens the side-by-side matrix with diff-highlighting.
3. Keyboard-only / screen-reader user navigates the server provisioning modal and confirms the cost approval step with full ARIA announcements and focus trapping.

---

## 4. 기술 스택 (Technology Stack)
Nuxt 3.13+, Vue 3.5 Composition API, TypeScript 5.6 (Strict mode), Pinia state stores, TanStack Table v8 (virtualized grid), TailwindCSS + Radix Vue headless accessible primitives, Playwright E2E.

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`comparison_lists`, `favorites`, `price_alerts`, `articles`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)
**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**
- `GET /api/v1/search/skus`
- `POST /api/v1/comparisons`
- `GET /api/v1/content/articles`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. UI/UX 설계 (User Interface & Experience)
Responsive design system (`05_UXUI/DESIGN_SYSTEM.md`) supporting Desktop (`1440px+`), Tablet (`768px–1024px`), and Mobile (`360px–430px`), with native Light/Dark mode (`prefers-color-scheme` + user override) and zero layout shift (`CLS < 0.05`).

---

## 8. 상세 처리 로직 (Core Processing Logic)
1. Hybrid Rendering Strategy (`routeRules` in `nuxt.config.ts`):
   - Public Catalog, Provider, Product, Comparison, Review, Tutorial pages: SSR with stale-while-revalidate edge caching (`swr: 300`).
   - `/dashboard/**` and `/admin/**`: CSR/SSR protected by auth middleware (`middleware/auth.global.ts`) with `X-Robots-Tag: noindex, nofollow`.
2. Persistent Compare Tray: Syncs to `localStorage` for anonymous visitors and merges into PostgreSQL `comparison_lists` upon login.

---

## 9. 보안 통제 (Security & Access Control)
Strict Content Security Policy (CSP) with nonce-based script execution; `HttpOnly; Secure; SameSite=Lax` session cookies; zero BYOC cloud API tokens ever stored in `localStorage`, `sessionStorage`, or Pinia state.

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
Granular `<NuxtErrorBoundary>` wrappers per widget (if benchmark chart fails, pricing table still renders); offline/degraded banner when API latency spikes.

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Consumes OpenAPI 3.1 contracts (`04_API_SPECIFICATION/openapi.yaml`) and Design System (`05_UXUI/DESIGN_SYSTEM.md`).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
Vitest component tests, Playwright E2E across Chromium/Firefox/WebKit, axe-core automated WCAG 2.1 AA accessibility audits, and Lighthouse CI (`Performance >= 92`, `SEO = 100`, `Accessibility >= 95`).

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] All 16 Public Web pages, 18 User Dashboard pages, and 23 Admin Console pages defined in Article 16 are implemented with responsive PC/Mobile layouts.
- [x] Lighthouse SEO score = 100 on all public pages; `/dashboard/*` and `/admin/*` strictly emit `noindex, nofollow`.
- [x] Zero sensitive API credentials appear in browser storage or network query strings.

---

## 14. 검증 증빙 (Required Verification Evidence)
Playwright test report HTML, Lighthouse CI JSON scorecards, and SSR curl HTML snapshots.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0451 ~ WO-CS-0550`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
