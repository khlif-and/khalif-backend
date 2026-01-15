-- Up Migration: Trigger Audit

CREATE OR REPLACE FUNCTION func_audit_admin_changes()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'UPDATE') THEN
        INSERT INTO admin_audit_logs (admin_id, action_type, old_data, new_data)
        VALUES (OLD.id, 'UPDATE', row_to_json(OLD), row_to_json(NEW));
        RETURN NEW;
    ELSIF (TG_OP = 'DELETE') THEN
        INSERT INTO admin_audit_logs (admin_id, action_type, old_data, new_data)
        VALUES (OLD.id, 'DELETE', row_to_json(OLD), NULL);
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_admin_update ON admins;

CREATE TRIGGER trg_admin_update
AFTER UPDATE OR DELETE ON admins
FOR EACH ROW
EXECUTE FUNCTION func_audit_admin_changes();
