# INDEX_DESIGN.md — PostgreSQL 16 Indexing, Partitioning & Query Performance Architecture

**Document ID:** `CS-DB-003`

---

## 1. Critical Hot-Path Indexes

| Index Name | Target Table | Indexed Columns / Expression | Index Type | Query Supported & SLA Target |
| :--- | :--- | :--- | :--- | :--- |
| `idx_product_skus_filter` | `product_skus` | `(category, provider_id, region_id, availability_status)` | B-Tree | Catalog filtering & OpenSearch fallback (`< 15ms`) |
| `idx_specs_vcpu_ram_storage` | `product_specs` | `(vcpu_count, ram_mb, storage_type, storage_gb)` | B-Tree | Hardware range queries (`vCPU >= 4 AND ram_mb >= 16384`) |
| `idx_specs_gpu` | `product_specs` | `(gpu_model, gpu_vram_gb, gpu_count) WHERE gpu_count > 0` | Partial B-Tree | Dedicated GPU Server filtering (`< 5ms`) |
| `idx_skus_trgm_name` | `product_skus` | `display_name gin_trgm_ops` | GIN Trigram | Fast substring/fuzzy SKU name search |
| `idx_price_snapshots_current` | `price_snapshots` | `(sku_id, is_current) WHERE is_current = TRUE` | Partial B-Tree | O(1) lookup of active normalized price per SKU |
| `idx_affiliate_clicks_ulid` | `affiliate_clicks` | `(click_ulid)` | Unique B-Tree | Instant conversion postback attribution join (`< 2ms`) |
| `idx_cloud_resources_org` | `cloud_resources` | `(organization_id, connection_id) WHERE deleted_at IS NULL` | Partial B-Tree | Tenant-isolated BYOC server list (`< 5ms`) |
| `idx_resource_ops_status` | `resource_operations` | `(organization_id, operation_status, created_at DESC)` | B-Tree | Real-time operation status polling & idempotency check |
| `idx_outbox_unprocessed` | `outbox_events` | `(created_at ASC) WHERE published_at IS NULL` | Partial B-Tree | Sub-50ms Transactional Outbox worker polling |

---

## 2. Migration Execution Order
1. `000001_extensions_and_types.up.sql`
2. `000002_identity_domain.up.sql`
3. `000003_providers_domain.up.sql`
4. `000004_pricing_domain.up.sql`
5. `000005_affiliate_domain.up.sql`
6. `000006_cloud_management_domain.up.sql`
7. `000007_search_and_content_domain.up.sql`
8. `000008_operations_domain.up.sql`
9. `000009_rls_and_security_policies.up.sql`
10. `000010_analytical_views_and_indexes.up.sql`
11. `000011_seed_providers_and_roles.up.sql`
12. `000012_down_migrations_reference.down.sql` (Rollback reference)
