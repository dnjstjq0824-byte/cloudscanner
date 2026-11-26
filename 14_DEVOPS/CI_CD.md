# CI_CD.md — Standalone GitHub Actions CI/CD & Quality Gates
**Document ID:** `CS-DO-001`
> TirpScanner pivot (2026-11-26): `safety-22-suite` stage retired with BYOC; remaining gates reframed for the deal scanner. See `TIRPSCANNER_PIVOT.md`.
## 1. Pipeline Stages (as implemented in `.github/workflows/ci.yml`)
1. `backend-gate`: TreeScan-coupling grep (Article 1) -> `go mod tidy` check -> `go test -race ./internal/... ./cmd/...` -> `go build` api + worker.
2. `frontend-gate`: `pnpm install --frozen-lockfile` -> `typecheck` (vue-tsc) -> `vitest` -> `pnpm build`.

## 2. Planned Stages (not yet implemented; target dates TBD)
3. `lint-and-contracts`: `golangci-lint`, `eslint`, `vacuum lint openapi.yaml`, `sqlfluff lint migrations/`.
4. `unit-and-integration`: `go test -race -cover ./...` with `testcontainers-go` (PG 16 + OpenSearch 2.17 + Redis 7.4).
5. `deal-affiliate-e2e`: search, deal-scan, affiliate attribution/reconciliation, and security E2E suites.
6. `e2e-playwright`: Desktop + Mobile Playwright suite + Lighthouse SEO/Accessibility gate.
7. `build-and-deploy`: Cosign-signed container images -> Staging rollout -> Independent Audit check -> Production canary rollout.
