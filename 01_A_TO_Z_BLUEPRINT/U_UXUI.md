# U — UX/UI — Complete 57-Screen Design System, Public Web, User Dashboard & Admin Console Specifications

**Document ID:** `CS-AZ-U`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0451 ~ WO-CS-0550, WO-CS-0701 ~ WO-CS-0750`

---

## 1. 기능 목적 (Functional Purpose)
Define the complete User Experience and User Interface architecture across all 57 screens spanning Public Web (16 screens), User Dashboard (18 screens), and Admin Console (23 screens) mandated by Article 16, including Desktop/Mobile wireframes, loading skeletons, error states, empty states, and RBAC visibility rules.

---

## 2. 사용자 역할 (User Roles & Actors)
UX/UI Agent, Frontend Agent, Product Designer.

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. First-time visitor lands on Home (`/`), uses the interactive Spec Pill Builder (`2 vCPU / 4 GB / NVMe / EU or Asia`), compares 4 providers in the sticky drawer, inspects the 36-field TCO breakdown, and clicks through an affiliate offer or connects their BYOC account.
2. Mobile user (`390x844` iPhone viewport) checks their running servers on `/dashboard/servers`, receives a budget warning notification, and safely reboots an unresponsive VPS with thumb-friendly touch targets (`>= 44x44px`).

---

## 4. 기술 스택 (Technology Stack)
TailwindCSS custom design tokens (`05_UXUI/DESIGN_SYSTEM.md`), Inter + JetBrains Mono typography (tabular numerals `font-variant-numeric: tabular-nums` for all prices and specs), Radix Vue accessible dialogs/drawers.

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`All user-facing and admin-facing tables`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)
**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**
- `All 22 API groups mapped per screen in 05_UXUI/SCREEN_INVENTORY.csv`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. UI/UX 설계 (User Interface & Experience)
Full screen-by-screen breakdown documented in `05_UXUI/USER_WEB.md`, `SEARCH_COMPARISON.md`, `PRICING_CALCULATOR.md`, `CLOUD_DASHBOARD.md`, `ADMIN_CONSOLE.md`, `DESIGN_SYSTEM.md`, and `SCREEN_INVENTORY.csv`.

---

## 8. 상세 처리 로직 (Core Processing Logic)
1. Tabular Numeric Alignment: All prices, CPU cores, RAM, bandwidth, and latency columns enforce `tabular-nums` right-aligned formatting for effortless vertical scanning.
2. Four Mandatory UI States per Component: `Loading (Skeleton)`, `Populated (Fresh vs Stale Data Badge)`, `Empty (Guided CTA)`, and `Error (Actionable Retry + Request ID)`.

---

## 9. 보안 통제 (Security & Access Control)
Sensitive UI actions (server deletion, API token rotation, affiliate payout approval) use high-contrast danger modals requiring explicit confirmation input.

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
Optimistic UI updates are strictly forbidden for billable cloud provisioning or server deletion; UI displays real-time verified server state from the backend job status.

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Implemented by Frontend (`F_FRONTEND.md`) and verified by E2E suites (`13_TESTING/`).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
Playwright visual regression and viewport layout checks (`1440x900` Desktop and `390x844` Mobile) across all 57 screens.

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] Every single page listed in Article 16 (16 Public + 18 User Dashboard + 23 Admin Console = 57 screens) is cataloged in `SCREEN_INVENTORY.csv` with exact route, API bindings, auth rules, and test IDs.
- [x] WCAG 2.1 AA contrast ratio (`>= 4.5:1` normal text, `>= 3:1` UI components) in both Light and Dark themes.

---

## 14. 검증 증빙 (Required Verification Evidence)
`05_UXUI/` specification files and `SCREEN_INVENTORY.csv`.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0451 ~ WO-CS-0550, WO-CS-0701 ~ WO-CS-0750`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
