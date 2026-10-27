# MONITORING.md — OpenTelemetry, Prometheus, Grafana & Alerting Architecture
**Document ID:** `CS-DO-002`
## 1. Golden Signals & Custom Metrics
- `cloudscanr_search_latency_ms` (Histogram, SLA: p95 < 50ms)
- `cloudscanr_affiliate_redirect_latency_ms` (Histogram, SLA: p99 < 15ms)
- `cloudscanr_pricing_sync_freshness_seconds` (Gauge per provider, Alert if > 21,600s)
- `cloudscanr_byoc_timeout_reconciliations_total` (Counter per provider)
- `cloudscanr_security_idor_attempts_total` (Counter, PagerDuty P1 if > 5/min)
