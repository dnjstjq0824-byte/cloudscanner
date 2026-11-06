# CloudScanr Deployment Guide

Two pieces ship separately: the Nuxt frontend goes to **Cloudflare Pages**,
the Go API runs on **any container/VM host** (Fly.io, Render, Hetzner VPS,
or `docker compose` on your own server). The Go API cannot run on Pages/Workers.

## 1. Frontend — Cloudflare Pages (dashboard)

Create a Pages project from this repository with these **Build settings**:

| Setting | Value |
|---|---|
| Framework preset | Nuxt |
| Root directory | _(repo root; leave empty)_ |
| Build command | `pnpm --dir web install --frozen-lockfile && pnpm --dir web build` |
| Build output directory | `web/dist` |
| Node.js version | 24 (set `NODE_VERSION=24` or via dashboard) |

**Environment variables** (Pages → Settings → Environment Variables):

| Variable | Value |
|---|---|
| `NITRO_PRESET` | `cloudflare-pages` |
| `NODE_VERSION` | `24` |
| `NUXT_PUBLIC_API_BASE` | `https://api.cloudscanr.com` (your deployed Go API origin, no trailing slash) |

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

```bash
cp .env.example .env
# Fill DATABASE_URL, AUTH_TOKEN_SECRET, AFFILIATE_SIGNING_SECRET, KMS_KEY_ID,
# and CORS_ALLOWED_ORIGINS=https://cloudscanr.com,https://www.cloudscanr.com
docker compose up --build -d
curl https://api.cloudscanr.com/api/v1/readiness
```

The API refuses production boot without a reachable database
(`APP_ENV=production` + failing `DATABASE_URL` exits non-zero).
`HTTP_PORT` (default 8080) and `CORS_ALLOWED_ORIGINS` (comma-separated,
exact match, no wildcards) configure the listener and browser access.

## 3. Wiring checklist

1. API serves `/api/v1/readiness` → `{"status":"ready",...}`.
2. `CORS_ALLOWED_ORIGINS` includes the exact Pages origin
   (e.g. `https://cloudscanr.pages.dev` for previews).
3. Pages env `NUXT_PUBLIC_API_BASE` points at the API origin.
4. `/go/{slug}` deal links resolve against the API origin automatically.
