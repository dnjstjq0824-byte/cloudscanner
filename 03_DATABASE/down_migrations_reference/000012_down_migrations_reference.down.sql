-- Migration 000012: Deterministic Rollback / Teardown Reference Script
DROP MATERIALIZED VIEW IF EXISTS mv_daily_affiliate_yield;
DROP TABLE IF EXISTS system_settings, feature_flags, admin_actions, outbox_events, background_jobs, audit_logs CASCADE;
DROP TABLE IF EXISTS provider_reviews, articles, price_alerts, favorites, comparison_lists, search_events CASCADE;
DROP TABLE IF EXISTS resource_reconciliations, cloud_usage_snapshots, resource_operations, operation_approvals, cloud_resources, budget_policies, provider_credentials_metadata, cloud_connections CASCADE;
DROP TABLE IF EXISTS commission_adjustments, commissions, affiliate_payouts, affiliate_conversions, affiliate_clicks, affiliate_links, affiliate_campaigns, affiliate_accounts, affiliate_programs CASCADE;
DROP TABLE IF EXISTS price_sync_jobs, tax_policies, fx_rates, promotions, price_components, price_snapshots, billing_terms, product_specs, product_skus CASCADE;
DROP TABLE IF EXISTS provider_data_sources, provider_capabilities, provider_products, provider_regions, cloud_providers CASCADE;
DROP TABLE IF EXISTS auth_identities, sessions, organization_members, roles, organizations, users CASCADE;
DROP TYPE IF EXISTS operation_status_enum, reward_type_enum, freshness_status_enum, product_category_enum, verification_status_enum;
