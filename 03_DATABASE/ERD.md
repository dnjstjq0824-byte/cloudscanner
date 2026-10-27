# ERD.md — CloudScanr PostgreSQL 16 Entity-Relationship Diagram (49 Tables Across 7 Domains)

**Document ID:** `CS-DB-001`  
**Database Engine:** PostgreSQL 16 (Dedicated Instance — Zero TreeScan Dependency)

---

## 1. Domain Entity-Relationship Diagram (Mermaid)

```mermaid
erDiagram
    users ||--o{ organizations : owns
    users ||--o{ organization_members : belongs_to
    organizations ||--o{ organization_members : has
    roles ||--o{ organization_members : assigns
    users ||--o{ sessions : maintains
    users ||--o{ auth_identities : authenticates_via

    cloud_providers ||--o{ provider_regions : operates_in
    cloud_providers ||--o{ provider_products : offers
    cloud_providers ||--|| provider_capabilities : supports
    cloud_providers ||--o{ provider_data_sources : syncs_from
    cloud_providers ||--o{ price_sync_jobs : executes

    provider_products ||--o{ product_skus : instantiates
    provider_regions ||--o{ product_skus : hosts
    product_skus ||--|| product_specs : specifies
    product_skus ||--o{ price_snapshots : records
    billing_terms ||--o{ price_snapshots : governs
    price_snapshots ||--o{ price_components : itemizes
    cloud_providers ||--o{ promotions : runs

    cloud_providers ||--|| affiliate_programs : monetizes_via
    affiliate_programs ||--o{ affiliate_accounts : registers
    affiliate_programs ||--o{ affiliate_campaigns : organizes
    affiliate_campaigns ||--o{ affiliate_links : generates
    affiliate_links ||--o{ affiliate_clicks : tracks
    affiliate_clicks ||--o{ affiliate_conversions : attributes
    affiliate_conversions ||--o{ commissions : posts
    commissions ||--o{ commission_adjustments : adjusts
    affiliate_payouts ||--o{ commissions : settles

    organizations ||--o{ cloud_connections : connects_byoc
    cloud_connections ||--|| provider_credentials_metadata : encrypts_kms
    cloud_connections ||--o{ cloud_resources : manages
    cloud_connections ||--o{ budget_policies : enforces
    cloud_connections ||--o{ operation_approvals : authorizes
    operation_approvals ||--o{ resource_operations : unlocks
    cloud_resources ||--o{ resource_operations : undergoes
    cloud_connections ||--o{ cloud_usage_snapshots : measures
    cloud_connections ||--o{ resource_reconciliations : reconciles

    users ||--o{ favorites : saves
    users ||--o{ comparison_lists : curates
    users ||--o{ price_alerts : subscribes
    users ||--o{ provider_reviews : writes
    users ||--o{ articles : authors
```

---

## 2. Domain Table Count & Compliance Summary (Article 18)

| Domain | Required Tables (Article 18) | Implemented SQL Migration | Total Tables |
| :--- | :--- | :--- | :--- |
| **1. Identity** | `users`, `organizations`, `organization_members`, `roles`, `sessions`, `auth_identities` | `000002_identity_domain.up.sql` | **6** |
| **2. Providers** | `cloud_providers`, `provider_regions`, `provider_products`, `provider_capabilities`, `provider_data_sources` | `000003_providers_domain.up.sql` | **5** |
| **3. Pricing** | `product_skus`, `product_specs`, `price_snapshots`, `price_components`, `billing_terms`, `promotions`, `fx_rates`, `tax_policies`, `price_sync_jobs` | `000004_pricing_domain.up.sql` | **9** |
| **4. Affiliate** | `affiliate_programs`, `affiliate_accounts`, `affiliate_campaigns`, `affiliate_links`, `affiliate_clicks`, `affiliate_conversions`, `commissions`, `commission_adjustments`, `affiliate_payouts` | `000005_affiliate_domain.up.sql` | **9** |
| **5. Cloud Management** | `cloud_connections`, `provider_credentials_metadata`, `cloud_resources`, `resource_operations`, `operation_approvals`, `cloud_usage_snapshots`, `budget_policies`, `resource_reconciliations` | `000006_cloud_management_domain.up.sql` | **8** |
| **6. Search & Content** | `search_events`, `comparison_lists`, `favorites`, `price_alerts`, `articles`, `provider_reviews` | `000007_search_and_content_domain.up.sql` | **6** |
| **7. Operations** | `audit_logs`, `background_jobs`, `outbox_events`, `admin_actions`, `feature_flags`, `system_settings` | `000008_operations_domain.up.sql` | **6** |
| **Total** | **All 49 Mandatory Domain Tables + 1 Materialized View (`mv_daily_affiliate_yield`)** | `000001` ~ `000012` | **49 Tables + 1 MV** |
