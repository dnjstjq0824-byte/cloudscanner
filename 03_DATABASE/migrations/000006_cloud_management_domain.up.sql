-- Migration 000006: BYOC Cloud Management Domain (Articles 12, 13, 18, 24)
-- Tables: cloud_connections, provider_credentials_metadata, cloud_resources, resource_operations, operation_approvals, cloud_usage_snapshots, budget_policies, resource_reconciliations
-- STRICT RULE: Never store plaintext API secrets in any column.

CREATE TABLE cloud_connections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    provider_id UUID NOT NULL REFERENCES cloud_providers(id) ON DELETE RESTRICT,
    created_by_user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    connection_name VARCHAR(120) NOT NULL,
    auth_mode VARCHAR(32) NOT NULL CHECK (auth_mode IN ('OAUTH2', 'SCOPED_API_TOKEN', 'IAM_ASSUME_ROLE')),
    access_level VARCHAR(24) NOT NULL DEFAULT 'READ_ONLY' CHECK (access_level IN ('READ_ONLY', 'READ_WRITE')),
    external_account_id VARCHAR(160),
    external_account_email VARCHAR(320),
    connection_status VARCHAR(32) NOT NULL DEFAULT 'PENDING_VERIFICATION' CHECK (connection_status IN ('PENDING_VERIFICATION', 'ACTIVE', 'EXPIRED', 'REVOKED', 'DISCONNECTED', 'FROZEN_SECURITY')),
    last_verified_at TIMESTAMPTZ,
    disconnected_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE provider_credentials_metadata (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    connection_id UUID NOT NULL REFERENCES cloud_connections(id) ON DELETE CASCADE UNIQUE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    kms_key_arn VARCHAR(255) NOT NULL,
    encrypted_dek BYTEA NOT NULL, -- 256-bit Data Encryption Key wrapped by KMS KEK
    ciphertext_blob BYTEA NOT NULL, -- AES-256-GCM encrypted credential payload
    gcm_nonce BYTEA NOT NULL, -- 12-byte GCM standard nonce
    aad_context VARCHAR(255) NOT NULL, -- Bound to org_id:connection_id:provider_id
    token_fingerprint_sha256 CHAR(64) NOT NULL, -- For deduplication without plaintext
    granted_scopes TEXT[] NOT NULL DEFAULT '{}',
    token_expires_at TIMESTAMPTZ,
    rotated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE budget_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    connection_id UUID REFERENCES cloud_connections(id) ON DELETE CASCADE,
    name VARCHAR(120) NOT NULL,
    monthly_limit_usd NUMERIC(18,2) NOT NULL CHECK (monthly_limit_usd > 0),
    max_single_server_monthly_usd NUMERIC(18,2) NOT NULL DEFAULT 100.00,
    alert_threshold_pct INT NOT NULL DEFAULT 80 CHECK (alert_threshold_pct BETWEEN 1 AND 100),
    hard_block_provisioning_on_exceed BOOLEAN NOT NULL DEFAULT TRUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE cloud_resources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    connection_id UUID NOT NULL REFERENCES cloud_connections(id) ON DELETE CASCADE,
    provider_id UUID NOT NULL REFERENCES cloud_providers(id) ON DELETE RESTRICT,
    region_id UUID REFERENCES provider_regions(id) ON DELETE SET NULL,
    sku_id UUID REFERENCES product_skus(id) ON DELETE SET NULL,
    external_resource_id VARCHAR(160) NOT NULL, -- Droplet ID, UpCloud UUID, Vultr Instance UUID, Hetzner Server ID
    resource_type VARCHAR(32) NOT NULL DEFAULT 'SERVER' CHECK (resource_type IN ('SERVER', 'VOLUME', 'LOAD_BALANCER', 'FIREWALL', 'SNAPSHOT', 'SSH_KEY')),
    hostname VARCHAR(255) NOT NULL,
    power_status VARCHAR(32) NOT NULL DEFAULT 'PROVISIONING' CHECK (power_status IN ('PROVISIONING', 'RUNNING', 'STOPPED', 'REBOOTING', 'RESIZING', 'ERROR', 'DELETED')),
    public_ipv4 INET,
    public_ipv6 INET,
    private_ipv4 INET,
    vcpu_count INT,
    ram_mb INT,
    disk_gb INT,
    os_image_slug VARCHAR(120),
    backups_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    estimated_hourly_cost_usd NUMERIC(18,8) NOT NULL DEFAULT 0,
    estimated_monthly_cost_usd NUMERIC(18,8) NOT NULL DEFAULT 0,
    tags_jsonb JSONB NOT NULL DEFAULT '{}'::jsonb,
    last_synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (connection_id, external_resource_id)
);

CREATE TABLE operation_approvals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    connection_id UUID NOT NULL REFERENCES cloud_connections(id) ON DELETE CASCADE,
    approved_by_user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    operation_type VARCHAR(64) NOT NULL CHECK (operation_type IN ('CREATE_SERVER', 'RESIZE_SERVER', 'ENABLE_BACKUP', 'DELETE_SERVER', 'REBUILD_SERVER')),
    target_resource_id UUID REFERENCES cloud_resources(id) ON DELETE SET NULL,
    target_sku_id UUID REFERENCES product_skus(id) ON DELETE SET NULL,
    estimated_hourly_cost_usd NUMERIC(18,8) NOT NULL DEFAULT 0,
    estimated_monthly_cost_usd NUMERIC(18,8) NOT NULL DEFAULT 0,
    destructive_confirmation_text VARCHAR(255),
    approval_signature_hmac CHAR(64) NOT NULL,
    step_up_mfa_verified BOOLEAN NOT NULL DEFAULT TRUE,
    expires_at TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE resource_operations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    connection_id UUID NOT NULL REFERENCES cloud_connections(id) ON DELETE CASCADE,
    resource_id UUID REFERENCES cloud_resources(id) ON DELETE SET NULL,
    approval_id UUID REFERENCES operation_approvals(id) ON DELETE RESTRICT,
    initiated_by_user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    idempotency_key VARCHAR(128) NOT NULL,
    operation_type VARCHAR(64) NOT NULL, -- CREATE_SERVER, POWER_ON, POWER_OFF, REBOOT, POWER_CYCLE, RESIZE, BACKUP_ENABLE, SNAPSHOT_CREATE, DELETE_SERVER
    operation_status operation_status_enum NOT NULL DEFAULT 'DRAFT',
    request_payload_jsonb JSONB NOT NULL DEFAULT '{}'::jsonb,
    before_state_jsonb JSONB,
    after_state_jsonb JSONB,
    external_operation_id VARCHAR(160),
    external_resource_id VARCHAR(160),
    error_code VARCHAR(64),
    error_message TEXT,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (organization_id, idempotency_key)
);

CREATE TABLE cloud_usage_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    connection_id UUID NOT NULL REFERENCES cloud_connections(id) ON DELETE CASCADE,
    resource_id UUID REFERENCES cloud_resources(id) ON DELETE SET NULL,
    accrued_Month_to_date_usd NUMERIC(18,8) NOT NULL DEFAULT 0,
    projected_month_end_usd NUMERIC(18,8) NOT NULL DEFAULT 0,
    bandwidth_used_gb NUMERIC(18,4) NOT NULL DEFAULT 0,
    cpu_utilization_avg_pct NUMERIC(5,2),
    ram_utilization_avg_pct NUMERIC(5,2),
    captured_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE resource_reconciliations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    connection_id UUID NOT NULL REFERENCES cloud_connections(id) ON DELETE CASCADE,
    operation_id UUID REFERENCES resource_operations(id) ON DELETE SET NULL,
    reconciliation_type VARCHAR(64) NOT NULL CHECK (reconciliation_type IN ('POST_TIMEOUT_CHECK', 'SCHEDULED_DRIFT_SCAN', 'MANUAL_SYNC')),
    discrepancy_found BOOLEAN NOT NULL DEFAULT FALSE,
    discrepancy_type VARCHAR(64), -- ORPHAN_REMOTE_SERVER, MISSING_REMOTE_SERVER, POWER_STATE_MISMATCH, PLAN_MISMATCH
    resolution_action VARCHAR(64), -- ADOPTED_INTO_DB, MARKED_DELETED_LOCALLY, UPDATED_STATE, FLAGGED_FOR_USER
    details_jsonb JSONB NOT NULL DEFAULT '{}'::jsonb,
    reconciled_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cloud_resources_org ON cloud_resources(organization_id, connection_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_resource_ops_status ON resource_operations(organization_id, operation_status, created_at DESC);
