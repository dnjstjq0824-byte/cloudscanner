# DEPLOYMENT.md — Production Topology, Multi-Environment Strategy & Zero-Downtime Rollouts

**Document ID:** `CS-ARCH-004`

> **Pivot addendum (2026-11-26):** Domains rebranded to TirpScanner (`tirpscanr.com`, `api.tirpscanr.com`, `www.tirpscanr.com` — apex registered on Cloudflare; DNS wiring pending). BYOC execution targets removed from every environment; live external surface is SaaS catalog/deal feeds plus affiliate links.

---

## 1. Environment Matrix

| Environment | Domain | Database & Search | External API Mode | Purpose |
| :--- | :--- | :--- | :--- | :--- |
| `local` | `localhost:3000` / `:8080` | Docker Compose (PG 16, OS 2.17, Redis 7.4) | Deterministic Deal-Feed Simulator (`httptest`) | Local developer & agent execution |
| `ci` | Ephemeral Container Network | `testcontainers-go` ephemeral instances | Deterministic Simulator + Recorded Fixtures | Automated Unit, Integration, E2E & Safety Suite |
| `staging` | `staging.tirpscanr.com` | Managed Staging PG 16 + OS 2.17 | Read-Only Live Deal Feeds + Sandbox Affiliate Links | Pre-release QA & Independent Audit |
| `production` | `tirpscanr.com` | Multi-AZ HA PG 16 (Primary + 2 Replicas) + 3-Node OS 2.17 | Live Deal Feeds + Live Affiliate Links | Global Production Service |

---

## 2. Zero-Downtime Deployment Sequence
1. **Expand-Contract Database Migrations**: `cs-migrator` runs backward-compatible DDL (`ADD COLUMN NULLABLE`, `CREATE INDEX CONCURRENTLY`) before rolling application pods.
2. **Rolling Pod Update**: Kubernetes `RollingUpdate` (`maxUnavailable: 0`, `maxSurge: 25%`) waits for `/api/v1/readiness` (verifying PG + Redis connectivity) before shifting traffic.
3. **Graceful Worker Drain**: On `SIGTERM`, `cs-worker-deals` and `cs-worker-affiliate` stop pulling new jobs and wait up to 90 seconds for active feed/API calls and state recordings to complete before exiting.
