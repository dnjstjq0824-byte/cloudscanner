// Canonical API contracts mirroring the Go JSON payloads in
// internal/domain/{catalog,search,estimation,identity}.
//
// Wire truth: every decimal money/score field arrives as a JSON STRING
// (shopspring decimal default). Such fields are typed `number | string` and
// templates must render them via num()/money() from ~/composables/useNum —
// never .toFixed() on the raw payload. Plain ints (vcpu, ram_mb, counts)
// arrive as JSON numbers.

export interface Region {
  region_code: string
  name: string
  country_code: string
  city: string
  continent: string
  is_available: boolean
}

export interface ProviderCapabilities {
  supports_provisioning: boolean
  supports_power_actions: boolean
  supports_resize: boolean
  supports_snapshots: boolean
  supports_firewalls: boolean
  supports_tag_filtering: boolean
}

export interface Provider {
  id: string
  slug: string
  name: string
  hq_country_code: string
  official_website_url: string
  official_api_docs_url: string
  tier_priority: string
  affiliate_status: string
  pricing_api_status: string
  management_api_status: string
  reseller_program_status: string
  supports_byoc_oauth: boolean
  supports_byoc_api_token: boolean
  regions?: Region[]
  capabilities: ProviderCapabilities
}

export interface Category {
  code: string
  name: string
  description: string
  sku_count: number
}

export interface SKUDetail {
  id: string
  provider_slug: string
  provider_name: string
  plan_code: string
  name: string
  category: string
  cpu_class: string
  vcpu: number
  ram_mb: number
  disk_gb: number
  disk_type: string
  included_bandwidth_gb: number
  ipv4_included_count: number
  has_monthly_cap: boolean
  cap_hours: number
  raw_hourly_price_usd: number | string
  raw_monthly_price_usd: number | string
  ipv4_monthly_addon_price: number | string
  egress_overage_per_gb: number | string
  backup_is_percentage: boolean
  backup_price_or_pct: number | string
  snapshot_per_gb_month: number | string
  freshness_status: string
  regions: string[]
  last_harvested_at: string
}

export interface TCOResult {
  base_compute_usd: number | string
  ipv4_cost_usd: number | string
  egress_cost_usd: number | string
  backup_cost_usd: number | string
  snapshot_cost_usd: number | string
  subtotal_usd: number | string
  subtotal_target: number | string
  tax_target: number | string
  total_target: number | string
  effective_hourly: number | string
  currency: string
}

export interface SearchResultItem extends SKUDetail {
  tco_result: TCOResult
  best_value_score: number | string
}

export interface FacetCounts {
  categories: Record<string, number>
  providers: Record<string, number>
  cpu_classes: Record<string, number>
  disk_types: Record<string, number>
  regions: Record<string, number>
}

export interface ParsedQuery {
  detected_vcpu?: number
  detected_ram_gb?: number
  detected_disk_type?: string
  detected_region?: string
  detected_max_price?: number
  detected_category?: string
}

export interface SearchResponse {
  data: SearchResultItem[]
  facets: FacetCounts
  parsed_query?: ParsedQuery
  total_count: number
  page: number
  limit: number
  execution_ms: number
}

export interface ListEnvelope<T> {
  data: T[]
  meta: { count: number }
}

export interface ItemEnvelope<T> {
  data: T
}

export interface AuthUser {
  id: string
  email: string
  full_name: string
  status: string
  created_at: string
  updated_at: string
}

export interface AuthOrg {
  id: string
  name: string
  slug: string
  owner_user_id: string
  created_at: string
}

export interface LoginResponse {
  user: AuthUser
  token: string
  expires_at: string
  org_id: string
  role: string
}

export interface RegisterResponse {
  user: AuthUser
  organization: AuthOrg
  token: string
  expires_at: string
}

export interface MeResponse {
  user: AuthUser
  organization: AuthOrg
  role: string
  expires_at: string
}

export interface ProblemDetails {
  type: string
  title: string
  status: number
  detail?: string
  instance?: string
  trace_id?: string
}

export interface CloudConnection {
  id: string
  organization_id: string
  provider_id: string
  provider_slug: string
  created_by_user_id: string
  connection_name: string
  auth_mode: string
  access_level: string
  external_account_id: string
  external_account_email: string
  connection_status: string
  created_at: string
  updated_at: string
}

export interface BudgetPolicy {
  id: string
  organization_id: string
  connection_id?: string
  name: string
  monthly_limit_usd: number | string
  max_single_server_monthly_usd: number | string
  alert_threshold_pct: number
  hard_block_provisioning_on_exceed: boolean
  is_active: boolean
  created_at: string
}

export interface OperationApproval {
  id: string
  organization_id: string
  connection_id: string
  approved_by_user_id: string
  operation_type: string
  target_resource_id?: string
  target_sku_id?: string
  estimated_hourly_cost_usd: number | string
  estimated_monthly_cost_usd: number | string
  approval_signature_hmac: string
  step_up_mfa_verified: boolean
  expires_at: string
  created_at: string
}

export interface CloudResource {
  id: string
  organization_id: string
  connection_id: string
  provider_slug: string
  external_resource_id: string
  resource_type: string
  hostname: string
  power_status: string
  public_ipv4?: string
  vcpu_count: number
  ram_mb: number
  disk_gb: number
  estimated_hourly_cost_usd: number | string
  estimated_monthly_cost_usd: number | string
  last_synced_at: string
  created_at: string
}

export interface EstimateScenario {
  usage_hours: number
  required_ipv4_count: number
  estimated_egress_gb: number
  enable_backups: boolean
  snapshot_storage_gb: number
  target_currency: string
  fx_rate_to_target: number
  tax_percentage: number
}

export interface EstimateMeta {
  sku_id: string
  provider_slug: string
  plan_code: string
}

export interface ProvisionSync {
  status: string
  detail: string
}

export interface ProvisionEnvelope {
  data: CloudResource
  provider_sync: ProvisionSync
}

export interface AdminOverview {
  users_total: number
  providers_total: number
  pricing_verified: number
  skus_indexed: number
  connections_total: number
  connections_active: number
  servers_managed: number
  approvals_signed: number
  identity_store: string
  maintenance_mode: boolean
  management_store: string
  provider_sync_coverage: string
}

export interface AdminUser {
  id: string
  email: string
  full_name: string
  status: string
  organization_id: string
  organization_name: string
  role: string
  created_at: string
}

export interface AuditEntry {
  id: string
  created_at: string
  actor_user_id: string
  org_id: string
  action: string
  target_type: string
  target_id: string
  detail: string
}

export interface AdminNotice {
  id: string
  title: string
  body: string
  severity: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface SiteSettings {
  maintenance_mode: boolean
  announcement: string
  updated_at: string
}
