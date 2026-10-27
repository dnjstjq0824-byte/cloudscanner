# DESIGN_SYSTEM.md — CloudScanr Design System, Typography, Tokens & Accessibility

**Document ID:** `CS-UX-001`

---

## 1. Visual Identity & Semantic Color Palette (WCAG 2.1 AA Verified)
- **Brand Primary (Cloud Cobalt)**: `#0F62FE` (Hover `#0043CE`, Active `#002D9C`, Dark Mode `#4589FF`)
- **Surface Neutrals (Light / Dark)**:
  - Canvas Background: `#F8FAFC` (Light) / `#0B0F19` (Dark)
  - Elevated Card / Table Surface: `#FFFFFF` (Light) / `#111827` (Dark)
  - Subtle Border: `#E2E8F0` (Light) / `#1E293B` (Dark)
- **Semantic Status Tokens**:
  - `FRESH_VERIFIED`: Emerald `#059669` (Text `#065F46`, Bg `#ECFDF5`)
  - `STALE_WARNING`: Amber `#D97706` (Text `#92400E`, Bg `#FFFBEB`)
  - `DESTRUCTIVE_DANGER`: Crimson `#DC2626` (Text `#991B1B`, Bg `#FEF2F2`)
  - `SPONSORED_BADGE`: Slate Outline `#475569` with mandatory `Sponsored` / `Affiliate Partner` label.

---

## 2. Typography & Tabular Numerals
- **UI & Editorial Sans**: `Inter` / `Pretendard Variable` (for seamless Korean/English CJK rendering).
- **Technical & Numeric Monospace**: `JetBrains Mono` (`font-variant-numeric: tabular-nums` enforced on all vCPU, RAM, Storage, Bandwidth, Price, FX, and Latency cells so columns align to the exact pixel vertically).

---

## 3. Responsive Breakpoints & Touch Ergonomics
- **Mobile (`360px–767px`)**: Sticky bottom filter/compare bar, horizontal swipeable comparison cards with pinned left SKU column, minimum `44x44px` touch targets.
- **Tablet (`768px–1279px`)**: Collapsible filter drawer, 3-SKU side-by-side comparison grid.
- **Desktop (`1280px–1920px+`)**: Persistent left facet sidebar (`280px`), full 12-column Engineer Data-Grid with virtualized rows (`TanStack Table`).
