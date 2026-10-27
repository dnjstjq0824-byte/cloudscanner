-- Migration 000007: Search, Personalization & Content Domain (Article 17 & 18)
-- Tables: search_events, comparison_lists, favorites, price_alerts, articles, provider_reviews

CREATE TABLE search_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    raw_query TEXT,
    parsed_filters_jsonb JSONB NOT NULL DEFAULT '{}'::jsonb,
    result_count INT NOT NULL DEFAULT 0,
    execution_time_ms INT NOT NULL DEFAULT 0,
    locale VARCHAR(10) NOT NULL DEFAULT 'en-US',
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE comparison_lists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    share_slug VARCHAR(64) NOT NULL UNIQUE,
    title VARCHAR(160) NOT NULL DEFAULT 'Cloud Server Comparison',
    sku_ids UUID[] NOT NULL,
    scenario_config_jsonb JSONB NOT NULL DEFAULT '{"usage_hours": 730, "egress_gb": 500, "ipv4_count": 1}'::jsonb,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE favorites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sku_id UUID NOT NULL REFERENCES product_skus(id) ON DELETE CASCADE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, sku_id)
);

CREATE TABLE price_alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sku_id UUID REFERENCES product_skus(id) ON DELETE CASCADE,
    category product_category_enum,
    target_max_monthly_usd NUMERIC(18,4) NOT NULL,
    min_vcpu INT,
    min_ram_mb INT,
    target_region_code VARCHAR(64),
    notification_channel VARCHAR(32) NOT NULL DEFAULT 'EMAIL' CHECK (notification_channel IN ('EMAIL', 'WEBHOOK')),
    webhook_url TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_triggered_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE articles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug VARCHAR(180) NOT NULL,
    locale VARCHAR(10) NOT NULL DEFAULT 'en-US',
    article_type VARCHAR(32) NOT NULL CHECK (article_type IN ('TUTORIAL', 'BENCHMARK_REPORT', 'PROVIDER_COMPARISON', 'METHODOLOGY', 'MARKET_ANALYSIS')),
    title VARCHAR(240) NOT NULL,
    summary TEXT NOT NULL,
    body_markdown TEXT NOT NULL,
    structured_data_jsonld JSONB NOT NULL DEFAULT '{}'::jsonb,
    canonical_url TEXT NOT NULL,
    author_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    has_affiliate_links BOOLEAN NOT NULL DEFAULT TRUE,
    status VARCHAR(24) NOT NULL DEFAULT 'PUBLISHED' CHECK (status IN ('DRAFT', 'REVIEW', 'PUBLISHED', 'ARCHIVED')),
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (slug, locale)
);

CREATE TABLE provider_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id UUID NOT NULL REFERENCES cloud_providers(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    verified_byoc_customer BOOLEAN NOT NULL DEFAULT FALSE,
    rating_overall INT NOT NULL CHECK (rating_overall BETWEEN 1 AND 5),
    rating_performance INT NOT NULL CHECK (rating_performance BETWEEN 1 AND 5),
    rating_pricing_transparency INT NOT NULL CHECK (rating_pricing_transparency BETWEEN 1 AND 5),
    rating_support INT NOT NULL CHECK (rating_support BETWEEN 1 AND 5),
    title VARCHAR(200) NOT NULL,
    pros_text TEXT NOT NULL,
    cons_text TEXT NOT NULL,
    moderation_status VARCHAR(24) NOT NULL DEFAULT 'PENDING' CHECK (moderation_status IN ('PENDING', 'APPROVED', 'REJECTED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (provider_id, user_id)
);
