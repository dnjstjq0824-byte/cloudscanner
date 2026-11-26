# U — UX/UI — Deal-Scanner Design System: Public Web, User Dashboard & Admin Console Specifications

**Document ID:** `CS-AZ-U` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (TirpScanner pivot reframe, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Pivot Note (2026-11-26):** Console server/provisioning screens retired; catalog/comparison/alert/deal/admin screens reframed for SaaS buying. See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0451 ~ WO-CS-0550, WO-CS-0701 ~ WO-CS-0750`

---

## 1. 기능 목적 (Functional Purpose)

Define the complete User Experience and User Interface architecture across Public Web, User Dashboard, and Admin Console screens, including Desktop/Mobile wireframes, loading skeletons, error states, empty states, and RBAC visibility rules — reframed for SaaS/AI-tool discovery, plan comparison, deal alerts, and affiliate operations.

---

## 2. User Roles & Actors

UX/UI Agent, Frontend Agent, Product Designer.

---

## 3. User Scenarios

1. First-time visitor lands on Home (`/`), uses the interactive Plan Pill Builder (`3 seats / Annual / AI Video / under $20 per seat`), compares 4 plans in the sticky drawer, inspects the effective-cost breakdown, and clicks through an affiliate deal or subscribes to a price-drop alert.
2. Mobile user (`390x844` iPhone viewport) checks tracked deals on `/dashboard/alerts`, receives a price-drop notification, and opens the verified coupon with thumb-friendly touch targets (`>= 44x44px`).

---

## 4. Technology Stack

TailwindCSS custom design tokens (`05_UXUI/DESIGN_SYSTEM.md`), Inter + JetBrains Mono typography (tabular numerals `font-variant-numeric: tabular-nums` for all prices and plan specs), Radix Vue accessible dialogs/drawers.

---

## 5. Database Schema & Tables

**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`All user-facing and admin-facing tables`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)

**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**

- `All live API groups mapped per screen in 05_UXUI/SCREEN_INVENTORY.csv (cloud groups retired)`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. User Interface & Experience

Full screen-by-screen breakdown documented in `05_UXUI/USER_WEB.md`, `SEARCH_COMPARISON.md`, `PRICING_CALCULATOR.md`, `ADMIN_CONSOLE.md`, `DESIGN_SYSTEM.md`, and `SCREEN_INVENTORY.csv`. (Cloud dashboard specs retired.)

---

## 8. Core Processing Logic

1. Tabular Numeric Alignment: All prices, seat counts, quotas, discount depths, and deal scores enforce `tabular-nums` right-aligned formatting for effortless vertical scanning.
2. Four Mandatory UI States per Component: `Loading (Skeleton)`, `Populated (Fresh vs Stale Data Badge)`, `Empty (Guided CTA)`, and `Error (Actionable Retry + Request ID)`.

---

## 9. Security & Access Control

Sensitive UI actions (credential rotation, affiliate payout approval, deal delisting) use high-contrast danger modals requiring explicit confirmation input.

---

## 10. Failure Modes & Resilience

Optimistic UI updates are strictly forbidden for payout approval or credential rotation; UI displays real-time verified backend state.

---

## 11. Inter-Module Dependencies

Implemented by Frontend (`F_FRONTEND.md`) and verified by E2E suites (`13_TESTING/`).

---

## 12. Testing & Verification Strategy

Playwright visual regression and viewport layout checks (`1440x900` Desktop and `390x844` Mobile) across all live screens.

---

## 13. Acceptance Criteria

- [x] Every live page (Public + User Dashboard + Admin Console, excluding retired console server screens) is cataloged in `SCREEN_INVENTORY.csv` with exact route, API bindings, auth rules, and test IDs.
- [x] WCAG 2.1 AA contrast ratio (`>= 4.5:1` normal text, `>= 3:1` UI components) in both Light and Dark themes.

---

## 14. Required Verification Evidence

`05_UXUI/` specification files and `SCREEN_INVENTORY.csv`.

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0451 ~ WO-CS-0550, WO-CS-0701 ~ WO-CS-0750`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
