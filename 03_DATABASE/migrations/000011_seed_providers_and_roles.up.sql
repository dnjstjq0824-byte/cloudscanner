-- Migration 000011: Seed System Roles, Tax Policies, FX Rates & 18 Verified Cloud Providers

INSERT INTO roles (code, name, description, permissions_jsonb, is_system_role) VALUES
('ORG_OWNER', 'Organization Owner', 'Full administrative, cloud connection, provisioning, deletion approval, and billing control', '["*"]', true),
('ORG_ADMIN', 'Organization Administrator', 'Manage cloud connections, approve server provisioning and deletion within budget limits', '["cloud:connect","cloud:provision","cloud:delete_approve","members:manage"]', true),
('CLOUD_OPERATOR', 'Cloud Infrastructure Operator', 'Provision and manage servers with valid cost approval', '["cloud:read","cloud:provision","cloud:power"]', true),
('BUDGET_APPROVER', 'Budget & Cost Approver', 'Approve server creation/resizing requests exceeding standard thresholds', '["cloud:read","budget:approve"]', true),
('ORG_VIEWER', 'Organization Read-Only Viewer', 'View server inventory, cost forecasts, and comparison lists without mutation rights', '["cloud:read"]', true),
('SUPER_ADMIN', 'Platform Super Administrator', 'Global CloudScanr platform operations and security governance', '["admin:*"]', true),
('CATALOG_ADMIN', 'Catalog & Pricing Administrator', 'Manage provider catalogs and resolve pricing sync quarantines', '["admin:catalog","admin:pricing"]', true),
('AFFILIATE_ADMIN', 'Affiliate & Finance Administrator', 'Manage affiliate campaigns, links, conversion reconciliation, and payouts', '["admin:affiliate","admin:ledger"]', true);

INSERT INTO tax_policies (country_code, country_name, tax_type, standard_rate_pct, b2b_reverse_charge_supported) VALUES
('US', 'United States', 'SALES_TAX', 0.000, true),
('KR', 'South Korea', 'VAT', 10.000, true),
('DE', 'Germany', 'VAT', 19.000, true),
('FI', 'Finland', 'VAT', 25.500, true),
('FR', 'France', 'VAT', 20.000, true),
('GB', 'United Kingdom', 'VAT', 20.000, true),
('JP', 'Japan', 'CONSUMPTION_TAX', 10.000, true),
('SG', 'Singapore', 'GST', 9.000, true);

INSERT INTO fx_rates (base_currency, target_currency, rate, source_provider, effective_at) VALUES
('USD', 'USD', 1.00000000, 'SYSTEM_IDENTITY', '2026-10-28T00:00:00Z'),
('USD', 'EUR', 0.92450000, 'ECB_OPENEXCHANGE', '2026-10-28T00:00:00Z'),
('USD', 'KRW', 1382.50000000, 'ECB_OPENEXCHANGE', '2026-10-28T00:00:00Z'),
('USD', 'JPY', 152.80000000, 'ECB_OPENEXCHANGE', '2026-10-28T00:00:00Z'),
('USD', 'GBP', 0.77100000, 'ECB_OPENEXCHANGE', '2026-10-28T00:00:00Z'),
('EUR', 'USD', 1.08166577, 'ECB_OPENEXCHANGE', '2026-10-28T00:00:00Z');

INSERT INTO cloud_providers (slug, name, hq_country_code, official_website_url, official_api_docs_url, tier_priority, affiliate_verification_status, pricing_api_status, management_api_status, reseller_program_status, supports_byoc_oauth, supports_byoc_api_token) VALUES
('digitalocean', 'DigitalOcean', 'US', 'https://www.digitalocean.com', 'https://docs.digitalocean.com/reference/api/reference/', 'TIER_1_PRIMARY', 'VERIFIED', 'VERIFIED', 'VERIFIED', 'VERIFIED', true, true),
('upcloud', 'UpCloud', 'FI', 'https://upcloud.com', 'https://upcloud.com/docs/tooling/api/', 'TIER_1_PRIMARY', 'VERIFIED', 'VERIFIED', 'VERIFIED', 'VERIFIED', false, true),
('cloudways', 'Cloudways', 'MT', 'https://www.cloudways.com', 'https://developers.cloudways.com/docs', 'TIER_1_PRIMARY', 'VERIFIED', 'VERIFIED', 'VERIFIED', 'VERIFIED', true, true),
('vultr', 'Vultr', 'US', 'https://www.vultr.com', 'https://www.vultr.com/api/', 'TIER_1_PRIMARY', 'VERIFIED', 'VERIFIED', 'VERIFIED', 'VERIFIED', false, true),
('hetzner', 'Hetzner Online', 'DE', 'https://www.hetzner.com/cloud/', 'https://docs.hetzner.cloud/reference/cloud', 'TIER_1_PRIMARY', 'VERIFIED', 'VERIFIED', 'VERIFIED', 'NOT_SUPPORTED', false, true),
('akamai_linode', 'Akamai Cloud (Linode)', 'US', 'https://www.akamai.com/cloud', 'https://techdocs.akamai.com/linode-api/reference/api', 'TIER_1_PRIMARY', 'REQUIRES_APPLICATION', 'VERIFIED', 'VERIFIED', 'VERIFIED', true, true),
('gcp', 'Google Cloud Platform', 'US', 'https://cloud.google.com', 'https://docs.cloud.google.com/billing/v1/how-tos/catalog-api', 'TIER_1_PRIMARY', 'VERIFIED', 'VERIFIED', 'VERIFIED', 'VERIFIED', true, true),
('aws', 'Amazon Web Services', 'US', 'https://aws.amazon.com', 'https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/price-changes.html', 'TIER_1_PRIMARY', 'CONTRACT_REQUIRED', 'VERIFIED', 'VERIFIED', 'VERIFIED', false, true),
('azure', 'Microsoft Azure', 'US', 'https://azure.microsoft.com', 'https://learn.microsoft.com/en-us/rest/api/cost-management/retail-prices/azure-retail-prices', 'TIER_2_SECONDARY', 'CONTRACT_REQUIRED', 'VERIFIED', 'VERIFIED', 'VERIFIED', true, true),
('contabo', 'Contabo', 'DE', 'https://contabo.com', 'https://contabo.com/en/contabo-api/', 'TIER_2_SECONDARY', 'VERIFIED', 'DOCUMENTATION_UNAVAILABLE', 'VERIFIED', 'CONTRACT_REQUIRED', true, true),
('hostinger', 'Hostinger', 'LT', 'https://www.hostinger.com', 'https://developers.hostinger.com/', 'TIER_2_SECONDARY', 'VERIFIED', 'REQUIRES_APPLICATION', 'VERIFIED', 'VERIFIED', false, true),
('ovhcloud', 'OVHcloud', 'FR', 'https://www.ovhcloud.com', 'https://api.eu.ovhcloud.com/console/', 'TIER_2_SECONDARY', 'REQUIRES_APPLICATION', 'VERIFIED', 'VERIFIED', 'VERIFIED', true, true),
('scaleway', 'Scaleway', 'FR', 'https://www.scaleway.com', 'https://www.scaleway.com/en/developers/api', 'TIER_2_SECONDARY', 'CONTRACT_REQUIRED', 'VERIFIED', 'VERIFIED', 'VERIFIED', false, true),
('kamatera', 'Kamatera', 'US', 'https://www.kamatera.com', 'https://www.kamatera.com/products/api/', 'TIER_2_SECONDARY', 'VERIFIED', 'VERIFIED', 'VERIFIED', 'VERIFIED', false, true),
('ionos', 'IONOS Cloud', 'DE', 'https://cloud.ionos.com', 'https://docs.ionos.com/cloud/', 'TIER_2_SECONDARY', 'VERIFIED', 'DOCUMENTATION_UNAVAILABLE', 'VERIFIED', 'VERIFIED', false, true),
('oracle_cloud', 'Oracle Cloud Infrastructure', 'US', 'https://www.oracle.com/cloud/', 'https://www.oracle.com/cloud/price-list/', 'TIER_2_SECONDARY', 'NOT_SUPPORTED', 'VERIFIED', 'VERIFIED', 'VERIFIED', false, true),
('cloudflare', 'Cloudflare', 'US', 'https://www.cloudflare.com', 'https://developers.cloudflare.com/api/', 'TIER_2_SECONDARY', 'CONTRACT_REQUIRED', 'DOCUMENTATION_UNAVAILABLE', 'VERIFIED', 'VERIFIED', false, true),
('runpod', 'RunPod', 'US', 'https://www.runpod.io', 'https://docs.runpod.io/sdks/graphql/manage-pods', 'TIER_2_SECONDARY', 'VERIFIED', 'VERIFIED', 'VERIFIED', 'REQUIRES_APPLICATION', false, true);
