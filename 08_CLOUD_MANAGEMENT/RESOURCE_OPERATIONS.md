> **RETIRED SCOPE — TirpScanner pivot (2026-11-26):** BYOC account connection, server provisioning, and cloud server management were removed from live scope. This file is CloudScanr-era history, preserved as record; do not build from it. Live scope (SaaS/AI-tool catalog, deal scan / price-drop detection, affiliate monetization) is defined in `TIRPSCANNER_PIVOT.md` at the repository root.

# RESOURCE_OPERATIONS.md — Day-2 Server Lifecycle Operations Specification
**Document ID:** `CS-CM-004`
## 1. Supported Day-2 Operations
- **Power Operations**: `POWER_ON`, `GRACEFUL_SHUTDOWN` (ACPI), `REBOOT`, `POWER_CYCLE` (hard reset).
- **Compute & Disk Resizing**: `RESIZE_CPU_RAM_ONLY` (reversible downgrade allowed) vs `RESIZE_WITH_DISK` (irreversible disk expansion; requires explicit destructive/cost confirmation).
- **Data Protection**: `ENABLE_BACKUPS`, `DISABLE_BACKUPS`, `CREATE_SNAPSHOT`, `RESTORE_SNAPSHOT`.
- **Network & Security**: `ATTACH_FIREWALL`, `ADD_SSH_KEY`, `ASSIGN_FLOATING_IP`.
- **Safe Termination**: `DELETE_SERVER` (requires hostname confirmation + Step-Up MFA + `operation_approvals` token).
