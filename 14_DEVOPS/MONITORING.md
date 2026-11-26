# MONITORING.md — OpenTelemetry, Prometheus, Grafana & Alerting Architecture
**Document ID:** `CS-DO-002`
> TirpScanner pivot (2026-11-26): metric namespace renamed to `tirpscanner_*`; BYOC reconciliation metric retired. See `TIRPSCANNER_PIVOT.md`.
## 1. Golden Signals & Custom Metrics
- `tirpscanner_search_latency_ms` (Histogram, SLA: p95 < 50ms)
- `tirpscanner_affiliate_redirect_latency_ms` (Histogram, SLA: p99 < 15ms)
- `tirpscanner_deal_sync_freshness_seconds` (Gauge per vendor/feed, Alert if > 21,600s)
- `tirpscanner_deal_ingestion_quarantines_total` (Counter per vendor/feed)
- `tirpscanner_security_idor_attempts_total` (Counter, PagerDuty P1 if > 5/min)
