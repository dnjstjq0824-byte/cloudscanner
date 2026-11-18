# CloudScanr Commercialization Ledger (2026-11-18 cut)

Authoritative implementation record for the production-readiness cut. Blueprint
directories (`00`–`15`) remain the design SSOT; where they predate this cut, the
implemented/planned status below wins. Per-path API status also lives in
`04_API_SPECIFICATION/openapi.yaml`.

## 1. Cut summary

- **Auth E2E repaired:** the Nuxt client now sends `Authorization: Bearer` +
  `X-Organization-Id` from session cookies, auto-attaches `Idempotency-Key` to
  mutations, and exposes the API origin (`useApi().baseURL`) so all `/go` deal
  links resolve. Development uses the same-origin Nitro proxy (target
  overridable via `NUXT_API_PROXY_TARGET`); any explicit `NUXT_PUBLIC_API_BASE`
  switches the client to direct absolute calls with exact-match CORS.
- **Foreign residue removed:** 59-file Nitro video-platform mock backend
  (`web/server/`), `web/types/{video,studio}.ts`, `web/composables/useFormat.ts`,
  and the `hls.js` dependency are deleted. The Go API is the only backend.
- **Servers domain shipped:** `POST /api/v1/cloud/servers` executes the full
  Article 24 gate (approval existence/expiry/signature/MFA binding, connection
  active + read-write, budget hard-block) and creates a single-use managed
  record; `GET` (one/list) and double-gated `DELETE` (hostname + MFA) complete
  CRUD. Provider-side execution stays pending live management adapters and is
  reported as such (`provider_sync.pending`), never faked as RUNNING.
- **Admin console shipped (API + UI):** overview KPIs, member list with
  status/role mutations (memory + PostgreSQL stores), provider status
  moderation, append-only audit trail (2,000-entry ring; all billable/admin
  mutations recorded), site notices CRUD with public banner feed, and runtime
  settings including a real maintenance-mode 503 gate.
- **Idempotency completed:** the middleware now captures outcomes and replays
  them (`X-Cache-Lookup: HIT-IDEMPOTENCY`); cloud/admin mutations require a
  16–64 char key. Auth, public catalog, and read endpoints are exempt.
- **KMS restart-safe:** sealing keys derive from the server secret
  (`NewMockKMSServiceWithSecret`); sealed credentials reopen after restarts.
  Managed KMS (AWS/Vault) remains the production follow-up behind the same
  interface. Safety-gate HMAC keys derive from the secret as well.
- **Worker wired:** `cmd/worker` registers all five pricing adapters and runs
  `HarvestAll` on boot plus every 15 minutes with structured run logs.
- **Affiliate honesty fixes:** click country comes from edge headers
  (`CF-IPCountry`/`X-Geo-Country`) with `UNKNOWN` fallback instead of a
  hardcoded `US`; the base adapter implements generic HMAC-SHA256 postback
  verification and real CSV reconciliation parsing (both tested).
- **Compose initdb hazard removed:** `000012_down_migrations_reference.down.sql`
  moved to `03_DATABASE/down_migrations_reference/` (content unchanged); fresh
  volumes no longer execute DROPs. The API runner only applies `*.up.sql`.
- **Contract alignment:** `POST /api/v1/pricing/estimate` is canonical with
  `/api/v1/estimate/tco` kept as alias; `GET /api/v1/pricing/fx-rates` reads
  PostgreSQL seed data (explicit 503 without a database); budgets use PUT;
  tokens documented as v1 HMAC-SHA256; router documented as `net/http`
  ServeMux; `Idempotency-Key` length fixed to 16–64.
- **UI/UX overhaul:** token-driven light/dark themes (dark default, persisted,
  pre-paint script, Pretendard loaded), 10 shared components, responsive
  layouts (mobile menu + drawer navigation at all breakpoints), live homepage
  statistics, notice banners, toasts, confirm modals, and full admin pages.
- **Toolchain unified:** Go 1.27 (module/CI/Docker), Node 24, pnpm 9.15.0
  pinned with lockfile v9.

## 2. Architecture truth (as built)

| Concern | Built state | Code anchor |
|---|---|---|
| HTTP router | `net/http` ServeMux, manual registration | `cmd/api/main.go` |
| Tokens | v1 HMAC-SHA256, 15m TTL, 32-byte min secret | `internal/domain/identity/token.go` |
| Identity store | PostgreSQL when connected, else in-memory | `cmd/api/main.go`, `pg_store.go` |
| Catalog/search/pricing | In-memory seed (18 providers, 11 SKUs) + NL parser + BestValue | `internal/domain/catalog`, `search` |
| BYOC state | In-memory + sealed credentials, tenant-checked | `internal/domain/management` |
| Admin/audit/notices | In-memory + audit ring (PG migration planned) | `internal/domain/admin`, `internal/platform/audit` |
| FX rates | PostgreSQL `fx_rates` seed, 503 without DB | `internal/domain/estimation/fx.go` |
| Middleware order | Panic → Security → RequestID → AccessLog → Maintenance → RateLimit → mux | `cmd/api/main.go` |
| Frontend | Nuxt 3 SSR/SWR public, SPA console/admin, Pinia auth, token CSS system | `web/` |

## 3. Endpoint ledger (implemented)

Auth: `POST /api/v1/auth/register`, `POST /api/v1/auth/login`,
`GET /api/v1/auth/me` (Bearer). Core: `GET /api/v1/health`,
`GET /api/v1/readiness`, `GET /metrics`. Catalog/search:
`GET /api/v1/providers`, `GET /api/v1/providers/{slug}`,
`GET /api/v1/catalog/categories`, `GET /api/v1/catalog/skus{,/{id}}`,
`GET /api/v1/search/skus`, `GET /api/v1/rankings/best-value`. Pricing:
`POST /api/v1/pricing/estimate`, `POST /api/v1/estimate/tco` (alias),
`GET /api/v1/pricing/fx-rates`. Affiliate: `GET /go/{slug}`,
`GET /api/v1/affiliate/disclosure`, `POST /api/v1/webhooks/affiliate/{network}`.
BYOC (auth + idempotent mutations): `/api/v1/cloud/connections` (GET/POST),
`/api/v1/cloud/connections/{id}` (GET/DELETE), `/api/v1/cloud/approvals`
(POST), `/api/v1/cloud/budgets` (GET/PUT), `/api/v1/cloud/resources` (GET),
`/api/v1/cloud/servers` (GET/POST), `/api/v1/cloud/servers/{id}`
(GET/DELETE). Admin (ADMIN+): `/api/v1/admin/overview`, `/api/v1/admin/users`
(GET), `/api/v1/admin/users/{id}` (PATCH), `/api/v1/admin/providers/{slug}`
(PATCH), `/api/v1/admin/audit-logs` (GET), `/api/v1/admin/notices`
(GET/POST), `/api/v1/admin/notices/{id}` (GET/PUT/DELETE),
`/api/v1/admin/settings` (GET/PUT). Public: `GET /api/v1/notices`,
`GET /api/v1/site-config`.

Planned (Phase 4): email verify, MFA/TOTP enrolment, logout/refresh, users/orgs
self-service, comparisons, promotions, costs/summary, price-alerts, analytics,
content, standalone audit API, notifications, server power/resize actions,
pricing quarantine, affiliate CSV reconcile endpoint, OpenSearch/pg_trgm
retrieval, Asynq fleet, OTel/Grafana/PagerDuty, K8s/Terraform topology.

## 4. UI inventory (22 pages, 3 layouts, 10 components)

Public: `/`, `/search`, `/calculator`, `/rankings`, `/providers`,
`/providers/[slug]`, `/auth/login`, `/auth/register`, `/affiliate-disclosure`,
`/privacy`, `/terms`. Console: `/dashboard`, `/dashboard/connections`,
`/dashboard/servers`, `/dashboard/provision`, `/dashboard/budgets`. Admin:
`/admin`, `/admin/users`, `/admin/providers`, `/admin/audit-logs`,
`/admin/notices`, `/admin/settings`. Layouts: `public` (mobile menu),
`dashboard` + `admin` (drawer under `lg`). Components: `PageHeader`,
`StatCard`, `StatusBadge`, `LoadingSkeleton`, `ErrorState`, `EmptyState`,
`ConfirmModal`, `ThemeToggle`, `NoticeBanner`, `ToastHost`; composables
`useApi`, `useTheme`, `useToast`; Pinia `auth` store; `auth`/`admin`/`guest`
middleware. Breakpoints verified: 360 / 768 / 1280 (see §6).

## 5. Data pipeline E2E

Register → login → Bearer session → connect cloud (KMS-sealed) → set budget →
estimate TCO → sign HMAC approval → execute provision (idempotent, gates
enforced) → inventory + audit trail → admin moderation + notices → storefront
banners. Pricing harvest runs on worker boot + 15-minute cadence. Seeds:
11 DDL migrations + provider/role/tax/FX seeds; in-memory catalog (18
providers/11 SKUs), 3 affiliate links. No screen renders empty without an
explicit loading/error/empty state.

## 6. Verification evidence (this cut)

- `go build ./...` — pass.
- `go test ./...` — pass (all packages; routing suite extended to 10
  scenarios: alias, provision/delete gates, admin console; new suites for
  admin, audit, KMS restart, affiliate HMAC/CSV/geo).
- `pnpm run typecheck` (vue-tsc, strict) — pass.
- `pnpm run build` (Nuxt production) — pass.
- Playwright chromium screenshots at 360×800 / 768×1024 / 1280×800 across
  public, console, and admin routes — captured in Phase 5 (see screenshot
  log in the final report); zero horizontal overflow asserted.
- `openapi.yaml` parses (js-yaml) after contract alignment edits.

## 7. Documentation scope decision

Updated living docs: `README.md`, `DEPLOY.md`, `COMMERCIALIZATION.md` (new),
`04_API_SPECIFICATION/openapi.yaml`, `03_DATABASE/{ERD,INDEX_DESIGN}.md`,
`02_SYSTEM_ARCHITECTURE/COMPONENTS.md` (migration path),
`15_FINAL_DELIVERY/{DELIVERY_CHECKLIST,FINAL_MASTER_REPORT}.md` (dated
addenda; original counts preserved as historical record). Frozen by design:
`11_WORK_ORDERS_0001_TO_1000/` (1,000 executed work orders + index/graph),
`MANIFEST.json` / `FILE_CHECKSUMS.sha256` (frozen 1,120-entry blueprint delivery
record — implementation files were never in its scope),
`15_FINAL_DELIVERY/REQUIREMENTS_TRACEABILITY.csv` (audit record; covered by the
checklist/report addenda), and the remaining A–Z/spec narratives, whose delta
list is §3's planned set. Rewriting executed history would falsify the audit
trail; the ledger above is the forward SSOT.

## 8. Pre-deploy checklist for operators

1. Provision PostgreSQL 16 and set `DATABASE_URL` (+ `APP_ENV=production`).
2. Set `AUTH_TOKEN_SECRET` (≥32 chars), `KMS_KEY_ID`, and exact-match
   `CORS_ALLOWED_ORIGINS` for the Pages origin.
3. Deploy the API (compose/VM), verify `/api/v1/readiness` → `ready`.
4. Set Pages env `NUXT_PUBLIC_API_BASE` to the API origin and deploy `web/`.
5. Register the first OWNER account, open `/admin`, publish a notice, and
   confirm the audit trail records the mutation.
6. Leave maintenance mode off except during windows; confirm `/api/v1/notices`
   banners render on the storefront.
