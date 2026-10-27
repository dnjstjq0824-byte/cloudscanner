-- Migration 000010: Analytical Materialized Views & High-Performance Composite Indexes

CREATE INDEX idx_specs_vcpu_ram_storage ON product_specs(vcpu_count, ram_mb, storage_type, storage_gb);
CREATE INDEX idx_specs_gpu ON product_specs(gpu_model, gpu_vram_gb, gpu_count) WHERE gpu_count > 0;
CREATE INDEX idx_skus_trgm_name ON product_skus USING gin (display_name gin_trgm_ops);

CREATE MATERIALIZED VIEW mv_daily_affiliate_yield AS
SELECT
    DATE_TRUNC('day', c.clicked_at) AS metric_day,
    c.provider_id,
    p.slug AS provider_slug,
    COUNT(DISTINCT c.id) FILTER (WHERE NOT c.is_bot_filtered) AS valid_clicks,
    COUNT(DISTINCT conv.id) FILTER (WHERE conv.conversion_status = 'APPROVED') AS approved_conversions,
    COALESCE(SUM(comm.normalized_amount_usd) FILTER (WHERE comm.reward_type = 'CASH_COMMISSION' AND comm.ledger_status IN ('APPROVED', 'SETTLED_PAID')), 0) AS approved_cash_commission_usd,
    COALESCE(SUM(comm.normalized_amount_usd) FILTER (WHERE comm.reward_type = 'HOSTING_CREDIT'), 0) AS approved_credit_value_usd
FROM affiliate_clicks c
JOIN cloud_providers p ON p.id = c.provider_id
LEFT JOIN affiliate_conversions conv ON conv.click_id = c.id
LEFT JOIN commissions comm ON comm.conversion_id = conv.id
GROUP BY 1, 2, 3;

CREATE UNIQUE INDEX idx_mv_daily_affiliate_yield ON mv_daily_affiliate_yield(metric_day, provider_id);
