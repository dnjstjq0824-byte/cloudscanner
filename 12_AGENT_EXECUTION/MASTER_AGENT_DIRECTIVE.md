# MASTER_AGENT_DIRECTIVE.md — Autonomous Multi-Agent Orchestration Directive (Article 22)
**Document ID:** `CS-AG-001`
## 1. Supreme Operating Principles for All 14 Agents
1. **Contract-First Execution**: No agent may alter `03_DATABASE/migrations/*.sql` or `04_API_SPECIFICATION/openapi.yaml` without a Lead Architect + Database/Backend Agent Contract Lock.
2. **Strict File Ownership**: Every repository path is owned by a single primary agent (`AGENT_ROLE_MATRIX.md`) to prevent parallel Git merge conflicts.
3. **Zero Unverified Claims**: A Work Order (`WO-CS-XXXX`) transitions to `VERIFIED_COMPLETE` only after Unit + Integration/E2E + Independent Audit Agent verification passes.
