# TABLE_DICTIONARY.md — Complete 49-Table Data Dictionary & Retention Policies

**Document ID:** `CS-DB-002`

---

## 1. Comprehensive Table Dictionary

| # | Domain | Table Name | Primary Key | Foreign Keys & Constraints | RLS Enabled | Retention / Deletion Policy |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | Identity | `users` | `id (UUID)` | `UNIQUE(email)` | No (Self) | Soft-delete (`PENDING_DELETION` -> 30d purge PII) |
| 2 | Identity | `organizations` | `id (UUID)` | `owner_user_id -> users(id) RESTRICT` | Yes | Restricted while active `cloud_resources` exist |
| 3 | Identity | `roles` | `id (UUID)` | `UNIQUE(code)` | No | System reference table |
| 4 | Identity | `organization_members` | `id (UUID)` | `organization_id`, `user_id`, `role_id` | Yes | Cascade on org/user delete |
| 5 | Identity | `sessions` | `id (UUID)` | `user_id -> users(id) CASCADE` | No | Expired/revoked purged after 90 days |
| 6 | Identity | `auth_identities` | `id (UUID)` | `user_id -> users(id) CASCADE` | No | Cascade on user delete |
| 7 | Providers | `cloud_providers` | `id (UUID)` | `UNIQUE(slug)` | No | Soft deactivation (`is_active = false`) |
| 8 | Providers | `provider_regions` | `id (UUID)` | `provider_id -> cloud_providers(id)` | No | Preserved for historical SKU references |
| 9 | Providers | `provider_products` | `id (UUID)` | `provider_id -> cloud_providers(id)` | No | Soft deactivation (`is_active = false`) |
| 10 | Providers | `provider_capabilities` | `id (UUID)` | `provider_id UNIQUE` | No | 1:1 capability matrix per provider |
| 11 | Providers | `provider_data_sources` | `id (UUID)` | `provider_id -> cloud_providers(id)` | No | Preserved for pricing provenance |
| 12 | Pricing | `product_skus` | `id (UUID)` | `provider_id`, `product_id`, `region_id` | No | Transitions to `DEPRECATED`; never hard-deleted |
| 13 | Pricing | `product_specs` | `id (UUID)` | `sku_id -> product_skus(id) UNIQUE` | No | 1:1 hardware specification record |
| 14 | Pricing | `billing_terms` | `id (UUID)` | `provider_id -> cloud_providers(id)` | No | Reference billing rules |
| 15 | Pricing | `price_snapshots` | `id (UUID)` | `sku_id`, `billing_term_id`, `data_source_id` | No | Append-only historical price time-series (5yr retention) |
| 16 | Pricing | `price_components` | `id (UUID)` | `price_snapshot_id -> price_snapshots(id)` | No | Itemized TCO rate components |
| 17 | Pricing | `promotions` | `id (UUID)` | `provider_id -> cloud_providers(id)` | No | Preserved with `expires_at` |
| 18 | Pricing | `fx_rates` | `id (UUID)` | `UNIQUE(base_currency, target_currency, effective_at)` | No | Historical FX time-series (7yr accounting retention) |
| 19 | Pricing | `tax_policies` | `id (UUID)` | `UNIQUE(country_code)` | No | Reference VAT/GST table |
| 20 | Pricing | `price_sync_jobs` | `id (UUID)` | `provider_id`, `data_source_id` | No | Purged after 180 days |
| 21 | Affiliate | `affiliate_programs` | `id (UUID)` | `provider_id UNIQUE` | No | Preserved for contract history |
| 22 | Affiliate | `affiliate_accounts` | `id (UUID)` | `program_id -> affiliate_programs(id)` | No | Restricted while unpaid commissions exist |
| 23 | Affiliate | `affiliate_campaigns` | `id (UUID)` | `program_id`, `account_id`, `promo_id` | No | Preserved for attribution history |
| 24 | Affiliate | `affiliate_links` | `id (UUID)` | `campaign_id`, `provider_id`, `sku_id`, `UNIQUE(slug)` | No | Soft deactivate (`is_active = false` -> uses `fallback_clean_url`) |
| 25 | Affiliate | `affiliate_clicks` | `id (UUID)` | `link_id`, `provider_id`, `UNIQUE(click_ulid)` | No | Partitioned monthly; 3yr attribution retention |
| 26 | Affiliate | `affiliate_conversions` | `id (UUID)` | `program_id`, `click_id`, `UNIQUE(program_id, external_conversion_id)` | No | 7yr financial tax audit retention |
| 27 | Affiliate | `commissions` | `id (UUID)` | `conversion_id`, `program_id`, `payout_id` | No | Immutable financial ledger (7yr retention) |
| 28 | Affiliate | `commission_adjustments` | `id (UUID)` | `commission_id -> commissions(id) RESTRICT` | No | Strictly append-only (`trg_commission_adjustments_immutable`) |
| 29 | Affiliate | `affiliate_payouts` | `id (UUID)` | `program_id`, `UNIQUE(external_payout_ref)` | No | 7yr financial tax audit retention |
| 30 | Cloud Mgmt | `cloud_connections` | `id (UUID)` | `organization_id`, `provider_id`, `created_by_user_id` | **Yes** | Disconnect wipes `provider_credentials_metadata` ciphertext immediately |
| 31 | Cloud Mgmt | `provider_credentials_metadata` | `id (UUID)` | `connection_id UNIQUE`, `organization_id` | **Yes** | **Zero plaintext secret columns**; cryptographic shred on disconnect |
| 32 | Cloud Mgmt | `budget_policies` | `id (UUID)` | `organization_id`, `connection_id` | **Yes** | Cascade on org/connection delete |
| 33 | Cloud Mgmt | `cloud_resources` | `id (UUID)` | `organization_id`, `connection_id`, `provider_id` | **Yes** | Soft delete (`deleted_at`) for billing/audit history |
| 34 | Cloud Mgmt | `operation_approvals` | `id (UUID)` | `organization_id`, `connection_id`, `approved_by_user_id` | **Yes** | 7yr operational & billing approval audit retention |
| 35 | Cloud Mgmt | `resource_operations` | `id (UUID)` | `organization_id`, `connection_id`, `approval_id`, `UNIQUE(organization_id, idempotency_key)` | **Yes** | 7yr operational audit retention |
| 36 | Cloud Mgmt | `cloud_usage_snapshots` | `id (UUID)` | `organization_id`, `connection_id`, `resource_id` | **Yes** | 2yr cost analytics retention |
| 37 | Cloud Mgmt | `resource_reconciliations` | `id (UUID)` | `organization_id`, `connection_id`, `operation_id` | **Yes** | 2yr reconciliation audit retention |
| 38 | Search/Content | `search_events` | `id (UUID)` | `user_id -> users(id) SET NULL` | No | Anonymized query telemetry (1yr retention) |
| 39 | Search/Content | `comparison_lists` | `id (UUID)` | `user_id`, `UNIQUE(share_slug)` | No | Cascade on user delete |
| 40 | Search/Content | `favorites` | `id (UUID)` | `UNIQUE(user_id, sku_id)` | No | Cascade on user delete |
| 41 | Search/Content | `price_alerts` | `id (UUID)` | `user_id`, `sku_id` | No | Cascade on user delete |
| 42 | Search/Content | `articles` | `id (UUID)` | `UNIQUE(slug, locale)` | No | Editorial CMS content |
| 43 | Search/Content | `provider_reviews` | `id (UUID)` | `UNIQUE(provider_id, user_id)` | No | Anonymized on user deletion if approved |
| 44 | Operations | `audit_logs` | `id (UUID)` | `UNIQUE(entry_hash)` | No | Strictly append-only (`trg_audit_logs_immutable`, 7yr retention) |
| 45 | Operations | `background_jobs` | `id (UUID)` | `UNIQUE(idempotency_key)` | No | Completed jobs archived after 90 days |
| 46 | Operations | `outbox_events` | `id (UUID)` | `idx_outbox_unprocessed` | No | Published events pruned after 14 days |
| 47 | Operations | `admin_actions` | `id (UUID)` | `admin_user_id`, `second_approver_user_id` | No | 7yr compliance retention |
| 48 | Operations | `feature_flags` | `key (VARCHAR)` | `updated_by_user_id` | No | System runtime config |
| 49 | Operations | `system_settings` | `key (VARCHAR)` | `updated_by_user_id` | No | System runtime config |
