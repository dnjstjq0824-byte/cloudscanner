# CI_CD.md — Standalone GitHub Actions CI/CD & Quality Gates
**Document ID:** `CS-DO-001`
> TirpScanner pivot (2026-11-26): `safety-22-suite` stage retired with BYOC; remaining gates reframed for the deal scanner. See `TIRPSCANNER_PIVOT.md`.
## 1. Pipeline Stages
1. `lint-and-contracts`: `golangci-lint`, `eslint`, `vue-tsc`, `vacuum lint openapi.yaml`, `sqlfluff lint migrations/`.
2. `unit-and-integration`: `go test -v -race -cover ./...` (with `testcontainers-go` PG 16 + OpenSearch 2.17 + Redis 7.4).
3. `deal-affiliate-e2e`: Runs search, deal-scan, affiliate attribution/reconciliation, and security E2E suites.
4. `e2e-playwright`: Desktop + Mobile Playwright suite + Lighthouse SEO/Accessibility gate.
5. `build-and-deploy`: Cosign-signed container images pushed to registry -> Staging rollout -> Independent Audit check -> Production canary rollout.
