# CI_CD.md — Standalone GitHub Actions CI/CD & Quality Gates
**Document ID:** `CS-DO-001`
## 1. Pipeline Stages
1. `lint-and-contracts`: `golangci-lint`, `eslint`, `vue-tsc`, `vacuum lint openapi.yaml`, `sqlfluff lint migrations/`.
2. `unit-and-integration`: `go test -v -race -cover ./...` (with `testcontainers-go` PG 16 + OpenSearch 2.17 + Redis 7.4).
3. `safety-22-suite`: Runs all 22 Article 24 BYOC provisioning fault-injection tests.
4. `e2e-playwright`: Desktop + Mobile Playwright suite + Lighthouse SEO/Accessibility gate.
5. `build-and-deploy`: Cosign-signed container images pushed to registry -> Staging rollout -> Independent Audit check -> Production canary rollout.
