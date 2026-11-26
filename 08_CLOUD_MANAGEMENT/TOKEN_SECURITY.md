> **RETIRED SCOPE — TirpScanner pivot (2026-11-26):** BYOC account connection, server provisioning, and cloud server management were removed from live scope. This file is CloudScanr-era history, preserved as record; do not build from it. Live scope (SaaS/AI-tool catalog, deal scan / price-drop detection, affiliate monetization) is defined in `TIRPSCANNER_PIVOT.md` at the repository root.

# TOKEN_SECURITY.md — KMS Envelope Encryption, Memory Zeroing & Token Lifecycle
**Document ID:** `CS-CM-002`
## 1. Cryptographic Storage Specification (`provider_credentials_metadata`)
- **Encryption Algorithm**: `AES-256-GCM` (96-bit random nonce, 128-bit authentication tag).
- **Key Hierarchy**:
  - Root **KEK (Key Encryption Key)** resides inside AWS KMS / HashiCorp Vault and never leaves the HSM.
  - Each `cloud_connections` record generates a unique 256-bit **DEK (Data Encryption Key)**.
  - **Additional Authenticated Data (AAD)** is set to `fmt.Sprintf("%s:%s:%s", orgID, connectionID, providerID)`.
- **Runtime Decryption & Memory Hygiene**:
  - Decryption occurs only inside `cs-worker-cloud` immediately before an outbound provider API call.
  - A `defer ZeroBytes(plaintextToken)` call overwrites the byte slice with `0x00` immediately after the HTTP request finishes.
