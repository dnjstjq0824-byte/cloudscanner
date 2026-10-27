-- Migration 000002: Identity Domain (users, organizations, organization_members, roles, sessions, auth_identities)

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(320) NOT NULL UNIQUE,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    email_verified_at TIMESTAMPTZ,
    password_hash TEXT NOT NULL, -- Argon2id encoded hash
    display_name VARCHAR(120) NOT NULL,
    preferred_locale VARCHAR(10) NOT NULL DEFAULT 'en-US',
    preferred_currency CHAR(3) NOT NULL DEFAULT 'USD',
    mfa_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    mfa_totp_encrypted_secret BYTEA,
    account_status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE' CHECK (account_status IN ('ACTIVE', 'SUSPENDED', 'PENDING_DELETION', 'DELETED')),
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug VARCHAR(80) NOT NULL UNIQUE,
    name VARCHAR(160) NOT NULL,
    owner_user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    billing_country_code CHAR(2) NOT NULL DEFAULT 'KR',
    default_currency CHAR(3) NOT NULL DEFAULT 'USD',
    require_mfa_for_cloud_ops BOOLEAN NOT NULL DEFAULT TRUE,
    require_dual_approval_for_delete BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(64) NOT NULL UNIQUE, -- ORG_OWNER, ORG_ADMIN, CLOUD_OPERATOR, BUDGET_APPROVER, ORG_VIEWER, SUPER_ADMIN, CATALOG_ADMIN, AFFILIATE_ADMIN
    name VARCHAR(120) NOT NULL,
    description TEXT NOT NULL,
    permissions_jsonb JSONB NOT NULL DEFAULT '[]'::jsonb,
    is_system_role BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE organization_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    can_connect_cloud_account BOOLEAN NOT NULL DEFAULT FALSE,
    can_provision_server BOOLEAN NOT NULL DEFAULT FALSE,
    can_approve_server_deletion BOOLEAN NOT NULL DEFAULT FALSE,
    can_approve_budget BOOLEAN NOT NULL DEFAULT FALSE,
    invited_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (organization_id, user_id)
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    active_organization_id UUID REFERENCES organizations(id) ON DELETE SET NULL,
    refresh_token_hash CHAR(64) NOT NULL UNIQUE,
    token_family_id UUID NOT NULL,
    device_name VARCHAR(160),
    user_agent TEXT,
    ip_hash CHAR(64) NOT NULL,
    step_up_mfa_verified_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    revoke_reason VARCHAR(120),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE auth_identities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider_type VARCHAR(32) NOT NULL CHECK (provider_type IN ('EMAIL', 'GOOGLE_OIDC', 'GITHUB_OAUTH', 'WEBAUTHN_PASSKEY')),
    provider_subject_id VARCHAR(255) NOT NULL,
    credential_public_key BYTEA, -- For WebAuthn Passkeys only
    sign_count BIGINT NOT NULL DEFAULT 0,
    metadata_jsonb JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (provider_type, provider_subject_id)
);

CREATE INDEX idx_org_members_user ON organization_members(user_id);
CREATE INDEX idx_sessions_user_active ON sessions(user_id, expires_at) WHERE revoked_at IS NULL;
