# CLOUD_DASHBOARD.md — BYOC Multi-Cloud Connection, Provisioning & Server Management Console UX

**Document ID:** `CS-UX-005`

---

## 1. Guided 4-Step Server Provisioning Wizard (`/dashboard/servers/new`)
- **Step 1: Select Connected Account & Capability Verification**: Shows only providers where the user has an active `READ_WRITE` connection and `can_create_server = true`.
- **Step 2: Region, Plan (SKU) & OS Image Selection**: Queries live regional availability (`AVAILABLE` vs `OUT_OF_STOCK`).
- **Step 3: Security & Resilience Options**: Select existing SSH keys, Firewall group, IPv6 toggle, Automated Backups (`+20%`), and Cloud-Init script.
- **Step 4: Mandatory Explicit Cost Approval Gate**:
  - Displays itemized Hourly Cost, Monthly Estimated Cost, Organization Budget Remaining, and Provider Direct Billing Notice (*"You will be billed directly by DigitalOcean on your connected account; CloudScanr charges $0 markup"*).
  - Requires checking `[x] I approve creating this server and incurring estimated charges of $24.00/month on my connected DigitalOcean account` + **Step-Up MFA / Passkey Verification** before the `Create Server` button unlocks.

## 2. Destructive Action Safety Modal (`Delete Server` / `Rebuild Server`)
- Displays server IP, hostname, attached volumes, and warning that data destruction is irreversible.
- Requires typing the exact hostname (e.g., `prod-api-fra-01`) + Step-Up MFA token.
