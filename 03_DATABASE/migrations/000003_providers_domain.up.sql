-- Migration 000003: Providers Domain (cloud_providers, provider_regions, provider_products, provider_capabilities, provider_data_sources)

CREATE TABLE cloud_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug VARCHAR(64) NOT NULL UNIQUE, -- digitalocean, upcloud, cloudways, vultr, hetzner, akamai_linode, gcp, aws, azure, etc.
    name VARCHAR(120) NOT NULL,
    legal_entity_name VARCHAR(200),
    hq_country_code CHAR(2) NOT NULL,
    official_website_url TEXT NOT NULL,
    official_api_docs_url TEXT,
    tier_priority VARCHAR(32) NOT NULL DEFAULT 'TIER_1_PRIMARY',
    affiliate_verification_status verification_status_enum NOT NULL DEFAULT 'NOT_VERIFIED',
    pricing_api_status verification_status_enum NOT NULL DEFAULT 'NOT_VERIFIED',
    management_api_status verification_status_enum NOT NULL DEFAULT 'NOT_VERIFIED',
    reseller_program_status verification_status_enum NOT NULL DEFAULT 'CONTRACT_REQUIRED',
    supports_byoc_oauth BOOLEAN NOT NULL DEFAULT FALSE,
    supports_byoc_api_token BOOLEAN NOT NULL DEFAULT TRUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE provider_regions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id UUID NOT NULL REFERENCES cloud_providers(id) ON DELETE CASCADE,
    region_code VARCHAR(64) NOT NULL, -- e.g., nyc1, fra1, eu-central, de-fra1
    display_name VARCHAR(120) NOT NULL,
    continent_code CHAR(2) NOT NULL CHECK (continent_code IN ('NA', 'EU', 'AS', 'SA', 'AF', 'OC')),
    country_code CHAR(2) NOT NULL,
    city_name VARCHAR(120) NOT NULL,
    latitude NUMERIC(9,6),
    longitude NUMERIC(9,6),
    availability_zones_count INT NOT NULL DEFAULT 1,
    green_energy_certified BOOLEAN NOT NULL DEFAULT FALSE,
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (provider_id, region_code)
);

CREATE TABLE provider_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id UUID NOT NULL REFERENCES cloud_providers(id) ON DELETE CASCADE,
    product_code VARCHAR(100) NOT NULL,
    name VARCHAR(160) NOT NULL,
    category product_category_enum NOT NULL,
    description TEXT,
    sla_uptime_pct NUMERIC(6,4) DEFAULT 99.9900,
    official_product_url TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (provider_id, product_code)
);

CREATE TABLE provider_capabilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id UUID NOT NULL REFERENCES cloud_providers(id) ON DELETE CASCADE UNIQUE,
    can_list_regions BOOLEAN NOT NULL DEFAULT FALSE,
    can_list_plans BOOLEAN NOT NULL DEFAULT FALSE,
    can_list_images BOOLEAN NOT NULL DEFAULT FALSE,
    can_create_server BOOLEAN NOT NULL DEFAULT FALSE,
    can_power_on_off BOOLEAN NOT NULL DEFAULT FALSE,
    can_reboot BOOLEAN NOT NULL DEFAULT FALSE,
    can_power_cycle BOOLEAN NOT NULL DEFAULT FALSE,
    can_resize_server BOOLEAN NOT NULL DEFAULT FALSE,
    can_manage_backups BOOLEAN NOT NULL DEFAULT FALSE,
    can_manage_snapshots BOOLEAN NOT NULL DEFAULT FALSE,
    can_manage_firewalls BOOLEAN NOT NULL DEFAULT FALSE,
    can_manage_ssh_keys BOOLEAN NOT NULL DEFAULT FALSE,
    can_manage_vpc_network BOOLEAN NOT NULL DEFAULT FALSE,
    can_assign_floating_ip BOOLEAN NOT NULL DEFAULT FALSE,
    can_fetch_metrics BOOLEAN NOT NULL DEFAULT FALSE,
    can_delete_server BOOLEAN NOT NULL DEFAULT FALSE,
    supports_resource_tagging BOOLEAN NOT NULL DEFAULT FALSE,
    supports_client_idempotency_token BOOLEAN NOT NULL DEFAULT FALSE,
    api_version VARCHAR(32) NOT NULL,
    rate_limit_requests_per_min INT NOT NULL DEFAULT 300,
    notes TEXT,
    verified_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE provider_data_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id UUID NOT NULL REFERENCES cloud_providers(id) ON DELETE CASCADE,
    source_type VARCHAR(32) NOT NULL CHECK (source_type IN ('PUBLIC_REST_API', 'AUTH_REST_API', 'BULK_JSON_FEED', 'GRAPHQL_API', 'CURATED_RATE_CARD')),
    endpoint_url TEXT NOT NULL,
    documentation_url TEXT NOT NULL,
    sync_interval_minutes INT NOT NULL DEFAULT 60,
    last_synced_at TIMESTAMPTZ,
    last_payload_sha256 CHAR(64),
    verification_status verification_status_enum NOT NULL DEFAULT 'VERIFIED',
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
