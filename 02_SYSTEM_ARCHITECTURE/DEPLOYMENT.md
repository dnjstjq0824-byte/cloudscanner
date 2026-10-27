# DEPLOYMENT.md — Production Topology, Multi-Environment Strategy & Zero-Downtime Rollouts

**Document ID:** `CS-ARCH-004`

---

## 1. Environment Matrix

| Environment | Domain | Database & Search | External Provider API Mode | Purpose |
| :--- | :--- | :--- | :--- | :--- |
| `local` | `localhost:3000` / `:8080` | Docker Compose (PG 16, OS 2.17, Redis 7.4, Local KMS Mock) | Deterministic Provider Simulator (`httptest`) | Local developer & agent execution |
| `ci` | Ephemeral Container Network | `testcontainers-go` ephemeral instances | Deterministic Simulator + Recorded Fixtures | Automated Unit, Integration, E2E & Safety Suite |
| `staging` | `staging.cloudscanr.com` | Managed Staging PG 16 + OS 2.17 | Read-Only Live Pricing APIs + Sandbox Test Cloud Accounts (Budget Capped $5) | Pre-release QA & Independent Audit |
| `production` | `cloudscanr.com` | Multi-AZ HA PG 16 (Primary + 2 Replicas) + 3-Node OS 2.17 | Live Pricing APIs + Live Affiliate Links + User-Authorized BYOC | Global Production Service |

---

## 2. Zero-Downtime Deployment Sequence
1. **Expand-Contract Database Migrations**: `cs-migrator` runs backward-compatible DDL (`ADD COLUMN NULLABLE`, `CREATE INDEX CONCURRENTLY`) before rolling application pods.
2. **Rolling Pod Update**: Kubernetes `RollingUpdate` (`maxUnavailable: 0`, `maxSurge: 25%`) waits for `/api/v1/readiness` (verifying PG + Redis connectivity) before shifting traffic.
3. **Graceful Worker Drain**: On `SIGTERM`, `cs-worker-cloud` stops pulling new jobs and waits up to 90 seconds for active provider API calls and state recordings to complete before exiting.
