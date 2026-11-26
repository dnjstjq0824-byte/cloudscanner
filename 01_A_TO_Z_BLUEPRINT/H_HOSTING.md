# H — Hosting — Independent Platform Infrastructure, Kubernetes Topology, CDN/WAF & Disaster Recovery

**Document ID:** `CS-AZ-H` (stable ID, unchanged since CloudScanr era)  
**Version:** `2.0.0` (TirpScanner pivot reframe, 2026-11-26)  
**Platform:** TirpScanner (`tirpscanr.com` — Global SaaS/AI-Tool & Software Affiliate Deal Scanner)  
**Pivot Note (2026-11-26):** TirpScanner's own hosting is unchanged in kind; BYOC egress-allowlist rationale retired (no user-cloud API calls remain). See `TIRPSCANNER_PIVOT.md`.  
**Responsible Work Orders:** `WO-CS-0075 ~ WO-CS-0100, WO-CS-0830 ~ WO-CS-0850`

---

## 1. 기능 목적 (Functional Purpose)

Specify TirpScanner's completely independent production hosting infrastructure, network topology, DNS/TLS management, CDN/WAF rules, object storage, and backup/disaster recovery procedures.

---

## 2. User Roles & Actors

DevOps Agent, Platform SRE, Security Architect.

---

## 3. User Scenarios

1. Production traffic spikes 10x during a viral SaaS price-hike announcement; Cloudflare Edge caches SSR comparison pages while Kubernetes HPA scales `tirpscanner-api` and `tirpscanner-web` pods from 3 to 20 replicas within 45 seconds.
2. Automated PostgreSQL WAL archiving streams continuous Point-In-Time Recovery (PITR) backups to S3/R2; a daily automated restore drill verifies RPO < 5 minutes and RTO < 30 minutes.

---

## 4. Technology Stack

Terraform / OpenTofu, Kubernetes (EKS / DOKS), Cloudflare DNS + Enterprise WAF + CDN + R2 Object Storage, PostgreSQL 16 HA (Patroni / Managed HA Multi-AZ), OpenSearch 2.17 3-node cluster, Redis 7.4 Sentinel, Prometheus + Grafana + Loki + Tempo.

---

## 3. Database Schema & Tables

**Target PostgreSQL 16 Tables (`03_DATABASE/migrations/`):**  
`system_settings`, `background_jobs`, `audit_logs`

- Enforces strict referential integrity (`FOREIGN KEY`), `UUID` / `ULID` identifiers, `NUMERIC(18,8)` monetary precision, and `TIMESTAMPTZ` UTC timestamps.
- Multi-tenant tables enforce PostgreSQL Row-Level Security (`RLS`) keyed on `organization_id`.

---

## 6. API 구조 (API Contracts & Endpoints)

**Primary OpenAPI 3.1 Endpoints (`04_API_SPECIFICATION/openapi.yaml`):**

- `GET /api/v1/health`
- `GET /api/v1/readiness`
- `GET /metrics`

- Standard JSON envelope (`data`, `meta`, `error`), RFC 9457 Problem Details for errors, rate-limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`), and `Idempotency-Key` support on all mutating `POST`/`PATCH`/`DELETE` routes.

---

## 7. User Interface & Experience

Admin Infrastructure & Deployment Monitor (`/admin/deployments`) displaying pod health, DB replica lag, OpenSearch cluster status (Green/Yellow/Red), Redis queue depth, and backup verification status.

---

## 8. Core Processing Logic

1. Dedicated VPC (`10.40.0.0/16`) with public ingress subnets, private application subnets (`api`, `web`, `worker`), and isolated data subnets (`postgres`, `opensearch`, `redis`) with zero public internet ingress.
2. NAT Gateway with static egress IPs for catalog/deal-ingestion workers so affiliate networks and feed sources with IP allowlists can optionally whitelist TirpScanner's egress IPs.

---

## 9. Security & Access Control

Cloudflare WAF OWASP Core Rule Set + Rate Limiting on `/api/v1/auth/*` and `/go/*`; Kubernetes NetworkPolicies blocking lateral pod movement; KMS hardware-backed root keys.

---

## 10. Failure Modes & Resilience

Multi-AZ deployment across 3 availability zones; automated pod disruption budgets (`PDB`); automated database failover via Patroni/Managed HA in <30s.

---

## 11. Inter-Module Dependencies

Supports Architecture (`A_ARCHITECTURE.md`) and DevOps (`14_DEVOPS/`).

---

## 12. Testing & Verification Strategy

Terraform `plan`/`validate`, Kubernetes manifest `kubeconform` validation, and automated PITR backup restore verification script.

---

## 13. Acceptance Criteria

- [x] Zero shared infrastructure, VPC peering, or database instances with TreeScan.
- [x] Static outbound NAT IPs documented and stable for feed/API IP allowlisting.
- [x] Automated daily PostgreSQL and OpenSearch backup verification passes with RPO < 5m and RTO < 30m.

---

## 14. Required Verification Evidence

Terraform validation logs, Kubernetes network policy manifests, and DR drill execution logs.

---

## 15. Assigned Work Orders

- **Work Order Range:** `WO-CS-0075 ~ WO-CS-0100, WO-CS-0830 ~ WO-CS-0850`
- **Cross-Reference Index:** `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`
