-- Migration 000009: Row-Level Security (RLS) Tenant Isolation & Immutability Triggers (Article 12 & 18)

ALTER TABLE cloud_connections ENABLE ROW LEVEL SECURITY;
ALTER TABLE provider_credentials_metadata ENABLE ROW LEVEL SECURITY;
ALTER TABLE cloud_resources ENABLE ROW LEVEL SECURITY;
ALTER TABLE resource_operations ENABLE ROW LEVEL SECURITY;
ALTER TABLE operation_approvals ENABLE ROW LEVEL SECURITY;
ALTER TABLE cloud_usage_snapshots ENABLE ROW LEVEL SECURITY;
ALTER TABLE budget_policies ENABLE ROW LEVEL SECURITY;
ALTER TABLE resource_reconciliations ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_cloud_connections ON cloud_connections
    USING (organization_id = NULLIF(current_setting('app.current_org_id', true), '')::uuid);

CREATE POLICY tenant_isolation_credentials_meta ON provider_credentials_metadata
    USING (organization_id = NULLIF(current_setting('app.current_org_id', true), '')::uuid);

CREATE POLICY tenant_isolation_cloud_resources ON cloud_resources
    USING (organization_id = NULLIF(current_setting('app.current_org_id', true), '')::uuid);

CREATE POLICY tenant_isolation_resource_ops ON resource_operations
    USING (organization_id = NULLIF(current_setting('app.current_org_id', true), '')::uuid);

CREATE POLICY tenant_isolation_op_approvals ON operation_approvals
    USING (organization_id = NULLIF(current_setting('app.current_org_id', true), '')::uuid);

CREATE POLICY tenant_isolation_usage_snapshots ON cloud_usage_snapshots
    USING (organization_id = NULLIF(current_setting('app.current_org_id', true), '')::uuid);

CREATE POLICY tenant_isolation_budget_policies ON budget_policies
    USING (organization_id = NULLIF(current_setting('app.current_org_id', true), '')::uuid);

CREATE POLICY tenant_isolation_reconciliations ON resource_reconciliations
    USING (organization_id = NULLIF(current_setting('app.current_org_id', true), '')::uuid);

-- Immutability Trigger for Audit Logs and Commission Adjustments
CREATE OR REPLACE FUNCTION prevent_immutable_modification()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'Table % is append-only immutable. Updates and deletes are strictly prohibited.', TG_TABLE_NAME;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_audit_logs_immutable
BEFORE UPDATE OR DELETE ON audit_logs
FOR EACH ROW EXECUTE FUNCTION prevent_immutable_modification();

CREATE TRIGGER trg_commission_adjustments_immutable
BEFORE UPDATE OR DELETE ON commission_adjustments
FOR EACH ROW EXECUTE FUNCTION prevent_immutable_modification();
