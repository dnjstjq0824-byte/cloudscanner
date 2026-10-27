# G — Globalization — Multi-Language (i18n), Multi-Currency FX, Regional Tax Policies & Timezones

**Document ID:** `CS-AZ-G`  
**Version:** `1.0.0`  
**Platform:** CloudScanr (`cloudscanr.com` — Independent Global Cloud Comparison, Affiliate & BYOC Management Platform)  
**Responsible Work Orders:** `WO-CS-0265 ~ WO-CS-0280, WO-CS-0771 ~ WO-CS-0790`

---

## 1. 기능 목적 (Functional Purpose)
Deliver a native global experience across major cloud buyer markets (English `en`, Korean `ko`, Japanese `ja`, German `de`, Spanish `es`) with accurate multi-currency conversion (`USD`, `EUR`, `KRW`, `JPY`, `GBP`, `SGD`, `AUD`), country-specific indirect tax rules (VAT/GST/Consumption Tax), and localized formatting.

---

## 2. 사용자 역할 (User Roles & Actors)
Globalization Engineer, Pricing Engine Agent, Content Localization Lead.

---

## 3. 핵심 사용자 시나리오 (User Scenarios)
1. A user visiting from Seoul (`ko-KR`) sees UI labels in Korean, prices in `KRW (₩)` alongside original `USD ($)`/`EUR (€)` provider billing currency, and a toggle for 10% Korean VAT.
2. A German user (`de-DE`) comparing Hetzner (`EUR`) and DigitalOcean (`USD`) sees unified `EUR` comparison figures with 19% German MwSt (VAT) clearly itemized and `hreflang="de"` canonical links.

---

## 4. 기술 스택 (Technology Stack)
`@nuxtjs/i18n` with lazy-loaded JSON locale bundles, Intl.NumberFormat / Intl.DateTimeFormat, PostgreSQL `fx_rates` & `tax_policies`, ECB / OpenExchangeRates / Bank of Korea FX ingestion worker.

---

## 5. 데이터베이스 구조 (Database Schema & Tables)
**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`fx_rates`, `tax_policies`, `articles`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)
**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**
- `GET /api/v1/pricing/fx-rates`
- `GET /api/v1/pricing/tax-policies`
- `PATCH /api/v1/users/me/preferences`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. UI/UX 설계 (User Interface & Experience)
Header Locale & Currency Selector (`Language: EN/KO/JA/DE/ES`, `Currency: USD/EUR/KRW/JPY/GBP`, `Tax: Excl. Tax / Incl. Local Tax`) with instant reactive recalculation and tooltip showing exact FX rate + timestamp.

---

## 8. 상세 처리 로직 (Core Processing Logic)
1. Dual-Currency Display Rule: Always show the converted target currency for cross-provider sorting AND the provider's native billing currency in parentheses (e.g., `₩33,200/월 (원청구 $24.00 USD)`) so users are never surprised by credit card FX spreads.
2. SEO `hreflang` Matrix: Subpath routing (`/en/...`, `/ko/...`, `/ja/...`, `/de/...`, `/es/...`) with bidirectional `<link rel="alternate" hreflang="...">` and `x-default` tags.

---

## 9. 보안 통제 (Security & Access Control)
Locale parameter sanitization against path traversal; tax calculation disclaimers stating that final tax depends on user billing address and B2B VAT ID status at the provider.

---

## 10. 실패 및 예외 처리 (Failure Modes & Resilience)
Fallback to `en-US` and `USD` if a translation key or FX pair is missing, logging a warning metric without breaking page rendering.

---

## 11. 작업 간 의존성 (Inter-Module Dependencies)
Integrates with Estimation Engine (`E_ESTIMATION_ENGINE.md`) and Frontend (`F_FRONTEND.md`).

---

## 12. 테스트 전략 (Testing & Verification Strategy)
Locale completeness linter (verifying 100% key parity across `en`, `ko`, `ja`, `de`, `es`), VAT/GST math verification tests, and `hreflang` sitemap XML validation.

---

## 13. 완료 승인 기준 (Acceptance Criteria)
- [x] 100% translation key parity across primary locales (`en`, `ko`, `ja`, `de`, `es`).
- [x] Every converted price displays the exact FX rate timestamp and original billing currency.
- [x] Valid `hreflang` tags rendered on all localized public routes.

---

## 14. 검증 증빙 (Required Verification Evidence)
i18n key parity check output, localized SSR HTML snapshots, and FX conversion test logs.

---

## 15. 실제 구현 담당 Work Order (Assigned Work Orders)
- **Work Order Range:** `WO-CS-0265 ~ WO-CS-0280, WO-CS-0771 ~ WO-CS-0790`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
