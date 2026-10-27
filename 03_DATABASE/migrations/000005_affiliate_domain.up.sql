-- Migration 000005: Affiliate & Commission Ledger Domain (Article 10 & 18)
-- Tables: affiliate_programs, affiliate_accounts, affiliate_campaigns, affiliate_links, affiliate_clicks, affiliate_conversions, commissions, commission_adjustments, affiliate_payouts

CREATE TABLE affiliate_programs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id UUID NOT NULL REFERENCES cloud_providers(id) ON DELETE CASCADE UNIQUE,
    network_name VARCHAR(80) NOT NULL, -- Impact, CJ, PartnerStack, HasOffers, Direct
    official_signup_url TEXT NOT NULL,
    individual_allowed BOOLEAN NOT NULL DEFAULT TRUE,
    kr_resident_eligible verification_status_enum NOT NULL DEFAULT 'VERIFIED',
    corporation_required BOOLEAN NOT NULL DEFAULT FALSE,
    reward_type reward_type_enum NOT NULL DEFAULT 'CASH_COMMISSION',
    commission_model VARCHAR(64) NOT NULL, -- CPA, REVSHARE, HYBRID, SLAB, CREDIT_ONLY
    default_cpa_amount NUMERIC(18,4) DEFAULT 0,
    default_revshare_pct NUMERIC(6,3) DEFAULT 0,
    cookie_duration_days INT NOT NULL DEFAULT 30,
    payout_min_threshold NUMERIC(18,2) NOT NULL DEFAULT 100.00,
    payout_currency CHAR(3) NOT NULL DEFAULT 'USD',
    tax_form_required VARCHAR(64) DEFAULT 'W-8BEN',
    trademark_bidding_prohibited BOOLEAN NOT NULL DEFAULT TRUE,
    self_referral_prohibited BOOLEAN NOT NULL DEFAULT TRUE,
    supports_subid BOOLEAN NOT NULL DEFAULT TRUE,
    supports_postback_webhook BOOLEAN NOT NULL DEFAULT FALSE,
    supports_reporting_api BOOLEAN NOT NULL DEFAULT FALSE,
    verification_status verification_status_enum NOT NULL DEFAULT 'VERIFIED',
    verified_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE affiliate_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id UUID NOT NULL REFERENCES affiliate_programs(id) ON DELETE CASCADE,
    external_publisher_id VARCHAR(120) NOT NULL,
    account_status VARCHAR(32) NOT NULL DEFAULT 'PENDING_APPROVAL' CHECK (account_status IN ('DRAFT', 'APPLIED', 'PENDING_APPROVAL', 'APPROVED_ACTIVE', 'REJECTED', 'SUSPENDED', 'EXPIRED')),
    tax_form_submitted BOOLEAN NOT NULL DEFAULT FALSE,
    payout_method_configured BOOLEAN NOT NULL DEFAULT FALSE,
    kms_api_credential_id UUID, -- Optional pointer if Reporting API key is used
    approved_at TIMESTAMPTZ,
    contract_expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE affiliate_campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id UUID NOT NULL REFERENCES affiliate_programs(id) ON DELETE CASCADE,
    account_id UUID REFERENCES affiliate_accounts(id) ON DELETE SET NULL,
    campaign_code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    target_category product_category_enum,
    promo_id UUID REFERENCES promotions(id) ON DELETE SET NULL,
    starts_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ends_at TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE affiliate_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id UUID NOT NULL REFERENCES affiliate_campaigns(id) ON DELETE CASCADE,
    provider_id UUID NOT NULL REFERENCES cloud_providers(id) ON DELETE CASCADE,
    sku_id UUID REFERENCES product_skus(id) ON DELETE SET NULL,
    slug VARCHAR(120) NOT NULL UNIQUE, -- Used in /go/{slug}
    destination_template_url TEXT NOT NULL, -- e.g. https://digitalocean.pxf.io/c/123/456/789?subId1={{click_id}}&subId2={{sig}}
    fallback_clean_url TEXT NOT NULL,
    subid_param_name VARCHAR(64) NOT NULL DEFAULT 'subId1',
    http_health_status INT NOT NULL DEFAULT 200,
    last_health_checked_at TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE affiliate_clicks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    click_ulid VARCHAR(32) NOT NULL UNIQUE, -- External ClickID passed in SubID
    link_id UUID NOT NULL REFERENCES affiliate_links(id) ON DELETE RESTRICT,
    provider_id UUID NOT NULL REFERENCES cloud_providers(id) ON DELETE RESTRICT,
    sku_id UUID REFERENCES product_skus(id) ON DELETE SET NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    subid_signature CHAR(64) NOT NULL,
    source_page_path TEXT NOT NULL,
    placement_slot VARCHAR(64) NOT NULL, -- compare_table, product_detail_hero, calculator_cta
    country_code CHAR(2),
    device_type VARCHAR(24) CHECK (device_type IN ('DESKTOP', 'MOBILE', 'TABLET', 'BOT_BLOCKED')),
    ip_subnet_hmac CHAR(64) NOT NULL,
    is_bot_filtered BOOLEAN NOT NULL DEFAULT FALSE,
    redirect_status_code INT NOT NULL DEFAULT 302,
    clicked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE affiliate_conversions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id UUID NOT NULL REFERENCES affiliate_programs(id) ON DELETE RESTRICT,
    click_id UUID REFERENCES affiliate_clicks(id) ON DELETE RESTRICT,
    external_conversion_id VARCHAR(160) NOT NULL,
    external_subid VARCHAR(120),
    conversion_source VARCHAR(32) NOT NULL CHECK (conversion_source IN ('POSTBACK_WEBHOOK', 'REPORTING_API', 'CSV_RECONCILIATION')),
    conversion_status VARCHAR(32) NOT NULL CHECK (conversion_status IN ('PENDING', 'APPROVED', 'REJECTED', 'REVERSED')),
    order_amount NUMERIC(18,8) NOT NULL DEFAULT 0,
    order_currency CHAR(3) NOT NULL DEFAULT 'USD',
    raw_event_jsonb JSONB NOT NULL DEFAULT '{}'::jsonb,
    converted_at TIMESTAMPTZ NOT NULL,
    reconciled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (program_id, external_conversion_id)
);

CREATE TABLE affiliate_payouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id UUID NOT NULL REFERENCES affiliate_programs(id) ON DELETE RESTRICT,
    external_payout_ref VARCHAR(160) NOT NULL UNIQUE,
    payout_status VARCHAR(32) NOT NULL CHECK (payout_status IN ('PENDING_SETTLEMENT', 'PAID', 'CANCELLED')),
    reward_type reward_type_enum NOT NULL DEFAULT 'CASH_COMMISSION',
    gross_amount NUMERIC(18,8) NOT NULL,
    withholding_tax_amount NUMERIC(18,8) NOT NULL DEFAULT 0,
    fee_amount NUMERIC(18,8) NOT NULL DEFAULT 0,
    net_amount NUMERIC(18,8) NOT NULL,
    payout_currency CHAR(3) NOT NULL,
    fx_rate_to_usd NUMERIC(18,8) NOT NULL DEFAULT 1.00000000,
    payout_method VARCHAR(64) NOT NULL,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE commissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversion_id UUID NOT NULL REFERENCES affiliate_conversions(id) ON DELETE RESTRICT,
    program_id UUID NOT NULL REFERENCES affiliate_programs(id) ON DELETE RESTRICT,
    payout_id UUID REFERENCES affiliate_payouts(id) ON DELETE RESTRICT,
    reward_type reward_type_enum NOT NULL DEFAULT 'CASH_COMMISSION',
    ledger_status VARCHAR(32) NOT NULL CHECK (ledger_status IN ('PENDING', 'APPROVED', 'SETTLED_PAID', 'VOIDED')),
    commission_amount NUMERIC(18,8) NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    normalized_amount_usd NUMERIC(18,8) NOT NULL,
    lock_expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE commission_adjustments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    commission_id UUID NOT NULL REFERENCES commissions(id) ON DELETE RESTRICT,
    adjustment_type VARCHAR(32) NOT NULL CHECK (adjustment_type IN ('CLAWBACK_REFUND', 'CHARGEBACK', 'TIER_BONUS', 'FX_ADJUSTMENT', 'MANUAL_RECONCILE')),
    delta_amount NUMERIC(18,8) NOT NULL, -- Negative for clawbacks, positive for tier bonuses
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    reason TEXT NOT NULL,
    recorded_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_affiliate_clicks_ulid ON affiliate_clicks(click_ulid);
CREATE INDEX idx_affiliate_clicks_time ON affiliate_clicks(clicked_at DESC);
CREATE INDEX idx_commissions_status ON commissions(ledger_status, reward_type);
