# ACCEPTANCE_CRITERIA.md — Global Platform Acceptance Criteria & Definition of Done
**Document ID:** `CS-TST-001`
> TirpScanner pivot (2026-11-26): provisioning safety gate retired; search/affiliate/security E2E gates stay. See `TIRPSCANNER_PIVOT.md`.
## 1. Release Gate Criteria
1. **100% Schema & Contract Parity**: All live PostgreSQL tables and OpenAPI groups match implementation with zero drift.
2. **Zero-Trust Honesty Compliance (Article 29)**: Zero fabricated affiliate approvals, zero unverified conversions, zero plaintext API secrets in DB or logs.
3. **~~22/22 Cloud Provisioning Safety Scenarios~~ (RETIRED)**: replaced by 100% pass on `SEARCH_E2E.md`, `AFFILIATE_E2E.md`, and `SECURITY_E2E.md`.
