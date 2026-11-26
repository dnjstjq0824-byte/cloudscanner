# TIRPSCANNER_PIVOT.md — Pivot Decision Record (2026-11-26)

**Document ID:** `TS-DOC-PIVOT-001`  
**Version:** `1.0.0`  
**Date:** `2026-11-26`  
**Status:** Authoritative pivot record for all living docs

---

## 1. Decision

The repository pivots from **CloudScanr** (cloud price comparison + BYOC
multi-cloud management) to **TirpScanner**, a global **SaaS/AI-tool &
software affiliate deal scanner**.

- Old brand: CloudScanr (`cloudscanr.com`, `api.cloudscanr.com`, `www.cloudscanr.com`)
- New brand: TirpScanner
- Name meaning: **Tirp** = Tool Indexing & Ranking Platform
- New domains (apex registered on Cloudflare 2026-11-26, wiring pending; subdomains ASSUMED):
  `tirpscanr.com` (storefront), `api.tirpscanr.com` (Go API origin),
  `www.tirpscanr.com` (www alias)
- New live scope: SaaS/AI-tool catalog, deal scan / price-drop detection,
  affiliate monetization, content hub, admin & affiliate operations.
- Retired scope: cloud-infrastructure catalog, BYOC account connection,
  server provisioning, server management, cloud resale.

## 2. Assumed domains

| Domain | Role | Status |
| :--- | :--- | :--- |
| `tirpscanr.com` | Public storefront (Nuxt on Cloudflare Pages) | **REGISTERED** (wiring pending) |
| `api.tirpscanr.com` | Go API origin (container/VM host) | ASSUMED |
| `www.tirpscanr.com` | www alias / redirect to apex | ASSUMED |

Apex `tirpscanr.com` is **REGISTERED** (Cloudflare, 2026-11-26); confirm
DNS, TLS, and CORS/Pages wiring in `DEPLOY.md` before launch. `api.`/`www.`
and previews (`tirpscanr.pages.dev`) remain ASSUMED.

## 3. Live-vs-retired scope map

| Area | Status | Notes |
| :--- | :--- | :--- |
| SaaS/AI-tool & software catalog | LIVE | Global search, categories, plan specs, availability |
| Deal scan / price-drop detection | LIVE | Coupon/promo tracking, plan normalization, alerts, freshness |
| Affiliate marketplace (`/go/*`, attribution, ledger) | LIVE | Postback/CSV reconciliation, disclosure badges |
| Creator & content hub | LIVE | Comparisons, rankings, reviews, buying guides |
| Admin & affiliate operations | LIVE | Catalog moderation, payout reconciliation, notices, settings |
| Platform (auth, security, testing, DevOps) | LIVE | Toolchain unchanged: Go 1.27, Node 24, pnpm 9.15.0 |
| Cloud-infrastructure catalog (VPS/GPU/bare-metal/TCO) | RETIRED | CloudScanr-era history only |
| Cloud Account Connection (BYOC) | RETIRED | No new connections; credential flows historical |
| Cloud Provisioning (server creation) | RETIRED | Approval/budget/idempotency gates historical |
| Cloud Management (power/resize/backup/firewall/SSH/IP/delete) | RETIRED | `08_CLOUD_MANAGEMENT/`, `/api/v1/cloud/*`, console server pages historical |
| Reseller adapter / direct resale | RETIRED | Never activated; requires executed contract if ever revived |

### Revised adapter law

| Adapter | Status | Role |
| :--- | :--- | :--- |
| Affiliate (`internal/adapters/affiliate`) | LIVE | Tracking links, SubID/ClickID, postbacks, ledger |
| Catalog/Deal-ingestion (`internal/adapters/pricing`) | LIVE | Catalog + deal-feed ingestion, plan normalization (`pricing` name historical) |
| Management/BYOC (`internal/adapters/management`) | RETIRED | Former user-account server ops; no new work |
| Reseller (`internal/adapters/reseller`) | RETIRED | Never activated |

## 4. Frozen-doc policy

The following are **FROZEN** and must never be edited by any lane:

- `11_WORK_ORDERS_0001_TO_1000/` (1,000 executed work orders + index/graph)
- `MANIFEST.json` (frozen blueprint-delivery inventory)
- `FILE_CHECKSUMS.sha256` (frozen integrity record)
- `15_FINAL_DELIVERY/REQUIREMENTS_TRACEABILITY.csv` (frozen audit record)

Where dated cuts matter (e.g., the 2026-11-18 commercialization cut),
**append a dated addendum** instead of rewriting history. All doc edits stay
in English and are docs-only — never edit code, config, or workflows.

## 5. Doc-change log (Lane A, 2026-11-26)

| File | Change |
| :--- | :--- |
| `README.md` | Rewritten for TirpScanner: deal-scanner services, revised adapter law, retired-scope notice, annotated directory map |
| `START_HERE.md` | Rewritten for TirpScanner: deal-scanner execution protocol (Steps 13–15 generalized, provisioning gates retired), updated safety rules |
| `DEPLOY.md` | Domains/branding updated to tirpscanr.com family; toolchain unchanged; pivot note + deal-ingestion worker wording |
| `COMMERCIALIZATION.md` | Dated 2026-11-26 pivot addendum appended (§9); 2026-11-18 cut preserved as history |
| `TIRPSCANNER_PIVOT.md` | Created: pivot decision record, assumed domains, scope map, frozen-doc policy, change log (this file) |
