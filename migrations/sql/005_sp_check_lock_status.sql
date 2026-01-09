-- Up Migration
-- SP to CHECK lock status WITHOUT incrementing failure counter

CREATE OR REPLACE PROCEDURE sp_check_lock_status(
    p_email VARCHAR,
    OUT p_is_locked BOOLEAN,
    OUT p_message VARCHAR
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_admin_id INT;
    v_locked_until TIMESTAMP;
BEGIN
    p_is_locked := FALSE;
    p_message := 'Account is active';

    SELECT id, locked_until 
    INTO v_admin_id, v_locked_until
    FROM admins 
    WHERE email = p_email;

    IF v_admin_id IS NULL THEN
        RETURN;
    END IF;

    IF v_locked_until IS NOT NULL AND v_locked_until > NOW() THEN
        p_is_locked := TRUE;
        p_message := 'Account is locked';
        RETURN;
    END IF;

    IF v_locked_until IS NOT NULL AND v_locked_until <= NOW() THEN
        UPDATE admins SET failed_login_attempts = 0, locked_until = NULL WHERE id = v_admin_id;
    END IF;
END;
$$;

-- Down Migration
-- DROP PROCEDURE IF EXISTS sp_check_lock_status;
