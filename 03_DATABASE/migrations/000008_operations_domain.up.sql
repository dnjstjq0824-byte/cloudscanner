-- Migration 000008: Operations, Audit, Queue & System Domain (Article 18)
-- Tables: audit_logs, background_jobs, outbox_events, admin_actions, feature_flags, system_settings

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id) ON DELETE SET NULL,
    actor_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    actor_type VARCHAR(32) NOT NULL CHECK (actor_type IN ('USER', 'ADMIN', 'SYSTEM_WORKER', 'WEBHOOK')),
    event_category VARCHAR(64) NOT NULL, -- AUTH, BYOC_CONNECTION, SERVER_LIFECYCLE, AFFILIATE_LEDGER, ADMIN_CONFIG, SECURITY_ALERT
    event_action VARCHAR(100) NOT NULL,
    target_entity_type VARCHAR(64) NOT NULL,
    target_entity_id VARCHAR(160) NOT NULL,
    request_id VARCHAR(64),
    ip_hash CHAR(64),
    metadata_jsonb JSONB NOT NULL DEFAULT '{}'::jsonb,
    prev_hash CHAR(64) NOT NULL,
    entry_hash CHAR(64) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE background_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    queue_name VARCHAR(64) NOT NULL DEFAULT 'default',
    job_type VARCHAR(100) NOT NULL,
    idempotency_key VARCHAR(160) NOT NULL UNIQUE,
    payload_jsonb JSONB NOT NULL DEFAULT '{}'::jsonb,
    job_status VARCHAR(32) NOT NULL DEFAULT 'QUEUED' CHECK (job_status IN ('QUEUED', 'PROCESSING', 'SUCCEEDED', 'RETRY_SCHEDULED', 'DEAD_LETTER')),
    attempt_count INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 5,
    last_error TEXT,
    scheduled_for TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE outbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_type VARCHAR(64) NOT NULL,
    aggregate_id VARCHAR(160) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload_jsonb JSONB NOT NULL,
    published_at TIMESTAMPTZ,
    retry_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE admin_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    second_approver_user_id UUID REFERENCES users(id) ON DELETE RESTRICT,
    action_code VARCHAR(100) NOT NULL,
    target_table VARCHAR(64) NOT NULL,
    target_record_id VARCHAR(160) NOT NULL,
    before_state_jsonb JSONB,
    after_state_jsonb JSONB,
    justification_reason TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE feature_flags (
    key VARCHAR(100) PRIMARY KEY,
    description TEXT NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    rollout_percentage INT NOT NULL DEFAULT 0 CHECK (rollout_percentage BETWEEN 0 AND 100),
    allowed_org_ids UUID[] NOT NULL DEFAULT '{}',
    updated_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE system_settings (
    key VARCHAR(100) PRIMARY KEY,
    value_jsonb JSONB NOT NULL,
    description TEXT NOT NULL,
    is_public_readable BOOLEAN NOT NULL DEFAULT FALSE,
    updated_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_outbox_unprocessed ON outbox_events(created_at ASC) WHERE published_at IS NULL;
CREATE INDEX idx_audit_logs_org_time ON audit_logs(organization_id, created_at DESC);
