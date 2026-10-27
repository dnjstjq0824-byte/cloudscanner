-- Migration 000004: Pricing & 36-Dimension Normalization Domain
-- Tables: product_skus, product_specs, price_snapshots, price_components, billing_terms, promotions, fx_rates, tax_policies, price_sync_jobs

CREATE TABLE product_skus (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id UUID NOT NULL REFERENCES cloud_providers(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES provider_products(id) ON DELETE CASCADE,
    region_id UUID NOT NULL REFERENCES provider_regions(id) ON DELETE CASCADE,
    external_sku_code VARCHAR(160) NOT NULL, -- e.g., s-2vcpu-4gb, cpx21, g6-standard-2, t4g.medium
    display_name VARCHAR(200) NOT NULL,
    category product_category_enum NOT NULL,
    availability_status VARCHAR(32) NOT NULL DEFAULT 'AVAILABLE' CHECK (availability_status IN ('AVAILABLE', 'LOW_STOCK', 'OUT_OF_STOCK', 'DEPRECATED')),
    freshness_status freshness_status_enum NOT NULL DEFAULT 'FRESH',
    value_score NUMERIC(6,2) DEFAULT 0.00,
    official_order_url TEXT NOT NULL,
    last_price_verified_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (provider_id, external_sku_code, region_id)
);

CREATE TABLE product_specs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku_id UUID NOT NULL REFERENCES product_skus(id) ON DELETE CASCADE UNIQUE,
    vcpu_count INT NOT NULL DEFAULT 1,
    cpu_architecture VARCHAR(32) NOT NULL DEFAULT 'x86_64' CHECK (cpu_architecture IN ('x86_64', 'arm64')),
    cpu_allocation_type VARCHAR(32) NOT NULL DEFAULT 'SHARED' CHECK (cpu_allocation_type IN ('SHARED', 'BURSTABLE', 'DEDICATED')),
    cpu_brand_model VARCHAR(160),
    cpu_base_clock_ghz NUMERIC(5,2),
    ram_mb INT NOT NULL,
    storage_type VARCHAR(32) NOT NULL DEFAULT 'NVME' CHECK (storage_type IN ('NVME', 'SSD', 'HDD', 'NETWORK_BLOCK', 'NONE')),
    storage_gb INT NOT NULL DEFAULT 0,
    disk_iops_read INT,
    disk_iops_write INT,
    disk_throughput_mbps INT,
    included_bandwidth_gb INT NOT NULL DEFAULT 0,
    bandwidth_billing_direction VARCHAR(32) NOT NULL DEFAULT 'OUTBOUND_ONLY' CHECK (bandwidth_billing_direction IN ('OUTBOUND_ONLY', 'BIDIRECTIONAL', 'UNMETERED_PORT')),
    network_speed_gbps NUMERIC(6,2) DEFAULT 1.00,
    ipv4_included_count INT NOT NULL DEFAULT 1,
    ipv6_supported BOOLEAN NOT NULL DEFAULT TRUE,
    gpu_model VARCHAR(120),
    gpu_count INT NOT NULL DEFAULT 0,
    gpu_vram_gb INT NOT NULL DEFAULT 0,
    gpu_interconnect VARCHAR(64),
    extra_specs_jsonb JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE billing_terms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id UUID NOT NULL REFERENCES cloud_providers(id) ON DELETE CASCADE,
    term_code VARCHAR(64) NOT NULL, -- ON_DEMAND_HOURLY_CAPPED, PURE_HOURLY, MONTHLY_FIXED, ANNUAL_COMMIT, RESERVED_1YR, SPOT
    name VARCHAR(120) NOT NULL,
    billing_increment_seconds INT NOT NULL DEFAULT 3600,
    standard_month_hours_cap INT, -- e.g., 672 (28 days) or 730 (mean month)
    minimum_charge_period_hours NUMERIC(8,2) NOT NULL DEFAULT 1.00,
    allows_hourly_cancellation BOOLEAN NOT NULL DEFAULT TRUE,
    commitment_months INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (provider_id, term_code)
);

CREATE TABLE price_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku_id UUID NOT NULL REFERENCES product_skus(id) ON DELETE CASCADE,
    billing_term_id UUID NOT NULL REFERENCES billing_terms(id) ON DELETE RESTRICT,
    data_source_id UUID REFERENCES provider_data_sources(id) ON DELETE SET NULL,
    native_currency CHAR(3) NOT NULL, -- Exact provider billing currency (USD, EUR)
    raw_hourly_price NUMERIC(18,8),
    raw_monthly_price NUMERIC(18,8),
    raw_annual_price NUMERIC(18,8),
    setup_fee NUMERIC(18,8) NOT NULL DEFAULT 0,
    ipv4_monthly_addon_price NUMERIC(18,8) NOT NULL DEFAULT 0,
    egress_overage_per_gb_price NUMERIC(18,8) NOT NULL DEFAULT 0,
    backup_monthly_price_or_pct NUMERIC(18,8) NOT NULL DEFAULT 0,
    backup_is_percentage BOOLEAN NOT NULL DEFAULT TRUE,
    snapshot_per_gb_month_price NUMERIC(18,8) NOT NULL DEFAULT 0,
    os_license_hourly_price NUMERIC(18,8) NOT NULL DEFAULT 0,
    tax_included BOOLEAN NOT NULL DEFAULT FALSE,
    normalized_hourly_usd NUMERIC(18,8) NOT NULL,
    normalized_monthly_730h_usd NUMERIC(18,8) NOT NULL,
    fx_rate_to_usd NUMERIC(18,8) NOT NULL DEFAULT 1.00000000,
    fx_rate_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    price_effective_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    collected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    raw_payload_sha256 CHAR(64) NOT NULL,
    verification_status verification_status_enum NOT NULL DEFAULT 'VERIFIED',
    is_current BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE price_components (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    price_snapshot_id UUID NOT NULL REFERENCES price_snapshots(id) ON DELETE CASCADE,
    component_type VARCHAR(64) NOT NULL CHECK (component_type IN ('COMPUTE_BASE', 'STORAGE_BLOCK', 'IPV4_ADDRESS', 'BANDWIDTH_EGRESS', 'BACKUP_SERVICE', 'SNAPSHOT_STORAGE', 'OS_LICENSE', 'K8S_CONTROL_PLANE', 'LOAD_BALANCER', 'SETUP_FEE')),
    unit_name VARCHAR(64) NOT NULL, -- hour, month, GB-month, GB-egress
    unit_price_native NUMERIC(18,8) NOT NULL,
    included_free_units NUMERIC(18,4) NOT NULL DEFAULT 0,
    tier_min_units NUMERIC(18,4) NOT NULL DEFAULT 0,
    tier_max_units NUMERIC(18,4)
);

CREATE TABLE promotions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id UUID NOT NULL REFERENCES cloud_providers(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    promo_code VARCHAR(80),
    credit_amount NUMERIC(18,2) NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    validity_days INT NOT NULL DEFAULT 60,
    new_customers_only BOOLEAN NOT NULL DEFAULT TRUE,
    requires_credit_card BOOLEAN NOT NULL DEFAULT TRUE,
    terms_summary TEXT NOT NULL,
    starts_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE fx_rates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    base_currency CHAR(3) NOT NULL DEFAULT 'USD',
    target_currency CHAR(3) NOT NULL,
    rate NUMERIC(18,8) NOT NULL CHECK (rate > 0),
    source_provider VARCHAR(64) NOT NULL DEFAULT 'ECB_OPENEXCHANGE',
    effective_at TIMESTAMPTZ NOT NULL,
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (base_currency, target_currency, effective_at)
);

CREATE TABLE tax_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_code CHAR(2) NOT NULL UNIQUE,
    country_name VARCHAR(120) NOT NULL,
    tax_type VARCHAR(32) NOT NULL CHECK (tax_type IN ('VAT', 'GST', 'CONSUMPTION_TAX', 'SALES_TAX', 'NONE')),
    standard_rate_pct NUMERIC(6,3) NOT NULL DEFAULT 0.000,
    b2b_reverse_charge_supported BOOLEAN NOT NULL DEFAULT TRUE,
    notes TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE price_sync_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id UUID NOT NULL REFERENCES cloud_providers(id) ON DELETE CASCADE,
    data_source_id UUID NOT NULL REFERENCES provider_data_sources(id) ON DELETE CASCADE,
    status VARCHAR(32) NOT NULL CHECK (status IN ('RUNNING', 'SUCCEEDED', 'QUARANTINED_ANOMALY', 'FAILED')),
    skus_processed INT NOT NULL DEFAULT 0,
    skus_created INT NOT NULL DEFAULT 0,
    skus_updated INT NOT NULL DEFAULT 0,
    skus_deprecated INT NOT NULL DEFAULT 0,
    anomalies_detected INT NOT NULL DEFAULT 0,
    error_message TEXT,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

CREATE INDEX idx_product_skus_filter ON product_skus(category, provider_id, region_id, availability_status);
CREATE INDEX idx_price_snapshots_current ON price_snapshots(sku_id, is_current) WHERE is_current = TRUE;
