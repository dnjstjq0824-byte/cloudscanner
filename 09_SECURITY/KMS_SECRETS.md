# KMS_SECRETS.md — Envelope Encryption, Key Rotation & Secret Redaction
**Document ID:** `CS-SEC-002`
## 1. Envelope Encryption & Annual Key Rotation
- Every `provider_credentials_metadata` record is encrypted with its own unique 256-bit DEK wrapped by KMS KEK `alias/cloudscanr-byoc-kek`.
- Automated rotation job (`WO-CS-0585`) re-wraps DEKs under the new KMS key version without exposing plaintext tokens.
- Structured log middleware (`internal/platform/observability/redactor.go`) scrubs any string matching token patterns (`dop_v1_*`, `linode_*`, `scw-*`, Bearer headers) before writing to stdout.
