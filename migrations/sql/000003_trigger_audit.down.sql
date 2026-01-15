-- Down Migration: Trigger Audit

DROP TRIGGER IF EXISTS trg_admin_update ON admins;
DROP FUNCTION IF EXISTS func_audit_admin_changes;
