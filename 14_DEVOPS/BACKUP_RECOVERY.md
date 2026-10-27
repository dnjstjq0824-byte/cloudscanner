# BACKUP_RECOVERY.md — PITR Database Backups, OpenSearch Snapshots & Disaster Recovery
**Document ID:** `CS-DO-003`
## 1. RPO & RTO Targets
- **PostgreSQL 16**: Continuous WAL archiving via `pgBackRest` to immutable Object Storage (S3/R2 Object Lock) + daily full base backups. **RPO < 5 minutes**, **RTO < 30 minutes**.
- **OpenSearch 2.17**: Hourly incremental snapshot to S3/R2; can also be deterministically rebuilt from PostgreSQL `product_skus` + `price_snapshots` in < 8 minutes.
