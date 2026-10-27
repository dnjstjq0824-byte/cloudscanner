# RESOURCE_OPERATIONS.md — Day-2 Server Lifecycle Operations Specification
**Document ID:** `CS-CM-004`
## 1. Supported Day-2 Operations
- **Power Operations**: `POWER_ON`, `GRACEFUL_SHUTDOWN` (ACPI), `REBOOT`, `POWER_CYCLE` (hard reset).
- **Compute & Disk Resizing**: `RESIZE_CPU_RAM_ONLY` (reversible downgrade allowed) vs `RESIZE_WITH_DISK` (irreversible disk expansion; requires explicit destructive/cost confirmation).
- **Data Protection**: `ENABLE_BACKUPS`, `DISABLE_BACKUPS`, `CREATE_SNAPSHOT`, `RESTORE_SNAPSHOT`.
- **Network & Security**: `ATTACH_FIREWALL`, `ADD_SSH_KEY`, `ASSIGN_FLOATING_IP`.
- **Safe Termination**: `DELETE_SERVER` (requires hostname confirmation + Step-Up MFA + `operation_approvals` token).
