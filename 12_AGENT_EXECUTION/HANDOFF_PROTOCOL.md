# HANDOFF_PROTOCOL.md — Inter-Agent Contract Change & Handoff Protocol
**Document ID:** `CS-AG-004`
## 1. Handoff State File & Contract Change Request (CCR)
1. When an upstream agent completes a prerequisite Work Order (e.g., `WO-CS-0551`), it records the Git commit SHA, migration version, and test log artifact path in the handoff registry.
2. If a downstream agent discovers a missing field in `openapi.yaml` or `migrations/*.sql`, it MUST NOT patch the file silently; it emits a `CCR` ticket reviewed by Lead Architect + Database/Backend Agent and runs `verify_blueprint.py` to ensure zero broken references.
