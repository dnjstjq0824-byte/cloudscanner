# TirpScanner Deployment Guide

Two pieces ship separately: the Nuxt frontend goes to **Cloudflare Pages**,
the Go API runs on **any container/VM host** (Fly.io, Render, Hetzner VPS,
or `docker compose` on your own server). The Go API cannot run on Pages/Workers.

## 1. Frontend — Cloudflare Pages (dashboard)

Create a Pages project from this repository with these **Build settings**:

| Setting | Value |
|---|---|
| Framework preset | Nuxt |
| Root directory | _(repo root; leave empty)_ |
| Build command | `npx -y pnpm@9.15.0 --dir web install --frozen-lockfile && npx -y pnpm@9.15.0 --dir web build` |
| Build output directory | `web/dist` |
| Node.js version | 24 (set `NODE_VERSION=24` or via dashboard) |

**Environment variables** (Pages → Settings → Environment Variables):

| Variable | Value |
|---|---|
| `NITRO_PRESET` | `cloudflare-pages` |
| `NODE_VERSION` | `24` |
| `NUXT_PUBLIC_API_BASE` | `https://api.tirpscanr.com` (your deployed Go API origin, no trailing slash) |

Do **not** use `npx wrangler deploy` as the build command: that targets Workers
and fails with "Could not detect a directory containing static files".
Pages builds the Nitro `cloudflare-pages` preset and serves `web/dist`
(static assets plus the `_worker.js` SSR bundle) automatically.

CLI alternative (from a machine with `wrangler login`):

```bash
cd web
NITRO_PRESET=cloudflare-pages pnpm build
pnpm run deploy:pages   # deploys ./dist
```

## 2. Backend — Go API (container/VM)

Development (zero-config; ports bind 127.0.0.1 only):

```bash
docker compose up --build -d
curl http://localhost:8080/api/v1/readiness
```

Production (real secrets required; missing values fail fast):

```bash
cp .env.example .env
# Fill APP_DOMAIN (tirpscanr.com family), DATABASE_URL, AUTH_TOKEN_SECRET
# (32+ bytes), AFFILIATE_SIGNING_SECRET, KMS_KEY_ID, POSTGRES_USER,
# POSTGRES_PASSWORD, POSTGRES_DB, and
# CORS_ALLOWED_ORIGINS=https://tirpscanr.com,https://www.tirpscanr.com
docker compose -f docker-compose.yml -f docker-compose.prod.yml up --build -d
curl https://api.tirpscanr.com/api/v1/readiness
```

The API refuses production boot without a reachable database
(`APP_ENV=production` + failing `DATABASE_URL` exits non-zero), without
a tirpscanr.com-family `APP_DOMAIN`, and without a valid 32+ byte
`AUTH_TOKEN_SECRET` (fail-closed, no fallback). `HTTP_PORT` (default
8080) and `CORS_ALLOWED_ORIGINS` (comma-separated, exact match, no
wildcards) configure the listener and browser access. Retired BYOC
`/cloud/*` routes stay unregistered unless `ENABLE_RETIRED_CLOUD_ROUTES`
is set to true.

## 3. Wiring checklist

1. API serves `/api/v1/readiness` → `{"status":"ready",...}`.
2. `CORS_ALLOWED_ORIGINS` includes the exact Pages origin
   (e.g. `https://tirpscanr.pages.dev` for previews).
3. Pages env `NUXT_PUBLIC_API_BASE` points at the API origin.
4. `/go/{slug}` deal links resolve against the API origin automatically
   (`useApi().baseURL`; the client reads `NUXT_PUBLIC_API_BASE`).

## 4. Routing truth (commercialization cut; domains updated 2026-11-26)

- **Development (default):** `NUXT_PUBLIC_API_BASE` is unset, so the Nuxt
  client issues same-origin calls through the Nitro `/api/**` proxy, whose
  target is `NUXT_API_PROXY_TARGET` (default `http://localhost:8080`).
- **Production:** `NUXT_PUBLIC_API_BASE=https://api.tirpscanr.com` makes the
  client call the API origin directly with `Authorization: Bearer` headers;
  exact-match CORS on the API authorizes the Pages origin. Pages needs no proxy.
- **Outbound deals:** `/go/*` is served only by the Go API, never by Nuxt, so
  deal links are always absolute against the API origin.

## 5. Operations notes

- **Toolchain:** Go `1.27` (module, CI, and `Dockerfile.api`), Node `24`,
  pnpm `9.15.0` (`packageManager` pinned; lockfile v9). Run web installs with
  `npx --yes pnpm@9.15.0 install` to match CI's `--frozen-lockfile` gate.
- **Frontend gates:** `typecheck` + `vitest` + `build` are the only web gates
  (no eslint configured); CI runs all three on every push/PR.
- **Compose initdb:** `docker-compose.yml` bind-mounts `03_DATABASE/migrations/`
  to initdb; the `000012` rollback reference was relocated to
  `03_DATABASE/down_migrations_reference/` so fresh volumes never execute DROP
  statements. Do not move `*.down.sql` files back into `migrations/`.
- **Maintenance mode:** `/api/v1/admin/settings` toggles a 503 gate over user
  traffic; health, readiness, metrics, admin, notices, and site-config stay
  reachable. The storefront shows active `/api/v1/notices` banners.
- **Worker:** `cmd/worker` (the `worker` compose service, dev and prod) harvests
  all catalog/deal-ingestion adapters on boot and every 15 minutes; logs
  provider counts per run. (`pricing` package naming is historical.)
- **Pivot note (2026-11-26):** brand is TirpScanner; `tirpscanr.com` is
  registered (wiring pending) and `api.`/`www.` are ASSUMED. Toolchain is
  unchanged (see §5). BYOC/provisioning/server-management scope is retired.
