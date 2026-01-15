-- Up Migration: SP Handle Login Failure

CREATE OR REPLACE PROCEDURE sp_handle_login_failure(
    p_email VARCHAR,
    OUT p_is_locked BOOLEAN,
    OUT p_message VARCHAR
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_admin_id INT;
    v_failed_attempts INT;
    v_max_attempts INT := 5;
    v_lockout_duration INTERVAL := '30 minutes';
    v_locked_until TIMESTAMP;
BEGIN
    p_is_locked := FALSE;
    p_message := 'Check complete';

    SELECT id, failed_login_attempts, locked_until 
    INTO v_admin_id, v_failed_attempts, v_locked_until
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
        v_failed_attempts := 0;
    END IF;

    v_failed_attempts := v_failed_attempts + 1;

    IF v_failed_attempts >= v_max_attempts THEN
        UPDATE admins 
        SET failed_login_attempts = v_failed_attempts,
            locked_until = NOW() + v_lockout_duration
        WHERE id = v_admin_id;
        
        p_is_locked := TRUE;
        p_message := 'Account locked due to excessive failed attempts';
    ELSE
        UPDATE admins 
        SET failed_login_attempts = v_failed_attempts
        WHERE id = v_admin_id;
    END IF;
END;
$$;
