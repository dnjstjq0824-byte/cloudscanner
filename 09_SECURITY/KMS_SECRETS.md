# KMS_SECRETS.md — Envelope Encryption, Key Rotation & Secret Redaction
**Document ID:** `CS-SEC-002`
> TirpScanner pivot (2026-11-26): BYOC provider-token specifics dropped; envelope encryption, rotation, and redaction stay and now protect affiliate/integration credentials. See `TIRPSCANNER_PIVOT.md`.
## 1. Envelope Encryption & Annual Key Rotation
- Every platform credential record (affiliate network API keys, deal-feed API keys, webhook signing secrets) is encrypted with its own unique 256-bit DEK wrapped by KMS KEK `alias/tirpscanner-platform-kek`.
- Automated rotation job re-wraps DEKs under the new KMS key version without exposing plaintext secrets.
- Structured log middleware (`internal/platform/core/module_0055.go`: `MaskingHandler`, `RedactingLoggingMiddleware`) scrubs any string matching secret patterns (affiliate/feed API keys, webhook signatures, `Authorization` / Bearer headers) before writing to stdout.
